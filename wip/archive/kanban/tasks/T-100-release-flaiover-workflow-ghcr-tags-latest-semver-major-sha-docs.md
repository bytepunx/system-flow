---
id: T-100
type: task
nature: feature
title: "release-flaiover workflow: GHCR tags latest, semver, major, sha; docs"
status: done
parent: S-015
owner: alex
created: 2026-09-17T05:29:33Z
updated: 2026-09-17T05:33:21Z
transitions:
  - to: ready
    at: 2026-09-17T05:33:21Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:33:21Z
    by: agent
  - to: done
    at: 2026-09-17T05:33:21Z
    by: agent
stream: S-015
tags: [dashboard, ci, docs]
---

# T-100 release-flaiover workflow: GHCR tags latest, semver, major, sha; docs

## Work
release-flaiover.yml on tags flaiover/v* and pushes to main: login to GHCR with GITHUB_TOKEN, docker/metadata-action with latest on main, semver and major from the tag, sha; build-push-action with the Dockerfile; docs/operators (private registry login, image tags) and design/tech/docker.md.

## Done when
Workflow valid; docs match; S-009's published-image criterion becomes verifiable after the first tag.

## Notes
