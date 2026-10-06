---
id: T-1061
type: task
nature: improvement
title: The conventions, their template copies, and the harness prompt send the agent to flai stream state
status: backlog
parent: S-0271
owner: alex
created: 2026-10-06T22:50:11Z
updated: 2026-10-06T22:50:11Z
transitions: []
stream: S-0271
tags: [conventions]
touches: [design/conventions/tooling.md, design/conventions/work-management.md, design/conventions/session-start.md, template/root/design/conventions/tooling.md, template/root/design/conventions/work-management.md, template/root/design/conventions/session-start.md, template/CHANGELOG.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
after: [T-1055, T-1056]
---
# T-1061 The conventions, their template copies, and the harness prompt send the agent to flai stream state

## Work

Criterion 3, for the rules the agent reads.

- `tooling.md` says the summary sections are "edited by hand". Change it to say that `## Current state` and `## Next steps` are written with `flai stream state`, or the MCP tool `stream_state`. The other summary sections stay hand-edited.
- `work-management.md` (the "Keep the narrative current" bullet) and `session-start.md` (the "Before any long-running" and "End every session" bullets): say to rewrite the two sections with `flai stream state`, never with Edit or Write.
- Make each `template/root` copy match its baseline above the marker. Add a line to `template/CHANGELOG.md`, in the same story, as CLAUDE.md says.
- In `flai/internal/harness/harness.go`, in the story agent's start prompt beside the `criteria tick` sentence, add one sentence: rewrite the two sections with `flai stream state` (or the flai MCP tool `stream_state`) at every task transition, never by editing the narrative. Update `harness_test.go`'s expected prompt.

## Done when

- [ ] No convention, in either copy, tells the agent to edit `## Current state` or `## Next steps` by hand.
- [ ] Each convention and its `template/root` copy are identical above the marker, and `template/CHANGELOG.md` records the change.
- [ ] The start prompt names `flai stream state`, and `harness_test.go` passes.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner for S-0271. Ticking is already sent to `flai criteria tick` in both the prompt and `work-management.md` (ADR-0089), so this task adds only the narrative state.
