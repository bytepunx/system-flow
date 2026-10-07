---
id: T-1231
type: task
nature: research
title: Which digest RepoDigests records for a flaiover image pulled by index and by platform digest is checked, and the digest list's format recorded in release-signing.md
status: backlog
parent: S-0234
owner: alex
created: 2026-10-07T22:42:56Z
updated: 2026-10-07T22:42:56Z
transitions: []
stream: S-0234
tags: [dashboard]
touches: [design/system/release-signing.md]
---
# T-1231 Which digest RepoDigests records for a flaiover image pulled by index and by platform digest is checked, and the digest list's format recorded in release-signing.md

## Work

Pull a published multi-platform `ghcr.io/bytepunx/flaiover` release twice on the host's Docker, once by its index digest and once by this platform's manifest digest (`docker buildx imagetools inspect --raw` lists them), and read `docker image inspect --format '{{json .RepoDigests}}'` and `docker inspect` of a container run from each. Note the Docker version and whether the containerd image store is on, since it changes what is recorded. Remove the images you pulled afterwards.

Record the finding in `design/system/release-signing.md`, under How the image's signature is published, and from it the format of `flaiover_<version>.digests`: one line per digest a pull can leave behind (the index and each platform's manifest), each with the reference it was pushed under, so that S-0236 and S-0238 can match any of them.

Waits for nothing: it decides what the workflow writes.

## Done when

- `release-signing.md` names the Docker version tried, which digest `RepoDigests` held after each pull, and the digest list's line format, with an example.

## Notes
