package saOs

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestUnzipRejectsPathTraversal(t *testing.T) {
	root := t.TempDir()
	zipPath := filepath.Join(root, "bad.zip")
	if err := writeZipEntry(zipPath, "../escape.txt", "owned"); err != nil {
		t.Fatal(err)
	}

	extractDir := filepath.Join(root, "extract")
	if err := os.Mkdir(extractDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	if err := os.Chdir(extractDir); err != nil {
		t.Fatal(err)
	}

	if err := Unzip(zipPath); err == nil {
		t.Fatal("Unzip accepted a path traversal entry")
	}
	if _, err := os.Stat(filepath.Join(root, "escape.txt")); !os.IsNotExist(err) {
		t.Fatalf("path traversal file was created, stat err: %v", err)
	}
}

func writeZipEntry(zipPath, name, body string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()
	entry, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = entry.Write([]byte(body))
	return err
}
