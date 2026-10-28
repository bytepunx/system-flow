---
title: Repository layout
updated: 2026-10-28
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
│   ├── experiments/        # what each experiment story found, one results document per story
│   ├── analysis/           # what the analyzer found, one report per run
│   ├── digest/             # what past stories learned, one fact per file, written through flai digest
│   ├── components/         # supporting files for components, one examples file per component
│   └── workflows/          # what flai runs as data
│       └── roles/          # one file per strategic role, its rules and harness definition (ADR-0207)
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
│   ├── threads/            # conversations between the designer and agents, anchored to a document or item
│   ├── messages/           # conversations between the agents of two open stories
│   └── archive/            # completed items and narratives, same layout as kanban/ and agents/
├── scripts/                # purpose-named shell scripts; the Makefile and CI call these
├── Makefile                # entry point for build, lint, check, test; targets call scripts/
├── <project-a>/            # code sub-projects, one folder each, at the root
├── <project-b>/
└── .github/                # CI and repo automation shipped by the template
```

## Folder contracts

### `design/`

Internal. Written for the people and agents building the system. Its subfolders by documentation type, each described in [design/README.md](../README.md). No other subfolders are added without an ADR.

`design/conventions/` is agent-facing: one short file per topic area stating how work is done in this repository, plus a `README.md` index in read order. The template ships the baseline; a project adds its own rules below the marker line in each file. Every agent session reads this folder first. See [conventions.md](conventions.md) and [ADR-0013](../adrs/0013-conventions-folder.md).

### `docs/`

Outward-facing. One subfolder per audience. The template ships `users`, `operators`, and `contributors`; a project may add audiences (for example `api`, `partners`). Each audience folder has an `index.md` that the dashboard uses as its landing page.

### `wip/`

Work in process. Everything in here is expected to change daily. See [work-hierarchy.md](work-hierarchy.md), [workflow.md](workflow.md), and [agent-narrative.md](agent-narrative.md).

- `kanban/` holds active items. An item is active from creation until it is archived.
- `agents/` holds one narrative file per active story.
- `threads/` holds the conversations between the designer and agents, one file per thread, `TH-nnnn-<slug>.md`, each anchored to a document, a heading, or a work item ([ADR-0020](../adrs/0020-files-plus-mcp.md)).
- `messages/` holds the conversations between the agents of two stories in progress or in review, one file per conversation, `MS-nnnn-<slug>.md`, kept apart from threads and closed when either story leaves ([ADR-0120](../adrs/0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md)).
- `archive/` mirrors `kanban/` and `agents/`. `flai archive` moves done and cancelled items here so the board stays small while history stays available for metrics.

`design/issues/` records recurring friction, defects, blockers, and inefficiencies with a count and cost per issue. See [continuous-improvement.md](continuous-improvement.md) and [ADR-0014](../adrs/0014-design-issues.md).

`design/experiments/` records what experiments found: one results document per `experiment` story, `<S-nnnn>-<slug>.md`, named as the story's file is, with front matter `title`, `updated`, `status`, and `story`, and the sections Hypothesis, Success measure, What was done, Results, and Recommendation, which says adopt, adapt, or drop. The template ships the folder's `README.md` and `template.md` to start one from. An experiment story is accepted only with its results document committed on its branch, and `flai check` validates every document in the folder. It is kept apart from `design/system/`, which says how the system is, because an experiment that is dropped does not describe it. See [workflow.md](workflow.md) and [ADR-0066](../adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md).

`design/analysis/` holds the analyzer's reports (S-0223): one per run, `<date>-<focus>.md`, the date the day it ran in UTC as `YYYY-MM-DD` and the focus `bottlenecks`, `intent`, or `risk`, or `all` for a run asked for none. Its front matter has `title`, `updated`, `status` (`draft` while the analyzer writes it, `active` once the run has ended, and `deprecated` when a later report replaces it), `focus`, and the window its metrics cover, `from` and `to`, as dates; it has one section per finding, with its evidence, its severity, and its estimated impact. The template ships the folder's `README.md`, which says what a report holds and indexes the reports, one row each. Only the analyzer writes here, and `flai check` validates every report and warns of one the index does not list and of a row with no report. It is kept apart from `design/system/`, which says how the system is, because a report says what one run found at one time. See [strategic-agents.md](strategic-agents.md).

`design/digest/` holds the digest (S-0428, [ADR-0191](../adrs/0191-the-digest-is-one-markdown-file-per-entry-under-the-design-folder-s-digest.md)): what past stories learned, one fact per file, `D-nnnn-<slug>.md`, so that a story's agent is primed with it rather than asking again. An entry's `kind` is `decision`, `answer`, `recommendation`, or `cause`; its front matter names its `sources` (work items, threads, and issues), the `topics` and `paths` it bears on, its `date`, and its `status`: `active`, `absorbed` (with `absorbed_by`, the document that now records it), or `compacted`. flai writes the front matter: `flai digest new` numbers an entry one past the highest D-nnnn in the checkout, on main, in every worktree, and on every story branch, and `flai digest edit` changes its body and lists. Only flai and the digest's keeper, the curator, write here, through those commands. The template ships the folder's `README.md`, which says what an entry holds; the folder keeps no index, since `flai digest list` prints one and parallel writers would conflict on a table. `flai check` reports `digest.entry` on an entry whose name or front matter is wrong, and `digest.secret` on text with the shape of a key or token: an entry says where a secret lives, never its value. Nothing here is deleted; absorption and compaction change an entry's status, and git keeps every earlier text.

`design/components/` holds supporting files for a project's components (S-0443): today one examples file per component, `design/components/<sub-project>/<path>`, the sub-project's manifest `name` and the component's path in it with `.svelte` replaced by `.examples.json`. Each is JSON, an exception to "Everything is markdown" as the design token file is, holding the component's examples, named sets of its props that flaiover's Components page lists and loads. The page writes them through `file.save`, on main or on the story it shows, and a story that changes a component keeps its examples in step; `flai check` reports `components.examples` for a file that does not fit the format or its component. The template ships the folder's `README.md`. See [component-library.md § Examples](component-library.md#examples-s-0443) and [ADR-0164](../adrs/0164-a-component-s-examples-are-named-sets-of-props-in-one-json-file-per-component.md).

`design/workflows/` holds what flai runs as data rather than as its own code (S-0414, [ADR-0207](../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md)): `roles/` now, `formulas/` when E-0024's formula stories land. `roles/` holds one file per strategic role, `<name>.md`, and a `README.md` that lists them. A role file's front matter names the role, the host action that gates it, its triggers, the items a run takes, its topics, the calls the guard allows it, and the fields flai generates its harness definition from; its body is that definition's prompt, with the baseline marker. The template ships `planner.md`, `orchestrator.md`, `analyzer.md`, and `curator.md`; a file it does not ship is the project's own. A role file is code the operator runs: it changes only through a story the operator accepts, and no role may edit the folder. `flai role list` prints the roles, and `flai check` reports `roles.file`, `roles.grant`, and `roles.definition`. See [strategic-agents.md § Role files](strategic-agents.md#role-files-s-0414).

### `scripts/`

Purpose-named shell scripts for common tasks. The `Makefile` is the entry point and its targets call these scripts, so local runs and CI execute the same code. See [conventions/tooling.md](../conventions/tooling.md).

Scripts are POSIX `sh` under `set -eu`, because the host shell may be zsh, and a sequence of more than a few commands is a script rather than a chain typed at the prompt (I-0006). The template ships `check.sh`, `lint-md.sh`, the three test tiers, and `close-out.sh` (S-0187): run in a story's worktree before review, it runs the lint, the tests, and `flai check --strict`, checks that the narrative's `## Current state` and `## Next steps` are written, and commits with the `git commit` options it is given, stopping at the first step that fails, so a failed step can no longer carry a chain on into a commit (I-0012). A project adds its own checks to it as steps of their own. This repository's copy picks the tests by what the branch changes: `flai-test.sh` for `flai/`, `template-test.sh` for `template/`, `flaiover-test.sh` for `flaiover/`. See [conventions/work-management.md](../conventions/work-management.md).

### Code sub-projects

Sub-projects sit at the repo root, one folder each, named by the project (`flai/`, `flaiover/`). They own their build files and code. Their design lives in `design/system/<project>.md` and their technology in `design/tech/`, not inside the sub-project folder. Each sub-project has a short `README.md` pointing at those documents. Sub-projects are listed in `system-flow.yaml` so tooling can enumerate them.

### What is not in the repo

- Generated metrics, search indexes, and caches. The dashboard builds these at runtime.
- Agent transcripts. The narrative in `wip/agents` is a curated summary, not a log dump.
- Secrets. Ever.
