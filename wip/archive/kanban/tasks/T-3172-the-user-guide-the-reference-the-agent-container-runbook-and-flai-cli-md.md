---
id: T-3172
type: task
nature: feature
title: The user guide, the reference, the agent-container runbook, and flai-cli.md describe flai image ensure and build as the way back for a missing harness image
status: done
parent: S-0740
owner: alex
created: 2026-11-12T15:59:57Z
updated: 2026-11-13T12:45:51Z
transitions:
  - to: ready
    at: 2026-11-13T12:42:22Z
    by: agent-S-0740
  - to: in-progress
    at: 2026-11-13T12:42:22Z
    by: agent-S-0740
  - to: done
    at: 2026-11-13T12:45:51Z
    by: agent-S-0740
stream: S-0740
tags: [cli]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, design/tech/docker.md]
after: [T-3171]
usage:
  source: log
  seconds: 209
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 11501
      cache_read: 2867038
      cache_write: 85063
      cost: 1.3054
---
# T-3172 The user guide, the reference, the agent-container runbook, and flai-cli.md describe flai image ensure and build as the way back for a missing harness image

## Work

- `docs/users/flai.md` § In a container: `ensure` and `build`, the pull and the fallback, the log file, and the exits.
- `docs/users/flai-reference.md`: `make flai-reference`, never by hand.
- `docs/operators/runbooks/agent-container.md` § Build (line 23 on): `flai image ensure` first, `make agent-images` as this repository's own build until the substitution story, and the step limit and retry.
- `design/system/flai-cli.md`, the section S-0733 wrote on where the images come from: mark the pull and the host build built, and bump `updated`.

It waits for T-3171, whose commands it describes.

## Done when

- The four files describe the commands as built.
- `TestReferenceIsCurrent`, `flai check --strict`, and the markdown lint pass.

## Notes
