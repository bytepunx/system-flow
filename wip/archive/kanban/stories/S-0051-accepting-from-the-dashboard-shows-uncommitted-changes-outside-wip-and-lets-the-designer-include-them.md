---
id: S-0051
type: story
nature: improvement
title: Accepting from the dashboard shows uncommitted changes outside wip and lets the designer include them
status: done
parent: E-0006
owner: alex
created: 2026-09-18T20:54:44Z
updated: 2026-09-18T22:54:14Z
transitions:
  - to: ready
    at: 2026-09-18T20:58:16Z
    by: alex
  - to: in-progress
    at: 2026-09-18T20:59:20Z
    by: alex
  - to: review
    at: 2026-09-18T21:09:55Z
    by: alex
  - to: done
    at: 2026-09-18T22:54:14Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover, flai/cmd, template]
---

# S-0051 Accepting from the dashboard shows uncommitted changes outside wip and lets the designer include them

## Goal
When a story is accepted from the dashboard and the working tree has uncommitted changes outside `wip/`, the designer sees which files they are before confirming and chooses there whether to include them in the acceptance commit or to cancel. Today the move is refused after the confirmation, with advice to pass `--yes`, which the dashboard has no way to do.

## Acceptance criteria
- [x] The acceptance preview reports uncommitted paths outside `wip/` instead of hiding them: `flai accept --dry-run` returns them in its result (for example `uncommitted: [...]`) whether or not `--yes` is given, and the dashboard's preview no longer differs from the move that follows it in what it checks
- [x] The confirmation lists those paths and offers an explicit choice, off by default: include them in the acceptance commit, or cancel and deal with them first. Accept stays disabled until the choice is made. With none, the confirmation is unchanged
- [x] Choosing to include them accepts the story with the same effect as `flai accept --yes` from a shell: the move endpoint passes the choice to flai, and nothing else gains a way to pass `--yes`
- [x] The container sees the ignore rules the host does: a file ignored only by the host user's global excludes file is not reported as uncommitted in the dashboard. How is decided in refinement (see notes), and documented
- [x] A project made from the template does not hit this for the common agent-local files: the template's `.gitignore`, and this repository's, cover `.claude/settings.local.json` and any other per-user agent file the conventions tell agents to create
- [x] Tests: the dry run result with and without uncommitted paths; the move endpoint passing the choice; the confirmation component with no paths, with paths and no choice, and with the choice made. `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, and `docs/users` updated

## Tasks
- T-0162 flai accept --dry-run reports uncommitted paths outside wip instead of refusing or hiding them
- T-0163 Dashboard: the preview lists uncommitted paths and the confirmation offers to include them
- T-0164 flai dashboard gives the container the host's global git excludes
- T-0165 Template and repository .gitignore cover agent-local files
- T-0166 Design and user documentation for the preview, the choice, and the excludes
- T-0167 Verify: all tiers and the choice through a real container

## Notes
Raised by the operator on 2026-09-18 accepting S-0048 from the board: "working tree has uncommitted changes outside wip (.claude/); commit or stash them so the acceptance commit holds only acceptance, or pass --yes to include them". The operator asked for a story to add the choice as a prompt in the dashboard. Recorded as I-0019.

What happened. `.claude/settings.local.json` was created by the agent that session to set `FLAI_AGENT`. The host ignores it through `~/.config/git/ignore`; the container has `HOME=/tmp` and no `core.excludesFile`, so git there reported `.claude/` as untracked, and `dirtyOutsideWip` in `flai/cmd/accept.go` refused the acceptance. Unblocked for this clone by adding the file to `.git/info/exclude`, which the container reads because it is inside the mounted repository.

Why the preview did not warn. `flaiover/src/routes/api/items/[id]/acceptance/+server.ts` runs `flai accept <id> --dry-run --yes`, and `--yes` skips the check; the move that follows (`.../move/+server.ts`, `flai move <id> done`) runs without it. So the first the designer hears of it is after confirming. The agent's own check of the preview on 2026-09-18 missed it for the same reason.

For the ignore rules, candidates: `flai dashboard` reads the host's `core.excludesFile` (or the default `~/.config/git/ignore`) and mounts it read-only in the container with `GIT_CONFIG_*` environment pointing `core.excludesFile` at it, the same way it passes the git identity; or flai copies its patterns into `.git/info/exclude`, which changes the clone and goes stale. The first is preferred. Either way the list in the confirmation is the backstop when the two still differ.

Including uncommitted files in an acceptance commit is a real choice with a cost: the acceptance commit stops holding only acceptance, which is why flai refuses by default. The prompt should say so in a line.

Not in the board `order`; the operator has not placed it against S-0040 to S-0043.

Verification, 2026-09-18. `make flai-test` passed (golangci-lint 0 issues; behavior, integration, smoke; markdown lint); flaiover lint, svelte-check 0 errors, 16 files and 110 tests, and a production build passed. A local image built from the branch was run as a second dashboard, on its own name and port, against a scratch git project with a story in review on a branch, started by the branch's `flai dashboard` with `XDG_CONFIG_HOME` pointing at a scratch global excludes file that lists `only-global.txt`; the operator's git configuration and dashboard were not touched. The container had the read-only mount at `/run/flaiover/gitignore` and the three `GIT_CONFIG_*` variables. With `only-global.txt` and `docs/stray.md` created at the project root: git in the container reported only `docs/stray.md`; the preview returned `uncommitted: ["docs/stray.md"]` with no blockers; `POST .../move` to done without the choice was refused with flai's message and the story stayed in review; with `include_uncommitted: true` it merged, moved to done, and archived, the acceptance commit held `docs/stray.md`, and `only-global.txt` stayed untracked. The confirmation itself was exercised by its component tests, not clicked in a browser. The refusal in the third step comes back as HTTP 500, not 400, as refusals that are not workflow rules already did; the confirmation no longer leads there, so it was left alone.
