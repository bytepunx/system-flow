---
id: T-3389
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0794
owner: alex
created: 2026-11-15T21:47:08Z
updated: 2026-11-16T00:39:21Z
transitions:
  - to: ready
    at: 2026-11-16T00:38:46Z
    by: agent-S-0794
  - to: in-progress
    at: 2026-11-16T00:38:46Z
    by: agent-S-0794
  - to: done
    at: 2026-11-16T00:39:21Z
    by: agent-S-0794
stream: S-0794
tags: []
after: [T-3388]
usage:
  source: log
  seconds: 35
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1192
      cache_read: 664129
      cache_write: 25242
      cost: 0.3586
---
# T-3389 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0794 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-16: the publish exited 0 in the main checkout and pushed.

| Component | From | To | Tag or publish | Items |
|-----------|------|----|----------------|-------|
| flai | 1.100.0 | 1.101.0 | `flai/v1.101.0` | S-0755 (minor), S-0790 (patch) |
| template | 1.1.42 | 1.1.43 | `git@github.com:bytepunx/system-flow-template.git` v1.1.43 (9da14f1) | S-0755 (patch) |
