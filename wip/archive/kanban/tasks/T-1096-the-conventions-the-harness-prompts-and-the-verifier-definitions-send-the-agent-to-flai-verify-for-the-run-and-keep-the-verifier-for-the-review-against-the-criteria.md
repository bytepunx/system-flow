---
id: T-1096
type: task
nature: improvement
title: The conventions, the harness prompts, and the verifier definitions send the agent to flai verify for the run and keep the verifier for the review against the criteria
status: done
parent: S-0270
owner: alex
created: 2026-10-06T22:53:06Z
updated: 2026-10-07T06:40:45Z
transitions:
  - to: ready
    at: 2026-10-07T03:49:40Z
    by: agent-S-0270
  - to: in-progress
    at: 2026-10-07T03:49:40Z
    by: agent-S-0270
  - to: done
    at: 2026-10-07T06:40:45Z
    by: agent-S-0270
stream: S-0270
tags: [flai, template, conventions]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, ".claude/agents/verifier.md", template/root/.claude/agents/verifier.md, template/CHANGELOG.md, design/adrs, design/conventions/code-quality.md, design/conventions/strategic-agents.md, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, template/root/design/conventions/code-quality.md, template/root/design/conventions/strategic-agents.md, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-1075]
usage:
  source: log
  seconds: 4140
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 190
      output: 79556
      cache_read: 13771940
      cache_write: 399724
      cost: 6.8407
    - model: claude-sonnet-5-5
      input: 14
      output: 3613
      cache_read: 147042
      cache_write: 61311
      cost: 0.2188
---
# T-1096 The conventions, the harness prompts, and the verifier definitions send the agent to flai verify for the run and keep the verifier for the review against the criteria

## Work

Criterion 4: the agent runs `flai verify` itself, and the verifier sub-agent is kept for judgement.

- Write an ADR that refines ADR-0059's hand-off. The story's agent runs the tests, the lint, and the check through `flai verify`, or the close-out that calls it, and reads the result. It hands the verifier only the review of the diff against the acceptance criteria and the conventions, when a review pays. Number it with `flai adr new`.
- In `flai/internal/harness/harness.go`, rewrite the story agent's `delegation()` paragraph on the run before review. It should say to run the close-out once, read its result, fix what it names, and run it again, with no verifier for the run. Do the same for the orchestrator's `accept_reviews` paragraph: run `flai verify` for the story, and have the verifier check the diff against each criterion and name the commit. Update the asserted sentences in `harness_test.go`.
- In `design/conventions/delegation.md`, change the verifier's part of "As an explorer or a verifier" in the same way. In `work-management.md`, make the close-out rule in the definition of done name `flai verify` as what the close-out runs. Copy both to `template/root` and add a line to `template/CHANGELOG.md`. The line says the release needs the flai that ships `flai verify`, as CLAUDE.md says.
- In `.claude/agents/verifier.md` and its template copy, say that the run is `flai verify`, read once, and that the verifier's work is the review. Claude Code asks the operator on a thread before it writes under `.claude/`, so make that edit yourself, last, when the rest of the task is done.
- Waits for T-1075, the command these texts name. Runs alongside T-1082 and T-1089: no path in common.

## Done when

- [ ] The ADR is written and linked from `delegation.md`.
- [ ] The story agent's and the orchestrator's prompts, `delegation.md`, `work-management.md`, and both verifier definitions send the run to `flai verify` and keep the verifier for the review.
- [ ] The convention copies in `template/root` match their originals, and `template/CHANGELOG.md` records the change.
- [ ] `scripts/flai-test.sh`, `flai check --strict`, and the markdown lint pass.

## Notes

Drafted by the planner. The `design/adrs` touch is a folder because the ADR's number is not known until it is written.
