---
id: S-0605
type: story
nature: improvement
title: settings.get finds the projects below the import folders once a request, not once a served project, so it answers in under 100 ms, not 825 ms
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:40:46Z
updated: 2026-11-07T22:40:46Z
transitions: []
tags: [cli]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_project.go, flai/internal/serve/below.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0605 settings.get finds the projects below the import folders once a request, not once a served project, so it answers in under 100 ms, not 825 ms

## Goal

Cause 15 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `settings.get` has 4,104 slow calls in 4 days at p50 825 ms. `servedView` (`flai/cmd/serve_actions.go`) asks `importRootOf` (`flai/cmd/serve_project.go`) for each of the 13 served projects, and each runs `serve.FindBelow` (`flai/internal/serve/below.go`), 45 to 85 ms: 13 walks of the import folders and 169 manifest loads a request. It takes no context, so it records no phases.

## Acceptance criteria

- [ ] `settings.get` works out the projects below the import folders once a request, or takes what the running flai serve already found, and answers what it answered before, which a test with three served projects under two import folders shows.
- [ ] `settings.get` records its phases.
- [ ] `settings.get` on the operator's host answers under 100 ms, and server-performance.md records the before and after.

## Tasks

## Notes
