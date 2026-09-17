---
id: T-0047
type: task
nature: feature
title: "Fix flai workflow: build lint from source, monorepo test errors only"
status: done
parent: S-0010
owner: alex
created: 2026-09-16T04:32:20Z
updated: 2026-09-16T04:32:52Z
transitions:
  - to: ready
    at: 2026-09-16T04:32:20Z
    by: agent
  - to: in-progress
    at: 2026-09-16T04:32:20Z
    by: agent
  - to: done
    at: 2026-09-16T04:32:52Z
    by: agent
stream: S-0010
tags: [ci, remediation]
---

# T-0047 Fix flai workflow: build lint from source, monorepo test errors only

## Work
The golangci-lint action downloads a prebuilt binary built with an older Go than flai/go.mod targets and refuses the config; switch the action to build from source with the workflow's Go. TestMonorepoIsClean fails on warnings, so a done story awaiting archive turns the integration tier red on every branch; make it fail on errors only and leave warnings to the strict check job.

## Done when
flai workflow green on main after push.

## Notes
