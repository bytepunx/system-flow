---
id: T-2547
type: task
nature: improvement
title: List the open issues the analyzer's report names
status: done
parent: S-0584
owner: alex
created: 2026-11-07T08:03:23Z
updated: 2026-11-07T08:33:49Z
transitions:
  - to: ready
    at: 2026-11-07T08:33:07Z
    by: agent-S-0584
  - to: in-progress
    at: 2026-11-07T08:33:10Z
    by: agent-S-0584
  - to: done
    at: 2026-11-07T08:33:49Z
    by: agent-S-0584
stream: S-0584
tags: []
after: [T-2545]
usage:
  source: log
  seconds: 39
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 1139
      cache_read: 1062925
      cache_write: 10352
      cost: 0.3182
---
# T-2547 List the open issues the analyzer's report names

## Work

Run `flai issue list --report latest --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai issue list --report latest --json` exits 0 in the main checkout.

## Notes

- 2026-11-07T08:33Z: `scripts/flai.sh issue list --report latest --json` in the main checkout exited 0.
- It names the report `design/analysis/2026-11-07-all.md` and 19 open issues: I-0230, I-0246, I-0257, I-0268, and I-0284 to I-0298. S-0584's Notes list each with its title and class.
