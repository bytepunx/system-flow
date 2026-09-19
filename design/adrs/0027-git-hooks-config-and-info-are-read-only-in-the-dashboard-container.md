---
id: ADR-0027
title: "Git hooks, config, and info are read-only in the dashboard container"
status: accepted
date: 2026-09-19
supersedes: []
superseded_by: []
refines: [ADR-0018, ADR-0022, ADR-0026]
---

# ADR-0027 Git hooks, config, and info are read-only in the dashboard container

## Context

`flai dashboard` mounts the clone into the container read-write, `.git` included, because an acceptance made from the board merges a story branch, commits, tags, and removes a worktree ([ADR-0016](0016-dashboard-delegates-to-flai.md), [ADR-0022](0022-repository-mounted-at-its-host-path.md)). S-0052 found what that also allows. A process in the container can write `.git/hooks/pre-push`, and the next `git push` the operator or an agent runs on the host executes it as the operator; tried in a scratch project, it ran (I-0022). The same holds for git settings that make git run a command or talk to a different remote, and for `.git/info/exclude`, which can hide a planted file from `git status`. None of this shows in `git status` or a diff.

It needs a compromised container, a flaw in flaiover or one of its dependencies, not just the dashboard token: the API never writes to `.git` except through git. But [ADR-0018](0018-dashboard-token.md) publishes the dashboard on every interface with one token as its lock, the operator reaches it over a public tunnel, and since [ADR-0026](0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md) the container may hold a push key. For a project that gives the container nothing, this route was the only way from the container to the operator's credentials.

S-0064 measured, with the published image run the way `flai dashboard` runs it, what the container's legitimate work writes under `.git`: for an acceptance, `ORIG_HEAD`, `COMMIT_EDITMSG`, `index`, the logs, the refs (the main branch moved, the story branch and its log removed, the tag created), `worktrees/S-nnnn` removed, objects, and `config`, which `git branch -d` rewrites with identical content and tolerates failing to rewrite. A document save, an item creation, and a board move write nothing under `.git` but a commit. Nothing the container does needs `.git/hooks`, `.git/info`, or a changed `.git/config`.

## Decision

`flai dashboard` mounts these paths a second time, read-only, over the read-write clone: `.git/hooks`, `.git/info`, `.git/config`, the dashboard token's file in `.flai-cache`, and the host's flai config when it is kept inside the clone. `.git/hooks` and `.git/info` are created on the host when missing, so the container cannot create them in their place. When the clone has no `.git` directory (not a repository, or a linked worktree) nothing extra is mounted and `flai dashboard` says so.

`flai dashboard` also names, at start and in `status`, what the clone already holds of the kinds the mounts now prevent: an enabled hook, and settings such as `core.hooksPath`, `core.fsmonitor`, `core.sshCommand`, an editor or pager, a credential helper, a shell alias, an include, a filter, diff, or merge driver, `url.*.insteadOf`, and a `pushurl`. They may be the operator's own; flai cannot tell, so it lists them. `status` says when a running container was started without the read-only paths.

## Consequences

- Tried with all three git paths read-only: writing a hook, `git config` for `core.sshCommand`, an alias, `core.hooksPath`, and `--worktree`, appending to and renaming the config, replacing the hooks directory, and writing `info/exclude` and `info/attributes` all fail. A hook written into `.git/worktrees/<name>/hooks` can be written and does not run: git uses the common directory's hooks. An acceptance with a merge, a branch deletion, a worktree removal, and a tag, a document save with its commit, an item creation, and a push with a key still work.
- A read-only bind mount of a file pins the file the host had at start. When the operator changes `.git/config` on the host, git replaces the file, and the container keeps reading the old one until the dashboard is restarted. The same was already true of the token. The operators' documentation says to restart after changing git settings.
- What stays open, and why. The container must write the work tree and the refs, so it can change a tracked file, including a script the operator runs, or move a branch. Those show in `git status`, `git diff`, and `git log`, acceptance refuses uncommitted changes outside `wip/` unless asked to include them (S-0051), and review is where they are caught; a mount cannot tell a legitimate merge from a malicious edit. Files that are ignored by git and executed on the host (a built binary, `node_modules/.bin`, tool caches under `.flai-cache`) are project specific and invisible to `git status`; flai does not know which a project has. This decision closes the routes git itself offers, not every route a writable work tree offers.
- The warning lives in `flai dashboard`, not in `flai check`: the check is the repository's gate and runs in CI, where a clone has no hooks, and a developer's own pre-commit hook is legitimate and would fail `--strict`.
- I-0022 is closed by this for the routes it names.

## Alternatives considered

- Mounting `.git` read-only: acceptance could not commit, tag, or remove a worktree. It would mean no acceptance from the board, which is what E-0006 exists to provide.
- Running the container's git with `core.hooksPath=/dev/null` and similar: that protects the container's git from the clone, which is the wrong direction; the threat is what the host's git later reads.
- Checking hooks and config on the host before every host-side git command: flai does not run the operator's git, and an agent's git is not flai's either.
- A separate clone for the dashboard that the operator pulls from: it removes the shared `.git` altogether and is the stronger boundary, at the cost of the live board ([ADR-0019](0019-story-branches-in-worktrees.md) keeps `wip/` in the main checkout so that the board is live). Worth revisiting if the dashboard is ever offered to people the operator does not trust.
