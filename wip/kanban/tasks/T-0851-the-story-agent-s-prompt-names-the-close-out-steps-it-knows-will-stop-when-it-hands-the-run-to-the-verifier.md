---
id: T-0851
type: task
nature: improvement
title: The story agent's prompt names the close-out steps it knows will stop when it hands the run to the verifier
status: backlog
parent: S-0266
owner: alex
created: 2026-10-05T00:19:04Z
updated: 2026-10-05T00:19:04Z
transitions: []
stream: S-0266
tags: [harness, flai]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/flai-cli.md]
after: [T-0849]
---
# T-0851 The story agent's prompt names the close-out steps it knows will stop when it hands the run to the verifier

## Work

Change `delegation()` in `flai/internal/harness/harness.go` where the prompt hands the close-out run to a fresh verifier. Add that the agent names, in the verifier's prompt, any step it already knows will stop, and why. One example is `flai check --strict` stopping on findings outside the story (I-0057). Add too that the verifier runs the close-out once and reads its last line, which gives the outcome and the step (T-0849). The verifier then reports the expected stop and goes on to the steps after it, rather than running the close-out again.

Assert the new sentences in the delegation tests in `flai/internal/harness/harness_test.go`. Record the change where `design/system/flai-cli.md` describes the prompt that flai serve builds.

This task waits for T-0849, so that the prompt names the line T-0849 prints. It shares no path with T-0850, so the two can run together.

S-0268 also edits `delegation()`. Whichever of the two is pulled second syncs over the first.

## Done when

- The harness tests pass, with the new sentences asserted.
- `design/system/flai-cli.md` describes the prompt as it now reads.
- The flai lint passes.

## Notes

Drafted by the planner for S-0266.
