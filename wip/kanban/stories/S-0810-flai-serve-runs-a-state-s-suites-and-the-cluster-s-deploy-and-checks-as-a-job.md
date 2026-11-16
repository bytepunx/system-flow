---
id: S-0810
type: story
nature: feature
title: flai serve runs a state's suites and the cluster's deploy and checks as a job when a story enters it, skips the cluster's steps where none is configured, and sends the story back when a step fails
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:06Z
updated: 2026-11-16T23:42:06Z
transitions: []
tags: [flai]
topics: [verify, workflow]
after: [S-0806, S-0804]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0810 flai serve runs a state's suites and the cluster's deploy and checks as a job when a story enters it, skips the cluster's steps where none is configured, and sends the story back when a step fails

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §4 and §5. The job is the one S-0804 makes for Rebase and resolve; this story starts it from a state's `runs:`. Today the cluster's deploy and checks run on entering any state of kind `waiting`, in code.

## Acceptance criteria

- [ ] When a story enters a state with `runs:`, flai serve rebases its branch onto the main branch and runs the state's suites and `kluster:` steps as a job, with its output kept and its steps recorded as they start and end.
- [ ] A job that passes leaves the story in the state; a rebase that stops on conflicts or a step that fails sends it back by the state's `send-back` transition, with the thread its agent is started on.
- [ ] A `kluster:` step is skipped, and the record says so and why, when the project configures no deploy, the kluster host action is off, or no agent namespace is available, and it never holds the story.
- [ ] The cluster's deploy and checks run where a state's `runs:` names them and on no other state, and the `kluster-verified` rule is evaluated by the move where a transition names it.
- [ ] On the baseline the cluster's deploy and checks run on entering `review` as before this story.

## Tasks

## Notes
