---
id: T-3229
type: task
nature: feature
title: The host API's kluster, components, people, import, board, and settings reads tell a story's place by kind
status: backlog
parent: S-0747
owner: alex
created: 2026-11-12T16:08:05Z
updated: 2026-11-12T16:08:05Z
transitions: []
stream: S-0747
tags: [cli]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_forward.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/people.go, flai/internal/hostapi/people_test.go, flai/internal/hostapi/import.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/settings.go]
---
# T-3229 The host API's kluster, components, people, import, board, and settings reads tell a story's place by kind

## Work

Replace the references to a state's name in these reads with a state's kind: the Deployments page's namespaces belong to stories in the wait state the acceptance leaves from, the components' worktrees to stories in active work, and the people page counts open work by kind. Change no answer's shape: the board's columns are S-0753's. Test the kluster and people reads on a workflow with a state added.

## Done when

These files name no state, and `kluster_test.go` and `people_test.go` pass on both workflows.

## Notes
