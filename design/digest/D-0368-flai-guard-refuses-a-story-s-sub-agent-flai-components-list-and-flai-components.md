---
id: D-0368
title: flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent
kind: cause
sources: [I-0165]
topics: [cli, components, flai, go, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/agent-adapters.md, design/system/component-library.md, design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/issues/summary.md, design/adrs/0186-flai-guard-lets-every-agent-run-flai-components-list-and-a-story-s-task-sub.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/guard/call.go]
date: 2026-10-22
status: absorbed
absorbed_by: design/system/component-library.md
asked_again: []
cited: []
---

# D-0368 flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent

I-0165 was closed on 2026-10-22.

The cause:

```text
flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent
```

The fix:

```text
S-0477 (ADR-0186): flai guard lets every agent run flai components list, and a story's task sub-agent run flai components serve, with or without --story; a read-only sub-agent and the strategic agents' own calls are refused serve, and config stays refused. TestASubAgentListsComponentsAndATaskSubAgentServesThem and TestTheStrategicAgentsListComponentsButDoNotServeThem in flai/internal/guard, and TestGuardLetsATaskSubAgentListAndServeComponents in flai/cmd, reproduce the refusal against main's guard.
```

Source: I-0165.
