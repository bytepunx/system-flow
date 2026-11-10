---
id: T-3070
type: task
nature: improvement
title: Document scripts/flai.sh's guard fallback on a failed build and the OpenCode plugin's 126 and 127 refusal
status: backlog
parent: S-0666
owner: alex
created: 2026-11-10T10:36:07Z
updated: 2026-11-10T10:36:07Z
transitions: []
stream: S-0666
tags: [guard, docs]
touches: [scripts/README.md, docs/users/flai.md, design/system/agent-adapters.md, docs/operators/index.md, design/system/strategic-agents.md]
after: [T-3066, T-3069]
---
# T-3070 Document scripts/flai.sh's guard fallback on a failed build and the OpenCode plugin's 126 and 127 refusal

## Work

Say what T-3066 and T-3069 changed where the fallback and the exit rule are described, linking the new ADR, and copy nothing it says at length:

- `scripts/README.md`: the `flai.sh` row names the guard's fallback on a failed build and the exit 127.
- `docs/users/flai.md`, `## Sub-agents`: the 127 refusal and the installed-flai fallback cover a failed build under `guard` too.
- `design/system/agent-adapters.md`, `### The guard`: below the exit table, how this repository's `scripts/flai.sh` reaches a verdict when its build fails, and that this repository's OpenCode plugin refuses on 126 and 127; `### A container harness of flai's own`: the fallback with `go` present.
- `docs/operators/index.md`, `### Strategic agents in a container (S-0509)`: what `scripts/flai.sh` runs when `go` is there but the build fails.
- `design/system/strategic-agents.md`, `### The live run (S-0567)`: I-0275 is fixed by S-0666.

Bump each document's `updated`. It waits for T-3066 and T-3069, whose behaviour it describes.

## Done when

- Each section above says what the code now does, and none still says a failed build fails open.
- `flai test` on the five files passes the markdown lint, and `flai check --strict` reports nothing new.

## Notes
