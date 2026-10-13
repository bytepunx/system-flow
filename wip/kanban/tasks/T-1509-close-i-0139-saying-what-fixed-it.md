---
id: T-1509
type: task
nature: remediation
title: Close I-0139 saying what fixed it
status: backlog
parent: S-0376
owner: alex
created: 2026-10-13T16:44:24Z
updated: 2026-10-13T16:44:24Z
transitions: []
stream: S-0376
tags: [issues]
touches: [design/issues/I-0139-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1508]
---
# T-1509 Close I-0139 saying what fixed it

## Work

- Run `flai issue close I-0139 --reason "<reason>"` in the story's worktree. The reason names the cause: `flai stream answer` left two blank lines where it removed the last open question, and an archived narrative could keep a hand edit's blank run that only a hand commit on main cleared. It also names the fix: T-1506 collapses the blank lines and guards the save, and T-1507 collapses blank runs as a narrative is archived.
- The command regenerates `design/issues/summary.md`.
- It waits for T-1508, so that the reason describes what was built and documented.

## Done when

- [ ] I-0139 is closed with a reason naming the cause and T-1506 and T-1507.
- [ ] `design/issues/summary.md` no longer lists I-0139, and `flai check --strict --story S-0376` is clean.

## Notes
