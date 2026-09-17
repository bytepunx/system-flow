---
id: S-0002
type: story
nature: feature
title: Template prototype under ./template
status: done
parent: E-0001
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T16:35:10Z
transitions:
  - to: ready
    at: 2026-09-15T16:10:18Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:16:51Z
    by: agent
  - to: review
    at: 2026-09-15T16:22:05Z
    by: agent
  - to: done
    at: 2026-09-15T16:35:10Z
    by: alex
tags: []
---

# S-0002 Template prototype under ./template

## Goal
A ./template folder that is a working prototype of the template repository: manifest, root tree with baseline CLAUDE.md, design, docs, and wip skeletons, and devex files, ready to be split into its own repo and rendered by flai.

## Acceptance criteria
- [x] template/template.yaml declares version, variables, layout, and render rules
- [x] template/root contains CLAUDE.md.tmpl, README.md.tmpl, system-flow.yaml.tmpl, and the three documentation folders with README files
- [x] Devex files present: .editorconfig, .gitignore, .gitattributes, .markdownlint.yaml, Makefile, CI workflow, PR template
- [x] Rendering with flai new --template ./template produces a repo that passes flai check (verified by S-0005 and by flai check in S-0008)

## Tasks
- T-0005 Write template.yaml
- T-0006 Write baseline CLAUDE.md.tmpl
- T-0007 Write root tree and devex files

## Notes
- Accepted 2026-09-15 with the fourth criterion deferred: verified by S-0005 once flai new renders the template.
