---
id: S-006
type: story
nature: feature
title: flai import for existing monorepos
status: done
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T02:01:33Z
transitions:
  - to: ready
    at: 2026-09-17T00:15:09Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:15:09Z
    by: agent
  - to: review
    at: 2026-09-17T00:19:51Z
    by: agent
  - to: done
    at: 2026-09-17T02:01:33Z
    by: alex
tags: []
---

# S-006 flai import for existing monorepos

## Goal
flai import analyses an existing repo, proposes the layout, prompts for folder names and markdown moves, and writes the manifest.

## Acceptance criteria
- [x] --dry-run prints the proposal without changes
- [x] Detects existing adr, docs, design, wip folders and sub-projects by build files
- [x] Offers move, leave, or skip for each existing markdown file, uses git mv when possible
- [x] Never overwrites an existing file without confirmation

## Tasks
- T-066 importer package: scan, detect folders and sub-projects, plan
- T-067 flai import command: prompts, dry-run, yes, apply with git mv
- T-068 Import tests on a synthetic repo, docs

## Notes
- FLAI_CACHE_DIR added after a second I-009 occurrence during testing.
