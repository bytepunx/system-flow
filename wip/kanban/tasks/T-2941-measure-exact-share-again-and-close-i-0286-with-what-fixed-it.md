---
id: T-2941
type: task
nature: remediation
title: Measure exact_share again and close I-0286 with what fixed it
status: ready
parent: S-0674
owner: alex
created: 2026-11-10T08:50:19Z
updated: 2026-11-11T01:43:21Z
transitions:
  - to: ready
    at: 2026-11-11T01:43:21Z
    by: system-flow
stream: S-0674
tags: [flai, issues]
touches: [design/issues/I-0286-declared-touches-match-what-a-story-changes-in-16-percent-of-completed-stories.md, design/issues/summary.md]
after: [T-2939, T-2940]
---
# T-2941 Measure exact_share again and close I-0286 with what fixed it

## Work

It waits for T-2939 and T-2940: it measures what both changed.

- Run the worktree's flai (`scripts/flai.sh stats --json`). Record in the story's Notes, next to the Planning table, `claims.weeks[].exact` and `exact_share` for W43 to W46, and the count of `claims.drift[]` entries with an `outside` file. The planning figures were 25 of 115 exact, with all 80 `outside` cases under flai's folders, and 97 of 115 expected after the drift change.
- Close I-0286 with `flai issue close I-0286 --reason "..."`. The reason names the ADR, the two changes, and the measured share.
- Tick each criterion you verified with `flai criteria tick`.

## Done when

- The story's Notes hold the measured figures.
- I-0286 is closed with a reason naming the ADR and the measured share, and `design/issues/summary.md` no longer lists it as open.
- Both criteria are ticked.

## Notes
