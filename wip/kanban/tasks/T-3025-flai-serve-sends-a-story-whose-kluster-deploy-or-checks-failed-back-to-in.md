---
id: T-3025
type: task
nature: remediation
title: flai serve sends a story whose kluster deploy or checks failed back to in-progress, starts its agent, and keeps its namespace
status: in-progress
parent: S-0694
owner: alex
created: 2026-11-10T10:32:46Z
updated: 2026-11-11T09:13:58Z
transitions:
  - to: ready
    at: 2026-11-11T09:13:58Z
    by: agent-S-0694
  - to: in-progress
    at: 2026-11-11T09:13:58Z
    by: agent-S-0694
stream: S-0694
tags: [flai, kluster, serve]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go]
after: [T-3024]
---
# T-3025 flai serve sends a story whose kluster deploy or checks failed back to in-progress, starts its agent, and keeps its namespace

## Work

This is ADR-0254 §5a's first bullet. Today `deploy()` and `checks()` in `flai/internal/serve/kluster.go` call `report()`, which opens one thread and moves nothing. `klusterFailureText`, `checksFailureText`, and the file's header comment say so.

- After the thread opens, move the story from review back to in-progress as `flai`, with a reason that names the job, the step, and the commit. Use `repo.Transition`, as `forgePass.sendBack()` in `forgestatus.go` does for ADR-0167. Do not move it when the action went off meanwhile, or when the story is no longer in review at that head.
- Start the story's agent on the thread, as `startAcceptSentBack()` in `agents_restart.go` does today. That function finds threads by their opener and title prefix, and the kluster thread's opener is `flai serve`. Either match that opener, or give `harness.SentBack` a kluster flag of its own. Name the choice in the narrative.
- Resolve the thread once the story is back in review and its new deploy passes. The pass already resolves the thread when a deploy passes; keep that behaviour.
- Keep the namespace. `lookLeft()` and `teardown()` tear down a story that left review. Skip that while the story is in-progress after a kluster send-back. Still tear down when it is cancelled or accepted.
- When the story returns to review at a new head, deploy it again in place, as a new head does already.
- Rewrite the failure texts and the header comment to say the story was sent back.

This task waits for T-3024: both change `kluster.go`, and the send-back reads the checks outcome T-3024 records.

## Done when

- `TestAFailedDeployOpensOneThreadMovesNothingAndAPassResolvesIt` becomes a test that a failed deploy opens one thread, moves the story to in-progress as `flai`, and starts its agent once.
- A failed checks run against the namespace does the same.
- A test shows that a sent-back story keeps its namespace, and that cancelling it then tears the namespace down.
- A test shows that the story's return to review deploys again in place, and that a pass resolves the thread.
- `flai test` passes on the five files.

## Notes
