---
id: T-3618
type: task
nature: research
title: Propose in an ADR refining ADR-0155 how the orchestrator keeps a recorded removal of main's lines, behind a permission off by default
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:57:24Z
updated: 2026-11-27T23:57:24Z
transitions: []
stream: S-0852
tags: [flai]
touches: [design/adrs/0322-with-keep-removals-on-the-orchestrator-keeps-a-recorded-removal-of-main-s-lines.md, design/adrs/README.md]
---
# T-3618 Propose in an ADR refining ADR-0155 how the orchestrator keeps a recorded removal of main's lines, behind a permission off by default

## Work

The story's goal asks for a solution proposed from I-0415's instances before it is built. ADR-0155 made `--keep-removal` the operator's alone and rejected letting the orchestrator keep a removal, since it cannot tell an intended removal from a lost change. S-0811's two instances show the case it did not weigh: the story's agent rewrote main's sentence on purpose, said so, and the orchestrator checked it, yet the story waited 148 hours.

Write a new ADR that refines ADR-0155, ADR-0093, and ADR-0278. The planner's recommendation, to confirm or change on the plan's thread before the code tasks start:

- A new permission, `orchestration.permissions.keep_removals`, off by default, needing `accept_reviews`.
- With it on, `flai accept --by orchestrator` takes `--keep-removal <path> --reason "<why>"` for a path the story's record holds. The MCP tool `accept` and `queue_add` take `keep_removals`, each a path and a reason. The queue records them with the entry, leaves a kept path out of `undone_main`, and its landing passes them.
- flai refuses the orchestrator's keep of a path the story's narrative does not name under `## Decisions`. That entry is the story's agent's own record that the removal is a rewrite, not a lost change.
- The kept reason is written under `### Removals kept` in the story's Notes with `by: orchestrator`, as the operator's is, and the orchestrator logs it with `activity_log`.
- With the permission off, nothing changes from ADR-0155.

Name the alternatives weighed: widening ADR-0297's word rule to a near match, which S-0811's sentence fails because it changes words, and which would hide real reverts; and an orchestrator recommendation that `confirm_after` confirms, which still leaves nobody to run the accept. Add the ADR to `design/adrs/README.md`. Take the next free ADR number when you write it; the number in this task's touches is the planner's guess.

This task waits for no other.

## Done when

- The new ADR is written with status `proposed` or `accepted` as the operator's answer on the plan's thread allows, refines ADR-0155, ADR-0093, and ADR-0278, and states the permission, what each of `flai accept`, `accept`, and `queue_add` takes, the narrative condition, and the record.
- `design/adrs/README.md` lists it.
- `flai check --strict` passes on it.

## Notes
