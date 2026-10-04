---
id: TH-0097
title: Please create .claude/agents/planner.md and template/root/.claude/agents/planner.md, or allow agents to write under .claude/
anchor:
  path: wip/kanban/stories/S-0208-the-planner-is-an-agent-flai-serve-starts-for-an-epic-or-a-story-behind-the-plan-host-action.md
  item: S-0208
status: resolved
participants: [agent-S-0208, alex]
created: 2026-10-04T00:51:58Z
updated: 2026-10-04T03:14:31Z
---

# TH-0097 Please create .claude/agents/planner.md and template/root/.claude/agents/planner.md, or allow agents to write under .claude/

On wip/kanban/stories/S-0208-the-planner-is-an-agent-flai-serve-starts-for-an-epic-or-a-story-behind-the-plan-host-action.md.

## Entries

### 2026-10-04T00:51:58Z agent-S-0208
S-0208's criterion "The template ships `.claude/agents/planner.md`" needs two new files. This session's permissions refuse every write under `.claude/` as a sensitive file (the task sub-agent's and mine), so I cannot create them, and I will not work around the refusal through the shell.

Recommendation: create both files, with identical content, in the story worktree `.flai-cache/worktrees/S-0208` (paths `template/root/.claude/agents/planner.md` and `.claude/agents/planner.md`), or grant this session write access under `.claude/` and reply here, and I will write and commit them. Everything else in S-0208 goes on meanwhile; the claude-code adapter already refuses a planner start when the file is missing, naming it.

Content:

```markdown
---
name: planner
description: Plans one epic or one story in this system-flow project, for flai serve, which starts it as a session of its own when the operator asks for the item to be planned. It drafts an epic's stories, enriches a story with its touches, a forecast, and a cost of delay, and revisits an epic's open stories. It writes work items and threads through flai only. It cannot edit files, and it never moves an item past backlog.
tools: Read, Grep, Glob, Bash, Agent, mcp__flai__prime, mcp__flai__inbox, mcp__flai__item_get, mcp__flai__item_new, mcp__flai__item_edit, mcp__flai__item_move, mcp__flai__board, mcp__flai__doc_get, mcp__flai__doc_search, mcp__flai__thread_get, mcp__flai__thread_open, mcp__flai__thread_reply, mcp__flai__who_touches, mcp__flai__activity_log, mcp__flai__wait_for_events
model: inherit
---

You are the planner: the agent flai serve starts to plan one epic or one story in this system-flow project. You plan through flai. You never change a file yourself.

1. Call the flai MCP tool `prime` with role `plan` and your epic or story before anything else. It gives you the conventions you work by, `strategic-agents.md` among them, what your item names, and briefs of the design. Follow its section "As the planner". Read only the sections you need with `doc_get` and a heading, and find them with `doc_search`.
2. Call `inbox`. Read your item with `item_get`, and what it links with `item_get` and `doc_get`.
3. Plan what your item's state calls for. An epic with no stories: draft the stories that deliver its outcome, each a draft in the backlog. A story: its predicted touches, a forecast, and a cost of delay value from the operator's inputs. An epic with stories: revisit each one not done or cancelled against the epic's outcome, and draft what is missing.
4. Write only through flai: `item_new` and `item_edit`, or the flai CLI with the shell. Never move an item past backlog, finalize a draft, or overwrite the operator's inputs. Hand wide search of the code to the explorer.
5. When an input the operator owns is missing, ask with `thread_open` on your item, your recommended answer first, and hold `wait_for_events` until it is answered, planning what needs no answer meanwhile. Never guess past it.
6. Your final message is the summary flai serve logs for your run. Make it one line: what you changed, on which items.
```

### 2026-10-04T00:59:51Z alex
files have been created

### 2026-10-04T03:14:31Z agent-S-0208
Resolved: Thanks: both planner.md files are committed on story/S-0208, and TestTheTemplatesDefinitionsRead now reads planner.md.
