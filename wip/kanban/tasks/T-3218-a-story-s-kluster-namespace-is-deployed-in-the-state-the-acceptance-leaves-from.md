---
id: T-3218
type: task
nature: feature
title: A story's kluster namespace is deployed in the state the acceptance leaves from and torn down on acceptance or cancellation, by effect
status: backlog
parent: S-0743
owner: alex
created: 2026-11-12T16:07:12Z
updated: 2026-11-12T16:07:12Z
transitions: []
stream: S-0743
tags: [cli]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
---
# T-3218 A story's kluster namespace is deployed in the state the acceptance leaves from and torn down on acceptance or cancellation, by effect

## Work

Replace the 16 references to a state's name in `serve/kluster.go` and the one in `kluster/deploy.go`: a story is deployed when it enters the state the acceptance transition leaves from, torn down on a transition with the acceptance or the cancellation effect, and sent back along the send-back transition when a deploy or check fails (ADR-0254 §5a). Test the deploy, the teardown, and the send back on a workflow with a state added between `in-progress` and `review`.

## Done when

Both files name no state, and their tests pass on the baseline and on the other workflow.

## Notes
