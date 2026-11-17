---
id: T-3585
type: task
nature: improvement
title: Ask again once, at serve's start, a close-out or rebase and resolve job a restart left interrupted
status: backlog
parent: S-0846
owner: alex
created: 2026-11-17T12:33:39Z
updated: 2026-11-17T12:33:39Z
transitions: []
stream: S-0846
tags: [serve, code]
touches: [flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go]
after: [T-3584]
---
# T-3585 Ask again once, at serve's start, a close-out or rebase and resolve job a restart left interrupted

## Work

Covers the restart the drain cannot hold off: the bound passed, an operator's `flai host upgrade`, or a crash. It waits for T-3584 because both change `verifyjob.go`, and for the ADR of T-3583 through it.

- In `verifyPass.sweep` (`verifyjob.go`), after it ends a job left running `interrupted` and announces it, write a request for the same story, kind, message, `all`, and `fresh`, marked as asked again by `flai serve`, unless the interrupted job was itself asked again. The next `look` takes it as any request; one the agent asks meanwhile is answered with that job's id, as for any running job.
- For a rebase and resolve (`resolvejob.go`), ask again only while the story is in review; otherwise end it interrupted as today.
- Journal the ask again, and say in the `verified` event's detail that the job was asked again, with the new request.
- Tests: an interrupted close-out is asked again once and runs; a job that was asked again and interrupted again is not asked a third time; an interrupted rebase and resolve whose story left review is not asked again.

## Done when

- The three tests above pass in `verifyjob_test.go` and `resolvejob_test.go`, with `TestAJobLeftRunningIsEndedInterruptedWhenServeStarts` changed to expect the request.
- `flai test` passes on the changed paths.

## Notes
