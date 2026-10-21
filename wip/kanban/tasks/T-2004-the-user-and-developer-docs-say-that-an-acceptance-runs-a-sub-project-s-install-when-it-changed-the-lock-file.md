---
id: T-2004
type: task
nature: remediation
title: The user and developer docs say that an acceptance runs a sub-project's install when it changed the lock file
status: done
parent: S-0483
owner: alex
created: 2026-10-21T12:11:20Z
updated: 2026-10-21T13:27:57Z
transitions:
  - to: ready
    at: 2026-10-21T13:17:04Z
    by: agent-S-0483
  - to: in-progress
    at: 2026-10-21T13:17:05Z
    by: agent-S-0483
  - to: done
    at: 2026-10-21T13:27:57Z
    by: agent-S-0483
stream: S-0483
tags: [flai, docs]
touches: [docs/users/flai.md, design/system/devex.md]
after: [T-2000]
usage:
  source: log
  seconds: 652
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 10503
      cache_read: 1886139
      cache_write: 95069
      cost: 1.1629
---
# T-2004 The user and developer docs say that an acceptance runs a sub-project's install when it changed the lock file

## Work

Document the behaviour T-1998's ADR decides, for users and builders, as `documentation.md` asks in the same story. It waits for T-2000 because the key's name and its `flai check` rules are fixed there. It runs with T-2003, which changes only Go files.

- In `docs/users/flai.md` under `## Accept and release`, say that `flai accept`, in every mode and from the merge queue, runs the `install` of each sub-project whose lock file or `package.json` the accepted change touched, in the main checkout, and what a failed install prints. Link the ADR and the `projects[].install` row of `docs/operators/settings.md`.
- In `design/system/devex.md`, in the Test tiers row, say that the main checkout's `flaiover/node_modules` is kept current by the acceptance through `scripts/flaiover-install.sh --if-needed`, so a test that links it finds what main's lock file pins.

## Done when

- Both documents say it, with the ADR linked, and their `updated` dates are bumped.
- `flai test docs/users/flai.md design/system/devex.md` passes.

## Notes
