---
id: T-3035
type: task
nature: remediation
title: The kluster design, the workflow, the orchestrator's blockers, and the user guide say a failed deploy sends the story back and blocks its acceptance
status: backlog
parent: S-0694
owner: alex
created: 2026-11-10T10:33:45Z
updated: 2026-11-10T10:33:45Z
transitions: []
stream: S-0694
tags: [docs, kluster]
touches: [design/system/kluster.md, design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-3025, T-3028, T-3029]
---
# T-3035 The kluster design, the workflow, the orchestrator's blockers, and the user guide say a failed deploy sends the story back and blocks its acceptance

## Work

Write down what T-3025, T-3028, and T-3029 built, as `documentation.md` asks.

- `design/system/kluster.md`: in § A failed deploy holds the story, remove the paragraph that says the send-back and the blocker are not built. Make the flai serve section's "a failure moves nothing" line say the story is sent back and keeps its namespace.
- `design/system/workflow.md` § What waits on an acceptance (S-0296): add the `kluster` blocker and `--without-kluster` beside `--without-queue-tiers`.
- `design/system/strategic-agents.md`: add a `kluster` row to the orchestrator's table of blockers, and say the queue's reason `kluster` keeps the story out of a batch.
- `design/system/flai-cli.md`: add `[--without-kluster]` to the `flai accept` row and the acceptance gate.
- `docs/users/flai.md` § Run a project's kluster: say what a failed deploy or checks run does to the story, what the blocker refuses, and how the operator lands a story with `--without-kluster`.
- Regenerate `docs/users/flai-reference.md` and the flag table in `docs/operators/settings.md` with `make flai-reference`. Never edit those two files by hand.

This task waits for T-3025, T-3028, and T-3029, whose behaviour it describes, and for T-3028's help text, which the regenerated reference reads.

## Done when

- No document under `design/system` or `docs/users` says the send-back or the `kluster` blocker is not built.
- `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` pass.
- The markdown lint passes on the seven files.

## Notes
