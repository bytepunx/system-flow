---
id: T-3642
type: task
nature: improvement
title: Add the verified event and the passed-in-review condition to the role trigger vocabulary
status: backlog
parent: S-0854
owner: alex
created: 2026-11-27T23:59:30Z
updated: 2026-11-27T23:59:30Z
transitions: []
stream: S-0854
tags: [flai]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go]
after: [T-3641]
---
# T-3642 Add the verified event and the passed-in-review condition to the role trigger vocabulary

## Work

This task waits for T-3641. It builds the event and the condition that T-3641's ADR names, with the names and rules the ADR settles.

In `flai/internal/roles/roles.go`, add `verified` to `Events`, so that a role file may name it. In `flai/internal/roles/conditions.go`:

- Give `Occurrence` the fields the condition reads: the verdict, whether the verdict came from the record, and the story's state when flai serve looked.
- Add `passed-in-review` to the `Conditions` table, for the `verified` event alone. Make `Trigger.Holds` match it only when the verdict is `passed`, the verdict did not come from the record, and the story is in review.
- Refuse the condition on any other event, as the table refuses `to` on an event other than `moved`.

Test it in `flai/internal/roles/conditions_test.go` and `flai/internal/roles/roles_test.go`. A role file naming `{event: verified, when: passed-in-review}` parses and validates. The condition holds for a passed job on a story in review. It does not hold for a verdict from the record, a verdict of `failed`, `stopped`, or `interrupted`, or a story that is in progress or done. It is refused on `moved`.

## Done when

- A role file may name the `verified` event and the `passed-in-review` condition, and validation refuses the condition on any other event.
- The tests above pass, and the existing roles tests pass unchanged.
- `flai test` on the changed paths passes.

## Notes
