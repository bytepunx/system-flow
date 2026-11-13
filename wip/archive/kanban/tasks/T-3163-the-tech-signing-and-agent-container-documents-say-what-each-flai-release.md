---
id: T-3163
type: task
nature: feature
title: The tech, signing, and agent-container documents say what each flai release publishes of the agent images and where
status: done
parent: S-0734
owner: alex
created: 2026-11-12T15:57:57Z
updated: 2026-11-13T08:08:13Z
transitions:
  - to: ready
    at: 2026-11-13T08:04:53Z
    by: agent-S-0734
  - to: in-progress
    at: 2026-11-13T08:04:53Z
    by: agent-S-0734
  - to: done
    at: 2026-11-13T08:08:13Z
    by: agent-S-0734
stream: S-0734
tags: [cli]
touches: [design/tech/ci.md, design/tech/docker.md, design/system/release-signing.md, docs/operators/runbooks/agent-container.md, design/tech/go.md, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/runbooks/install.md]
after: [T-3162]
usage:
  source: log
  seconds: 200
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 39
      output: 14458
      cache_read: 3347161
      cache_write: 92251
      cost: 1.6108
---
# T-3163 The tech, signing, and agent-container documents say what each flai release publishes of the agent images and where

## Work

Say what T-3162's job publishes, and bump each file's `updated`:

- `design/tech/ci.md`: the release row and the GoReleaser row.
- `design/tech/docker.md` § Agent images: the published images and their tags, the label, and that they are signed. Keep the local builds, which a development tree still makes.
- `design/system/release-signing.md` § The digest list: the agent images' list beside flaiover's.
- `docs/operators/runbooks/agent-container.md` § The images (lines 11 to 21): where the published images are. Replace the line that says nothing pushes them and no registry publishes them.

It waits for T-3162 because it describes the job that task writes.

## Done when

- None of the four files says the agent images are unpublished, and each names the registry and the list as T-3162 publishes them.
- `flai check --strict` and the markdown lint pass on them.

## Notes
