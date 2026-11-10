---
id: T-3069
type: task
nature: improvement
title: This repository's OpenCode flai-guard plugin refuses a call when scripts/flai.sh exits 126 or 127 in a session flai started
status: done
parent: S-0666
owner: alex
created: 2026-11-10T10:35:59Z
updated: 2026-11-10T11:20:51Z
transitions:
  - to: ready
    at: 2026-11-10T11:15:49Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T11:15:49Z
    by: system-flow
  - to: done
    at: 2026-11-10T11:20:51Z
    by: agent-S-0666
stream: S-0666
tags: [guard, harness]
touches: [".opencode/plugins/flai-guard.js", flai/tests/integration/opencode_test.go]
after: [T-3061]
usage:
  source: log
  seconds: 302
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 4917
      cache_read: 1243656
      cache_write: 41021
      cost: 0.6344
---
# T-3069 This repository's OpenCode flai-guard plugin refuses a call when scripts/flai.sh exits 126 or 127 in a session flai started

## Work

This repository's `.opencode/plugins/flai-guard.js` runs `scripts/flai.sh guard --harness opencode`. It treats only a spawn error as a missing guard, and passes every exit but 2. So when `scripts/flai.sh` exits 127, as it will with neither `bin/flai` nor an installed flai after a failed build, the call passes, against what the Claude Code hooks do. Treat an exit of 126 or 127 as a missing guard: refuse the call in a session flai started (`FLAI_STARTED_BY`, `FLAI_STORY`, or `FLAI_ROLE` set), with `flai guard did not run (exit N)` and the output, and pass it with the one warning otherwise, as a spawn error does now. Update the plugin's header comment to match. Leave the template's plugin as it is: it runs `flai` from `PATH`, so a missing flai is a spawn error there already.

`TestFlaiGuardPluginUnderNode` in `flai/tests/integration/opencode_test.go` loads only the template's plugin. Load this repository's plugin too, with a fake `scripts/flai.sh` that exits 127, and check that a session flai started is refused and one it did not start passes. The test skips where `node` is not on `PATH`.

It waits for T-3061, whose ADR records the decision this code follows. It runs beside T-3066: the two share no path.

## Done when

- In a session flai started, the plugin refuses a call when `scripts/flai.sh` exits 126 or 127, naming the exit.
- The new case in `TestFlaiGuardPluginUnderNode` passes where `node` is on `PATH`.
- `flai test .opencode/plugins/flai-guard.js flai/tests/integration/opencode_test.go` passes.

## Notes
