---
id: T-1149
type: task
nature: improvement
title: I-0076 is closed with flai issue close, saying that a close-out records no wip.overlap and that wip.overlap compares claims
status: done
parent: S-0279
owner: alex
created: 2026-10-07T01:21:01Z
updated: 2026-10-07T09:42:31Z
transitions:
  - to: ready
    at: 2026-10-07T09:42:24Z
    by: agent-S-0279
  - to: in-progress
    at: 2026-10-07T09:42:24Z
    by: agent-S-0279
  - to: done
    at: 2026-10-07T09:42:31Z
    by: agent-S-0279
stream: S-0279
tags: [flai]
touches: [design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1145]
usage:
  source: log
  seconds: 7
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2492
      cache_read: 396810
      cache_write: 15095
      cost: 0.2267
---
# T-1149 I-0076 is closed with flai issue close, saying that a close-out records no wip.overlap and that wip.overlap compares claims

## Work

Close I-0076 with `flai issue close I-0076 --reason`, run in the story's worktree, as the story's second criterion asks. The reason names what fixed it:

- T-1145: a close-out records no `wip.overlap`, and it lists only the overlaps that name the story.
- T-1143: `wip.overlap` compares the claims of the stories in progress, one finding per pair.
- T-1141's ADR, by number.

It waits for T-1145, because the issue is closed only once the fix is in. It shares no file with T-1150, so the two run together.

## Done when

- I-0076's status is `closed`, and its reason names the ADR and the two changes.
- `design/issues/summary.md` no longer lists I-0076, as `flai issue close` regenerates it.
- The story's own close-out records no new I-0076 instance and opens no new `wip.overlap` issue.

## Notes
