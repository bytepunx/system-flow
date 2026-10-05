---
id: T-0899
type: task
nature: feature
title: The orchestrator's prompt and definition say how it plans, finalizes, promotes, and orders by its permissions, logging each action with its policy figure
status: backlog
parent: S-0219
owner: alex
created: 2026-10-05T04:47:25Z
updated: 2026-10-05T04:47:48Z
transitions: []
stream: S-0219
tags: [flai, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml]
after: [T-0893, T-0896]
---
# T-0899 The orchestrator's prompt and definition say how it plans, finalizes, promotes, and orders by its permissions, logging each action with its policy figure

## Work

S-0218 gives the orchestrator its prompt in `harness.go` (beside `planPrompt`) and its definition in `.claude/agents/orchestrator.md`: prime, read the board and inbox, act within permissions, log, wait. This task adds what it does with each of this story's permissions, in that prompt, the definition, and the convention's "As the orchestrator", each saying the same in short:

- **`plan_backlog_epics`.** Run `flai plan --candidates` and start the planner (`plan`) for each, one at a time.
- **`finalize_drafts`.** Run `flai promote --drafts`. Finalize a complete draft whose criteria, touches, forecast, and value it judges consistent; for any other, open a thread on the story saying what is missing or inconsistent, once, and leave it.
- **`promote_to_ready`.** Run `flai promote --candidates` and move them to `ready` in its order while the ready limit has room. Never a draft, never a held story.
- **`order_ready`.** After each change to the ready column, its own or another's, run `flai order --by <policy> --apply`, which keeps a story the operator placed by hand within the last day.
- **Logging.** Log every action with `activity_log`: what it did, to which item, and the policy figure that justified it (the value, the value over duration, the forecast, or the candidate's rank). A refusal from `flai guard` or from flai ends that attempt: log it with the refusal, and do not retry it until something changes.

Convention changes land in `design/conventions` and in the template in the same task, with a template version and a CHANGELOG entry. Edit the baseline above the marker only as a template release, as `conventions.md` says.

This task waits for T-0893 and T-0896, for the calls that pass and the refusals the prompt must name. It runs with T-0901, whose paths it does not share.

## Done when

- `harness_test.go` pins the orchestrator's prompt naming each permission's command, the policy figure in the log, and ending an attempt on a refusal
- `orchestrator.md` and the convention, here and in the template, say the same; the template's version and CHANGELOG are bumped
- `go test ./internal/harness/` and the markdown lint pass

## Notes
