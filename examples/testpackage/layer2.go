package testpackage

import (
	"crypto"
	"fmt"
	"github.com/matthewmcnew/archtest/examples/dependency"
)

func What(a crypto.Decrypter) {
	fmt.Println(dependency.Item)
}
