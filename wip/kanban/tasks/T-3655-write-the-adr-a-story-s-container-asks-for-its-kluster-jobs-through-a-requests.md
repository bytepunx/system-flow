---
id: T-3655
type: task
nature: remediation
title: "Write the ADR: a story's container asks for its kluster jobs through a requests folder in its agent cache folder, which flai serve also takes from"
status: backlog
parent: S-0849
owner: alex
created: 2026-11-28T00:04:37Z
updated: 2026-11-28T00:04:37Z
transitions: []
stream: S-0849
tags: [flai, kluster, adr]
touches: [design/adrs/README.md]
---
# T-3655 Write the ADR: a story's container asks for its kluster jobs through a requests folder in its agent cache folder, which flai serve also takes from

## Work

Record the fix before building it, as the story's goal asks. Write it with `flai adr new`, then `flai adr topics` with `kluster` and `harness`. It refines ADR-0292 §1 and ADR-0319. Confirm the decision on the plan's thread on S-0849 first: if the operator chose the alternative, a container's `--job` refusing at once and naming the MCP tool, write that one instead and say so in this task's notes.

The decision to record:

- **Where a container asks.** Where the main checkout's `.flai-cache/kluster` is absent, as in a story's container, `flai kluster deploy --job`, `teardown --job`, and `flai checks run --job` write their request in `.flai-cache/agents/S-nnnn/kluster-requests/`. That is the story's agent cache folder, which the container mounts read-write already (ADR-0319), so no mount is added. The folder sits outside `kluster/`, the copy, because a teardown that passes removes the copy whole. The file names and the protocol are ADR-0292 §1's: request, taken, answer, and `S-nnnn.job.json`.
- **What flai serve takes.** Each look of the request pass also lists the `kluster-requests/` folder of each story's agent cache folder, `S-nnnn` folders only, never a role's. It refuses a request whose story is not the folder's. It writes its answer in the folder it took the request from. It writes the story's job state there as well as in the main checkout's `requests/`, so `--wait` and the found job's state read from the container. The sweep of interrupted jobs at start covers both.
- **On the host** nothing changes. The MCP tools keep writing under the main checkout's `requests/`.
- **Alternatives.** Mounting `.flai-cache/kluster/requests` is one: it holds every story's requests and job states, and a folder bind-mount breaks when the folder is removed and made again. A container's `--job` refusing and naming the MCP tool is another: the documented command would stay broken where agents now run.

It waits for nothing: it is the first layer.

## Done when

- An accepted ADR under `design/adrs/` records the decision, with its context from I-0405's instance and its alternatives, and `refines: [ADR-0292, ADR-0319]`.
- `design/adrs/README.md` lists it, as flai generates it.
- `flai check --strict` and the markdown lint pass on the ADR.

## Notes
