---
id: T-1850
type: task
nature: feature
title: flai file show and flai file save, with --story, and flai design writes the stylesheet whenever it writes the token file
status: done
parent: S-0436
owner: alex
created: 2026-10-16T21:31:39Z
updated: 2026-10-16T22:20:14Z
transitions:
  - to: ready
    at: 2026-10-16T22:05:57Z
    by: agent-S-0436
  - to: in-progress
    at: 2026-10-16T22:05:58Z
    by: agent-S-0436
  - to: done
    at: 2026-10-16T22:20:14Z
    by: agent-S-0436
stream: S-0436
tags: [cli]
touches: [flai/cmd/file.go, flai/cmd/file_test.go, flai/cmd/root.go, flai/cmd/design.go, flai/cmd/design_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1849]
usage:
  source: log
  seconds: 856
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 58
      output: 265
      cache_read: 2815071
      cache_write: 128198
      cost: 1.3671
---
# T-1850 flai file show and flai file save, with --story, and flai design writes the stylesheet whenever it writes the token file

## Work

- `flai/cmd/file.go`: `flai file show <path> [--story] [--json]` and `flai file save <path> --hash <h> [--story] [--message]`, the content on standard input, as `flai doc show` and `flai doc save` in `flai/cmd/doc.go` are, with exit 3 for a conflict and exit 4 for a refusal. Register the group in `root.go`.
- `flai/cmd/design.go`: every subcommand that writes the token file, such as `flai design init`, writes the stylesheet the manifest names in the same step, through T-1849's save path.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Waits for T-1849, whose functions both commands call.

## Done when

Command tests cover show and save as text and JSON, each exit code, `--story`, and `flai design init` leaving the token file and the stylesheet in step with no `design.stale`. `flai test flai/cmd` passes.

## Notes

None.
