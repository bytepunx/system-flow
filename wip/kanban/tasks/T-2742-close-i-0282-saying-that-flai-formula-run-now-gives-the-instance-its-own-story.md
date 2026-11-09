---
id: T-2742
type: task
nature: remediation
title: Close I-0282 saying that flai formula run now gives the instance its own story file as its touch
status: backlog
parent: S-0620
owner: alex
created: 2026-11-09T17:06:12Z
updated: 2026-11-09T17:06:12Z
transitions: []
stream: S-0620
tags: [issues]
touches: [design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md, design/issues/summary.md]
after: [T-2740, T-2741]
---
# T-2742 Close I-0282 saying that flai formula run now gives the instance its own story file as its touch

## Work

Run `flai issue close I-0282 --reason "<reason>"` in the story's worktree. The reason names S-0620, the change to `formulaRun` in `flai/cmd/formula.go`, and the test that reproduces the hold. `flai issue close` writes the issue file and `design/issues/summary.md`.

Waits on T-2740 and T-2741, because the issue closes only once the fix is built, tested, and documented.

## Done when

- I-0282 is `closed`, its reason naming what fixed it.
- `design/issues/summary.md` no longer lists I-0282.
- `flai check --strict` passes.

## Notes
