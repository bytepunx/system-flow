---
id: T-029
type: task
nature: feature
title: Docs and design updates for check and stats
status: done
parent: S-008
owner: alex
created: 2026-09-15T18:00:58Z
updated: 2026-09-15T18:09:35Z
transitions:
  - to: ready
    at: 2026-09-15T18:09:35Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:09:35Z
    by: agent
  - to: done
    at: 2026-09-15T18:09:35Z
    by: agent
stream: S-008
tags: [docs]
---

# T-029 Docs and design updates for check and stats

## Work
flai check and flai stats commands with table and --json output; root workflow system-flow-check.yml runs flai check --strict; docs/users/flai.md, design/system/flai-cli.md, metrics.md, and documentation-standard.md updated for anything refined (README.md front matter exemption).

## Done when
flai check --strict passes on this repo and the docs match the CLI.

## Notes
