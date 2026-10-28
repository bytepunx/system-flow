---
id: D-0299
title: "S-0374 plan: the close-out syncs and verifies again, reusing the tiers main's new paths do not select (new ADR), five tasks in four layers: Recommendation: approve the plan as drafted, with its…"
kind: answer
sources: [TH-0424]
topics: [cli, conventions, template, testing, flai, go, code, all]
paths: [design/adrs/README.md, flai/internal/verify/story.go, flai/internal/verify/paths.go, flai/internal/verify/story_test.go, flai/internal/verify/paths_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md, design/issues/summary.md, design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md, docs/users/flai-reference.md, flai/cmd/verify.go]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md#decision"
asked_again: []
cited: []
---

# D-0299 S-0374 plan: the close-out syncs and verifies again, reusing the tiers main's new paths do not select (new ADR), five tasks in four layers: Recommendation: approve the plan as drafted, with its…

The operator, alex, answered TH-0424 at 2026-10-14T20:56:59Z.

What was asked:

```text
**Recommendation:** approve the remedy below, which a new ADR refining ADR-0135 and ADR-0110 records.
```

The answer:

```text
**Recommendation:** approve the plan as drafted, with its three-round bound and the baseline change to `work-management.md` in both copies.

## Why

- **The remedy keeps ADR-0135's safeguard.** It does not pass over main's commits outright, because a main change can break the branch without touching its files. It reuses only tiers whose selection main's new paths do not touch, `{importers}` included, and only while the branch's own patch is unchanged (`git patch-id --stable`).
- **It cuts I-0136's cost.** That is about six minutes of tiers run again per acceptance during a close-out, three times within 1.5 hours.
- **The bound and the prefix rule keep it predictable:** at most three rounds, and once one tier runs again, every tier after it runs.

## Why this is yours

It rewrites a baseline rule that the template releases to every project, and it adds an ADR refining ADR-0135 and ADR-0110. The figures are flai's: 1h, raised from 23m with a stated reason, and 45 USD a week from flai's 18m input. Every touch is a file. I leave S-0374 a draft until you answer.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md
```

Source: TH-0424.
