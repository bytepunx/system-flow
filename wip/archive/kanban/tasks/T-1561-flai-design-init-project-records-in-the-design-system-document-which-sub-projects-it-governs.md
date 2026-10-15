---
id: T-1561
type: task
nature: feature
title: flai design init --project records in the design system document which sub-projects it governs
status: done
parent: S-0386
owner: alex
created: 2026-10-13T17:49:10Z
updated: 2026-10-15T04:40:13Z
transitions:
  - to: ready
    at: 2026-10-15T04:21:11Z
    by: agent-S-0386
  - to: in-progress
    at: 2026-10-15T04:21:12Z
    by: agent-S-0386
  - to: done
    at: 2026-10-15T04:40:13Z
    by: agent-S-0386
stream: S-0386
tags: [cli]
touches: [flai/cmd/design.go, flai/cmd/design_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 1141
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 7899
      cache_read: 1130862
      cache_write: 43442
      cost: 0.6383
---
# T-1561 flai design init --project records in the design system document which sub-projects it governs

## Work

In `flai/cmd/design.go`, `flai design init --project <name>`, repeatable, writes the named sub-projects to the document's front matter as `governs` and their tags to its `topics`, refusing a name the manifest does not list. Without `--project` the document governs none and has the topic `design`. Tests cover each. It waits for nothing in this story; it extends S-0381's `init`.

## Done when

- `go test ./cmd/ -run Design` passes with the new cases.

## Notes

- Layer 1, beside T-1562 and T-1563. T-1562 reads the `governs` key this task writes.
