---
id: T-3444
type: task
nature: feature
title: ADR-0312 is accepted, or superseded by an ADR recording the operator's changes
status: done
parent: S-0805
owner: alex
created: 2026-11-16T23:48:23Z
updated: 2026-11-17T07:47:18Z
transitions:
  - to: ready
    at: 2026-11-16T23:53:09Z
    by: agent-S-0805
  - to: in-progress
    at: 2026-11-16T23:53:09Z
    by: agent-S-0805
  - to: done
    at: 2026-11-17T07:47:18Z
    by: agent-S-0805
stream: S-0805
tags: [flai]
touches: [design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md, design/adrs/README.md]
usage:
  source: log
  seconds: 3764
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 231
      output: 50552
      cache_read: 15604117
      cache_write: 744137
      cost: 9.1443
---
# T-3444 ADR-0312 is accepted, or superseded by an ADR recording the operator's changes

## Work

ADR-0312 is `proposed`. Ask the operator on a thread on S-0805 whether to accept it as written, recommending acceptance, and name the two points this story builds on: §1's `suite` key defaulting to the tier's own name, and the template's labels `static`, `behavior`, `integration`, and `smoke`. On their word, accept it with `flai adr accept 0312`, which sets its status and its row in `design/adrs/README.md`. When they change it, write the new ADR with `flai adr new`, superseding ADR-0312, and set `superseded_by` on ADR-0312. This task waits for nothing; every other task waits for it, since the decision is what they build.

## Done when

- [ ] ADR-0312's status is `accepted`, or an ADR that supersedes it records the operator's changes and ADR-0312 names it in `superseded_by`.
- [ ] The ADR index row matches the ADR's status, and `flai check --strict` is clean.

## Notes
