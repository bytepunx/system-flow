---
id: T-0667
type: task
nature: feature
title: "Listing paths read work items, threads, and issues with unknown front-matter fields, warn, and keep them on write; flai check reports them"
status: done
parent: S-0181
owner: arobson
created: 2026-10-01T10:43:14Z
updated: 2026-10-01T10:46:36Z
transitions:
  - to: ready
    at: 2026-10-01T10:43:33Z
    by: agent-S-0181
  - to: in-progress
    at: 2026-10-01T10:43:33Z
    by: agent-S-0181
  - to: done
    at: 2026-10-01T10:46:36Z
    by: agent-S-0181
stream: S-0181
tags: []
usage:
  source: log
  seconds: 183
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 71
      output: 20138
      cache_read: 4935332
      cache_write: 84206
      cost: 1.9446
---

# T-0667 Listing paths read work items, threads, and issues with unknown front-matter fields, warn, and keep them on write; flai check reports them

## Work

Decode work items, threads, and issues without `yaml.Strict()`, keeping each top-level front-matter key the struct does not know, verbatim, on the parsed value. The paths that list or serve them log a warning naming the field and the file; `Marshal` writes the kept keys back; `flai check` reports each as an error.

## Done when

- An item, thread, and issue with an unknown field are read by the listing paths with a warning, and written back with the field unchanged
- `flai check` reports the unknown field on each, naming the field and the file
- Tests cover both

## Notes
