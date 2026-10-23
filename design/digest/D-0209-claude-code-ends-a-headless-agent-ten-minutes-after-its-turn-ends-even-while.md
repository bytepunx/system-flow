---
id: D-0209
title: Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent
kind: cause
sources: [I-0084]
topics: [cli, conventions, flai, go, template, code, all]
paths: [flai/internal/harness, flai/internal/mcpserver, flai/internal/guard, flai/cmd/guard.go, design/conventions/delegation.md, template, ".claude/settings.json", design/system/flai-cli.md, design/system/agent-context.md, docs/users, docs/operators, design/issues, design/adrs, flai/cmd/guard_test.go]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0209 Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent

I-0084 was closed on 2026-10-07.

The cause:

```text
flai serve runs a story's agent headless, `claude -p`. When the agent ends a turn while a sub-agent it launched in the background is still running, Claude Code 2.1.290 keeps the process for ten minutes. A sub-agent that finishes in that time begins the agent's next turn with its notice. One that does not is ended with the process, mid-work, and the session is over. In the one case timed to the second, the process ended 10m00s after the turn did.

flai serve then records the run as failed and does nothing more. The story stays in progress with no agent, and every ready story it holds by overlap stays held, until someone restarts the agent: `flai serve agent restart`, or the story page's Retry button. On 2026-10-06 that was ten ready stories behind S-0220.

Ending the turn is what the start prompt's sentence on sub-agents leads to when it is followed: "Wait for a sub-agent run in the background through the harness's notice that it has finished". The other habit, holding `wait_for_events`, keeps the process alive and wastes the wait instead ([I-0083](I-0083-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md)). The agent has no way to wait for a long sub-agent that is both safe and prompt.
```

The fix:

```text
Fixed in two parts. S-0285 (ADR-0092) gave a story's agent a safe way to wait: it launches sub-agents in the foreground and never ends its turn on a background one. S-0294 (ADR-0108) has flai serve restart a story's agent that ended with its story in progress, up to the host's agent.auto_restarts (2 by default), and then open one thread to the operator.
```

Source: I-0084.
