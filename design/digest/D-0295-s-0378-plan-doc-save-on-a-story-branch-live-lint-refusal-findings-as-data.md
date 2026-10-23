---
id: D-0295
title: "S-0378 plan: doc save on a story branch, live lint, refusal findings as data, editor over any page (new ADR), eight tasks in five layers: Recommendation: approve the plan as drafted, and finalize…"
kind: answer
sources: [TH-0436]
topics: [markdown, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd/doc.go, flai/cmd/doc_test.go, flai/cmd/doc_lint.go, flai/cmd/doc_lint_test.go, flai/cmd/accept_markdown.go, flai/cmd/accept.go, flai/cmd/accept_fix_test.go, flai/internal/docedit/docedit.go, flai/internal/docedit/docedit_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/docs.go, flai/internal/hostapi/docs_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/writes.test.ts, flaiover/src/routes/api/docs/edit/+server.ts, flaiover/src/routes/api/docs/file/+server.ts, flaiover/src/routes/api/docs/file/save.test.ts, flaiover/src/routes/api/docs/lint/+server.ts, flaiover/src/routes/api/docs/lint/lint.test.ts, flaiover/src/lib/components/DocEditor.svelte, flaiover/src/lib/components/DocEditor.svelte.test.ts, flaiover/src/lib/components/EditorDialog.svelte, flaiover/src/lib/components/EditorDialog.svelte.test.ts, flaiover/src/lib/editor.svelte.ts, flaiover/src/lib/edit.ts, flaiover/src/lib/edit.test.ts, flaiover/src/routes/+layout.svelte, "flaiover/src/routes/edit/[...path]/+page.svelte", flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts, flaiover/src/lib/markdown.ts, flaiover/src/lib/markdown.test.ts, flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/api/items/[id]/accept/+server.ts", "flaiover/src/routes/api/items/[id]/accept/accept.test.ts", flaiover/src/routes/board/+page.svelte, "flaiover/src/routes/items/[id]/+page.svelte", design/adrs, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, docs/operators/settings.md, flaiover/src/lib/components/DismissibleNotice.svelte, flaiover/src/lib/components/DismissibleNotice.svelte.test.ts, design/issues/I-0145-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md, design/issues/summary.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md]
date: 2026-10-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0295 S-0378 plan: doc save on a story branch, live lint, refusal findings as data, editor over any page (new ADR), eight tasks in five layers: Recommendation: approve the plan as drafted, and finalize…

The operator, alex, answered TH-0436 at 2026-10-14T18:23:40Z.

What was asked:

```text
**Plan:** eight tasks in five layers. Three are in flai, four in flaiover, and one is docs with an ADR refining ADR-0023. S-0378 stays a draft for you to finalize.
```

The answer:

```text
**Recommendation:** approve the plan as drafted, and finalize S-0378. The tasks cover the four criteria: T-1519 for links (1), T-1515 and T-1518 for live findings (2), T-1514 and T-1517 for saving through flai on the branch (3), and T-1520 for accepting again (4). The figures are flai's: 1h45m (1h24m raised a quarter for two tiers) and 25 USD/week from the input set on TH-0435.

I leave the thread open for you because it holds two decisions not yet recorded:

- **Saving on a story's branch.** ADR-0023 § Decision commits a dashboard save "on the main checkout's current branch". Under T-1514 and T-1521, a save with `story` commits that one path on `story/S-nnnn`, in a worktree a story's agent may be working in, through a new ADR. My recommendation is yes: criterion 3 asks for it, and the save keeps ADR-0023's hash conflict and strict check.
- **A new exit code for `flai accept --json`.** T-1516 gives the ADR-0140 refusal exit 4, with the findings as data. That is a contract the dashboard reads, so T-1521's ADR or the design should record it too.

The one folder touch, `design/adrs`, stays only for the ADR's unknown file name.

Source: design/adrs/0023-documents-are-saved-through-flai.md § Decision
```

Source: TH-0436.
