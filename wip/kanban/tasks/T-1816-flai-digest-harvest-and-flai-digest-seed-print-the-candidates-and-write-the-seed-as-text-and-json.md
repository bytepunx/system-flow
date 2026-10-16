---
id: T-1816
type: task
nature: feature
title: flai digest harvest and flai digest seed print the candidates and write the seed, as text and JSON
status: backlog
parent: S-0429
owner: alex
created: 2026-10-16T06:04:11Z
updated: 2026-10-16T06:04:11Z
transitions: []
stream: S-0429
tags: [cli]
touches: [flai/cmd/digest.go, flai/cmd/digest_test.go, docs/users/flai-reference.md]
after: [T-1815]
---
# T-1816 flai digest harvest and flai digest seed print the candidates and write the seed, as text and JSON

## Work

Add to `flai/cmd/digest.go`:

- `flai digest harvest [--since <date>] [--kind <k>]... [--json]`, through `digest.Harvest`: one line per candidate with its kind, source, date, and the first line of its text, and the entries citing it; the JSON gives every field.
- `flai digest seed [--since <date>] [--dry-run] [--json]`, through `digest.Seed`: what it wrote and skipped, with counts by kind.

Regenerate `docs/users/flai-reference.md`.

Waits for T-1815, whose functions it calls.

## Done when

`flai/cmd/digest_test.go` covers both commands in text and JSON on a fixture, and a seed refused by the check; `flai test flai/cmd` passes.

## Notes

The curator's guard (S-0431) passes `digest harvest` as a read.
