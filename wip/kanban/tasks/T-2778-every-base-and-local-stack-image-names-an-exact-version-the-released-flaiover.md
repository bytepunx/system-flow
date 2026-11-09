---
id: T-2778
type: task
nature: remediation
title: Every base and local-stack image names an exact version, the released flaiover image its digest too, and dependabot watches the agent image folders
status: backlog
parent: S-0630
owner: alex
created: 2026-11-09T21:58:31Z
updated: 2026-11-09T21:58:31Z
transitions: []
stream: S-0630
tags: [ci, security, docker]
touches: [flaiover/Dockerfile, flai/images/bedrock/Dockerfile, flai/images/system-flow/Dockerfile, flaiover/compose.yaml, scripts/litellm.sh, ".github/dependabot.yml"]
---
# T-2778 Every base and local-stack image names an exact version, the released flaiover image its digest too, and dependabot watches the agent image folders

## Work

- `flaiover/Dockerfile:9, 15, 26`: change `ARG NODE_VERSION=24` to the current exact Node 24 Alpine tag, `node:<24.x.y>-alpine<a.b>`. Add that tag's index digest to both `FROM` lines, `node:${NODE_VERSION}-alpine<a.b>@sha256:<digest>`, since a signed `flaiover/v*` release builds from this file. Do the same for any `golang` stage the file has. Read the digest with `docker buildx imagetools inspect <image>:<tag>`.
- `flai/images/bedrock/Dockerfile:11-15, 31`: make `GO_VERSION` and `ALPINE_VERSION` exact, such as `1.26.0` and `3.24.<n>`. Make the `golang` stage `golang:${GO_VERSION}-alpine${ALPINE_VERSION%.*}` or an equivalent exact tag. `scripts/agent-image.sh` passes `GO_VERSION` from `flai/go.mod`, so the default only matters when it is built by hand.
- `flai/images/system-flow/Dockerfile:14`: pin its `golang:${GO_VERSION}-alpine` stage to an exact Alpine release, as above. `koalaman/shellcheck-alpine:${SHELLCHECK_VERSION}` names `v0.10.0` already, so leave it.
- `flaiover/compose.yaml:29`: replace `grafana/otel-lgtm:latest` with its current release tag.
- `scripts/litellm.sh:37, 87`: replace the default `ghcr.io/berriai/litellm:main-stable` with a current release tag, and say in the comment at line 37 that `LITELLM_IMAGE` overrides it. Check the proxy still starts with `scripts/litellm.sh up` where Docker and a key are at hand. Otherwise say in the task's Notes that it was not run.
- `.github/dependabot.yml`: make the `docker` entry cover `/flaiover`, `/flai/images/bedrock`, and `/flai/images/system-flow`, with `directories:`, so the new pins move weekly as the actions do.

No agent image is built or published by a workflow (`design/tech/docker.md`), so only `flaiover/Dockerfile` takes a digest. The others take exact tags that dependabot moves.

It waits for no task: no other task of S-0630 changes these files.

## Done when

- No `FROM` or `image:` in the five files names `latest`, `main-stable`, or a major-only or minor-only tag.
- Both `FROM` lines in `flaiover/Dockerfile` carry an `@sha256:` digest.
- `scripts/flaiover-image.sh` (`make flaiover-image`) builds where Docker is at hand. Otherwise the task's Notes say it was not run.
- `.github/dependabot.yml` lists the three image folders under its `docker` entry, and `scripts/shellcheck.sh scripts/litellm.sh` passes.

## Notes
