---
id: T-2958
type: task
nature: remediation
title: I-0339 is closed with the reason that scripts/agent-image.sh builds with no login when the credential helper fails
status: backlog
parent: S-0713
owner: alex
created: 2026-11-10T10:24:11Z
updated: 2026-11-10T10:24:11Z
transitions: []
stream: S-0713
tags: [issues]
touches: [design/issues/I-0339-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md, design/issues/summary.md]
after: [T-2956, T-2957]
---
# T-2958 I-0339 is closed with the reason that scripts/agent-image.sh builds with no login when the credential helper fails

## Work

Run `flai issue close I-0339 --reason "<reason>"`. The reason says that `scripts/agent-image.sh` now probes the `credsStore` helper the Docker configuration names. When the helper fails, the script builds with a scratch `DOCKER_CONFIG` of `{}`, so the public base images are pulled with no login, and leaves the operator's configuration as it was. It also names `scripts/agent-image-test.sh`, which the smoke tier runs, as the reproduction. This rewrites the issue and its row in `design/issues/summary.md`. It waits for T-2956 and T-2957, so that the reason can name the test and the docs as they are.

## Done when

- I-0339 is closed, with that reason.
- `design/issues/summary.md` shows it closed.
- `flai check --strict` is clean.

## Notes
