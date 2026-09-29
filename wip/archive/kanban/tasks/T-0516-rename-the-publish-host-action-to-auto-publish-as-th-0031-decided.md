---
id: T-0516
type: task
nature: remediation
title: Rename the publish host action to auto-publish, as TH-0031 decided
status: done
parent: S-0144
owner: alex
created: 2026-09-29T02:01:55Z
updated: 2026-09-29T02:06:08Z
transitions:
  - to: ready
    at: 2026-09-29T02:01:55Z
    by: agent-S-0144
  - to: in-progress
    at: 2026-09-29T02:01:55Z
    by: agent-S-0144
  - to: done
    at: 2026-09-29T02:06:08Z
    by: agent-S-0144
stream: S-0144
tags: []
touches: [flai/cmd/push.go, flai/cmd/push_test.go, flai/internal/hostapi/writes.go, design/adrs, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/conventions/git.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
---
# T-0516 Rename the publish host action to auto-publish, as TH-0031 decided

## Work

- Rename `ActionPublish` to `ActionAutoPublish` with the value `auto-publish` in `flai/internal/hostapi/writes.go`, and every use, test, and document that names `publish` as a host action; ADR-0048 (still proposed) with its file name and index row.
- Extend `TestPushPendingReleasesNothingByDefault` so that a later `flai release --pending` tags and pushes the release the push left pending (TH-0031's question).
- Accept ADR-0048, since the operator agreed in TH-0031.

## Done when

- No document or code names a `publish` host action; `make flai-test` passes but for the overlap and TH-0029 warnings; committed.

## Notes
