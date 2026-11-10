---
id: T-2788
type: task
nature: improvement
title: Record in an ADR how the orchestrator's judgement reaches the publish formula's evaluate step
status: done
parent: S-0625
owner: alex
created: 2026-11-09T22:01:47Z
updated: 2026-11-10T05:19:09Z
transitions:
  - to: ready
    at: 2026-11-10T05:14:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:14:43Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:18:24Z
    by: agent-S-0625
stream: S-0625
tags: [flai, design]
touches: [design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/adrs/0273-under-judgement-the-orchestrator-s-judgement-reaches-the-publish-formula-as-an.md]
usage:
  source: log
  seconds: 221
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 6063
      cache_read: 5011949
      cache_write: 125747
      cost: 1.8503
---
# T-2788 Record in an ADR how the orchestrator's judgement reaches the publish formula's evaluate step

## Work

Write a new ADR, refining ADR-0240 and ADR-0094, that decides how a judgement the orchestrator makes reaches the publish formula. Its context is I-0291 and TH-0811. Under `judgement`, `flai release --evaluate --require-met` always exits 2. The publish formula's step 2 runs it and, since S-0621 (ADR-0271), stops cleanly on that exit. `release_publish` refuses while the formula exists. So the `publish` permission publishes nothing.

Decide, as the story's Planning note proposes:

- `flai release --evaluate --judged "<reason>"`. Under `judgement`, a non-empty reason makes the evaluation met when something is accepted and not yet released and `whole_epics` holds nothing back. The evaluation's reason names the judgement. Under `threshold` and `theme`, the flag changes nothing. An empty reason is the same as none.
- The publish formula declares the input `judgement`, with `default: ""`. Step 2 runs `release --evaluate --require-met --judged "{{judgement}}"`.
- Under `judgement`, the orchestrator runs `flai formula run publish --ready --autocommit --input judgement="<its reasoning>"`. The reasoning is one sentence with no double quote, backtick, or dollar sign, because flai fills the input into a shell line as text (`docs/operators/index.md`, the formula inputs note).
- The guard does not change: `formulaRunFlags` already passes `--input`.

Record the alternatives considered:

- A judgement record the orchestrator writes and `--evaluate` reads.
- Under `judgement`, letting `--require-met` pass on any non-empty batch.
- Letting `release_publish` run beside the formula.

Then update the design to match:

- `design/system/strategic-agents.md` § Publishing a release (S-0222): the formula paragraph and the `judgement` row.
- `design/system/flai-cli.md`: the `release --evaluate` row.
- `design/system/project-manifest.md`: where it says judgement is never met by itself.

Add the ADR to `design/adrs/README.md`. The ADR's file name in this task's touches is predicted: use the number and slug `flai adr new` gives.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`. It names the flag, the input, the orchestrator's command, and the alternatives.
- The three design documents describe the `--judged` flag and the input, and none still says that judgement can never pass the formula's evaluate step.
- `flai check --strict` is clean.

## Notes

Layer 1: waits for no task. The tasks after it build what it decides.
