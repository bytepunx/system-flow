---
id: S-0035
type: story
nature: feature
title: flai dashboard runs from a private registry or a local build
status: done
parent: E-0003
owner: alex
created: 2026-09-17T07:49:28Z
updated: 2026-09-17T19:02:33Z
transitions:
  - to: ready
    at: 2026-09-17T07:49:29Z
    by: alex
  - to: in-progress
    at: 2026-09-17T07:49:29Z
    by: alex
  - to: review
    at: 2026-09-17T18:21:19Z
    by: alex
  - to: done
    at: 2026-09-17T19:02:33Z
    by: alex
tags: [dashboard, cli]
---

# S-0035 flai dashboard runs from a private registry or a local build

## Goal
`flai dashboard` gets the operator to a running flaiover on the first try. While the image on GHCR is private, the command logs Docker into the registry with the same token sources as the installer, explains the missing `read:packages` scope when the token lacks it, and offers `--build` to build the image from the source tree instead of pulling.

## Acceptance criteria
- [x] When `docker pull` is refused, `flai dashboard` logs into the registry with `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` and retries once; the login name comes from `gh api user` or `GITHUB_ACTOR`, falling back to a placeholder
- [x] When no token is available or the login is refused, the error names the fix: `gh auth refresh -h github.com -s read:packages`, or make the package public
- [x] `flai dashboard --build` builds the image from `flaiover/` in the monorepo (or refuses with a clear message elsewhere) and runs it as `flaiover:local`; `scripts/flaiover-image.sh` stays the script behind it
- [x] Tests cover the login retry, the scope hint, and the build path with the fake runner
- [x] docs/operators/index.md has an "Access to the image" section covering token scopes, package visibility, and the local build
- [x] The container publishes on `0.0.0.0` by default (operator's call, 2026-09-17) so the dashboard is reachable from other hosts; `dashboard.bind` in `system-flow.yaml` and config, and `--bind`, restrict it; compose and the security posture docs say so

## Tasks
- T-0112 Registry login on unauthorized pull with token sources, login name resolution, one retry, scope hint
- T-0113 flai dashboard --build: build flaiover:local from the monorepo via docker build, then run it
- T-0114 Tests with the fake runner; operator docs section on image access
- T-0115 Bind the dashboard to all interfaces by default: dashboard.bind in manifest and config, --bind flag, compose, docs

## Notes
- Found 2026-09-17: `docker pull ghcr.io/bytepunx/flaiover:latest` returns unauthorized because the package is private and Docker has no ghcr.io login; the operator's gh token lacks `read:packages`. The dashboard runs today from `make flaiover-image` and `flai dashboard --image flaiover --tag local`.
