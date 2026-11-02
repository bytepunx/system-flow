---
id: D-0468
title: flai guard refuses a story's sub-agent a write to scripts/agent-container.sh, scripts/agent-image.sh, or flai/images/ while auto-approve is off, so the story's agent writes those files from the content its sub-agent returns
kind: decision
sources: [S-0517]
topics: [harness, docker, planning]
paths: [scripts/agent-container.sh, scripts/agent-image.sh, flai/images, flai/internal/protected/protected.go, flai/internal/guard/call.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0468 flai guard refuses a story's sub-agent a write to scripts/agent-container.sh, scripts/agent-image.sh, or flai/images/ while auto-approve is off, so the story's agent writes those files from the content its sub-agent returns

S-0508 added the container harness's files to the paths flai protects: `flai/images/`, `scripts/agent-image.sh`, and `scripts/agent-container.sh`, matched from the repository's root (`flai/internal/protected/protected.go`). Since ADR-0102 and ADR-0106, `flai guard` refuses a story's sub-agent an `Edit` or `Write` of any protected path while auto-approve is off, and tells it to put the file's whole new content in its final message for the story's agent to write.

S-0517 met this on T-2220, from its Decisions (`wip/archive/agents/S-0517.md`):

```text
The T-2220 script edits were mine: the guard protects `scripts/agent-container.sh` from a sub-agent, which returned the whole new file.
```

For a later story that changes these files (S-0510, S-0518, S-0520): plan the edit as the story's agent's, or expect the task's sub-agent to hand back the file, rather than leaving the write to a sub-agent's layer. A story whose branch changes a protected path is also accepted by the operator only (ADR-0106).

Source: S-0517.
