---
title: Developer experience and operations
updated: 2026-09-15
status: active
---

# Developer experience and operations

What the template ships so a conforming repo is pleasant to work in from the first commit, and what this monorepo uses on top.

## Shipped by the template

| Component | File | Purpose |
|-----------|------|---------|
| Editor config | `.editorconfig` | UTF-8, LF, final newline, 2-space indent for yaml/md/json, tabs for Go and Makefile |
| Git ignore | `.gitignore` | OS and editor noise, `node_modules`, Go build output, `.flai-cache/` |
| Git attributes | `.gitattributes` | LF normalisation, markdown treated as text with diff, lockfiles marked generated |
| Markdown lint | `.markdownlint.yaml` | Relaxed line length, consistent headings and lists |
| Make | `Makefile` | `make check` runs `flai check` and markdownlint, `make dashboard` runs `flai dashboard`, `make board`, `make stats` |
| CI | `.github/workflows/system-flow-check.yml` | Runs `flai check --strict` and markdownlint on pull requests |
| Agent instructions | `CLAUDE.md` | The baseline described in [template.md](template.md) |
| Pull request template | `.github/pull_request_template.md` | Asks for the story ID and confirms the definition of done |

## Used by this monorepo

| Concern | Choice | Detail |
|---------|--------|--------|
| Go build and lint | `go build`, `golangci-lint` | Config in `flai/.golangci.yaml` |
| Go release | GoReleaser | `flai/.goreleaser.yaml`, tags `flai/v*` produce GitHub release binaries for linux, darwin, windows |
| Node package manager | pnpm | Lockfile committed, `corepack` pins the version |
| Dashboard image | GitHub Actions, `docker/build-push-action` | Tags `flaiover/v*` and `main` push to GHCR with `latest` and semver tags |
| Task runner | Root `Makefile` | Delegates into `flai/` and `flaiover/` |
| Versioning | Independent per sub-project | Git tags prefixed with the project name |

Details and versions are in `design/tech/`.
