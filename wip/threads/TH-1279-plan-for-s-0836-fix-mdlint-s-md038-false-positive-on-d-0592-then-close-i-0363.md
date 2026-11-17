---
id: TH-1279
title: "Plan for S-0836: fix mdlint's MD038 false positive on D-0592, then close I-0363"
anchor:
  path: wip/kanban/stories/S-0836-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md
  item: S-0836
status: answered
participants: [planner-S-0836, orchestrator]
created: 2026-11-17T08:29:41Z
updated: 2026-11-17T08:32:19Z
---

# TH-1279 Plan for S-0836: fix mdlint's MD038 false positive on D-0592, then close I-0363

On wip/kanban/stories/S-0836-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md.

## Entries

### 2026-11-17T08:29:41Z planner-S-0836
S-0836 is planned with two tasks in two layers.

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-3529: mdlint parses D-0592's line 17 as markdownlint-cli2 does, with a fixture that reproduces it | none | `flai/internal/mdlint/inline.go`, `rules.go`, `mdlint_test.go`, `testdata/cases/code-spans.md`, `testdata/cases/expected.txt` |
| 2 | T-3531: close I-0363 with the cause and the fix | T-3529, whose cause and fix its reason names | `design/issues/I-0363-…md`, `design/issues/summary.md` |

Figures: forecast 45m, up from flai forecast's 18m, because the cause must be found against markdownlint-cli2 0.20.0 before the fix and the close-out runs the integration tier. Delivery 2026-11-17T16:00:00Z. Cost of delay value 12.50 USD a week, from the operator's 5m-per-cycle input, unchanged.

Proposed solution: fix the lint, not the digest entry. markdownlint passes the line, so mdlint's finding is a false positive, and a rewritten entry would leave the next such line to fail.

Assumptions:

- The issue says mdlint ends the span at an escaped backtick. That account is unverified. CommonMark gives a backslash no escaping power inside a code span, so the gap may be in which backtick opens the next span, or in `md038`'s padding rule. T-3529 establishes the cause first, by running both linters on the line, and records it in the narrative.
- The fix stays within `flai/internal/mdlint`. If the cause is in `doc.go` or `mdlint.go`, T-3529 widens its touches.
- No `after` on another story and no folder touches: the story reads no other open story's output.
- No thread named "asked" exists on S-0836, and it had no earlier plan, so this is a first plan.

I would split, merge, or drop nothing.

### 2026-11-17T08:32:19Z orchestrator
Approved. The plan fits S-0836:

- T-3529 finds the cause against markdownlint-cli2 and adds a fixture and a test that reproduce it, which meets criterion 1. T-3531 closes I-0363, which meets criterion 2, and waits for T-3529 because its reason names the cause T-3529 finds.
- Fixing the lint rather than the digest entry is right: markdownlint passes the line, so mdlint's finding is a false positive.
- Every touch is a file under `flai/internal/mdlint` or the issue's own file. There is no folder touch and no `after`.
- The 45m forecast covers finding the cause and the integration tier. The 12.50 USD a week value comes from `flai cod`.

I am finalizing S-0836 under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator
