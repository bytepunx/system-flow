---
id: T-3519
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0826
owner: alex
created: 2026-11-17T08:16:46Z
updated: 2026-11-17T08:19:38Z
transitions:
  - to: ready
    at: 2026-11-17T08:19:14Z
    by: agent-S-0826
  - to: in-progress
    at: 2026-11-17T08:19:14Z
    by: agent-S-0826
  - to: done
    at: 2026-11-17T08:19:38Z
    by: agent-S-0826
stream: S-0826
tags: []
after: [T-3518]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 95
      cache_read: 530537
      cache_write: 3036
      cost: 0.2033
---
# T-3519 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0826 --step publish --exit 2`, giving it the step's output.

## Notes

Exited 0 in the main checkout, `pushed: true`.

| Component | From | To | Tag or version | Items |
|-----------|------|----|----------------|-------|
| flai | 1.103.0 | 1.104.0 (minor) | `flai/v1.104.0` | S-0805, S-0750 |
| template | 1.1.44 | 1.1.45 (patch) | `template/template.yaml`, published to the template repository as v1.1.45 (315ea50) | S-0805 |
