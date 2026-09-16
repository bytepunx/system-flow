---
id: S-006
type: story
nature: feature
title: flai import for existing monorepos
status: backlog
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T16:22:05Z
transitions: []
tags: []
---

# S-006 flai import for existing monorepos

## Goal
flai import analyses an existing repo, proposes the layout, prompts for folder names and markdown moves, and writes the manifest.

## Acceptance criteria
- [ ] --dry-run prints the proposal without changes
- [ ] Detects existing adr, docs, design, wip folders and sub-projects by build files
- [ ] Offers move, leave, or skip for each existing markdown file, uses git mv when possible
- [ ] Never overwrites an existing file without confirmation

## Tasks

## Notes
