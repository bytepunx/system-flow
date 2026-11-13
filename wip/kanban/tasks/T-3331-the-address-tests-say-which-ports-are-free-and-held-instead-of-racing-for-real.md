---
id: T-3331
type: task
nature: remediation
title: The address tests say which ports are free and held instead of racing for real ones
status: done
parent: S-0779
owner: alex
created: 2026-11-13T20:33:43Z
updated: 2026-11-13T20:48:04Z
transitions:
  - to: ready
    at: 2026-11-13T20:44:27Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T20:44:27Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T20:48:04Z
    by: agent-S-0779
stream: S-0779
tags: [flai]
touches: [flai/cmd/serve_mcp_test.go]
after: [T-3330]
usage:
  source: log
  seconds: 216
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 124
      cache_read: 1271604
      cache_write: 12547
      cost: 0.482
---
# T-3331 The address tests say which ports are free and held instead of racing for real ones

## Work

Every test in `flai/cmd/serve_mcp_test.go` that takes a port from `unusedPort`, directly or through `mcpHostFor`, counts on that port staying free after it closed it. Give each the free-port check T-3330 adds as a value: a fake that says which addresses are held, so that no test races another for a real port. They are `TestMCPAddrsForAProjectFlaiServeStarts`, `TestTheHostGivesTwoProjectsThatRememberOnePortOneEach`, `TestMCPAddrsForListensOnTheBindAddress`, `TestTheHostStartsEachProjectsServerOnTheBindAddress`, the three that call `mcpHostFor`, and `TestTheHostMovesAProjectOffAPortAnotherProcessHolds`. The default range they probe through `mcpAddrsFor` goes through the fake too.

`TestTheHostMovesAProjectOffAPortAnotherProcessHolds` keeps showing what I-0271 fixed: once the fake holds the remembered port, the project moves off it, keeps the new one, and the log names both.

Remove `unusedPort`, or correct its comment, which says its port is outside the range the parallel tests probe: that holds for flai's default range and not for the kernel's. Keep one test of `freePort` against a real listener the test holds, which is false for as long as it holds it, so that what `freePort` does with the operating system stays covered.

Waits for T-3330, which adds the value these tests pass, and changes the same file.

## Done when

- No test in `flai/cmd/serve_mcp_test.go` counts on a port it closed, or never held, being free when flai probes or listens on it.
- `TestTheHostMovesAProjectOffAPortAnotherProcessHolds` shows the move off a held port, the new port kept, and the log line.
- A test covers `freePort` against a real listener it holds.
- `flai test` on `flai/cmd/serve_mcp_test.go` passes.

## Notes
