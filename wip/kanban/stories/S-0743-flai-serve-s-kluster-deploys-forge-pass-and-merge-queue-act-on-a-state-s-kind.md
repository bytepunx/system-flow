---
id: S-0743
type: story
nature: feature
title: flai serve's kluster deploys, forge pass, and merge queue act on a state's kind and a transition's effect instead of review and done
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:52Z
updated: 2026-11-12T15:59:52Z
transitions: []
tags: [cli]
topics: [cli, kluster]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/forge.go, flai/internal/serve/forge_test.go, flai/internal/serve/forgemerge.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/state.go, flai/internal/kluster/deploy.go, flai/internal/forge/state.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0743 flai serve's kluster deploys, forge pass, and merge queue act on a state's kind and a transition's effect instead of review and done

## Goal

What flai serve does around review follows the workflow. A story is deployed into its kluster namespace, pushed to the forge with its pull request, and taken by the merge queue when it is in a state of the wait kind the acceptance leaves from, and a send back takes it to the state the workflow's send-back transition names, so a workflow with a state of its own before `review` is not deployed, pushed, or queued there.

## Acceptance criteria

- [ ] `flai/internal/serve/kluster.go`, `forgestatus.go`, `forge.go`, `forgemerge.go`, and `mergequeue.go`, `flai/internal/mergequeue`, `flai/internal/kluster/deploy.go`, and `flai/internal/forge/state.go` name no state: each reads a state's kind or a transition's effect from the workflow
- [ ] A story is deployed, pushed with its pull request, and eligible for the merge queue in the state the acceptance transition leaves from, and torn down on the transitions with the acceptance or the cancellation effect, as in `review`, `done`, and `cancelled` today
- [ ] A failing required check, a requested change, a bisected batch, and a failed deploy send the story back along the transition with the send-back effect, with today's reason and thread
- [ ] Tests run the deploy, the forge pass, eligibility, and the send back against a workflow with a state renamed and one added between `in-progress` and `review`, with the outcome the kinds and effects give
- [ ] With the baseline, every existing `serve`, `mergequeue`, `kluster`, and `forge` test passes unchanged

## Tasks

## Notes

- `serve/kluster.go` holds 16 of serve's references, `forgestatus.go` 6, `mergequeue.go` 5; `mergequeue` holds 7.
