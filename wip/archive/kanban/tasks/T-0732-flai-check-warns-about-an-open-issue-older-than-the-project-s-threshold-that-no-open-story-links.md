---
id: T-0732
type: task
nature: feature
title: flai check warns about an open issue older than the project's threshold that no open story links
status: done
parent: S-0198
owner: arobson
created: 2026-10-03T01:28:25Z
updated: 2026-10-03T01:51:11Z
transitions:
  - to: ready
    at: 2026-10-03T01:28:52Z
    by: agent-S-0198
  - to: in-progress
    at: 2026-10-03T01:46:41Z
    by: agent-S-0198
  - to: done
    at: 2026-10-03T01:51:11Z
    by: agent-S-0198
stream: S-0198
tags: []
touches: [flai/internal/check, flai/internal/manifest, design/system/project-manifest.md, docs/operators/settings.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-0731]
usage:
  source: log
  seconds: 270
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 51
      output: 17052
      cache_read: 2914186
      cache_write: 59730
      cost: 1.2797
---
# T-0732 flai check warns about an open issue older than the project's threshold that no open story links

## Work

- `system-flow.yaml` gets `issues.story_after`, a Go duration (`manifest.Issues`). How an unset key behaves follows the designer's answer on TH-0074.
- `flai check` warns with `issues.no-story` about each open issue older than that, by `first_reported`, that no open story links (`issues.NoStory`). The message names the command that makes the story, `flai issue story I-nnnn`.
- Docs in the same change: `design/system/project-manifest.md`, `docs/operators/settings.md` (the settings test fails on a manifest key with no row), the `flai check` row of `design/system/flai-cli.md`, and the check part of `docs/users/flai.md`.

Waits for the second task: both edit `design/system/flai-cli.md` and `docs/users/flai.md`.

## Done when

- Check tests cover an old unlinked issue (warned), a linked one, a closed one, a young one, and the key turned off.
- `go test -race -short ./internal/check/ ./internal/manifest/` and the settings doc test in `./cmd/` pass.

## Notes
