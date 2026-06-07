package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
	"github.com/ayoub3bidi/solar-commons/internal/types"
	"gopkg.in/yaml.v3"
)

// Parse reads solar.manifest.yaml from dest and returns the parsed manifest.
func Parse(dest string) (types.Manifest, error) {
	data, err := os.ReadFile(filepath.Join(dest, "solar.manifest.yaml"))
	if err != nil {
		return types.Manifest{}, fmt.Errorf("%w: read manifest: %w", apperrors.ErrInvalidManifest, err)
	}
	var m types.Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return types.Manifest{}, fmt.Errorf("%w: parse manifest: %w", apperrors.ErrInvalidManifest, err)
	}
	if err := m.Validate(); err != nil {
		return types.Manifest{}, err
	}
	return m, nil
}
