---
id: T-3454
type: task
nature: feature
title: The verify report records each suite's tiers with the head and base they passed at, merges suites per commit, and reuses them across a clean rebase
status: done
parent: S-0807
owner: alex
created: 2026-11-16T23:49:27Z
updated: 2026-11-18T09:47:32Z
transitions:
  - to: ready
    at: 2026-11-18T09:09:27Z
    by: agent-S-0807
  - to: in-progress
    at: 2026-11-18T09:09:27Z
    by: agent-S-0807
  - to: done
    at: 2026-11-18T09:47:32Z
    by: agent-S-0807
stream: S-0807
tags: [flai]
touches: [flai/internal/verifyrecord/record.go, flai/internal/verifyrecord/record_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/record.go, flai/internal/verify/record_test.go]
usage:
  source: log
  seconds: 2285
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 96
      output: 39598
      cache_read: 7825697
      cache_write: 360591
      cost: 4.3248
---
# T-3454 The verify report records each suite's tiers with the head and base they passed at, merges suites per commit, and reuses them across a clean rebase

## Work

ADR-0312 §7 and criterion 3. This is the first layer: the rule and the selection both read what it records.

- Move the report's types and a reader that answers "which suites passed at this commit, with which tiers, head, and base" into a new leaf package, `flai/internal/verifyrecord`. `verify` imports `workitem`, so `Repo.Move` cannot import `verify`. `verify` keeps using the same types, so the record has one shape.
- Record per suite in the report: its tiers, and the head and base each passed at. A run of fewer suites at a commit merges into that commit's record and erases no suite another run passed there. Today `KeepAt` and `SaveReport` replace the whole record.
- Make `reusable()` reuse a passed suite's tiers across a clean rebase under ADR-0156's gates (same patch id, fingerprint, and tier order per suite). A tier that main's new paths select runs again, and so does every tier after it.
- Make `recorded()` say which suites a recorded verdict covers, so that a later answer from the record can match on suites.
- Keep reading a report written before this task: it has no suites, and counts as every tier under its own name's suite.

## Done when

- `flai/internal/verifyrecord` and `flai/internal/verify` tests show that two runs of different suites at one commit both stay in the record.
- They show that a suite passed before a clean rebase is reused while main's new paths do not select its tiers, and runs again when they do.
- They show that a report from before this task still reads.
- `flai test` passes on the changed paths.

## Notes
