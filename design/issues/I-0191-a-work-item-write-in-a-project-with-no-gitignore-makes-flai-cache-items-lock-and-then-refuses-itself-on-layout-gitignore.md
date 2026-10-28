---
id: I-0191
title: A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-10-22T20:40:04Z
last_reported: 2026-10-22T20:40:04Z
updated: 2026-10-28T22:58:42Z
---

# I-0191 A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore

## Description
A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore

## Instances

### 2026-10-22T20:40:04Z
Story: S-0484.
S-0479's items lock (ADR-0184) creates .flai-cache/items.lock on every work item write. In a project with no .gitignore, such as the fixture flai/internal/metrics/testdata/good copied to a temp folder, the write's scoped check then finds layout.gitignore (.flai-cache/ exists but there is no .gitignore) and refuses the write: `flai edit S-004 --title ...` on a fresh copy fails with "flai check has 1 finding(s) with this change". Two tests of flaiover/src/lib/server/writes.test.ts failed on main and stopped S-0484's close-out at the flaiover tier; S-0479's own close-out did not select that tier. S-0484 gave the two temp projects a .gitignore (T-2029). Remedy to choose: keep the lock outside the project, leave layout.gitignore out of a write's scoped check, or have flai write the .gitignore entry when it makes .flai-cache.

## Remediation

Story S-0525 remediates this issue, created from it at 2026-10-28T12:53:57Z.
Closed 2026-10-28T22:58:42Z: Fixed by S-0525 (ADR-0213): a checked creation (itemnew.Create, behind story/task/epic new with --body-stdin, --after, or --autocommit, item_new with a body or after, and flai issue story) scopes its second check run to the new item's story through check.ScopeToItem, as ADR-0197 scopes a checked edit, so the layout.gitignore that making .flai-cache/items.lock brings in a project with no .gitignore is a warning outside the story and the item stands. TestTheFirstCreationInAProjectWithNoGitignoreIsKeptWithAWarning (itemnew) reproduces it and fails on the code before the fix.
