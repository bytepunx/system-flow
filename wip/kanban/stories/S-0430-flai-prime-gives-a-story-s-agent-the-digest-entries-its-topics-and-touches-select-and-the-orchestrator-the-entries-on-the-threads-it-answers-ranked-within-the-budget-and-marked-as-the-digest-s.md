---
id: S-0430
type: story
nature: feature
title: flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:58:39Z
updated: 2026-10-21T07:58:04Z
transitions: []
tags: [cli, template]
topics: [cli, conventions, orchestration]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/conventions/work-management.md, design/system/agent-context.md, design/system/agent-narrative.md, docs/users/flai.md, flai/cmd/prime_test.go, flai/internal/context/context.go, flai/internal/context/design.go, flai/internal/context/digest.go, flai/internal/context/digest_test.go, flai/internal/context/items.go, flai/internal/context/story.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/mcpserver/prime_test.go, flai/internal/mcpserver/server_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/template.yaml]
after: [S-0429]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 48m
  delivery: 2026-10-21T22:55:00Z
  basis: "Its own forecast of 48m; 27th in the pull order with an in-progress limit of 5, behind S-0456, S-0458, S-0463, S-0464, S-0465, S-0459, S-0461, S-0462, S-0460, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428 and S-0429."
  by: flai
  at: 2026-10-21T07:58:04Z
---
# S-0430 flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's

## Goal

What a past story decided, or the operator answered, reaches the next agent that meets the same question without anyone writing it into a design document first. `flai prime --story` and the MCP `prime` add a digest step: the active and compacted entries whose topics meet the story's topics or whose paths meet its claim, ranked against the story with the rest within the budget, each marked as the digest's with its reason. An absorbed entry is never primed. `flai prime --role orchestrate` gives the orchestrator the `answer` and `recommendation` entries ranked against the threads awaiting the operator, so that its recommendation can cite the answer the operator gave last time as its `source`. The conventions tell an agent to cite an entry's ID when it relies on one, which is what keeps an entry from being compacted.

## Acceptance criteria

- [ ] An ADR, refining ADR-0049 and ADR-0075, records the digest step: which entries a story's pack selects (topics, paths against the claim), where they rank among the ranked sections, how a compacted entry loads (its one-line title, as a brief), that an absorbed entry never loads, how much of the budget the step may take, and what the orchestrator's pack selects; the sub-agents' packs take entries as briefs only.
- [ ] `flai/internal/context` adds the step: each entry carries `step: digest`, its reason (`topics: ...`, `path: ...`, `rank n`), and its size in text and `--json`; a test packs a fixture story whose topics select one entry, whose touches select another, and whose budget leaves out a third, and a fixture with an absorbed entry that never loads.
- [ ] `flai prime --role orchestrate`, as text, `--json`, and the MCP `prime` with role `orchestrate`, includes the answer and recommendation entries ranked against the titles and first entries of the threads awaiting the operator, within the budget; a test shows an entry ranked first for a thread that asks its question again.
- [ ] A `source` on `thread_reply` may name an entry's path, and the orchestrator's prompt and the convention `strategic-agents.md` say to cite the entry when one answers the question.
- [ ] The convention `work-management.md` (template and here) and `design/system/agent-narrative.md` say an agent cites `D-nnnn` in its narrative or thread when it relies on an entry, and asks the operator only what no entry answers; `design/system/agent-context.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` describe the step.
- [ ] Measured on this repository after the seed: the pack sizes of three archived stories with and without the step, recorded in `agent-context.md`, and no pack over its budget by the step alone.

## Tasks
- T-1818 An ADR records the pack's digest step: what a story's pack and the orchestrator's select, how entries rank and load, and what part of the budget they may take
- T-1819 A story's pack and a sub-agent's add the digest step: entries selected by topics and claim, ranked within the budget, each marked as the digest's
- T-1820 The orchestrator's pack adds the answer and recommendation entries ranked against the threads awaiting the operator
- T-1821 flai prime and the MCP prime carry the digest step, a thread_reply source may cite an entry, and the orchestrator's prompt says to cite the answer the operator gave last time
- T-1822 The conventions tell an agent to cite a digest entry it relies on and to ask only what no entry answers, and agent-context.md describes and measures the step

## Notes

- Drafted by the planner for E-0025, after S-0429 so that there are entries to prime with and measure.
- Parts (ADR-0104) are unchanged: entries are items like the rest and go into parts as they fit.
- The epic names a convention `agent-narrative.md`, which does not exist. The narrative's rules live in `work-management.md` and `design/system/agent-narrative.md`, so criterion 5 names those.

### Planning

- **Touches.** The story declared none. Every touch comes from its five tasks, file by file.
  - **From layout:** a new `flai/internal/context/digest.go` with its test; `items.go`, `design.go`, `context.go`, and `story.go` for the step's constant, order, print, and use; `strategic.go` and its test for the orchestrator; the tests of `flai/cmd/prime.go` and of the MCP `prime` and `thread_reply`; `flai/internal/harness/harness.go` and its test for the orchestrator's prompt.
  - **From design:** the ADR and its index; `work-management.md` and `strategic-agents.md` in both copies, `agent-narrative.md`, `agent-context.md`, `docs/users/flai.md`, `template/CHANGELOG.md`, and `template/template.yaml` for criteria 4 to 6.
  - **Co-change, added:** `flai/internal/context/context.go` (55% of the commits that changed the pack's files), `design.go` (45%), `flai/cmd/prime_test.go` (36%), and `flai/internal/mcpserver/prime_test.go` (36%).
  - **Co-change, not added:** `flai/cmd/prime.go` (55%), since it prints what `context.go` gives; `flai/internal/mcpserver/folder.go` and `server.go` (45%), since no tool is added, and T-1821 changes `server.go` only if the source check refuses an entry's path, when `flai task done` widens the touches; `flai/internal/conventions/conventions.go` (36%), since no role is added here; `budget.go` and `parts_test.go`, since the budget and the parts are unchanged.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** Kept at flai's 48m. Measuring three packs balances the template copies, which inflate the size. Delivery is flai's play-out moved by S-0429's raise.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
