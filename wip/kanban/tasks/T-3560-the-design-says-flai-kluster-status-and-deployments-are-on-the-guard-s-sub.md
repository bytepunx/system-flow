---
id: T-3560
type: task
nature: remediation
title: The design says flai kluster status and deployments are on the guard's sub-agent reads
status: backlog
parent: S-0829
owner: alex
created: 2026-11-17T12:28:56Z
updated: 2026-11-17T12:28:56Z
transitions: []
stream: S-0829
tags: [docs, guard]
touches: [design/system/flai-cli.md, design/system/kluster.md]
after: [T-3559]
---
# T-3560 The design says flai kluster status and deployments are on the guard's sub-agent reads

## Work

In `design/system/flai-cli.md`, add `kluster status` and `kluster deployments` (S-0829) to the `flai guard` row's list of the flai commands a sub-agent may run. In `design/system/kluster.md`, under the guard's history in the paragraph that says `flai kluster status` is a read any agent may run, add that since S-0829 the guard's sub-agent read table lists `status` and `deployments`, so a sub-agent's call no longer falls through to the generic refusal. Bump each file's `updated`.

It waits for T-3559 because it describes the read table as T-3559 leaves it: the subcommands named must be the ones the fix added.

## Done when

- Both documents name `kluster status` and `kluster deployments` as reads open to a sub-agent, citing S-0829.
- `flai test design/system/flai-cli.md design/system/kluster.md` passes the markdown lint.

## Notes
