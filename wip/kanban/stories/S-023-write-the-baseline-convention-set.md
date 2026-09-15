---
id: S-023
type: story
nature: feature
title: Write the baseline convention set
status: backlog
parent: E-005
owner: alex
created: 2026-09-15T18:26:31Z
updated: 2026-09-15T18:26:31Z
transitions: []
tags: [conventions]
---

# S-023 Write the baseline convention set

## Goal
Write the baseline convention files in the template and adopt them in this repository, replacing the norms currently scattered across CLAUDE.md and the design.

## Acceptance criteria
- [ ] template/root/design/conventions/ has README.md (index in read order) and one file per agreed topic, each with front matter and a project-extension marker
- [ ] Each file is short enough to be read at session start (target under 120 lines) and states rules, not rationale; rationale links to design/system or an ADR
- [ ] This repository has the same files under design/conventions/ with its project additions below the marker
- [ ] Norms that lived only in CLAUDE.md now live in a convention file and CLAUDE.md links to them
- [ ] flai check --strict passes

## Tasks

## Notes
