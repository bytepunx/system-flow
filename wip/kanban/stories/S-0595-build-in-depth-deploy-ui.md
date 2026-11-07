---
id: S-0595
type: story
nature: feature
title: Build in-depth Deploy UI
status: backlog
owner: alex
created: 2026-11-07T22:00:05Z
updated: 2026-11-07T22:00:05Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-07T22:00:05Z
---
# S-0595 Build in-depth Deploy UI

## Goal

Under Ops/Deployments, there is only a single button and when pressed, a new project shows the error:

```text
the kluster: block of system-flow.yaml names no deploy command; write kluster.deploy, the argument list that deploys into a namespace, and run flai check
```

The deployment screen should have the controls necessary to set up build and deploy machinery for a project's services and configuration.

## Acceptance criteria
- [ ] The operator never has to edit system-flow.yaml directly to control or setup deployments
- [ ] The UI/UX allow the operator to create or make any changes needed to control their projects deployments to the kluster managed kubernetes cluster
- [ ] The Deployment page should use the vertical nav rail in page to group controls, livestream logs, and other status indicators for a given deployment
- [ ] a Build sub-page should show both the status and an expandable live feed of each component's build. If there are 3 Docker images (for example), there should be 3 different component build panels with status indicators so that the operator can see which steps have completed, are still running, or have failed.

## Tasks

## Notes
