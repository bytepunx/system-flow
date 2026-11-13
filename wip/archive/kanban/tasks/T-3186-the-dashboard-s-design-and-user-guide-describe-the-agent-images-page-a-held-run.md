---
id: T-3186
type: task
nature: feature
title: The dashboard's design and user guide describe the Agent images page, a held run, and the retry behind the agent host action
status: done
parent: S-0757
owner: alex
created: 2026-11-12T16:04:18Z
updated: 2026-11-13T17:05:41Z
transitions:
  - to: ready
    at: 2026-11-13T17:00:30Z
    by: agent-S-0757
  - to: in-progress
    at: 2026-11-13T17:00:30Z
    by: agent-S-0757
  - to: done
    at: 2026-11-13T17:05:41Z
    by: agent-S-0757
stream: S-0757
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-3184, T-3185]
usage:
  source: log
  seconds: 311
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 22228
      cache_read: 3330215
      cache_write: 125212
      cost: 1.8437
---
# T-3186 The dashboard's design and user guide describe the Agent images page, a held run, and the retry behind the agent host action

## Work

- `design/system/flaiover-dashboard.md`: the Ops › Agent images page, its methods (`images.get`, `images.log`, `images.ensure`), the `agent` host action that gates the retry, and a story's agent held on a job. Bump `updated`.
- `docs/users/flaiover.md`: what the operator sees while an image comes back, how to read a failure, and how to retry.

It waits for T-3184 and T-3185, whose page and held run it describes.

## Done when

- Both files describe the page and the held run as built.
- `design-routes` (`flaiover/src/lib/server/design-routes.test.ts`), `flai check --strict`, and the markdown lint pass.

## Notes
