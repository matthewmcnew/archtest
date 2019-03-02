package dep_test

import (
	"github.com/mattmcnew/archtest/examples/testfiledeps/testpkgdependency"
	"testing"
)

func Test(t *testing.T) {
	testpkgdependency.OohNoBadCode()
}
