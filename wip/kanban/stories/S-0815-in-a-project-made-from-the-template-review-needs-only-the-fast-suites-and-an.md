---
id: S-0815
type: story
nature: feature
title: In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:11Z
updated: 2026-11-16T23:42:11Z
transitions: []
tags: [flai, template, flaiover]
topics: [verify, workflow]
after: [S-0807, S-0808, S-0809, S-0810, S-0811, S-0812, S-0814, S-0756]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites

## Goal

The end-to-end proof of [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md), as S-0756 is of the epic: the operator's example of 2026-11-16, made in the workflow file alone.

## Acceptance criteria

- [ ] A test makes a project from the template, edits its workflow file so that the move to review needs `static` and `behavior`, adds an `integration` state between `review` and `done` whose `runs:` names `kluster: [deploy, checks]` and the `integration` and `smoke` suites, and sets `publish` to the two slow suites.
- [ ] In that project a story's close-out runs only the fast suites, its entry to `integration` runs the slow ones once with the cluster's steps skipped for want of a cluster, a failing tier there sends it back, and its acceptance from `integration` is not held for a change to main that only the fast suites' paths do not select.
- [ ] The same project with a cluster configured deploys and runs its checks on entering `integration`.
- [ ] Anything that does not follow the file is recorded as an issue or a story under E-0034.
- [ ] The user guide has a worked example of the file.

## Tasks

## Notes
