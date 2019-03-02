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

func (t *PackageTest) ShouldNotDependOn(d string) {
	for i := range t.findDeps(t.packages) {
		if i.name == d {
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
		for _, p := range expand(packages) {

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
		fmt.Printf("build import error: %+v", err)
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

func skip(cache map[string]struct{}, pkg string) bool {
	if strings.HasPrefix(pkg, "internal/") || pkg == "C" {
		return true
	}

	_, seen := cache[pkg]
	return seen
}

func expand(ps []string) []string {
	cfg := &packages.Config{
		Mode:       packages.LoadImports,
		Tests:      false,
		BuildFlags: []string{},
	}

	loadedPs, err := packages.Load(cfg, ps...)
	if err != nil {
		fmt.Printf("packages error: %+v", err)
		return nil
	}

	ls := make([]string, 0, len(loadedPs))

	for _, p := range loadedPs {
		ls = append(ls, p.PkgPath)
	}

	return ls
}
