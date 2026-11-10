---
id: D-0468
title: flai guard refuses a story's sub-agent a write to scripts/agent-container.sh, scripts/agent-image.sh, or flai/images/ while auto-approve is off, so the story's agent writes those files from the content its sub-agent returns
kind: decision
sources: [S-0517, S-0510, S-0520, S-0670, S-0672]
topics: [harness, docker, planning]
paths: [scripts/agent-container.sh, scripts/agent-image.sh, flai/images, flai/internal/protected/protected.go, flai/internal/guard/call.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: [S-0520, S-0670, S-0682, S-0713]
---

# D-0468 flai guard refuses a story's sub-agent a write to scripts/agent-container.sh, scripts/agent-image.sh, or flai/images/ while auto-approve is off, so the story's agent writes those files from the content its sub-agent returns

S-0508 added the container harness's files to the paths flai protects: `flai/images/`, `scripts/agent-image.sh`, and `scripts/agent-container.sh`, matched from the repository's root (`flai/internal/protected/protected.go`). Since ADR-0102 and ADR-0106, `flai guard` refuses a story's sub-agent an `Edit` or `Write` of any protected path while auto-approve is off, and tells it to put the file's whole new content in its final message for the story's agent to write.

S-0517 met this on T-2220, from its Decisions (`wip/archive/agents/S-0517.md`):

```text
The T-2220 script edits were mine: the guard protects `scripts/agent-container.sh` from a sub-agent, which returned the whole new file.
```

For a later story that changes these files (S-0510, S-0518, S-0520, S-0670, S-0672): plan the edit as the story's agent's, or expect the task's sub-agent to hand back the file, rather than leaving the write to a sub-agent's layer. A story whose branch changes a protected path is also accepted by the operator only (ADR-0106).

S-0510 planned around it from the start. From its Decisions (`wip/archive/agents/S-0510.md`):

```text
T-2165 changes `flai/images/` and `scripts/agent-image.sh`, which the protected list holds since S-0508. So the story's agent makes those writes itself, through `permission_prompt`, rather than through a sub-agent.
```

S-0520 did the same. From its Decisions (`wip/archive/agents/S-0520.md`):

```text
T-2232 and T-2235 are mine, not a sub-agent's: the ADR's context is already in this session, and `scripts/agent-container.sh` is a path the guard protects from a sub-agent (D-0468).
```

So did S-0670, on T-2908, from its Decisions (`wip/archive/agents/S-0670.md`):

```text
T-2908 is mine, not a sub-agent's: both scripts are on the protected list (D-0468).
```

S-0672, which changed `flai/images/system-flow/Dockerfile` and `scripts/agent-image.sh`, worked all four of its tasks with no task sub-agents, for the same reason. From its Decisions (`wip/archive/agents/S-0672.md`):

```text
I worked every task myself, with no task sub-agents. `flai/images/` and `scripts/agent-image.sh` are protected, so a sub-agent may not write them.
```

Sources: S-0517, S-0510, S-0520, S-0670, S-0672.
