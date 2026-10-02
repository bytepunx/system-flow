---
id: TH-0060
title: "S-0194: ADR-0066 and the experiment convention edits are uncommitted in the main checkout; may they be committed to main so my branch can sync them?"
anchor:
  path: wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md
  item: S-0194
status: resolved
participants: [agent-S-0194, alex]
created: 2026-10-02T10:07:51Z
updated: 2026-10-02T10:32:09Z
---

# TH-0060 S-0194: ADR-0066 and the experiment convention edits are uncommitted in the main checkout; may they be committed to main so my branch can sync them?

On wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md.

## Entries

### 2026-10-02T10:07:51Z agent-S-0194
The story is built on ADR-0066 and on your edits to `git.md` and `documentation.md` (here and in `template/root/`). Right now they exist only in the main checkout: ADR-0066 to ADR-0069 are untracked, and the conventions, `design/adrs/README.md`, and others are modified but not committed. My worktree branches from `main` HEAD (05b67f9), so it has neither the ADR nor the new convention text. That matters in two ways. First, the design documents I have to edit link ADR-0066, and `flai check` in the worktree would report those links as broken. Second, if I copy the ADR or the convention text onto my branch, accepting S-0194 would merge files that are also uncommitted in the main checkout, and git refuses that kind of merge.

Recommended: commit your uncommitted design edits to `main` yourself (ADRs 0066–0069, the conventions here and in `template/root/`, `design/adrs/README.md`, `design/system/continuous-improvement.md`). Then I run `flai stream sync S-0194` to pick them up, and my branch leaves those files alone.

Alternative: I copy only ADR-0066 onto my branch, so its links resolve, and leave the conventions to you. Before accepting, you would have to delete the untracked copy in main, or git refuses the merge.

I'm working on the code (accept, release plan, check, dashboard) in the meantime, since none of it depends on the answer.

### 2026-10-02T10:11:37Z alex
done, I have added the items and committed locally in main

### 2026-10-02T10:13:55Z agent-S-0194
Thank you. 2d4e090 has reached my branch. It carries the template's conventions, `design/adrs/README.md`, and `design/system/continuous-improvement.md`, but not the ADR files themselves or this repository's own conventions. In the main checkout these are still uncommitted: `design/adrs/0066-…` to `0069-…` (untracked), and `design/conventions/*.md` (modified). As a result, `design/adrs/README.md` on `main` now links four ADRs that are not committed, and my design edits that link ADR-0066 point to a file not yet on the branch.

Could you commit those too: the four ADR files and `design/conventions/*.md`? I'll sync once they're on `main`. In the meantime I'm going on with the dashboard and the docs, and my links to ADR-0066 use its path in `main`, so they resolve once it is committed.

### 2026-10-02T10:17:11Z alex
done

### 2026-10-02T10:21:28Z agent-S-0194
Thank you. f45f959 brought ADR-0066, ADR-0069, and the conventions onto my branch, and every link from S-0194 now resolves. One thing is still open: `design/adrs/0067-…` and `0068-…` are untracked in the main checkout, while `design/adrs/README.md` on `main` indexes both. `flai check --strict` therefore warns `adr.index` twice, on `main` and on every story branch, and S-0194's close-out stops on it. Could you commit those two as well? If ADR-0067 is meant to stay proposed for now, that is fine: I changed the dashboard's ADR-list test, which required every ADR in this repository to be accepted, so that it allows a proposed one (T-0687). I'll run the close-out once they are on `main`.

### 2026-10-02T10:22:55Z alex
done

### 2026-10-02T10:32:09Z agent-S-0194
Resolved: The designer committed ADR-0066 to ADR-0069 and the convention edits to main (2d4e090, f45f959, b541527); S-0194 is synced onto them and flai check --strict is clean.
