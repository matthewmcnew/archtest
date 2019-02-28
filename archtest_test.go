package archtest_test

import (
	"github.com/mattmcnew/archtest"
	"testing"
)

func TestPackage_ShouldNotDependOn(t *testing.T) {
	t.Run("Fails on dependencies", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/dependency")

		if !mockT.errored() {
			t.Fatal("archtest did not fail on dependency")
		}

		expected := "Error:\ngithub.com/mattmcnew/archtest/examples/testpackage\n	github.com/mattmcnew/archtest/examples/dependency\n"
		if mockT.message() != expected {
			t.Errorf("expected %s got error message: %s", expected, mockT.message())
		}
	})

	t.Run("Fails on transative dependencies", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/transative")

		if !mockT.errored() {
			t.Fatal("archtest did not fail on dependency")
		}

		expected := "Error:\ngithub.com/mattmcnew/archtest/examples/testpackage\n	github.com/mattmcnew/archtest/examples/dependency\n		github.com/mattmcnew/archtest/examples/transative\n"
		if mockT.message() != expected {
			t.Errorf("expected %s got error message: %s", expected, mockT.message())
		}

	})

	t.Run("Supports multiple packages", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/dontdependonanything", "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/dependency")

		if !mockT.errored() {
			t.Fatal("archtest did not fail on dependency")
		}

		expected := "Error:\ngithub.com/mattmcnew/archtest/examples/testpackage\n	github.com/mattmcnew/archtest/examples/dependency\n"
		if mockT.message() != expected {
			t.Fatalf("expected %s got error message: %s", expected, mockT.message())
		}
	})

	t.Run("Supports wildcard matching", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/...").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/nodependency")

		if mockT.errored() {
			t.Fatalf("archtest with a wildcard should not have failed on %s", "github.com/mattmcnew/archtest/examples/nodependency")
		}

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage/...").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/nesteddependency")

		if !mockT.errored() {
			t.Fatal("archtest did not fail on dependency")
		}

		expected := "Error:\ngithub.com/mattmcnew/archtest/examples/testpackage/nested/dep\n	github.com/mattmcnew/archtest/examples/nesteddependency\n"
		if mockT.message() != expected {
			t.Errorf("expected %s got error message: %s", expected, mockT.message())
		}
	})

	t.Run("Succeeds on non dependencies", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/nodependency")

		if mockT.errored() {
			t.Fatalf("archtest should not fail")
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
