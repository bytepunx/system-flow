---
title: Repository layout
updated: 2026-09-15
status: active
---

# Repository layout

A conforming monorepo has this shape. Folder names are defaults; a project may rename the three top-level documentation folders and records the chosen names in [system-flow.yaml](project-manifest.md). Tooling always resolves folders through the manifest, never by hard-coded name.

```
<repo>/
├── CLAUDE.md               # agent operating instructions, baseline from the template
├── README.md               # human entry point
├── system-flow.yaml        # project manifest, marks a conforming repo
├── design/                 # internal documentation
│   ├── README.md
│   ├── adrs/               # point-in-time architecture decisions
│   ├── system/             # living design, always current
│   └── tech/               # active technology choices with versions
├── docs/                   # outward-facing documentation, one subfolder per audience
│   ├── README.md
│   ├── users/
│   ├── operators/
│   └── contributors/
├── wip/                    # all work in process
│   ├── README.md
│   ├── agents/             # agent narrative per active work stream
│   ├── kanban/             # work items for human consumption
│   │   ├── board.md        # columns, WIP limits, policies
│   │   ├── epics/
│   │   ├── stories/
│   │   └── tasks/
│   └── archive/            # completed items and narratives, same layout as kanban/ and agents/
├── <project-a>/            # code sub-projects, one folder each, at the root
├── <project-b>/
└── .github/                # CI and repo automation shipped by the template
```

## Folder contracts

### `design/`

Internal. Written for the people and agents building the system. Three subfolders by documentation type, each described in [design/README.md](../README.md). No other subfolders are added without an ADR.

### `docs/`

Outward-facing. One subfolder per audience. The template ships `users`, `operators`, and `contributors`; a project may add audiences (for example `api`, `partners`). Each audience folder has an `index.md` that the dashboard uses as its landing page.

### `wip/`

Work in process. Everything in here is expected to change daily. See [work-hierarchy.md](work-hierarchy.md), [workflow.md](workflow.md), and [agent-narrative.md](agent-narrative.md).

- `kanban/` holds active items. An item is active from creation until it is archived.
- `agents/` holds one narrative file per active story.
- `archive/` mirrors `kanban/` and `agents/`. `flai archive` moves done and cancelled items here so the board stays small while history stays available for metrics.

### Code sub-projects

Sub-projects sit at the repo root, one folder each, named by the project (`flai/`, `flaiover/`). They own their build files and code. Their design lives in `design/system/<project>.md` and their technology in `design/tech/`, not inside the sub-project folder. Each sub-project has a short `README.md` pointing at those documents. Sub-projects are listed in `system-flow.yaml` so tooling can enumerate them.

### What is not in the repo

- Generated metrics, search indexes, and caches. The dashboard builds these at runtime.
- Agent transcripts. The narrative in `wip/agents` is a curated summary, not a log dump.
- Secrets. Ever.
