---
id: S-0275
type: story
nature: improvement
title: "Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:33Z
updated: 2026-10-06T22:53:55Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions]
touches: [flai/cmd/issue.go, flai/cmd/issue_test.go, flai/cmd/adr.go, flai/cmd/adr_test.go, flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/internal/issues/record.go, flai/internal/issues/record_test.go, flai/internal/adr/adr.go, flai/internal/storygit/commit.go, flai/internal/storygit/commit_test.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/widen_test.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go, flai/internal/mcpserver/adr.go, flai/internal/mcpserver/adr_test.go, flai/internal/mcpserver/folder.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/continuous-improvement.md, design/conventions/decisions.md, template/root/design/conventions/continuous-improvement.md, template/root/design/conventions/decisions.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, design/system/dashboard-host-channel.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0252, S-0245]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 41
  by: planner-E-0017
  at: 2026-10-06T11:36:19Z
forecast:
  duration: 49m
  delivery: 2026-10-07T09:39:00Z
  basis: flai forecast's 83 s per unit over 25 done large improvement stories on claude-opus-5-5, times size 35 (3 criteria, 32 touches), 29th in the pull order with an in-progress limit of 3.
  by: planner-S-0275
  at: 2026-10-06T22:53:55Z
finalized:
  by: alex
  at: 2026-10-06T22:49:06Z
---
# S-0275 Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch

## Goal

Recording friction and decisions costs several turns each: `flai issue bump` or `new`, then `git add` and `commit`, then `flai touches` to claim the issue file; `flai adr new` the same, and ADR and issue numbers taken from the worktree alone collide between parallel branches (S-0245, S-0250, S-0252 fix the numbering). `flai issue bump|new|close` and `flai adr new`, run in a story's worktree with `--commit`, number from the whole repository, commit the files they wrote on the story branch with the story's prefix, and widen the story's touches to them, in one call; over MCP they are `issue_bump`, `issue_new`, `issue_close`, and `adr_new`, and the host channel has the same, so the dashboard can file an issue against a story.

## Acceptance criteria
- [ ] `flai issue bump`, `new`, and `close` and `flai adr new` take `--commit`, which commits what they wrote on the story branch and widens the story's touches to it, in one call, as text and `--json`
- [ ] The same operations exist over MCP and on the host channel
- [ ] The conventions, the harness prompt, `design/system/flai-cli.md`, and the user guide send the agent to the one call

## Tasks
- T-1072 A shared helper commits the files a command wrote on the story branch and widens the story's touches to them
- T-1074 The issues package returns the paths each new, bump, and close wrote
- T-1077 flai issue bump, new, and close and flai adr new take --commit
- T-1083 MCP issue_new and issue_bump take commit, and issue_close and adr_new are new tools
- T-1093 The host channel files, bumps, and closes an issue against a story
- T-1100 The conventions and the harness prompt send a story's agent to the one call with --commit
- T-1105 The design and the user guide describe --commit, the new MCP tools, and the host channel operations

## Notes

Waits for S-0250 and S-0252, which make the numbering safe across branches; this story adds only the single-call shape.

### Planning

`after` is S-0252 and S-0245. S-0250 was cancelled as a duplicate of S-0245, which still holds ADR numbering's fix in the backlog. S-0252, issue numbering, is done. `adr.NextNumber` (`flai/internal/adr/adr.go`) still reads only the checkout.

What the code shows (explorer, 2026-10-06):

- `flai issue` already has `close`.
- `flai adr new` has `--autocommit` and `--trailer`; `flai issue new`, `bump`, and `close` have neither.
- `adr.New` returns the paths it wrote (`Result.Changed`); the `issues` package's `New`, `BumpWith`, `Close`, and `NewOrBump` do not.
- Over MCP, `issue_new`, `issue_bump`, and `issue_story` exist; `issue_close` and `adr_new` do not.
- The host channel has `adr.new`, `adr.accept`, `issue.list`, and `issue.story`, but no `issue.new`, `bump`, or `close`.

Tasks, in four layers (`after` between them):

1. T-1072 (commit-and-widen helper) and T-1074 (issue paths returned).
2. T-1077 (CLI `--commit`) and T-1083 (MCP tools), after both of layer 1.
3. T-1093 (host channel), after T-1077, since both change `flai/cmd/issue.go`; T-1100 (conventions and harness prompt), after T-1077 and T-1083.
4. T-1105 (design and user guide), after T-1077, T-1083, and T-1093.

Touches, all files; no folder touch kept. The three folder touches of the first plan (`flai/internal/issues`, `flai/internal/adr`, `flai/internal/storygit`) are narrowed to the files the tasks name. `flai touches suggest S-0275` was run again: its co-change list adds nothing a task reaches.

- `flai/cmd/issue.go`, `issue_test.go`, `adr.go`, `adr_test.go`: declared and co-change. T-1077 and T-1093.
- `flai/internal/issues/issues.go`, `issues_test.go`, `record.go`, `record_test.go`: narrowed from the declared folder; co-change. T-1074.
- `flai/internal/adr/adr.go`: narrowed from the declared folder. T-1077, only if `Result` needs more.
- `flai/internal/storygit/commit.go`, `commit_test.go`: narrowed from the declared folder; layout, new files. T-1072.
- `flai/internal/itemedit/widen.go`, `widen_test.go`: layout, new. T-1072. `flai touches` widens through `itemedit.ClaimWatch`.
- `flai/internal/mcpserver/issues.go`, `issues_test.go`, `folder.go`: declared. `adr.go` and `adr_test.go`: layout, new. T-1083.
- `flai/internal/hostapi/writes.go`, `writes_test.go`: declared. T-1093.
- `flai/internal/harness/harness.go`, `harness_test.go`, the four convention files, and `template/CHANGELOG.md`: declared. T-1100.
- `design/system/continuous-improvement.md`, `flai-cli.md`, `docs/users/flai.md`, `flai-reference.md`: declared. `design/system/dashboard-host-channel.md`: design, added for the host channel operations. T-1105.
- Left out: `flai/internal/guard`. It gives sub-agents `issue list` only and does not stand in a story agent's way. S-0287 changes which `flai adr` commands a task sub-agent may run, and overlaps here if it widens. Also left out: `flai/internal/docedit` and `flai/internal/check`, co-changed but not reached.
- Overlap to watch: S-0269 (`flai task done`) also widens touches. T-1072 reuses its helper if it lands first.

Forecast: 49m, flai's figure, and it stands. It rose from 39m because the size rose from 26 to 35. Part of that is narrowing three folders to files, which adds no work. The rest is scope the code showed: T-1074, `issue_close` over MCP, and the host channel's issue operations with their `--autocommit`. Delivery 2026-10-07T09:39Z is flai's, after S-0245's.

Cost of delay: 41 USD a week, kept against flai's 90.46. flai shares E-0017's 900 USD a week by forecast time, 39m of 6h28m. planner-E-0017 shared it by the turns each story removes. git shows 94 issue and ADR files added from 2026-09-23 to 2026-10-04, plus bumps. At three turns each (write, commit, claim), that is about 280 turns. The longer forecast changes the work, not the turns removed, so the share stands. The figure assumes three turns a record, which no log classification has measured, so this share is the least certain in the epic.
