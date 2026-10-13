---
id: T-1506
type: task
nature: remediation
title: flai stream answer leaves one blank line where it removes an open question, and refuses a save the markdown lint rejects
status: done
parent: S-0376
owner: alex
created: 2026-10-13T16:43:59Z
updated: 2026-10-13T17:18:49Z
transitions:
  - to: ready
    at: 2026-10-13T17:17:40Z
    by: agent-S-0376
  - to: in-progress
    at: 2026-10-13T17:17:41Z
    by: agent-S-0376
  - to: done
    at: 2026-10-13T17:18:49Z
    by: agent-S-0376
stream: S-0376
tags: [flai, narrative]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go]
usage:
  source: log
  seconds: 68
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1678
      cache_read: 529921
      cache_write: 17773
      cost: 0.2818
---
# T-1506 flai stream answer leaves one blank line where it removes an open question, and refuses a save the markdown lint rejects

## Work

- I-0139's five instances are one file, `wip/archive/agents/S-0232.md`, whose `## Open questions` was left as the heading, two blank lines, and `## Log`. Before the acceptance commit `7801de00` it held one hand-written bullet between two blank lines, under the mirrored thread block. Removing that bullet from the template's layout leaves exactly that.
- `AnswerOpenQuestion` in `flai/internal/workitem/narrative.go` removes the bullet with `lines = append(lines[:bStart], lines[bEnd:]...)`. It does not collapse the blank lines on either side. It saves with `n.Save()`, without the lint guard that `LogStream` and `SetStreamState` go through.
- Collapse the blank lines around the removed bullet to one, so `## Open questions`, blank, `- Q`, blank, `## Log` becomes `## Open questions`, blank, `## Log`.
- Save through `r.LintGuard(n.Path, was, n.Marshal())`, as `LogStream` does, so that an answer that would bring a markdown finding is refused with its rule and line.
- Add a test in `flai/internal/workitem/narrative_test.go` that reproduces I-0139: the template's layout with one bullet, answered, has no `\n\n\n` and reads `## Open questions\n\n## Log`. Add a second case with a bullet between two others.
- It waits for no task: no other task of S-0376 changes these paths.

## Done when

- [ ] Answering the only open question of a narrative in the template's layout leaves one blank line between `## Open questions` and `## Log`, checked by a test that fails before the change.
- [ ] `AnswerOpenQuestion` saves through the lint guard.
- [ ] `flai test` on the changed paths passes.

## Notes
