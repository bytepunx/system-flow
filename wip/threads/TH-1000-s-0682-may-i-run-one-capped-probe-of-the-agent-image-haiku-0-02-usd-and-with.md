---
id: TH-1000
title: "S-0682: may I run one capped probe of the agent image (haiku, 0.02 USD), and with which Claude login?"
anchor:
  path: wip/kanban/stories/S-0682-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s.md
  item: S-0682
status: answered
participants: [agent-S-0682, alex]
created: 2026-11-10T08:50:46Z
updated: 2026-11-10T08:53:31Z
---

# TH-1000 S-0682: may I run one capped probe of the agent image (haiku, 0.02 USD), and with which Claude login?

On wip/kanban/stories/S-0682-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s.md.

## Entries

### 2026-11-10T08:50:46Z agent-S-0682
S-0682 adds `scripts/agent-container.sh --probe S-nnnn`, which runs the agent image as `probe-S-nnnn` with a prompt that only reads and reports, and with no `FLAI_STORY`. T-2900 checks it on the rebuilt image.

The free check (`--probe S-0682 -- env`) needs no model. The capped trial calls a paid model.

Recommendation: one run of `FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-0682`, then I check that S-0682's narrative and items are unchanged.

The script mounts `FLAI_CLAUDE_LOGIN`, default `~/.claude`. There is no `~/.claude/.credentials.json` on this host. Logins exist in `~/.claude-me`, `~/.claude-rmh`, and `~/.claude-virago`.

Answer one of:

1. Yes, with `FLAI_CLAUDE_LOGIN=<folder>` (name it).
2. No: skip the paid trial; criterion 1 rests on the pilot test that reproduces I-0224 and the free check.

### 2026-11-10T08:53:31Z alex
1 - use ~/.claude-me
