---
id: S-0808
type: story
nature: feature
title: The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:04Z
updated: 2026-11-16T23:42:04Z
transitions: []
tags: [flai]
topics: [verify, workflow]
after: [S-0807]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0808 The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §3. The operator's note of 2026-11-16: every story in review has to be rebased and verified again even when the change to main cannot break it. With the gates reading the file, a project whose acceptance needs fewer suites is held for fewer.

## Acceptance criteria

- [ ] `verify_behind_main` sends back only tiers of the suites the transition with the `accept` effect lists in `verify:`, and lists no blocker when it lists none of the tiers main's paths select.
- [ ] `queue_tiers` and the merge queue's eligibility read the same suites.
- [ ] A story whose last verify did not run a suite the acceptance needs is refused with words naming the suite and the close-out that runs it.
- [ ] On the baseline every gate judges what it judged before this story, shown by the tests of ADR-0291, ADR-0303, and ADR-0221 passing unchanged.

## Tasks

## Notes

- After S-0807: the gates read the suites the verify report records, which S-0807 adds.
