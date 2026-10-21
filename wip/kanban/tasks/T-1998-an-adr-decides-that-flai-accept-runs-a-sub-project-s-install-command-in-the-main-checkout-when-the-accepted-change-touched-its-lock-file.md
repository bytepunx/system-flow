---
id: T-1998
type: task
nature: remediation
title: An ADR decides that flai accept runs a sub-project's install command in the main checkout when the accepted change touched its lock file
status: done
parent: S-0483
owner: alex
created: 2026-10-21T12:10:03Z
updated: 2026-10-21T13:06:56Z
transitions:
  - to: ready
    at: 2026-10-21T13:05:47Z
    by: agent-S-0483
  - to: in-progress
    at: 2026-10-21T13:05:49Z
    by: agent-S-0483
  - to: done
    at: 2026-10-21T13:06:56Z
    by: agent-S-0483
stream: S-0483
tags: [flai]
touches: [design/adrs/README.md, design/system/project-manifest.md, design/system/flai-cli.md, design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md]
usage:
  source: log
  seconds: 67
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 33
      output: 9386
      cache_read: 1685477
      cache_write: 84955
      cost: 1.0391
---
# T-1998 An ADR decides that flai accept runs a sub-project's install command in the main checkout when the accepted change touched its lock file

## Work

The story's goal asks for a fix proposed from I-0179's instances before it is built. This task is that proposal, recorded as an ADR with `flai adr new`, so the tasks after it build what it decides. It waits for nothing.

Read I-0179 and I-0178, and `flai/cmd/accept.go` `acceptItem`, which merges (local, pull-request, and `--batch` from the merge queue all pass through it) and then commits the acceptance in the main checkout. Nothing there installs a sub-project's dependencies, and no manifest key runs a command after an acceptance.

Decide, and write in the ADR:

- A sub-project's `install` in the manifest: a command, as a list like a test tier's `command`, run from the main checkout's root. For flaiover here: `[scripts/flaiover-install.sh, --if-needed]`.
- When `flai accept` runs it: after the merge, for each sub-project whose lock file (`pnpm-lock.yaml`, `package-lock.json`, `yarn.lock`) or `package.json` in its folder the accepted change touched, in every acceptance mode.
- What a failed install does: the acceptance stands, since the merge is made; flai says so with the command and its last lines, in the text output and in `--json`.
- Why a project's command and not flai's own install by lock file (`components.Install`): the project's script carries its environment, such as the pnpm store `scripts/env.sh` sets. pnpm run with another store aborts with no TTY (I-0177).
- How it relates to S-0482's fix, the integration tier installing in a story worktree: that frees close-outs from the main checkout's `node_modules`; this keeps the main checkout itself current for its own tests and the component dev server.

Add the ADR's row to `design/adrs/README.md`, and say what the key does in `design/system/project-manifest.md` and under the acceptance in `design/system/flai-cli.md`.

## Done when

- The ADR is accepted with `flai adr accept` and indexed in `design/adrs/README.md`.
- `design/system/project-manifest.md` and `design/system/flai-cli.md` describe `install` and the acceptance step, linking the ADR.
- `flai test` on the changed markdown passes.

## Notes
