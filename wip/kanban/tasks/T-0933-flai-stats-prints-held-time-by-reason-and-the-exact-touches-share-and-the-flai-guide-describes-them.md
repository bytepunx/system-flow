---
id: T-0933
type: task
nature: feature
title: flai stats prints held time by reason and the exact-touches share, and the flai guide describes them
status: backlog
parent: S-0214
owner: alex
created: 2026-10-05T05:45:15Z
updated: 2026-10-05T05:45:15Z
transitions: []
stream: S-0214
tags: [flai, docs]
touches: [flai/cmd/stats.go, flai/cmd/check_stats_test.go, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0929]
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
