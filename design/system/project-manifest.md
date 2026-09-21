---
title: Project manifest
updated: 2026-09-19
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
projects:                                    # releasable components at the repo root
  - name: flai
    path: flai
    kind: go
    tags: [cli]                              # story tags that mean "delivers to flai"
  - name: flaiover
    path: flaiover
    kind: sveltekit
    tags: [dashboard]
  - name: template
    path: template
    kind: template                           # released by bumping template.yaml and CHANGELOG.md, not by tag
    tags: [template, conventions]
dashboard:
  image: ghcr.io/bytepunx/flaiover
  tag: latest
  port: 4242
  autocommit: true                           # optional, default true: commit documents saved from the dashboard (ADR-0023)
  notify_url: ""                             # optional, default unset: POST new inbox entries here as JSON (S-0042)
```

Rules:

- `flai` refuses to run project commands in a directory tree with no `system-flow.yaml` above the current directory, except `flai new` and `flai import`.
- `layout` is the only place folder names live. Everything else resolves through it. Subfolders such as `design/conventions` are fixed names under their layout folder.
- `template.version` is the version `flai upgrade` compares against; `system-flow.lock.yaml` beside the manifest records the hash of every rendered file so upgrade can tell project edits from baseline (ADR-0015).
- `projects` are the components `flai release` versions, one item at a time or, since S-0087, batched by `flai release --pending`: code kinds get `<name>/vX.Y.Z` tags, kind `template` gets its version file bumped. `flai accept` versions nothing. `tags` are aliases a story or epic tag may use to say which component it delivers to.
- `dashboard.notify_url`, when set, makes the dashboard's server POST `{ project, entry: { key, kind, title, href, at } }` to that URL for each inbox entry that appears after it started: one attempt, a short timeout, a warning in the log on failure. The token and file contents are never sent. Unset by default.
- `dashboard.autocommit: false` leaves documents saved from the dashboard uncommitted; the default commits each save on the main checkout, one path per commit (ADR-0023).
- The manifest is human-edited YAML. `flai` rewrites only the keys it owns (`template.*`, `projects`) and preserves comments where the YAML library allows it.
