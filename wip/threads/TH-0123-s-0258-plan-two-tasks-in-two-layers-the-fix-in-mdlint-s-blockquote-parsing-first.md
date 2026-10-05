---
id: TH-0123
title: "S-0258 plan: two tasks in two layers, the fix in mdlint's blockquote parsing first"
anchor:
  path: wip/kanban/stories/S-0258-flai-s-wip-markdown-lint-does-not-flag-an-ordered-list-item-numbered-from-other-than-1-inside-a-blockquote-so-a-thread-entry-reached-main-and-failed-a-story-s-close-out.md
  item: S-0258
status: open
participants: [planner-S-0258]
created: 2026-10-05T03:13:59Z
updated: 2026-10-05T03:13:59Z
---

# TH-0123 S-0258 plan: two tasks in two layers, the fix in mdlint's blockquote parsing first

On wip/kanban/stories/S-0258-flai-s-wip-markdown-lint-does-not-flag-an-ordered-list-item-numbered-from-other-than-1-inside-a-blockquote-so-a-thread-entry-reached-main-and-failed-a-story-s-close-out.md.

## Entries

### 2026-10-05T03:13:59Z planner-S-0258
**Remedy proposed:** `flai/internal/mdlint/doc.go` marks quoted lines `kQuote` and never parses what is inside them, so MD029 never sees `> 3.`. The fix parses a blockquote's content as blocks of their own, with the `>` markers taken out. Findings are reported on the original lines, so MD029 judges quoted ordered lists as markdownlint-cli2 does.

**Tasks and layers:**

1. Layer 1, **T-0860** (touches `flai/internal/mdlint`):
   - Adds the parser change and a fixture, `testdata/cases/blockquote.md`, made from TH-0101's `> 3.` and `> 6.` entry.
   - Regenerates `expected.txt` with `scripts/mdlint-fixtures.sh`.
   - Also checks that MD004, MD007, MD022, and MD031 agree on quoted content, and that the repository still lints clean.
2. Layer 2, **T-0861** (touches `design/system/flai-cli.md` and `design/issues`, after T-0860):
   - Updates the design's `mdlint/` line.
   - Runs `flai issue close I-0070 --reason` from the worktree.

**Story enrichment:**

- **Touches:** `flai/internal/mdlint`, `design/issues`, `design/system/flai-cli.md`. The story declared none.
- **Forecast:** 30m, raised from flai's 8m, using S-0240's 21m for an earlier mdlint fix as the reference. Delivery about 03:50Z.
- **Cost of delay:** 25.00 USD/week, from flai's inputs, unadjusted.

**Assumptions:**

- markdownlint-cli2 at `lint-md.sh`'s version flags `> 3.` under MD029 style one_or_ordered, as I-0070's instance says. T-0860 verifies this by regenerating the reference.
- npx is available to the story agent to run `scripts/mdlint-fixtures.sh`.
- `docs/users/flai.md` needs no change, since it describes the lint only as checking "lists".
- One task holds both the fix and its reproducing fixture, so that the branch's tests are never red between tasks.

Nothing is proposed to split, merge, or drop.
