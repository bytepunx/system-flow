---
id: T-0025
type: task
nature: feature
title: "Check: narratives, board, and documentation front matter"
status: done
parent: S-0008
owner: alex
created: 2026-09-15T18:00:58Z
updated: 2026-09-15T18:09:33Z
transitions:
  - to: ready
    at: 2026-09-15T18:09:33Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:09:33Z
    by: agent
  - to: done
    at: 2026-09-15T18:09:33Z
    by: agent
stream: S-0008
tags: [cli, check]
---

# T-0025 Check: narratives, board, and documentation front matter

## Work
Report active stories without a narrative, narratives whose stream does not match the filename or whose story is gone, missing narrative sections, index.md rows out of date, board order entries that are not ready stories, WIP limit breaches, and design/docs markdown lacking title and updated (README.md exempt) plus ADR front matter and numbering.

## Done when
Each rule has a name, a level, and a test case.

## Notes
