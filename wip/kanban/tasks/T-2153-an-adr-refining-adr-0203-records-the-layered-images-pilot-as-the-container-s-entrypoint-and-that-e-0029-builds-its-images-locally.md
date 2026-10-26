---
id: T-2153
type: task
nature: feature
title: An ADR refining ADR-0203 records the layered images, pilot as the container's entrypoint, and that E-0029 builds its images locally
status: backlog
parent: S-0507
owner: alex
created: 2026-10-26T04:47:38Z
updated: 2026-10-26T04:47:38Z
transitions: []
stream: S-0507
tags: [cli, docker]
touches: [design/adrs/README.md, design/tech/docker.md]
---
# T-2153 An ADR refining ADR-0203 records the layered images, pilot as the container's entrypoint, and that E-0029 builds its images locally

## Work

Write an ADR with `flai adr new`, refining ADR-0203, that records the two refinements E-0029's notes make and the facts S-0505 settled: the images are layered, a bedrock layer (a small Linux base, `git`, `flai`, `pilot`) and one layer per inner harness, `claude-code` first, each built from `flai/images/<layer>/Dockerfile`; `pilot`, a second Go binary of the flai module, is the container's entrypoint in place of `flai agent run`, reads one start request on standard input, starts the inner harness through flai's adapters, and relays its log unchanged so flai's readers stay per harness (ADR-0130); and the images are built locally by `scripts/agent-image.sh` within E-0029, with the signed, published image and its release pipeline left to a later story. Give its consequences, among them that a project's layer (S-0510) builds over the harness layer. Update `design/tech/docker.md` to list the layered images beside `flaiover`'s. This task waits for nothing.

## Done when

The ADR is accepted and indexed in `design/adrs/README.md`, `design/tech/docker.md` lists the bedrock and claude-code images with their Dockerfiles and the `pilot` entrypoint, and `flai check --strict` passes.

## Notes

The ADR's file is named by `flai adr new`; the story keeps the `design/adrs` folder touch for it.
