---
id: T-2813
type: task
nature: improvement
title: Say in the design and user docs that a stopped design or docs markdown file merges by sections
status: done
parent: S-0637
owner: alex
created: 2026-11-09T22:09:28Z
updated: 2026-11-10T08:12:38Z
transitions:
  - to: ready
    at: 2026-11-10T08:06:39Z
    by: agent-S-0637
  - to: in-progress
    at: 2026-11-10T08:06:39Z
    by: agent-S-0637
  - to: done
    at: 2026-11-10T08:12:38Z
    by: agent-S-0637
stream: S-0637
tags: [docs, design]
touches: [design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-2811]
usage:
  source: log
  seconds: 359
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 23
      output: 221
      cache_read: 3226090
      cache_write: 22135
      cost: 1.1952
---
# T-2813 Say in the design and user docs that a stopped design or docs markdown file merges by sections

## Work

Update each place that says the section merge reaches only the markdown files inside `claims.shared`, linking the new ADR:

- `design/system/workflow.md` § Branches and collisions: the shared paths bullet and the S-0401 paragraph that names ADR-0175.
- `design/system/agent-coordination.md`: the ADR-0175 paragraph.
- `design/system/flai-cli.md`: the `flai stream sync` row and the merge queue's build, where they list the merge steps of a stop.
- `docs/users/flai.md` § Shared documents: which files a stopped rebase now merges by sections, and that a section both stories changed on the same lines still stops it.

Bump each file's `updated`. Say what is, not what changed: S-0637 and I-0315 belong in one sentence at most.

Waits for T-2811, whose ADR these documents link and describe. It runs beside T-2812: no file in common.

## Done when

- No document under `design/` or `docs/` says the section merge is limited to `claims.shared`.
- `flai check --strict` and the markdown lint pass on the four files.

## Notes
