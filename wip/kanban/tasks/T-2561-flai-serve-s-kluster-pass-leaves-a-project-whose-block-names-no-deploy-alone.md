---
id: T-2561
type: task
nature: improvement
title: flai serve's kluster pass leaves a project whose block names no deploy alone
status: in-progress
parent: S-0586
owner: alex
created: 2026-11-07T08:30:03Z
updated: 2026-11-07T08:39:05Z
transitions:
  - to: ready
    at: 2026-11-07T08:39:02Z
    by: agent-S-0586
  - to: in-progress
    at: 2026-11-07T08:39:05Z
    by: agent-S-0586
stream: S-0586
tags: [kluster, cli]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go]
---
# T-2561 flai serve's kluster pass leaves a project whose block names no deploy alone

## Work

Waits for T-2559, whose ADR decides what a block with no `deploy` means. It shares no path with the tasks of its layer.

With the scaffold's block and the action on, the kluster pass in `flai/internal/serve/kluster.go` would try a deploy at every story's head, bring the cluster up for it, fail with the refusal for no `deploy`, and open a thread per story. Make the pass queue no deploy and no promotion while `kluster.deploy` is empty. It asks at every look, as it asks for the block and the action today, so a `deploy` the operator sets is picked up at the next look.

A story already in review then gets its deploy at that look, because its head is neither on a record nor tried. A teardown of a namespace with a record still runs. Bring-up stays what the operator asks for from the page.

## Done when

- A test with a block that names no `deploy` and the action on finds, with a story in review and a story accepted, no job queued, no `kluster up` run, and no thread opened.
- A test that then sets `deploy` finds the story in review deployed at the next look.
- A namespace with a record whose story left review is still torn down with no `deploy` set.
- `flai test` passes on the paths this task changes.

## Notes
