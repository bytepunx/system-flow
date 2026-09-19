---
id: S-0064
type: story
nature: remediation
title: A dashboard container cannot leave git hooks or configuration that run on the host
status: backlog
parent: E-0006
owner: alex
created: 2026-09-19T07:54:42Z
updated: 2026-09-19T07:55:41Z
transitions: []
tags: []
touches: [flai/cmd, docs/operators]
---

# S-0064 A dashboard container cannot leave git hooks or configuration that run on the host

## Goal
A process inside the dashboard container cannot leave anything in the mounted clone that git later executes on the host as the operator. Today it can: the clone is mounted read-write with `.git` included, and a hook written from the container runs at the operator's next git command.

## Acceptance criteria
- [ ] The routes are enumerated first, in the story's notes or the living design: hooks in `.git/hooks` and in each worktree's git directory, `core.hooksPath`, `core.fsmonitor`, `core.sshCommand`, `core.editor` and pager settings, credential helpers, aliases, `include` paths, filter and diff drivers set through `.git/config` or `.git/info/attributes`, and anything else found; each marked as closed by this story, already closed, or left open with the reason
- [ ] `flai dashboard` closes them without breaking what the container legitimately does (commit, merge a story branch, remove a worktree, tag, and with S-0062 push): for example hooks and config mounted read-only over the read-write clone, with whatever acceptance needs to write kept writable. An acceptance with a story branch and worktree, a document save, and a board move all still work in the container, tried end to end
- [ ] Tried in a scratch project, never this repository: a hook and a config setting written from inside the container either cannot be written or do not run on the host, where before the change the hook ran as the host user at the next push
- [ ] `flai check` or `flai dashboard` warns when the clone already holds a hook or a config setting of the kinds above that the template did not put there, so something planted before the fix is seen
- [ ] ADR-0018's consequences are not edited (accepted ADRs are immutable); the operators' documentation and `flaiover-dashboard.md` say what the container can and cannot write, and an ADR records the mount layout if it changes what `flai dashboard` gives the container

## Tasks

## Notes
Found by S-0052 on 2026-09-19 while weighing credentials for the container, and recorded in `design/system/pushing-from-the-board.md` under "What is at stake". Tried in a scratch project with a scratch remote: from inside a container of the published image, run as `flai dashboard` runs it, a `pre-push` hook was written to `.git/hooks`; the next `git push` on the host ran it as the host user. No credential in the container was involved.

It needs a compromised container, not just the dashboard token: the API never writes to `.git` except through git. It matters most for the default set-up, where the container holds no credential and this is the only route to the operator's. With S-0062's opt-in the container holds a key anyway, and this still stops the compromise spreading from the key to a shell on the host.

Not researched: which paths under `.git` acceptance has to write (objects, refs, logs, index, `worktrees/`, `config` when a worktree is added or removed with `extensions.worktreeConfig` or relative paths, ADR-0022). That decides whether a read-only `config` is possible or whether the answer is `safe` git options on the host side instead. Start there.
