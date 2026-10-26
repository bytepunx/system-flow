---
id: T-2174
type: task
nature: feature
title: The bedrock image builds pilot beside flai and sets it as the entrypoint, the claude-code layer keeps it, and a container started with a request on stdin runs Claude Code through pilot
status: backlog
parent: S-0507
owner: alex
created: 2026-10-26T04:51:35Z
updated: 2026-10-26T04:51:35Z
transitions: []
stream: S-0507
tags: [cli, docker]
touches: [flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/agent-image.sh, scripts/agent-container.sh]
after: [T-2168]
---
# T-2174 The bedrock image builds pilot beside flai and sets it as the entrypoint, the claude-code layer keeps it, and a container started with a request on stdin runs Claude Code through pilot

## Work

In `flai/images/bedrock/Dockerfile`, build `pilot` in the Go stage beside `flai` (`go build ./pilot`), copy it to `/usr/local/bin/pilot`, and set `ENTRYPOINT ["tini", "--", "pilot"]`; keep `flai` on `PATH` for the guard and the CLI. Make sure `flai/images/claude-code/Dockerfile` inherits the entrypoint and does not reset it. Change `scripts/agent-container.sh` from S-0505 to feed a start request JSON on standard input in place of starting `claude` itself, so the spike's script becomes the by-hand way to run what `flai serve` will run, and keep `scripts/agent-image.sh` building both images. Run one container by hand on a story's worktree with a request naming `inner: claude-code` and check Claude Code reaches the HTTP MCP server through pilot's file. Waits for T-2168, which gives the binary.

## Done when

`make agent-images` builds both images with `pilot` as the entrypoint, `echo '{}' | docker run --rm -i flai-agent-bedrock:local` exits 2 with one `pilot:` line, and `scripts/agent-container.sh S-nnnn` runs Claude Code through pilot against the HTTP MCP server.

## Notes

The images stay local and unsigned, as T-2153's ADR records.
