---
title: Project manifest
updated: 2026-09-15
status: active
---

# Project manifest

A conforming repo has `system-flow.yaml` at its root. It is how `flai` and `flaiover` recognise a project, find its folders, and know which template it came from.

```yaml
# system-flow.yaml
version: 1                                   # manifest schema version
name: system-flow                            # project name, used in titles
key: sf                                      # short key, used in generated IDs when a project opts into prefixed IDs
description: Agentic lean project management system
template:
  repo: https://github.com/bytepunx/system-flow-template
  ref: main                                  # branch, tag, or commit
  version: 0.1.0                             # template version applied, from template.yaml
  applied: 2026-09-15T16:00:00Z
layout:                                      # folder names, defaults shown, renameable at import
  design: design
  docs: docs
  wip: wip
projects:                                    # code sub-projects at the repo root
  - name: flai
    path: flai
    kind: go
  - name: flaiover
    path: flaiover
    kind: sveltekit
dashboard:
  image: ghcr.io/bytepunx/flaiover
  tag: latest
  port: 4242
```

Rules:

- `flai` refuses to run project commands in a directory tree with no `system-flow.yaml` above the current directory, except `flai new` and `flai import`.
- `layout` is the only place folder names live. Everything else resolves through it. Subfolders such as `design/conventions` are fixed names under their layout folder.
- `template.version` lets `flai upgrade` (future) diff the applied template against a newer one.
- The manifest is human-edited YAML. `flai` rewrites only the keys it owns (`template.*`, `projects`) and preserves comments where the YAML library allows it.
