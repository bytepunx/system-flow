---
id: TH-0578
title: "Plan for S-0483: flai accept installs a sub-project's dependencies in the main checkout when its lock file changed"
anchor:
  path: wip/kanban/stories/S-0483-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md
  item: S-0483
status: open
participants: [planner-S-0483]
created: 2026-10-21T12:13:02Z
updated: 2026-10-21T12:13:02Z
---

# TH-0578 Plan for S-0483: flai accept installs a sub-project's dependencies in the main checkout when its lock file changed

On wip/kanban/stories/S-0483-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md.

## Entries

### 2026-10-21T12:13:02Z planner-S-0483
S-0483 is planned: five tasks in four layers, 1h forecast, cost of delay 37.50 USD a week.

## Proposed fix

`flai accept` runs a sub-project's `install` command in the main checkout when the accepted change touched that sub-project's lock file or `package.json`. `install` is a new manifest key. For flaiover it is `scripts/flaiover-install.sh --if-needed`. Every acceptance mode passes through `acceptItem`, so the step covers local, pull-request, and merge-queue (`--batch`) acceptances.

## Tasks and layers

| Layer | Task | After |
|-------|------|-------|
| 1 | T-1998 ADR for the `install` key and the acceptance step | none |
| 2 | T-2000 manifest key, `flai check` rule, settings row, flaiover's `install` in `system-flow.yaml` | T-1998 |
| 3 | T-2003 the step in `flai/cmd/accept.go`, with tests that reproduce I-0179 | T-2000 |
| 3 | T-2004 `docs/users/flai.md` and `design/system/devex.md` | T-2000 |
| 4 | T-2006 run the integration tier, close I-0179 | T-2003, T-2004 |

T-2003 and T-2004 share no path, so they run together.

## Assumptions

- The install runs the project's own command, not flai's built-in install by lock file (`components.Install`). The project's script sets the pnpm store from `scripts/env.sh`. pnpm run with a different store aborts when there is no TTY (I-0177).
- A failed install leaves the acceptance standing, because the merge is already made. flai reports the failure with the command's last lines. T-1998's ADR settles this, and the timeout.
- An install in the main checkout can briefly disturb a close-out that is linking that `node_modules` at the same moment. I-0179's 08:30Z instance is one example. S-0482's fix removes that link for story worktrees.

## Proposal: S-0482 and S-0483 overlap

S-0482 (I-0178) and S-0483 (I-0179) record the same incident: S-0458 added `@lucide/svelte`. S-0482's planner is planning, in parallel:

- T-1994: `scripts/integration.sh` installs flaiover's dependencies in a worktree.
- T-1996: `flaioverModules` refuses a stale `node_modules`.

That fixes the close-out side. S-0483 fixes the main checkout side. The two plans share no path except `design/issues/summary.md`.

Recommendation: keep both stories. If S-0482's fix is enough for you, cancel S-0483 instead and have S-0482 close I-0179 too. Its words are yours to change. I have not changed either story's words.
