---
id: D-0514
title: OpenCode's files, AGENTS.md, opencode.json, opencode.jsonc, and .opencode/, are protected on every harness, so a story whose branch changes one is accepted by the operator alone, on claude-code too
kind: decision
sources: [S-0406]
topics: [harness, agents, review, planning]
paths: [flai/internal/protected/protected.go, AGENTS.md, opencode.json, opencode.jsonc, ".opencode", template/root/opencode.json, docs/operators/index.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0514 OpenCode's files, AGENTS.md, opencode.json, opencode.jsonc, and .opencode/, are protected on every harness, so a story whose branch changes one is accepted by the operator alone, on claude-code too

S-0406 built OpenCode's protected list from Claude Code's, plus `.opencode/`, `opencode.json`, `opencode.jsonc`, and `AGENTS.md` (`flai/internal/protected/protected.go`). `protected.Path` answers for every list, so these paths are now protected whatever the harness, `claude-code` included.

From S-0406's Decisions (`wip/archive/agents/S-0406.md`): "`protected.Path` answers for every list, so these paths are now protected on claude-code too. That makes S-0406 itself operator-accepted (ADR-0106), since it adds `template/root/opencode.json`."

What follows for a later story: while `auto-approve` is off, a story's agent writes one of these files only with the operator's approval on a thread, and the orchestrator cannot accept or enqueue a story whose branch changes one ([ADR-0106](../adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)). Plan for that on the rest of E-0023, such as S-0407's guard plugin under `.opencode/plugins/`. The operator guide says the same in `docs/operators/index.md` under "Protected paths on OpenCode".

See D-0468 for the container harness's protected files, and D-0458 for an earlier story the operator had to accept for a protected path.

Source: S-0406.
