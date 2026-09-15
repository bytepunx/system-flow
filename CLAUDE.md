# system-flow

Agentic lean project management system: a documentation and work-in-process standard for monorepos, the flai CLI that applies it, and the flaiover dashboard that visualises it.

This repository follows the system-flow standard. Read this file fully before making changes. The detailed rules live in `design/system/`; this file is the short version every agent must know.

## Layout

| Path | Purpose | You may |
|------|---------|---------|
| `design/adrs/` | Architecture decisions, one per file, immutable once accepted | Add a new ADR. Never edit an accepted one except to set `superseded_by`. |
| `design/system/` | The living design, always current | Edit whenever a conversation resolves a new detail, direction, standard, or decision. |
| `design/tech/` | Every technology in use, its version, where, and why | Edit whenever a dependency is added, upgraded, or removed. |
| `docs/<audience>/` | Outward-facing documentation per audience | Edit whenever user-facing behaviour changes. |
| `wip/kanban/` | Work items: `epics/`, `stories/`, `tasks/`, `board.md` | Create and transition items, preferably with `flai`. |
| `wip/agents/` | Your narrative, one file per story | Keep it current. See below. |
| `wip/archive/` | Finished items and narratives | Do not edit. `flai archive` moves things here. |
| `system-flow.yaml` | Project manifest | Edit only `name`, `description`, `dashboard`. `flai` owns the rest. |

All documentation is markdown with YAML front matter. Diagrams are Mermaid. Dates are ISO 8601 UTC.

## Work items

Epics span stories. Stories are incremental deliverables. Tasks are the pieces of work for a story. Every item has a `nature`: `feature`, `improvement`, `remediation`, `research`, or `experiment`. States are `backlog`, `ready`, `in-progress`, `review`, `done`, `cancelled`. Every state change is appended to `transitions` with a UTC timestamp. Blocked is a timestamped interval in `blocked`, not a state.

Use the CLI so front matter stays valid:

```bash
flai story new --epic E-001 "Title"
flai task new --story S-004 "Title"
flai move T-021 in-progress
flai block T-021 --reason "waiting on X"
flai check
```

Full schema: `design/system/work-hierarchy.md`. Board policies and WIP limits: `wip/kanban/board.md` and `design/system/workflow.md`.

## How to start a session

1. Read `wip/agents/index.md`. If a stream is active, read its `## Current state` and `## Next steps`, check `git status`, append a log entry, and continue from the first next step.
2. Otherwise read `wip/kanban/board.md`, pull the top `ready` story to `in-progress` if the WIP limit allows, and open its narrative with `flai stream open S-nnn`.
3. Work tasks in order. Transition each task as you go.

## Your narrative obligations

For every story you touch, `wip/agents/<story-id>.md` must let a fresh agent resume without asking the human anything.

- Rewrite `## Current state` and `## Next steps` after every task transition and before any long-running operation.
- Append a `## Log` entry with a timestamp at every task transition, decision, and blocker. `flai stream log S-nnn "entry"` does this.
- Record decisions in `## Decisions`. Architectural ones also get an ADR.
- Put questions for the human in `## Open questions` with a date, block the item, and move on to other work if any exists.
- Never write secrets, transcripts, or raw tool output there. Summaries and paths only.

## Definition of done for a story

- All acceptance criteria checked and all tasks done or cancelled.
- Decisions captured: an ADR if architectural, an edit to `design/system` or `design/tech` if the living design changed.
- User-facing changes reflected under `docs/`.
- `flai check` passes.
- Narrative has a closing log entry. Move the story to `review`; a human moves it to `done`.

## Conventions

- Commit only when asked. Reference the story ID in commit messages.
- Do not create new top-level folders. Sub-projects are listed in `system-flow.yaml`.
- Prefer editing the living design over adding a new document.

<!-- system-flow:end-of-baseline -->
<!-- Project-specific instructions go below this line. flai upgrade replaces only the section above it. -->

## About this repository

system-flow builds itself. The original brief is in `design/system/brief.md`; the epics in `wip/kanban/epics` come from it. Read `design/system/overview.md` for the four parts and the delivery order.

- `template/` is the development copy of the template repository. `system-flow.yaml` points `flai` at it with a local path. Changes to conventions land in `design/system` first, then in `template/`, in the same story.
- `flai/` is the Go CLI. Module `github.com/bytepunx/system-flow/flai`. Design in `design/system/flai-cli.md`, libraries in `design/tech/go-libraries.md`. Run tests with `go test -race ./...` and lint with `golangci-lint run` from inside `flai/`.
- `flaiover/` is the SvelteKit dashboard. Design in `design/system/flaiover-dashboard.md`. Not started; see E-003.
- Until `flai` is functional, keep front matter valid by hand and follow the schema in `design/system/work-hierarchy.md` exactly. `flai check` will validate the history later.
- Metrics definitions in `design/system/metrics.md` are the contract between `flai stats` and the dashboard. Change them only with an ADR.
