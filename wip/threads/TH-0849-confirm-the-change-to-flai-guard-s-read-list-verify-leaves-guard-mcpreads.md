---
id: TH-0849
title: "Confirm the change to flai guard's read list: verify leaves guard.MCPReads"
anchor:
  path: wip/kanban/stories/S-0556-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md
  item: S-0556
status: resolved
participants: [agent-S-0556, orchestrator, alex]
created: 2026-11-07T06:50:26Z
updated: 2026-11-07T06:52:19Z
---

# TH-0849 Confirm the change to flai guard's read list: verify leaves guard.MCPReads

On wip/kanban/stories/S-0556-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md.

## Entries

### 2026-11-07T06:50:26Z agent-S-0556
S-0556's notes ask for your confirmation before the close-out. The change is committed on `story/S-0556`, and nothing has been closed out yet.

**The question.** Do you confirm that `verify` leaves `guard.MCPReads`, and that `flai verify --job` is refused to sub-agents? I recommend yes.

**What changes:**

| Caller | MCP `verify` | `flai verify --job` | `flai verify --last`, plain `flai verify` |
|--------|--------------|---------------------|-------------------------------------------|
| Story's agent | passes | passes | passes |
| Orchestrator | passes, named in its role file and `guard.MCPOrchestrates` | not named; its prompt does not need it | passes |
| Sub-agents | refused | refused | passes |
| Planner, analyzer, curator | refused | refused | `--last` passes |

**Why.** The tool asks flai serve for a close-out job that commits in the story's worktree (ADR-0229), so any sub-agent could start a committing run. The convention only asked them not to. D-0446 anticipated this narrowing and asked to keep two calls. The orchestrator's verify of a story in review stays through the tool. A reader's look at the record stays through `flai verify --last`.

**The one cost.** A reader that called the tool at a verified commit, to get the record, now gets a refusal and uses `flai verify --last` instead.

**Upgrade.** A project whose `orchestrator.md` predates this gets `roles.definition` from `flai check`, and its orchestrator is refused the tool, until `flai upgrade` replaces the baseline. `template/CHANGELOG.md` says so.

**What each answer leads to:**

- **Yes:** I ask for the close-out, review the diff, and move S-0556 to review.
- **No, or a change:** I rework the guard, its tests, and the documents as you say, and reopen I-0240's reason to match.

The tiers for the changed paths passed. The whole suite has not run. Criteria 1 and 2 are ticked.

### 2026-11-07T06:51:36Z orchestrator (recommendation)
**Recommended: yes.** Confirm that `verify` leaves `guard.MCPReads` and that sub-agents are refused `flai verify --job`. agent-S-0556 can then ask for the close-out and move S-0556 to review.

**Why it is yours:** S-0556's Notes ask you to confirm this before the close-out. The change also narrows what S-0546 chose, which D-0446 records. So I recommend and do not answer.

**Why yes:** I checked the diff on `story/S-0556` against D-0446.

- **D-0446 expects this change.** It says S-0556 narrows the guard and must keep two calls. Both stay:
  - The orchestrator verifying a story in review. `verify` is now in `MCPOrchestrates` in `flai/internal/guard/guard.go`, and the three copies of `orchestrator.md` name it. The orchestrator's plain `flai verify <S-nnnn> --json` still passes as a read.
  - A reader asking at a commit the record already holds. That goes through `flai verify --last` now, which runs nothing.
- **The gap is real.** The only change to `flagReads` is `--job` joining `--record-issues` among `verify`'s writing flags. So no sub-agent can start the committing job of ADR-0229 any more. Until now only the delegation convention asked them not to.
- **The one cost is small.** A reader that used the tool at a verified commit now uses `flai verify --last`.

On TH-0827 I approved the same plan, and I stand behind it as the answer.

Source: design/digest/D-0446-s-0546-kept-the-mcp-verify-tool-among-flai-guard-s-reads-since-a-call-at-a.md

### 2026-11-07T06:51:54Z alex
Confirmed the recommendation of 2026-11-07T06:51:36Z orchestrator.

Source: design/digest/D-0446-s-0546-kept-the-mcp-verify-tool-among-flai-guard-s-reads-since-a-call-at-a.md

### 2026-11-07T06:52:19Z agent-S-0556
Resolved: The operator confirmed the orchestrator's recommendation at 2026-11-07T06:51:54Z: verify leaves guard.MCPReads and sub-agents are refused flai verify --job, as built on story/S-0556. The close-out follows.
