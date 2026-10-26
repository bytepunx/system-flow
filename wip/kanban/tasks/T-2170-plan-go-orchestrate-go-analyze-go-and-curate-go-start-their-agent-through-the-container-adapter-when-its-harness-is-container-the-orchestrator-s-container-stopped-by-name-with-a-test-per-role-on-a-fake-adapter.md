---
id: T-2170
type: task
nature: feature
title: plan.go, orchestrate.go, analyze.go, and curate.go start their agent through the container adapter when its harness is container, the orchestrator's container stopped by name, with a test per role on a fake adapter
status: backlog
parent: S-0509
owner: alex
created: 2026-10-26T04:50:48Z
updated: 2026-10-26T04:50:48Z
transitions: []
stream: S-0509
tags: [cli, docker]
touches: [flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2163]
---
# T-2170 plan.go, orchestrate.go, analyze.go, and curate.go start their agent through the container adapter when its harness is container, the orchestrator's container stopped by name, with a test per role on a fake adapter

## Work

The four strategic start paths call `harness.For` and `startOf` already; make each pass the container adapter what it needs: the role's cache folder created under `.flai-cache/agents/<role>` before the start, the request's `Role`, `Item`, and `Focus` kept, `FLAI_ROLE` passed through `-e`, and the refusals of `startAllowed` (docker, image, login) applied as for a story. In `agents.go`, make the stop of a strategic run stop its container by name as T-2162 does for a story, so the orchestrator's container ends when the `orchestrate` host action is turned off, and the activity log reader (`activity.go` reads the claude-code reader directly) read a container run's log through the inner harness's reader. Test each role in its `_test.go` with a fake adapter that records the request and in `agents_test.go` that the orchestrator's stop reaches `docker stop`. Waits for T-2163, which gives the adapter the role's mounts.

## Done when

Each of the four roles starts through the container adapter in its test with the role's request and cache folder, the orchestrator's stop stops the container, and `flai test flai/internal/serve` passes.

## Notes

`claudecheck.go` probes the host's `claude` before a start; for `harness: container` the ADR of T-2158 says whether it probes the image's instead, and this task follows it.
