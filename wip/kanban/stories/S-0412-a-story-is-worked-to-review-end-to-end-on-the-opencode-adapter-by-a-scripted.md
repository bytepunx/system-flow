---
id: S-0412
type: story
nature: feature
title: A story is worked to review end to end on the opencode adapter by a scripted session in OpenRouter's format, with the guard and the permission prompt on, on fixtures alone
status: cancelled
parent: E-0023
owner: alex
created: 2026-10-16T05:20:15Z
updated: 2026-11-07T21:53:43Z
transitions:
  - to: ready
    at: 2026-11-07T20:43:27Z
    by: alex
  - to: cancelled
    at: 2026-11-07T21:53:43Z
    by: alex
tags: [cli]
topics: [agents, cli]
touches: [flai/tests/integration/opencode_test.go, flai/tests/integration/testdata/opencode/session.jsonl, flai/tests/integration/testdata/opencode/generation.json, design/system/agent-adapters.md, docs/operators/index.md]
after: [S-0407, S-0408, S-0410, S-0411]
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
    - kind: orchestrator
      seconds: 270
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 456
          cache_read: 3705983
          cache_write: 192722
          cost: 0.9755
cost_of_delay:
  value: 74.63
  by: planner-E-0023
  at: 2026-11-07T09:36:57Z
forecast:
  duration: 50m
  delivery: 2026-11-08T04:28:00Z
  basis: "Its own forecast of 50m; 3rd in the pull order with an in-progress limit of 5, behind S-0408, S-0409, S-0410, S-0591 and S-0411."
  by: flai
  at: 2026-11-07T21:49:38Z
finalized:
  by: orchestrator
  at: 2026-11-07T09:37:35Z
---
# S-0412 A story is worked to review end to end on the opencode adapter by a scripted session in OpenRouter's format, with the guard and the permission prompt on, on fixtures alone

## Goal

The epic's proof, on fixtures alone. The operator runs no live OpenCode and no OpenRouter key in this repository for now (TH-0533, TH-0532, D-0506), so the proof that flai's contracts compose on a second harness is a test: `flai serve` starts a story's agent on the `opencode` adapter in a scratch repository, a stub `opencode` replays a scripted session of a non-Claude model in OpenRouter's response shape, and the story is pulled, has a task closed, and reaches review, with `flai guard` and the permission prompt on. The operator learns which of the epic's pieces hold together, and `agent-adapters.md` says what only a live run would prove.

## Acceptance criteria

- [ ] A scripted session, `flai/tests/integration/testdata/opencode/session.jsonl`, follows the pinned OpenCode release's event types (S-0406) and names a non-Claude OpenRouter model as `openrouter/<vendor>/<model>`; `generation.json` holds OpenRouter's documented `/api/v1/generation` answer for its calls. Neither holds a key or anything recorded from a paid call.
- [ ] An integration test starts the story's agent through `flai serve` on a stub `opencode` that replays the session: it starts the story, closes one task with `task_done`, and moves the story to review, each through `flai mcp --agent <name>`. It runs each tool call through `flai guard --harness opencode` with the input the plugin hands it (S-0407).
- [ ] In the same run, a sub-agent's `flai` write is refused by the guard and the story's own agent's is let through, and a protected write is held and allowed by an owner's answer on the story's permission thread, or refused by the deny rule, as S-0408's ADR settled.
- [ ] The story's `usage` after the run carries the session's tokens per model and is `priced_by: gateway`, its cost read from `generation.json` through a fake OpenRouter on `httptest`, as `flai/internal/usage/openrouter_test.go` fakes it (S-0410); its card shows the session's last step from the stream (S-0411).
- [ ] The test needs no `opencode` binary, no network, and no key, and runs in the integration tier, so `make integration` and a close-out that selects it run it on any host.
- [ ] `design/system/agent-adapters.md` § OpenCode and `docs/operators/index.md` § Harnesses say what the fixture run proves, and what stays unproven until a live run against a real `opencode` and OpenRouter, which is not part of this repository for now.

## Tasks
- T-1716 The scripted OpenCode session and its OpenRouter generation fixture are written from the pinned release's event types and OpenRouter's documented answer
- T-1717 An integration test drives a story of a scratch repository to review on a stub opencode that replays the session, through flai mcp and flai guard
- T-1718 The integration test asserts the guard's refusal, the permission prompt's answer, the gateway-priced usage, and the card's status from the run
- T-1719 agent-adapters.md and the operator docs say what the fixture run proves and what only a live run would

## Notes

- Drafted by the planner for E-0023, last. Re-planned on 2026-11-07 after the operator's answer on TH-0533: "We will only be using OpenRouter fixtures, no direct integration will be part of this repository for now". It was an `experiment` that ran a real non-Claude model through OpenRouter on a funded key; it is now a `feature`, a test on fixtures, and spends nothing.
- A defect the test finds in an earlier E-0023 story's code is fixed here when it is small, or recorded with `flai issue new` from the story's worktree and named in the Notes.
- A live run, a real `opencode` on the host and an OpenRouter key with a limit, is left for a later story outside this epic, when the operator wants direct integration.

### Planning

- **Touches.** Re-planned on 2026-11-07. Every touch is the planner's:
  - **Layout:** `flai/tests/integration/opencode_test.go`, new, in the integration tier beside `mcp_http_test.go`, since the run needs real git and worktrees; the stub `opencode` is the test binary re-run as a helper process, so it needs no file of its own. The two fixtures under `flai/tests/integration/testdata/opencode/`, new, named by what they hold.
  - **Design:** `design/system/agent-adapters.md` and `docs/operators/index.md` for criterion 6.
  - **Dropped with the live run:** the results document under `design/experiments/` and its `README.md` entry, `design/issues`, and `design/issues/summary.md`. An issue the test finds is filed by `flai issue new`, which the close-out commits.
  - **Co-change:** `flai touches suggest` lists hub documents (`docs/users/flai-reference.md`, `design/adrs/README.md`, `template/CHANGELOG.md`, `strategic-agents.md`, `project-manifest.md`) from the old experiment's seed commits. This story adds no command, ADR, or template file, so none was added.
- **After.** It reads four stories' output on main: S-0407's `flai guard --harness opencode` and the plugin's input shape, which the stub hands each call to; S-0408's ADR and hold-and-ask, which decide whether the protected write is held or denied; S-0410's OpenCode usage reader and its pricing from OpenRouter's spend log; and S-0411's OpenCode stream, which the card's status comes from. S-0406's adapter reaches it through all four.
- **Forecast.** 50m against flai's 22m, which rests on only 5 done medium feature stories. A stub `opencode` that speaks to `flai mcp` over standard input and hands each call to `flai guard`, with a fake OpenRouter and assertions over four stories' output in one run, more than doubles it. The delivery, 2026-11-07T18:19Z, is flai's 17:51Z play-out moved by the 28m raise.
- **Cost of delay.** 74.63 USD a week. E-0023 carries `penalty_per_week: 500`, the operator's answer to TH-0467 question 1, confirmed on TH-0470 on 2026-10-20. `flai cod` shares it among the epic's seven open stories by forecast, 50m of 5h35m. Kept as flai gives it.
- 2026-11-07T21:53:43Z: moved to cancelled: E-0023 cancelled: OpenCode isn't going to receive first-class support given its aversion to supporting before hooks for tool calls (amongst other things)
