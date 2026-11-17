---
id: T-3601
type: task
nature: improvement
title: release_publish fetches, evaluates the policy again, publishes only what is pending, names the release that covers each component, and no longer refuses for a publish formula
status: backlog
parent: S-0837
owner: alex
created: 2026-11-17T12:37:35Z
updated: 2026-11-17T12:37:35Z
transitions: []
stream: S-0837
tags: [flai]
touches: [flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_test.go]
after: [T-3600]
---
# T-3601 release_publish fetches, evaluates the policy again, publishes only what is pending, names the release that covers each component, and no longer refuses for a publish formula

## Work

This task waits for T-3600, because it builds ADR-0316 §1 and §2 as accepted.

- `releasePublish` in `flai/internal/mcpserver/orchestrate.go` evaluates once on local refs today, and `flai release --pending` only asks the remote with `ls-remote`. Before the evaluation, fetch the remote's branch and tags in the main checkout, as the formula's `fetch` step did (`git fetch --tags`). Then evaluate the policy again with `EvaluateRepo` or, under `judgement`, the judged evaluation the ADR keeps. Publish only when the evaluation is met and something is pending.
- When another publish already released what was pending, answer which release covers each component and change nothing. `flai release --pending --require-pending` writes `covered_by` (exit 2). The MCP path runs `release --pending --json` without it, and `publishExits` in `flai/internal/hostapi/writes_outcome.go` maps only exit 3. Give `release_publish` the flag and map exit 2 to an answer. Leave the dashboard's `publish.run` unchanged (ADR-0316 §6): if `publishArgs` is shared, give the orchestrator's path its own arguments.
- Remove `noPublishFormula` and its call, its test (`orchestrate_test.go`, the formula refusal), and the tool description's words on the formula.
- Keep a moved remote's refusal and the `push` host action's refusal word for word.
- Add the criterion's test beside `TestReleasePublishPublishesWhenThePolicyIsMet`, on the `publishProject` fixture. Read the board, the stories, the tasks, the branches, and the worktrees before and after one `release_publish` that publishes, and assert nothing was made. Add tests for the fetch (a tag only the remote has counts as released) and for the `covered_by` answer.

## Done when

- `release_publish` fetches, evaluates again, and publishes only when met and pending.
- `release_publish` answers `covered_by` with nothing changed when another publish got there first.
- It no longer refuses for a formula named `publish`.
- A test shows one call makes no story, task, branch, or worktree.
- `flai test` passes on the changed paths.

## Notes

- Drafted by planner-S-0837.
- The explorer found that the fake runner in `orchestrate_test.go` asserts the publish arguments exactly, so it changes with the arguments.
