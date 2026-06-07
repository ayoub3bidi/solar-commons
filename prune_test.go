package solarcommons

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPruneFeature(t *testing.T) {
	dir := t.TempDir()

	// Create OIDC path subtree
	oidcDir := filepath.Join(dir, "src", "routes", "oidc")
	if err := os.MkdirAll(oidcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oidcDir, "google.py"), []byte("# google oidc"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create restful_ressources.py with mixed lines
	resources := filepath.Join(dir, "src", "restful_ressources.py")
	if err := os.MkdirAll(filepath.Dir(resources), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "from routes import auth, health\nfrom routes.oidc import google\nfrom routes.user import user\n"
	if err := os.WriteFile(resources, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	f := Feature{
		Name:        "oidc",
		RemovePaths: []string{"src/routes/oidc"},
		RemoveLines: []LineRemoval{
			{File: "src/restful_ressources.py", Pattern: "oidc"},
		},
	}

	if err := PruneFeature(dir, f); err != nil {
		t.Fatal(err)
	}

	// OIDC dir must be gone
	if _, err := os.Stat(oidcDir); !os.IsNotExist(err) {
		t.Error("src/routes/oidc should be removed")
	}

	// OIDC lines must be gone, normal lines must remain
	data, _ := os.ReadFile(resources)
	got := string(data)
	if strings.Contains(got, "oidc") {
		t.Errorf("expected no 'oidc' lines, got:\n%s", got)
	}
	if !strings.Contains(got, "from routes import auth, health") {
		t.Errorf("expected normal lines to remain, got:\n%s", got)
	}
}
