---
id: T-3530
type: task
nature: remediation
title: flai guard reads the command and subcommand after flai as cobra finds them, so a value flag before the command takes its value with it
status: backlog
parent: S-0832
owner: alex
created: 2026-11-17T08:29:26Z
updated: 2026-11-17T08:29:53Z
transitions: []
stream: S-0832
tags: [flai, guard]
touches: [flai/internal/guard/guard_refusals.go, flai/internal/guard/guard_rules.go, flai/internal/guard/authors.go, flai/internal/guard/guard_tables.go, flai/internal/guard/guard.go, flai/internal/guard/roles.go, flai/internal/guard/guard_test.go]
---
# T-3530 flai guard reads the command and subcommand after flai as cobra finds them, so a value flag before the command takes its value with it

## Work

It waits for no task: it is the fix I-0358 asks for.

The proposed fix, from I-0358's one instance: stop listing the flags that take a value, and list instead the flags that take none, as cobra does. cobra's `Command.Find` strips the flags before each command word with `stripFlags`: a word that starts with `-` and has no `=` takes the next word with it unless it is a flag of that command, its persistent flags merged, whose `NoOptDefVal` is set, a bool flag; `--` ends the flags. So `flai --wait 5m kluster deploy`, `flai --by alex thread reply`, and `flai --reason x move` all run the command after the value, whatever flag stands first.

- In `flai/internal/guard/guard_refusals.go`, make `subcommands` read the words after `flai` that way, at each of its two levels, from a table of the flags that take no value. Keep its use for other programs, such as `gitSubcommand`'s callers, as it is.
- In `flai/internal/guard/guard_rules.go`, put that table beside `flaiValues`: flai's root flags that take no value, `--json`, `--yes`, `-y`, `--verbose`, `-v`, and `--help`, `-h`, as the built-in default. Give `Guard` in `flai/internal/guard/guard.go` a field for the no-value flags of flai's command tree, per command, which falls back to that default while nil; T-3532 fills it from the tree.
- Read the command the same way where a caller takes it by position: `names` in `flai/internal/guard/authors.go`, whose `authorValues` lacks `--wait` (so `flai --wait 5m thread reply TH-0001 --by other` passes the author check today), `readsRoles` and `hostAPIMethod` in `flai/internal/guard/guard_tables.go`, the planner's moves in `flai/internal/guard/guard.go`, and the plans in `flai/internal/guard/roles.go`. Keep `positionals` for the words after the command, where pflag parses the found command's own flags.
- Export the reading, for example `guard.FlaiCommand(words, noValue)`, so that T-3532 can compare it with cobra's.
- Update the comments of `flaiValues` and `authorValues` to say what each still reads.

Add tests to `flai/internal/guard/guard_test.go` that fail before the change: under the kluster rule, `flai --by alex kluster deploy S-0001` and `flai --reason x kluster deploy S-0001` are refused as deploys; `flai --json kluster deploy S-0001`, `flai -y kluster deploy S-0001`, and `flai --wait=5m kluster deploy S-0001` read `kluster deploy`; `flai --wait 5m thread reply TH-0001 --by other` is refused to a served agent as another author. Keep every existing guard test passing, `kluster_test.go`'s `--wait` lines among them.

## Done when

- `subcommands` reads `--by`, `--reason`, `--wait`, and any other flag with a value before flai's command as taking its value, and a no-value flag as taking none.
- The callers listed above read the command flai runs, and the new tests in `guard_test.go` pass where they failed before.
- `flai test` passes on the changed paths.

## Notes
