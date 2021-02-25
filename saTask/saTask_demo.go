package saTask

import "github.com/saxon134/go-utils/saBroker/saTrigger"

func HandleDemo() (err error) {
	err = saTrigger.Remote("print_test", "123abc")
	return err
}
