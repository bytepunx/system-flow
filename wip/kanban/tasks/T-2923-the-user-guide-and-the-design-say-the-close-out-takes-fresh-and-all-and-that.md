---
id: T-2923
type: task
nature: improvement
title: The user guide and the design say the close-out takes --fresh and --all and that the job passes them on
status: done
parent: S-0671
owner: alex
created: 2026-11-10T08:48:29Z
updated: 2026-11-10T09:41:59Z
transitions:
  - to: ready
    at: 2026-11-10T09:41:33Z
    by: agent-S-0671
  - to: in-progress
    at: 2026-11-10T09:41:35Z
    by: agent-S-0671
  - to: done
    at: 2026-11-10T09:41:59Z
    by: agent-S-0671
stream: S-0671
tags: [docs, close-out]
touches: [docs/users/flai.md, design/system/flai-cli.md, design/system/devex.md]
after: [T-2918, T-2920]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 421
      cache_read: 477245
      cache_write: 4331
      cost: 0.1385
---
# T-2923 The user guide and the design say the close-out takes --fresh and --all and that the job passes them on

## Work

- In `docs/users/flai.md`, under `flai verify`: replace "The close-out runs `flai verify --record-issues` without `--fresh`" with what is true now, that it does so unless asked for `--fresh`, and say how to ask, by hand and as a job. In the paragraph on the MCP tool `verify`, say that `--all` and `--fresh` reach the worktree's script as `CLOSE_OUT_ALL` and `CLOSE_OUT_FRESH`.
- In `design/system/flai-cli.md`, in the close-out job's "What it runs" bullet, say the same, and why the job passes variables rather than flags: a worktree's older script ignores them.
- In `design/system/devex.md`, give the Close-out row the new usage, `scripts/close-out.sh [--fresh] [--all] S-nnnn -m "<message>"`.
- Waits for T-2918 and T-2920, since it describes what each built, the options and the variables, as they landed.

## Done when

- The three documents describe the options and variables as T-2918 and T-2920 built them, and none still says the close-out cannot run fresh.
- `flai test docs/users/flai.md design/system/flai-cli.md design/system/devex.md` passes, the markdown lint included.

## Notes
