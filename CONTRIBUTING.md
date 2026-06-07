# Contributing

Thanks for helping improve `solar-commons`.

`solar-commons` is the scaffold engine for the planned `sun` CLI. Keep this repository focused on deterministic template operations: clone, strip git history, parse manifests, rewrite files, prune disabled features, and initialize a clean git repository.

## Local Setup

Use the Go version declared in `go.mod`.

```bash
go test ./...
go vet ./...
```

Networked template checks are intentionally separate:

```bash
go test -tags integration -timeout 120s ./...
```

## Contribution Rules

- Keep the root package as the only public import surface.
- Put implementation logic under `internal/`.
- Keep `sun` concerns out of this repository: prompts, Bubble Tea, Cobra, installers, registry UI, and product labels belong in `sun`.
- Keep default tests offline. Any GitHub/network test must use the `integration` build tag.
- Preserve path safety for manifest-controlled paths.
- Preserve contextual errors and exported sentinel errors.
- Do not add Mars, Mercury, or Venus-specific branching logic. Template behavior belongs in manifests.

## Pull Requests

Before opening a pull request:

1. Run `go test ./...`.
2. Run `go vet ./...`.
3. Run the integration test if the change affects clone/scaffold behavior.
4. Update `README.md` and `Agents.md` if public behavior or architecture changes.

Small, focused pull requests are easier to review.
