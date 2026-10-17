---
id: S-0412
type: story
nature: experiment
title: A non-Claude model works one story of this repository end to end on OpenCode through OpenRouter, with the guard and the permission prompt on, and the results record what broke and what it cost
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:20:15Z
updated: 2026-10-17T17:33:42Z
transitions: []
tags: [cli]
topics: [agents, cli]
touches: [design/experiments/S-0412-a-non-claude-model-works-one-story-of-this-repository-end-to-end-on-opencode-through-openrouter-with-the-guard-and-the-permission-prompt-on-and-the-results-record-what-broke-and-what-it-cost.md, design/experiments/README.md, design/issues, design/issues/summary.md, design/system/agent-adapters.md, docs/operators/index.md]
after: [S-0411]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 30m
  delivery: 2026-10-18T03:05:00Z
  basis: "Its own forecast of 30m; 19th in the pull order with an in-progress limit of 5, behind S-0395, S-0396, S-0397, S-0241, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410 and S-0411."
  by: flai
  at: 2026-10-17T17:33:42Z
---
# S-0412 A non-Claude model works one story of this repository end to end on OpenCode through OpenRouter, with the guard and the permission prompt on, and the results record what broke and what it cost

## Goal

The epic's proof: a story of this repository, small and real, is pulled, worked, closed out, and moved to review by an agent on OpenCode running a non-Claude model through OpenRouter, with `flai guard` and the permission prompt on. The operator learns whether flai's contracts hold on a second harness and a second model family, what broke, and what the run cost, and has an issue for each defect to fix.

## Acceptance criteria

- [ ] The target story and the model are chosen with the operator on a thread on this story, recommendation first: a small open story of this repository, and an OpenRouter model that is not Claude. The target's `agent` is set to `harness: opencode`, the host's OpenRouter `provider`, and that model, by the operator.
- [ ] The target story's agent, started by `flai serve`, works it on OpenCode with the guard and the permission prompt on: it pulls it, writes or reviews its tasks, closes each with `task_done`, runs the close-out, and moves it to review, or stops, and the results say at which step and why.
- [ ] A results document under `design/experiments/`, from its template, records the model, the `opencode` and flai versions, the target story, each step it reached, every refusal, thread, and restart, what broke, and the run's tokens and cost with its `priced_by`, against what a Claude Code story of the same size cost.
- [ ] Each defect the run found is an issue in `design/issues` with its impact, linked from the results.
- [ ] `docs/operators/index.md` § Harnesses and `agent-adapters.md` § OpenCode say, from the run, what a story on OpenCode has and lacks.

## Tasks
- T-1716 The target story and the non-Claude model are chosen with the operator, and the results document records the setup
- T-1717 The target story runs on OpenCode under flai serve, and the results document records each step it reached, every refusal, thread, and restart, and its tokens and cost
- T-1718 Each defect the OpenCode run found is an issue with its impact, linked from the results
- T-1719 The operator docs and agent-adapters.md say, from the run, what a story on OpenCode has and lacks

## Notes

- Drafted by the planner for E-0023, last: it needs the guard (S-0407), the permission prompt (S-0408), the roles (S-0409), the usage (S-0410), and the status (S-0411).
- An `experiment` story is accepted like any other and records its results under `design/experiments` (ADR-0066).
- Spend: the OpenRouter key's limit is the operator's (TH-0467, question 3). The target story's own changes are accepted or sent back on their own merits, apart from this story.

### Planning

- **Touches.** None were declared. Every touch is the planner's:
  - **Layout:** the results document, `design/experiments/` with this story's file name, as S-0176's is named; `design/experiments/README.md`, which indexes it; `design/issues/summary.md`, which `flai issue new` regenerates.
  - **Design:** `design/system/agent-adapters.md` and `docs/operators/index.md` for criterion 5.
  - **Folder kept:** `design/issues`, because the issues' file names are not known until `flai issue new` makes them.
  - **Co-change:** `flai touches suggest` lists hub documents; none is reached by this story's criteria, so none was added.
- **Forecast.** 30m against flai's 17m, which matched all 234 done stories since no experiment rung has three. Following a whole story's run through `flai show`, its narrative, and its threads, and filing an issue per defect, nearly doubles it. The target's own run is the target's agent time, not this story's. The delivery, 2026-10-16T17:24Z, is flai's play-out moved by the raise times the cycle factor.
- **Cost of delay.** Not set: E-0023 has no inputs, and `flai cod` refuses. Asked on TH-0467.
