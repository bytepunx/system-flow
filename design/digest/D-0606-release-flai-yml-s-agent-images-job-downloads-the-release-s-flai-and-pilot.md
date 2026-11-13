---
id: D-0606
title: release-flai.yml's agent-images job downloads the release's flai and pilot archives with gh and checks them against the signed checksums.txt with openssl, rather than taking GoReleaser's dist folder by an artifact action, which would need a new action pin
kind: decision
sources: [S-0734]
topics: [release, ci, docker, security]
paths: [".github/workflows/release-flai.yml", scripts/ci-pins.sh, design/tech/ci.md]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0606 release-flai.yml's agent-images job downloads the release's flai and pilot archives with gh and checks them against the signed checksums.txt with openssl, rather than taking GoReleaser's dist folder by an artifact action, which would need a new action pin

S-0734 made each `flai/v*` tag build and publish the agent images from the release's own Linux binaries. The `agent-images` job runs after the `goreleaser` job, so it needs that job's archives. S-0734 chose to fetch them from the published release rather than pass the dist folder between jobs. From its Decisions (`wip/archive/agents/S-0734.md`):

```text
The release jobs use only actions already pinned in the repository and download the archives from the published release with gh, checked against checksums.txt and its signature with openssl, rather than handing GoReleaser's dist folder over by an artifact action that would need a new pin.
```

`design/tech/ci.md` describes the job: it downloads with `gh` and checks the signature with `openssl`. It does not give the reason. Every action the workflows use is pinned, and `scripts/ci-pins.sh` checks the pins. So a later story that changes `release-flai.yml` keeps the download from the release. It adds an artifact action only if it also adds that action's pin.

Source: S-0734.
