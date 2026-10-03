---
id: T-0738
type: task
nature: improvement
title: An ADR, the design, and the docs say a full review holds the pull, and I-0007 is closed
status: done
parent: S-0243
owner: arobson
created: 2026-10-03T02:58:51Z
updated: 2026-10-03T03:26:26Z
transitions:
  - to: ready
    at: 2026-10-03T02:59:05Z
    by: agent-S-0243
  - to: in-progress
    at: 2026-10-03T03:12:44Z
    by: agent-S-0243
  - to: done
    at: 2026-10-03T03:26:26Z
    by: agent-S-0243
stream: S-0243
tags: []
touches: [design/adrs, design/system/workflow.md, design/system/flai-cli.md, design/system/agent-narrative.md, design/system/flaiover-dashboard.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, template/template.yaml, docs/users/flai.md, docs/users/flaiover.md, docs/operators/index.md, docs/operators/runbooks/migrate.md, design/issues]
after: [T-0736, T-0737]
usage:
  source: log
  seconds: 822
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 128
      output: 816
      cache_read: 5539758
      cache_write: 112492
      cost: 2.3476
---
# T-0738 An ADR, the design, and the docs say a full review holds the pull, and I-0007 is closed

## Work

Record TH-0078's decision in a new ADR beside ADR-0043, which says the in-progress limit is the only thing that holds back a ready story. Bring `workflow.md`, `flai-cli.md`, the other design files, the work-management convention (with its template copy), and the user and operator docs up to date: a full review holds the pull, and `--strict` passes over the review column over its limit. Then close I-0007 with `flai issue close I-0007 --reason`, saying what fixed it.

## Done when

- No design file or doc says that the in-progress limit alone holds the pull, or that `--strict` fails on every warning.
- I-0007 is closed with a reason that names the fix.

## Notes
