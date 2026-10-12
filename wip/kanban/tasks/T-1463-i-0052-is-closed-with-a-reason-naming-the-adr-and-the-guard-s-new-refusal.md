---
id: T-1463
type: task
nature: remediation
title: I-0052 is closed with a reason naming the ADR and the guard's new refusal
status: backlog
parent: S-0368
owner: alex
created: 2026-10-12T16:57:45Z
updated: 2026-10-12T16:57:45Z
transitions: []
stream: S-0368
tags: [guard, issues]
touches: [design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md, design/issues/summary.md]
after: [T-1461, T-1462]
---
# T-1463 I-0052 is closed with a reason naming the ADR and the guard's new refusal

## Work

Close I-0052 once the fix is built and documented. It waits for T-1461 and T-1462 because the reason names what they shipped.

- In the story's worktree, run `flai issue close I-0052 --reason "<reason>"`. Issue commands write under the checkout they run in.
- The reason names the ADR, what the guard now refuses a read-only sub-agent, and the tests that reproduce the issue.
- Check that `design/issues/summary.md` no longer lists I-0052.

## Done when

- [ ] I-0052's status is closed, and its `## Remediation` holds the reason.
- [ ] `design/issues/summary.md` is regenerated without it.
- [ ] `flai check --strict` passes.

## Notes

This task ticks the story's second criterion. T-1461's tests tick the first.
