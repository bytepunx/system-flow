---
id: D-0516
title: "A flai/cmd test that expects an agent start to be refused names harness: command with no command set, since an agent naming no harness now falls back to claude-code and would exec the host's real claude"
kind: decision
sources: [S-0587, T-2582]
topics: [flai, cli, go, testing]
paths: [flai/cmd/serve_actions_test.go, flai/internal/serve/harness.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0516 A flai/cmd test that expects an agent start to be refused names harness: command with no command set, since an agent naming no harness now falls back to claude-code and would exec the host's real claude

S-0587 made an agent that names no harness fall back, after the project's default and the host's command, to flai's built-in default, `claude-code` ([ADR-0261](../adrs/0261-a-story-or-strategic-agent-that-names-no-harness-falls-back-after-the-project-s.md)). Before that, a story with no harness was refused, and tests in `flai/cmd` relied on that refusal.

T-2582's sub-agent found that those tests "would now exec the host's real `claude`". S-0587 widened T-2582, "and its tests now name `harness: command`, which refuses without a command" (wip/archive/agents/S-0587.md, Decisions).

So a later test that needs a start to be refused names `harness: command` and sets no command. An empty agent block no longer refuses: it starts the host's `claude`.

Sources: S-0587's narrative, § Decisions; T-2582.
