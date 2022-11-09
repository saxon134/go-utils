package saLog

import (
	"errors"
	"fmt"
	"testing"
)

func TestSaLog(t *testing.T) {
	Init(InfoLevel, ZapType)
	err := f1()
	Err(err)
}

func f1() error {
	fmt.Println("f1")
	return f2()
}

func f2() error {
	fmt.Println("f2")
	return f3()
}

func f3() error {
	return errors.New("function 1")
}
