package saTask

import "gitee.com/go-utils/saBroker/saTrigger"

func HandleDemo() (err error) {
	err = saTrigger.Remote("print_test", "123abc")
	return err
}
