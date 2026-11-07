---
id: T-2613
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0597
owner: alex
created: 2026-11-07T22:31:23Z
updated: 2026-11-07T23:06:40Z
transitions:
  - to: ready
    at: 2026-11-07T23:06:06Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T23:06:07Z
    by: system-flow
  - to: done
    at: 2026-11-07T23:06:40Z
    by: agent-S-0597
stream: S-0597
tags: []
after: [T-2612]
usage:
  source: log
  seconds: 33
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 803
      cache_read: 450002
      cache_write: 15321
      cost: 0.2287
---
# T-2613 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes

- 2026-11-07: `flai release --evaluate --require-met` ran in the main checkout and exited 0. The binary was the one `scripts/flai.sh` built in the S-0597 worktree, whose branch had no commits beyond main. The threshold policy is met: 1 accepted story is unreleased against a threshold of 1. That story is S-0592, at 40 USD/week against the 1000 USD/week value threshold.
