---
id: TH-0454
title: "S-0383 was never replanned: criterion 3 saves front matter through doc.save, but the tokens are a DTCG JSON file and file.save does not exist yet"
anchor:
  path: wip/kanban/stories/S-0383-the-operator-changes-the-design-system-s-tokens-on-the-design-page-sees-the-preview-follow-before-saving-and-saves-through-flai.md
  item: S-0383
status: resolved
participants: [agent-S-0383, orchestrator, alex]
created: 2026-10-15T01:57:56Z
updated: 2026-10-15T02:54:34Z
---

# TH-0454 S-0383 was never replanned: criterion 3 saves front matter through doc.save, but the tokens are a DTCG JSON file and file.save does not exist yet

On wip/kanban/stories/S-0383-the-operator-changes-the-design-system-s-tokens-on-the-design-page-sees-the-preview-follow-before-saving-and-saves-through-flai.md.

## Entries

### 2026-10-15T01:57:56Z agent-S-0383
**Recommendation:** let me build S-0383 against ADR-0147's contract. That means a dashboard route `/api/design/tokens` that reads through `file.get` and saves through `file.save`, and answers 501 until the host flai has them. Rewrite criterion 3 to say so. Make the `file.get`/`file.save` story (design-system-tool.md § A new story: file.get and file.save) to build the flai side, and verify the save end to end there.

## What I found

- **Criterion 3 cannot be met as written.** It says the save writes "the tokens into the document's front matter through flai (`doc.save`, ADR-0023)". S-0381, accepted, made the tokens a DTCG JSON file (ADR-0146), and `doc.save` refuses a file that is not `.md`.
- **The method the design names does not exist.** design-system-tool.md § Epic › S-0383 Editing on the page says the save goes through `file.save` on main, with flai writing the stylesheet in the same commit (ADR-0147, ADR-0150). `flai/internal/hostapi` has `design.get` and `doc.save`, but no `file.get` or `file.save`. No story on the board builds them.
- **The planner never applied § Epic to S-0383.** Its title, criteria, tasks, and touches are the draft from before TH-0445. TH-0438 asked for that replan.

## Options

1. **(Recommended) Build the dashboard side now against ADR-0147.**
   - T-1543 edits the DTCG JSON and keeps the keys it does not know.
   - A new task adds `flaiover/src/routes/api/design/tokens/+server.ts` with its test. It calls `file.get`/`file.save` and answers 501 from an older flai.
   - The page shows a 501 as "this flai cannot save tokens yet" and keeps the edits.
   - Criterion 3 is verified against a stubbed host. The end-to-end save is verified in the file story.
   - You rewrite criterion 3. If you want ADR-0151's "unsaved values kept across a reload", add that criterion too.
2. **Widen S-0383 to build `file.get`/`file.save` in flai as well.** That adds `flai/cmd/file.go`, `docedit`, `hostapi/writes.go`, and the docs. It roughly doubles the story, and it claims `flai/cmd/root.go`.
3. **Block S-0383** until the planner replans E-0020 and the file story is accepted.

## What I am doing meanwhile

I am working the parts that hold under every option:
- the edit model on the DTCG JSON (T-1543, retargeted from front matter)
- the in-browser contrast check (T-1544)
- the token editor (T-1545)

The save in T-1546 waits for your answer.

### 2026-10-15T01:58:09Z orchestrator (recommendation)
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

### 2026-10-15T02:38:32Z alex
Confirmed the recommendation of 2026-10-15T01:58:09Z orchestrator.

Source: design/adrs/0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md § Decision

### 2026-10-15T02:54:34Z agent-S-0383
Resolved: Option 1 is built. T-1594 added /api/design/tokens, which calls file.get and file.save and answers 501 from an older flai. Criterion 3 now has the confirmed wording. E-0020's planner still has to draft the flai story for file.get/file.save.
