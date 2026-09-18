---
id: S-0051
type: story
nature: improvement
title: Accepting from the dashboard shows uncommitted changes outside wip and lets the designer include them
status: backlog
parent: E-0006
owner: alex
created: 2026-09-18T20:54:44Z
updated: 2026-09-18T20:54:44Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover, flai/cmd, template]
---

# S-0051 Accepting from the dashboard shows uncommitted changes outside wip and lets the designer include them

## Goal
When a story is accepted from the dashboard and the working tree has uncommitted changes outside `wip/`, the designer sees which files they are before confirming and chooses there whether to include them in the acceptance commit or to cancel. Today the move is refused after the confirmation, with advice to pass `--yes`, which the dashboard has no way to do.

## Acceptance criteria
- [ ] The acceptance preview reports uncommitted paths outside `wip/` instead of hiding them: `flai accept --dry-run` returns them in its result (for example `uncommitted: [...]`) whether or not `--yes` is given, and the dashboard's preview no longer differs from the move that follows it in what it checks
- [ ] The confirmation lists those paths and offers an explicit choice, off by default: include them in the acceptance commit, or cancel and deal with them first. Accept stays disabled until the choice is made. With none, the confirmation is unchanged
- [ ] Choosing to include them accepts the story with the same effect as `flai accept --yes` from a shell: the move endpoint passes the choice to flai, and nothing else gains a way to pass `--yes`
- [ ] The container sees the ignore rules the host does: a file ignored only by the host user's global excludes file is not reported as uncommitted in the dashboard. How is decided in refinement (see notes), and documented
- [ ] A project made from the template does not hit this for the common agent-local files: the template's `.gitignore`, and this repository's, cover `.claude/settings.local.json` and any other per-user agent file the conventions tell agents to create
- [ ] Tests: the dry run result with and without uncommitted paths; the move endpoint passing the choice; the confirmation component with no paths, with paths and no choice, and with the choice made. `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, and `docs/users` updated

## Tasks

## Notes
Raised by the operator on 2026-09-18 accepting S-0048 from the board: "working tree has uncommitted changes outside wip (.claude/); commit or stash them so the acceptance commit holds only acceptance, or pass --yes to include them". The operator asked for a story to add the choice as a prompt in the dashboard. Recorded as I-0019.

What happened. `.claude/settings.local.json` was created by the agent that session to set `FLAI_AGENT`. The host ignores it through `~/.config/git/ignore`; the container has `HOME=/tmp` and no `core.excludesFile`, so git there reported `.claude/` as untracked, and `dirtyOutsideWip` in `flai/cmd/accept.go` refused the acceptance. Unblocked for this clone by adding the file to `.git/info/exclude`, which the container reads because it is inside the mounted repository.

Why the preview did not warn. `flaiover/src/routes/api/items/[id]/acceptance/+server.ts` runs `flai accept <id> --dry-run --yes`, and `--yes` skips the check; the move that follows (`.../move/+server.ts`, `flai move <id> done`) runs without it. So the first the designer hears of it is after confirming. The agent's own check of the preview on 2026-09-18 missed it for the same reason.

For the ignore rules, candidates: `flai dashboard` reads the host's `core.excludesFile` (or the default `~/.config/git/ignore`) and mounts it read-only in the container with `GIT_CONFIG_*` environment pointing `core.excludesFile` at it, the same way it passes the git identity; or flai copies its patterns into `.git/info/exclude`, which changes the clone and goes stale. The first is preferred. Either way the list in the confirmation is the backstop when the two still differ.

Including uncommitted files in an acceptance commit is a real choice with a cost: the acceptance commit stops holding only acceptance, which is why flai refuses by default. The prompt should say so in a line.

Not in the board `order`; the operator has not placed it against S-0040 to S-0043.
