---
id: T-0852
type: task
nature: improvement
title: delegation.md, the design, and the user docs say how the verifier runs the close-out once and what the story agent tells it
status: backlog
parent: S-0266
owner: alex
created: 2026-10-05T00:19:12Z
updated: 2026-10-05T00:19:12Z
transitions: []
stream: S-0266
tags: [conventions, docs, template]
touches: [design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/system/workflow.md, docs/users/flai.md]
after: [T-0850, T-0851]
---
# T-0852 delegation.md, the design, and the user docs say how the verifier runs the close-out once and what the story agent tells it

## Work

Change the baseline of `design/conventions/delegation.md` and of `template/root/design/conventions/delegation.md`, keeping the same text in both:

- In the close-out bullet under "While working", the story agent names any step it knows will stop, and why, when it hands the run to the verifier.
- Under "As an explorer or a verifier", state T-0850's rule. The verifier runs a long script once, in one command, without a pipe or a file, and with the longest timeout. It reads the last line, reports an expected stop, and goes on to the steps after it.

In `design/system/workflow.md`, where close-out is described, record that close-out ends every run with one line giving the outcome and the step it stopped at. Record there too that the verifier reads that line from a single run. Say the same in `docs/users/flai.md`, where it describes `scripts/close-out.sh`.

The story also names `design/system/devex.md`, "or where close-out is described". That place is `workflow.md`. S-0267 rewrites the test tiers in `devex.md`, so this task leaves that file alone.

This task waits for T-0850 and T-0851, so that the convention repeats the words of the definition and of the prompt.

## Done when

- The two copies of `delegation.md` have the same baseline.
- `workflow.md` and `docs/users/flai.md` describe the last line and the single run.
- The markdown lint and `flai check --strict` pass.

## Notes

Drafted by the planner for S-0266.
