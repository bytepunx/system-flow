---
id: T-1117
type: task
nature: improvement
title: The story agent's start prompt, session-start.md, work-management.md, CLAUDE.md, and the template's copies begin the loop with flai story start
status: done
parent: S-0274
owner: alex
created: 2026-10-06T22:53:52Z
updated: 2026-10-07T08:06:10Z
transitions:
  - to: ready
    at: 2026-10-07T08:00:04Z
    by: agent-S-0274
  - to: in-progress
    at: 2026-10-07T08:00:05Z
    by: agent-S-0274
  - to: done
    at: 2026-10-07T08:06:10Z
    by: agent-S-0274
stream: S-0274
tags: [conventions, template, go]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/session-start.md, design/conventions/work-management.md, template/root/design/conventions/session-start.md, template/root/design/conventions/work-management.md, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/CHANGELOG.md]
after: [T-1109, T-1111]
usage:
  source: log
  seconds: 365
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 80
      output: 28605
      cache_read: 4128681
      cache_write: 131934
      cost: 2.1311
---
# T-1117 The story agent's start prompt, session-start.md, work-management.md, CLAUDE.md, and the template's copies begin the loop with flai story start

## Work

Make the first step of a story's loop one call wherever it is written down:

- **The prompt.** In `Prompt` in `flai/internal/harness/harness.go`, the story agent's start prompt tells it to prime with `flai prime --story` and then to open the story with `flai stream open`. It now tells the agent to begin with the MCP tool `story_start` (or `flai story start S-nnnn`), which moves the story to in-progress, opens the stream, primes, and answers the inbox, and to work in the worktree it answers. The reconcile prompt for a story already in progress with no agent on this host keeps `flai stream open`, because `story start` refuses a story that is not ready.
- **The conventions.** The pulling step in `design/conventions/session-start.md` and the ordered list in `design/conventions/work-management.md` ("When you pull a story, in this order") become `flai story start` (`story_start`) followed by writing or reviewing the tasks. Change their `template/root` copies the same way, the baseline above the marker, as one template change.
- **The map.** Step 3 of "Prime your session" in `CLAUDE.md` and `template/root/CLAUDE.md.tmpl` pulls with `flai story start`.
- **The changelog.** Add an entry to `template/CHANGELOG.md` naming the convention change.

It waits for T-1109 and T-1111, the command and the tool it names.

## Done when

- `harness_test.go` asserts that the story prompt names `story_start` and `flai story start`, and that the reconcile prompt still names `flai stream open`.
- The conventions, their template copies, `CLAUDE.md`, and `CLAUDE.md.tmpl` name `flai story start` as the first step of pulling a story, and the template's copies match the project's above the marker.
- `flai check --strict`, the markdown lint, and `scripts/flai-test.sh` pass.

## Notes
