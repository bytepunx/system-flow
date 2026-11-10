---
id: TH-1040
title: "Plan for S-0686: reproduce I-0255 with git's own message, guard it in the smoke tier, drop GOFLAGS=-buildvcs=false, four tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0686-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md
  item: S-0686
status: resolved
participants: [planner-S-0686, orchestrator]
created: 2026-11-10T10:35:14Z
updated: 2026-11-10T11:10:08Z
---

# TH-1040 Plan for S-0686: reproduce I-0255 with git's own message, guard it in the smoke tier, drop GOFLAGS=-buildvcs=false, four tasks in four layers

On wip/kanban/stories/S-0686-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md.

## Entries

### 2026-11-10T10:35:14Z planner-S-0686
S-0686 is planned: 4 tasks in 4 layers, 9 file touches, forecast 1h30m, cost of delay 35 USD a week.

## Tasks

| Layer | Task | Nature | Touches |
|-------|------|--------|---------|
| 1 | T-3032 The cause of I-0255's exit status 128 is found in the layer's image and a fix is proposed | research | the I-0255 issue file |
| 2 | T-3036 The smoke tier's agent-container test runs a VCS-stamped go build in the layer's image (after T-3032) | remediation | `scripts/agent-container-test.sh` |
| 3 | T-3039 A story's container stamps go builds again, and the layer drops GOFLAGS=-buildvcs=false (after T-3036) | remediation | `flai/images/system-flow/Dockerfile`, `scripts/agent-container.sh`, `flai/internal/harness/container.go`, `container_test.go` |
| 4 | T-3041 The runbook and docker.md say how a story's container stamps go builds, and I-0255 is closed (after T-3039) | remediation | `docs/operators/runbooks/agent-container.md`, `design/tech/docker.md`, the I-0255 issue file, `design/issues/summary.md` |

The chain is strict. T-3036's case runs the command and mounts T-3032 found. T-3039 is checked by T-3036's case. T-3041 describes T-3039's fix and changes I-0255 after T-3032.

## Assumptions

- The orchestrator's request named a thread `asked`. No thread of that name exists, and none is on S-0686 or I-0255. So I planned the story from I-0255 alone, as a first plan. If you meant another thread, name it and I will plan again.
- The cause may be gone. I-0255's one instance, S-0510 on 2026-11-02, came before ADR-0253's main checkout mounts and `safe.directory` (2026-11-06). If T-3032 finds it no longer reproduces, T-3036's case guards against a return, T-3039 only drops the flag, and the issue is closed with the change that fixed it.
- If the cause lies in Go or git and cannot be fixed here, the flag stays with its cause written in the Dockerfile's comment and in the runbook. I-0255 is closed on that record.
- If the fix changes the container's mounts or what a story's agent may write, the story's agent asks on this story before building it, since that is a security boundary (ADR-0253).
- `scripts/agent-image.sh` (`--check` sets the same `safe.directory`) is not in the touches. The story's agent widens them if the fix reaches it.
- Forecast raised from flai's 25m to 1h30m for the live container runs and the layer rebuild. Cost of delay kept at flai's 35 USD a week.
- No `after`. S-0667 also changes `scripts/agent-container.sh` and the runbook, and S-0713 the runbook. A trial merge judges that.

Nothing waits on you but finalizing the draft.

### 2026-11-10T11:09:35Z orchestrator
Approved.

- The strict chain fits a cause that is not yet known: reproduce it with git's own message (T-3032), guard it in the smoke tier (T-3036), fix it and drop the flag (T-3039), then document and close (T-3041).
- Each case is covered: the cause gone, the cause found, and a cause in Go or git.
- A fix that changes the container's mounts or what a story's agent may write is asked on S-0686 first (ADR-0253).
- There are nine file touches. The 1h30m forecast states its reason. The value is 35 USD a week, from flai's 14m input. There is no `after`.
- The `asked` thread did not exist, so a first plan from I-0255 was right.

I am finalizing S-0686.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:10:08Z orchestrator
Resolved: Plan approved; S-0686 finalized
