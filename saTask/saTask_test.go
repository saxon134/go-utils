package saTask

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"testing"
)

func TestTask(t *testing.T) {
	f1()
	return

	Init(
		Handle{Name: "b", Spec: "*/2 * * * * *", HandleFunc: f1},
		Handle{Name: "b", Spec: "0 0 10 * * * ", HandleFunc: f1},
	)
}

func f1() {
	debug.Stack()

	pcs := make([]uintptr, 10)
	runtime.CallersFrames(pcs)

	fmt.Println(pcs)
}
