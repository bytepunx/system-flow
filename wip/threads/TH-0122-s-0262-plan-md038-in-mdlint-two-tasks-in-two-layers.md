---
id: TH-0122
title: "S-0262 plan: MD038 in mdlint, two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md
  item: S-0262
status: open
participants: [planner-S-0262]
created: 2026-10-05T03:13:40Z
updated: 2026-10-05T03:13:40Z
---

# TH-0122 S-0262 plan: MD038 in mdlint, two tasks in two layers

On wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md.

## Entries

### 2026-10-05T03:13:40Z planner-S-0262
The fix: add MD038, spaces inside a code span, to `flai/internal/mdlint`, checked against markdownlint-cli2 on fixtures. This is how S-0240 added MD007.

Tasks, in two layers:

1. Layer 1, **T-0856**: mdlint reports MD038 as markdownlint does. It touches `flai/internal/mdlint`.
   - It adds the rule and a fixture of T-0822's original body with code-span cases.
   - It regenerates `expected.txt` with `scripts/mdlint-fixtures.sh` and adds a test that the guard refuses the body.
2. Layer 2, **T-0857**: the design and the user guide name MD038, and I-0072 is closed. It touches `design/system/flai-cli.md`, `docs/users/flai.md` and `design/issues`. It waits for T-0856 because it describes that rule and closes the issue only once the rule is in.

Enrichment:

- Touches: `flai/internal/mdlint`, `design/system/flai-cli.md`, `docs/users/flai.md` and `design/issues`. The story is no longer held for want of touches.
- Tag: `flai`.
- Topics: `cli` and `go`.
- Forecast: 20m, delivery about 03:40Z. I raised flai's 9m to the 21m S-0240 took.
- Cost of delay: 25.00 USD a week, as `flai cod` gives from your input.

Assumptions:

- MD038 follows markdownlint-cli2 0.20.0. A span of only spaces and a span padded by one space on both sides are not findings.
- ADR-0061 stays as it is, and no new ADR is needed, as with MD007.
- S-0258, the MD029-in-blockquote fix, will claim `flai/internal/mdlint` too once it declares touches. The two will then run one after the other, not together. Merging them would also do, but I propose keeping them apart: their issues are separate.
