---
id: S-0265
type: story
nature: remediation
title: flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out
status: backlog
owner: alex
created: 2026-10-05T00:03:14Z
updated: 2026-10-06T11:14:34Z
transitions: []
tags: []
topics: [cli]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases, docs/users/flai.md, design/issues/I-0056-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 128
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 284
          output: 60
          cache_read: 1477393
          cache_write: 131606
          cost: 0.3254
        - model: claude-opus-5-5
          input: 251
          output: 13833
          cache_read: 7240759
          cache_write: 234172
          cost: 4.2409
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-05T00:03:14Z
  value: 25
  by: planner-S-0265
  at: 2026-10-05T05:49:47Z
forecast:
  duration: 25m
  delivery: 2026-10-06T23:59:00Z
  basis: "Its own forecast of 25m; 30th in the pull order with an in-progress limit of 3, behind S-0221, S-0222, S-0224, S-0226, S-0223, S-0227, S-0229, S-0284, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0261 and S-0264."
  by: flai
  at: 2026-10-06T11:14:34Z
---
# S-0265 flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out

## Goal

This story remediates [I-0056](../../../design/issues/I-0056-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md), "flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0056 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0056 is closed with `flai issue close I-0056 --reason` saying what fixed it

## Tasks
- T-0981 flai's MD034 reports a bare email address as markdownlint does
- T-0982 The user guide says flai's lint checks bare email addresses
- T-0983 Close I-0056 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0056. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-02T16:11:21Z, 2.3 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from the issue's instance: flai's MD034 (`flai/internal/mdlint/inline.go`, `bareURL` and its caller in the inline parse) knows only `http://` and `https://` literals, while markdownlint-cli2 also reports GFM's extended email autolink, a bare `local@domain.tld`. Recognise that literal in the same parse, outside code spans and the `<…>` autolinks `emailRe` already takes, and add it to the URLs MD034 reports. The fixture that markdownlint-cli2 lints (`scripts/mdlint-fixtures.sh`) settles the edges: the characters allowed before and after the address, a trailing dot, and an address inside link text.

Touches, and where each came from:

- `flai/internal/mdlint/inline.go`: design and layout; the issue names `bareURL` there.
- `flai/internal/mdlint/mdlint_test.go`: co-change, 2 of 2 commits with `inline.go`; the regression test of I-0056 lives here.
- `flai/internal/mdlint/testdata/cases`: co-change (`expected.txt`, 2 of 2); a new fixture of bare addresses and its markdownlint findings.
- `docs/users/flai.md`: layout; its `flai check` paragraph lists what the lint checks as "bare URLs", to say email addresses too.
- `design/issues/I-0056-….md` and `design/issues/summary.md`: the second criterion, `flai issue close`.
- Left out: `doc.go`, `mdlint.go`, and `rules.go` co-changed with `inline.go` in both commits, but those added a rule; MD034 is registered already and `md034` reports whatever the parse collects.

Topics: `cli` added, as the work is in flai and its user guide.

Figures:

- Forecast 25m, adjusted from flai's 5m (134 s per unit of size times size 2, with no touches counted). The three sibling lint remediations, S-0262 (MD038), S-0240 (MD007), and S-0258 (blockquotes), took 16m, 21m, and 23m of agent time, including regenerating fixtures and the close-out. Delivery 2026-10-06T00:17Z is flai's 2026-10-05T23:57Z, 33rd in the pull order, plus the extra 20m.
- Cost of delay 25.00 USD a week, as flai works it out: 10m lost per 168h cycle at 150 USD an hour, from flai's inputs, which stand.
