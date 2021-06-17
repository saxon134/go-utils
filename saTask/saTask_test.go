package saTask

import "testing"

func TestTask(t *testing.T) {
	Init(
		Handle{Name: "b", Spec: "*/2 * * * * *", HandleFunc: f1},
		Handle{Name: "b", Spec: "0 0 10 * * * ", HandleFunc: f1},
	)
}

func f1() {

}
