---
id: T-3045
type: task
nature: improvement
title: Record in an ADR refining ADR-0247 that a strategic agent's container mounts .git/config read-only
status: backlog
parent: S-0688
owner: alex
created: 2026-11-10T10:34:46Z
updated: 2026-11-10T10:34:46Z
transitions: []
stream: S-0688
tags: [adr, container, security]
touches: [design/adrs/README.md]
---
# T-3045 Record in an ADR refining ADR-0247 that a strategic agent's container mounts .git/config read-only

## Work

Write the remedy as a decision before building it, as the story's goal asks. Run `flai adr new` for an ADR that refines ADR-0247 and ADR-0278: a strategic agent's container mounts the main checkout's `.git/config` read-only over the read-write checkout, as it mounts `.git/hooks`. Give its context from I-0257's two instances: `.git/config` can name commands git runs on the host (`core.fsmonitor`, `core.sshCommand`, `core.hooksPath`, `core.pager`, an alias), and `flai guard` does not see a write through the shell. ADR-0247 kept it writable only for `flai accept`'s branch deletion, which ADR-0278 moved to the host. Name what still runs there and writes no config: `flai subproject --autocommit` and `flai formula run --autocommit`. Weigh the alternatives: a check of `.git/config` after each strategic run, `.git` read-only as a whole, and doing nothing. Say what the read-only file mount does to git's lock-and-rename write (it fails, which is the point) and that a host rewrite of the file while a container runs is not seen inside it. Accept it with `flai adr accept` and its topics with `flai adr topics`.

This task waits for nothing: the decision comes first.

## Done when

- The ADR is accepted, refines ADR-0247 and ADR-0278, and `design/adrs/README.md` lists it.
- `flai check --strict` passes.

## Notes
