# solar-commons

Core Go library powering the [sun](https://github.com/ayoub3bidi/sun) CLI. Handles the full scaffold pipeline: clone a Solar template at a pinned ref, apply variable rewrites, prune optional features, and initialise a clean git history.

## Pipeline

```
Scaffold  →  clone  →  strip .git  →  manifest  →  rewrite  →  prune  →  git init
```

## API

```go
import sc "github.com/ayoub3bidi/solar-commons"

// Clone a template at a semver tag and strip its git history
sc.CloneTemplate("https://github.com/ayoub3bidi/mars", "v0.1.1", "./my-api")
sc.StripGit("./my-api")

// Read the template's manifest
m, _ := sc.ParseManifest("./my-api")

// Substitute variables (resolves {{.var_name}} tokens)
sc.ApplyRewrites("./my-api", m, map[string]string{
    "project_name": "my-api",
    "app_title":    "My API",
    "api_version":  "v1",
})

// Remove an optional feature (paths + import lines)
sc.PruneFeature("./my-api", m.Features[0])

// Fresh single-commit repo
sc.InitGit("./my-api", "Initial commit")
```

## Manifest format

Each Solar template repo carries a `solar.manifest.yaml` at its root:

```yaml
id: mars
repo: https://github.com/ayoub3bidi/mars
default_ref: v0.1.1
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

## Templates

| ID | Repo | Stack |
|----|------|-------|
| `mars` | [ayoub3bidi/mars](https://github.com/ayoub3bidi/mars) | FastAPI + Alembic (recommended) |
| `mercury` | [ayoub3bidi/mercury](https://github.com/ayoub3bidi/mercury) | FastAPI + Flyway |
| `venus` | [ayoub3bidi/venus](https://github.com/ayoub3bidi/venus) | Docusaurus + Tailwind |

## Testing

```bash
# Unit tests
go test ./...

# Full pipeline integration test (clones Mars v0.1.1 from GitHub)
go test -tags integration -timeout 120s ./...
```

## License

MIT
