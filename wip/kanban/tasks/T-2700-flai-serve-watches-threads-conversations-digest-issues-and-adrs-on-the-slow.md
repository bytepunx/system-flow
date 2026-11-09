---
id: T-2700
type: task
nature: feature
title: flai serve watches threads, conversations, digest, issues, and ADRs on the slow tier, and a flai issue and ADR write are seen at the next fast look
status: in-progress
parent: S-0606
owner: alex
created: 2026-11-09T07:35:43Z
updated: 2026-11-09T07:48:04Z
transitions:
  - to: ready
    at: 2026-11-09T07:48:03Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T07:48:04Z
    by: system-flow
stream: S-0606
tags: []
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
after: [T-2699]
---

# T-2700 flai serve watches threads, conversations, digest, issues, and ADRs on the slow tier, and a flai issue and ADR write are seen at the next fast look

## Work

TH-0908's answer: `slowPaths` in `flai/internal/serve/serve.go` adds the manifest's `wip` folder's `threads` and `messages`, and its `design` folder's `digest`, `issues`, and `adrs`, beside `wip/archive`. A test in `flai/internal/serve/serve_test.go` checks the paths, under default and renamed layout folders. It also shows that a flai issue write and an ADR write, through `issues` and `adr` with the watcher running on these tiers, are reported at the next fast look and not only at the slow interval.

Waits for T-2699: the test relies on its atomic writes.

## Done when

- [ ] `slowPaths` names the six folders, and the test of the paths passes.
- [ ] A test shows a flai issue write and an ADR write each reported within two fast looks while `SlowEvery` is long.

## Notes
