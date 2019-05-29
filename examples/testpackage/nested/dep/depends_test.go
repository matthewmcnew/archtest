package dep

import (
	"github.com/matthewmcnew/archtest/examples/testfiledeps/testonlydependency"
	"testing"
)

func TestDoIBreakYou(t *testing.T) {
	testonlydependency.OohNoBadCode()
}
