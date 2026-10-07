---
id: I-0100
title: The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai
class: defect
status: open
count: 1
cost: 2m
first_reported: 2026-10-06T23:31:35Z
last_reported: 2026-10-06T23:31:35Z
updated: 2026-10-07T01:07:13Z
---

# I-0100 The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai

## Description

`flai accept` runs its steps in order: rebase and merge, move to done, archive, commit, then the overlap notices. The commit is a `git commit` in the main checkout, where flai serve's replans, planners' autocommits, and other acceptances also commit. When one of them holds `.git/index.lock` at that moment, the commit fails, and the acceptance exits with git's error. Everything before the commit has happened and stays: the branch is merged, the story is done, the items are moved to the archive and staged. Run again, the acceptance refuses ("already done"), and no other command finishes it. The operator, or someone watching, has to commit by hand.

The first lines of git's error are all the record there is, on the terminal of whoever ran the acceptance; flai serve's journal records the run as done because the move to done succeeded.

## Instances

### 2026-10-06T23:31:35Z
Story: S-0299.
flai accept S-0299 at 23:03:32Z on 2026-10-06 merged the branch, moved the story to done, and archived eight items, then its git commit failed and the command exited, leaving the archive moves staged. Run again it said 'S-0299 is already done' and changed nothing. Claude, watching the board, checked the staged set and committed it as the acceptance would have, at 23:04:31Z. Commits land in the main checkout from flai serve's replans and from planners' and stories' autocommits; a planner was running at the time. The error's text beyond its first line was not kept.

## Remediation

Directions to weigh: retry the commit a few times when git reports the index lock held, since the other writer is brief; or make an acceptance that is done and archived but whose commit is missing resumable, as one that is done but not archived already is ("completed from step 0 without a second transition"); and record the commit's failure in the journal.

Story S-0307 remediates this issue, created from it at 2026-10-07T01:07:13Z.
