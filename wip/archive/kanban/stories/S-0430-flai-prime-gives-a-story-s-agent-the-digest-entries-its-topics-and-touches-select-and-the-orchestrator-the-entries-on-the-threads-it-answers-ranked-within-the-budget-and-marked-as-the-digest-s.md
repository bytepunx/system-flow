---
id: S-0430
type: story
nature: feature
title: flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's
status: done
parent: E-0025
owner: alex
created: 2026-10-16T05:58:39Z
updated: 2026-10-23T21:00:06Z
transitions:
  - to: ready
    at: 2026-10-23T19:11:25Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-23T19:12:05Z
    by: agent-S-0430
  - to: review
    at: 2026-10-23T20:47:22Z
    by: agent-S-0430
  - to: done
    at: 2026-10-23T21:00:06Z
    by: orchestrator
tags: [cli, template]
topics: [cli, conventions, orchestration]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/conventions/work-management.md, design/system/agent-context.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/context/context.go, flai/internal/context/design.go, flai/internal/context/digest.go, flai/internal/context/digest_test.go, flai/internal/context/items.go, flai/internal/context/rank.go, flai/internal/context/story.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go, flai/internal/mcpserver/server_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/template.yaml, flai/internal/inbox/inbox.go, flai/internal/threads/threads.go, flai/internal/mcpserver/server.go]
after: [S-0429]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3338
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 6
      hand_edits: 4
      work: 86
  models:
    - model: claude-opus-5-5
      input: 442
      output: 173087
      cache_read: 32458775
      cache_write: 873755
      cost: 15.2786
  strategic:
    - kind: orchestrator
      seconds: 2210
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 228
          output: 3842
          cache_read: 56977739
          cache_write: 75818
          cost: 12.5487
        - model: claude-sonnet-5-5
          input: 10
          output: 46
          cache_read: 150417
          cache_write: 61631
          cost: 0.1461
cost_of_delay:
  value: 93.33
  by: planner-S-0430
  at: 2026-10-23T16:52:11Z
forecast:
  duration: 1h3m
  delivery: 2026-10-24T07:53:00Z
  basis: "Its own forecast of 1h3m; 21st in the pull order with an in-progress limit of 5, behind S-0429, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424 and S-0425."
  by: flai
  at: 2026-10-23T17:56:28Z
finalized:
  by: orchestrator
  at: 2026-10-23T16:53:02Z
---
# S-0430 flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's

## Goal

What a past story decided, or the operator answered, reaches the next agent that meets the same question without anyone writing it into a design document first. `flai prime --story` and the MCP `prime` add a digest step: the active and compacted entries whose topics meet the story's topics or whose paths meet its claim, ranked against the story with the rest within the budget, each marked as the digest's with its reason. An absorbed entry is never primed. `flai prime --role orchestrate` gives the orchestrator the `answer` and `recommendation` entries ranked against the threads awaiting the operator, so that its recommendation can cite the answer the operator gave last time as its `source`. The conventions tell an agent to cite an entry's ID when it relies on one, which is what keeps an entry from being compacted.

## Acceptance criteria

- [x] An ADR, refining ADR-0049 and ADR-0075, records the digest step: which entries a story's pack selects (topics, paths against the claim), where they rank among the ranked sections, how a compacted entry loads (its one-line title, as a brief), that an absorbed entry never loads, how much of the budget the step may take, and what the orchestrator's pack selects; the sub-agents' packs take entries as briefs only.
- [x] `flai/internal/context` adds the step: each entry carries `step: digest`, its reason (`topics: ...`, `path: ...`, `rank n`), and its size in text and `--json`; a test packs a fixture story whose topics select one entry, whose touches select another, and whose budget leaves out a third, and a fixture with an absorbed entry that never loads.
- [x] `flai prime --role orchestrate`, as text, `--json`, and the MCP `prime` with role `orchestrate`, includes the answer and recommendation entries ranked against the titles and first entries of the threads awaiting the operator, within the budget; a test shows an entry ranked first for a thread that asks its question again.
- [x] A `source` on `thread_reply` may name an entry's path, and the orchestrator's prompt and the convention `strategic-agents.md` say to cite the entry when one answers the question.
- [x] The convention `work-management.md` (template and here) and `design/system/agent-narrative.md` say an agent cites `D-nnnn` in its narrative or thread when it relies on an entry, and asks the operator only what no entry answers; `design/system/agent-context.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` describe the step.
- [x] Measured on this repository after the seed: the pack sizes of three archived stories with and without the step, recorded in `agent-context.md`, and no pack over its budget by the step alone.

## Tasks
- T-1818 An ADR records the pack's digest step: what a story's pack and the orchestrator's select, how entries rank and load, and what part of the budget they may take
- T-1819 A story's pack and a sub-agent's add the digest step: entries selected by topics and claim, ranked within the budget, each marked as the digest's
- T-1820 The orchestrator's pack adds the answer and recommendation entries ranked against the threads awaiting the operator
- T-1821 flai prime and the MCP prime carry the digest step, a thread_reply source may cite an entry, and the orchestrator's prompt says to cite the answer the operator gave last time
- T-1822 The conventions tell an agent to cite a digest entry it relies on and to ask only what no entry answers, and agent-context.md describes and measures the step
- T-2090 The digest step's over-budget rule follows TH-0624, and agent-context.md records the packs of three archived stories with and without the step

## Notes

- Drafted by the planner for E-0025, after S-0429 so that there are entries to prime with and measure.
- Parts (ADR-0104) are unchanged: entries are items like the rest and go into parts as they fit.
- The epic names a convention `agent-narrative.md`, which does not exist. The narrative's rules live in `work-management.md` and `design/system/agent-narrative.md`, so criterion 5 names those.

### Planning

Planned again on 2026-10-23 (planner-S-0430) against the tasks of 2026-10-16 and the code as it stands, with S-0428's `flai/internal/digest` on its branch. The five tasks and their layers stand; five touches were added, and the figures were written.

- **Touches.** The story declared none at first; every touch is a task's, file by file (29 in all).
  - **Declared, kept:** the 24 touches of the first plan.
  - **From layout, added:** `flai/internal/context/rank.go` (T-1819), where `Candidates` builds the BM25 index the ranked step uses and the digest step's ranking goes beside it; `flai/cmd/prime.go` and `flai/internal/mcpserver/folder.go` (T-1821), whose `--help` text and `prime` tool description name the pack's steps in prose and must name the digest step too; `docs/users/flai-reference.md` (T-1821), which `make flai-reference` regenerates from `prime.go`'s help text.
  - **From design, added:** `design/system/flai-cli.md` (T-1822), whose `flai prime` row describes the pack and changes with it.
  - **Co-change, kept out:** `flai/internal/mcpserver/server.go` (8%), since `threads.ResolveAnchor` accepts any repository path that exists, so an entry's path already passes as a `source` and T-1821 only tests it; `flai/internal/threads/threads.go`, since a bare `D-nnnn` as a source is not a criterion; `flai/internal/conventions/conventions.go`, since no role is added; `budget.go` and `parts_test.go`, since the budget and the parts are unchanged; `design/system/project-manifest.md`, `docs/operators/settings.md`, and `flai/internal/manifest`, on the assumption that the step's share of the budget is a constant in `flai/internal/context` beside `DefaultBudget`, not a manifest key (T-1818 decides; if it chooses a key, `flai task done` widens the touches).
  - **Folder touch kept:** `design/adrs` (T-1818), because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** flai's 1h3m kept: median 107 s per unit over 85 done feature stories on claude-opus-5-5 in the large band, size 35 (6 criteria, 29 touches). The five touches added raised it from 54m; four of them are one-sentence edits, but an ADR, three code tasks with their tests, and a docs task that measures three packs fill the hour, so nothing was taken off. Delivery 2026-10-24T07:40Z is flai's play-out, 22nd in the pull order.
- **Cost of delay.** 93.33 USD a week, flai's figure: the story's share of E-0025's 600 USD a week (`time_lost_per_cycle: 4h`, set by planner-S-0428 after the operator confirmed TH-0476), 1h3m of 6h45m forecast over the epic's 8 open stories without inputs. Not adjusted: the story's weight in the epic is its duration, and no input of its own was given.

### Accepted by the orchestrator

- Verified: 8842e6999e6254137f18ff03acdc3833e7a30325
- Batch: 20261023-204831
- Enqueued verified: 8842e6999e6254137f18ff03acdc3833e7a30325
- At: 2026-10-23T21:00:06Z

Verdict: pass. All six criteria are met and ticked, the code and ADR-0194 follow the operator's answer A on TH-0624, and nothing breaks a convention.
- 1: design/adrs/0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md, design/adrs/README.md
- 2: flai/internal/context/digest.go, flai/internal/context/digest_test.go, flai/internal/context/context.go, flai/internal/context/story.go, flai/internal/context/items.go, flai/internal/context/design.go
- 3: flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/inbox/inbox.go, flai/internal/threads/threads.go, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/mcpserver/prime_test.go, flai/internal/mcpserver/server_test.go
- 4: flai/internal/mcpserver/server.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md
- 5: design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/agent-narrative.md, design/system/agent-context.md, docs/users/flai.md, template/CHANGELOG.md
- 6: design/system/agent-context.md
