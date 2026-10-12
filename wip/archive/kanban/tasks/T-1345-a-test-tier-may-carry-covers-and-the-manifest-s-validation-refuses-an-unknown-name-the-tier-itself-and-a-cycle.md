---
id: T-1345
type: task
nature: improvement
title: A test tier may carry covers, and the manifest's validation refuses an unknown name, the tier itself, and a cycle
status: done
parent: S-0342
owner: alex
created: 2026-10-08T08:05:23Z
updated: 2026-10-08T21:02:33Z
transitions:
  - to: ready
    at: 2026-10-08T20:58:08Z
    by: agent-S-0342
  - to: in-progress
    at: 2026-10-08T20:58:08Z
    by: agent-S-0342
  - to: done
    at: 2026-10-08T21:02:33Z
    by: agent-S-0342
stream: S-0342
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md]
usage:
  source: log
  seconds: 265
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 7624
      cache_read: 1249001
      cache_write: 50075
      cost: 0.7144
---
# T-1345 A test tier may carry covers, and the manifest's validation refuses an unknown name, the tier itself, and a cycle

## Work

Add `Covers []string` to `manifest.TestTier`, with the tags `yaml:"covers,omitempty"` and `json:"covers,omitempty"`, so that both the YAML and `flai manifest set tests='<JSON>'` take it. Its doc comment says that a tier that covers another runs everything the other runs, so that when both are selected only the covering tier runs.

Extend `Manifest.TestErrors` so that each problem is a `manifest.tests` error in the existing form (index, field, tier name):

- a name in `covers` that names no tier;
- a tier that names itself;
- a name given twice in one tier's `covers`;
- a cycle through `covers`, reported once, naming the tiers on it in list order.

`TestTiers` refuses such a list as it refuses any other error, so `flai test` and `flai verify` run none of it. Add cases to `TestTestErrors` in `manifest_test.go` for each refusal and for a valid chain (`a` covers `b`, `b` covers `c`).

`flai manifest set tests=` writes the list through `settings.go`: make `testsLines` render `covers` and `sameTiers` compare it, with a case in `settings_test.go` that sets a tier with `covers` and reads it back.

`cmd/settings_doc_test.go` fails when a manifest key has no row in `docs/operators/settings.md`, so add the row `tests[].covers` beside `tests[].all_command`: default none; the tiers this one runs everything of, which are skipped when it is selected with them.

It waits for nothing: it is the key every other task reads.

## Done when

- `manifest.TestTier` has `Covers`, read from YAML and JSON and written back by `flai manifest set tests=`.
- `TestErrors` refuses an unknown name, the tier itself, a duplicate, and a cycle, each with its index, field, and name; a valid chain gives no error.
- `docs/operators/settings.md` has the `tests[].covers` row.
- `flai test ./internal/manifest/ ./cmd/` passes.

## Notes

Layer 1 of S-0342's plan.
