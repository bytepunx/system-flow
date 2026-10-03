---
title: Agent narratives
updated: 2026-10-03
status: active
---

# agents

One file per active story, named by story ID. Read `index.md` first. The full convention is in [design/system/agent-narrative.md](../../design/system/agent-narrative.md).

`planner.md`, `orchestrator.md`, and `analyzer.md` are the strategic agents' activity documents, one entry per activity with totals in front matter. flai writes them; do not edit them by hand.

Resuming after a crash: read `index.md`, then each active narrative's `## Current state` and `## Next steps`, reconcile with `git status`, append a log entry, continue.
