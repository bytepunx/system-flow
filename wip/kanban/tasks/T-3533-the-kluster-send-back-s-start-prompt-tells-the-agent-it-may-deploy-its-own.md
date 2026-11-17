---
id: T-3533
type: task
nature: improvement
title: The kluster send-back's start prompt tells the agent it may deploy its own story before review
status: done
parent: S-0831
owner: alex
created: 2026-11-17T08:29:42Z
updated: 2026-11-17T09:14:18Z
transitions:
  - to: ready
    at: 2026-11-17T09:07:03Z
    by: agent-S-0831
  - to: in-progress
    at: 2026-11-17T09:07:03Z
    by: agent-S-0831
  - to: done
    at: 2026-11-17T09:14:18Z
    by: agent-S-0831
stream: S-0831
tags: [flai, harness, kluster]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 435
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 13
      cache_read: 575419
      cache_write: 1599
      cost: 0.2201
---
# T-3533 The kluster send-back's start prompt tells the agent it may deploy its own story before review

## Work

Rewrite the `Kluster` case of `SentBack.then` in `flai/internal/harness/harness.go`, which both the resumed and the fresh send-back prompts end with. Keep what it says of the thread and of a cause outside the story. Add that, once the change is fixed and committed, the agent may ask for its own deploy at the branch's head with the MCP tool `kluster_deploy` or `flai kluster deploy <story> --job`, which runs the checks once the deploy passes while `checks` is on, and wait for the story's `kluster` event. Add that a deploy and checks that pass at the head resolve the thread, and that back in review flai serve deploys the story again in place only when they have not passed at its head ([kluster.md](../../../design/system/kluster.md) § The send-back, ADR-0292 §3). Keep the sentences short.

Add a kluster case beside `TestAnAgentAnAcceptanceSentBackIsToldItIsAcceptedAgain` in `flai/internal/harness/harness_test.go`. It reproduces I-0357: the prompt of a story sent back with `Kluster`, resumed and fresh, names `kluster_deploy` and `flai kluster deploy <story> --job`. Update the text `TestAStoryTheKlusterSentBackGetsItsAgentOnce` in `flai/internal/serve/agents_test.go` expects.

It waits for no task: it shares no path with the thread's text.

## Done when

- The kluster send-back prompt, resumed and fresh, names `kluster_deploy` and `flai kluster deploy <story> --job`, says a pass at the head resolves the thread, and says the review deploys again only what has not passed at the head.
- The new harness test fails on the old text and passes on the new one, and `flai test` passes on the three files.

## Notes
