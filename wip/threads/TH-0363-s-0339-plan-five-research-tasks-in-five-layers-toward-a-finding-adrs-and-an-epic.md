---
id: TH-0363
title: "S-0339 plan: five research tasks in five layers, toward a finding, ADRs, and an epic"
anchor:
  path: wip/kanban/stories/S-0339-explore-adapters-for-litellm-and-openrouter.md
  item: S-0339
status: resolved
participants: [planner-S-0339, alex]
created: 2026-10-08T07:21:57Z
updated: 2026-10-08T07:22:40Z
---

# TH-0363 S-0339 plan: five research tasks in five layers, toward a finding, ADRs, and an epic

On wip/kanban/stories/S-0339-explore-adapters-for-litellm-and-openrouter.md.

## Entries

### 2026-10-08T07:21:57Z planner-S-0339
The plan for S-0339 is ready. It has five tasks in five layers, one after another, because each edits `design/system/agent-adapters.md` after the one before.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1333 Record where flai is tied to Anthropic and Claude Code today | `design/system/agent-adapters.md`, `design/system/README.md` | none |
| 2 | T-1334 Set out what LiteLLM and OpenRouter offer and which agent harnesses can drive them | `design/system/agent-adapters.md` | T-1333 |
| 3 | T-1335 Set out the vendor-neutral abstractions flai needs, with options and a recommendation | `design/system/agent-adapters.md` | T-1334 |
| 4 | T-1336 Report the options to the operator and record the decision in ADRs | `design/system/agent-adapters.md`, `design/adrs`, `design/adrs/README.md` | T-1335 |
| 5 | T-1337 Create the epic for the LiteLLM and OpenRouter adapters and link it from the finding | `design/system/agent-adapters.md` | T-1336 |

The criteria map to tasks like this:

- Criterion 1, the epic: T-1337.
- Criterion 2, the ADRs: T-1336.
- Criterion 3, the documentation: T-1333 to T-1335, closed by T-1337.

## Figures

- Touches: the four paths above. `design/adrs` is a folder touch, kept because ADR file names exist only after the decision. It is a shared path, so it holds no story.
- Topics added: `cli` and `template`.
- Forecast: 1h, delivery 2026-10-08T09:45:00Z. flai gave 39m; I raised it for the wider scope compared with S-0193, which took 35m.
- Cost of delay: 500 USD a week, from your penalty input alone.

The reasons are under `### Planning` in the story's Notes.

## Assumptions

- The finding goes in `design/system/agent-adapters.md`, as S-0193 and S-0124 put theirs in `design/system/`. It is not a new folder.
- The research reads the vendors' and harnesses' documentation. It makes no paid call to LiteLLM or OpenRouter unless you say so on T-1336's thread.
- The research treats the harness, the provider or gateway, and the model name as separate options to weigh. Claude Code pointed at a gateway is one option. Another harness, such as Codex CLI or OpenCode, or a loop of flai's own, is another. Neither is decided ahead of your answer.
- The story's agent creates the epic, but no stories under it: criterion 3 leaves those to the planner. The epic gets no cost of delay inputs, since they are yours.
- No code and no living design, such as `flai-cli.md`, changes in this story. They change when the epic is built.

I propose no split, merge, or drop.

### 2026-10-08T07:22:40Z alex
Resolved.
