package manifest

import (
	"testing"

	"github.com/ayoub3bidi/solar-commons/internal/testutil"
)

func TestParse(t *testing.T) {
	m, err := Parse(testutil.FixtureTemplatePath())
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "fixture" {
		t.Errorf("ID: got %q, want %q", m.ID, "fixture")
	}
	if m.DefaultRef != "v0.1.0" {
		t.Errorf("DefaultRef: got %q, want %q", m.DefaultRef, "v0.1.0")
	}
	if len(m.Variables) != 2 || m.Variables[0].Name != "project_name" {
		t.Errorf("Variables: unexpected %+v", m.Variables)
	}
	if len(m.Rewrites) != 2 || m.Rewrites[0].File != ".env.dist" {
		t.Errorf("Rewrites: unexpected %+v", m.Rewrites)
	}
	if len(m.Rewrites[0].Replacements) != 2 {
		t.Errorf("Replacements: got %d, want 2", len(m.Rewrites[0].Replacements))
	}
	if len(m.Features) != 1 || m.Features[0].Name != "oidc" {
		t.Errorf("Features: unexpected %+v", m.Features)
	}
	f := m.Features[0]
	if len(f.RemovePaths) != 2 {
		t.Errorf("RemovePaths: got %d, want 2", len(f.RemovePaths))
	}
	if len(f.RemoveLines) != 2 || f.RemoveLines[0].Pattern != "oidc" {
		t.Errorf("RemoveLines: unexpected %+v", f.RemoveLines)
	}
}
