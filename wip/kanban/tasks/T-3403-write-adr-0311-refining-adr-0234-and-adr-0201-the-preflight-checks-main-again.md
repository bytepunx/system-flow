---
id: T-3403
type: task
nature: remediation
title: "Write ADR-0311 refining ADR-0234 and ADR-0201: the preflight checks main again with a flai built from main when the running flai is behind, else holds naming the release to install"
status: done
parent: S-0797
owner: alex
created: 2026-11-16T01:44:13Z
updated: 2026-11-16T03:14:19Z
transitions:
  - to: ready
    at: 2026-11-16T03:03:03Z
    by: agent-S-0797
  - to: in-progress
    at: 2026-11-16T03:03:04Z
    by: agent-S-0797
  - to: done
    at: 2026-11-16T03:14:19Z
    by: agent-S-0797
stream: S-0797
tags: [adr, merge-queue]
touches: [design/adrs/0311-the-merge-queue-s-preflight-checks-main-again-with-a-flai-built-from-main-when.md, design/adrs/README.md]
usage:
  source: log
  seconds: 675
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 2912
      cache_read: 486680
      cache_write: 18692
      cost: 0.2708
---
# T-3403 Write ADR-0311 refining ADR-0234 and ADR-0201: the preflight checks main again with a flai built from main when the running flai is behind, else holds naming the release to install

## Work

Write the ADR the story's goal asks for with `flai adr new`, with `refines: [ADR-0234]`, and accept it.

- Context: I-0382. `preflightCheck` in `flai/internal/mergequeue/preflight.go` runs the running flai's own `check.Run` on the main checkout. On 2026-11-15 the installed flai 1.99.0 did not know the role files' `when: {to: <state>}` form that S-0748 wrote. It reported `roles.file` on two role files that were right, and the hold thread told the operator to fix them. The queue built no batch for about 3h, until 1.100.0 was installed.
- Decision, to settle in the ADR:
  1. The running flai is behind main when main's flai sub-project, the manifest's project at `flai`, has commits the running flai was not built from. When `buildinfo` gives a commit in main's history, that means `git log <commit>..HEAD -- flai/` lists one. When the commit is `unknown` or not in main's history, compare against the newest `flai/v*` tag the running version is not below.
  2. When the preflight's own `flai check` finds errors and the running flai is behind, the preflight runs the check again with the manifest's `flai.command` (`check --json`) in its scratch worktree at main's head, as `confirmLint` checks the lint, so no file of the main checkout changes. The errors that check gives replace the running flai's.
  3. When the running flai is behind and no `flai.command` is named, or the command gives no result, the hold keeps its findings and records the cause: the running version and commit, and the count of main's flai commits since. The thread then names that cause and the remedy: publish with `flai release --pending`, install the release, and restart `flai serve`. It never tells the operator to fix files the running flai cannot judge.
  4. Left out: holding an acceptance that changes a format the running flai reads, the issue's last proposal. The plan's thread proposes it as a story of its own.
- Consequences: a hold like I-0382's either never starts, or names its real remedy and clears on the first look after the release is installed.

The number is the next free one when the task runs. If another story took 0309, use the next one and let `flai task done` correct the touches. This task waits for none: the decision is written here, so the code tasks can start beside it.

## Done when

- The ADR is accepted under `design/adrs/`, refines ADR-0234, and states the three rules above and what it leaves out.
- `design/adrs/README.md` lists it, and `flai check --strict` and the markdown lint pass on both files.

## Notes
