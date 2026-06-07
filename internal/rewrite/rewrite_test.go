package rewrite

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
	"github.com/ayoub3bidi/solar-commons/internal/types"
)

func TestApply(t *testing.T) {
	dir := t.TempDir()
	envDist := filepath.Join(dir, ".env.dist")
	original := "POSTGRES_HOST=mars_db\nPOSTGRES_USER=mars\nAPP_TITLE=\"Mars API Docs\"\n"
	if err := os.WriteFile(envDist, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	rewrites := []types.FileRewrite{{
		File: ".env.dist",
		Replacements: []types.Replacement{
			{From: "mars", To: "{{.project_name}}"},
			{From: "Mars API Docs", To: "{{.app_title}}"},
		},
	}}
	vars := map[string]string{"project_name": "my-project", "app_title": "My API"}

	if err := Apply(dir, rewrites, vars); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(envDist)
	got := string(data)
	if strings.Contains(got, "mars") {
		t.Errorf("expected no 'mars' in output, got:\n%s", got)
	}
	if !strings.Contains(got, "my-project") {
		t.Errorf("expected 'my-project' in output, got:\n%s", got)
	}
	if !strings.Contains(got, "My API") {
		t.Errorf("expected 'My API' in output, got:\n%s", got)
	}
}

func TestApplyRequiresVariableValues(t *testing.T) {
	dir := t.TempDir()
	envDist := filepath.Join(dir, ".env.dist")
	if err := os.WriteFile(envDist, []byte("APP_TITLE=Fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Apply(dir, []types.FileRewrite{{
		File: ".env.dist",
		Replacements: []types.Replacement{
			{From: "Fixture", To: "{{.app_title}}"},
		},
	}}, nil)
	if err == nil {
		t.Fatal("expected missing variable error")
	}
	if !errors.Is(err, apperrors.ErrMissingVariable) {
		t.Fatalf("error: got %v, want ErrMissingVariable", err)
	}
}

func TestApplyRejectsEscapingPaths(t *testing.T) {
	err := Apply(t.TempDir(), []types.FileRewrite{{
		File: "../outside",
		Replacements: []types.Replacement{
			{From: "old", To: "new"},
		},
	}}, nil)
	if err == nil {
		t.Fatal("expected path traversal error")
	}
	if !errors.Is(err, apperrors.ErrUnsafePath) {
		t.Fatalf("error: got %v, want ErrUnsafePath", err)
	}
}
