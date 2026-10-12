---
title: CI and repository automation
updated: 2026-10-12
status: active
topics: [code]
---

# CI and repository automation

| Component | Use |
|-----------|-----|
| GitHub Actions | All CI |
| `markdownlint-cli2` | Markdown style in monorepo and template output; `scripts/lint-md.sh` runs the same globs locally as part of `make smoke` |
| `flai check --strict` | Standard conformance, run from a built `flai` in CI |
| GoReleaser action | Release `flai` on tags `flai/v*`. GoReleaser OSS cannot strip a monorepo tag prefix, so `flai/.goreleaser.yaml` derives the version with `trimprefix .Tag "flai/v"` in every template and the workflow passes `--skip=validate`. `git.ignore_tags` excludes `flaiover/*` when finding the previous tag. |
| `docker/build-push-action` with `docker/metadata-action` | Build and push `flaiover` on tags `flaiover/v*` and on `main`. The action's digest output, the image index's, is what the digest list, the cosign signature, and the attestation of a tag's build name |
| GHCR | Image registry, public |
| cosign v3.1.3, installed by `sigstore/cosign-installer` v4.1.2 | Used by both release workflows. Signs flai's `checksums.txt`, flaiover's digest list, and the flaiover image with the release key, an ECDSA P-256 pair whose private half is the secrets `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD` of the environment `release` ([runbook](../../docs/operators/runbooks/release-key.md)), through GoReleaser's `signs` in `release-flai.yml` and `cosign sign-blob` and `cosign sign --key` in `release-flaiover.yml`, with no Rekor entry, so flai and `install.sh` verify with the standard library and `openssl` and nothing else ([ADR-0070](../adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md)). Pinned: its `--output-signature` and `--tlog-upload` flags are deprecated in v3 and go in v4, and the installer has no moving `v4` tag. Chosen over minisign, which `install.sh` could rarely verify, and over keyless signing, which needs Sigstore at verification ([release-signing.md](../system/release-signing.md)). |
| `actions/attest-build-provenance` v4 | A keyless SLSA provenance statement for flai's archives and `checksums.txt`, and for the flaiover image by its index digest, pushed to GHCR beside it, for anyone with `gh attestation verify`, beside the key signature. Skipped while `bytepunx/system-flow` is private: on the Team plan artifact attestations are for public repositories only. |
| Releases | Automatic on acceptance, public or private, per `design/conventions/git.md`: for the component the item delivers to, feature story minor (an epic none of its own since S-0200, ADR-0078), remediation, improvement, docs-only, or non-breaking dependency update patch; components touched incidentally with additive changes get a patch. Research and experiment stay on branches. Tags are per sub-project (`flai/v*`, `flaiover/v*`); `flai release` (S-0029) will compute and push them. Remotes are created with `gh` after asking the operator for organization and visibility. |
| Dependabot | Weekly for Go modules, npm, GitHub Actions, Docker base images |

Workflows in this monorepo:

| Workflow | Trigger | Does |
|----------|---------|------|
| `system-flow-check.yml` | pull request, push to main | markdownlint, build flai, install GoReleaser v2.18.1 (the GoReleaser action with `install-only`, at the version `scripts/install-tools.sh` pins), then `scripts/smoke.sh`: render `./template` and check the result, `flai check --strict`, the markdown lint, and install and self-upgrade flai from releases built from the tree and served on `127.0.0.1`, with no token and no request to GitHub (S-0340) |
| `install-published.yml` | push to main, daily at 06:17 UTC, by hand | build flai, then `scripts/install-published-test.sh`: install the latest published release from GitHub with `install.sh` and `flai self-upgrade`, with the workflow's `GITHUB_TOKEN`. No story's close-out runs this check; a failure fails the run, which is the record, and opens no issue |
| `flai.yml` | changes under `flai/` | lint, test with race detector on Linux and macOS (S-0180: macOS's temp dir is a symlink, its `/bin/sh` is bash, and its git works out an identity from the hostname, which broke tests only a Mac ran), build |
| `flaiover.yml` | changes under `flaiover/` | lint, unit tests, build, e2e |
| `release-flai.yml` | tag `flai/v*` | In the environment `release`: fails when the release key is missing, installs cosign, GoReleaser (builds, archives, `checksums.txt` and its signature `checksums.txt.sig`), then the provenance attestation when the repository is public |
| `release-flaiover.yml` | tag `flaiover/v*`, push to main touching flaiover or flai | Three jobs. `image`: buildx build of `flaiover/Dockerfile` from the repo root for linux/amd64 and linux/arm64, push to GHCR with latest (main), semver and major (tag), and sha tags, flai version from the latest `flai/v*` tag. On a push to main it is the only job, and its images are unsigned. On a tag, `release-key` runs first, in the environment `release`, and fails the run before anything is pushed when the release key is missing; after `image`, `release`, in the same environment, writes the digest list `flaiover_<version>.digests` with `scripts/flaiover-digests.sh`, signs it into `flaiover_<version>.digests.sig`, signs the image with `cosign sign --key` by its index digest (no Rekor entry), adds the provenance attestation when the repository is public, and publishes the release `flaiover vX.Y.Z`, not marked latest, with the list, its signature, and `cosign.pub` ([release-signing.md](../system/release-signing.md#the-digest-list)) |
