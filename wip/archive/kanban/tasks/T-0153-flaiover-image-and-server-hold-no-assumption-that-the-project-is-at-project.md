---
id: T-0153
type: task
nature: remediation
title: flaiover image and server hold no assumption that the project is at /project
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:24Z
updated: 2026-09-18T19:51:45Z
transitions:
  - to: ready
    at: 2026-09-18T19:50:39Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:50:39Z
    by: alex
  - to: done
    at: 2026-09-18T19:51:45Z
    by: alex
stream: S-0050
tags: []
touches: [flaiover]
---

# T-0153 flaiover image and server hold no assumption that the project is at /project

## Work
Read the `flaiover/Dockerfile`, `flaiover/src/hooks.server.ts`, `src/lib/server/repo.ts`, `src/lib/server/flai.ts`, the compose file, and `scripts/flaiover-*.sh` for anything that assumes `/project`: the `VOLUME`, comments, defaults, git `safe.directory`, readiness. `PROJECT_DIR` already drives the server and the flai environment; change only what breaks when it is another path, and keep `/project` as the image default for anyone running the image by hand.

## Done when
- Every assumption found is fixed or listed in the narrative as harmless, with the reason
- flaiover lint, type check, and tests pass

## Notes
