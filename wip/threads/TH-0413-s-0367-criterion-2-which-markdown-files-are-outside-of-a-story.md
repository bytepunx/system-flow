---
id: TH-0413
title: "S-0367 criterion 2: which markdown files are \"outside of a story\"?"
anchor:
  path: wip/kanban/stories/S-0367-perform-lint-fix-on-every-markdown-file-a-story-touches-as-part-of-validation.md
  item: S-0367
status: resolved
participants: [agent-S-0367, alex]
created: 2026-10-12T16:52:00Z
updated: 2026-10-12T17:07:30Z
---

# TH-0413 S-0367 criterion 2: which markdown files are "outside of a story"?

On wip/kanban/stories/S-0367-perform-lint-fix-on-every-markdown-file-a-story-touches-as-part-of-validation.md.

## Entries

### 2026-10-12T16:52:00Z agent-S-0367
**Recommendation:** read "markdown files outside of a story that get added" as the markdown that `flai accept` commits from the main checkout, the `wip/` items, narratives, threads, and messages flai writes there. `flai accept` runs the markdown fix over the main checkout's changed markdown just before its `git add -A` and commit.

## The plan for S-0367

- A test tier in `system-flow.yaml` gets an optional `fix_command`, such as `[scripts/lint-md.sh, --fix, "{files}"]`, run over the files the tier's paths select.
- Criterion 1: `flai task done` runs it on the paths it is about to commit, and every clean `flai stream sync` runs it on every file the branch changes against main, committing what it fixed.
- Criterion 3: that sync step runs right after the rebase, so a file main changed and the branch changes is fixed as it now stands.
- Criterion 2: `flai accept` runs it in the main checkout before it commits, as above.

A fix that fails or leaves findings never stops a commit, a sync, or an acceptance; the markdown lint tier still reports what is left.

## The question

Is that what criterion 2 means? If not, name which additions you meant, for example commits made by hand on main, or `--autocommit` writes by `flai issue new` and `flai adr new`. I am building criteria 1 and 3, and the accept step, now. If your answer differs, I will change the accept step.

### 2026-10-12T16:58:39Z alex
that's correct. I do have a one adjustment, however. If the markdown can't be fixed, it should block the review being accepted with the files, lines, and errors so that the user can correct them. Lets also make a story that adds a markdown editor that can be launched from inlined file links in reports and other areas of the site so that the operator can correct the issue and unblock the review. I understand this is a significant increase in scope, so the editor can become a separate story for the planner.

### 2026-10-12T17:07:29Z agent-S-0367
Recorded in S-0367's narrative, under `## Decisions`.

- **The block.** `flai accept` will refuse an acceptance when the markdown tier still fails after the fix. The refusal names each file, line, rule, and message. It checks in two places:
  - **At preflight**, the main checkout's uncommitted markdown, before anything changes.
  - **After the rebase and its fix**, the files the branch changes, before the merge.
  
  The story stays in review until you correct the files and accept again.
- **The editor.** It is the draft story S-0378, with no epic, for the planner to plan once you finalize it.

### 2026-10-12T17:07:30Z agent-S-0367
Resolved: Criterion 2 is the main checkout's markdown that flai accept commits. Markdown the fix cannot fix blocks the acceptance with its findings. The editor is S-0378.
