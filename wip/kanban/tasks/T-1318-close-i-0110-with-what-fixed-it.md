---
id: T-1318
type: task
nature: remediation
title: Close I-0110 with what fixed it
status: backlog
parent: S-0324
owner: alex
created: 2026-10-08T00:29:35Z
updated: 2026-10-08T00:29:35Z
transitions: []
stream: S-0324
tags: [issues]
touches: [design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/summary.md]
after: [T-1317]
---
# T-1318 Close I-0110 with what fixed it

## Work

Second layer: it waits for T-1317, because the reason names the fix and the test that T-1317 lands.

- In the story's worktree, run `flai issue close I-0110 --reason "<reason>"`. The reason says that `parseRange` now records a bare `www.` literal in `out.urls`, so MD034 reports it as markdownlint-cli2 0.20.0 does. It also names the `www.md` fixture and the I-0110 test.
- Check that `design/issues/summary.md` no longer lists I-0110.

## Done when

- I-0110's status is closed, and its reason names the fix and the test.
- `design/issues/summary.md` is current, and `flai check --strict` is clean.

## Notes

Drafted by the planner for S-0324. Run `flai issue` from the worktree, which writes under the checkout it runs in.
