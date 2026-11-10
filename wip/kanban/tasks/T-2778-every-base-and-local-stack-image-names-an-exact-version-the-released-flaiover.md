---
id: T-2778
type: task
nature: remediation
title: Every base and local-stack image names an exact version, the released flaiover image its digest too, and dependabot watches the agent image folders
status: done
parent: S-0630
owner: alex
created: 2026-11-09T21:58:31Z
updated: 2026-11-10T05:53:10Z
transitions:
  - to: ready
    at: 2026-11-10T05:47:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:47:13Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:53:10Z
    by: agent-S-0630
stream: S-0630
tags: [ci, security, docker]
touches: [flaiover/Dockerfile, flai/images/bedrock/Dockerfile, flai/images/system-flow/Dockerfile, flaiover/compose.yaml, scripts/litellm.sh, ".github/dependabot.yml"]
usage:
  source: log
  seconds: 342
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 25
      output: 171
      cache_read: 2954169
      cache_write: 24192
      cost: 1.0945
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

Pins taken 2026-11-10 from the registries. flaiover's two `FROM` lines are literal `node:24.22.2-alpine3.24@sha256:a34ab6d7…`, so dependabot can move them; nothing passed `NODE_VERSION`, so the argument is gone. The Go build stages use `golang:1.26.0-alpine3.23` through a new `GO_ALPINE_VERSION`, the newest Alpine golang 1.26.0 ships on, and bedrock's runtime is `alpine:3.24.4`, so the claude-code layer's Alpine Node 24 is unchanged. otel-lgtm is pinned at 0.40.0, LiteLLM at v1.108.4, and LiteLLM's Postgres at 16.15-alpine3.24. `scripts/flaiover-image.sh` built (exit 0, about 2.5 minutes; it replaced the host's `flaiover:local`). `docker build --check` passes on all three Dockerfiles. `scripts/litellm.sh up` was not run: it needs a paid key. The bedrock and system-flow images were not built.
