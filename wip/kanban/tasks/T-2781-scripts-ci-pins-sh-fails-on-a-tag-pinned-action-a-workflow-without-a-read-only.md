---
id: T-2781
type: task
nature: remediation
title: scripts/ci-pins.sh fails on a tag-pinned action, a workflow without a read-only top-level permissions block, an unpinned image, or a second flaiover lockfile, and smoke runs it
status: in-progress
parent: S-0630
owner: alex
created: 2026-11-09T21:58:48Z
updated: 2026-11-10T05:53:58Z
transitions:
  - to: ready
    at: 2026-11-10T05:53:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:53:58Z
    by: system-flow
stream: S-0630
tags: [ci, security, test]
touches: [scripts/ci-pins.sh, scripts/smoke.sh, scripts/README.md]
after: [T-2774, T-2778, T-2779]
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
