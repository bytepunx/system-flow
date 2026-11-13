---
id: T-3302
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0771
owner: alex
created: 2026-11-13T13:41:40Z
updated: 2026-11-13T14:37:38Z
transitions:
  - to: ready
    at: 2026-11-13T14:37:10Z
    by: agent-S-0771
  - to: in-progress
    at: 2026-11-13T14:37:10Z
    by: agent-S-0771
  - to: done
    at: 2026-11-13T14:37:38Z
    by: agent-S-0771
stream: S-0771
tags: []
after: [T-3301]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1202
      cache_read: 523735
      cache_write: 22793
      cost: 0.3112
---
# T-3302 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0771 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran as `scripts/flai.sh release --pending --require-pending --json` (the manifest's `flai.command`) in the main checkout; exited 0 with `pushed: true`.

| Component | From | To | Tag or version | Items |
|-----------|------|----|----------------|-------|
| flai | 1.91.0 | 1.92.0 (minor) | `flai/v1.92.0` | S-0661, S-0745 |
| template | 1.1.36 | 1.1.37 (patch) | `template/template.yaml`; published to `git@github.com:bytepunx/system-flow-template.git` as `v1.1.37` (8693665) | S-0745 |
