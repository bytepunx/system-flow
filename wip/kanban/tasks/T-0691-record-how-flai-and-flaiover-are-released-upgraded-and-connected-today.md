---
id: T-0691
type: task
nature: research
title: Record how flai and flaiover are released, upgraded, and connected today
status: done
parent: S-0193
owner: arobson
created: 2026-10-02T12:15:33Z
updated: 2026-10-02T12:19:01Z
transitions:
  - to: ready
    at: 2026-10-02T12:16:22Z
    by: claude-fable-5-1
  - to: in-progress
    at: 2026-10-02T12:16:22Z
    by: claude-fable-5-1
  - to: done
    at: 2026-10-02T12:19:01Z
    by: claude-fable-5-1
stream: S-0193
tags: []
touches: [design/system/release-signing.md, design/system/README.md]
usage:
  source: log
  seconds: 159
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 256
      output: 103
      cache_read: 892184
      cache_write: 24522
      cost: 0
---
# T-0691 Record how flai and flaiover are released, upgraded, and connected today

## Work

Start `design/system/release-signing.md` as the finding of S-0193, in the shape of `dashboard-host-channel.md`: a `## Today` section that states, from the code and the workflows, how `flai` is released (GoReleaser, `checksums.txt`, no signature), how `flaiover` is built and pushed (buildx to GHCR, no attestation), how `flai self-upgrade` and `install.sh` verify what they download (SHA-256 against an unsigned `checksums.txt`), how `flai dashboard` chooses the image it runs (`dashboard.image` and `dashboard.tag`, `--build` for a local image), and how the two authenticate each other on the channel (HMAC proof of the agent credential, versions exchanged in `hello`, close codes). Add the row to `design/system/README.md`. Nothing in this task depends on a decision, so it waits for no task.

## Done when

- `design/system/release-signing.md` exists with front matter and a `## Today` section whose every statement names the file or workflow it comes from.
- `design/system/README.md` lists it.

## Notes
