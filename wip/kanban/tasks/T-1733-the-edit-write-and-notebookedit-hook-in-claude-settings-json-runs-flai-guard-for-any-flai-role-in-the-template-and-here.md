---
id: T-1733
type: task
nature: feature
title: The Edit, Write, and NotebookEdit hook in .claude/settings.json runs flai guard for any FLAI_ROLE, in the template and here
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:23Z
updated: 2026-10-16T05:44:29Z
transitions: []
stream: S-0416
tags: [template]
touches: [".claude/settings.json", template/root/.claude/settings.json]
after: [T-1729]
---
# T-1733 The Edit, Write, and NotebookEdit hook in .claude/settings.json runs flai guard for any FLAI_ROLE, in the template and here

## Work

The second `PreToolUse` entry in `.claude/settings.json` exits at once unless `FLAI_ROLE` is `plan`, `orchestrate`, or `analyze`. Change its test to "`FLAI_ROLE` is set", and keep the `FLAI_STORY` case, in the template's file and in this repository's. The template's runs the installed `flai guard`; this repository's runs `scripts/flai.sh guard`. Each write under `.claude/` is a protected write: make it through `permission_prompt`, which asks the owner.

It waits for T-1729, whose role path decides the edits this hook hands over. It runs beside T-1730 and T-1732, since they share no path.

## Done when

A test, or a session with `FLAI_ROLE=reviewer`, shows the guard refusing a fixture role's `Write` outside the files its role file lists. `flai check --strict` passes.

## Notes
