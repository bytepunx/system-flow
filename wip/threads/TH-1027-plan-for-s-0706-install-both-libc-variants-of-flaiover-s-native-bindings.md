---
id: TH-1027
title: "Plan for S-0706: install both libc variants of flaiover's native bindings"
anchor:
  path: wip/kanban/stories/S-0706-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md
  item: S-0706
status: open
participants: [planner-S-0706]
created: 2026-11-10T10:31:41Z
updated: 2026-11-10T10:31:41Z
---

# TH-1027 Plan for S-0706: install both libc variants of flaiover's native bindings

On wip/kanban/stories/S-0706-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md.

## Entries

### 2026-11-10T10:31:41Z planner-S-0706
The plan: one `flaiover/node_modules` that holds both the glibc and the musl native bindings, so that the host and the Alpine agent container can share it. This is pnpm's `supportedArchitectures` with `libc: [glibc, musl]`, the first remedy I-0259 names.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2993 Set `supportedArchitectures` for flaiover | `flaiover/package.json` | none |
| 1 | T-2994 `flaiover-install.sh --if-needed` installs again when the install configuration is newer than the stamp | `scripts/flaiover-install.sh`, `scripts/README.md` | none |
| 2 | T-2997 Offline reproduction in `flaiover-install-test.sh` | `scripts/flaiover-install-test.sh` | T-2993, T-2994 |
| 3 | T-2999 Document it in `design/tech/node.md` and the agent-container runbook, and close I-0259 | `design/tech/node.md`, `docs/operators/runbooks/agent-container.md`, the I-0259 file, `design/issues/summary.md` | T-2997 |

T-2993 and T-2994 share no path, so they run together.

## Figures

- Forecast: 45m, delivery 2026-11-10T17:21Z. flai gave 23m. I raised it for the offline pnpm fixture and the check from both sides.
- Cost of delay: 22.50 USD a week, flai's value from the 9m input. Unchanged.

## Assumptions

- `flaiover/pnpm-lock.yaml` already lists both variants of rolldown, rollup, lightningcss, and tailwind oxide, so the setting leaves it unchanged. T-2993 checks this with a frozen install.
- Without T-2994, the setting would never reach an existing install, because `--if-needed` compares only the lockfile's mtime. That covers the main checkout, open worktrees, and flai accept's install.
- I did not plan a separate `node_modules` per platform. It would change the container mounts in `scripts/agent-container.sh` and `flai/internal/harness/container.go`, and nothing in the instance needs it.
- If pnpm cannot resolve libc-specific packages offline, T-2997 falls back to testing the `--if-needed` rule and flaiover's declared configuration, and says why in its Notes.
- I found no earlier plan of S-0706 and no thread on it, so this is its first plan.

I would split, merge, or drop nothing.
