---
id: T-3168
type: task
nature: feature
title: "The flai binary carries the bedrock and claude-code Dockerfiles through a go:embed package at flai/images"
status: backlog
parent: S-0740
owner: alex
created: 2026-11-12T15:59:20Z
updated: 2026-11-12T15:59:20Z
transitions: []
stream: S-0740
tags: [cli]
touches: [flai/images/embed.go]
---
# T-3168 The flai binary carries the bedrock and claude-code Dockerfiles through a go:embed package at flai/images

## Work

Add `flai/images/embed.go`, package `images`, which embeds `bedrock/Dockerfile` and `claude-code/Dockerfile` with `//go:embed` and gives them by inner harness. `go:embed` cannot reach a parent folder, so the package sits beside the Dockerfiles. Leave `system-flow/Dockerfile` out: it is this repository's layer, not flai's.

Give each embedded file's sha256, which the label of S-0733's ADR may record as the Dockerfile's version for the staleness rule.

`flai/images` is a protected folder for a container's agent (`flai/internal/protected/protected.go`, lines 75 to 78), so the story's agent's write there asks the owner. It waits for no task.

## Done when

- `go build ./...` passes, and a test in the package reads each embedded Dockerfile and its sum.
- `go list -f '{{.EmbedFiles}}' ./images` names exactly the two Dockerfiles.

## Notes
