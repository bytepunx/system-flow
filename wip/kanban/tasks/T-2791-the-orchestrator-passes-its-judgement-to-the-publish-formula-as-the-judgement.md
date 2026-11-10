---
id: T-2791
type: task
nature: improvement
title: The orchestrator passes its judgement to the publish formula as the judgement input
status: done
parent: S-0625
owner: alex
created: 2026-11-09T22:02:09Z
updated: 2026-11-10T05:44:50Z
transitions:
  - to: ready
    at: 2026-11-10T05:19:15Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:19:16Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:44:50Z
    by: agent-S-0625
stream: S-0625
tags: [flai, template, conventions]
touches: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/harness/harness_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
after: [T-2788]
usage:
  source: log
  seconds: 1534
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 37
      output: 314
      cache_read: 4799703
      cache_write: 310611
      cost: 1.878
---
# T-2791 The orchestrator passes its judgement to the publish formula as the judgement input

## Work

Tell the orchestrator how to publish under `judgement`, as T-2788's ADR decides. It waits for T-2788 because it names the input and the command the ADR fixes.

- The publish paragraph of the orchestrator's role file, above the marker, in all three copies: `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, and `flai/internal/roles/baseline/orchestrator.md`.
  - Under `judgement`, once it judges the unreleased work coherent and complete, it runs `flai formula run publish --ready --autocommit --input judgement="<its reasoning>"`.
  - The reasoning is one sentence with no double quote, backtick, or dollar sign.
  - Under `threshold` and `theme`, it gives no input.
  - It logs the reasoning with `activity_log`.
- `design/conventions/strategic-agents.md` § As the orchestrator, and its template copy: the `publish` rule names the `judgement` input under that policy.
- `flai/internal/harness/harness_test.go`: `TestTheOrchestratorsPromptSaysHowItPublishes` asserts the new phrase.

## Done when

- The three role file copies match above the marker and name the input under `judgement`.
- The two convention copies match above the marker.
- `TestTheOrchestratorsPromptSaysHowItPublishes` passes, and `flai check --strict` is clean.

## Notes

Layer 2. It runs beside T-2790; the two share no path.
