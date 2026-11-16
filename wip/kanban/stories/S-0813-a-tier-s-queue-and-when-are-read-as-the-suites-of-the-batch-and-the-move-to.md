---
id: S-0813
type: story
nature: feature
title: A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:09Z
updated: 2026-11-16T23:42:09Z
transitions: []
tags: [flai, template]
topics: [verify, workflow]
after: [S-0808, S-0809]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §8. `queue: true` leaves a tier to the batch and `when:` runs it at the close-out anyway; with the file naming the suites of each step, the two keys say the same thing twice.

## Acceptance criteria

- [ ] An ADR says how `queue` and `when` map onto `verify:` and `verify_at.batch`, and what of ADR-0174 and ADR-0233 it supersedes.
- [ ] `flai check` advises a manifest that sets `queue` or `when` of the lines to put in the workflow file in their place.
- [ ] A project that sets neither key and names suites in its file runs what the file says; one that still sets them runs what it ran before, until it removes them.
- [ ] This repository's manifest and file are migrated with no change to what any step runs.

## Tasks

## Notes
