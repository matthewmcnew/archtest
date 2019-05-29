package dependency

import "fmt"
import "github.com/matthewmcnew/archtest/examples/transative"

const Item = "depend on me"

func Somemethod() {
	fmt.Println(transative.NowYouDependOnMe)
}
