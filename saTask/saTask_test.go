package saTask

func initTask() {
	Init(
		Handle{Name: "b", Spec: "", HandleFunc: f1},
		Handle{Name: "b", Spec: "", HandleFunc: f1},
	)
}

func f1() {

}
