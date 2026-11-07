---
id: D-0518
title: claudeCodeProject counts an agent that names no harness as claude-code even when the host sets a command, an over-approximation S-0587's review kept since it costs only a version check
kind: decision
sources: [S-0587]
topics: [flai, go, code]
paths: [flai/internal/serve/claudecheck.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0518 claudeCodeProject counts an agent that names no harness as claude-code even when the host sets a command, an over-approximation S-0587's review kept since it costs only a version check

S-0587 made an agent that names no harness resolve to `claude-code` after the host's command ([ADR-0261](../adrs/0261-a-story-or-strategic-agent-that-names-no-harness-falls-back-after-the-project-s.md)). It changed `claudeCodeProject` in `flai/internal/serve/claudecheck.go` to count such an agent as claude-code, so flai checks the installed Claude Code for that project.

When the host does set a command, that agent actually runs the command, not claude-code. The verifier flagged it, and the story kept it: "`claudeCodeProject` counts an agent with no harness as claude-code even when a command is set. That over-approximation costs only a version check." (wip/archive/agents/S-0587.md, Decisions, Verifier review)

A later review that finds this behaviour finds a known, accepted trade-off, not a defect.

Source: S-0587's narrative, § Decisions.
