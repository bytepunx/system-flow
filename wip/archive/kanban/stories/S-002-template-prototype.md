---
id: S-002
type: story
nature: feature
title: Template prototype under ./template
status: done
parent: E-001
owner: agent
created: 2026-09-15T16:10:00Z
updated: 2026-09-15T18:05:00Z
transitions:
  - to: ready
    at: 2026-09-15T16:11:00Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:30:00Z
    by: agent
  - to: review
    at: 2026-09-15T16:45:00Z
    by: agent
  - to: done
    at: 2026-09-15T18:05:00Z
    by: alex
tags: []
---

# S-002 Template prototype under ./template

## Goal
A ./template folder that is a working prototype of the template repository: manifest, root tree with baseline CLAUDE.md, design, docs, and wip skeletons, and devex files, ready to be split into its own repo and rendered by flai.

## Acceptance criteria
- [x] template/template.yaml declares version, variables, layout, and render rules
- [x] template/root contains CLAUDE.md.tmpl, README.md.tmpl, system-flow.yaml.tmpl, and the three documentation folders with README files
- [x] Devex files present: .editorconfig, .gitignore, .gitattributes, .markdownlint.yaml, Makefile, CI workflow, PR template
- [ ] Rendering with flai new --template ./template produces a repo that passes flai check (blocked on S-005)

## Tasks
- T-005 Write template.yaml
- T-006 Write baseline CLAUDE.md.tmpl
- T-007 Write root tree and devex files

## Notes
- Accepted 2026-09-15 with the fourth criterion deferred: verified by S-005 once flai new renders the template.
