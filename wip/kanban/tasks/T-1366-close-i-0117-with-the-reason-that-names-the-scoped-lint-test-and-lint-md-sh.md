---
id: T-1366
type: task
nature: remediation
title: Close I-0117 with the reason that names the scoped lint test and lint-md.sh
status: backlog
parent: S-0345
owner: alex
created: 2026-10-08T08:40:07Z
updated: 2026-10-08T08:40:07Z
transitions: []
stream: S-0345
tags: [flai]
touches: [design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/summary.md]
after: [T-1363, T-1364]
---
# T-1366 Close I-0117 with the reason that names the scoped lint test and lint-md.sh

## Work

Close the issue the story remediates, in the story's worktree, so the close lands on the story branch: `flai issue close I-0117 --reason "<what fixed it>"`. The reason names S-0345, that `TestRepositoryLintsClean` and `scripts/lint-md.sh` now scope themselves to the story when `CLOSE_OUT_STORY` is set, as `TestMonorepoIsClean` and `scripts/check.sh` do (ADR-0085), and that `make lint-md` and CI still lint main's `wip/`.

Waits for T-1363 and T-1364, since the issue closes only once the fix is in. Runs beside T-1365: they share no path.

## Done when

- I-0117 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict --story S-0345` reports nothing inside the story.

## Notes

Drafted by the planner for S-0345.
