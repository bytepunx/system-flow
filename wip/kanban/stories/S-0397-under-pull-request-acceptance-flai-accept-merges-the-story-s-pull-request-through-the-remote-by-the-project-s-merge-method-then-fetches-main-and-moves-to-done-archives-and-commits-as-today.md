---
id: S-0397
type: story
nature: feature
title: Under pull-request acceptance, flai accept merges the story's pull request through the remote by the project's merge method, then fetches main and moves to done, archives, and commits as today
status: backlog
parent: E-0021
owner: alex
created: 2026-10-16T04:39:46Z
updated: 2026-10-17T01:08:46Z
transitions: []
tags: [cli, template]
topics: [release, orchestration]
touches: [flai/cmd/accept.go, flai/cmd/accept_remote.go, flai/cmd/accept_remote_test.go, flai/internal/serve/forgemerge.go, flai/internal/serve/forgemerge_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, design/adrs, design/system/workflow.md, design/system/flai-cli.md, design/conventions/git.md, template/root/design/conventions/git.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md]
after: [S-0395]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 45m
  delivery: 2026-10-17T06:46:00Z
  basis: "Its own forecast of 45m; 13th in the pull order with an in-progress limit of 5, behind S-0440, S-0441, S-0439, S-0442, S-0444, S-0443, S-0241, S-0392, S-0393, S-0394, S-0395 and S-0396."
  by: flai
  at: 2026-10-17T01:08:46Z
---
# S-0397 Under pull-request acceptance, flai accept merges the story's pull request through the remote by the project's merge method, then fetches main and moves to done, archives, and commits as today

## Goal

A team whose main branch takes changes only through pull requests can accept stories with flai. Under pull-request acceptance, `flai accept` merges the story's pull request on the forge by the project's merge method, rather than on the host. It then brings the local main branch up to the remote's, and moves the story to done, archives it, and commits, as today. The acceptance commit and the archive read the same in both modes, so metrics and history do too. Whoever accepts, the operator, the dashboard, or the orchestrator, the merge goes through flai serve with the token the host names.

## Acceptance criteria

- [ ] Under `acceptance.mode: pull-request`, once S-0395's gate passes, `flai accept S-nnnn` merges the story's pull request on the forge with `acceptance.merge_method`, and deletes the remote story branch. Then:
  - it fetches main and brings the local main branch up to `origin`'s, rebasing onto it any local commits not yet published
  - it moves the story to done, archives it, and commits, with the acceptance commit's subject and the archive as under `local`
  - it removes the story's local branch and worktree as today
- [ ] Where the forge token's variable is unset, as it is for every agent, `flai accept` asks flai serve on the host to make the merge and waits for its answer. The dashboard's acceptance merges in flai serve directly. Either way it refuses, naming the host action, while `forge` is off.
- [ ] A pull request already merged on the forge, for example by a reviewer, is accepted without a second merge: `flai accept` fetches and goes on with done, archive, and commit.
- [ ] A merge the forge refuses, for a conflict, a branch out of date, or a rule, stops the acceptance before anything changes locally, naming the forge's reason. A local main that has unpublished commits outside `wip/` also stops it before the merge, naming those commits.
- [ ] The dry run says that it would merge through the remote, by which method, or that the request is already merged.
- [ ] An ADR records the merge through the remote and the hand-off to flai serve. `design/system/workflow.md` (review to done), `design/conventions/git.md` and the template's (accepting an item), `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md` describe both modes.

## Tasks
- T-1633 An ADR records that under pull-request acceptance the merge goes through the remote, by flai serve when the token is not at hand, and that the acceptance commit reads as under local acceptance
- T-1634 flai serve merges a story's pull request on the forge when asked, by the project's method, deletes the remote story branch, and answers whether it merged or why not
- T-1635 flai accept under pull-request acceptance merges through the remote, directly or through flai serve, then fetches, rebases unpublished wip commits, and archives and commits as today
- T-1636 The dashboard's acceptance under pull-request acceptance merges through the remote in flai serve, and refuses while the forge host action is off
- T-1637 The workflow, the git convention and the template's, and the CLI and operator docs describe acceptance through the remote beside the local merge

## Notes

- Drafted by the planner for E-0021, after S-0395, whose gate it passes first.
- Under pull-request acceptance the story's code reaches the remote's main branch at the merge, but the acceptance commit, with the archive, reaches it only when published, as ADR-0067 has it. A main branch that refuses every direct push would also refuse that publish. The plan's thread raises this. This story does not solve it.

### Planning

- **Touches.** None were declared. Every touch is the planner's, and each is a file except `design/adrs`:
  - **Layout:**
    - `flai/cmd/accept.go`, with a new `accept_remote.go` and its test, beside `mergeStoryBranch` in `branch.go`. `branch.go` is called, not changed. Criteria 1 and 3 to 5.
    - `flai/internal/serve/forgemerge.go`, with its test: the merge in flai serve and the hand-off. Criterion 2.
    - `flai/internal/hostapi/writes.go`, with its test: `accept.run` for the dashboard. Criterion 2.
  - **Design:** the workflow, the git convention and the template's, the CLI docs, and the operators' guide, for criterion 6.
- **Folder touch kept.** `design/adrs`, a shared path, for an ADR `flai adr new` names.
- **Co-change left out.**
  - `docs/users/flai-reference.md`: no new flag.
  - `docs/operators/settings.md`: no new setting.
  - The flaiover files: the dashboard's acceptance changes in flai serve, and the confirmation already lists blockers since S-0395.
- **Paths it shares with other stories.** S-0395 changes `flai/cmd/accept.go` and `writes.go` before it, by `after`. E-0022's S-0400 changes `accept.go` too. TH-0464 recommends this story before S-0400.
- **Layers.**
  - First: T-1633 and T-1634 together.
  - Second: T-1635, after T-1634.
  - Third: T-1636, after T-1634 and T-1635.
  - Fourth: T-1637, after T-1633, T-1635, and T-1636.
- **Forecast: 45m.** `flai forecast` gave 34m: 96 s per unit over 61 done feature stories on claude-opus-5-5 in the large band, times size 21. Raised by a third for git against a scratch remote, a merge by three methods, a fetch, and a rebase of unpublished commits, plus the hand-off to flai serve. flai's replan sets the delivery from this duration.
- **Cost of delay: not yet set.** `flai cod` refuses until E-0021 has inputs, asked for on TH-0463.
