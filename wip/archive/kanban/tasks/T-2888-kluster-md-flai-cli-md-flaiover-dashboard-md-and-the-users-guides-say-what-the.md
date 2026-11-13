---
id: T-2888
type: task
nature: feature
title: kluster.md, flai-cli.md, flaiover-dashboard.md, and the users' guides say what the move, the close-out, the review page, and the pack say of the namespace
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:28:04Z
updated: 2026-11-13T15:04:55Z
transitions:
  - to: ready
    at: 2026-11-13T14:52:11Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T14:52:11Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T15:04:55Z
    by: agent-S-0662
stream: S-0662
tags: [cli, dashboard]
touches: [design/system/kluster.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/verify.go, flai/cmd/accept.go, docs/users/flai-reference.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/review_test.go]
after: [T-2878, T-2880, T-2881, T-2883, T-2885, T-2886, T-3306]
usage:
  source: log
  seconds: 764
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 247
      output: 95479
      cache_read: 18345254
      cache_write: 785605
      cost: 9.8686
---
# T-2888 kluster.md, flai-cli.md, flaiover-dashboard.md, and the users' guides say what the move, the close-out, the review page, and the pack say of the namespace

## Work

Say what T-2880, T-2881, T-2885, T-2886, T-2883, and T-2878 built:

- `design/system/kluster.md`: replace "The close-out and `flai verify` do not use the cluster" with what the close-out records, and say what the move to review warns and what the pack names.
- `design/system/flai-cli.md`: the move's warnings, `flai verify --last`, `item_get`'s field, and the evidence form.
- `design/system/flaiover-dashboard.md`: the review page's new part and its links.
- `docs/users/flai.md` and `docs/users/flaiover.md`: all of it, for a reader with a kluster project.

It waits for every code task, whose behaviour it describes.

## Done when

- Each document says what the code does, with no part left as not built that is built.
- The markdown lint and `flai check --strict` pass on the five files.

## Notes

- Criterion 8 of S-0662. If a flag was added, add `docs/users/flai-reference.md` with `flai touches` and regenerate it with `make flai-reference`.
