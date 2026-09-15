---
title: Operators guide
updated: 2026-09-15
status: draft
---

# system-flow for operators

## Running the dashboard

`flai dashboard` runs `ghcr.io/bytepunx/flaiover` with the repository mounted read-write at `/project`, bound to `127.0.0.1:4242` by default. Change the image, tag, or port in `~/.flai/config.json` or per project in `system-flow.yaml` under `dashboard`.

Without `flai`:

```bash
docker run --rm -p 127.0.0.1:4242:3000 -v "$PWD:/project" -e PROJECT_DIR=/project --user "$(id -u):$(id -g)" ghcr.io/bytepunx/flaiover:latest
```

## Security posture

The dashboard has no authentication and can write to the mounted repository. Keep it bound to localhost. Do not put it behind a public reverse proxy. If a team needs shared access, run it on a host they reach over a private network or VPN, and treat the mount as you would a shared working copy.

## Requirements

- Docker Engine 24 or newer on `PATH`
- `git` on `PATH` for template cloning
