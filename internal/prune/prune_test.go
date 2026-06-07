package prune

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

	oidcDir := filepath.Join(dir, "src", "routes", "oidc")
	if err := os.MkdirAll(oidcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oidcDir, "google.py"), []byte("# google oidc"), 0o644); err != nil {
		t.Fatal(err)
	}

	resources := filepath.Join(dir, "src", "restful_ressources.py")
	if err := os.MkdirAll(filepath.Dir(resources), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "from routes import auth, health\nfrom routes.oidc import google\nfrom routes.user import user\n"
	if err := os.WriteFile(resources, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	feature := types.Feature{
		Name:        "oidc",
		RemovePaths: []string{"src/routes/oidc"},
		RemoveLines: []types.LineRemoval{
			{File: "src/restful_ressources.py", Pattern: "oidc"},
		},
	}

	if err := Apply(dir, feature); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(oidcDir); !os.IsNotExist(err) {
		t.Error("src/routes/oidc should be removed")
	}

	data, _ := os.ReadFile(resources)
	got := string(data)
	if strings.Contains(got, "oidc") {
		t.Errorf("expected no 'oidc' lines, got:\n%s", got)
	}
	if !strings.Contains(got, "from routes import auth, health") {
		t.Errorf("expected normal lines to remain, got:\n%s", got)
	}
}

func TestApplyRejectsEscapingPaths(t *testing.T) {
	err := Apply(t.TempDir(), types.Feature{
		Name:        "unsafe",
		RemovePaths: []string{"../outside"},
	})
	if err == nil {
		t.Fatal("expected path traversal error")
	}
	if !errors.Is(err, apperrors.ErrUnsafePath) {
		t.Fatalf("error: got %v, want ErrUnsafePath", err)
	}
}
