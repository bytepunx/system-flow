---
id: TH-0088
title: S-0201 and S-0203 conflict when merged
anchor:
  path: wip/kanban/stories/S-0203-a-story-created-from-an-issue-is-a-draft-and-carries-the-issue-s-cost-of-delay-inputs.md
  item: S-0203
status: resolved
participants: [flai, agent-S-0203, alex]
created: 2026-10-03T18:10:44Z
updated: 2026-10-03T18:33:22Z
---

# TH-0088 S-0201 and S-0203 conflict when merged

On wip/kanban/stories/S-0203-a-story-created-from-an-issue-is-a-draft-and-carries-the-issue-s-cost-of-delay-inputs.md.

## Entries

### 2026-10-03T18:10:44Z flai
A trial merge of story/S-0201 with story/S-0203 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md`
- `design/issues/summary.md`
- `design/system/flai-cli.md`
- `design/system/work-hierarchy.md`

Whichever of S-0201 and S-0203 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-03T18:10:57Z agent-S-0203
None of these five paths is in S-0203's diff: story/S-0203 changes only flai/internal/itemnew and flai/internal/workitem/create.go (and their tests) over main. The conflicts are main's own changes from S-0200's acceptance against story/S-0201, which has not synced since; TH-0087 settled them (S-0201 keeps both sides when it rebases). S-0203 neither narrows nor waits. Recorded as an issue: the trial merge compares the other branch with this one including main, so it names a story for conflicts it did not make.

### 2026-10-03T18:33:22Z alex
Resolved.
