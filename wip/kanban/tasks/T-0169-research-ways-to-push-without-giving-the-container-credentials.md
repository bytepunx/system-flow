---
id: T-0169
type: task
nature: research
title: "Research: ways to push without giving the container credentials"
status: backlog
parent: S-0052
owner: alex
created: 2026-09-19T01:55:25Z
updated: 2026-09-19T01:55:25Z
transitions: []
stream: S-0052
tags: []
---

# T-0169 Research: ways to push without giving the container credentials

## Work
Work out the ways a push could follow a board acceptance with the container still holding no credentials: a host-side `flai` command or watcher the operator starts that pushes when main is ahead with an acceptance commit; a git hook in the clone, and whether hooks run for commits made inside the container; the agent session pushing when `flai mcp` reports the change, as happens by hand today; and leaving the push manual while the board shows a standing "accepted, not pushed" state with the command and a way to check again. For each: what runs where, who has to have started it, what happens when nobody is there, and how a failed push is shown. Prototype only in a scratch repository with a scratch remote.

## Done when
- A section of the finding covers each approach under the same headings as T-0168 where they apply
- Each claim is marked tried or not tried

## Notes
