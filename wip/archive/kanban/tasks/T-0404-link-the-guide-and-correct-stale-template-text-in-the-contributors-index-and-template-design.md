---
id: T-0404
type: task
nature: feature
title: Link the guide and correct stale template text in the contributors index and template design
status: done
parent: S-0018
owner: alex
created: 2026-09-24T08:12:56Z
updated: 2026-09-24T08:15:47Z
transitions:
  - to: ready
    at: 2026-09-24T08:13:18Z
    by: agent-S-0018
  - to: in-progress
    at: 2026-09-24T08:14:59Z
    by: agent-S-0018
  - to: done
    at: 2026-09-24T08:15:47Z
    by: agent-S-0018
stream: S-0018
tags: []
---

# T-0404 Link the guide and correct stale template text in the contributors index and template design

## Work
Link the guide from `docs/contributors/index.md` and `docs/README.md`. Replace the index's release and publishing text, which still says `flai accept` bumps and publishes the template (S-0087 moved that to `flai release --pending` and `flai push --pending --publish`). Correct `design/system/template.md` where the guide found it out of date with the code.

## Done when
No page under `docs/contributors` or `design/system/template.md` contradicts the guide or the code.

## Notes
