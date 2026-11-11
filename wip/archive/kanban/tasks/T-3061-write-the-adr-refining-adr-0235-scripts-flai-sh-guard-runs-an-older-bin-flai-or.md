---
id: T-3061
type: task
nature: improvement
title: "Write the ADR refining ADR-0235: scripts/flai.sh guard runs an older bin/flai or the installed flai when the build fails, and exits 127 with neither"
status: done
parent: S-0666
owner: alex
created: 2026-11-10T10:35:32Z
updated: 2026-11-10T11:15:41Z
transitions:
  - to: ready
    at: 2026-11-10T11:15:14Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T11:15:14Z
    by: system-flow
  - to: done
    at: 2026-11-10T11:15:41Z
    by: agent-S-0666
stream: S-0666
tags: [guard, adr]
touches: [design/adrs/README.md, design/adrs/0280-scripts-flai-sh-guard-runs-an-older-bin-flai-or-else-the-installed-flai-when.md]
usage:
  source: log
  seconds: 27
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1620
      cache_read: 749616
      cache_write: 20994
      cost: 0.3383
---
# T-3061 Write the ADR refining ADR-0235: scripts/flai.sh guard runs an older bin/flai or the installed flai when the build fails, and exits 127 with neither

## Work

Write a new ADR with `flai adr new`, refining ADR-0235, which says a failed build of this tree fails open and that `scripts/flai.sh` never runs an older installed flai where `go` is present. The ADR records the fix S-0666 proposes from I-0275's instance.

- Context: I-0275. `scripts/flai.sh` runs `scripts/flai-build.sh` when `go` is present and `bin/flai` is stale. A failed `go build` exits 1 under `set -eu`, and the hooks in `.claude/settings.json` allow every exit but 2, 126, and 127, so the call passes unguarded. This repository's `.opencode/plugins/flai-guard.js` treats only a spawn error as missing, so it passes an exit 127 as well.
- Decision: under `guard`, when the build fails, `scripts/flai.sh` runs the existing `bin/flai`, older than the sources, or else the installed `flai`, each with a warning on standard error, and exits 127 when neither is there, so the hook refuses. Every other command still fails on a failed build. The hooks' exit table stays as ADR-0235 has it. The OpenCode plugin refuses on 126 and 127 in a session flai started, as the Claude Code hooks do.
- Alternatives: fail closed on every non-zero exit, which ADR-0235 rejected; fall back for every command, which hides a broken build from the agent working on flai.

The ADR's file is new, and its number is the next free one when you write it; the story's `design/adrs` touch covers it. Nothing waits before this task: the code and the docs follow the decision it records.

## Done when

- A new ADR refining ADR-0235 states the decision above, with its context, consequences, and alternatives, and `design/adrs/README.md` lists it.
- `flai check --strict` reports nothing new on it.

## Notes
