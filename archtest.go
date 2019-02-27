package archtest

import (
	"fmt"
	"go/build"
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
		if i == d {
			p.t.Error("blah")
		}
	}
}

func findDeps(packageName string) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)

		importCache := make(map[string]struct{})
		read(c, packageName, importCache)
	}()
	return c
}

func read(packages chan string, name string, importCache map[string]struct{}) {
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
		packages <- i
	}

	for _, i := range newImports {
		read(packages, i, importCache)
	}
}
