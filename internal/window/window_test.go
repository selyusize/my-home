package window

import (
	"path/filepath"
	"testing"
)

func TestFolderPathRejectsRelative(t *testing.T) {
	t.Parallel()
	if _, err := folderPath("relative"); err == nil {
		t.Fatal("expected error")
	}
}

func TestFolderPathRejectsMissing(t *testing.T) {
	t.Parallel()
	_, err := folderPath(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFolderPathAcceptsDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	got, err := folderPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != dir {
		t.Fatalf("got %q want %q", got, dir)
	}
}
