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

func TestInitGoTypeDoesNotPanic(t *testing.T) {
	log = nil
	defer func() {
		if e := recover(); e != nil {
			t.Fatalf("Init GoType panicked: %v", e)
		}
	}()
	Init(InfoLevel, GoType)
	if log == nil {
		t.Fatal("Init GoType left log nil")
	}
}

func TestInfoWithoutInitDoesNotPanic(t *testing.T) {
	oldLog := log
	oldLevel := logLevel
	log = nil
	logLevel = InfoLevel
	defer func() {
		log = oldLog
		logLevel = oldLevel
		if e := recover(); e != nil {
			t.Fatalf("Info panicked without initialized logger: %v", e)
		}
	}()

	Info("hello")
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
