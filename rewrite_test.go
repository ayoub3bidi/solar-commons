package solarcommons

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyRewrites(t *testing.T) {
	dir := t.TempDir()
	envDist := filepath.Join(dir, ".env.dist")
	original := "POSTGRES_HOST=mars_db\nPOSTGRES_USER=mars\nAPP_TITLE=\"Mars API Docs\"\n"
	if err := os.WriteFile(envDist, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	m := Manifest{
		Rewrites: []FileRewrite{
			{
				File: ".env.dist",
				Replacements: []Replacement{
					{From: "mars", To: "{{.project_name}}"},
					{From: "Mars API Docs", To: "{{.app_title}}"},
				},
			},
		},
	}
	vars := map[string]string{"project_name": "my-project", "app_title": "My API"}

	if err := ApplyRewrites(dir, m, vars); err != nil {
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

func TestApplyRewritesRequiresVariableValues(t *testing.T) {
	dir := t.TempDir()
	envDist := filepath.Join(dir, ".env.dist")
	if err := os.WriteFile(envDist, []byte("APP_TITLE=Fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ApplyRewrites(dir, Manifest{
		Rewrites: []FileRewrite{{
			File: ".env.dist",
			Replacements: []Replacement{
				{From: "Fixture", To: "{{.app_title}}"},
			},
		}},
	}, nil)
	if err == nil {
		t.Fatal("expected missing variable error")
	}
	if !errors.Is(err, ErrMissingVariable) {
		t.Fatalf("error: got %v, want ErrMissingVariable", err)
	}
}

func TestApplyRewritesRejectsEscapingPaths(t *testing.T) {
	err := ApplyRewrites(t.TempDir(), Manifest{
		Rewrites: []FileRewrite{{
			File: "../outside",
			Replacements: []Replacement{
				{From: "old", To: "new"},
			},
		}},
	}, nil)
	if err == nil {
		t.Fatal("expected path traversal error")
	}
	if !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("error: got %v, want ErrUnsafePath", err)
	}
}
