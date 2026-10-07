---
id: TH-0248
title: "The MCP server reads orchestration.permissions as off; the committed manifest has them on"
anchor:
  path: system-flow.yaml
status: answered
participants: [orchestrator, alex]
created: 2026-10-07T07:29:15Z
updated: 2026-10-07T07:33:56Z
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
