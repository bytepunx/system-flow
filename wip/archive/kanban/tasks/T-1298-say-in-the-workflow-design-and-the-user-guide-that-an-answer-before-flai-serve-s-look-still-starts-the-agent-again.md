---
id: T-1298
type: task
nature: remediation
title: Say in the workflow design and the user guide that an answer before flai serve's look still starts the agent again
status: done
parent: S-0317
owner: alex
created: 2026-10-08T00:01:56Z
updated: 2026-10-08T00:17:17Z
transitions:
  - to: ready
    at: 2026-10-08T00:16:25Z
    by: agent-S-0317
  - to: in-progress
    at: 2026-10-08T00:16:26Z
    by: agent-S-0317
  - to: done
    at: 2026-10-08T00:17:17Z
    by: agent-S-0317
stream: S-0317
tags: [flai]
touches: [design/system/workflow.md, docs/users/flai.md]
after: [T-1297]
usage:
  source: log
  seconds: 51
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 14
      output: 3674
      cache_read: 618510
      cache_write: 29177
      cost: 0.3933
---
# T-1298 Say in the workflow design and the user guide that an answer before flai serve's look still starts the agent again

## Work

Add one sentence, naming S-0317, to the **An answered agent** bullet in `design/system/workflow.md` (§ Transitions and who makes them › The pull order), after the S-0335 sentences: a run whose agent asked on a thread of its story during the run, and whose question was answered before flai serve judged its end, is recorded `asked` too, and started again at once in its session, so an answer that comes before the look is not lost and costs no automatic restart.

Say the same for the operator in `docs/users/flai.md`, in the paragraphs on an agent that ended asking a question (the S-0272 and S-0335 paragraphs under flai serve): answering at once, while the agent is ending, still starts it again.

It waits for T-1297 so the words describe what the code does, including the `Why` it records.

## Done when

- Both documents say that an answer given before flai serve's look still starts the agent again, and name S-0317.
- `flai test design/system/workflow.md docs/users/flai.md` passes the markdown lint.

## Notes
