---
id: TH-0409
title: "S-0364 is ready, but accepting it breaks the installed flai 1.40.1 until the host is upgraded: say when"
anchor:
  path: wip/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md
  item: S-0364
status: resolved
participants: [orchestrator, alex]
created: 2026-10-12T15:43:37Z
updated: 2026-10-12T15:56:20Z
---

# TH-0409 S-0364 is ready, but accepting it breaks the installed flai 1.40.1 until the host is upgraded: say when

On wip/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md.

## Entries

### 2026-10-12T15:43:37Z orchestrator
**Recommendation:** reply "go" when you can run `flai host upgrade` within a few minutes. I will then accept S-0364 and publish it with S-0232. Run `flai host upgrade` as soon as the release workflow finishes.

## Why I am holding it

S-0364 changes this repository's `vet` tier in `system-flow.yaml` to `{importers}`. The installed flai, 1.40.1, which runs the MCP server, rejects that placeholder. In S-0364's worktree:

```text
manifest.tests (error) system-flow.yaml:50: tests[1].command (tier "vet") argument {importers} is not a placeholder flai fills
```

Once that is on main, every installed-flai `flai check --strict` fails. That includes the check inside MCP `task_done`, and MCP `test` and `verify` refuse to run any tier. S-0337 and S-0366 are in progress, and their agents would stop at their next task close until the host runs a flai with S-0364. Close-outs are not affected, since `scripts/close-out.sh` builds flai from the story's tree.

A second reason: I cannot verify S-0364 at its head. MCP `verify` refuses on the same error.

- The stored verify passed at f91d7773.
- The head, eef5130b, adds only the close-out's own commit: an issue record (MD012 outside the story) and `design/issues/summary.md`.
- The dry-run at eef5130b has no blocker.

## What else is ready

- Pending release: S-0232 (2.63 USD a week), held for TH-0408, now resolved.
- S-0364 would add 25 USD a week.

If you would rather do it yourself: `scripts/flai.sh accept S-0364 --by alex`, then `flai release --pending`, then `flai host upgrade`.

### 2026-10-12T15:56:20Z alex
Resolved: S-0364 was accepted
