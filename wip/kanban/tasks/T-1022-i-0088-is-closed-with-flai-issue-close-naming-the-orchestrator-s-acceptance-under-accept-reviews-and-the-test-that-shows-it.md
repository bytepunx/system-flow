---
id: T-1022
type: task
nature: remediation
title: I-0088 is closed with flai issue close, naming the orchestrator's acceptance under accept_reviews and the test that shows it
status: backlog
parent: S-0296
owner: alex
created: 2026-10-06T11:50:17Z
updated: 2026-10-06T11:50:17Z
transitions: []
stream: S-0296
tags: [flai]
touches: [design/issues/I-0088-a-story-waits-in-review-for-the-operator-while-nothing-else-can-start-so-the-board-stands-idle-until-a-person-accepts.md, design/issues/summary.md]
after: [T-1020, T-1021]
---
# T-1022 I-0088 is closed with flai issue close, naming the orchestrator's acceptance under accept_reviews and the test that shows it

## Work

In the story's worktree, run `flai issue close I-0088 --reason` with a reason that names what fixed it: S-0221's acceptance by the orchestrator under `orchestration.permissions.accept_reviews`, S-0286's operator-only acceptance of a story that changes a protected path, T-1020's test, and the design and guides T-1021 wrote. Commit it on the story branch.

It waits for T-1020 and T-1021 because the reason cites what they deliver, and the issue is closed only once the fix is shown. If the operator has not turned `accept_reviews` on for this project by then, say so in the reason, as the plan's thread records their answer.

## Done when

- I-0088's status is closed, its reason as above, and `design/issues/summary.md` agrees.
- `flai check --strict` passes.

## Notes
