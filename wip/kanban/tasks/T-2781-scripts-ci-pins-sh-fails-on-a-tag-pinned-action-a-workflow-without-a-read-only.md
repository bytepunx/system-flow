---
id: T-2781
type: task
nature: remediation
title: scripts/ci-pins.sh fails on a tag-pinned action, a workflow without a read-only top-level permissions block, an unpinned image, or a second flaiover lockfile, and smoke runs it
status: done
parent: S-0630
owner: alex
created: 2026-11-09T21:58:48Z
updated: 2026-11-10T05:55:47Z
transitions:
  - to: ready
    at: 2026-11-10T05:53:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:53:58Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:55:47Z
    by: agent-S-0630
stream: S-0630
tags: [ci, security, test]
touches: [scripts/ci-pins.sh, scripts/smoke.sh, scripts/README.md]
after: [T-2774, T-2778, T-2779]
usage:
  source: log
  seconds: 109
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 4925
      cache_read: 2566378
      cache_write: 41314
      cost: 0.9306
---
# T-2781 scripts/ci-pins.sh fails on a tag-pinned action, a workflow without a read-only top-level permissions block, an unpinned image, or a second flaiover lockfile, and smoke runs it

## Work

- Write `scripts/ci-pins.sh`, POSIX `sh` under `set -eu`, taking the repository root as an optional argument (default: the script's parent folder). It prints one line per finding, `<path>:<line>: <reason>`, and exits 1 when there is any. The findings are:
  - a `uses:` in `.github/workflows/*.yml` whose ref is not 40 hexadecimal characters (a local `./` action is exempt);
  - a workflow with no top-level `permissions:`, or one whose top-level block grants `write` or is `write-all`;
  - a `FROM` in `flaiover/Dockerfile` or `flai/images/*/Dockerfile`, or an `image:` in `flaiover/compose.yaml`, that names `latest`, `main-stable`, no tag, or a tag with fewer than three version parts. A build argument's default stands for `${...}`, and `FROM ${BASE}` naming a local `:local` image is exempt;
  - `scripts/litellm.sh`'s default `LITELLM_IMAGE` tag, checked the same way;
  - more than one of `pnpm-lock.yaml`, `package-lock.json`, and `yarn.lock` tracked under `flaiover/`, or a `package-lock.json` at the root.
- Call it from `scripts/smoke.sh` beside `scripts/check.sh`, so that CI's `system-flow-check.yml` and the close-out's smoke tier run it. Add it to `scripts/README.md`.
- Show that it reproduces I-0297. Run it against the story's base, `git archive <base> | tar -x -C <tmp>` followed by `scripts/ci-pins.sh <tmp>`, and record in the task's Notes that it exits 1 and names each instance the analysis lists.

It waits for T-2774, T-2778, and T-2779: it must pass on the tree they pin, and it reads the files they change.

## Done when

- `scripts/ci-pins.sh` exits 0 on the story's tree.
- Run on the story's base, it exits 1 and names the tag-pinned actions, the three workflows with no `permissions:`, the top-level write grants, the tag-pinned images, and the two `package-lock.json` files.
- `scripts/shellcheck.sh scripts/ci-pins.sh scripts/smoke.sh` passes, and `scripts/smoke.sh` runs it.

## Notes

On the story's tree, `scripts/ci-pins.sh` exits 0. On the story's base, `git archive main | tar -x -C <tmp>` then `scripts/ci-pins.sh <tmp>`, it exits 1 with 51 findings: 42 tag-pinned `uses:` across the six workflows, the missing top-level `permissions:` of `flai.yml`, `flaiover.yml`, and `system-flow-check.yml`, the top-level writes of `release-flai.yml` (contents, id-token, attestations) and `release-flaiover.yml` (packages), `node:24-alpine` twice, `golang:1.26-alpine`, `alpine:3.24`, `golang:1.26.0-alpine` (an `-alpine` variant with no Alpine release, a rule added beyond the plan), `grafana/otel-lgtm:latest`, LiteLLM's `main-stable`, the second flaiover lockfile, and the root `package-lock.json`. Outside a git checkout it reads lockfiles from disk, inside one from `git ls-files`. Written for mawk too (no regex intervals). `flai test` on the three paths: shellcheck and markdown passed.
