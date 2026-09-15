---
id: S-005
type: story
nature: feature
title: flai new from template
status: review
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T16:38:47Z
transitions:
  - to: ready
    at: 2026-09-15T16:36:28Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:37:47Z
    by: agent
  - to: review
    at: 2026-09-15T16:38:47Z
    by: agent
tags: []
---

# S-005 flai new from template

## Goal
flai new renders the template into a new directory with prompted or flag-supplied variables and writes system-flow.yaml.

## Acceptance criteria
- [x] Clones and caches a git template, or uses a local path
- [x] Renders .tmpl files, copies others, rewrites layout folder names
- [x] --defaults and --var produce a repo without prompts
- [x] Rendered repo passes flai check (flai check does not exist yet; the integration test checks the tree, front matter validity, absence of template syntax, and manifest parsing. Re-verify with flai check in S-008.)

## Tasks
- T-012 Template source resolution and manifest parsing
- T-013 Renderer
- T-014 flai new command
- T-015 flai template show, update, use
- T-016 Render ./template end to end and update docs

## Notes
- Also delivered `flai template show|update|use`, which the design lists but no story owned.
- Added a `quote` template function and an absolute `.template.repo` for local sources after the first end-to-end render exposed unquoted empty YAML values and a relative path in the manifest.
- Confirms the deferred S-002 criterion: ./template renders to a valid project (integration test TestRenderPrototypeTemplate).
- Timestamps in wip/ before this story were written ahead of the wall clock and were rescaled to real times in this story; the sequence is preserved, the durations are approximate.
