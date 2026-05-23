package saImg

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestMakeThumbnailReturnsErrorForMissingSource(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatalf("MakeThumbnail panicked for missing source: %v", e)
		}
	}()

	err := MakeThumbnail("missing-source.png", "unused-output.png", 100)
	if err == nil {
		t.Fatal("MakeThumbnail returned nil error for missing source")
	}
}

func TestMakeThumbnailReturnsErrorForInvalidOutputPath(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src.png")
	if err := writeTestPNG(src); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if e := recover(); e != nil {
			t.Fatalf("MakeThumbnail panicked for invalid output path: %v", e)
		}
	}()

	err := MakeThumbnail(src, filepath.Join(t.TempDir(), "missing", "out.png"), 100)
	if err == nil {
		t.Fatal("MakeThumbnail returned nil error for invalid output path")
	}
}

func writeTestPNG(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.White)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
