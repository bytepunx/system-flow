---
title: CI and repository automation
updated: 2026-09-15
status: active
---

# CI and repository automation

| Component | Use |
|-----------|-----|
| GitHub Actions | All CI |
| `markdownlint-cli2` | Markdown style in monorepo and template output |
| `flai check --strict` | Standard conformance, run from a built `flai` in CI |
| GoReleaser action | Release `flai` on tags `flai/v*`. GoReleaser OSS cannot strip a monorepo tag prefix, so `flai/.goreleaser.yaml` derives the version with `trimprefix .Tag "flai/v"` in every template and the workflow passes `--skip=validate`. `git.ignore_tags` excludes `flaiover/*` when finding the previous tag. |
| `docker/build-push-action` with `docker/metadata-action` | Build and push `flaiover` on tags `flaiover/v*` and on `main` |
| GHCR | Image registry, public |
| Dependabot | Weekly for Go modules, npm, GitHub Actions, Docker base images |

Workflows in this monorepo:

| Workflow | Trigger | Does |
|----------|---------|------|
| `system-flow-check.yml` | pull request, push to main | markdownlint, build flai, `flai check --strict`, render `./template` and check the result |
| `flai.yml` | changes under `flai/` | lint, test with race detector, build |
| `flaiover.yml` | changes under `flaiover/` | lint, unit tests, build, e2e |
| `release-flai.yml` | tag `flai/v*` | GoReleaser |
| `release-flaiover.yml` | tag `flaiover/v*`, push to main | build and push image |
