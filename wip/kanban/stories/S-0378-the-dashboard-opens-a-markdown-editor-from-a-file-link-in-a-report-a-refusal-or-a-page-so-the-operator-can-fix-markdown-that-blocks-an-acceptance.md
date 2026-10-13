---
id: S-0378
type: story
nature: feature
title: The dashboard opens a markdown editor from a file link in a report, a refusal, or a page, so the operator can fix markdown that blocks an acceptance
status: backlog
owner: alex
created: 2026-10-12T17:07:24Z
updated: 2026-10-13T17:29:30Z
transitions: []
tags: [dashboard, cli]
topics: [markdown]
touches: [flai/cmd/doc.go, flai/cmd/doc_test.go, flai/cmd/doc_lint.go, flai/cmd/doc_lint_test.go, flai/cmd/accept_markdown.go, flai/cmd/accept.go, flai/cmd/accept_fix_test.go, flai/internal/docedit/docedit.go, flai/internal/docedit/docedit_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/docs.go, flai/internal/hostapi/docs_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/writes.test.ts, flaiover/src/routes/api/docs/edit/+server.ts, flaiover/src/routes/api/docs/file/+server.ts, flaiover/src/routes/api/docs/file/save.test.ts, flaiover/src/routes/api/docs/lint/+server.ts, flaiover/src/routes/api/docs/lint/lint.test.ts, flaiover/src/lib/components/DocEditor.svelte, flaiover/src/lib/components/DocEditor.svelte.test.ts, flaiover/src/lib/components/EditorDialog.svelte, flaiover/src/lib/components/EditorDialog.svelte.test.ts, flaiover/src/lib/editor.svelte.ts, flaiover/src/lib/edit.ts, flaiover/src/lib/edit.test.ts, flaiover/src/routes/+layout.svelte, "flaiover/src/routes/edit/[...path]/+page.svelte", flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts, flaiover/src/lib/markdown.ts, flaiover/src/lib/markdown.test.ts, flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/api/items/[id]/accept/+server.ts", "flaiover/src/routes/api/items/[id]/accept/accept.test.ts", flaiover/src/routes/board/+page.svelte, "flaiover/src/routes/items/[id]/+page.svelte", design/adrs, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
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
    - kind: orchestrator
      seconds: 615
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 85
          output: 1496
          cache_read: 7698952
          cache_write: 25719
          cost: 1.9476
cost_of_delay:
  inputs:
    penalty_per_week: 25
    by: orchestrator
    at: 2026-10-13T16:49:35Z
  value: 25
  by: planner-S-0378
  at: 2026-10-13T16:58:49Z
forecast:
  duration: 1h45m
  delivery: 2026-10-13T20:30:00Z
  basis: "Its own forecast of 1h45m; 10th in the pull order with an in-progress limit of 5, behind S-0356, S-0376, S-0379, S-0357, S-0241, S-0304, S-0305, S-0306, S-0358, S-0371, S-0372 and S-0374."
  by: flai
  at: 2026-10-13T17:26:49Z
finalized:
  by: alex
  at: 2026-10-13T17:29:30Z
---
# S-0378 The dashboard opens a markdown editor from a file link in a report, a refusal, or a page, so the operator can fix markdown that blocks an acceptance

## Goal

When an acceptance is refused because markdown the lint fix cannot fix still has findings (S-0367, TH-0413), the operator opens the file from the link the refusal names, corrects it in a markdown editor in the dashboard, saves it through flai, and accepts again.

## Acceptance criteria

- [ ] A repository file link shown inline in a report, a refusal, a thread, or a page opens the file in a markdown editor in the dashboard, at the line the link names when it names one.
- [ ] The editor shows the markdown lint findings of the file as it is edited, with their lines and rules.
- [ ] Saving writes the file through flai on the host (ADR-0023), on the branch or checkout the link names, and refuses a save that conflicts with a change made meanwhile.
- [ ] After the findings that blocked an acceptance are corrected and saved, the acceptance can be made again from the same page.

## Tasks
- T-1514 flai doc show and doc save take --story and edit the file in that story's worktree, committing on its branch
- T-1515 flai doc lint and the host method doc.lint answer the markdown findings of content as it is edited
- T-1516 flai accept's refusal for markdown the fix left carries each finding's file, line, rule, and checkout as data, through accept.run to the dashboard
- T-1517 The dashboard's document endpoints pass a story to doc.show and doc.save, and /api/docs/lint answers doc.lint
- T-1518 The document editor opens over any page at a line, on a story's branch when given one, and shows the lint findings as the operator types
- T-1519 Repository file links in rendered markdown, threads, reports, and plain refusal text open the editor at the line they name
- T-1520 A refused acceptance lists the markdown the fix left as editor links, and the operator accepts again from the same page once they are saved
- T-1521 An ADR refining ADR-0023 and the design and user docs say how the editor opens from a link, lints, and saves on a story's branch

## Notes

- Asked for by the operator on TH-0413, as a story of its own for the planner, beside S-0367, which makes `flai accept` refuse with the files, lines, and errors of markdown the fix cannot fix.

### Planning

- **Touches.** S-0378 declared none, so `flai touches suggest` had nothing to start from. Every touch is the planner's, from the layout, read by the explorer, and the design the story links:
  - `flai/cmd/doc.go`, `flai/internal/docedit/docedit.go`, and `flai/internal/hostapi/writes.go`, with their tests: `flai doc show` and `doc save`, and the host methods `doc.show` and `doc.save` (ADR-0023). None takes a worktree today. Criterion 3.
  - `flai/cmd/doc_lint.go` and `flai/internal/hostapi/docs.go`, with their tests: new. No command lints content today; `flai/internal/mdlint` does, for wip only. Criterion 2.
  - `flai/cmd/accept_markdown.go`, `flai/cmd/accept.go`, and `flai/cmd/accept_fix_test.go`: the ADR-0140 refusal, text today, gains its findings as data. Criterion 4.
  - `flaiover/src/lib/server/` and `flaiover/src/routes/api/docs/`: the endpoints behind `DocEditor.svelte` (S-0040), with `/api/docs/lint` new.
  - `flaiover/src/lib/components/DocEditor.svelte`, `EditorDialog.svelte` (new), `flaiover/src/lib/editor.svelte.ts` (new), `flaiover/src/lib/edit.ts`, `flaiover/src/routes/+layout.svelte`, and `flaiover/src/routes/edit/[...path]/+page.svelte`: the editor, at a line, over any page.
  - `flaiover/src/lib/filelinks.ts` (new), `flaiover/src/lib/markdown.ts`, `Threads.svelte`, and `AnalyzerPanel.svelte`: file links in pages, threads, and reports. Criterion 1.
  - `Review.svelte`, `flaiover/src/routes/api/items/[id]/accept/`, `flaiover/src/routes/board/+page.svelte`, and `flaiover/src/routes/items/[id]/+page.svelte`: where acceptances are refused today. Criterion 4.
  - The docs `documentation.md` asks for: `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, and the generated `docs/users/flai-reference.md`.
- **Co-change.** Rerun with these touches, `flai touches suggest` lists files that change often with anything: `docs/operators/index.md`, `docs/operators/settings.md`, `design/issues/summary.md`, and `design/system/workflow.md`. None is added: the story adds no setting and changes no workflow.
- **Folder touch kept.** `design/adrs`, because T-1521's ADR gets its file name only when `flai adr new` numbers it.
- **Overlaps.** `docs/users/flai.md` and `design/system/flai-cli.md` are claimed by S-0375 and by E-0019's stories, in progress. They are the usual shared docs, touched by T-1521 last.
- **Forecast.** `flai forecast` gave 1h24m: 95 s per unit of size over 50 done feature stories on claude-opus-5-5 in the large band, times size 53. Raised to 1h45m, and the delivery from 20:00Z to 20:25Z, because the story spans flai and flaiover: its tasks run both test tiers and add new Svelte components with component tests, which a rate per unit of size does not tell apart.
- **Cost of delay.** 25.00 USD a week, as `flai cod` works it out from the input the orchestrator set on TH-0435, `penalty_per_week` 25, the figure recommended there. It stands as worked out.
