---
id: T-0815
type: task
nature: feature
title: flai touches suggest lists files often changed with the ones named, from git history
status: done
parent: S-0210
owner: alex
created: 2026-10-04T19:45:32Z
updated: 2026-10-04T19:59:35Z
transitions:
  - to: ready
    at: 2026-10-04T19:46:38Z
    by: agent-S-0210
  - to: in-progress
    at: 2026-10-04T19:46:38Z
    by: agent-S-0210
  - to: done
    at: 2026-10-04T19:59:35Z
    by: agent-S-0210
stream: S-0210
tags: [flai]
touches: [flai/internal/storygit, flai/internal/planning/cochange.go, flai/internal/planning/cochange_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go]
usage:
  source: log
  seconds: 777
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 81
      output: 30848
      cache_read: 4011436
      cache_write: 114769
      cost: 2.097
---
# T-0815 flai touches suggest lists files often changed with the ones named, from git history

## Work

Add `flai touches suggest <S-nnnn> [path...]`, a read-only subcommand of `flai touches`. Its seeds are the story's touches, its own and those of its tasks not cancelled, plus the paths given. `storygit.CommitFiles` reads each non-merge commit on the main branch with the files it changed, leaving out files under the wip folder and flai's bookkeeping commits (accept and archive, create, publish, default agent, import), then drops commits with no file left. `planning.CoChanged` scores every file outside the seeds by the number of seed commits it appears in, with its share of them, and keeps those at or above a minimum count, sorted by count then path, up to a limit. Text and `--json` output.

Waits for nothing: it shares no path with the forecast task and runs beside it in its own task worktree.

## Done when

- [ ] `flai touches suggest S-nnnn [path...]` prints the suggestions with counts and shares, `--json` the same structured, `--limit` and `--min` set the cut
- [ ] Bookkeeping commits and wip files are left out; a story with no seeds and no paths given is refused with what to pass
- [ ] Behavior tests pin the scoring on fixed commit lists and the bookkeeping filter on fixed subjects; the command test runs against a fixture repository or a fake runner

## Notes
