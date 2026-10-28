---
id: T-1720
type: task
nature: feature
title: "An ADR records the role file: where it lives, its front matter, its body as the prompt, the marker, and the named conditions it may reference"
status: in-progress
parent: S-0414
owner: alex
created: 2026-10-16T05:42:52Z
updated: 2026-10-28T13:07:27Z
transitions:
  - to: ready
    at: 2026-10-28T13:07:25Z
    by: agent-S-0414
  - to: in-progress
    at: 2026-10-28T13:07:27Z
    by: agent-S-0414
stream: S-0414
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
---
# T-1720 An ADR records the role file: where it lives, its front matter, its body as the prompt, the marker, and the named conditions it may reference

## Work

Write the ADR with `flai adr new`, on TH-0471's answers to questions 2, 4, and 6. It records:

- Where a role file lives: `design/workflows/roles/<name>.md`, in a `design/workflows/` folder that formulas will share as `design/workflows/formulas/` (answer 2). `design/system/repository-layout.md` adds no `design/` subfolder without an ADR, so this ADR adds `design/workflows/`.
- Each front matter field and its meaning: `name` (the activity document's and the page's), `role` (the verb `flai prime --role` and `FLAI_ROLE` take), `definition`, `action`, `trigger`, `item`, `topics`, `permissions`, `tools`, `commands`, `files`, `works_story`, and `charge`, which S-0425 uses.
- The body as the harness prompt, and the baseline marker.
- The named conditions flai implements for rules that are not plain allow-lists. Draw the list from `flai/internal/guard/guard.go`'s planner, orchestrator, and analyzer rules.

First layer: every other task of the story reads the format this decides.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`. `flai check --strict` reports nothing on it.

## Notes

The ADR's file name is not known until `flai adr new` makes it, so the task touches the shared `design/adrs` folder.

- Answer 4: a role's agent is `roles.<name>.agent` in the manifest, merged over `agent` as `planning.agent` is.
- Answer 6: flai generates the role's harness definition from its role file, so the file no longer names a hand-written one. The ADR records which fields the definition comes from (description, tools, model, and the definition's prompt) and how it reaches Claude Code. Two ways are open. The first is the `--agents` JSON that `subAgents` in `flai/internal/harness/adapters.go` already builds from `.claude/agents`, which writes nothing. The second is a symlink from `.claude/agents/<name>.md` to a file flai generates outside `.claude/`, where Claude Code follows a link. The ADR also says where the definition's prompt goes and where the start prompt (`planPrompt` and the others in `flai/internal/harness/harness.go`) goes. Generating the definitions is a later story's work (TH-0567).
