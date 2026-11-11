---
id: T-2979
type: task
nature: improvement
title: The design and docs say a defined role runs on the container harness, and I-0258 is closed
status: done
parent: S-0705
owner: alex
created: 2026-11-10T10:29:20Z
updated: 2026-11-11T16:04:21Z
transitions:
  - to: ready
    at: 2026-11-11T15:59:00Z
    by: agent-S-0705
  - to: in-progress
    at: 2026-11-11T15:59:00Z
    by: agent-S-0705
  - to: done
    at: 2026-11-11T16:04:21Z
    by: agent-S-0705
stream: S-0705
tags: [docs, issues]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/system/agent-adapters.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0258-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md, design/issues/summary.md]
after: [T-2975, T-2976]
usage:
  source: log
  seconds: 321
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 16981
      cache_read: 2867628
      cache_write: 143925
      cost: 1.7519
---
# T-2979 The design and docs say a defined role runs on the container harness, and I-0258 is closed

## Work

Third layer. It waits for T-2975 and T-2976, since it describes what they built: the field pilot reads, and the container's and cache folder's names.

- Replace each sentence that says a role a project defines is refused on the container harness, with what it does now: its name and inputs travel in pilot's start request, pilot reads its file from the mounted checkout, and its container and cache folder are named for the role.
  - `design/system/strategic-agents.md` § In a container (S-0509): the paragraph after the mounts table and the I-0258 line under The live run (S-0567). Add the role's name to the cache folder row.
  - `design/system/flai-cli.md` § the container harness's adapter paragraph (`containerSubject`, `roleStart`).
  - `design/system/agent-adapters.md`: drop the role a project defines from what is still not done.
  - `docs/users/flai.md` § In a container, the paragraph on the strategic agents' container.
  - `docs/operators/index.md` § What it cannot do yet, in the container section.
- Bump `updated` in each design and docs file changed.
- Close the issue with `flai issue close I-0258 --reason "<what fixed it>"`, naming S-0705, the start request's role name and inputs, and pilot's read of the file. That command writes the issue file and `design/issues/summary.md`.

## Done when

- No document under `design/system` or `docs` says a defined role is refused on the container harness.
- I-0258 is closed with its reason.
- `flai test` on the changed paths passes, the markdown lint among its tiers.

## Notes
