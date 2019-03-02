package archtest

import (
	"fmt"
	"go/build"
	"golang.org/x/tools/go/packages"
	"strings"
)

type PackageTest struct {
	packages     []string
	t            TestingT
	includeTests bool
}

type TestingT interface {
	Error(args ...interface{})
}

func Package(t TestingT, packageName ...string) *PackageTest {
	return &PackageTest{packageName, t, false}
}

func (t PackageTest) IncludeTests() *PackageTest {
	t.includeTests = true
	return &t
}

func (t *PackageTest) ShouldNotDependOn(d ...string) {
	dl := t.expand(d)

	for i := range t.findDeps(t.packages) {
		if contains(dl, i.name) {
			chain, _ := i.chain()
			msg := fmt.Sprintf("Error:\n%s", chain)
			t.t.Error(msg)
		}
	}
}

type dep struct {
	name   string
	parent *dep
	xtest  bool
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

func (t *PackageTest) findDeps(packages []string) <-chan *dep {
	c := make(chan *dep)
	go func() {
		defer close(c)

		importCache := map[string]struct{}{}
		for _, p := range t.expand(packages) {

			t.read(c, &dep{name: p}, importCache)
		}
	}()
	return c
}

func (t *PackageTest) read(pChan chan *dep, d *dep, cache map[string]struct{}) {
	context := build.Default
	var importMode build.ImportMode

	pkg, err := context.Import(d.name, ".", importMode)
	if err != nil {
		e := fmt.Sprintf("Error reading: %s", d.name)
		t.t.Error(e)
		return
	}

	newImports := make([]*dep, 0, len(pkg.Imports)+len(pkg.TestImports)+len(pkg.XTestImports))

	for _, i := range pkg.Imports {
		if skip(cache, i) {
			continue
		}

		dep := &dep{name: i, parent: d}
		cache[dep.name] = struct{}{}
		pChan <- dep
		newImports = append(newImports, dep)
	}

	if t.includeTests {
		for _, i := range pkg.TestImports {
			if skip(cache, i) {
				continue
			}

			dep := &dep{name: i, parent: d}
			cache[dep.name] = struct{}{}
			pChan <- dep
			newImports = append(newImports, dep)
		}

		for _, i := range pkg.XTestImports {
			if skip(cache, i) {
				continue
			}

			dep := &dep{name: i, parent: d.asxtest()}
			cache[dep.name] = struct{}{}
			pChan <- dep
			newImports = append(newImports, dep)
		}

	}

	for _, v := range newImports {
		t.read(pChan, v, cache)
	}
}

func (t *PackageTest) expand(ps []string) []string {
	if !needExpansion(ps) {
		return ps
	}

	cfg := &packages.Config{
		Mode:       packages.LoadImports,
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

func skip(cache map[string]struct{}, pkg string) bool {
	if strings.HasPrefix(pkg, "internal/") || pkg == "C" {
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
