package saVideo

import (
	"testing"
)

func TestScreenshot(t *testing.T) {
	first, _ := Load("1.mp4")
	_, _ = first.Screenshot(1, "simple-screen.png")
}
