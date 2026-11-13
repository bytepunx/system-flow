---
id: S-0659
type: story
nature: feature
title: "An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once"
status: done
parent: E-0035
owner: alex
created: 2026-11-10T08:09:07Z
updated: 2026-11-13T07:29:40Z
transitions:
  - to: ready
    at: 2026-11-12T15:59:54Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-12T16:00:03Z
    by: agent-S-0659
  - to: review
    at: 2026-11-13T07:26:47Z
    by: agent-S-0659
  - to: done
    at: 2026-11-13T07:29:40Z
    by: orchestrator
blocked:
  - from: 2026-11-12T16:09:00Z
    until: 2026-11-13T07:25:04Z
    reason: awaiting the operator on TH-1103
tags: [cli]
topics: [kluster, deployment]
touches: [design/system/kluster.md, design/adrs/README.md, design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/issues/I-0354-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 640
  priced_by: harness
  turns:
    - day: 2026-11-12
      ceremony: 2
      work: 31
    - day: 2026-11-13
      ceremony: 2
      work: 11
  models:
    - model: claude-opus-5-5
      input: 102
      output: 41676
      cache_read: 12014694
      cache_write: 622786
      cost: 8.2192
  strategic:
    - kind: planner
      seconds: 81
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 66
          output: 189
          cache_read: 2487345
          cache_write: 136497
          cost: 0.1793
        - model: claude-opus-5-5
          input: 69
          output: 8377
          cache_read: 4438131
          cache_write: 92398
          cost: 1.6379
    - kind: orchestrator
      seconds: 394
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 88
          output: 584
          cache_read: 3441484
          cache_write: 386173
          cost: 1.0699
cost_of_delay:
  value: 263.16
  by: planner-S-0659
  at: 2026-11-12T15:53:54Z
forecast:
  duration: 20m
  delivery: 2026-11-12T19:38:00Z
  basis: "flai forecast: median 116 s per unit of size over 5 done feature stories on claude-opus-5-5, times size 10 (7 criteria, 3 touches); 2nd in the pull order behind S-0683, S-0709, and S-0241 with an in-progress limit of 3."
  by: planner-S-0659
  at: 2026-11-12T15:53:54Z
finalized:
  by: alex
  at: 2026-11-12T07:06:20Z
---
# S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once

## Goal

An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request, so that the stories after it build one design. It amends [ADR-0254](../../../design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md), and [kluster.md](../../../design/system/kluster.md) says the same.

## Acceptance criteria

- [x] An accepted ADR says how a story's agent asks for a deploy and a teardown of its own story's namespace, that flai serve runs the job, and that the `kluster` host action gates it
- [x] It says what the agent may run against the namespace: the project's checks, its own commands with the namespace's variables, and what `flai guard` refuses, another story's namespace and the default one among them
- [x] It says when a namespace is made, deployed again, and torn down across in-progress, review, a send back, acceptance, and cancellation, in a table that replaces ADR-0254 §4's
- [x] It sets how many story namespaces may live at once, where the project sets it, and what a request over the limit gets
- [x] It says how the agent waits for a job behind the queue, what it is told on a failure, and when a failure reaches the operator
- [x] It says what is recorded of a deploy and a checks run before review, and where the move to review, the close-out, and an acceptance read it
- [x] ADR-0254 names it in `superseded_by` or the new ADR names the sections it amends, as `design/conventions/decisions.md` has it, and `design/system/kluster.md` describes the design as decided and marks each part not yet built

## Tasks
- T-2860 Write the ADR, proposed, that amends ADR-0254 §4, §5, §7, and §8 for a deploy a story's agent asks for while its story is in progress
- T-2861 Describe in design/system/kluster.md the deploy a story's agent asks for as decided, marking each part not yet built
- T-2862 Ask the operator to confirm the ADR's decisions on a thread, and accept it with their answer

## Notes

- A design story: it changes no code.
- The hand-off `flai verify --job` makes to flai serve is the pattern to weigh first.

### Planning

- Touches:
  - `design/system/kluster.md`: design. The goal links it and the last criterion has it describe the decided design.
  - `design/adrs/README.md`: co-change (`flai touches suggest` counts 328 of 3801 commits on the story's three files). `flai adr new` and `flai adr accept` write the index again.
  - `design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md`: design. It changes only if the ADR supersedes it and `superseded_by` is set; with `--refines`, as T-2860 expects, it stays as it is.
  - No folder touch on `design/adrs`. The new ADR's file takes the number `flai adr new` gives it when the story runs, and `flai adr new --commit` adds that file to the story's touches itself.
  - Left out from `flai touches suggest`: the docs and `flai-cli.md`, `workflow.md`, `project-manifest.md`, `flaiover-dashboard.md`, and `design/tech/kluster.md`. They describe what is built; S-0660 to S-0663 change them as they build it.
- Since the first plan (2026-11-10), S-0694 built ADR-0254 §5a: a failed deploy or checks run sends the story back to in-progress and keeps its namespace there, and `kluster.md` § What happens when gained its rows. ADR-0259 and ADR-0263 also refine ADR-0254 (deploy optional, components). T-2860 and T-2861 now say to start from these.
- Topics `kluster` and `deployment` were added, the topics of kluster.md and ADR-0254, so the story's pack briefs them.
- No `after`: the story reads no other story's output. S-0660 to S-0663 build on this one's ADR.
- Forecast 20m stands as `flai forecast` gives it, delivery 2026-11-12T19:38Z. S-0576, which wrote ADR-0254 and its design, took 710 s of agent time; this ADR amends four of its sections and edits the same design, so a little more is fair. The wait on the operator in T-2862 is wall time, not agent time.
- Cost of delay 263.16 USD a week, as `flai cod` gives it: S-0659's share by forecast (20m of 1h35m) of E-0035's `revenue_per_week: 1250`, which the operator confirmed on TH-0982. No adjustment: the story is E-0035's design story and its share is fair.
- Assumptions: the limit on story namespaces is set in the manifest's `kluster:` block; the third criterion's "ADR-0254 §4's" table means §5's, the only table of when a namespace is deployed and torn down (§4 is the variables).

### Accepted by the orchestrator

- Verified: df3e026279382bd0e2473609b90e931412d221bd
- At: 2026-11-13T07:29:40Z

Verdict: met at df3e026279382bd0e2473609b90e931412d221bd
- 1: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/system/kluster.md
- 2: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/system/kluster.md
- 3: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/system/kluster.md
- 4: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/system/kluster.md
- 5: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/system/kluster.md
- 6: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/system/kluster.md
- 7: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md, design/adrs/README.md, design/system/kluster.md
