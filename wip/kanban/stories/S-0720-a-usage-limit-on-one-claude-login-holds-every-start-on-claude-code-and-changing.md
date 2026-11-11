---
id: S-0720
type: story
nature: remediation
title: A usage limit on one Claude login holds every start on claude-code, and changing the harness's program or login does not lift it
status: backlog
owner: alex
created: 2026-11-11T05:05:59Z
updated: 2026-11-11T05:05:59Z
transitions: []
tags: [serve, budget]
touches: [flai/internal/serve/budget.go, flai/cmd/serve_agent.go, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---

# S-0720 A usage limit on one Claude login holds every start on claude-code, and changing the harness's program or login does not lift it

## Goal

When the operator changes what a harness runs as (its program, its arguments, or the container's login), a usage limit recorded for the old login stops holding starts, and a limit on one login never holds a project that runs on another. The operator can also clear a budget state by hand.

## Acceptance criteria
- [ ] `flai serve agent harness <name>` (the dashboard's harness setting included), when it changes the program, arguments, or login a project's or the host's harness runs with, ends the hold on that harness's budget state for what it changed, so the next start goes as the probe; the change is journaled under `serve.budget` with why, and sent as a notice like any other clear.
- [ ] A budget state is kept under what tells one login from another (the provider, else the harness with its resolved program and login), so a run of `rmhclod` ending on the limit does not hold a project whose `claude-code` program is `meclod`; a `budget.json` in the old shape is still read.
- [ ] `flai serve budget clear [<provider>]` clears a state by hand and `flai serve budget` shows it; the hold's message (`budgetHeldWhy`) and the dashboard name that command.
- [ ] Tests cover: a harness change lifting the hold, two programs on one harness held apart, and the manual clear; `design/system/flai-cli.md` and the operator docs describe all three.

## Tasks

## Notes

What happened, 2026-11-11 (UTC): S-0692's run on `rmhclod` (`CLAUDE_CONFIG_DIR=~/.claude-rmh`) ended on `usage_limit_reached` (429) at 02:33:30Z, and flai serve set `~/.flai/serve/budget.json` key `claude-code`, held until the API's reset, 2026-11-13T22:00:00Z. At 02:35:06Z the operator set system-flow's `claude-code` program to `meclod` (another account) from the dashboard (`settings.harness`). The hold stayed: S-0690, S-0692, and S-0693 were logged "agents not started: the budget state holds them" through a flai serve restart at 05:01Z. Removing `budget.json` by hand at 05:04Z let all three start on `meclod` at the next pass.

Why: `budgetKey` (`flai/internal/serve/budget.go`) is the run's provider, or else the harness's name, `claude-code`. The program and login are not part of it, so every project on the host shares one state whichever account its program logs in with (signet was on `meclod` already and was held too). The state clears only when a probe ends clean, and a probe only starts after `until`; nothing else clears it, and there is no command for it. The error text in `Dir.budgets` tells the operator to remove the file.

Still true after the manual clear: the host default and authstar run `rmhclod`, so their next run ends on the limit and sets the host-wide hold again, system-flow and signet included. That is the second criterion.
