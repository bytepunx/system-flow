---
id: TH-0779
title: S-0506 criterion 3 still names flai serve settings, which does not exist, so the draft is not finalized
anchor:
  path: wip/kanban/stories/S-0506-flai-s-http-mcp-server-serves-an-agent-in-a-container-as-the-stdio-server-serves-one-on-the-host-every-tool-permission-prompt-under-auto-approve-for-an-agent-flai-serve-started-and-a-listener-a-container-reaches.md
  item: S-0506
status: resolved
participants: [orchestrator]
created: 2026-10-31T07:58:27Z
updated: 2026-10-31T08:01:17Z
---

# TH-0779 S-0506 criterion 3 still names flai serve settings, which does not exist, so the draft is not finalized

On wip/kanban/stories/S-0506-flai-s-http-mcp-server-serves-an-agent-in-a-container-as-the-stdio-server-serves-one-on-the-host-every-tool-permission-prompt-under-auto-approve-for-an-agent-flai-serve-started-and-a-listener-a-container-reaches.md.

## Entries

### 2026-10-31T07:58:27Z orchestrator
**Recommended:** reword criterion 3 to say "when the host setting `mcp.bind` says so, loopback by default: `flai serve mcp bind <address>` sets it, `flai serve mcp` shows it". I finalize S-0506 under finalize_drafts once it reads that way.

**What the edit of 2026-10-31T07:56Z changed.** It replaced the setting's name, `mcp.bind`, with `flai serve mcp`. The phrase TH-0664's confirmed answer asked to reword, "`flai serve settings` shows it", is still there.

**Why this blocks finalizing.**
- **`flai serve settings` does not exist.** `flai serve --help` lists no `settings` subcommand, and no task builds one. T-2152 assumes it exists ("show it in `flai serve settings`").
- **The setting is now named as a command.** Criterion 3 calls the setting `flai serve mcp`. T-2152 and T-2167 name it `mcp.bind` in `config.go` and `docs/operators/settings.md`. The agent would have to guess which one is meant.

**Alternative:** keep "`flai serve settings` shows it" and add a task that builds `flai serve settings`. That widens the scope, so it is your call.

Nothing else is missing. `flai promote --drafts` lists S-0506 as complete: forecast 1h30m, value 16.67 USD/week. It also waits on S-0505 by `after`.

### 2026-10-31T08:01:17Z orchestrator
Resolved: The concern is gone. The operator's edit of criterion 3 keeps `flai serve settings` as what shows the setting, and planner-S-0506 added T-2418, which builds that command, with its touches on the story. Criterion 3 can now be ticked against the diff.
