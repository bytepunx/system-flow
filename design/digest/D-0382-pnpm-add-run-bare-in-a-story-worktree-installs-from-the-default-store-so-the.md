---
id: D-0382
title: pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY
kind: cause
sources: [I-0177]
topics: [all]
paths: [design/tech/node.md, scripts/README.md, scripts/flaiover-install-test.sh, scripts/flaiover-install.sh, scripts/flaiover-pnpm.sh, scripts/flaiover-test.sh, design/issues/I-0177-pnpm-add-run-bare-in-a-story-worktree-installs-from-the-default-store-so-the-next-flaiover-install-sh-aborts-with-err-pnpm-aborted-remove-modules-dir-no-tty.md, design/issues/summary.md]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0382 pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY

I-0177 was closed on 2026-10-23.

The cause:

```text
pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY
```

The fix:

```text
Fixed by S-0481. scripts/flaiover-install.sh installs with --config.confirm-modules-purge=false, so an install linked from another store, such as a bare pnpm add makes from the default store, is purged and made again from env.sh's store instead of aborting with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY. scripts/flaiover-pnpm.sh add <package> adds a dependency with the pinned pnpm and the project's store. scripts/flaiover-install-test.sh, which the flaiover tier runs first, reproduces the store change.
```

Source: I-0177.
