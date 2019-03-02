package dep

import (
	"github.com/mattmcnew/archtest/examples/testfiledeps/testonlydependency"
	"testing"
)

func TestDoIBreakYou(t *testing.T) {
	testonlydependency.OohNoBadCode()
}
