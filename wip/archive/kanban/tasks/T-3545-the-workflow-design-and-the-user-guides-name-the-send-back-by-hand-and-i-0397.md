---
id: T-3545
type: task
nature: remediation
title: The workflow design and the user guides name the send back by hand, and I-0397 is closed
status: done
parent: S-0823
owner: alex
created: 2026-11-17T08:31:08Z
updated: 2026-11-17T10:52:30Z
transitions:
  - to: ready
    at: 2026-11-17T10:48:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T10:48:52Z
    by: system-flow
  - to: done
    at: 2026-11-17T10:52:30Z
    by: agent-S-0823
stream: S-0823
tags: [flai]
touches: [design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0397-a-story-the-operator-moves-from-review-to-in-progress-by-hand-gets-no-agent-and.md, design/issues/summary.md]
after: [T-3541, T-3543]
usage:
  source: log
  seconds: 218
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 56
      output: 9561
      cache_read: 2016221
      cache_write: 68537
      cost: 1.0659
---
# T-3545 The workflow design and the user guides name the send back by hand, and I-0397 is closed

## Work

Bring the documents in step with what T-3541 and T-3543 built, in their words:

- `design/system/workflow.md`: in Transitions and who makes them, the review to in-progress row says that a send back by hand, from `flai move` or the dashboard, by anyone but `flai` and the story's agent, opens the `Acceptance sent back:` thread with the reason, that `flai serve` starts the story's agent on it and resolves it once the story is back in review; the pull order's paragraph on agents started on a send back names it beside the others. Say that the `reason-given` rule, which the baseline's send-back transition lists, is why no board marker is needed.
- `design/system/flai-cli.md`: the `flai move` section, the same in a sentence, linking the workflow section.
- `docs/users/flai.md`: the `flai move ... in-progress --reason` example and its text.
- `docs/users/flaiover.md`: moving a story out of review on the board, its reason now reaching the story's agent on a thread.

Then close the issue with `flai issue close I-0397 --reason "<what fixed it>"`, naming T-3541 and T-3543, which rewrites the issue file and `design/issues/summary.md`.

This task waits for T-3541 and T-3543, whose behaviour and words it documents, and the issue is closed only once the fix is in.

## Done when

- The four documents say what a send back by hand does, and the markdown lint passes on them.
- I-0397 is closed with a reason naming what fixed it, and `flai check --strict` is clean.

## Notes
