---
id: D-0300
title: flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories
kind: cause
sources: [I-0098]
topics: [planning, flai, cli, go, code, all]
paths: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0098-flai-serve-replans-a-story-on-its-own-agent-s-touches-edit-while-the-agent-works-it-and-the-planner-s-new-touches-hold-other-ready-stories.md, design/issues/summary.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md]
date: 2026-10-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0300 flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories

I-0098 was closed on 2026-10-14.

The cause:

```text
flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories
```

The fix:

```text
Fixed by S-0305: flai serve's replanner gives no trigger for an edit by the story's own agent, as for a planner's, and drops and logs a queued planner whose story is in progress or in review or whose agent runs, whatever queued it. The operator's flai plan is unchanged. Tests in flai/internal/serve/replan_test.go reproduce the instance.
```

Source: I-0098.
