---
id: T-2743
type: task
nature: improvement
title: An ADR and the design say a formula step may name the exits that stop its instance cleanly, and how the instance then reaches done
status: done
parent: S-0621
owner: alex
created: 2026-11-09T17:06:27Z
updated: 2026-11-09T17:54:08Z
transitions:
  - to: ready
    at: 2026-11-09T17:43:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:43:12Z
    by: system-flow
  - to: done
    at: 2026-11-09T17:54:08Z
    by: agent-S-0621
stream: S-0621
tags: [design]
touches: [design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 648
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 146
      output: 36739
      cache_read: 8465566
      cache_write: 205094
      cost: 3.773
---
# T-2743 An ADR and the design say a formula step may name the exits that stop its instance cleanly, and how the instance then reaches done

## Work

I-0283 names no remedy. Its three instances, and the operator's answer on TH-0854 (D-0498), point to one. Write it down before the code, since ADR-0238 is immutable and a new step field changes the formula format.

- Write a new ADR with `flai adr new`, refining ADR-0238 and ADR-0239. Decide in it:
  - the step field, proposed `stop`: a list of exit codes other than 0, on a `flai` or `script` step only, never a role step;
  - the command the story's agent runs on such an exit, proposed `flai formula stop <S-nnnn> --step <id> --exit <n>`, with the step's output on standard input or `--output`;
  - what it writes: the step's task done, with the exit and the last 40 lines of output in its `## Notes`; each later step task cancelled with the reason `step <id> stopped the run with exit <n>`; and each of their criteria ticked and marked skipped with that reason;
  - that the orchestrator's evidence for a skipped criterion is `- <n>: T-nnnn skipped`, and flai takes a step task cancelled by a stop as evidence;
  - what it refuses: a story that is no instance, a step whose formula does not list that exit, a step not in progress or not the next in its layers.
- Say why each alternative was not chosen: the agent cancelling by hand, against the instance rule; closing the skipped steps as done, which misreports them; leaving their criteria unticked, which blocks acceptance.
- Find out why `story/S-0585` was gone when the orchestrator came to accept it (TH-0868, the third instance). Record the cause in the ADR's context, and the fix in its decision if it belongs to this story.
- Update `design/system/strategic-agents.md` § Formulas (S-0421) and § Working an instance (S-0422): the field, a step that stops, and its acceptance. Update `design/system/flai-cli.md` § Commands with the new command.

This task waits for none: the tasks after it read the field and the command it settles.

## Done when

- The ADR is accepted, indexed in `design/adrs/README.md`, and names the field, the command, what it writes, its refusals, and the evidence form.
- `design/system/strategic-agents.md` and `design/system/flai-cli.md` describe them and link the ADR.
- `flai test` on the changed paths passes.

## Notes
