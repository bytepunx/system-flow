---
id: S-0831
type: story
nature: improvement
title: The kluster send-back's texts to the agent do not say it can ask for its own deploy before review
status: in-progress
owner: alex
created: 2026-11-17T08:17:25Z
updated: 2026-11-17T09:02:53Z
transitions:
  - to: ready
    at: 2026-11-17T08:34:20Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T09:02:53Z
    by: agent-S-0831
tags: [flai, serve, kluster]
topics: [kluster]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, design/system/kluster.md, design/issues/I-0357-the-kluster-send-back-s-texts-to-the-agent-do-not-say-it-can-ask-for-its-own.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 17
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 8
          output: 13
          cache_read: 694306
          cache_write: 6963
          cost: 0.0479
        - model: claude-opus-5-5
          input: 43
          output: 2876
          cache_read: 3001208
          cache_write: 42427
          cost: 0.9988
    - kind: orchestrator
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 48
          cache_read: 354304
          cache_write: 7666
          cost: 0.1123
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-17T08:17:25Z
  value: 12.5
  by: planner-S-0831
  at: 2026-11-17T08:29:20Z
forecast:
  duration: 17m
  delivery: 2026-11-17T09:16:00Z
  basis: "Its own forecast of 17m; 1st in the pull order with an in-progress limit of 3, behind S-0752, S-0806 and S-0827."
  by: flai
  at: 2026-11-17T08:34:32Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:33:51Z
---
# S-0831 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review

## Goal

This story remediates [I-0357](../../../design/issues/I-0357-the-kluster-send-back-s-texts-to-the-agent-do-not-say-it-can-ask-for-its-own.md), "The kluster send-back's texts to the agent do not say it can ask for its own deploy before review". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0357 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0357 is closed with `flai issue close I-0357 --reason` saying what fixed it

## Tasks
- T-3533 The kluster send-back's start prompt tells the agent it may deploy its own story before review
- T-3534 The kluster failure thread of a story sent back tells its agent it may deploy the story before review
- T-3536 kluster.md says what the send-back tells the agent of its own deploy, and I-0357 is closed

## Notes

Cost of delay inputs set by flai from I-0357. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T08:42:12Z, 4 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0357's one instance (S-0660's docs review). Two texts tell the agent of a story the kluster sent back what follows its fix:

- `SentBack.then` for `Kluster` in `flai/internal/harness/harness.go`, the agent's start prompt, resumed or fresh.
- `heldText` in `flai/internal/serve/kluster.go`, the sent-back branch, which the deploy and the checks failure threads end with.

Both say only that `flai serve` deploys the story again in place once it is back in review. Since S-0660 ([ADR-0292](../../../design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) §3, [kluster.md](../../../design/system/kluster.md)), the agent may ask for its own deploy while the story is in progress: `kluster_deploy` or `flai kluster deploy S-nnnn --job`, which runs the checks once the deploy passes while `checks` is on. It then waits for the story's `kluster` event. A deploy and checks that pass at the head resolve the send-back's threads. Back in review, the story is deployed again only when its deploy, and its checks while `checks` is on, have not passed at its head. Each text gains that, in place of saying flai serve always deploys again. `heldText`'s branch for a story not sent back stays as it is: that story is still in review, where the agent asks for nothing. `docs/users/flai.md` § the send-back already says it, so it is not touched.

Touches, where each came from:

| Touch | From |
|-------|------|
| `flai/internal/harness/harness.go` | I-0357 names `harness.SentBack`'s kluster text |
| `flai/internal/harness/harness_test.go` | layout: the send-back prompt tests (`TestASentBackAgentIsToldTheThread` and those of each kind) have no kluster case |
| `flai/internal/serve/agents_test.go` | layout: `TestAStoryTheKlusterSentBackGetsItsAgentOnce` quotes the text |
| `flai/internal/serve/kluster.go` | I-0357 names `heldText` |
| `flai/internal/serve/kluster_test.go` | layout: `TestAFailedDeployOpensOneThreadSendsTheStoryBackAndAPassResolvesIt` and `TestFailingChecksAgainstTheNamespaceSendTheStoryBack` read the thread's text |
| `design/system/kluster.md` | design: § The send-back (S-0694) says what the thread's entry and `harness.SentBack.Kluster` tell the agent |
| `design/issues/I-0357-…md` | criterion 2 closes it |

`flai touches suggest` from these paths lists co-changed files, led by `flai/internal/serve/agents.go` (27%), `design/system/flai-cli.md`, and `docs/users/flai.md` (13–14%). None is added: neither text is built in `agents.go`, and neither `flai-cli.md` nor `docs/users/flai.md` quotes them. No folder touch. No `after`: the story reads nothing an open story writes. S-0810 changes `kluster.go` too, and a trial merge judges that.

Figures:

- Forecast 17m, delivered 2026-11-17T15:41:00Z: `flai forecast`'s figure, kept. Its median, 109 s per unit of size over 39 done medium-band improvement stories, times size 9, fits two texts, their tests, and one design paragraph.
- Cost of delay 12.50 USD a week: `flai cod` from flai's input, 5m lost per 168h cycle at 150 USD an hour, kept as worked out.
