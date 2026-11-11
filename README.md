# solar-commons

Core Go library powering the [`sun`](https://github.com/ayoub3bidi/sun) CLI. It handles Solar template scaffolding: clone a pinned template ref, strip source git history, parse `solar.manifest.yaml`, apply variable rewrites, prune disabled features, and initialize a clean git repository.

> **Solar Stack** — see the [ecosystem map](https://github.com/ayoub3bidi/sun/blob/v1.0.0/docs/SOLAR_STACK.md) in the `sun` repo.

`solar-commons` is intentionally UI-agnostic. CLI prompts, Bubble Tea flows, template registry decisions, install scripts, and user-facing labels belong in `sun`.

## Pipeline

```
Scaffold  →  clone  →  strip .git  →  manifest  →  rewrite  →  prune  →  git init
```

## API

```go
import sc "github.com/ayoub3bidi/solar-commons"

result, err := sc.Scaffold(sc.ScaffoldConfig{
    Repo:        "https://github.com/ayoub3bidi/mars",
    Ref:         "v0.1.5",
    Destination: "./my-api",
    Variables: map[string]string{
        "project_name": "my-api",
        "app_title":    "My API",
        "api_version":  "v1",
    },
    DisabledFeatures: []string{"oidc"},
    InitialCommit:    "Initial commit",
    Author: sc.GitAuthor{
        Name:  "sun",
        Email: "sun@solar",
    },
})
if err != nil {
    // handle error
}
_ = result
```

Primitive helpers are still exposed for focused tests and advanced callers:

```go
CloneTemplate(repo, ref, dest string) error
StripGit(dest string) error
ParseManifest(dest string) (Manifest, error)
ApplyRewrites(dest string, manifest Manifest, vars map[string]string) error
PruneFeature(dest string, feature Feature) error
InitGit(dest, message string) error
```

Use `Scaffold` as the canonical API for `sun`.

## Errors

Public sentinel errors are available for `errors.Is`:

| Error | Meaning |
|-------|---------|
| `ErrInvalidConfig` | Required scaffold config is missing |
| `ErrDestinationExists` | Destination already exists |
| `ErrInvalidManifest` | Manifest is missing, malformed, or invalid |
| `ErrMissingVariable` | A rewrite references a variable not provided by the caller |
| `ErrUnknownFeature` | Caller asked to disable a feature not declared in the manifest |
| `ErrUnsafePath` | Manifest-controlled path is absolute or escapes the destination root |

## Manifest format

Each Solar template repo carries a `solar.manifest.yaml` at its root:

```yaml
id: mars
repo: https://github.com/ayoub3bidi/mars
default_ref: v0.1.5
variables:
  - name: project_name
  - name: app_title
  - name: api_version
rewrites:
  - file: .env.dist
    replacements:
      - from: "mars"
        to: "{{.project_name}}"
features:
  - name: oidc
    default: false
    remove_paths:
      - src/routes/oidc
      - src/controllers/oidc
    remove_lines:
      - file: src/restful_ressources.py
        pattern: "oidc"
```

Manifest rules:

- `id`, `repo`, and `default_ref` are required.
- `default_ref` should pin a stable semver tag for end-user scaffolds.
- Variables and features must have unique, non-empty names.
- Rewrites are literal string replacements, not regexes.
- Replacement values can reference variables with `{{.project_name}}` syntax.
- Missing variable values fail fast.
- Manifest-controlled paths must be relative and stay inside the destination root.

## Templates

| ID | Repo | Stack |
|----|------|-------|
| `mars` | [ayoub3bidi/mars](https://github.com/ayoub3bidi/mars) | FastAPI + Alembic (recommended) |
| `mercury` | [ayoub3bidi/mercury](https://github.com/ayoub3bidi/mercury) | FastAPI + Flyway |
| `venus` | [ayoub3bidi/venus](https://github.com/ayoub3bidi/venus) | Docusaurus + Tailwind |

## Testing

```bash
# Offline unit and local E2E tests
go test ./...
go vet ./...

# Networked integration tests against the pinned Mars, Mercury, and Venus releases
go test -tags integration -timeout 300s ./...
```

The default test suite must stay offline. Networked tests belong behind the `integration` build tag.

## Architecture

The root package is the only public API surface. Public structs are exposed as root aliases, while their canonical definitions and implementation details live under `internal/`:

| Path | Role |
|------|------|
| `internal/types/` | Canonical scaffold and manifest structs |
| `internal/scaffold/` | Full clone, rewrite, prune, and git-init orchestration |
| `internal/manifest/` | Manifest file loading and YAML parsing |
| `internal/git/` | go-git clone, strip, and init implementation |
| `internal/fs/` | Safe path handling for manifest-controlled paths |
| `internal/rewrite/` | Literal rewrite engine and variable resolution |
| `internal/prune/` | Disabled feature path and line pruning |
| `internal/fixtures/template/` | Local scaffold fixture used by offline E2E tests |
| `internal/testutil/` | Shared test helpers for fixture repositories |
