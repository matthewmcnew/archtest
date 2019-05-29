package dep

import (
	"fmt"
	"github.com/matthewmcnew/archtest/examples/nesteddependency"
)

func init() {
	fmt.Println(nesteddependency.Item)
}
