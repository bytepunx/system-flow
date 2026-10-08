---
id: T-1364
type: task
nature: remediation
title: lint-md.sh's whole-repository run leaves main's wip out in a close-out and lints only the wip files the story changes
status: in-progress
parent: S-0345
owner: alex
created: 2026-10-08T08:39:56Z
updated: 2026-10-08T09:09:25Z
transitions:
  - to: ready
    at: 2026-10-08T09:09:24Z
    by: agent-S-0345
  - to: in-progress
    at: 2026-10-08T09:09:25Z
    by: agent-S-0345
stream: S-0345
tags: [flai]
touches: [scripts/lint-md.sh]
---
# T-1364 lint-md.sh's whole-repository run leaves main's wip out in a close-out and lints only the wip files the story changes

## Work

`scripts/smoke.sh` runs `scripts/lint-md.sh` with no files, which lints every markdown file the branch holds, `wip/` included. The close-out runs smoke for any story that changes `flai/`, so a line any agent commits to main's `wip/` fails its smoke tier the same way it fails integration (I-0117, and I-0027 before it).

- In `scripts/lint-md.sh`, when `CLOSE_OUT_STORY` is set and no files are given, add `!wip/**` to the globs. Then add back the `.md` files under `wip/` that the branch changes against the main branch, so the story's own items and narrative are still linted.
- Leave the run with files given, as `flai test`'s markdown tier calls it, unchanged.
- Leave the run without `CLOSE_OUT_STORY`, as in `make lint-md` and CI, unchanged, so main's `wip/` is still linted there. `scripts/check.sh` scopes `flai check` the same way.
- Do not change `scripts/smoke.sh`: S-0340 changes it now, and the variable reaches `lint-md.sh` through the environment.
- Keep the script POSIX `sh` under `set -eu`, with every variable quoted, as `tooling.md` says.

Runs in the first layer: it waits for no task and shares no path with the `TestRepositoryLintsClean` task.

## Done when

- With `CLOSE_OUT_STORY` set and no files given, `scripts/lint-md.sh` passes in a story worktree while main's `wip/` holds a bare `www.` line the story does not change. It fails when the story changes that file.
- Without `CLOSE_OUT_STORY`, it still lints `wip/` and fails on that line.
- `flai test scripts/lint-md.sh` passes.

## Notes

Drafted by the planner for S-0345.
