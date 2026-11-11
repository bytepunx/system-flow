---
id: S-0706
type: story
nature: remediation
title: A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse
status: review
owner: alex
created: 2026-11-10T09:11:37Z
updated: 2026-11-11T09:28:03Z
transitions:
  - to: ready
    at: 2026-11-11T08:51:57Z
    by: alex
  - to: in-progress
    at: 2026-11-11T09:17:19Z
    by: agent-S-0706
  - to: review
    at: 2026-11-11T09:28:03Z
    by: agent-S-0706
tags: [flaiover, scripts, agent-container]
touches: [flaiover/pnpm-workspace.yaml, system-flow.yaml, scripts/flaiover-install.sh, scripts/flaiover-install-test.sh, scripts/README.md, design/tech/node.md, docs/operators/runbooks/agent-container.md, design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 663
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 4
      hand_edits: 1
      work: 46
  models:
    - model: claude-opus-5-5
      input: 110
      output: 28717
      cache_read: 11584131
      cache_write: 264416
      cost: 5.0069
  strategic:
    - kind: planner
      seconds: 52
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 37
          output: 5665
          cache_read: 1691229
          cache_write: 39748
          cost: 0.7334
    - kind: orchestrator
      seconds: 90
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 29
          output: 195
          cache_read: 1232226
          cache_write: 187505
          cost: 0.3964
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T09:11:37Z
  value: 22.5
  by: planner-S-0706
  at: 2026-11-10T10:31:25Z
forecast:
  duration: 45m
  delivery: 2026-11-11T10:23:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 3, behind S-0694, S-0698 and S-0702."
  by: flai
  at: 2026-11-11T08:58:19Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:46Z
---
# S-0706 A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse

## Goal

This story remediates [I-0259](../../../design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md), "A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0259 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0259 is closed with `flai issue close I-0259 --reason` saying what fixed it

## Tasks
- T-2993 Make pnpm install both the glibc and the musl native bindings for flaiover
- T-2994 Make flaiover-install.sh --if-needed install again when flaiover's install configuration is newer than the last install
- T-2997 Reproduce I-0259 in flaiover-install-test.sh: one install holds both libc variants, and a newer configuration installs again
- T-2999 Say in node.md and the agent-container runbook that flaiover's node_modules holds both libc variants, and close I-0259

## Notes

Cost of delay inputs set by flai from I-0259. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:22:35Z, 7.6 days before this story).

### Planning

The proposed remedy is the first one the issue's instance names: pnpm's `supportedArchitectures` with `libc: [glibc, musl]`. With it, one `flaiover/node_modules` holds both native bindings, and the host and the Alpine agent container can share it. `flaiover/pnpm-lock.yaml` already records both variants, so the lockfile is not expected to change. Keeping `node_modules` apart per platform would need a per-platform mount or path in `scripts/agent-container.sh` and `flai/internal/harness/container.go`. That is a larger change, and nothing in the instance asks for it.

The setting alone does not reach an existing install. `scripts/flaiover-install.sh --if-needed` compares only the lockfile's mtime with `node_modules/.modules.yaml`. flai accept runs flaiover's `install` only for a change to `package.json` or `pnpm-lock.yaml`. So T-2994 widens both: the script's check, and the `install` `paths` in `system-flow.yaml`.

Replanned on 2026-11-11 at the orchestrator's request. The only thread on the story, TH-1027, is resolved and holds no entry from the operator. The replan changed three things:

- The setting's home is `flaiover/pnpm-workspace.yaml`, which already holds `onlyBuiltDependencies`. `package.json` stays the fallback.
- The accept `install` `paths` were added to T-2994.
- The test's README row was added to T-2997.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flaiover/package.json` | declared; design: the fallback home of the setting, watched already by the accept `install` |
| `flaiover/pnpm-workspace.yaml` | layout: it already holds flaiover's pnpm settings (`onlyBuiltDependencies`), so `supportedArchitectures` goes there |
| `system-flow.yaml` | layout: flaiover's accept `install` `paths` name only `package.json` and `pnpm-lock.yaml`; the `install` key is written by hand, since `flai subproject` does not set it |
| `scripts/flaiover-install.sh` | declared; layout: its `--if-needed` stamp check skips a configuration change |
| `scripts/flaiover-install-test.sh` | declared; layout: the offline install test that reproduces I-0177, where I-0259's reproduction fits; `scripts/flaiover-test.sh` runs it |
| `scripts/README.md` | declared; layout: describes the install step and the test |
| `design/tech/node.md` | declared; design: the store and frozen-lockfile notes |
| `docs/operators/runbooks/agent-container.md` | declared; design: the mounts the worktree's `node_modules` is shared through |
| `design/issues/I-0259-…md`, `design/issues/summary.md` | declared by the criteria: `flai issue close` writes both |

`flai touches suggest` listed only files that change with nearly everything (`docs/users/flai.md` 5%, `design/system/flai-cli.md` 5%, and issue files from other stories). This story changes none of them, so none was added.

Forecast: flai gave 23m (median 135 s per unit of size over 59 done medium remediation stories on claude-opus-5-5, size 10: 2 criteria and 8 touches). It is raised to 45m. The test needs an offline pnpm fixture with libc-specific optional packages, and confirming the fix means one `node_modules` read from both the container and the host. Delivery is flai's 2026-11-11T10:31Z moved by the 22m added.

Cost of delay: 22.50 USD a week, flai's value from the operator's input of 9m lost per 168h cycle at 150 USD an hour. It stands. I-0259 has one instance, and its cost recurs only for a story whose flaiover tier runs both in the container and on the host.

No `after`: the story reads no output of another open story.
