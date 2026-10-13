---
id: TH-0436
title: "S-0378 plan: doc save on a story branch, live lint, refusal findings as data, editor over any page (new ADR), eight tasks in five layers"
anchor:
  path: wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md
  item: S-0378
status: open
participants: [planner-S-0378, orchestrator]
created: 2026-10-13T16:59:05Z
updated: 2026-10-13T16:59:18Z
---

# TH-0436 S-0378 plan: doc save on a story branch, live lint, refusal findings as data, editor over any page (new ADR), eight tasks in five layers

On wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md.

## Entries

### 2026-10-13T16:59:05Z planner-S-0378
**Plan:** eight tasks in five layers. Three are in flai, four in flaiover, and one is docs with an ADR refining ADR-0023. S-0378 stays a draft for you to finalize.

## Tasks and layers

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-1514 | `flai doc show` and `doc save` take `--story`, edit in its worktree, and commit on its branch; `doc.show` and `doc.save` take `story` | none |
| 1 | T-1515 | `flai doc lint` and the read method `doc.lint` lint content with flai's mdlint and the project's config | none |
| 2 | T-1516 | `flai accept --json` refuses with exit 4 and the ADR-0140 findings as data (path, line, rule, checkout), passed through `accept.run` | T-1514 (`hostapi/writes.go`) |
| 2 | T-1517 | dashboard endpoints pass `story`; `POST /api/docs/lint` | T-1514, T-1515 |
| 2 | T-1518 | `DocEditor` at a line, with a story, live findings; `EditorDialog` over any page from `+layout.svelte` | T-1515 |
| 3 | T-1519 | `filelinks.ts` and `markdown.ts`: repository paths in text and rendered markdown open the editor at their line; threads and reports | T-1518 |
| 4 | T-1520 | the review page lists refusal findings as editor links, and Accept is enabled again after a save; the board and item notices are linked too | T-1516, T-1517, T-1519 |
| 5 | T-1521 | ADR refining ADR-0023, `flaiover-dashboard.md`, `flai-cli.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, regenerated reference | T-1520 |

Every touch names a file except `design/adrs`, kept for the ADR's unknown file name.

## Figures

- Forecast 1h45m, delivery 2026-10-13T20:25Z: `flai forecast`'s 1h24m, raised a quarter because the story runs both the Go and the flaiover tiers.
- Cost of delay 25.00 USD a week, from the `penalty_per_week` 25 set on TH-0435.

## Assumptions

1. **The editor opens over the page.** It opens in a dialog, not on a page of its own, so after a save you are still on the review page, and Accept works again there (criterion 4). `/edit/<path>?line=&story=` still serves a new tab.
2. **Ordinary document links keep showing the rendered document.** A relative link to a markdown file, or a code-span path, gains an edit control beside it. A path in plain refusal text opens the editor directly.
3. **Live lint is flai's own subset** of the markdownlint rules (ADR-0061), fast enough to run while you type. The markdown tier (`scripts/lint-md.sh`) stays the authority at acceptance, and the editor says so. Running markdownlint-cli2 on every pause would start npx each time.
4. **"The branch the link names" means the story's worktree** for files the branch changes, and the main checkout otherwise. A save with `--story` commits that one path on `story/S-nnnn`. A wip path is never saved there, since wip lives on main.
5. **Saving on a story branch refines ADR-0023**, which commits on the main checkout only, so T-1521 writes an ADR for it.
6. **`doc.lint` is optional for the dashboard.** An older flai on the host still connects; the editor then shows no live findings.

Nothing to split, merge, or drop.

### 2026-10-13T16:59:18Z orchestrator (recommendation)
**Recommendation:** approve the plan as drafted, and finalize S-0378. The tasks cover the four criteria: T-1519 for links (1), T-1515 and T-1518 for live findings (2), T-1514 and T-1517 for saving through flai on the branch (3), and T-1520 for accepting again (4). The figures are flai's: 1h45m (1h24m raised a quarter for two tiers) and 25 USD/week from the input set on TH-0435.

I leave the thread open for you because it holds two decisions not yet recorded:

- **Saving on a story's branch.** ADR-0023 § Decision commits a dashboard save "on the main checkout's current branch". Under T-1514 and T-1521, a save with `story` commits that one path on `story/S-nnnn`, in a worktree a story's agent may be working in, through a new ADR. My recommendation is yes: criterion 3 asks for it, and the save keeps ADR-0023's hash conflict and strict check.
- **A new exit code for `flai accept --json`.** T-1516 gives the ADR-0140 refusal exit 4, with the findings as data. That is a contract the dashboard reads, so T-1521's ADR or the design should record it too.

The one folder touch, `design/adrs`, stays only for the ADR's unknown file name.

Source: design/adrs/0023-documents-are-saved-through-flai.md § Decision
