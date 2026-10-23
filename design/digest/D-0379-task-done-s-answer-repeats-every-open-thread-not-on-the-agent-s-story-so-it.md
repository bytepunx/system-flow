---
id: D-0379
title: task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result
kind: cause
sources: [I-0154]
topics: [flai, cli, go, code, all]
paths: [flai/internal/mcpserver/task_test.go, flai/internal/taskdone/taskdone_test.go, design/issues/I-0154-task-done-s-answer-repeats-every-open-thread-not-on-the-agent-s-story-so-it-overflows-a-tool-result.md, design/issues/summary.md]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0379 task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result

I-0154 was closed on 2026-10-23.

The cause:

```text
task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result
```

The fix:

```text
S-0470: the cause was removed by S-0450 (ADR-0179, flai v1.67.1), which made task_done and story_start list only the story's threads, count the rest in threads_omitted, and fit each answer to 40,000 bytes, and by S-0453 (ADR-0181, flai v1.67.6), which made a thread on another item the agent never wrote to read awaiting operator, neither listed nor counted. TestTaskDoneCarriesNoneOfTheOperatorsThreads and TestRunCarriesNoneOfTheOperatorsThreads reproduce I-0154's instance and name it.
```

Source: I-0154.
