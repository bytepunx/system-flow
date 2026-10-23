---
id: T-2052
type: task
nature: improvement
title: The design and the user docs describe the orchestrator's acceptance on a dirty main, the send-back on a stop it cannot settle, and the thread that does not gate it
status: done
parent: S-0467
owner: alex
created: 2026-10-22T21:00:58Z
updated: 2026-10-23T05:36:27Z
transitions:
  - to: ready
    at: 2026-10-23T05:20:48Z
    by: agent-S-0467
  - to: in-progress
    at: 2026-10-23T05:20:49Z
    by: agent-S-0467
  - to: done
    at: 2026-10-23T05:36:27Z
    by: agent-S-0467
stream: S-0467
tags: [flai]
touches: [design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-2039, T-2042, T-2047]
usage:
  source: log
  seconds: 938
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 19719
      cache_read: 3494228
      cache_write: 171842
      cost: 2.0423
---
# T-2052 The design and the user docs describe the orchestrator's acceptance on a dirty main, the send-back on a stop it cannot settle, and the thread that does not gate it

## Work

Waits for the three code tasks, so that the design describes what flai does.

- `design/system/strategic-agents.md` § Accepting a story (S-0221): step 4 and the blocker table. The "No thread is open" row leaves out a thread the orchestrator opened that nobody else wrote on, listed as a note; a new row or paragraph for the uncommitted file the story's change would overwrite, as ADR-0178's landing; the resumed acceptance the orchestrator finishes when it made it; and the send-back when the rebase stops on a conflict flai cannot settle, with the thread `flai accept` opens and `flai serve` starting the agent. Link the ADR the first task wrote.
- `design/system/workflow.md` § What waits on an acceptance (S-0296): the sentence "A wait remains only for a story the orchestrator cannot vouch for" gains what the orchestrator no longer waits on, and the eligibility list's "no thread ... is open" bullet notes the thread the orchestrator alone wrote. § Transitions, where the queue's send-back is listed, gains the acceptance's send-back by `flai`.
- `design/system/flai-cli.md`: the `flai accept` row says what `--by orchestrator` does on a dirty main, on a stop it cannot settle, and with a resumed acceptance of its own.
- `docs/users/flai.md` § When the orchestrator accepts: the same, in the user's terms, with what the operator still clears (an uncommitted file the change would overwrite, a resumed acceptance someone else made, an undone path, a thread anyone else wrote on).

Bump each document's `updated`. Keep the merge queue's sections as they are: the queue already settles these.

## Done when

- The four documents describe the three remedies as built, link the ADR, and `flai check --strict` and the markdown lint pass on them.

## Notes

Layer 4, beside the prompt-and-convention task; no path in common with it. `docs/users/flai-reference.md` is generated from the flags and needs no change: no flag is added.
