---
id: T-2186
type: task
nature: remediation
title: flai guard's help names the host configuration rule and the reference is regenerated
status: backlog
parent: S-0511
owner: alex
created: 2026-10-28T08:19:55Z
updated: 2026-10-28T08:19:55Z
transitions: []
stream: S-0511
tags: [cli, guard, docs]
touches: [flai/cmd/guard.go, docs/users/flai-reference.md]
after: [T-2184]
---
# T-2186 flai guard's help names the host configuration rule and the reference is regenerated

## Work

In the long help of `flai guard` in `flai/cmd/guard.go`, after the paragraph that says no sub-agent and no session flai serve starts changes the manifest's shared paths, add one in the same voice saying that none of them runs a `flai serve` sub-command that writes the host configuration or turns a host action on or off, `flai serve enable`, `disable`, `agent set`, `clear`, `harness` and `provider` with a setting or `--reset`, `checks set`, `timeout`, and `clear`, `components idle` and `max-servers`, `forge` with a setting or `--reset`, `notify set` and `clear`, `import add` and `remove`, nor `flai config set` or `flai template use`, which write the same file, because the agent's flai reads the operator's `~/.flai/config.json`; the reading forms stay open, and only the operator's own session changes the host configuration. Cite the ADR T-2184 recorded. Then run `make flai-reference` so `docs/users/flai-reference.md` carries the new help, and commit both. Waits for T-2184 for the ADR's number; runs beside T-2185, with which it shares no path.

## Done when

`flai guard --help` prints the paragraph, `docs/users/flai-reference.md` matches the help (`make flai-reference` changes nothing when run again), and `flai test flai/cmd` passes.

## Notes

`flai/cmd/guard_test.go` tests the command's behaviour, not its help text; add nothing there unless a test reads the help.
