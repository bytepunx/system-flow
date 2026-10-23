---
id: D-0327
title: task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator
kind: cause
sources: [I-0167]
topics: [mcp, inbox, flai, cli, go, template, conventions, code, all]
paths: [design/adrs/README.md, design/system/agent-narrative.md, design/system/flai-cli.md, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/cmd/story_start.go, flai/cmd/story_start_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/session-start.md, template/root/design/conventions/session-start.md, design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md, design/issues/summary.md, design/adrs/0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0327 task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator

I-0167 was closed on 2026-10-20.

The cause:

```text
task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator
```

The fix:

```text
Fixed by S-0450 (ADR-0179): story_start and task_done, over MCP and with --json, now carry an inbox that lists only the threads on the story and its tasks and counts the rest in threads_omitted, and each whole answer is fitted to 40,000 bytes of JSON, cutting the oldest changes, then the pack to its header, then the listed threads' last entries. Tests reproduce S-0398's instance (33 threads elsewhere, task_done) and S-0399's (32 threads elsewhere, story_start). S-0450's own task_done answered 54,219 bytes before and 2,873 after, with 33 threads omitted.
```

Source: I-0167.
