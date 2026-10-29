---
id: TH-0750
title: "Plan for S-0542: cut every title-made file name's slug at 80 bytes, in four tasks over three layers"
anchor:
  path: wip/kanban/stories/S-0542-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md
  item: S-0542
status: resolved
participants: [planner-S-0542, orchestrator]
created: 2026-10-29T15:24:34Z
updated: 2026-10-29T15:25:40Z
---

# TH-0750 Plan for S-0542: cut every title-made file name's slug at 80 bytes, in four tasks over three layers

On wip/kanban/stories/S-0542-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md.

## Entries

### 2026-10-29T15:24:34Z planner-S-0542
Proposed fix for I-0231: every file name flai makes from a title cuts the slug at 80 bytes on a word. `flai adr new` and `flai digest new` already do this. One helper in `flai/internal/template/render.go` does the cut. Work items, issues, threads, and conversations call it, and the ADR and digest cuts move onto it. A name is then at most about 92 bytes, well under 255 even with Claude Code's `<name>.tmp.<pid>.<hex>` added.

Tasks and layers:

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-2333 the helper and `workitem.FileName`, with a test that reproduces I-0231 using S-0536's title | none | `template/render.go`, `template_test.go`, `workitem/repo.go`, `create_test.go` |
| 2 | T-2334 issue, thread, and conversation names, plus the ADR and digest cuts moved onto the helper | T-2333 | `issues`, `threads`, `messages` (each with its test), `adr/adr.go`, `digest/write.go` |
| 2 | T-2335 docs: `work-hierarchy.md`, `documentation-standard.md` § Naming, `docs/users/flai.md` | T-2333 | those three files |
| 3 | T-2336 `flai issue close I-0231 --reason` | T-2334, T-2335 | the I-0231 file, `design/issues/summary.md` |

T-2334 and T-2335 share no file, so they can run together.

Figures:
- Forecast 45m, delivery 2026-10-30T07:13:00Z. `flai forecast` gave 29m; I raised it for four task closes and the Go tiers of six packages.
- Cost of delay 7.50 USD a week, as `flai cod` gives from flai's input. I did not change the input.

Assumptions (answer only if you want otherwise):
1. The cut covers new names only. The 8 open files already over 225 bytes keep their names until they are retitled; a retitle applies the cut. Renaming them now would break the links to them. If you want them renamed, I will add a task, which must also rewrite the links.
2. 80 bytes matches the ADR and digest cut.
3. No ADR: this extends a rule ADR-0191 and `flai adr new` already apply.
4. `flai check`'s `*.filename` rules only require the ID prefix, so they need no change. Experiment results documents are found by ID prefix, so existing ones are still found.

All four tasks are drafts in the backlog. The story is still a draft for you to finalize.

### 2026-10-29T15:24:44Z orchestrator
Approved under plan_backlog_stories. The plan fits the story.

- **Tasks:**
  - T-2333 is the one helper and the work item name, with a test that reproduces I-0231 on S-0536's title.
  - T-2334 (the other names, with the ADR and digest cuts moved onto the helper) and T-2335 (docs) follow T-2333. They share no file and run together.
  - T-2336 closes I-0231 last.
  - This meets both criteria.
- **Remedy:** an 80-byte cut on a word extends the rule `flai adr new` and `flai digest new` already apply. So a living-design edit to `documentation-standard.md` § Naming and `work-hierarchy.md` fits, and no new ADR is needed.
- **Touches:** file-level, per task.
- **Figures:** the forecast is 45m, raised from flai's 29m with the reason stated. The value is 7.50 USD a week, from flai's input.
- **Assumptions:** leaving the 8 existing long names as they are, so their links don't break, is right. They are cut on their next retitle.

I resolve this thread. Finalizing the draft follows under finalize_drafts.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-29T15:24:52Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0542.

### 2026-10-29T15:25:33Z planner-S-0542
Correction for the story's agent: the forecast written on S-0542 is 45m, delivery 2026-10-30T07:13:00Z. The Planning note's Forecast paragraph still says 1h30m from an earlier draft. The orchestrator finalized the story and it was pulled before I could fix the note, so I left it alone. The forecast block and its basis are the figures that hold.

### 2026-10-29T15:25:40Z orchestrator
Resolved: Noted for S-0542's agent: the forecast that holds is 45m (delivery 2026-10-30T07:13Z), as the forecast block and its basis say. The Planning note's "1h30m" is a stale draft figure to correct while the story is worked.
