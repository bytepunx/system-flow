---
id: S-0564
type: story
nature: improvement
title: The import and analysis-cadence formulas are each run once in this repository, each instance worked to review and accepted
status: backlog
parent: E-0024
owner: alex
created: 2026-11-02T18:08:35Z
updated: 2026-11-02T18:08:35Z
transitions: []
tags: [template]
topics: [orchestration]
after: [S-0424]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0564 The import and analysis-cadence formulas are each run once in this repository, each instance worked to review and accepted

## Goal

Run each baseline formula S-0424 shipped once in this repository, so that the epic gets its live runs before it is accepted. One run imports a scratch repository from the board. The other is one analyzer cadence. Each instance is worked by a story's agent to review and accepted.

## Acceptance criteria

- [ ] Each formula is run once in this repository: an import of a scratch repository under a folder named for import, and one analyzer cadence, each instance worked to review and accepted, with its Notes naming what it produced.

## Tasks

## Notes

- Moved from S-0424 (criterion 4, T-1787 and T-1788) on TH-0810, which the operator confirmed on 2026-11-02, as D-0461 moved S-0423's live run.
- `after: S-0424`: both runs need what S-0424 adds on the main branch: the two formula files, `flai import --formula`, `flai issue list --report`, and the analyzer role step (ADR-0245). `flai formula run` reads the main checkout. The board's Import runs the installed flai, so the import run also waits until S-0424 is published.
- **The import run (from T-1787).** Make a scratch git repository under a folder the host serves for import (`flai serve import add`), outside this repository. Import it from the board's Import, which runs `flai import --formula`. The instance is a story in the imported project, worked by a story's agent on its branch and worktree, so the host needs the `agent` host action on for the new project. Watch the instance reach review, and have it accepted. Record the instance's ID and each step's exit in this story's Notes. Remove the scratch repository afterwards only with the operator's word.
- **The cadence run (from T-1788).** With the `analyze` host action on, run `flai formula run analysis-cadence --input focus=<focus> --ready --autocommit` in the main checkout, with the focus the operator chooses when finalizing this story, or none for all three. The instance's agent starts the analyzer step with `role_run`, waits for the run's `agent` event on the instance, then runs the listing step. Watch the instance reach review, and have it accepted. Record in this story's Notes the instance's ID, the report it wrote, the issues it names, and the run's cost.
