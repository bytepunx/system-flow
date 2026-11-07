---
id: D-0484
title: TestHostActionPush in flai/cmd/serve_actions_test.go pins the name of every host action, so a story that adds one to hostapi.Actions changes that test too
kind: decision
sources: [S-0577, T-2473]
topics: [host-actions, flai, go, code]
paths: [flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0484 TestHostActionPush in flai/cmd/serve_actions_test.go pins the name of every host action, so a story that adds one to hostapi.Actions changes that test too

S-0577 added the `kluster` host action to `hostapi.Actions` (T-2473). `TestHostActionPush` in `flai/cmd/serve_actions_test.go` checks the error `flai serve enable` gives for an unknown action, and that error lists every host action by name, in alphabetical order. So the new action broke the test, and T-2473 widened its touches to the test file.

From S-0577's narrative, `## Decisions`:

```text
T-2473's new action changed the action list `TestHostActionPush` pins in `flai/cmd/serve_actions_test.go`, so T-2473 took that path too.
```

A story that adds a host action puts the name in that list, in its alphabetical place, and claims `flai/cmd/serve_actions_test.go` in its touches. D-0460 names the other place nothing checks: the `host_actions` row of `docs/operators/settings.md`.

Sources: S-0577 (wip/archive/agents/S-0577.md), T-2473.
