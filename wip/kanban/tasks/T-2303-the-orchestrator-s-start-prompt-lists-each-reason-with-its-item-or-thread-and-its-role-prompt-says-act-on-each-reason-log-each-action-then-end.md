---
id: T-2303
type: task
nature: feature
title: The orchestrator's start prompt lists each reason with its item or thread, and its role prompt says act on each reason, log each action, then end
status: done
parent: S-0536
owner: alex
created: 2026-10-28T21:46:04Z
updated: 2026-10-29T08:02:12Z
transitions:
  - to: ready
    at: 2026-10-29T07:46:23Z
    by: agent-S-0536
  - to: in-progress
    at: 2026-10-29T07:46:24Z
    by: agent-S-0536
  - to: done
    at: 2026-10-29T08:02:12Z
    by: agent-S-0536
stream: S-0536
tags: [flai, harness, orchestrator, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md]
after: [T-2300, T-2302]
usage:
  source: log
  seconds: 948
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 21545
      cache_read: 4874355
      cache_write: 156088
      cost: 2.3612
---
# T-2303 The orchestrator's start prompt lists each reason with its item or thread, and its role prompt says act on each reason, log each action, then end

## Work

ADR-0209 §3, criterion 2 (the reasons in the prompt) and criterion 5 (the orchestrator's prompt). Two prompts change:

- `orchestratePrompt` in `flai/internal/harness/harness.go`, the start prompt: list each of `Request.Reasons` on its own line with its event and its item or thread, and replace "wait on wait_for_events between decisions, without ending" with: act on each reason, log each action with `activity_log`, then end; flai serve starts you again when there is more. Rewrite its doc comment to match.
- The role file's body, steps 15 and 16 of `design/workflows/roles/orchestrator.md` and its template copy: step 15 becomes act on each reason, log each action, then end; step 16's "hold `wait_for_events` ... Do not end" goes, or becomes: when `wait_for_events` answers `end: true`, end. Update the definition's `description` in the front matter, which says it runs for as long as the action is on.
- Bring `.claude/agents/orchestrator.md` and `template/root/.claude/agents/orchestrator.md`, the definitions Claude Code runs, in step with the role file. A write under `.claude/` is asked of the operator through `permission_prompt`; say so in the narrative before the write.

Until S-0537 lands, the prompt keeps priming the general `orchestrate` pack.

Waits for T-2300, which changes the same role files' trigger, and T-2302, which adds `Request.Reasons` and gives the reasons.

## Done when

- A test in `harness_test.go` builds the orchestrator's start with three reasons and finds each in the prompt, and finds no "without ending".
- Both copies of the role file and both definitions say act on each reason, log each action, then end.
- `flai check --strict` and `flai test` pass on the changed paths.

## Notes
