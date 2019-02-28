package dep

import (
	"fmt"
	"github.com/mattmcnew/archtest/examples/nesteddependency"
)

func init() {
	fmt.Println(nesteddependency.Item)
}
