---
title: CI and repository automation
updated: 2026-10-08
status: active
topics: [code]
---

# CI and repository automation

| Component | Use |
|-----------|-----|
| GitHub Actions | All CI |
| `markdownlint-cli2` | Markdown style in monorepo and template output; `scripts/lint-md.sh` runs the same globs locally as part of `make smoke` |
| `flai check --strict` | Standard conformance, run from a built `flai` in CI |
| GoReleaser action | Release `flai` on tags `flai/v*`. GoReleaser OSS cannot strip a monorepo tag prefix, so `flai/.goreleaser.yaml` derives the version with `trimprefix .Tag "flai/v"` in every template and the workflow passes `--skip=validate`. `git.ignore_tags` excludes `flaiover/*` when finding the previous tag. `system-flow-check.yml` uses it with `install-only` to install GoReleaser v2.18.1, the version `scripts/install-tools.sh` pins, so the smoke tier can build releases from the tree (S-0340). |
| `docker/build-push-action` with `docker/metadata-action` | Build and push `flaiover` on tags `flaiover/v*` and on `main` |
| GHCR | Image registry, public |
| Releases | Automatic on acceptance, public or private, per `design/conventions/git.md`: for the component the item delivers to, feature story minor (an epic none of its own since S-0200, ADR-0078), remediation, improvement, docs-only, or non-breaking dependency update patch; components touched incidentally with additive changes get a patch. Research and experiment stay on branches. Tags are per sub-project (`flai/v*`, `flaiover/v*`); `flai release` (S-0029) will compute and push them. Remotes are created with `gh` after asking the operator for organization and visibility. |
| Dependabot | Weekly for Go modules, npm, GitHub Actions, Docker base images |

Workflows in this monorepo:

| Workflow | Trigger | Does |
|----------|---------|------|
| `system-flow-check.yml` | pull request, push to main | markdownlint, build flai, install GoReleaser, then `scripts/smoke.sh`: render `./template` and check the result, `flai check --strict`, the markdown lint, and install and self-upgrade flai from releases built from the tree and served on `127.0.0.1`, with no token and no request to GitHub (S-0340) |
| `install-published.yml` | push to main, daily at 06:17 UTC, by hand | build flai, then `scripts/install-published-test.sh`: install the latest published release from GitHub with `install.sh` and `flai self-upgrade`, with the workflow's `GITHUB_TOKEN`. No story's close-out runs this check; a failure fails the run, which is the record, and opens no issue |
| `flai.yml` | changes under `flai/` | lint, test with race detector on Linux and macOS (S-0180: macOS's temp dir is a symlink, its `/bin/sh` is bash, and its git works out an identity from the hostname, which broke tests only a Mac ran), build |
| `flaiover.yml` | changes under `flaiover/` | lint, unit tests, build, e2e |
| `release-flai.yml` | tag `flai/v*` | GoReleaser |
| `release-flaiover.yml` | tag `flaiover/v*`, push to main touching flaiover or flai | buildx build of `flaiover/Dockerfile` from the repo root, push to GHCR with latest, semver, major, and sha tags, flai version from the latest `flai/v*` tag |
