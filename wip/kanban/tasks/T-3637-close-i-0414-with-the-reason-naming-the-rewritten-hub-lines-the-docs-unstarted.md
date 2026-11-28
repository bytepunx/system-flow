---
id: T-3637
type: task
nature: remediation
title: Close I-0414 with the reason naming the rewritten hub lines, the docs.unstarted-story check, and the project additions
status: backlog
parent: S-0856
owner: alex
created: 2026-11-27T23:58:42Z
updated: 2026-11-27T23:58:42Z
transitions: []
stream: S-0856
tags: [issues]
touches: [design/issues/I-0414-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md, design/issues/summary.md]
after: [T-3627, T-3632, T-3635]
---
# T-3637 Close I-0414 with the reason naming the rewritten hub lines, the docs.unstarted-story check, and the project additions

## Work

- Run `flai issue close I-0414 --reason "<reason>"` from the story's worktree. The reason names what fixed it: the hub sentences in `design/system/workflow.md` and `docs/operators/runbooks/change-the-board.md` rewritten as one sentence per step, the advisory `docs.unstarted-story` and its ADR, and the project additions to `documentation.md` and `strategic-agents.md`.
- `flai issue close` rewrites `design/issues/summary.md`; commit both.
- Waits for the other three tasks, since the reason names what each of them did.

## Done when

- I-0414's status is closed and its reason names the three fixes.
- `design/issues/summary.md` no longer lists I-0414, and `flai check --strict` reports no `issues.summary`.

## Notes
