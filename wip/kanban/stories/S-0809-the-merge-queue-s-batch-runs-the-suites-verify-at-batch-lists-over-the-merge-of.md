---
id: S-0809
type: story
nature: feature
title: The merge queue's batch runs the suites verify_at.batch lists over the merge of its stories
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:05Z
updated: 2026-11-16T23:42:05Z
transitions: []
tags: [flai]
topics: [verify, workflow]
after: [S-0806]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0809 The merge queue's batch runs the suites verify_at.batch lists over the merge of its stories

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §6. Today a batch runs every tier its paths select.

## Acceptance criteria

- [ ] A batch and its bisection run the tiers of the suites `verify_at.batch` lists that the batch's paths select, and every such tier under `all`.
- [ ] `flai queue` and the batch's record name the suites a batch ran.
- [ ] On the baseline a batch runs what it ran before this story.

## Tasks

## Notes
