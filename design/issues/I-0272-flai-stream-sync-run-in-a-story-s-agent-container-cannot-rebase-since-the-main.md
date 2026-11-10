---
id: I-0272
title: flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached
class: defect
status: open
count: 2
cost: 6m
first_reported: 2026-11-07T04:41:52Z
last_reported: 2026-11-07T06:52:30Z
updated: 2026-11-10T08:49:00Z
---

# I-0272 flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached

## Description
flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached

## Instances

### 2026-11-07T04:41:52Z
Story: S-0558.
S-0558's agent ran on harness container, whose mounts (ADR-0246) keep the main checkout's .git read-only. Its start prompt says to run flai stream sync before review. Run through Bash in the container, flai stream sync S-0558 exited 1 with "the rebase onto main stopped" and "resolve them and git rebase --continue", though no path conflicted: the first pick could not write its commit object and was rescheduled, with HEAD on main's commit and a clean tree. The message's other way out, git rebase --abort, exited 0 yet printed "cannot lock ref refs/heads/story/S-0558 ... Read-only file system", and left HEAD detached on main's commit with the branch's changes staged. The branch itself stayed at its head, so nothing was lost: the agent checked that the index and files matched the branch's tip, then ran git symbolic-ref HEAD refs/heads/story/S-0558 in its own registration. The sync that works from the container is the host's, as task_done runs it, and no MCP tool syncs alone. The same mounts refuse flai move, whose items lock is under the main checkout (the MCP item_move works), and hand edits of the narrative's Decisions, which the prompt and work-management.md ask for and which S-0558 recorded in the story's Notes instead. The sync should say what stopped it and, in a container, refuse at once or go through the host; the start prompt on this harness should name what works there.

### 2026-11-07T06:52:30Z
Story: S-0556.
S-0556's agent ran flai stream sync S-0556 through Bash in its container with main not ahead, so no rebase was needed. It exited 0 and said the branch is rebased, but warned "story branch checks failed after the rebase": the trial merge with story/S-0565, git merge-tree --write-tree, exited 128 with "unable to create temporary file: Read-only file system". So in the container the sync can report success with its trial merge against other open stories not run. The MCP tool task_done, whose sync runs on the host, reported that same trial merge clean each time in this session.

## Remediation

Story S-0689 remediates this issue, created from it at 2026-11-10T08:49:00Z.
