---
id: S-0119
type: story
nature: feature
title: Build a version of the flaiover dashboard for arm64
status: backlog
owner: alex
created: 2026-09-26T03:14:08Z
updated: 2026-09-26T03:14:08Z
transitions: []
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0119 Build a version of the flaiover dashboard for arm64

## Goal

CI does not presently build an arm64 compatible docker image. This forces people to build their own image locally or do without and breaks the dashboard --pull command.

## Acceptance criteria
- [ ] CI builds include an arm64 image build
- [ ] flai dashboard --pull fetches the correct image from ghcr.io based on the platform

## Tasks

## Notes
