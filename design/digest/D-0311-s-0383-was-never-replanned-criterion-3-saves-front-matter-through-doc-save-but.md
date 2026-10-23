---
id: D-0311
title: "S-0383 was never replanned: criterion 3 saves front matter through doc.save, but the tokens are a DTCG JSON file and file.save does not exist yet: Recommendation: take option 1, and do two things on…"
kind: answer
sources: [TH-0454]
topics: [design, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/lib/design/tokens.ts, flaiover/src/lib/design/tokens.test.ts, flaiover/src/lib/design/preview.ts, flaiover/src/lib/design/preview.test.ts, flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/design/tokens/+server.ts, flaiover/src/routes/api/design/tokens/tokens.test.ts, design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md, design/issues/summary.md]
date: 2026-10-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0311 S-0383 was never replanned: criterion 3 saves front matter through doc.save, but the tokens are a DTCG JSON file and file.save does not exist yet: Recommendation: take option 1, and do two things on…

The operator, alex, answered TH-0454 at 2026-10-15T02:38:32Z.

What was asked:

```text
**Recommendation:** let me build S-0383 against ADR-0147's contract. That means a dashboard route `/api/design/tokens` that reads through `file.get` and saves through `file.save`, and answers 501 until the host flai has them. Rewrite criterion 3 to say so. Make the `file.get`/`file.save` story (design-system-tool.md § A new story: file.get and file.save) to build the flai side, and verify the save end to end there.
```

The answer:

```text
**Recommendation:** take option 1, and do two things on your side.

1. **Rewrite S-0383's criterion 3.** Proposed wording: "Save writes the edited tokens to the DTCG token file the manifest names, through the host method `file.save` (ADR-0147), on the main branch. A host flai without it answers 501, and the page says so and keeps the edits." Add "unsaved values survive a reload" (ADR-0151) if you want it.
2. **Start E-0020's planner** (Plan, or `flai plan E-0020`) to draft the `file.get`/`file.save` story, as recommended on TH-0438. The end-to-end save is verified there.

## Why option 1

- **ADR-0147 § Decision already fixes the contract.** It defines `file.get {path, story?}` → content and SHA-256, and `file.save {path, content, hash, story?, message?}` with doc.save's hash, check, and commit rules. The dashboard can build against that exactly, and the 501 path follows how the dashboard already treats a host method an older flai lacks.
- **Option 2 doubles the story** and claims `flai/cmd/root.go` and `hostapi/writes.go`. S-0385 and S-0387 touch those files too, so it would hold them.
- **Option 3 stops work** that holds under every option. The agent is already doing that work (T-1543 to T-1545).

## Why this needs you

- **S-0383 is finalized.** Its criteria are yours to change, not mine or its agent's.
- **The file story is a new story under E-0020.** Only the epic's planner drafts it, and I cannot start that planner: `flai plan --candidates` does not list E-0020.

While T-1546 waits, the agent carries on with the parts that hold under every option.

Source: design/adrs/0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md § Decision
```

Source: TH-0454.
