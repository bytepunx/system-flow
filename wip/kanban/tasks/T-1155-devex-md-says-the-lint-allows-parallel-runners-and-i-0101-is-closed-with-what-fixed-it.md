---
id: T-1155
type: task
nature: remediation
title: devex.md says the lint allows parallel runners, and I-0101 is closed with what fixed it
status: backlog
parent: S-0308
owner: alex
created: 2026-10-07T02:19:57Z
updated: 2026-10-07T02:19:57Z
transitions: []
stream: S-0308
tags: [docs, issues]
touches: [design/system/devex.md, design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md, design/issues/summary.md]
after: [T-1153]
---
# T-1155 devex.md says the lint allows parallel runners, and I-0101 is closed with what fixed it

## Work

- In `design/system/devex.md`, add to the "Go build and lint" row that `flai/.golangci.yaml` sets `run.allow-parallel-runners`, so lint runs in parallel story worktrees on one host do not stop each other (I-0101, S-0308). Bump `updated`.
- Close the issue from the story worktree: `flai issue close I-0101 --reason` naming the config key, the test, and S-0308. That regenerates `design/issues/summary.md`.
- Waits for T-1153: the design row and the close reason describe what T-1153 settled, and the issue closes only once the fix and its test are in.

## Done when

- `design/system/devex.md` names the setting and why.
- I-0101 is closed with a reason saying what fixed it, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes

Run `flai issue close` in the story worktree: issue commands write under the checkout they run in.
