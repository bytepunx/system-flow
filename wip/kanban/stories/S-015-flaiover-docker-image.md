---
id: S-015
type: story
nature: feature
title: Docker image build and publish
status: review
parent: E-003
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T05:33:22Z
transitions:
  - to: ready
    at: 2026-09-17T05:29:33Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:29:33Z
    by: agent
  - to: review
    at: 2026-09-17T05:33:22Z
    by: agent
tags: []
---

# S-015 Docker image build and publish

## Goal
Multi-stage Dockerfile and a workflow that publishes to ghcr.io/bytepunx/flaiover.

## Acceptance criteria
- [x] Image bundles a flai binary built from the same commit at /usr/local/bin/flai (ADR-0016)
- [x] Image runs as an arbitrary non-root UID (flai dashboard passes `--user uid:gid`; verified in S-009 that an image needing root or a fixed user exits at once) and serves on 3000
- [x] Tags: latest, semver, major, sha
- [ ] flai dashboard runs the published image end to end (verified with the local image in S-015; the published one is checked after the first tag push, at acceptance)

## Tasks
- T-098 Multi-stage Dockerfile: flai from the same commit, flaiover build, node:24-alpine runtime as any UID on 3000
- T-099 Local image script and Makefile target; run it with flai dashboard end to end
- T-100 release-flaiover workflow: GHCR tags latest, semver, major, sha; docs

## Notes
- Verified end to end with the local image; the published image is checked at acceptance. Image size 667 MB, slimming is a follow-up.
