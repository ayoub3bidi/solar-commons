# Public release readiness

Before making this repository public:

- [ ] `go test ./...` and `go vet ./...` pass
- [ ] Integration tests pass for Mars (`TestFullPipelineMars`)
- [ ] Mercury/Venus integration pass after template tags include manifests
- [ ] README and Agents.md match pinned template refs
- [ ] `sun` consumes released module version (no `replace` in go.mod)

## Release

```bash
git tag v0.2.0
git push origin v0.2.0
```

The [release workflow](.github/workflows/release.yml) runs tests and creates a GitHub Release.
