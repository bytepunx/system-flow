---
id: T-0013
type: task
nature: feature
title: Renderer
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

# T-0013 Renderer

## Work
Walk root/, render .tmpl files with text/template (funcs initials, slug, upper, lower; built-ins layout, today, now, template), rewrite first path segment for renamed layout folders, copy other files byte for byte, preserve exec bits, refuse to overwrite unless force, honour render.ignore.

## Done when
Tests cover each rule with a testdata template.

## Notes
