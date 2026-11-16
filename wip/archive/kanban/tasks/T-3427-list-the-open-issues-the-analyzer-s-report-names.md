---
id: T-3427
type: task
nature: improvement
title: List the open issues the analyzer's report names
status: done
parent: S-0803
owner: alex
created: 2026-11-16T06:00:53Z
updated: 2026-11-16T06:18:46Z
transitions:
  - to: ready
    at: 2026-11-16T06:18:09Z
    by: agent-S-0803
  - to: in-progress
    at: 2026-11-16T06:18:10Z
    by: agent-S-0803
  - to: done
    at: 2026-11-16T06:18:46Z
    by: agent-S-0803
stream: S-0803
tags: []
after: [T-3426]
usage:
  source: log
  seconds: 36
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 1609
      cache_read: 826031
      cache_write: 28717
      cost: 0.4272
---
# T-3427 List the open issues the analyzer's report names

## Work

Run `scripts/flai.sh issue list --report latest --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh issue list --report latest --json` exits 0 in the main checkout.

## Notes

- Exit: 0, run in `/home/alex/git/bytepunx/system-flow` on 2026-11-16. A first attempt redirected the output into the worktree's `bin/S-0803/`, which did not exist yet, so the shell failed before the command ran; the command was then run again without the redirect.
- Report: `design/analysis/2026-11-16-all.md`. It names nine open issues, I-0386 to I-0394, listed with their titles and classes in S-0803's Notes.
