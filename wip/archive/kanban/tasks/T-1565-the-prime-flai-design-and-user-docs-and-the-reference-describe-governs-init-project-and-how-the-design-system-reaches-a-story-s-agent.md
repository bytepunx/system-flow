---
id: T-1565
type: task
nature: feature
title: The prime, flai design, and user docs and the reference describe governs, init --project, and how the design system reaches a story's agent
status: done
parent: S-0386
owner: alex
created: 2026-10-13T17:49:27Z
updated: 2026-10-15T04:50:54Z
transitions:
  - to: ready
    at: 2026-10-15T04:40:44Z
    by: agent-S-0386
  - to: in-progress
    at: 2026-10-15T04:40:45Z
    by: agent-S-0386
  - to: done
    at: 2026-10-15T04:50:54Z
    by: agent-S-0386
stream: S-0386
tags: [cli]
touches: [design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/prime.go]
after: [T-1561, T-1562]
usage:
  source: log
  seconds: 609
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 100
      output: 34024
      cache_read: 4871101
      cache_write: 187121
      cost: 2.7492
---
# T-1565 The prime, flai design, and user docs and the reference describe governs, init --project, and how the design system reaches a story's agent

## Work

Describe in `design/system/agent-context.md` how the pack briefs the design system and its convention; in `design/system/flai-cli.md` and `docs/users/flai.md` `flai design init --project` and `governs`; and regenerate `docs/users/flai-reference.md`. It waits for T-1561 and T-1562, whose behaviour it describes.

## Done when

- The four documents describe it, and the markdown lint passes on them.

## Notes

- Layer 2, beside T-1564.
