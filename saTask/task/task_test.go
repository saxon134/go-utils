package task

import "testing"

func TestSetTaskEnableUpdatesTaskState(t *testing.T) {
	Init()
	AddTask("job", NewTask("job", "@hourly", "", func(key string, params string) error {
		return nil
	}))

	SetTaskEnable("job", false)
	if AdminTaskList["job"].IsEnable {
		t.Fatal("SetTaskEnable did not disable the task")
	}

	SetTaskEnable("job", true)
	if !AdminTaskList["job"].IsEnable {
		t.Fatal("SetTaskEnable did not enable the task")
	}
}
