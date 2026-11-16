---
id: S-0812
type: story
nature: feature
title: flai test --at ci runs the suites verify_at.ci lists, and the template's and this repository's CI call it
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:08Z
updated: 2026-11-16T23:42:08Z
transitions: []
tags: [flai, template]
topics: [verify, workflow]
after: [S-0806]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0812 flai test --at ci runs the suites verify_at.ci lists, and the template's and this repository's CI call it

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §6. Today the CI workflow files hold commands of their own, apart from the manifest's tiers.

## Acceptance criteria

- [ ] `flai test --at <point>` runs the tiers of the suites the file lists for `ci`, `batch`, or `publish`, over all paths.
- [ ] The template's CI workflow calls `flai test --at ci` once flai is installed in it, and a project made from the template runs its suites in CI with no command of its own.
- [ ] This repository's CI workflows call it where a suite stands for what the job ran, and a note in the design says which jobs keep commands of their own and why.
- [ ] The user guide says how a project's CI runs its suites.

## Tasks

## Notes
