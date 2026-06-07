# Security Policy

## Supported Versions

`solar-commons` is pre-1.0. Security fixes will target the latest commit on `main` and the latest released `v0.x` tag when releases exist.

## Reporting a Vulnerability

Please do not open a public issue for vulnerabilities.

Use GitHub private vulnerability reporting or open a private security advisory for this repository. Include:

- affected version or commit
- reproduction steps
- expected and actual behavior
- impact assessment
- any suggested fix

The maintainer will acknowledge valid reports as soon as practical and coordinate a fix before public disclosure.

## Security-Relevant Areas

Pay special attention to:

- manifest-controlled path handling
- cloning untrusted repositories or refs
- rewrite behavior that could leak unresolved template values
- feature pruning that removes or edits generated project files
- git initialization behavior in generated projects
