---
id: T-3612
type: task
nature: improvement
title: List the open issues the analyzer's report names
status: done
parent: S-0850
owner: alex
created: 2026-11-23T06:00:22Z
updated: 2026-11-27T22:49:38Z
transitions:
  - to: ready
    at: 2026-11-27T22:49:23Z
    by: agent-S-0850
  - to: in-progress
    at: 2026-11-27T22:49:23Z
    by: agent-S-0850
  - to: done
    at: 2026-11-27T22:49:38Z
    by: agent-S-0850
stream: S-0850
tags: []
after: [T-3611]
usage:
  source: log
  seconds: 15
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 480
      cache_read: 548733
      cache_write: 10638
      cost: 0.2045
---
# T-3612 List the open issues the analyzer's report names

## Work

Run `scripts/flai.sh issue list --report latest --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh issue list --report latest --json` exits 0 in the main checkout.

## Notes

- Exit: 0, in the main checkout, on 2026-11-27.
- Report: `design/analysis/2026-11-27-all.md`. It names ten open issues:

| Issue | Class | Count | Cost | Open story |
|-------|-------|-------|------|------------|
| I-0388 | efficiency | 4 | 41m | S-0846 |
| I-0390 | defect | 2 | 15m | none |
| I-0411 | efficiency | 1 | 76h | none |
| I-0412 | defect | 1 | 14h | none |
| I-0413 | efficiency | 1 | 48h | none |
| I-0414 | efficiency | 3 | 39m | none |
| I-0415 | efficiency | 2 | 74h10m | none |
| I-0416 | impression | 1 | 5m | none |
| I-0417 | impression | 1 | 5m | none |
| I-0418 | impression | 1 | 5m | none |
