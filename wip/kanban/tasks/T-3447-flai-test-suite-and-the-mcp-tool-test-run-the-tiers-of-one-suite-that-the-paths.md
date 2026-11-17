---
id: T-3447
type: task
nature: feature
title: flai test --suite and the MCP tool test run the tiers of one suite that the paths select, or under --all every one but a manual tier
status: done
parent: S-0805
owner: alex
created: 2026-11-16T23:48:45Z
updated: 2026-11-17T00:39:22Z
transitions:
  - to: ready
    at: 2026-11-17T00:13:45Z
    by: agent-S-0805
  - to: in-progress
    at: 2026-11-17T00:13:45Z
    by: agent-S-0805
  - to: done
    at: 2026-11-17T00:39:22Z
    by: agent-S-0805
stream: S-0805
tags: [flai]
touches: [flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/cmd/test.go, flai/cmd/test_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-3446]
usage:
  source: log
  seconds: 1537
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 43
      output: 8236
      cache_read: 3888054
      cache_write: 151265
      cost: 2.0087
---
# T-3447 flai test --suite and the MCP tool test run the tiers of one suite that the paths select, or under --all every one but a manual tier

## Work

Carry the suite into `verify.Tier` (`flai/internal/verify/verify.go`) through `FromManifest` (`manifest.go`), with the default filled from `SuiteOrDefault()`, and update `TestCheckoutTiers`. Add `Suite` to `verify.Options` and filter in `selectTiers` (`select.go`), after the tiers are read and before the run, so that `covers` still sees every tier: with paths, the suite's tiers those paths select; with `--all`, every tier of the suite but a `manual` one, which still runs only by name. Give `flai test` a `--suite <name>` flag (`flai/cmd/test.go`), refused beside a tier name as `namedTier` refuses a name with `--all`, and refused, naming the suites there are, when no tier has that suite. Give the MCP tool `test` a `suite` argument (`flai/internal/mcpserver/test.go`) and its description. Regenerate the reference and the flag index with `make flai-reference`, which writes `docs/users/flai-reference.md` and the flag block of `docs/operators/settings.md`. Waits for T-3446, which adds the manifest key this reads.

## Done when

- [ ] `flai test --suite static <paths>` runs only the `static` tiers those paths select, and `flai test --suite integration --all` runs every `integration` tier but a manual one, each covered by a test in `select_test.go` and `flai/cmd/test_test.go`.
- [ ] An unknown suite, and `--suite` with a tier name, are refused with what to do, and tested.
- [ ] The MCP tool `test` takes `suite` and passes it on, and its test checks the description.
- [ ] `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` pass, and `flai test` on the changed paths passes.

## Notes
