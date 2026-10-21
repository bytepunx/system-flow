---
id: T-1935
type: task
nature: remediation
title: Document in the user docs and the generated reference what inbox and wait_for_work report as awaiting you, operator, and other
status: backlog
parent: S-0453
owner: alex
created: 2026-10-20T22:21:14Z
updated: 2026-10-20T22:21:14Z
transitions: []
stream: S-0453
tags: [flai]
touches: [docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1930, T-1932, T-1933]
---
# T-1935 Document in the user docs and the generated reference what inbox and wait_for_work report as awaiting you, operator, and other

## Work

- In `docs/users/flai.md`, where it describes `inbox`, `wait_for_work`, and the inbox `story_start` and `task_done` carry, say that a thread awaits an agent only when it takes part in it or it is on the agent's own story, that the others read `awaiting: "operator"` or `"other"`, and that the orchestrator's `inbox` lists those awaiting the operator by default. Link the ADR of T-1927.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`, so it carries the tool descriptions T-1930 changed.

It waits for T-1930, T-1932, and T-1933, so that it documents what they built and the reference is generated from their final text.

## Done when

- `docs/users/flai.md` states the rule.
- `make flai-reference` leaves no diff after it runs, and the markdown lint passes.

## Notes

The operator's view, the dashboard and `flai thread list`, is unchanged, so `docs/operators` needs nothing.
