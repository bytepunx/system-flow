---
title: Repository layout
updated: 2026-10-02
status: active
topics: [all]
---

# Repository layout

A conforming monorepo has this shape. Folder names are defaults; a project may rename the three top-level documentation folders and records the chosen names in [system-flow.yaml](project-manifest.md). Tooling always resolves folders through the manifest, never by hard-coded name.

```text
<repo>/
├── CLAUDE.md               # agent operating instructions, baseline from the template
├── README.md               # human entry point
├── system-flow.yaml        # project manifest, marks a conforming repo
├── system-flow.lock.yaml   # hashes and variable values of what the template rendered, owned by flai (ADR-0015)
├── design/                 # internal documentation
│   ├── README.md
│   ├── adrs/               # point-in-time architecture decisions
│   ├── system/             # living design, always current
│   ├── tech/               # active technology choices with versions
│   ├── conventions/        # how agents work here, one file per topic, primed every session
│   ├── issues/             # recurring friction with counts and cost, see continuous-improvement.md
│   └── experiments/        # what each experiment story found, one results document per story
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
├── scripts/                # purpose-named shell scripts; the Makefile and CI call these
├── Makefile                # entry point for build, lint, check, test; targets call scripts/
├── <project-a>/            # code sub-projects, one folder each, at the root
├── <project-b>/
└── .github/                # CI and repo automation shipped by the template
```

## Folder contracts

### `design/`

Internal. Written for the people and agents building the system. Six subfolders by documentation type, each described in [design/README.md](../README.md). No other subfolders are added without an ADR.

`design/conventions/` is agent-facing: one short file per topic area stating how work is done in this repository, plus a `README.md` index in read order. The template ships the baseline; a project adds its own rules below the marker line in each file. Every agent session reads this folder first. See [conventions.md](conventions.md) and [ADR-0013](../adrs/0013-conventions-folder.md).

### `docs/`

Outward-facing. One subfolder per audience. The template ships `users`, `operators`, and `contributors`; a project may add audiences (for example `api`, `partners`). Each audience folder has an `index.md` that the dashboard uses as its landing page.

### `wip/`

Work in process. Everything in here is expected to change daily. See [work-hierarchy.md](work-hierarchy.md), [workflow.md](workflow.md), and [agent-narrative.md](agent-narrative.md).

- `kanban/` holds active items. An item is active from creation until it is archived.
- `agents/` holds one narrative file per active story.
- `archive/` mirrors `kanban/` and `agents/`. `flai archive` moves done and cancelled items here so the board stays small while history stays available for metrics.

`design/issues/` records recurring friction, defects, blockers, and inefficiencies with a count and cost per issue. See [continuous-improvement.md](continuous-improvement.md) and [ADR-0014](../adrs/0014-design-issues.md).

`design/experiments/` records what experiments found: one results document per `experiment` story, `<S-nnnn>-<slug>.md`, named as the story's file is, with front matter `title`, `updated`, `status`, and `story`, and the sections Hypothesis, Success measure, What was done, Results, and Recommendation, which says adopt, adapt, or drop. The template ships the folder's `README.md` and `template.md` to start one from. An experiment story is accepted only with its results document committed on its branch, and `flai check` validates every document in the folder. It is kept apart from `design/system/`, which says how the system is, because an experiment that is dropped does not describe it. See [workflow.md](workflow.md) and [ADR-0066](../adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md).

### `scripts/`

Purpose-named shell scripts for common tasks. The `Makefile` is the entry point and its targets call these scripts, so local runs and CI execute the same code. See [conventions/tooling.md](../conventions/tooling.md).

Scripts are POSIX `sh` under `set -eu`, because the host shell may be zsh, and a sequence of more than a few commands is a script rather than a chain typed at the prompt (I-0006). The template ships `check.sh`, `lint-md.sh`, the three test tiers, and `close-out.sh` (S-0187): run in a story's worktree before review, it runs the lint, the tests, and `flai check --strict`, checks that the narrative's `## Current state` and `## Next steps` are written, and commits with the `git commit` options it is given, stopping at the first step that fails, so a failed step can no longer carry a chain on into a commit (I-0012). A project adds its own checks to it as steps of their own. This repository's copy picks the tests by what the branch changes: `flai-test.sh` for `flai/`, `template-test.sh` for `template/`, `flaiover-test.sh` for `flaiover/`. See [conventions/work-management.md](../conventions/work-management.md).

### Code sub-projects

Sub-projects sit at the repo root, one folder each, named by the project (`flai/`, `flaiover/`). They own their build files and code. Their design lives in `design/system/<project>.md` and their technology in `design/tech/`, not inside the sub-project folder. Each sub-project has a short `README.md` pointing at those documents. Sub-projects are listed in `system-flow.yaml` so tooling can enumerate them.

### What is not in the repo

- Generated metrics, search indexes, and caches. The dashboard builds these at runtime.
- Agent transcripts. The narrative in `wip/agents` is a curated summary, not a log dump.
- Secrets. Ever.
