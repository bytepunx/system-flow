---
id: T-0987
type: task
nature: improvement
title: The context pack is fitted so that the prime result stays under the tool result limit, with a test that reproduces I-0068
status: done
parent: S-0261
owner: alex
created: 2026-10-05T05:52:15Z
updated: 2026-10-06T23:13:31Z
transitions:
  - to: ready
    at: 2026-10-06T23:07:03Z
    by: agent-S-0261
  - to: in-progress
    at: 2026-10-06T23:07:03Z
    by: agent-S-0261
  - to: done
    at: 2026-10-06T23:13:31Z
    by: agent-S-0261
stream: S-0261
tags: [flai]
touches: [flai/internal/context]
after: [T-0986]
usage:
  source: log
  seconds: 388
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 41
      output: 17628
      cache_read: 1464627
      cache_write: 60940
      cost: 1.0201
---
# T-0987 The context pack is fitted so that the prime result stays under the tool result limit, with a test that reproduces I-0068

## Work

Build T-0986's ADR in `flai/internal/context`. Depending on the fix it chose, that means one of these:

- measuring in `Pack.size` the result the MCP tool encodes;
- splitting the pack into parts;
- briefing fewer ADRs in `AddDesign` once the briefs exceed the budget, the way `AddBriefs` counts what it leaves out;
- changing `DefaultBudget` in `budget.go`.

Add a test beside `TestTheBriefsAreKeptWhenTheyTakeThePackOverTheBudget` in `brief_test.go` that reproduces the cause. It builds a pack whose briefs, or encoding, took the result over the limit, as S-0205's, S-0210's, and S-0211's did, and shows that the result now fits. If the ADR changes the rule that test checks, change the test with it.

It waits for T-0986 because the ADR decides what to build.

## Done when

- A test under `flai/internal/context` fails on the code as it was and passes with the fix.
- The pack the package returns for a story whose briefs exceed the budget fits the limit the ADR sets.
- `scripts/flai-test.sh` passes.

## Notes
