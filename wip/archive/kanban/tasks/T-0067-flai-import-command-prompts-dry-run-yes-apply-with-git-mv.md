---
id: T-0067
type: task
nature: feature
title: "flai import command: prompts, dry-run, yes, apply with git mv"
status: done
parent: S-0006
owner: alex
created: 2026-09-17T00:15:09Z
updated: 2026-09-17T00:19:50Z
transitions:
  - to: ready
    at: 2026-09-17T00:19:50Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:19:50Z
    by: agent
  - to: done
    at: 2026-09-17T00:19:50Z
    by: agent
stream: S-0006
tags: [cli, import]
---

# T-0067 flai import command: prompts, dry-run, yes, apply with git mv

## Work
cmd/import.go: flai import [dir] [--dry-run] [--yes] [--layout k=v] [--var k=v] [--template] [--force]; refuses a tree that already has system-flow.yaml unless --force; interactive prompts (huh) for layout names, folder moves, and per-file markdown destinations (design/system, design/adrs, docs/<audience>, leave, skip all); non-interactive uses defaults and leaves markdown in place; renders the template without overwriting (skips reported); git mv when the file is tracked, rename otherwise; writes the manifest with detected projects; runs check; prints a summary.

## Done when
In-process tests exercise dry-run, --yes apply, git mv, and the no-overwrite rule.

## Notes
