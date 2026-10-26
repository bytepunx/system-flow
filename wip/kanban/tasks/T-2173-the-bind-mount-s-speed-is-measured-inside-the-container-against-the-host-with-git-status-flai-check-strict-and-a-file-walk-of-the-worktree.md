---
id: T-2173
type: task
nature: research
title: The bind mount's speed is measured inside the container against the host with git status, flai check --strict, and a file walk of the worktree
status: backlog
parent: S-0505
owner: alex
created: 2026-10-26T04:51:25Z
updated: 2026-10-26T04:51:25Z
transitions: []
stream: S-0505
tags: [cli, docker]
touches: [design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md]
after: [T-2166]
---
# T-2173 The bind mount's speed is measured inside the container against the host with git status, flai check --strict, and a file walk of the worktree

## Work

In the story's worktree, time three commands on the host and inside the container on the same bind mount, three runs each, warm: `git status`, `flai check --strict`, and a file walk (`find . -type f | wc -l`). Run the container with `docker run --rm` on `flai-agent-claude-code:local` with the worktree mounted as `scripts/agent-container.sh` mounts it. Note the host filesystem the worktree is on (the Linux filesystem of this WSL2 distribution, which ADR-0203 expects to be fast, or a Windows path, which it expects to be slow) and Docker's storage driver. Record the medians in a table in the results document's "Results" with the ratio inside over host, and say against ADR-0203's limitation whether the bind mount at the host path is fast enough for a Go test suite or whether the clone path is needed from the start. Waits for T-2166, which writes the same results document.

## Done when

The results document holds the table of the three commands, host and container, with the ratios and the filesystem named, and one sentence on whether the bind mount is fast enough.

## Notes

If Go is on the host, `go test ./internal/workitem/...` run on the host alone gives the suite's own time to compare the ratio against.
