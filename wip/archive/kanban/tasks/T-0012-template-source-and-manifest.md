---
id: T-0012
type: task
nature: feature
title: Template source resolution and manifest parsing
status: done
parent: S-0005
owner: agent
created: 2026-09-15T16:36:28Z
updated: 2026-09-15T16:38:47Z
transitions:
  - to: ready
    at: 2026-09-15T16:36:28Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:38:47Z
    by: agent
  - to: done
    at: 2026-09-15T16:38:47Z
    by: agent
stream: S-0005
tags: [cli, template]
---

# T-0012 Template source resolution and manifest parsing

## Work
internal/template: parse template.yaml with goccy/go-yaml; resolve a source as local directory or git URL; clone git sources into cache_dir/templates/<hash> with git on PATH; Update re-clones.

## Done when
Unit tests for manifest parsing and local resolution; integration test clones a local bare repo when git is available.

## Notes
