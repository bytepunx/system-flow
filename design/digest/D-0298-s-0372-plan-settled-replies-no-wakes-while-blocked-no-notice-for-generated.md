---
id: D-0298
title: "S-0372 plan: settled replies, no wakes while blocked, no notice for generated paths (new ADR), seven tasks in five layers: Recommendation: approve the plan as drafted, including the new ADR refining…"
kind: answer
sources: [TH-0422]
topics: [cli, conventions, flai, go, template, code, all]
paths: [design/adrs/README.md, design/system/agent-coordination.md, design/system/workflow.md, design/system/agent-narrative.md, design/system/flai-cli.md, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/folder.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/summary.md, design/adrs/0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, docs/operators/settings.md, design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md#decision"
asked_again: []
cited: []
---

# D-0298 S-0372 plan: settled replies, no wakes while blocked, no notice for generated paths (new ADR), seven tasks in five layers: Recommendation: approve the plan as drafted, including the new ADR refining…

The operator, alex, answered TH-0422 at 2026-10-14T20:56:22Z.

What was asked:

```text
**Recommendation:** approve the remedy below and the change to the baseline of `work-management.md` that it brings. S-0372 stays a draft until it is finalized.
```

The answer:

```text
**Recommendation:** approve the plan as drafted, including the new ADR refining ADR-0120 and ADR-0108, and the baseline change to `work-management.md` in both copies.

## Why

- **The remedy matches the instances.** I-0125 records settled conversations that keep waking each agent in turn. A settled mark that awaits neither story stops that and keeps the conversation open for the next real message. It still wakes an agent whose run ended waiting on that very reply, so nothing that needs an answer is lost.
- **The baseline rule causes the loop.** § Rules says "Reply on the conversation either way, since it awaits your story until you do". That reply is what wakes the other side. "Reply settled when the change does not break your work" keeps the duty to answer and ends the ping-pong.
- **No wakes while blocked** matches the blocked rule in the same section: a blocked story waits on the operator's answer, not on its peers.
- **Figures.** Every touch names a file. Raising flai's 34m forecast to 1h has a stated reason: seven tasks, an ADR, and a convention change. The 90 USD a week is flai cod's, from flai's 36m input.

## Why this is yours

Both points change recorded decisions: ADR-0120's definition of `awaiting`, and a baseline convention rule, which the template releases to every project. I leave S-0372 a draft until you answer. It is held behind S-0367 either way.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/conventions/work-management.md § Rules
```

Source: TH-0422.
