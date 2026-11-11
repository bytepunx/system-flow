---
id: S-0706
type: story
nature: remediation
title: A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse
status: backlog
owner: alex
created: 2026-11-10T09:11:37Z
updated: 2026-11-11T05:25:48Z
transitions: []
tags: [flaiover, scripts, agent-container]
touches: [flaiover/package.json, scripts/flaiover-install.sh, scripts/flaiover-install-test.sh, scripts/README.md, design/tech/node.md, docs/operators/runbooks/agent-container.md, design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 2
          cache_read: 102139
          cache_write: 491
          cost: 0.0281
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
  delivery: 2026-11-11T19:15:00Z
  basis: "Its own forecast of 45m; 22nd in the pull order with an in-progress limit of 3, behind S-0678, S-0686, S-0691, S-0693, S-0721, S-0696, S-0697, S-0700, S-0698, S-0702, S-0694, S-0705, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0689, S-0699, S-0701 and S-0703."
  by: flai
  at: 2026-11-11T05:25:48Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:46Z
---
# S-0706 A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse

## Goal

This story remediates [I-0259](../../../design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md), "A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0259 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0259 is closed with `flai issue close I-0259 --reason` saying what fixed it

## Tasks
- T-2993 Make pnpm install both the glibc and the musl native bindings for flaiover
- T-2994 Make flaiover-install.sh --if-needed install again when flaiover's install configuration is newer than the last install
- T-2997 Reproduce I-0259 in flaiover-install-test.sh: one install holds both libc variants, and a newer configuration installs again
- T-2999 Say in node.md and the agent-container runbook that flaiover's node_modules holds both libc variants, and close I-0259

## Notes

Cost of delay inputs set by flai from I-0259. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:22:35Z, 7.6 days before this story).

### Planning

The proposed remedy is the first one the issue's instance names: pnpm's `supportedArchitectures` with `libc: [glibc, musl]`. With it, one `flaiover/node_modules` holds both native bindings, and the host and the Alpine agent container can share it. `flaiover/pnpm-lock.yaml` already records both variants, so the lockfile is not expected to change. Keeping `node_modules` apart per platform would need a per-platform mount or path in `scripts/agent-container.sh` and `flai/internal/harness/container.go`. That is a larger change, and nothing in the instance asks for it. The setting alone does not reach an existing install: `scripts/flaiover-install.sh --if-needed` compares only the lockfile's mtime with `node_modules/.modules.yaml`. So T-2994 widens that check.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flaiover/package.json` | design: the `pnpm` field the pinned pnpm 10.34.5 reads `supportedArchitectures` from, watched already by the accept `install` in `system-flow.yaml` |
| `scripts/flaiover-install.sh` | layout: its `--if-needed` stamp check skips a configuration change |
| `scripts/flaiover-install-test.sh` | layout: the offline install test that reproduces I-0177, where I-0259's reproduction fits |
| `scripts/README.md` | layout: describes the install step |
| `design/tech/node.md` | design: the store and frozen-lockfile notes |
| `docs/operators/runbooks/agent-container.md` | design: the mounts the worktree's `node_modules` is shared through |
| `design/issues/I-0259-…md`, `design/issues/summary.md` | declared by the criteria: `flai issue close` writes both |

`flai touches suggest` listed only files that change with nearly everything (`design/system/flai-cli.md` 5%, `docs/users/flai.md` 5%, and issue files from other stories). This story changes none of them, so none was added. If T-2993 finds `flaiover/pnpm-workspace.yaml` the better home, that file and `system-flow.yaml` join its touches, as its Work says.

Forecast: flai gave 23m (median 134 s per unit of size over 54 done medium remediation stories on claude-opus-5-5, size 10: 2 criteria and 8 touches). It is raised to 45m. The test needs an offline pnpm fixture with libc-specific optional packages, and confirming the fix means one `node_modules` read from both the container and the host. Delivery is flai's 2026-11-10T16:59Z moved by the 22m added.

Cost of delay: 22.50 USD a week, flai's value from the operator's input of 9m lost per 168h cycle at 150 USD an hour. It stands. I-0259 has one instance, and its cost recurs only for a story whose flaiover tier runs both in the container and on the host.

No `after`: the story reads no output of another open story.
