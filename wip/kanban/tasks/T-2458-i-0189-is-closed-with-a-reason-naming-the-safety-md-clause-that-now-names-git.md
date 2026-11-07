---
id: T-2458
type: task
nature: remediation
title: I-0189 is closed with a reason naming the safety.md clause that now names git.md's force-with-lease exception
status: in-progress
parent: S-0558
owner: alex
created: 2026-11-05T01:34:24Z
updated: 2026-11-07T04:38:53Z
transitions:
  - to: ready
    at: 2026-11-07T04:38:50Z
    by: agent-S-0558
  - to: in-progress
    at: 2026-11-07T04:38:53Z
    by: agent-S-0558
stream: S-0558
tags: [conventions]
touches: [design/issues/I-0189-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai-serve-s-force-with-lease-push-of-a-story-s-branch.md, design/issues/summary.md]
after: [T-2457]
---
# T-2458 I-0189 is closed with a reason naming the safety.md clause that now names git.md's force-with-lease exception

## Work

Close I-0189 with `flai issue close I-0189 --reason "<reason>"` in the story's worktree. The reason says what fixed it: S-0558 rewrote the clause in both copies of `safety.md` § Rules so that an agent's force push still always needs confirmation and `flai serve`'s `--force-with-lease` push of a story's own branch under pull-request acceptance is named as the one exception, linking `git.md`.

It waits for T-2457, because the reason names the clause T-2457 writes, and the issue is closed only once that clause is committed.

Record in the story's Notes why no automated test reproduces the cause: it is a contradiction between two sentences of prose, and the template tier's render and `flai check --strict` of the rendered project is the check that applies. Tick criterion 1 only once T-2457's diff shows both copies carry the clause.

## Done when

- I-0189's status is closed, with the reason above.
- `design/issues/summary.md` no longer lists I-0189 as open, and `flai check --strict` is clean on `design/issues`.
- The story's Notes say why no automated test fits.

## Notes
