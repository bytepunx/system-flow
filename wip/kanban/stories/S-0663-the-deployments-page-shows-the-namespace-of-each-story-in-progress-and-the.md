---
id: S-0663
type: story
nature: feature
title: The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:11Z
updated: 2026-11-10T08:09:11Z
transitions: []
tags: [dashboard]
after: [S-0660]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it

## Goal

The dashboard's Deployments page shows the namespace of each story in progress beside those in review, and the operator deploys, checks, and tears down from it. The operator sees what agents have running in the cluster and can clear it.

## Acceptance criteria

- [ ] The Deployments page lists each story namespace with its story's column, its commit, whether that is the branch's head, its last deploy and checks run with their outcomes and logs, and who asked for each
- [ ] The operator deploys a story in progress, runs its checks, and tears it down from the page, each through flai, behind the `kluster` host action
- [ ] The page shows how many story namespaces live against the limit, and which request was refused for it
- [ ] `flai kluster deployments` and its `--json` report the same
- [ ] A port-forward to a story in progress starts and stops as one to a story in review does
- [ ] The component examples under `design/components/` are in step with the components changed
- [ ] `design/system/kluster.md`, `design/system/flaiover-dashboard.md`, and `docs/users/flaiover.md` say all of it

## Tasks

## Notes

- S-0581 built the page.
