# system-flow

Agentic lean project management system: a documentation and work-in-process standard for monorepos, the flai CLI that applies it, and the flaiover dashboard that visualises it.

This repository follows the system-flow standard. This file is the map. The norms are in `design/conventions/`; the design is in `design/system/`; the work is in `wip/`.

## Prime your session

Before any change, read in this order:

1. `design/conventions/README.md` and every file it lists, in its order. These are the rules you work by. `flai prime --cat` prints them in one call.
2. `wip/agents/index.md`. If a stream is active, read its `## Current state` and `## Next steps`, run `git status`, reconcile, read its `## Open questions` (threads are mirrored there; `flai thread list` shows them all), and append a log entry. With the `flai` MCP server connected (`.mcp.json`), call `inbox` here and at the start of every turn, at every task transition, and before moving a story to review: it lists threads awaiting you, the ready stories in pull order, and what others changed since you last looked. Whenever you have no story of your own in progress, hold `wait_for_work` and do what it answers: it names the next story to pull as soon as one is ready and the WIP limit leaves room.
3. `wip/kanban/board.md`. If no stream is active, pull the top `ready` story in `order` within the WIP limit and open its narrative with `flai stream open`, write its tasks if it has none, then work in the worktree it prints and run `flai stream sync` at every task transition.

Precedence when rules conflict: an explicit instruction from the operator in this conversation, then the project additions below each convention's marker, then the convention baseline, then your own defaults. When the operator's instruction conflicts with a convention, follow the operator, record the conflict in the narrative's `## Decisions`, and raise it in `## Open questions` if it looks like the convention should change. When a convention seems wrong, propose the edit in `## Open questions` and keep following it until the operator or a template release changes it. Never deviate silently and never edit a baseline rule yourself.

## Layout

| Path | Purpose | You may |
|------|---------|---------|
| `design/conventions/` | How agents work here, one file per topic | Read all of it every session. Add project rules below the marker only. |
| `design/adrs/` | Architecture decisions, one per file, immutable once accepted | Add a new ADR. Never edit an accepted one except to set `superseded_by`. |
| `design/system/` | The living design, always current | Edit whenever a conversation resolves a new detail, direction, standard, or decision. |
| `design/tech/` | Every technology in use, its version, where, and why | Edit whenever a dependency is added, upgraded, or removed. |
| `design/issues/` | Recurring friction, defects, blockers, with counts and cost | Record occurrences; keep `summary.md` current. |
| `docs/<audience>/` | Outward-facing documentation per audience | Edit whenever user-facing behaviour changes. |
| `wip/kanban/` | Work items: `epics/`, `stories/`, `tasks/`, `board.md` | Create and transition items with `flai`. |
| `wip/agents/` | Your narrative, one file per story | Keep `## Current state` and `## Next steps` true. |
| `wip/archive/` | Finished items and narratives | Do not edit. `flai archive` moves things here. |
| `scripts/` | Purpose-named shell scripts behind the `Makefile` | Add a script rather than inlining a command sequence. |
| `system-flow.yaml` | Project manifest | Edit only `name`, `description`, `dashboard`. `flai` owns the rest. |

Do not create top-level folders beyond these and the sub-projects listed in `system-flow.yaml`.

## Work items

Epics span stories; stories are incremental deliverables; tasks are the pieces of a story. Every item has a nature (`feature`, `improvement`, `remediation`, `research`, `experiment`) and moves through `backlog`, `ready`, `in-progress`, `review`, `done`, or `cancelled`, with every transition timestamped. Use the CLI so front matter stays valid:

```bash
flai story new --epic E-0001 "Title"
flai task new --story S-0004 "Title"
flai move T-0021 in-progress
flai block T-0021 --reason "waiting on X"
flai check --strict
```

Schema: `design/system/work-hierarchy.md`. Policies and WIP limits: `wip/kanban/board.md` and `design/system/workflow.md`. Definitions of ready and done, narrative obligations, commits, and releases: `design/conventions/work-management.md` and `design/conventions/git.md`.

<!-- system-flow:end-of-baseline -->
<!-- Project-specific instructions go below this line. flai upgrade replaces only the section above it. -->

## About this repository

system-flow builds itself. The original brief is in `design/system/brief.md`; the epics in `wip/kanban/epics` come from it. Read `design/system/overview.md` for the four parts and the delivery order.

- `template/` is the development copy of the template repository. `system-flow.yaml` points `flai` at it with a local path. Changes to conventions land in `design/system` first, then in `template/`, in the same story.
- `flai/` is the Go CLI. Module `github.com/bytepunx/system-flow/flai`. Design in `design/system/flai-cli.md`, libraries in `design/tech/go-libraries.md`. Build and run it with `scripts/flai.sh`; test and lint with `scripts/flai-test.sh` (`make flai-test`).
- `flaiover/` is the SvelteKit dashboard. Design in `design/system/flaiover-dashboard.md`. Not started; see E-0003.
- Run flai through `scripts/flai.sh`, which builds `bin/flai` from this tree and uses `.flai-cache/config.json`. The MCP server in `.mcp.json` is the exception: it runs the installed `flai` from `PATH` (`flai mcp`), like any project made from the template, so it needs flai 1.2.1 or newer installed (`install.sh` or `flai self-upgrade`). Set `FLAI_AGENT` and `FLAI_SESSION` first. Hand-edit item bodies, not front matter that a command owns.
- Metrics definitions in `design/system/metrics.md` are the contract between `flai stats` and the dashboard. Change them only with an ADR.
