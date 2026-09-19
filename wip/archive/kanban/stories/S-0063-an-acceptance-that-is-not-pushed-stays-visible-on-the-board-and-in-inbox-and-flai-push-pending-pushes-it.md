---
id: S-0063
type: story
nature: feature
title: An acceptance that is not pushed stays visible on the board and in inbox, and flai push --pending pushes it
status: done
parent: E-0006
owner: alex
created: 2026-09-19T07:54:42Z
updated: 2026-09-19T10:00:05Z
transitions:
  - to: ready
    at: 2026-09-19T08:13:32Z
    by: alex
  - to: in-progress
    at: 2026-09-19T08:34:07Z
    by: system-flow
  - to: review
    at: 2026-09-19T09:18:43Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:00:05Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal/workitem, flai/internal/mcpserver, flaiover]
---

# S-0063 An acceptance that is not pushed stays visible on the board and in inbox, and flai push --pending pushes it

## Goal
When an acceptance has been made and not pushed, that stays visible until it is no longer true: on the board and the item page for the operator, in `inbox` for agents, with one command on the host that pushes what is pending. Today the dashboard says it once, after the acceptance, and the fact is gone on the next page load.

## Acceptance criteria
- [x] flai can tell, offline and with no credential, whether `main` is ahead of its remote-tracking branch by commits that include an acceptance, and which tags point at commits not yet on the remote-tracking branch; `flai board --json` and the MCP `board` tool carry it as `unpushed: { commits, acceptances, tags }`, absent when there is nothing or no remote
- [x] The board and the item page of an accepted story show a standing "accepted, not pushed" notice with the counts, the tags, and the command to run, on every load until it stops being true; it clears by itself after a push from this clone, and the notice says that a push made from another clone is not seen until someone fetches here
- [x] `inbox` reports unpushed acceptances to agents on every call while they exist, like ready work and unlike one-time changes, and the server instructions and the work-management convention say what an agent does about it: fetch, push `main` and those tags, never force
- [x] `flai push --pending` on the host pushes `main` and the tags on unpushed commits with the host's own credentials, only when the commits ahead include an acceptance, says "nothing pending" otherwise, refuses when `main` has diverged from the remote and says to fetch and merge, and never forces; `--dry-run` prints what it would push
- [x] No retry button on the board: the operator chose the standing state without one
- [x] Tests for the detection (nothing, ahead without an acceptance, ahead with one, tags, no remote, diverged), the command against a scratch bare remote, the MCP field, and the notice; `docs/users/flai.md`, `docs/users/flaiover.md`, `docs/operators/index.md`, `flai-cli.md`, and `flaiover-dashboard.md` updated

## Tasks
- T-0225 Detect an unpushed acceptance offline, and carry it in flai board --json
- T-0226 flai push --pending on the host
- T-0227 The standing notice on the board and the item page
- T-0228 inbox and the MCP board tool report it, with what an agent does about it; documentation and criteria

## Notes
From S-0052's finding, `design/system/pushing-from-the-board.md`, and ADR-0026. The operator chose the standing state on 2026-09-19. It is worth building whether or not a project gives its container a key (S-0062): it is the default for every project that does not, and the fallback when a push fails.

Tried in the research: from a container with no credential, `git rev-list --count origin/main..main` and a loop over `git tag` with `git merge-base --is-ancestor <tag>^{commit} origin/main` gave the right answer offline; asking the remote fails without a credential; after a push from the same clone the container saw zero ahead at once, because it shares the clone. A shell prototype of the pending push worked against a scratch remote.

Acceptance commits are recognised today by their subject, `chore: [S-nnnn] accept and archive`. Between 2026-09-18 and 2026-09-19 an agent session pushed by hand after ten board acceptances in a row; this story turns that into something `inbox` tells an agent to do.
