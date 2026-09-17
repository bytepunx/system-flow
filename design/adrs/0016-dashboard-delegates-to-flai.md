---
id: ADR-0016
title: flaiover delegates writes and metrics to the flai binary
status: accepted
date: 2026-09-17
supersedes: [ADR-0007]
superseded_by: []
---

# ADR-0016 flaiover delegates writes and metrics to the flai binary

## Context

ADR-0007 planned for the dashboard server to port flai's validation and metrics rules to TypeScript and keep the two aligned with a fixture test. By the time the board story started, the Go side held the transition rules, the byte-stable front matter writer, board order maintenance, narrative index regeneration, the archive planner, and the metrics reference, all with tests. Porting them would double the surface that has to stay correct, and every rule change would need two implementations and a fixture update.

## Decision

flaiover performs every mutation (move, block, unblock, narrative log) and every metric by invoking the `flai` binary with `--json` inside the mounted project, and reports flai's rule text verbatim when a write is refused. The server finds the binary through `FLAI_BIN`, then `PATH`; the Docker image bundles a `flai` built from the same commit. Reads stay in the dashboard's own reader, which parses the same front matter and is fixture-tested against flai's output. Without a binary the dashboard runs read-only and says so.

This supersedes ADR-0007 in part: the "server ports the validation and metrics rules" consequence is withdrawn. The SPA-with-node-adapter shape, the mount, the port, and the no-authentication posture stand.

## Consequences

- One reference implementation for rules and metrics; the dashboard cannot drift from the CLI.
- The image carries a second binary and the build is a two-stage, two-language Dockerfile (S-0015).
- Writes cost a process spawn each, acceptable for a local single-user tool.
- The dashboard's version and flai's move together; an image tag implies the flai it bundles.

## Alternatives considered

- Port the rules and test alignment: the original plan; rejected as duplication that would drift.
- A long-running flai server mode: cleaner than spawning, but a new interface for no measured need; revisit if spawn latency matters.
