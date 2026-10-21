---
id: I-0179
title: Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs
class: blocker
status: open
count: 1
cost: 15m
first_reported: 2026-10-21T08:44:04Z
last_reported: 2026-10-21T08:44:04Z
updated: 2026-10-21T12:00:17Z
---

# I-0179 Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs

## Description
Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs

## Instances

### 2026-10-21T08:44:04Z
Story: S-0463.
S-0458 added @lucide/svelte (748a07c5) for NavRail.svelte and was accepted at 08:32Z. TestEveryListedFlaioverComponentCompilesInTheHarness links the main checkout's flaiover/node_modules into its fixture, since a story worktree has none, and that folder has no @lucide/svelte. S-0463's close-out at 08:41Z stopped at the integration tier with flai/internal/components/flaiover_test.go:93 'src/lib/components/NavRail.svelte does not compile for the harness: 500'. S-0463 changes only flai/internal/host/host_test.go and design/issues. The round before, at 08:30Z on a main without S-0458, the same test failed with 'the setup src/routes/layout.css does not compile: 404'. That was while S-0458's close-out removed and reinstalled that node_modules (I-0177). Nothing installs flaiover's dependencies in the main checkout when an acceptance changes flaiover/package.json or pnpm-lock.yaml.

## Remediation

Story S-0483 remediates this issue, created from it at 2026-10-21T12:00:17Z.
