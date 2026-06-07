# Solar Commons Agent Guide

## Purpose

`solar-commons` is the Go library behind the planned `sun` CLI. It owns template filesystem operations only: clone a pinned template, remove source git history, parse and validate `solar.manifest.yaml`, apply variable rewrites, prune disabled features, and initialize a clean git repository.

Keep this module UI-agnostic. Bubble Tea prompts, Cobra commands, template selection, and user-facing labels belong in `sun`, not here.

## Current Shape

| File | Role |
|------|------|
| `clone.go`, `git.go`, `rewrite.go`, `prune.go` | Root package wrappers for the public API |
| `manifest.go` | Manifest structs, YAML parsing, and validation |
| `internal/git/` | go-git clone, strip, and init implementation |
| `internal/files/` | Safe path handling, rewrites, and pruning implementation |
| `internal/fixtures/template/` | Local scaffold fixture for offline E2E tests |
| `.github/workflows/ci.yml` | Offline Go test and vet validation |
| `.github/workflows/integration.yml` | Manual/weekly networked Mars integration validation |
| `pipeline_test.go` | Networked Mars integration test, gated by the `integration` build tag |

## Solar Stack Remote Context

Known GitHub repositories:

| Repo | Visibility | Role | Latest public release |
|------|------------|------|-----------------------|
| `ayoub3bidi/mars` | Public | FastAPI + Alembic template, default API path | `v0.1.1` |
| `ayoub3bidi/mercury` | Public | FastAPI + Flyway template, alternative API path | `v0.4.0` |
| `ayoub3bidi/venus` | Public | Docusaurus + Tailwind docs template | `v1.1.1` |
| `ayoub3bidi/solar-commons` | Private | Go scaffold engine consumed by `sun` | none yet |

Mars was derived from Mercury, so local Mars checkouts can contain stale Mercury tags such as `v0.4.0`. Trust `git ls-remote --tags origin` or GitHub releases when choosing `default_ref`; Mars' current public release line is `v0.1.x`.

Template manifest readiness:

- Mars has `solar.manifest.yaml` and is the first integration target.
- Mercury and Venus still need `solar.manifest.yaml` before `sun` should expose them.
- Keep end-user scaffolds pinned to public release tags, not `main`.

## Public Contract

The package currently exposes these primitives:

```go
Scaffold(config ScaffoldConfig) (ScaffoldResult, error)
CloneTemplate(repo, ref, dest string) error
StripGit(dest string) error
ParseManifest(dest string) (Manifest, error)
ApplyRewrites(dest string, manifest Manifest, vars map[string]string) error
PruneFeature(dest string, feature Feature) error
InitGit(dest, message string) error
```

Public sentinel errors are available for `errors.Is`: `ErrInvalidConfig`, `ErrDestinationExists`, `ErrInvalidManifest`, `ErrMissingVariable`, `ErrUnknownFeature`, and `ErrUnsafePath`.

Callers are expected to build a `ScaffoldConfig` first and call `Scaffold`. Primitive helpers remain available for tests and advanced callers. Do not make product or UX decisions inside this package.

## Manifest Rules

`solar.manifest.yaml` is the contract between public templates and `sun`.

- `id`, `repo`, and `default_ref` are required.
- `default_ref` should be a stable semver tag for end-user scaffolds.
- `variables[].name` and `features[].name` must be unique and non-empty.
- Rewrites are simple literal replacements; they are not regexes.
- Replacement values may reference variables with `{{.project_name}}` syntax.
- Missing variable values are errors, not silent no-ops.
- Manifest-controlled paths must be relative and must not escape the destination root.

## Testing

Run from this directory:

```bash
go test ./...
```

The default test suite must stay offline. Any test that clones GitHub or another network remote belongs behind the integration build tag:

```bash
go test -tags integration -timeout 120s ./...
```

## Development Notes

- Use `gofmt` after edits.
- Run `go test ./...` and `go vet ./...` before every commit.
- Keep errors contextual; callers should be able to report useful CLI failures.
- Preserve file permissions when rewriting or pruning file contents.
- Avoid shelling out to `git`; this package uses go-git so `sun` remains self-contained.
- Do not add Bubble Tea, Cobra, or installer concerns here.
- Do not embed Mercury, Mars, or Venus-specific logic. Template differences belong in manifests.
- Be careful with broad string rewrites. The main product risk is incomplete or over-broad manifest replacement coverage.
