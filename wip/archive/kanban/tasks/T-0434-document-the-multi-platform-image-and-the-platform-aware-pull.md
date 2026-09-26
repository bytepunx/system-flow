---
id: T-0434
type: task
nature: feature
title: Document the multi-platform image and the platform-aware pull
status: done
parent: S-0119
owner: arobson
created: 2026-09-26T03:17:08Z
updated: 2026-09-26T03:22:50Z
transitions:
  - to: ready
    at: 2026-09-26T03:17:26Z
    by: agent-S-0119
  - to: in-progress
    at: 2026-09-26T03:21:40Z
    by: agent-S-0119
  - to: done
    at: 2026-09-26T03:22:50Z
    by: agent-S-0119
stream: S-0119
tags: []
---

# T-0434 Document the multi-platform image and the platform-aware pull

## Work

- `design/system/flaiover-dashboard.md`: the image is published for linux/amd64 and linux/arm64, and the build stage runs on the build platform.
- `docs/users/flai.md` and the operators' install runbook: `flai dashboard` pulls the daemon's platform and replaces a present image of another platform.
- `flai/CHANGELOG` or release notes, where the release tooling does not write them.

## Done when

- The docs describe the behavior as built, and `flai check --strict` is clean.

## Notes
