---
title: Users guide
updated: 2026-09-24
status: draft
---

# system-flow for users

system-flow is a way of running a software project from inside its own repository so that people and coding agents share one source of truth. Three folders carry everything:

| Folder | Holds | Who reads it |
|--------|-------|--------------|
| `design/` | Decisions (`adrs/`), the current design (`system/`), technology choices (`tech/`), how agents work here (`conventions/`), recurring friction (`issues/`) | The people and agents building the system |
| `docs/` | Documentation per audience | Everyone else |
| `wip/` | Work items on a kanban board, and the agents' running narrative | Anyone who wants to know what is happening right now |

Two tools make it practical:

- **flai** creates a conforming repository, converts an existing one, manages work items, checks the repo, and prints metrics. See [flai.md](flai.md).
- **flaiover** is a local dashboard that renders the documentation, searches it, shows the board, and charts how work flows. See [flaiover.md](flaiover.md).

## Work items in one minute

Work is an **epic** (a deliverable spanning several stories), a **story** (one incremental deliverable), or a **task** (a piece of work for a story). Each has a nature, one of `feature`, `improvement`, `remediation`, `research`, `experiment`, and moves through `backlog`, `ready`, `in-progress`, `review`, `done`. Every move is timestamped in the item's front matter, which is where the charts come from.

A story looks like this:

```markdown
---
id: S-0004
type: story
nature: feature
title: CLI scaffold and config
status: in-progress
parent: E-0002
owner: agent
created: 2026-09-15T16:10:00Z
updated: 2026-09-16T09:02:00Z
transitions:
  - to: ready
    at: 2026-09-15T16:40:00Z
    by: alex
  - to: in-progress
    at: 2026-09-16T09:02:00Z
    by: agent
tags: [cli]
---

# S-0004 CLI scaffold and config

## Goal
...

## Acceptance criteria
- [ ] ...

## Tasks
- T-0021 Read and write ~/.flai/config.json

## Notes
```

The layout, the work items, the workflow, and the narrative are explained in the [conventions guide](conventions.md).

## Getting started

```bash
go install github.com/bytepunx/system-flow/flai@latest
flai new my-project
cd my-project
flai epic new "First deliverable"
flai dashboard
```

Or convert an existing repository:

```bash
cd existing-repo
flai import
```

Or build one by hand, to see what each piece is for: [Set up a conforming repository by hand](conventions.md#set-up-a-conforming-repository-by-hand).
