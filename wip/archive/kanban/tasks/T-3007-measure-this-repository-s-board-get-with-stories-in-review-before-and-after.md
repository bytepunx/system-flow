---
id: T-3007
type: task
nature: improvement
title: Measure this repository's board.get with stories in review before and after, record it in the design, and close I-0333
status: done
parent: S-0698
owner: alex
created: 2026-11-10T10:31:24Z
updated: 2026-11-11T08:53:10Z
transitions:
  - to: ready
    at: 2026-11-11T08:52:58Z
    by: agent-S-0698
  - to: in-progress
    at: 2026-11-11T08:52:58Z
    by: agent-S-0698
  - to: done
    at: 2026-11-11T08:53:10Z
    by: agent-S-0698
stream: S-0698
tags: [cli, docs]
touches: [design/system/server-performance.md, design/system/flai-cli.md, design/issues/I-0333-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md, design/issues/summary.md]
after: [T-3003]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 22
      output: 5454
      cache_read: 3151388
      cache_write: 41971
      cost: 1.0522
---
# T-3007 Measure this repository's board.get with stories in review before and after, record it in the design, and close I-0333

## Work

Measure `board.get` as S-0636 did, in `design/system/server-performance.md` § 13 boards at once, measured again after S-0602 and S-0606 (S-0636). Use one in-process program through `hostapi.MethodsFor`, with the file watcher and the main-refs watcher running as flai serve runs them. Compare `main` at the story's base with the story's branch, with at least five stories in review.

- Where the agent can reach the operator's host's `~/.flai/serve/projects.json`, run the 13 projects at once.
- Otherwise run this repository's board alone, and say so in the section.

Record the `merge.queue` phase and the git processes a call ran.

Add a section to `design/system/server-performance.md`, Cause 22 fixed (S-0698). It gives the cause, the fix, and the before and after, with the host's load average. Add a sentence to the merge queue paragraph of `design/system/flai-cli.md` on what `Compute` reads with git: one heads read, and `git status` only while a story in review has a passing verify result.

Close the issue with `flai issue close I-0333 --reason "<what fixed it>"`. That command rewrites the issue file and `design/issues/summary.md`.

If this repository's calls still miss 300 ms only because of the host's load, ask on a thread on S-0698 before ticking criterion 1, as D-0535 and D-0556 did. Do not tick it on a reworded reading.

It waits for T-3003, whose change it measures.

## Done when

- `design/system/server-performance.md` records the before and after, how they were measured, and the load.
- `design/system/flai-cli.md` says what the queue reads with git.
- `flai issue close I-0333` has run with a reason naming T-3000 and T-3003's change. The summary no longer lists I-0333 as open.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
