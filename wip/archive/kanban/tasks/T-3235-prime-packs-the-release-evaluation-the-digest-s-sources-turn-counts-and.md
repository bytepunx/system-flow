---
id: T-3235
type: task
nature: feature
title: Prime packs, the release evaluation, the digest's sources, turn counts, and statsread read a state's kind
status: done
parent: S-0749
owner: alex
created: 2026-11-12T16:08:29Z
updated: 2026-11-15T19:57:41Z
transitions:
  - to: ready
    at: 2026-11-15T19:45:27Z
    by: agent-S-0749
  - to: in-progress
    at: 2026-11-15T19:45:27Z
    by: agent-S-0749
  - to: done
    at: 2026-11-15T19:57:41Z
    by: agent-S-0749
stream: S-0749
tags: [cli]
touches: [flai/internal/context/reasons.go, flai/internal/context/reasons_figures.go, flai/internal/context/story.go, flai/internal/release/evaluate.go, flai/internal/digest/sources.go, flai/internal/statsread/statsread.go, flai/internal/context/context_workflow_test.go, flai/internal/release/evaluate_workflow_test.go, flai/internal/statsread/statsread_workflow_test.go]
usage:
  source: log
  seconds: 734
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 77
      output: 24566
      cache_read: 4128826
      cache_write: 214809
      cost: 2.5632
---
# T-3235 Prime packs, the release evaluation, the digest's sources, turn counts, and statsread read a state's kind

## Work

Replace the one to three references to a state's name in each of these files with a read of the kind: the orchestrator's reasons and a story's pack pick their items by kind, the release evaluation counts the stories in accepted states, the digest harvests from accepted and cancelled items, and the turn counts and statsread tell open from closed by kind.

## Done when

These files name no state, and the existing tests of these packages pass unchanged.

## Notes
