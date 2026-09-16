---
id: T-051
type: task
nature: feature
title: conventions loader and check rules
status: done
parent: S-025
owner: alex
created: 2026-09-16T23:03:59Z
updated: 2026-09-16T23:07:54Z
transitions:
  - to: ready
    at: 2026-09-16T23:07:54Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:07:54Z
    by: agent
  - to: done
    at: 2026-09-16T23:07:54Z
    by: agent
stream: S-025
tags: [cli, check]
---

# T-051 conventions loader and check rules

## Work
internal/conventions loads <design>/conventions with front matter, order, marker count, and project-additions detection; Validate applies conventions.missing, front-matter, order, marker, length, and index rules; check calls it and skips the folder in the generic documentation walk; design layout check expects a conventions subfolder.

## Done when
Good fixture clean; bad fixture triggers every rule; monorepo clean.

## Notes
