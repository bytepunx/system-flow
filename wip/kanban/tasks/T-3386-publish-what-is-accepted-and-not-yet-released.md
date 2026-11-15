---
id: T-3386
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0793
owner: alex
created: 2026-11-15T21:17:19Z
updated: 2026-11-15T21:32:57Z
transitions:
  - to: ready
    at: 2026-11-15T21:31:22Z
    by: agent-S-0793
  - to: in-progress
    at: 2026-11-15T21:31:22Z
    by: agent-S-0793
blocked:
  - from: 2026-11-15T21:31:38Z
    until: 2026-11-15T21:32:57Z
    reason: "step publish (flai release --pending --require-pending --json) exited 1: git commit found .git/index.lock held by another git process"
stream: S-0793
tags: []
after: [T-3385]
usage:
  source: log
  seconds: 58
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 4648
      cache_read: 1997861
      cache_write: 83289
      cost: 1.159
---
# T-3386 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0793 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-15T21:31:24Z: ran in the main checkout and exited 1. The formula's stop lists only exit 2, so the step failed. The task is blocked and the question is on a thread on S-0793.
- Output (the whole of it, one line):

  ```text
  {"ts":"2026-11-15T21:31:24.165Z","level":"FATAL","msg":"command failed","component":"cmd","err":"git commit -q -m chore: publish flai 1.99.0: exit status 128\nfatal: Unable to create '/home/alex/git/bytepunx/system-flow/.git/index.lock': File exists.\n\nAnother git process seems to be running in this repository, e.g.\nan editor opened by 'git commit'. Please make sure all processes\nare terminated then try again. If it still fails, a git process\nmay have crashed in this repository earlier:\nremove the file manually to continue."}
  ```

- State found afterwards, read-only: `.git/index.lock` no longer exists, so another git process held it only briefly. No tag `flai/v1.99.0` exists (the newest is `flai/v1.98.0`), HEAD is still `2fd4f33d4`, and nothing was pushed. The main checkout's `wip/` and `design/digest/` changes are now staged in the index, not committed.
