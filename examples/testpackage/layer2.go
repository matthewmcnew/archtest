package testpackage

import (
	"crypto"
	"fmt"
	"github.com/mattmcnew/archtest/examples/dependency"
)

func What(a crypto.Decrypter) {
	fmt.Println(dependency.Item)
}
