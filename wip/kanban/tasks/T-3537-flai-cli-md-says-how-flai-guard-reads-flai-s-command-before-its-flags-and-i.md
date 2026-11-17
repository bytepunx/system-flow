---
id: T-3537
type: task
nature: remediation
title: flai-cli.md says how flai guard reads flai's command before its flags, and I-0358 is closed
status: done
parent: S-0832
owner: alex
created: 2026-11-17T08:30:00Z
updated: 2026-11-17T10:54:38Z
transitions:
  - to: ready
    at: 2026-11-17T10:54:08Z
    by: agent-S-0832
  - to: in-progress
    at: 2026-11-17T10:54:08Z
    by: agent-S-0832
  - to: done
    at: 2026-11-17T10:54:38Z
    by: agent-S-0832
stream: S-0832
tags: [docs, guard]
touches: [design/system/flai-cli.md, design/issues/I-0358-flai-guard-reads-a-subcommand-s-value-flag-before-flai-s-command-as-the-command.md, design/issues/summary.md]
after: [T-3532]
usage:
  source: log
  seconds: 26
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 8
      output: 61
      cache_read: 1004006
      cache_write: 2485
      cost: 0.3842
---
# T-3537 flai-cli.md says how flai guard reads flai's command before its flags, and I-0358 is closed

## Work

It waits for T-3532, because it describes the reading as T-3530 and T-3532 build it, and the issue closes only once the fix and its test against cobra are in.

- In `design/system/flai-cli.md`, in the `flai guard` row of the command table, add one sentence: the guard reads the command and subcommand after `flai` as cobra's `Find` does, a flag before them taking the next word unless flai's command tree defines it with no value or it carries its value after `=`, so that `flai --wait 5m kluster deploy` and `flai --by alex thread reply` are judged as the commands flai runs (S-0832, I-0358). Replace any wording there that names `--config` and `--wait` as the only flags read before the command.
- Close the issue: `flai issue close I-0358 --reason "S-0832: flai guard reads flai's command as cobra's Find does, from the flags flai's command tree defines with no value, so a value flag before the command no longer hides it; flai/cmd/guard_test.go checks the reading against cobra"`.

## Done when

- `flai-cli.md`'s `flai guard` row says how the guard reads flai's command, and passes the markdown lint.
- I-0358 is closed with its reason, and `design/issues/summary.md` no longer lists it.
- `flai test` passes on the changed paths.

## Notes
