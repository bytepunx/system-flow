---
id: D-0328
title: The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close
kind: cause
sources: [I-0157]
topics: [flai, cli, go, code, all]
paths: [design/adrs/README.md, flai/internal/buildinfo/buildinfo.go, flai/internal/buildinfo/buildinfo_test.go, flai/cmd/version.go, flai/cmd/version_test.go, flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/update.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/adrs/0180-flai-task-done-refuses-an-older-flai-only-when-the-newest-flai-release-raised.md, flai/internal/buildinfo/task-done-rules.txt, docs/users/flai-reference.md, flai/cmd/task_done.go]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0328 The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close

I-0157 was closed on 2026-10-20.

The cause:

```text
The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close
```

The fix:

```text
S-0445: flai task done refuses an older flai only when the newest flai release raised the task-done rules number (ADR-0180), kept in flai/internal/buildinfo/task-done-rules.txt, embedded as buildinfo.TaskDoneRules, printed by flai version --json as task_done_rules, and read at each flai tag's commit; a tag without it keeps ADR-0142's minor rule, and the MCP task_done's installed-flai fallback catches up by the same number. A minor release that keeps the number, as in every instance of I-0157, no longer stops a close: TestRunClosesAMinorVersionBehindAReleaseThatKeptTheRules (taskdone) and TestTaskDoneClosesAMinorVersionBehindAReleaseThatKeptTheRules (mcpserver) reproduce it.
```

Source: I-0157.
