package saError

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	err := e1()
	fmt.Println(err)
}

func e1() error {
	fmt.Println("111")
	return Stack(e2())
}

func e2() error {
	fmt.Println("222")
	return Stack(e3())
}

func e3() error {
	fmt.Println("333")
	return Stack(errors.New("saError测试"))
}
