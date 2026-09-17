---
id: E-0002
type: epic
nature: feature
title: flai CLI
status: done
owner: alex
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T03:59:08Z
transitions:
  - to: ready
    at: 2026-09-17T03:59:08Z
    by: alex
  - to: in-progress
    at: 2026-09-17T03:59:08Z
    by: alex
  - to: review
    at: 2026-09-17T03:59:08Z
    by: alex
  - to: done
    at: 2026-09-17T03:59:08Z
    by: alex
tags: []
---

# E-0002 flai CLI

## Outcome
A polished Go CLI that creates conforming projects, imports existing ones interactively, manages work items and narratives, validates the repo, prints metrics, and runs the dashboard.

## Stories
- S-0004 CLI scaffold and config
- S-0005 flai new from template
- S-0006 flai import for existing monorepos
- S-0007 Work item and narrative commands
- S-0008 flai check and flai stats
- S-0009 flai dashboard
- S-0010 Release pipeline for flai
- S-0029 flai release computes the bump, tags, and pushes on acceptance
- S-0031 flai structured logging and --verbose

## Notes
Defined from the brief in the root CLAUDE.md.
