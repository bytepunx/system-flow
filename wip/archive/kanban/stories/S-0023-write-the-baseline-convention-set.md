---
id: S-0023
type: story
nature: feature
title: Write the baseline convention set
status: done
parent: E-0005
owner: alex
created: 2026-09-15T18:26:31Z
updated: 2026-09-15T23:00:55Z
transitions:
  - to: ready
    at: 2026-09-15T18:40:51Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:40:51Z
    by: agent
  - to: review
    at: 2026-09-15T18:43:49Z
    by: agent
  - to: done
    at: 2026-09-15T23:00:55Z
    by: alex
tags: [conventions]
---

# S-0023 Write the baseline convention set

## Goal
Write the baseline convention files in the template. The operator reviews and edits them there; adopting them in this repository is S-0026, pulled when the operator asks.

## Acceptance criteria
- [x] template/root/design/conventions/ has README.md (index in read order) and one file per agreed topic, each with front matter and a project-extension marker
- [x] Each file is short enough to be read at session start (target under 120 lines) and states rules, not rationale; rationale links to design/system or an ADR
- [x] Rendering the template with flai new produces the folder and flai check --strict passes on the result
- [x] flai check --strict passes

## Tasks
- T-0036 README index, session-start, communication
- T-0037 work-management, decisions, documentation
- T-0038 code-quality, git, safety, tooling
- T-0039 Render check and design note

## Notes
- Longest file is under 60 lines; the whole set is 341 lines.
- The operator edits the template files directly; S-0026 copies them here afterwards.
