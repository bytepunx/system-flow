---
id: D-0302
title: A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other
kind: cause
sources: [I-0125]
topics: [cli, conventions, flai, go, template, code, all]
paths: [design/adrs/README.md, design/system/agent-coordination.md, design/system/workflow.md, design/system/agent-narrative.md, design/system/flai-cli.md, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/folder.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/summary.md, design/adrs/0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, docs/operators/settings.md, design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md#decision"
asked_again: []
cited: []
---

# D-0302 A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other

I-0125 was closed on 2026-10-14.

The cause:

```text
A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other
```

The fix:

```text
S-0372, ADR-0154. A reply that asks nothing is sent settled (flai message reply --settled, message_reply settled), and the conversation then awaits neither story, so an agreement no longer wakes the other story's agent (instances of MS-0021, MS-0041, MS-0101, MS-0109, MS-0110; TestASettledReplyAwaitsNeitherStory, TestASettledReplyDoesNotStartAnAgentThatEndedForAThread, TestASettledReplyStartsTheAgentWaitingOnItOnce, TestWaitForEventsCountsNoSettledConversationAsAwaitingAReply). flai task done leaves the paths sync regenerates, such as design/issues/summary.md, out of its notice (MS-0053; TestRunTellsNoOneOfThePathsSyncRegenerates). While a story is blocked, flai serve starts its agent again only for the answer to its thread (S-0337's 12 restarts; TestABlockedStorysAgentIsStartedOnlyForTheAnswerToItsQuestion). The work-management convention says to reply settled when a change does not break your work.
```

Source: I-0125.
