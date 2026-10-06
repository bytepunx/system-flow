---
id: S-0275
type: story
nature: improvement
title: "Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:33Z
updated: 2026-10-06T20:20:49Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions]
touches: [flai/cmd/issue.go, flai/cmd/issue_test.go, flai/cmd/adr.go, flai/cmd/adr_test.go, flai/internal/issues, flai/internal/adr, flai/internal/storygit, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go, flai/internal/mcpserver/folder.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/continuous-improvement.md, design/conventions/decisions.md, template/root/design/conventions/continuous-improvement.md, template/root/design/conventions/decisions.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0252, S-0245]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 41
  by: planner-E-0017
  at: 2026-10-06T11:36:19Z
forecast:
  duration: 39m
  delivery: 2026-10-07T09:42:00Z
  basis: "Its own forecast of 39m; 32nd in the pull order with an in-progress limit of 3, behind S-0223, S-0224, S-0227, S-0229, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0261, S-0264, S-0265, S-0269, S-0270, S-0271, S-0272, S-0273 and S-0274."
  by: flai
  at: 2026-10-06T20:20:49Z
---
# S-0275 Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch

## Goal

Recording friction and decisions costs several turns each: `flai issue bump` or `new`, then `git add` and `commit`, then `flai touches` to claim the issue file; `flai adr new` the same, and ADR and issue numbers taken from the worktree alone collide between parallel branches (S-0245, S-0250, S-0252 fix the numbering). `flai issue bump|new|close` and `flai adr new`, run in a story's worktree with `--commit`, number from the whole repository, commit the files they wrote on the story branch with the story's prefix, and widen the story's touches to them, in one call; over MCP they are `issue_bump`, `issue_new`, `issue_close`, and `adr_new`, and the host channel has the same, so the dashboard can file an issue against a story.

## Acceptance criteria
- [ ] `flai issue bump`, `new`, and `close` and `flai adr new` take `--commit`, which commits what they wrote on the story branch and widens the story's touches to it, in one call, as text and `--json`
- [ ] The same operations exist over MCP and on the host channel
- [ ] The conventions, the harness prompt, `design/system/flai-cli.md`, and the user guide send the agent to the one call

## Tasks

## Notes

Waits for S-0250 and S-0252, which make the numbering safe across branches; this story adds only the single-call shape.

### Planning

`after` changed from S-0250 and S-0252 to S-0252 and S-0245. S-0250 was cancelled as a duplicate of S-0245, which still holds ADR numbering's fix in the backlog. S-0252, issue numbering, is done. `adr.NextNumber` (`flai/internal/adr/adr.go`) still reads only the checkout.

Touches, none declared before. `flai touches suggest S-0275` was seeded with `flai/cmd/issue.go` and `flai/cmd/adr.go`, which 9 commits changed:

- `flai/cmd/issue.go`, `issue_test.go`, `adr.go`, `adr_test.go`: declared seeds and co-change (`issue_test.go` 5 of 9, `adr_test.go` 2 of 9). The `--commit` flag.
- `flai/internal/issues`, `flai/internal/adr`: co-change (3 and 2 of 9). Writing returns the paths written.
- `flai/internal/storygit`: layout. Committing on the story branch with its prefix, beside `FolderNames`.
- `flai/internal/mcpserver/issues.go`, `issues_test.go`, `folder.go`: co-change and layout. The issue tools exist; `adr_new` is new.
- `flai/internal/hostapi/writes.go`, `writes_test.go`: layout (criterion 2).
- `flai/internal/harness/harness.go`, `harness_test.go`: design. The prompt names `flai issue new` and `bump`.
- `design/conventions/continuous-improvement.md`, `decisions.md`, their `template/root` copies, and `template/CHANGELOG.md`: layout. These name `flai issue` and `flai adr new`.
- `design/system/continuous-improvement.md`, `flai-cli.md`, `docs/users/flai.md`, `flai-reference.md`: co-change (`flai-cli.md` and `flai.md` 4 of 9) and criterion 3.
- Left out: `flai/internal/guard`. S-0287 changes which `flai adr` commands a task sub-agent may run, and it overlaps here if it widens. Also left out: `flai/internal/docedit` and `flai/internal/check`, co-changed but not reached.

Forecast: flai gave 39m (89 s per unit over 21 done large improvement stories, times size 26), and it stands. `--commit` is one shared helper over four commands, plus an MCP `adr_new`. The delivery, 2026-10-07T02:10Z, is flai's and comes after S-0245's.

Cost of delay: 41 USD a week, against flai's 102.63. This is E-0017's 900 USD a week shared by the turns each story removes. The epic's evidence gives no count for issue and ADR turns. git shows 94 issue and ADR files added from 2026-09-23 to 2026-10-04, plus bumps. At three turns each (write, commit, claim) that is about 280 turns. The figure assumes three turns a record, which no log classification has measured, so this share is the least certain in the epic.
