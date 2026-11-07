---
id: S-0586
type: story
nature: improvement
title: Provide first time wizard for kluster setup
status: backlog
owner: alex
created: 2026-11-07T08:24:44Z
updated: 2026-11-07T08:24:44Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-07T08:24:44Z
---
# S-0586 Provide first time wizard for kluster setup

## Goal

When enabling kluster in a project's feature settings, the following error appears at the bottom of the features page:

```text
this project's system-flow.yaml has no kluster: block, so flai runs no cluster for it; add one naming the kluster config and the deploy command, such as kluster: {config: integration/kluster.yaml, deploy: [scripts/kluster-deploy.sh]}, as design/system/kluster.md § The kluster: block describes, then run this again
```

After a refresh, the `kluster` feature is no longer enabled.

This is confusing and goes against the goal of making the dashboard the single place the operator needs access to. If there are hosted versions of flai in the future, we shouldn't assume shell access.

When enabling kluster, provide a scaffold for system-flow.yaml that is either functionally a no-op, or that starts with the empty kubernetes cluster.

The Klusters page should allow the user to manage all features completely without the need for CLI commands or file edits.

## Acceptance criteria
- [ ] When an operator enables the kluster feature, provide a stub/no-op/defaults only block in system-flow.yaml
- [ ] The Klusters page should have everything the user needs to set up, manage, observe, and control the kluster for that project
- [ ] There are no operations, configuration, or commands the operator needs to issue outside of the dashboard to make full use of this feature

## Tasks

## Notes
