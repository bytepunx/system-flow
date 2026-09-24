---
id: S-0108
type: story
nature: remediation
title: "CI after S-0106: a stopped child reported before its process ends, and a test that pinned the host actions"
status: review
owner: alex
created: 2026-09-24T03:22:00Z
updated: 2026-09-24T03:34:32Z
transitions:
  - to: ready
    at: 2026-09-24T03:22:25Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T03:22:25Z
    by: system-flow
  - to: review
    at: 2026-09-24T03:34:32Z
    by: system-flow
tags: [cli, dashboard]
touches: [flai/internal/host]
---
# S-0108 CI after S-0106: a stopped child reported before its process ends, and a test that pinned the host actions

## Goal

CI is green again on main after S-0106's push: flai, release flai v1.15.0, flaiover, and the system-flow check all failed.

## Acceptance criteria

- [x] `flai host` reports a child as stopped only once its process has ended; until then it is stopping, with its PID, and `TestTheHostKeepsTheMCPServersServeAsksFor` no longer races on a slow runner
- [x] flaiover's push-refusal test no longer breaks when a host action is added
- [x] `flai check --strict` has no warning on main
- [ ] the flai, flaiover, and system-flow check workflows pass on main, and flai's next release publishes

## Tasks
- T-0382 flai host reports a child stopping, with its PID, until its process has ended
- T-0383 flaiover's push-refusal test holds whatever host actions there are
- T-0384 CI green on main; the release published

## Notes

The runs that failed: 35949437943 (flai), 35949437971 (release flai flai/v1.15.0, whose pre-release test hook hit the same test), 35949437986 (flaiover), and 35949437955 (system-flow check).
