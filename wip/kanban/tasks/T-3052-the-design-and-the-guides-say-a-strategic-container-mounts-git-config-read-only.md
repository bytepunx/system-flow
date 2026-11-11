---
id: T-3052
type: task
nature: improvement
title: The design and the guides say a strategic container mounts .git/config read-only
status: done
parent: S-0688
owner: alex
created: 2026-11-10T10:35:05Z
updated: 2026-11-10T19:43:58Z
transitions:
  - to: ready
    at: 2026-11-10T19:32:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T19:32:51Z
    by: system-flow
  - to: done
    at: 2026-11-10T19:43:58Z
    by: agent-S-0688
stream: S-0688
tags: [docs, container]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, design/system/agent-adapters.md]
after: [T-3045]
usage:
  source: log
  seconds: 667
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 6471
      cache_read: 2566542
      cache_write: 75574
      cost: 1.2151
---
# T-3052 The design and the guides say a strategic container mounts .git/config read-only

## Work

Add the read-only `.git/config` row to each strategic container mount table and link the ADR T-3045 wrote:

- `design/system/strategic-agents.md` § In a container (S-0509): the mount table, the paragraph on what the guard does not see, and the line naming I-0257, which now says it is closed and how.
- `design/system/flai-cli.md`: the `roleStart` section's mount table.
- `docs/operators/index.md`: the strategic container section's mount table and its `.git/config` note. Its last paragraph still says the orchestrator's shell `flai verify` and `flai accept` run in its container; correct it to what ADR-0278 §5 says, since it is the same passage.
- `docs/users/flai.md`: the strategic container mount table.

Bump each document's `updated`. It waits for T-3045, whose ADR these documents link and describe; it runs beside the code task, since they share no file.

## Done when

- The four documents name the read-only `.git/config` mount and link the new ADR, and none says `.git/config` is writable in a strategic container.
- The markdown lint and `flai check --strict` pass.

## Notes
