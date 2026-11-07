---
id: T-2502
type: task
nature: feature
title: "The live planner run is checked: its container, its plan, its activity entry's duration and cost, and a guard refusal inside the container"
status: cancelled
parent: S-0567
owner: alex
created: 2026-11-07T03:44:47Z
updated: 2026-11-07T03:49:22Z
transitions:
  - to: cancelled
    at: 2026-11-07T03:49:22Z
    by: agent-S-0567
stream: S-0567
tags: [docker]
after: [T-2498]
---
# T-2502 The live planner run is checked: its container, its plan, its activity entry's duration and cost, and a guard refusal inside the container

## Work

Once the operator has answered T-2498's thread and the planner's run on the chosen story has ended, check the run against criteria 1 to 3. This task waits for T-2498, since there is no run to check before the operator starts it.

- From `flai serve journal` and `docker ps -a`, show the run started in the project's layer image, `flai-agent-sf-claude-code:local`, as `flai-agent-plan-<session>`.
- Show the chosen story gained the planner's thread and drafted tasks, written through flai (`flai show`, `flai thread list`).
- Show `wip/agents/planner.md` gained the run's activity entry with its duration and cost.
- Find a guard refusal in the run's log under `.flai-cache/agents/plan`; where there is none, run `flai guard` by hand in the image with `FLAI_ROLE=plan` on a call the planner is refused on the host, such as a file edit, and record its verdict.

Record each finding in S-0567's Notes under `### Live run`, with `flai edit S-0567`.

## Done when

- S-0567's Notes name the run's container and image, the plan's thread and tasks, the activity entry's duration and cost, and a guard refusal inside the container.

## Notes
- 2026-11-07T03:49:22Z: moved to cancelled: duplicates T-2499 and T-2500, the planner's split (TH-0838)
