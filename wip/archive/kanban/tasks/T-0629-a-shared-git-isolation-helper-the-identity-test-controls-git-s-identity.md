---
id: T-0629
type: task
nature: feature
title: "A shared git-isolation helper; the identity test controls git's identity"
status: done
parent: S-0180
owner: arobson
created: 2026-10-01T08:22:54Z
updated: 2026-10-01T08:23:04Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
  - to: in-progress
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
  - to: done
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
stream: S-0180
tags: []
usage:
  source: log
  seconds: 0
  models: []
---

# T-0629 A shared git-isolation helper; the identity test controls git's identity

## Work

Add `flai/internal/gittest`: `Isolate` (a global config of the test's own, no system config, no identity from the environment), `Identity`, and `NoIdentity` (`user.useConfigOnly`, so git does not work one out from the hostname). Use it in `accept_resume_test.go` and `doc_test.go`.

## Done when

`TestAcceptRefusesBeforeChangingAnythingWithoutIdentity` passes on macOS, and the helper has a test of its own.

## Notes
