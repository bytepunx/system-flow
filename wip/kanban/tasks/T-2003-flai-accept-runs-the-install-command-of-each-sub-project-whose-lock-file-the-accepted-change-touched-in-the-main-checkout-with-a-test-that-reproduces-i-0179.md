---
id: T-2003
type: task
nature: remediation
title: flai accept runs the install command of each sub-project whose lock file the accepted change touched, in the main checkout, with a test that reproduces I-0179
status: done
parent: S-0483
owner: alex
created: 2026-10-21T12:10:56Z
updated: 2026-10-21T13:27:42Z
transitions:
  - to: ready
    at: 2026-10-21T13:17:03Z
    by: agent-S-0483
  - to: in-progress
    at: 2026-10-21T13:17:04Z
    by: agent-S-0483
  - to: done
    at: 2026-10-21T13:27:41Z
    by: agent-S-0483
stream: S-0483
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_install.go, flai/cmd/accept_install_test.go, design/system/flai-cli.md]
after: [T-2000]
usage:
  source: log
  seconds: 637
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 59
      output: 16595
      cache_read: 2980183
      cache_write: 150213
      cost: 1.8374
---
# T-2003 flai accept runs the install command of each sub-project whose lock file the accepted change touched, in the main checkout, with a test that reproduces I-0179

## Work

Build the acceptance step T-1998's ADR decides, on the key T-2000 adds. It waits for T-2000 because it reads `Project.install`.

- In `flai/cmd/accept.go` `acceptItem`, note main's head before step 0, the merge, and after the merge list the paths the accepted change touched between the two heads. Put the step where every mode passes: local, pull-request (`acceptThroughRemote`), and `--batch` from the merge queue.
- In a new `flai/cmd/accept_install.go`, for each sub-project with `install` whose folder's `package.json` or lock file is among those paths, run its command from the main checkout's root, with the timeout the ADR sets. Report each run, and a failure with the command and its last lines, in the text output and under `--json`. A failure leaves the acceptance standing, as the ADR says.
- In a new `flai/cmd/accept_install_test.go`, reproduce I-0179 with the `gitIn` and `runIn` fixtures of `accept_fix_test.go`: a story whose branch changes `flaiover/pnpm-lock.yaml` is accepted, and the fixture's install command ran in the main checkout. A story that changes no lock file runs none. A failing command leaves the story done and is reported.

## Done when

- The three tests pass, and the existing accept tests pass unchanged.
- `flai test flai/cmd` passes.

## Notes
