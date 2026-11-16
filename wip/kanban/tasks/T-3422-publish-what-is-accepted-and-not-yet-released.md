---
id: T-3422
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0801
owner: alex
created: 2026-11-16T01:54:24Z
updated: 2026-11-16T03:53:32Z
transitions:
  - to: ready
    at: 2026-11-16T03:53:06Z
    by: agent-S-0801
  - to: in-progress
    at: 2026-11-16T03:53:06Z
    by: agent-S-0801
  - to: done
    at: 2026-11-16T03:53:32Z
    by: agent-S-0801
stream: S-0801
tags: []
after: [T-3421]
usage:
  source: log
  seconds: 26
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 1400
      cache_read: 829839
      cache_write: 29613
      cost: 0.4309
---
# T-3422 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0801 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-16: the publish ran in the main checkout and exited 0, `pushed: true`. It released S-0795, S-0796, S-0798, and S-0800 as patches: flai 1.101.0 → 1.101.1 (tag `flai/v1.101.1`, on `origin`), and template 1.1.43 → 1.1.44 (published to the template repository as `v1.1.44`, 6466ec6). The commit is `f3d6d8711 chore: publish flai 1.101.1, template 1.1.44`, and main is level with `origin/main`.
