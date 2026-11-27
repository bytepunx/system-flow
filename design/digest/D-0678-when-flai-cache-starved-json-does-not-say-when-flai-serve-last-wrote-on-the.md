---
id: D-0678
title: When .flai-cache/starved.json does not say when flai serve last wrote on the starved board's thread, the coalescing window counts from flai serve's newest entry there, or as passed with none, and a starved test waits out the 15-minute window before it expects a thread
kind: decision
sources: [S-0847, T-3573]
topics: [threads, ready, orchestration, go, testing]
paths: [flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/readyempty_test.go]
date: 2026-11-27
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0678 When .flai-cache/starved.json does not say when flai serve last wrote on the starved board's thread, the coalescing window counts from flai serve's newest entry there, or as passed with none, and a starved test waits out the 15-minute window before it expects a thread

S-0847 built ADR-0321's settle window in `flai/internal/serve/starved.go`. ADR-0321 §4 says only that the window starts again from the look that finds `.flai-cache/starved.json` unreadable. T-3573's sub-agent settled two cases it leaves open, and the story's agent kept both. From S-0847's Decisions (`wip/archive/agents/S-0847.md`):

```text
T-3573's sub-agent decided the following, and I reviewed and kept each choice. When the state does not say when flai serve last wrote, the time comes from flai serve's newest entry on the thread; with no such entry, the window counts as passed. When the state file is lost, the window starts again at that look. Four existing starved tests expected a thread after a single look, so they now wait out the window.
```

So a story that changes the starved board's thread, or writes a test of it, keeps these fallbacks. A test that expects the thread must look again at least 15 minutes after its first starved look, as the four tests in `flai/internal/serve/starved_test.go` and `readyempty_test.go` now do. `TestAFlippingBoardOfI0389OpensOneThreadAndWritesOnceAWindow` reproduces I-0389.

Sources: S-0847 (`wip/archive/agents/S-0847.md`, Decisions); T-3573; ADR-0321.
