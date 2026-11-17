---
id: T-3532
type: task
nature: remediation
title: flai guard takes the flags that take no value from flai's command tree, and a test holds its reading of the command to cobra's Find
status: backlog
parent: S-0832
owner: alex
created: 2026-11-17T08:29:35Z
updated: 2026-11-17T08:29:44Z
transitions: []
stream: S-0832
tags: [flai, guard]
touches: [flai/cmd/guard.go, flai/cmd/guard_test.go]
after: [T-3530]
---
# T-3532 flai guard takes the flags that take no value from flai's command tree, and a test holds its reading of the command to cobra's Find

## Work

It waits for T-3530, because it fills the `Guard` field T-3530 adds and tests the reading T-3530 exports.

- In `flai/cmd/guard.go`, where the guard's `Commands` are taken from `cmd.Root().Commands()`, fill the field T-3530 adds with the flags of flai's command tree that take no value: for the root and for each command, its own and its persistent flags, merged from its parents, whose `NoOptDefVal` is set, with their shorthands. The guard then reads a flag before the command as cobra does, for every flag flai defines now or adds later, with no list to keep in step by hand.
- Add a test to `flai/cmd/guard_test.go` that reproduces I-0358 against cobra itself: for lines such as `--wait 5m kluster deploy S-0001`, `--by alex thread reply TH-0001`, `--reason x move S-0001 cancelled`, `--config c.json --json kluster deploy S-0001`, `-y kluster deploy S-0001`, `--wait=5m checks run S-0001`, and `kluster --json deploy S-0001`, the command and subcommand the exported reading gives with the tree's no-value flags equal the names of the command flai's root `Find` returns for the same words.

## Done when

- A served agent's guard built by `flai guard` reads the command and subcommand flai would run for each line of the test.
- The test in `flai/cmd/guard_test.go` passes, and its `--wait 5m` and `--by alex` lines read `5m` and `alex` as the command when run against the reading before T-3530.
- `flai test` passes on the changed paths.

## Notes
