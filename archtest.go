package archtest

import (
	"fmt"
	"go/build"
	"strings"
)

type PackageTest struct {
	packageName string
	t           TestingT
}

type TestingT interface {
	Error(args ...interface{})
}

func Package(t TestingT, packageName string) *PackageTest {
	return &PackageTest{packageName, t}
}

func (p *PackageTest) ShouldNotDependOn(d string) {
	for i := range findDeps(p.packageName) {
		if i.name == d {
			chain, _ := i.chain()
			msg := fmt.Sprintf("Error:\n%s", chain)
			p.t.Error(msg)
		}
	}
}

type dep struct {
	name   string
	parent *dep
}

func (d *dep) chain() (string, int) {
	if d.parent == nil {
		return d.name + "\n", 1
	}

	c, tabs := d.parent.chain()

	return c + strings.Repeat("\t", tabs) + d.name + "\n", tabs + 1
}

func findDeps(packageName string) <-chan *dep {
	c := make(chan *dep)
	go func() {
		defer close(c)

		importCache := make(map[string]struct{})
		read(c, &dep{packageName, nil}, packageName, importCache)
	}()
	return c
}

func read(packages chan *dep, parent *dep, name string, importCache map[string]struct{}) {
	context := build.Default
	var importMode build.ImportMode

	pkg, err := context.Import(name, ".", importMode)
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}

	newImports := make([]string, 0, len(pkg.Imports))
	for _, i := range pkg.Imports {
		if _, seen := importCache[i]; seen {
			continue
		}
		newImports = append(newImports, i)
		importCache[i] = struct{}{}
	}

	for _, i := range newImports {
		packages <- &dep{i, parent}
	}

	for _, i := range newImports {
		read(packages, &dep{i, parent}, i, importCache)
	}
}
