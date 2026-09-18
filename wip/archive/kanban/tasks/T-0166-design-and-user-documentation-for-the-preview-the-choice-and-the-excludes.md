---
id: T-0166
type: task
nature: improvement
title: Design and user documentation for the preview, the choice, and the excludes
status: done
parent: S-0051
owner: alex
created: 2026-09-18T21:00:24Z
updated: 2026-09-18T21:06:31Z
transitions:
  - to: ready
    at: 2026-09-18T21:05:50Z
    by: alex
  - to: in-progress
    at: 2026-09-18T21:05:50Z
    by: alex
  - to: done
    at: 2026-09-18T21:06:31Z
    by: alex
stream: S-0051
tags: []
touches: [design/system, docs]
---

# T-0166 Design and user documentation for the preview, the choice, and the excludes

## Work
Update `design/system/flaiover-dashboard.md` (the acceptance and move endpoints, the confirmation, what `flai dashboard` passes into the container), `design/system/flai-cli.md` (`flai accept --dry-run`, `flai dashboard`), `docs/users/flaiover.md` (what the confirmation shows and the choice), `docs/users/flai.md`, and `docs/operators/index.md` (the excludes mount and the environment variables). This is a refinement of S-0046 and ADR-0022's mount contract, not a new decision; say so in the narrative.

## Done when
- Each document states the behaviour as built
- Markdown lint and `flai check --strict` pass

## Notes
