---
id: T-0854
type: task
nature: improvement
title: git.md says how to show a test fails without the change and never to use git stash in a worktree, and delegation.md says to batch independent calls
status: backlog
parent: S-0268
owner: alex
created: 2026-10-05T00:29:27Z
updated: 2026-10-05T00:29:27Z
transitions: []
stream: S-0268
tags: [conventions, template]
touches: [design/conventions/git.md, template/root/design/conventions/git.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md]
---
# T-0854 git.md says how to show a test fails without the change and never to use git stash in a worktree, and delegation.md says to batch independent calls

## Work

Change the baseline of `design/conventions/git.md`, above its marker, and make the same change in `template/root/design/conventions/git.md`. Add a rule on showing that a new test fails without the change under test:

- Check the files the change touches out from the main branch into a scratch copy outside the worktree, such as with `git show main:<path>`, and run the test against that copy. Or build the old binary from the main branch and run the test against it.
- Never use `git stash` in a worktree. Every worktree of a clone shares one stash stack, so a stash pushed or popped in one session can take another session's work. S-0248's agent pushed and dropped a stash to prove its test failed.

Change the baseline of `design/conventions/delegation.md`, under "While working", and its template copy the same way. Say to make independent edits and commands in one turn, and to move a task just written to ready and in-progress in one command. These are the sentences that T-0853 puts in the prompt, so that an agent started without flai serve's prompt reads them too.

This task waits for none. It shares no path with T-0853, so the two can run together. Leave `template/CHANGELOG.md` and `template/template.yaml` to the release tooling.

## Done when

- Both copies of `git.md` carry the rule on showing a test fails, and the ban on `git stash` with its reason, and they read the same above the marker.
- Both copies of `delegation.md` carry the batching sentence and the one-command move, and they read the same above the marker.
- The markdown lint passes.

## Notes

Drafted by the planner for S-0268.
