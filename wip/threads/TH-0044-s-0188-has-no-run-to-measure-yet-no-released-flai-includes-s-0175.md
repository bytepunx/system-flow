---
id: TH-0044
title: "S-0188 has no run to measure yet: no released flai includes S-0175"
anchor:
  path: wip/kanban/stories/S-0188-measure-what-delegating-to-sub-agents-saves-on-two-stories-run-with-the-released-prompt.md
  item: S-0188
status: resolved
participants: [agent-S-0188, alex]
created: 2026-10-01T08:31:59Z
updated: 2026-10-01T09:25:17Z
---

# TH-0044 S-0188 has no run to measure yet: no released flai includes S-0175

On wip/kanban/stories/S-0188-measure-what-delegating-to-sub-agents-saves-on-two-stories-run-with-the-released-prompt.md.

## Entries

### 2026-10-01T08:31:59Z agent-S-0188
The criterion asks for two stories run by `flai serve` with a flai that includes S-0175. None exist yet:

- The installed flai is 1.26.3 (commit 6cd3686, built 2026-09-30T01:28Z). The newest tag is `flai/v1.26.3`. S-0175 was accepted at 08:30Z and is unpushed and unreleased (`inbox` lists it under `unpushed`).
- What S-0175 changed in flai is the claude-code prompt's delegation paragraph (`harness.delegation`). The agent definitions, `delegation.md`, and the guard hook are already on main, so a session `flai serve` starts now has the explorer and verifier and gets the convention through `prime`. It does not get the prompt's instruction to delegate.
- The kept logs since S-0175 are S-0179 and S-0180, both started at 08:18Z, before S-0175 was merged, and neither started a sub-agent. This run (S-0188) is the first with the definitions on main.

Recommended: release a flai that includes S-0175 (`flai push --pending`, then `flai release --pending` or Publish), install it (`flai self-upgrade`), and restart `flai serve`, because the running serve builds the prompt and keeps the old binary until it restarts. I then block S-0188 until two stories have run under it and delegated, and you unblock it, or answer here, when they have.

Alternative: count runs made with the definitions and convention on main but the 1.26.3 prompt, starting with the next two stories flai serve runs. That measures whether the convention alone makes agents delegate, not the released prompt, and the table would say so.

Which do you want?

### 2026-10-01T09:00:17Z alex
go with the recommendation

### 2026-10-01T09:00:44Z agent-S-0188
Thank you. No block was needed: flai 1.26.4, which includes S-0175, was installed and `flai serve` restarted at 08:39Z. Since then S-0185 (08:48Z) and S-0184 (08:54Z) have started with the delegation prompt, and S-0185 called the verifier and S-0184 the explorer. I'll measure them when their runs end, against non-delegating runs from the kept logs.

### 2026-10-01T09:25:17Z agent-S-0188
Resolved: flai 1.26.4 with S-0175 was released and installed. S-0184 and S-0185 ran on it and delegated, and they are measured in design/system/agent-context.md § Sub-agents › Measured.
