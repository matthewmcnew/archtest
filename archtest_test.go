package archtest_test

import (
	"github.com/mattmcnew/archtest"
	"testing"
)

func TestPackage_ShouldNotDependOn(t *testing.T) {

	t.Run("Fails on dependencies", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/testdata/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/testdata/dependency")

		if !mockT.errored() {
			t.Error("archtest did not fail on dependency")
		}

		expected := "Error:\ngithub.com/mattmcnew/archtest/testdata/testpackage\n	github.com/mattmcnew/archtest/testdata/dependency\n"
		if mockT.message() != expected {
			t.Errorf("expected %s got error message: %s", expected, mockT.message())
		}
	})

	t.Run("Fails on transative dependencies", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/testdata/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/testdata/transative")

		if !mockT.errored() {
			t.Error("archtest did not fail on dependency")
		}

		expected := "Error:\ngithub.com/mattmcnew/archtest/testdata/testpackage\n	github.com/mattmcnew/archtest/testdata/dependency\n		github.com/mattmcnew/archtest/testdata/transative\n"
		if mockT.message() != expected {
			t.Errorf("expected %s got error message: %s", expected, mockT.message())
		}

	})

	t.Run("Succeeds on non dependencies", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/testdata/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/testdata/nodependency")

		if mockT.errored() {
			t.Error("archtest should not fail")
		}
	})
}

type testingT struct {
	errors [][]interface{}
}

func (t *testingT) Error(args ...interface{}) {
	t.errors = append(t.errors, args)
}

func (t testingT) errored() bool {
	if len(t.errors) != 0 {
		return true
	}

	return false

}

func (t *testingT) message() interface{} {
	return t.errors[0][0]
}
