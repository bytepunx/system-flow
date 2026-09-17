---
id: T-0065
type: task
nature: feature
title: Document the sync procedure and point defaults at the published repo
status: done
parent: S-0003
owner: alex
created: 2026-09-17T00:07:13Z
updated: 2026-09-17T00:11:38Z
transitions:
  - to: ready
    at: 2026-09-17T00:11:37Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:11:38Z
    by: agent
  - to: done
    at: 2026-09-17T00:11:38Z
    by: agent
stream: S-0003
tags: [template, docs]
---

# T-0065 Document the sync procedure and point defaults at the published repo

## Work
docs/contributors: sync procedure (edit in ./template, subtree push on template changes, bump version and changelog, tag) until S-0021 automates it; design/system/template.md status; flai default source already points at the repo; this repository keeps the local path in system-flow.yaml per the project addition.

## Done when
Docs match what was done; check clean.

## Notes
