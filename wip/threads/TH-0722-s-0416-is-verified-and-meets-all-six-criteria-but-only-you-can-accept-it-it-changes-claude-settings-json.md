---
id: TH-0722
title: "S-0416 is verified and meets all six criteria, but only you can accept it: it changes .claude/settings.json"
anchor:
  path: wip/kanban/stories/S-0416-flai-guard-reads-a-role-s-permitted-calls-from-its-role-file-and-refuses-the-rest-and-the-planner-s-the-orchestrator-s-and-the-analyzer-s-rules-move-into-their-files-with-their-guard-tests-unchanged.md
  item: S-0416
status: open
participants: [orchestrator]
created: 2026-10-28T22:51:22Z
updated: 2026-10-28T22:51:22Z
---

# TH-0722 S-0416 is verified and meets all six criteria, but only you can accept it: it changes .claude/settings.json

On wip/kanban/stories/S-0416-flai-guard-reads-a-role-s-permitted-calls-from-its-role-file-and-refuses-the-rest-and-the-planner-s-the-orchestrator-s-and-the-analyzer-s-rules-move-into-their-files-with-their-guard-tests-unchanged.md.

## Entries

### 2026-10-28T22:51:22Z orchestrator
**Recommended:** accept S-0416 yourself, at 6c25aad7.

The merge queue refused it, because it changes paths Claude Code protects: `.claude/settings.json` and `template/root/.claude/settings.json` (`protected_paths`, ADR-0106). There is no other blocker.

## What I checked

- **`flai verify`** passed every step at 6c25aad7, the branch head. Integration and smoke are left to the queue, so they run only if you accept through the queue, or `flai accept` runs them.
- **The verifier** finds all six criteria met. All 38 changed files are inside the story's touches.
- **No refusal is looser** for the planner, the orchestrator, or the analyzer, and their existing guard tests are unchanged.
- **The settings.json change** widens the Edit/Write hook's condition to `[ -n "$FLAI_ROLE" ] || [ -n "$FLAI_STORY" ]`. That is criterion 5, which you allowed on TH-0712 and TH-0720.

## Worth knowing before you accept

1. **The planner gets stricter.** Its `item_edit` is now `ungated-edit`: it refuses `draft: false`, and it refuses setting or clearing a cost of delay input. The planner never sets inputs today, but confirm that is what you want.
2. **Criterion 4 is read a little narrowly.** A built-in role with no role file falls back to flai's embedded baseline copy instead of being refused. A defined role with no file is still refused. The design and the changelog say so.
3. **One stale file, harmless.** `flai/internal/roles/baseline/orchestrator.md` lacks item 11, which S-0502 added to the template's copy. The guard reads only the front matter, so no decision changes, and the test compares front matter only. A follow-up can copy the body over.

Accepting S-0416 frees the holds on S-0514 and on most of the backlog.
