---
id: T-2499
type: task
nature: feature
title: "The planner's container run is watched to its end: it plans the chosen story through flai, and wip/agents/planner.md gains its entry with duration and cost from the container's log"
status: backlog
parent: S-0567
owner: alex
created: 2026-11-07T03:43:22Z
updated: 2026-11-07T03:43:22Z
transitions: []
stream: S-0567
tags: [docker]
after: [T-2498]
---
# T-2499 The planner's container run is watched to its end: it plans the chosen story through flai, and wip/agents/planner.md gains its entry with duration and cost from the container's log

## Work

Waits for T-2498: the run starts only once the operator has set the harnesses and asked for the plan.

Hold `wait_for_events` until the planner's run on the chosen story reports `started`, then `ended` or `failed` (kind `agent`, role `plan`). While it runs, record with `docker ps` the container's name (`flai-agent-plan-<session>`) and its image, and with `docker inspect` its mounts against ADR-0247's four: the main checkout read-write, `.git/hooks` read-only, `.flai-cache/agents/plan` read-write, and the login read-only at `/run/claude-login`. Nothing else may be mounted.

Once it ends, read the chosen story with `item_get` and its plan's thread with `thread_get`: the drafted tasks and the thread must be written through flai (criterion 1). Read the new entry at the end of `wip/agents/planner.md`: its duration and cost, and whether the cost is priced as a host run's is, against the entries of host runs before it (criterion 2). Read the serve journal (`flai serve journal`) for the run's start and end.

Record each finding, and each limitation met, in S-0567's Notes under `### Live run`, with `flai edit S-0567`. When the run fails, record why, and ask the operator on S-0567's thread whether to run it again after a fix.

## Done when

- The container's name, image, and mounts are recorded, and the mounts are ADR-0247's four, or each difference is listed as a limitation.
- The chosen story has the planner's tasks and plan thread, written through flai.
- `wip/agents/planner.md` has the run's entry with a duration and a cost, and its pricing is compared with host runs in the Notes.

## Notes
