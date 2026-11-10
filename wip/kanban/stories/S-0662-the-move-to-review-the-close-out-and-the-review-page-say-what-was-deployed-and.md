---
id: S-0662
type: story
nature: feature
title: The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:10Z
updated: 2026-11-10T08:09:10Z
transitions: []
tags: [cli, dashboard, template]
after: [S-0661]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first

## Goal

The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there before review. A criterion that needs the running system reaches review ticked, or unticked with a reason the operator can weigh.

## Acceptance criteria

- [ ] While the `kluster` host action is on and a `deploy` is set, `flai move <story> review` and `item_move` warn when the story's namespace is not deployed at its branch's head, or its checks did not pass there, naming the commands that do both; they do not refuse
- [ ] Their warning for unticked acceptance criteria says, in such a project, whether the namespace was deployed and checked at the head
- [ ] The close-out records the namespace's last deploy and checks run, each with its commit and outcome, and `flai verify --last` and the MCP `item_get` show them
- [ ] The dashboard's review page shows the same beside the criteria, with a link to each log
- [ ] The orchestrator's acceptance reads the record as evidence for a criterion, as `design/system/strategic-agents.md` says
- [ ] The work-management convention, in `design/system`, in `design/conventions`, and in `template/`, tells a story's agent in a project with a kluster to deploy, check, and tick before review, and to leave unticked only what it could not verify, with the reason
- [ ] A story's prime pack in such a project names the deploy and checks commands
- [ ] `docs/users/flai.md` and `docs/users/flaiover.md` say all of it

## Tasks

## Notes

- `Repo.Move` in `flai/internal/workitem/rules.go` makes the warning for unticked criteria today (S-0282).
