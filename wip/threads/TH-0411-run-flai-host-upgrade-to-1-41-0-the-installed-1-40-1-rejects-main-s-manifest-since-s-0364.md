---
id: TH-0411
title: "Run flai host upgrade to 1.41.0: the installed 1.40.1 rejects main's manifest since S-0364"
anchor:
  path: wip/archive/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md
  item: S-0364
status: open
participants: [orchestrator]
created: 2026-10-12T15:57:27Z
updated: 2026-10-12T18:14:07Z
---

# TH-0411 Run flai host upgrade to 1.41.0: the installed 1.40.1 rejects main's manifest since S-0364

On wip/archive/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md.

## Entries

### 2026-10-12T15:57:27Z orchestrator
**Recommendation:** once the `release-flai.yml` run for `flai/v1.41.0` has finished, run `flai host upgrade`. Then reply here, and I resolve this thread.

## Why

- I published S-0232, S-0364, and S-0366 as flai 1.41.0, flaiover 0.42.4, and template 1.0.87.
- Main's `vet` tier now uses `{importers}`, from S-0364.
- The installed flai, 1.40.1, rejects that placeholder. Until the host runs 1.41.0, MCP `task_done`, `test`, and `verify` fail for S-0234, S-0313, S-0337, and S-0349.
- 1.41.0 is the first release signed in CI (S-0232). If the workflow fails at "Release key present", a secret is missing from the `release` environment; see `docs/operators/runbooks/release-key.md`.

### 2026-10-12T17:32:37Z orchestrator
Resolved: The host now runs flai 1.43.0, which reads main's {importers} manifest.

### 2026-10-12T17:37:26Z orchestrator
Correction to my resolve: the `flai` on PATH is 1.43.0, but the MCP server is still the old process. MCP `verify` still refuses main's `{importers}` tier with the 1.40.1 message, and so would MCP `task_done` and `test` for every story agent.

**What you need to do:** restart the MCP server, or `flai host`, so it runs 1.43.x. Until then I verify through the CLI.

### 2026-10-12T17:46:50Z orchestrator
**Urgent:** restart the MCP server, or `flai host`, now. `system-flow.yaml` requires flai 1.43.0 or newer, and the MCP server is still 1.40.1, so its MCP tools now refuse work. My `activity_log` just got:

```text
system-flow.yaml needs flai 1.43.0 or newer, and this is flai 1.40.1, which may not know the fields its items carry
```

Every story agent's MCP `task_done`, `item_move`, and `verify` are likely refused the same way. That affects S-0237, S-0368, and S-0373.

Meanwhile I published S-0233, S-0359, and S-0367 as flai 1.45.0, flaiover 0.43.0, and template 1.0.89, worth 44.22 USD a week.

### 2026-10-12T17:59:59Z orchestrator
**The flow has stalled.** Restart `flai host` so the running processes are the installed flai. The binary on PATH is 1.43.0, but flai serve and the MCP server are still the 1.40.1 processes. Since `system-flow.yaml` set `flai.minimum: 1.43.0`, they refuse work:

- No ready story has started. S-0373 has waited 29 minutes, and S-0363, S-0370, S-0351, and S-0236 are behind it, with nothing in progress.
- MCP `plan`, `activity_log`, and `verify` refuse with "needs flai 1.43.0 or newer, and this is flai 1.40.1". I cannot start the planners for S-0375 to S-0379, and my decision log has a gap from 17:46.

Unless you choose otherwise, run `flai host upgrade` to take 1.45.0, the newest published, then check `flai serve status`. Keep S-0237's compatibility break in mind: its release is not published yet (TH-0427), so 1.45.0 still talks to today's dashboard.

### 2026-10-12T18:14:07Z orchestrator
flai serve is running again, and S-0363, S-0373, S-0236, and S-0369 started at 18:13. Thank you.

One process is still on 1.40.1: the MCP server of this orchestrator session. Its `plan` and `activity_log` still refuse on `flai.minimum: 1.43.0`.

**What you need to do:** turn the orchestrate host action off and on again (`flai serve disable orchestrate`, then `flai serve enable orchestrate`) so I restart on the current flai. Until then I cannot start the planners for S-0375 to S-0379 or write my decision log. Accepting, publishing, promoting, and threads still work.
