---
id: S-0814
type: story
nature: feature
title: The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:10Z
updated: 2026-11-16T23:42:10Z
transitions: []
tags: [flaiover]
topics: [verify, workflow]
after: [S-0810]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §4. The review page's verification section shows the last report's tiers, and after S-0804 a Rebase and resolve job's stream.

## Acceptance criteria

- [ ] The verification section groups the report's tiers by suite and says which suites the story's next transitions need and which have not passed at its head.
- [ ] A job a state's `runs:` started is followed in the collapsible activity stream S-0804 adds, with its current step, its skipped cluster steps and why, and how it ended.
- [ ] The host API gives the dashboard the suites of each transition, state, and point, and the board's column for a state with `runs:` says what it runs.
- [ ] The component examples and the dashboard design are in step.

## Tasks

## Notes
