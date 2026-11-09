---
id: S-0630
type: story
nature: improvement
title: CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles
status: ready
owner: alex
created: 2026-11-07T23:40:47Z
updated: 2026-11-09T22:33:17Z
transitions:
  - to: ready
    at: 2026-11-09T22:05:06Z
    by: alex
tags: [flaiover, template, ci, security]
topics: [ci, security, release]
touches: [".github/workflows/flai.yml", ".github/workflows/flaiover.yml", ".github/workflows/system-flow-check.yml", ".github/workflows/release-flai.yml", ".github/workflows/release-flaiover.yml", ".github/workflows/install-published.yml", ".github/dependabot.yml", template/root/.github/workflows/system-flow-check.yml, flaiover/Dockerfile, flai/images/bedrock/Dockerfile, flai/images/system-flow/Dockerfile, flaiover/compose.yaml, scripts/litellm.sh, flaiover/package-lock.json, package-lock.json, scripts/ci-pins.sh, scripts/smoke.sh, scripts/README.md, design/tech/ci.md, design/tech/docker.md, design/tech/node.md, design/tech/observability.md, design/system/release-signing.md, design/issues/I-0297-ci-actions-and-base-images-are-pinned-by-tag-three-workflows-have-no.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 81
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 48
          cache_read: 473430
          cache_write: 12681
          cost: 0.1293
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759049
          cache_write: 13854
          cost: 0.278
forecast:
  duration: 1h
  delivery: 2026-11-10T00:39:00Z
  basis: "Its own forecast of 1h; 7th in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0633, S-0634, S-0645, S-0644, S-0638, S-0636, S-0640 and S-0625."
  by: flai
  at: 2026-11-09T22:33:17Z
finalized:
  by: alex
  at: 2026-11-09T22:04:57Z
---
# S-0630 CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles

## Goal

This story remediates [I-0297](../../../design/issues/I-0297-ci-actions-and-base-images-are-pinned-by-tag-three-workflows-have-no.md), "CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0297 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0297 is closed with `flai issue close I-0297 --reason` saying what fixed it

## Tasks
- T-2774 Every GitHub Actions step names a commit SHA, and every workflow grants contents read at the top level and writes only on the jobs that need them
- T-2778 Every base and local-stack image names an exact version, the released flaiover image its digest too, and dependabot watches the agent image folders
- T-2779 The root and flaiover package-lock.json files, which nothing reads, are removed, leaving flaiover/pnpm-lock.yaml the one lockfile
- T-2781 scripts/ci-pins.sh fails on a tag-pinned action, a workflow without a read-only top-level permissions block, an unpinned image, or a second flaiover lockfile, and smoke runs it
- T-2783 The tech and release-signing documents say how actions, images, and the lockfile are pinned and which jobs write, and I-0297 is closed

## Notes

The analyzer found I-0297: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0297: analyzer 41s, 0.8156 USD.

### Planning

The proposed solution, from the instances in the analysis's Finding R5. The plan's thread on this story gives it for review.

- Actions: every `uses:` in `.github/workflows/` names a full commit SHA with its version in a trailing comment (`# v7.0.1`). `.github/dependabot.yml` already watches `github-actions` and updates SHA pins and their comments. GoReleaser's `version: "~> v2"` in `release-flai.yml` becomes the exact `v2.18.1` that `scripts/install-tools.sh` pins.
- Permissions: `flai.yml`, `flaiover.yml`, and `system-flow-check.yml` get a top-level `permissions: contents: read`. The release workflows keep `contents: read` at the top level and grant `contents`, `packages`, `id-token`, and `attestations` write only on the jobs that need them. The template's `system-flow-check.yml` gets the same `permissions:` block. It keeps its version tags, because a project made from it has no dependabot to move a SHA.
- Images: each `FROM` and `image:` names an exact version: `golang:<x.y.z>-alpine<a.b>`, `alpine:<a.b.c>`, `node:<x.y.z>-alpine<a.b>`, `grafana/otel-lgtm:<version>`, and a `ghcr.io/berriai/litellm` release tag in place of `main-stable`. `flaiover/Dockerfile` also adds the image digest, since a signed release builds from it. Dependabot's `docker` entry is widened from `/flaiover` to the `flai/images` folders.
- Lockfiles: `flaiover/package-lock.json` and the root `package-lock.json`, read by nothing (`design/tech/node.md`), are removed. The caret ranges in `flaiover/package.json` stay: `--frozen-lockfile` on `flaiover/pnpm-lock.yaml` is the pin.
- The test: `scripts/ci-pins.sh` fails on a `uses:` that is not a SHA, a workflow with no top-level `permissions:` or with a write grant there, an image with no exact version, and a second tracked lockfile under `flaiover/`. `scripts/smoke.sh` runs it, so CI's `system-flow-check.yml` and the close-out's smoke tier run it.
- Out of scope: the signing key's `--tlog-upload=false` and `--insecure-ignore-tlog`, which [ADR-0070](../../../design/adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md) decided. Also out of scope: the agent images, which no workflow builds or publishes.

Where each touch came from:

| Touch | Source |
|-------|--------|
| `.github/workflows/*.yml` (six files), `flaiover/Dockerfile`, `flai/images/bedrock/Dockerfile`, `flaiover/compose.yaml`, `scripts/litellm.sh`, `flaiover/package-lock.json`, `package-lock.json` | design: the instances in I-0297 and the analysis's Finding R5. `install-published.yml` is added from the layout: it has a `permissions:` block, but its actions are tag-pinned too |
| `flai/images/system-flow/Dockerfile` | layout: its `golang` and `koalaman/shellcheck-alpine` stages are pinned by tag as the bedrock image's are |
| `template/root/.github/workflows/system-flow-check.yml` | layout: the template ships the same workflow with no `permissions:` block |
| `.github/dependabot.yml` | layout: SHA and image pins need its `docker` entry to reach the `flai/images` folders |
| `scripts/ci-pins.sh`, `scripts/smoke.sh`, `scripts/README.md` | criterion 1: the test that reproduces the cause, run by the smoke tier and CI, and indexed in the scripts README |
| `design/tech/ci.md`, `design/tech/docker.md`, `design/tech/node.md`, `design/tech/observability.md`, `design/system/release-signing.md` | design: they name the actions, the image tags, the two lockfiles, the `otel-lgtm` stack, and the release workflows' `permissions:` |
| `design/issues/I-0297-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

`flai touches suggest S-0630` listed by co-change only hub files: `docs/users/flai.md`, `design/system/flai-cli.md`, `docs/operators/index.md`, and close-out issues. None of them names a pin or a workflow permission, so none was taken. No folder touch is kept. The story has no `after`: it reads no other open story's output.

The figures:

- Forecast: `flai forecast` gave 43m, the median of 98 s per unit of size over 142 large-band improvement stories, times size 27. It is raised to 1h. Each of about 14 action pins needs its commit looked up, and so does each of 6 image versions. The template and `flai/images` paths also select the integration and smoke tiers at close-out. Delivery moves by the same 17m, to 2026-11-09T23:16:00Z.
- Cost of delay: none yet. The story has no inputs and no epic. The plan's thread asks the operator for them.
