---
id: I-0178
title: A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness
class: defect
status: open
count: 1
cost: 20m
first_reported: 2026-10-21T08:42:54Z
last_reported: 2026-10-21T08:42:54Z
updated: 2026-10-21T12:00:09Z
---

# I-0178 A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness

## Description
A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness

## Instances

### 2026-10-21T08:42:54Z
Story: S-0456.
S-0458 added @lucide/svelte (748a07c5) for NavRail.svelte and was accepted while S-0456's close-out ran. The close-out synced, and its integration tier failed in TestEveryListedFlaioverComponentCompilesInTheHarness (flai/internal/components/flaiover_test.go:93): src/lib/components/NavRail.svelte does not compile for the harness: 500. A story worktree has no flaiover/node_modules, so flaioverModules (integration_test.go:972) links the main checkout's, which lacks @lucide/svelte because nothing installs it there after an acceptance. So main itself fails the test, and so does every story's close-out on this host. Running scripts/flaiover-install.sh --if-needed in S-0456's worktree, from the lockfile, cleared it for this story only. The merge queue's worktree presumably falls back to the main checkout's node_modules too.

## Remediation

Story S-0482 remediates this issue, created from it at 2026-10-21T12:00:09Z.
