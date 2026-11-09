---
id: T-2783
type: task
nature: remediation
title: The tech and release-signing documents say how actions, images, and the lockfile are pinned and which jobs write, and I-0297 is closed
status: backlog
parent: S-0630
owner: alex
created: 2026-11-09T21:59:02Z
updated: 2026-11-09T21:59:02Z
transitions: []
stream: S-0630
tags: [docs, ci]
touches: [design/tech/ci.md, design/tech/docker.md, design/tech/node.md, design/tech/observability.md, design/system/release-signing.md, design/issues/I-0297-ci-actions-and-base-images-are-pinned-by-tag-three-workflows-have-no.md, design/issues/summary.md]
after: [T-2774, T-2778, T-2779, T-2781]
---
# T-2783 The tech and release-signing documents say how actions, images, and the lockfile are pinned and which jobs write, and I-0297 is closed

## Work

- `design/tech/ci.md`: say that every action is pinned by commit SHA with its version in a comment, and that dependabot moves the pins. Say that each workflow grants `contents: read` at the top level and that the release jobs name their own writes. Say that `scripts/ci-pins.sh`, run by `scripts/smoke.sh`, fails a pin or a grant that breaks this. Correct the `release-flai.yml` row's GoReleaser version.
- `design/tech/docker.md` rows Image and Agent images: give the exact tags and the flaiover base digest that T-2778 wrote, and say that dependabot watches the three image folders.
- `design/tech/node.md` line 20: the `package-lock.json` files are gone. `flaiover/pnpm-lock.yaml` is the one lockfile, and the caret ranges in `package.json` are held by it.
- `design/tech/observability.md` line 15: name the `grafana/otel-lgtm` version `flaiover/compose.yaml` now pins.
- `design/system/release-signing.md` lines 18 and 24: say which job of each release workflow holds which write grant, as T-2774 left them.
- Bump each document's `updated`.
- Close the issue with `flai issue close I-0297 --reason "<what fixed it>"`, naming the SHA pins, the permissions blocks, the exact image tags and digest, the removed lockfiles, and `scripts/ci-pins.sh`. That command writes the issue file and `design/issues/summary.md`.

It waits for T-2774, T-2778, T-2779, and T-2781: the documents state the pins and grants those tasks write, and the issue closes only once the check that reproduces it passes.

## Done when

- The five documents name no tag-only pin, no top-level write grant, and no second lockfile, and each bumps `updated`.
- `flai issue show I-0297` says closed, with the reason.
- `flai check --strict` and `scripts/lint-md.sh` pass on the changed files.

## Notes
