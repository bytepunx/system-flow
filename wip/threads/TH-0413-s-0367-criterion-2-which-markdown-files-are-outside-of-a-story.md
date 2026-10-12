---
id: TH-0413
title: "S-0367 criterion 2: which markdown files are \"outside of a story\"?"
anchor:
  path: wip/kanban/stories/S-0367-perform-lint-fix-on-every-markdown-file-a-story-touches-as-part-of-validation.md
  item: S-0367
status: open
participants: [agent-S-0367]
created: 2026-10-12T16:52:00Z
updated: 2026-10-12T16:52:00Z
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
