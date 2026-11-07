---
id: D-0517
title: "flai does not refuse a harness whose program is not on PATH: S-0587 dropped that check, since the exec failure already reports it and a lookup would differ per harness kind"
kind: decision
sources: [S-0587]
topics: [flai, cli, go]
paths: [flai/internal/serve/harness.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0517 flai does not refuse a harness whose program is not on PATH: S-0587 dropped that check, since the exec failure already reports it and a lookup would differ per harness kind

S-0587's planner proposed refusing a harness whose program is not on `PATH`. The story's agent dropped it when it reviewed the tasks:

"I dropped the planner's idea of refusing a harness whose program is not on `PATH`. The criteria do not ask for it, it needs a lookup per harness kind (the container harness runs `docker`), and an exec failure already reports it." (wip/archive/agents/S-0587.md, Decisions)

[ADR-0261](../adrs/0261-a-story-or-strategic-agent-that-names-no-harness-falls-back-after-the-project-s.md) says only that a new project starts an agent "provided `claude` is on the host's `PATH` or the claude-code harness's program is set". It does not say that flai checks `PATH` before a start, and flai does not. A later story that wants that check has to look up each harness kind's own program.

Source: S-0587's narrative, § Decisions.
