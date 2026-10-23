---
id: I-0177
title: pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY
class: efficiency
status: closed
count: 1
cost: 15m
first_reported: 2026-10-21T08:25:41Z
last_reported: 2026-10-21T08:25:41Z
updated: 2026-10-23T07:20:18Z
---

# I-0177 pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY

## Description
pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY

## Instances

### 2026-10-21T08:25:41Z
Story: S-0458.
S-0458 added @lucide/svelte with a bare `pnpm add` in flaiover/ of its worktree. pnpm used ~/.local/share/pnpm/store, not the npm_config_store_dir scripts/env.sh sets (.flai-cache/pnpm-store), and left pnpm-lock.yaml newer than node_modules/.modules.yaml. Every later flai test and the close-out's flaiover tier then ran flaiover-install.sh, which asked pnpm to purge node_modules for the store change and aborted with no TTY. Sub-agents worked around it with touch node_modules/.modules.yaml; the close-out needed node_modules removed and reinstalled. A scripts/flaiover-add.sh that sources env.sh would add a dependency with the project's store.

## Remediation

Story S-0481 remediates this issue, created from it at 2026-10-21T12:00:00Z.
Closed 2026-10-23T07:20:18Z: Fixed by S-0481. scripts/flaiover-install.sh installs with --config.confirm-modules-purge=false, so an install linked from another store, such as a bare pnpm add makes from the default store, is purged and made again from env.sh's store instead of aborting with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY. scripts/flaiover-pnpm.sh add <package> adds a dependency with the pinned pnpm and the project's store. scripts/flaiover-install-test.sh, which the flaiover tier runs first, reproduces the store change.
