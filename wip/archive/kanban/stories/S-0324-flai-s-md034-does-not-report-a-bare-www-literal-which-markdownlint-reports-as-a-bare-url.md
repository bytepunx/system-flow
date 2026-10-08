---
id: S-0324
type: story
nature: remediation
title: "flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL"
status: done
owner: alex
created: 2026-10-07T18:59:55Z
updated: 2026-10-08T04:34:30Z
transitions:
  - to: ready
    at: 2026-10-08T00:29:56Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T00:30:39Z
    by: agent-S-0324
  - to: review
    at: 2026-10-08T04:33:46Z
    by: agent-S-0324
  - to: done
    at: 2026-10-08T04:34:30Z
    by: orchestrator
tags: []
topics: [cli]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/www.md, flai/internal/mdlint/testdata/cases/expected.txt, design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/summary.md, design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1963
  turns:
    - day: 2026-10-08
      ceremony: 8
      test_runs: 9
      hand_edits: 5
      work: 64
  models:
    - model: claude-opus-5-5
      input: 182
      output: 51721
      cache_read: 15083094
      cache_write: 411547
      cost: 7.3441
  strategic:
    - kind: orchestrator
      seconds: 940
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 52
          output: 664
          cache_read: 17821630
          cache_write: 48104
          cost: 4.4028
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T18:59:55Z
  value: 12.5
  by: planner-S-0324
  at: 2026-10-08T00:29:11Z
forecast:
  duration: 12m
  delivery: 2026-10-08T00:44:00Z
  basis: "Its own forecast of 12m; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0316 and S-0317."
  by: flai
  at: 2026-10-08T00:29:58Z
finalized:
  by: orchestrator
  at: 2026-10-08T00:29:52Z
---
# S-0324 flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL

## Goal

This story remediates [I-0110](../../../design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md), "flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0110 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0110 is closed with `flai issue close I-0110 --reason` saying what fixed it

## Tasks
- T-1317 flai's MD034 reports a bare `www.` literal as markdownlint does
- T-1318 Close I-0110 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0110. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T08:55:00Z, 0.4 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0110's instance: in `parseRange` in `flai/internal/mdlint/inline.go`, the `c == 'w'` case skips a GFM `www.` literal without recording it. Record its position in `out.urls`, as the `c == 'h'` case does, gated the same way on `!link && !unclosed`, so `md034` reports it. A new fixture, `www.md`, settles the edges against markdownlint-cli2 0.20.0, and a guard test reproduces the issue.

Touches, all files, no folder touch:

| Touch | Source |
|-------|--------|
| `flai/internal/mdlint/inline.go` | design: I-0110 names `parseRange`'s `www.` case |
| `flai/internal/mdlint/mdlint_test.go` | co-change: 3 of 3 commits that changed `inline.go` |
| `flai/internal/mdlint/testdata/cases/www.md` | design: I-0110 asks for a `www.` fixture; a new file, named now |
| `flai/internal/mdlint/testdata/cases/expected.txt` | co-change: 3 of 3; `scripts/mdlint-fixtures.sh` regenerates it |
| `design/issues/I-0110-…-bare-url.md` | layout: criterion 2 closes it |
| `design/issues/summary.md` | layout: `flai issue close` rewrites it |

Left out: `doc.go`, `mdlint.go`, and `rules.go`, which `flai touches suggest` lists at 67%. `md034` in `rules.go` already reports every entry in `out.urls`, so it needs no change. `docs/users/flai.md` already says flai checks bare URLs, which takes in a `www.` literal, so the docs need no change.

Forecast: flai gave 6m. It is raised to 12m because S-0265, the same kind of fix for email addresses in the same files, took 701 s of agent time with the fixture regenerated through `npx` and a close-out. The delivery, 2026-10-08T07:54:00Z, is flai's and stands: the 26 stories ahead in the pull order decide it, not this story's duration.

Cost of delay: flai cod gave 12.50 USD a week from the operator's input of 5m lost per 168h cycle at 150 USD an hour. The value stands: one occurrence, no penalty, no revenue.

Topics: `cli` was added, the topic of S-0265, the story that changed the same lint.

### Accepted by the orchestrator

- Verified: 54ec667e6884ed00d8f24573bce3cf4a00ef92ec
- At: 2026-10-08T04:34:30Z

Verdict: accept; both criteria met (verifier at 54ec667e6884ed00d8f24573bce3cf4a00ef92ec; flai verify passed every step at that commit). The I-0114 bump and the new I-0117, I-0118, I-0119 are the close-outs' own records, all within the touches.

- 1: flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/www.md, flai/internal/mdlint/testdata/cases/expected.txt
- 2: design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/summary.md
