---
id: T-2188
type: task
nature: remediation
title: I-0210 is closed with the reason naming the ADR and the guard rule, and the criteria are ticked
status: backlog
parent: S-0511
owner: alex
created: 2026-10-28T08:20:34Z
updated: 2026-10-28T08:20:34Z
transitions: []
stream: S-0511
tags: [cli, guard, issues]
touches: [design/issues/I-0210-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md, design/issues/summary.md]
after: [T-2185, T-2186, T-2187]
---
# T-2188 I-0210 is closed with the reason naming the ADR and the guard rule, and the criteria are ticked

## Work

In the story's worktree, run `flai issue close I-0210 --reason "..."` with a reason that names the ADR T-2184 recorded and the rule T-2185 added: `flai guard` now refuses every sub-agent and every session `flai serve` started the `flai serve` sub-commands that write the host configuration or turn a host action on or off, and `flai config set` and `flai template use`. The command writes the issue file and regenerates `design/issues/summary.md`; commit both with the task. Then verify each acceptance criterion against the branch and tick it with `flai criteria tick S-0511 <n>`: criterion 1 by `flai test flai/internal/guard` passing with `hostconfig_test.go`, criterion 2 by the issue's `status: closed` and its close entry. Waits for T-2185, T-2186, and T-2187, so the close names a fix that is whole and the close-out follows at once.

## Done when

I-0210 reads `status: closed` with a reason naming the ADR and the rule, `design/issues/summary.md` no longer lists it as open, both criteria of S-0511 are ticked, and `flai check --strict` is clean.

## Notes

The issue's remediation section already names S-0511; `flai issue close` appends the close line below it. Run it in the worktree, not the main checkout, so the change lands on the story's branch.
