---
id: T-0933
type: task
nature: feature
title: flai stats prints held time by reason and the exact-touches share, and the flai guide describes them
status: done
parent: S-0214
owner: alex
created: 2026-10-05T05:45:15Z
updated: 2026-10-07T08:41:09Z
transitions:
  - to: ready
    at: 2026-10-07T08:33:39Z
    by: agent-S-0214
  - to: in-progress
    at: 2026-10-07T08:33:40Z
    by: agent-S-0214
  - to: done
    at: 2026-10-07T08:41:09Z
    by: agent-S-0214
stream: S-0214
tags: [flai, docs]
touches: [flai/cmd/stats.go, flai/cmd/check_stats_test.go, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0929]
usage:
  source: log
  seconds: 449
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 271
      cache_read: 1853653
      cache_write: 80738
      cost: 0.8689
---
# T-0933 flai stats prints held time by reason and the exact-touches share, and the flai guide describes them

## Work

In `printClaims` (`flai/cmd/stats.go`), add held time by reason over the window and the share of stories with exact touches to the claims aggregates, and name the new `--json` sections in the command's help. Mention them in `docs/users/flai.md` where it describes what S-0205 added to `flai stats`, and bring `docs/users/flai-reference.md` in line with the help. Waits for T-0929, whose values it prints.

## Done when

- [ ] `flai stats` prints held time per reason and the exact-touches share, and `flai stats --json` carries the new values
- [ ] `check_stats_test.go` pins the printed lines
- [ ] The help, `docs/users/flai.md`, and `docs/users/flai-reference.md` agree
- [ ] `scripts/flai-test.sh` passes

## Notes
