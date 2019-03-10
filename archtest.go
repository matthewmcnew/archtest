package archtest

import (
	"container/list"
	"fmt"
	"go/build"
	"golang.org/x/tools/go/packages"
	"strings"
)

type PackageTest struct {
	packages     []string
	excluding    map[string]interface{}
	t            TestingT
	includeTests bool
}

type TestingT interface {
	Error(args ...interface{})
}

func Package(t TestingT, packageName ...string) *PackageTest {
	return &PackageTest{packages: packageName, t: t, includeTests: false}
}

func (t PackageTest) IncludeTests() *PackageTest {
	t.includeTests = true
	return &t
}

func (t PackageTest) Excluding(e ...string) *PackageTest {
	set := make(map[string]interface{})

	for v := range t.excluding {
		set[v] = struct{}{}
	}

	for _, v := range t.expand(e) {
		set[v] = struct{}{}
	}

	t.excluding = set
	return &t
}

func (t PackageTest) ShouldNotDependDirectlyOn(pkgs ...string) {
	t.shouldNotDependOnPackageWithFilter(func(d *dep) bool {
		return d.depth() > 1
	}, pkgs)
}

func (t *PackageTest) ShouldNotDependOn(pkgs ...string) {
	t.shouldNotDependOnPackageWithFilter(noOpFilter, pkgs)
}

func (t *PackageTest) shouldNotDependOnPackageWithFilter(filter depFilter, d []string) {
	dl := t.expand(d)
	for i := range t.findDeps(t.packages) {
		if filter(i) {
			continue
		}

		if i.isDependencyOn(dl) {
			chain, _ := i.chain()
			msg := fmt.Sprintf("Error:\n%s", chain)
			t.t.Error(msg)
		}
	}
}

type depFilter func(*dep) bool

var noOpFilter depFilter = func(i *dep) bool {
	return false
}

type dep struct {
	name   string
	parent *dep
	xtest  bool
}

func (d *dep) depth() int {
	if d.parent == nil {
		return 0
	}
	return d.parent.depth() + 1
}

func (d *dep) chain() (string, int) {
	name := d.name
	if d.xtest {
		name = d.name + "_test"
	}

	if d.parent == nil {
		return name + "\n", 1
	}

	c, tabs := d.parent.chain()

	return c + strings.Repeat("\t", tabs) + name + "\n", tabs + 1
}

func (d dep) asxtest() *dep {
	d.xtest = true
	return &d
}

func (d *dep) isDependencyOn(dl []string) bool {
	if d.parent == nil {
		return false
	}

	if contains(dl, d.name) {
		return true
	}
	return false
}

func (t *PackageTest) findDeps(packages []string) <-chan *dep {
	c := make(chan *dep)
	go func() {
		defer close(c)

		importCache := map[string]struct{}{}
		for _, p := range t.expand(packages) {

			t.read(c, &dep{name: p, parent: nil}, importCache)
		}
	}()
	return c
}

func (t *PackageTest) read(pChan chan *dep, d *dep, cache map[string]struct{}) {
	queue := list.New()

	context := build.Default
	var importMode build.ImportMode

	queue.PushBack(d)
	for queue.Len() > 0 {
		front := queue.Front()
		queue.Remove(front)
		d, _ := (front.Value).(*dep)

		if t.skip(cache, d.name) {
			continue
		}

		cache[d.name] = struct{}{}
		pChan <- d

		pkg, err := context.Import(d.name, ".", importMode)
		if err != nil {
			e := fmt.Sprintf("Error reading: %s", d.name)
			t.t.Error(e)

			continue
		}
		if pkg.Goroot {
			continue
		}

		for _, i := range pkg.Imports {
			queue.PushBack(&dep{name: i, parent: d})
		}

		if t.includeTests {
			for _, i := range pkg.TestImports {
				queue.PushBack(&dep{name: i, parent: d})
			}

			for _, i := range pkg.XTestImports {
				queue.PushBack(&dep{name: i, parent: d.asxtest()})
			}
		}
	}
}

func (t *PackageTest) expand(ps []string) []string {
	if !needExpansion(ps) {
		return ps
	}

	cfg := &packages.Config{
		Mode:       packages.LoadFiles,
		Tests:      false,
		BuildFlags: []string{},
	}

	loadedPs, err := packages.Load(cfg, ps...)
	if err != nil {
		e := fmt.Sprintf("Error reading: %s, err: %s", ps, err)
		t.t.Error(e)
		return nil
	}
	if len(loadedPs) == 0 {
		e := fmt.Sprintf("Error reading: %s, did not match any packages", ps)
		t.t.Error(e)
		return nil

	}

	ls := make([]string, 0, len(loadedPs))

	for _, p := range loadedPs {
		ls = append(ls, p.PkgPath)
	}

	return ls
}

func (t PackageTest) skip(cache map[string]struct{}, pkg string) bool {

	if _, excluded := t.excluding[pkg]; excluded ||
		pkg == "C" {
		return true
	}

	_, seen := cache[pkg]
	return seen

}

func needExpansion(ps []string) bool {
	for _, p := range ps {
		if strings.Contains(p, "...") {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
