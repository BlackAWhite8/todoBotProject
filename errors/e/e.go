package e

import "fmt"

func WrapErr(errMsg string, err error) {
	fmt.Errorf(errMsg, err)
}
