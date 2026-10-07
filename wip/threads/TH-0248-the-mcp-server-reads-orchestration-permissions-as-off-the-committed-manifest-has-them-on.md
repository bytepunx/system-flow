---
id: TH-0248
title: "The MCP server reads orchestration.permissions as off; the committed manifest has them on"
anchor:
  path: system-flow.yaml
status: resolved
participants: [orchestrator, alex]
created: 2026-10-07T07:29:15Z
updated: 2026-10-07T08:26:21Z
---

# TH-0248 The MCP server reads orchestration.permissions as off; the committed manifest has them on

On system-flow.yaml.

## Entries

### 2026-10-07T07:29:15Z orchestrator
Recommended: restart this project's `flai mcp` server, the one `flai serve` runs, so that it rereads `system-flow.yaml`. Then I act on the permissions you set in 9c28ced8 (07:03Z).

Since that commit, the CLI reads `policy: throughput` and every permission on. The MCP server still answers as if there were no `orchestration` block: it orders by fifo, evaluates the release under judgement, and refuses these calls:

- `item_edit S-0265 draft:false` and `item_edit S-0279 draft:false`: "the orchestrator finalizes a draft only with orchestration.permissions.finalize_drafts, which is off"
- `item_move S-0308 ready`: "the orchestrator moves a story to ready only with orchestration.permissions.promote_to_ready, which is off"

The guard's verdict holds, so I have not finalized, promoted, ordered, or published anything. I have not worked around it through the CLI.

Second question. `orchestration.release` has `value: 1000` but no `policy`, so `flai release --evaluate` uses judgement: 9 accepted stories not yet released, worth 1063.5 USD/week. Recommended: add `policy: threshold` if you meant a threshold of 1000, which this batch meets. Under judgement I would publish the 9 too: each is a complete, accepted increment.

Meanwhile:

- S-0271 is in review. `flai accept --dry-run` is blocked only by TH-0247, the trial-merge conflict with S-0212 in `design/issues/I-0078-…md`. I will verify and accept it once that thread resolves.
- S-0265 and S-0279 are complete drafts I judge consistent, ready to finalize.
- S-0280, S-0287, S-0309, and S-0310 lack touches, a forecast, and a cost of delay value.

### 2026-10-07T07:33:56Z alex
the MCPs are restarted, please continue the review and approve

### 2026-10-07T08:15:03Z orchestrator
Recommended: restart the orchestrator, by turning the orchestrate host action off and on. Its own `flai mcp`, the stdio server its Claude session started, is still flai 1.34.0 with the old manifest. The restart did not reach it.

Since your reply, that server still:

- refuses `item_edit S-0265 draft:false` because finalize_drafts is off
- refuses `activity_log`, because the manifest needs flai 1.34.1
- reports `flai_outdated`: running 1.34.0, newest 1.34.2

The CLI on PATH has been upgraded and reads the permissions as on. I use it only where the CLI is the prescribed path, which is acceptance.

Where acceptance stands:

- S-0212: verified twice. Every criterion is met, and every tier passes at 94d8b7bb. The branch was 14 commits behind main, so `flai verify` stopped at its sync step. Its agent is rebasing now and stopped on the I-0078 conflict. I will verify again at the synced head. TH-0249 (answered, which its opener or you resolve) and my TH-0251 also block the acceptance.
- Release: nothing is pending since your publish.

Decisions I could not log while this server refuses `activity_log`: the S-0212 non-acceptances, the recommendation on TH-0249, and the refused finalize and promote calls.

### 2026-10-07T08:25:35Z alex
I've restarted the orchestrate host action, can you verify it's now flai 1.34.2?

### 2026-10-07T08:25:46Z orchestrator
Yes: the orchestrator's `flai mcp` (pid 563181, started 08:15Z) now runs `/home/alex/.flai/bin/flai`, which is flai 1.34.2 (commit 7772bd34). The binary is the current file, not a deleted one.

It reads the permissions as on. Since the restart it took these calls it refused before:

- finalized S-0265 and S-0279 (`item_edit draft:false`)
- promoted S-0308, S-0265, and S-0246 to ready
- logged each action with `activity_log`

One leftover: S-0213's agent still runs a `flai mcp` on a deleted, older binary (pid 79835). It picks up 1.34.2 when that agent restarts. Nothing for you to do unless S-0213's agent hits a refusal.

The release-policy question has moved to TH-0256. This thread is settled for me; resolve it when you are content.

### 2026-10-07T08:26:21Z alex
Resolved.
