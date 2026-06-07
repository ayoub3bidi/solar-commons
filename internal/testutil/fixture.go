package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func FixtureTemplatePath() string {
	return filepath.Join(internalDir(), "fixtures", "template")
}

func CopyFixtureTemplate(t *testing.T, dest string) {
	t.Helper()

	source := FixtureTemplatePath()
	err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relativePath == "." {
			return nil
		}
		target := filepath.Join(dest, relativePath)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
	if err != nil {
		t.Fatal(err)
	}
}

func internalDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "internal"
	}
	return filepath.Dir(filepath.Dir(filename))
}
