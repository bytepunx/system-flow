---
id: T-1994
type: task
nature: remediation
title: scripts/integration.sh installs flaiover's dependencies in a worktree before the Go run, as flaiover-unit.sh does
status: done
parent: S-0482
owner: alex
created: 2026-10-21T12:08:19Z
updated: 2026-10-21T12:22:01Z
transitions:
  - to: ready
    at: 2026-10-21T12:20:18Z
    by: agent-S-0482
  - to: in-progress
    at: 2026-10-21T12:20:19Z
    by: agent-S-0482
  - to: done
    at: 2026-10-21T12:22:01Z
    by: agent-S-0482
stream: S-0482
tags: [scripts, tests]
touches: [scripts/integration.sh, scripts/README.md, design/system/devex.md]
usage:
  source: log
  seconds: 102
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 9
      output: 48
      cache_read: 628663
      cache_write: 6238
      cost: 0.215
---
# T-1994 scripts/integration.sh installs flaiover's dependencies in a worktree before the Go run, as flaiover-unit.sh does

## Work

I-0178's cause: a story worktree has no `flaiover/node_modules`, so the components integration tests link the main checkout's, which nothing reinstalls after an acceptance that changes `flaiover/pnpm-lock.yaml`. Give each worktree its own install from its own lockfile instead.

- In `scripts/integration.sh`, before `go test`, run `"$ROOT/scripts/flaiover-install.sh" --if-needed` when `CACHE_ROOT` differs from `ROOT`, as `scripts/flaiover-unit.sh` does. That covers a story worktree and the merge queue's batch worktree, both under `.flai-cache/worktrees/`. The main checkout is still never installed into, as `design/system/component-library.md` says.
- Keep the install's output on stderr, so the `-json` stream `flai test`'s integration tier reads stays clean on stdout.
- Say what the script now does in its header comment, in its row in `scripts/README.md`, and in the Test tiers row of `design/system/devex.md`, naming I-0178.

It waits for no task: it shares no path with the test-helper task.

## Done when

- `scripts/integration.sh` run in a story worktree without `flaiover/node_modules` installs it from that worktree's lockfile, then runs the Go tests.
- Run in the main checkout, it installs nothing.
- `scripts/README.md` and `design/system/devex.md` say so.
- `flai test scripts/integration.sh design/system/devex.md scripts/README.md` passes.

## Notes
