package testpackage

import (
	"C"
	"fmt"
	"github.com/mattmcnew/archtest/examples/dependency"
	"runtime/debug"
)

func What(a debug.BuildInfo) {
	fmt.Println(dependency.Item)
}
