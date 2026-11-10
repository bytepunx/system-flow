---
id: S-0660
type: story
nature: feature
title: A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:08Z
updated: 2026-11-10T08:09:08Z
transitions: []
tags: [cli]
after: [S-0659]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled

## Goal

A story's agent asks flai serve to deploy its story from its worktree while the story is in progress, and the namespace stays until the story is accepted or cancelled. The agent has the running system to work against, and runs no kluster command that writes.

## Acceptance criteria

- [ ] With the `kluster` host action on and a `deploy` set, a story's agent asks for a deploy of its own story in progress with one flai command and one MCP tool, and flai serve builds in the story's worktree at its branch's head, loads, and deploys into `s-nnnn`, as it does for a story in review
- [ ] The request waits for the job and answers with its outcome, the commit, each step's result, and the log's path; a request that outlasts its wait says the job is still queued or running and how to read it later
- [ ] A request is refused, with the reason and what to do, for another story, a story not in progress or review, uncommitted changes in the worktree, the action off, no `deploy`, and a request over the limit on story namespaces
- [ ] `flai guard` still refuses every agent `flai kluster up`, `down`, `deploy`, `teardown`, and `promote`, and refuses a sub-agent the new request
- [ ] A story in progress keeps its namespace: flai serve tears it down only when the story is accepted or cancelled, or when its agent or the operator asks
- [ ] A story that enters review at the head already deployed is not deployed again, and one at a new head is deployed in place
- [ ] A failed deploy asked for by the agent opens no thread to the operator: the agent gets the step, the exit code, and the log
- [ ] Each job the agent asked for is in the host journal and the namespace's record with who asked
- [ ] `design/system/kluster.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` and its reference say all of it

## Tasks

## Notes

- `flai/internal/serve/kluster.go` is the pass, and `flai/internal/guard/guard_rules.go` holds `klusterWrites`.
