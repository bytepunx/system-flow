---
title: Operators guide
updated: 2026-09-15
status: draft
---

# system-flow for operators

## Running the dashboard

`flai dashboard` runs `ghcr.io/bytepunx/flaiover` detached as `flaiover-<project>` with the repository mounted read-write at `/project`, bound to `127.0.0.1:4242` by default, as the invoking user. `flai dashboard status`, `logs`, and `stop` manage it. Change the image, tag, or port in `~/.flai/config.json` or per project in `system-flow.yaml` under `dashboard`.

Without `flai`, the equivalent is:

```bash
docker run --detach --rm --name flaiover-myproject \
  --publish 127.0.0.1:4242:3000 --volume "$PWD:/project" --env PROJECT_DIR=/project \
  --user "$(id -u):$(id -g)" ghcr.io/bytepunx/flaiover:latest
```

## The image

`ghcr.io/bytepunx/flaiover` is built from `flaiover/Dockerfile` at the repository root: a `flai` binary from the same commit at `/usr/local/bin/flai`, the SvelteKit build on `node:24-alpine`, listening on `3000`, running as an unprivileged user by default and working under any `--user`. Tags: `latest` (main), `X.Y.Z` and `X` from `flaiover/vX.Y.Z` release tags, and `sha-<commit>`. While the repository is private the package is too: `docker login ghcr.io` with a token that has `read:packages` before `flai dashboard` can pull it. The container writes flai's config and cache under the mounted project's `.flai-cache/`, which is git-ignored.

Build locally with `make flaiover-image` (tag `flaiover:local`) and run it with `flai dashboard --image flaiover --tag local`.

## Security posture

The dashboard has no authentication and can write to the mounted repository. Keep it bound to localhost. Do not put it behind a public reverse proxy. If a team needs shared access, run it on a host they reach over a private network or VPN, and treat the mount as you would a shared working copy.

## Requirements

- Docker Engine 24 or newer on `PATH`
- `git` on `PATH` for template cloning
