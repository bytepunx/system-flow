---
id: S-0052
type: story
nature: research
title: Work out whether and how a board move to done can push the acceptance
status: backlog
parent: E-0006
owner: alex
created: 2026-09-19T01:55:25Z
updated: 2026-09-19T07:23:18Z
transitions: []
blocked:
  - from: 2026-09-19T01:57:55Z
    until: 2026-09-19T07:23:18Z
    reason: research stories cannot presently be moved to done
tags: [dashboard, cli]
touches: [design/adrs, design/system]
---

# S-0052 Work out whether and how a board move to done can push the acceptance

## Goal
Decide, with evidence, whether accepting a story from the board should also push the acceptance commit and its release tags, and if so how. Today the container has no git credentials by design: an acceptance from the board is committed and tagged locally, and the dashboard prints the `git push` command for someone to run on the host. The deliverable is a written finding with a recommendation, the operator's decision recorded, and the follow-up work queued. "Do not automate it" is an acceptable outcome. No push mechanism is built in this story.

## Acceptance criteria
- [ ] The finding describes each way the push could happen, from inside the container and from outside it, and for each: what has to be configured, what someone who holds the dashboard token could then do to the remote, how the access is limited and revoked, which hosts it works on (Linux, WSL2, macOS and Windows with Docker Desktop), and how a failed push is shown and retried
- [ ] Anything the finding says works was tried in a scratch repository against a scratch remote, never against this repository's remote, and says so; anything not tried is marked as not tried
- [ ] The finding includes the option of not pushing from the board and what would make that less easy to miss, such as a standing "accepted, not pushed" state on the board with the command to run
- [ ] The operator's answers on exposure and unattended pushes are recorded in the finding, from a conversation held through a thread on this story or in the session
- [ ] The finding ends with one recommendation and the reasons the others lost, in a document under `design/system` or as a proposed ADR
- [ ] The operator's decision is recorded: an ADR when it changes what the container is given or the security posture in ADR-0018 and the operator documentation, a note in the living design otherwise
- [ ] The work that follows from the decision is queued as one or more stories under E-0006, or S-0041 is amended, with goal and acceptance criteria; nothing about pushing is implemented here

## Tasks
- T-0168 Research: what pushing from the container would need, and what each way exposes
- T-0169 Research: ways to push without giving the container credentials
- T-0170 Conversation with the operator: who can reach the dashboard, and what an unattended push may do
- T-0171 Write the finding with a recommendation, get the decision, record it, and queue what follows

## Notes
Raised by the operator on 2026-09-18 after accepting S-0051 from the board: the dashboard said "done. run: git push origin HEAD flaiover/v0.11.4 flai/v1.2.4" and the change had not reached the remote, which was a surprise. The same had happened with S-0048. Both times the agent pushed `main` and the tags from the host afterwards.

Why it is like this. `flai dashboard` passes the host's git identity and, since S-0051, its global excludes into the container, and never a key or a token; `docs/operators/index.md` says an acceptance from the dashboard is committed and tagged locally and pushed from a shell. The dashboard is published on every interface of the host by default and protected by one project token (ADR-0018), and the hub work in S-0043 is about reaching it over a tunnel. A container that can push turns that token into write access to the remote, including release tags that start the release workflows and publish images. That is the trade to weigh, not an argument already settled.

Starting points for the research, none of them checked yet. Inside the container: forwarding the host's SSH agent socket; a deploy key mounted read-only; an HTTPS token (a fine-grained personal access token or `gh auth token`) handed to a git credential helper. Outside the container: something on the host that pushes when an acceptance appears, such as a `flai` command or watcher the operator runs, a git hook on the host, or the agent session doing it as it does now through `flai mcp` events; and the board merely showing that a push is pending.

The operator asked for this story to carry research and conversational tasks, so four tasks are written now, while the story is backlog. That departs from the convention that the pulling agent writes the tasks (ADR-0021); it is followed as an operator instruction, and the agent that pulls the story may add to them.

Two things to know before this story is accepted. A `research` story cannot be accepted from the board today: `release.LevelFor` refuses the nature ("research and experiment stay on a branch and do not release from main"), which fails the preview and the move; `flai accept S-0052 --no-release` from a shell works. And `git.md` says research stays on a branch, while this story's finding and ADR are meant to land on main; the agent that pulls it should raise that in its open questions rather than work around it.

Not in the board `order`; the operator has not placed it against S-0040 to S-0043.
