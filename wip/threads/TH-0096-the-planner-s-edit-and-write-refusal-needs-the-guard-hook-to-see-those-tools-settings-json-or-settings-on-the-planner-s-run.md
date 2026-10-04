---
id: TH-0096
title: "The planner's Edit and Write refusal needs the guard hook to see those tools: settings.json, or --settings on the planner's run?"
anchor:
  path: wip/kanban/stories/S-0208-the-planner-is-an-agent-flai-serve-starts-for-an-epic-or-a-story-behind-the-plan-host-action.md
  item: S-0208
status: resolved
participants: [agent-S-0208, alex]
created: 2026-10-04T00:50:10Z
updated: 2026-10-04T03:34:56Z
---

# TH-0096 The planner's Edit and Write refusal needs the guard hook to see those tools: settings.json, or --settings on the planner's run?

On wip/kanban/stories/S-0208-the-planner-is-an-agent-flai-serve-starts-for-an-epic-or-a-story-behind-the-plan-host-action.md.

## Entries

### 2026-10-04T00:50:10Z agent-S-0208
T-0785 taught `flai guard` to refuse the planner's Edit, Write, and NotebookEdit calls (with `FLAI_ROLE=plan`). Claude Code runs the hook only for tools its matcher names, and the matcher in `.claude/settings.json` and `template/root/.claude/settings.json` is `Bash|mcp__flai__.*`. Changing either file is refused by this session's permissions, so I have not changed them.

Recommendation: leave both settings files alone, and have the claude-code adapter pass, for a planner run only, `--settings` with a PreToolUse hook on `Edit|Write|NotebookEdit` that runs the run's flai `guard`. The refusal then reaches the planner in every project without a template change, and a story's agent pays no extra hook per edit. `planner.md` also leaves Edit, Write, and NotebookEdit out of its tools, so the session cannot call them in the first place.

Alternative: widen the matcher in both settings files to `Bash|Edit|Write|NotebookEdit|mcp__flai__.*`, which you would need to make or allow, plus a template release.

I am going ahead with the recommendation in T-0784; reply if you want the settings-file change instead.

### 2026-10-04T00:56:58Z alex
I want you to update the settings files as the planner will need the ability to read, create, and edit: epics, stories, tasks, notes, and threads (or it can't do its job). I grant you the permissions to change the files necessary.

### 2026-10-04T03:14:34Z agent-S-0208
Going ahead as T-0790. Both settings files get a PreToolUse entry on Edit|Write|NotebookEdit that runs the guard only when FLAI_ROLE=plan, so a story's agent pays no hook per edit (and here, no rebuild of bin/flai). The adapter drops its --settings. The guard lets the planner Edit and Write files under wip/kanban/ (epics, stories, tasks, and their notes) and wip/threads/, and refuses every other path, code included. planner.md lists Edit and Write, and its prompt says: edit an item's or thread's body by hand, and its front matter (status, links, touches) through flai, so moves stay held at backlog. Reply if you want any other path open to it.

### 2026-10-04T03:16:14Z agent-S-0208
Correction to my last entry: your planner.md says the planner writes items and threads "through flai only" and never edits files through the file system, and item_edit already changes an item's goal, criteria, notes, and body. So the guard stays as T-0785 made it: the planner's Edit, Write, and NotebookEdit are refused everywhere, and its flai writes (item_new, item_edit, thread_open, thread_reply, flai story/epic new, edit, touches, thread new/reply, move to backlog) pass. T-0790 only widens the settings and drops the adapter's --settings.

This session is still refused writes to both settings files (the permission prompt is not granted for .claude/), so please make this one change in the worktree .flai-cache/worktrees/S-0208, or grant write access under .claude/ there, and reply here. In each file add a second PreToolUse entry after the existing one (the existing entry is unchanged). It runs the guard only in a planner session, so a story's agent pays no hook per edit:

.claude/settings.json:
      {
        "matcher": "Edit|Write|NotebookEdit",
        "hooks": [
          {
            "type": "command",
            "command": "[ \"$FLAI_ROLE\" = plan ] || exit 0; out=$(\"$CLAUDE_PROJECT_DIR\"/scripts/flai.sh guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"
          }
        ]
      }

template/root/.claude/settings.json: the same, with the command
            "command": "[ \"$FLAI_ROLE\" = plan ] || exit 0; out=$(flai guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"

I go on with T-0787 and T-0788 meanwhile, and commit T-0790 once the files are there.

### 2026-10-04T03:21:57Z alex
I have made the changes to both files as requested

### 2026-10-04T03:22:11Z agent-S-0208
Thank you. One thing to fix before I commit: in both files the new entry replaced the existing `Bash|mcp__flai__.*` entry instead of joining it, so as they stand no sub-agent is guarded on Bash or flai's MCP tools, and the planner's flai calls are not held either. Please make each file hold both entries. Whole contents:

.claude/settings.json:
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash|mcp__flai__.*",
        "hooks": [
          {
            "type": "command",
            "command": "out=$(\"$CLAUDE_PROJECT_DIR\"/scripts/flai.sh guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"
          }
        ]
      },
      {
        "matcher": "Edit|Write|NotebookEdit",
        "hooks": [
          {
            "type": "command",
            "command": "[ \"$FLAI_ROLE\" = plan ] || exit 0; out=$(\"$CLAUDE_PROJECT_DIR\"/scripts/flai.sh guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"
          }
        ]
      }
    ]
  }
}

template/root/.claude/settings.json:
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash|mcp__flai__.*",
        "hooks": [
          {
            "type": "command",
            "command": "out=$(flai guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"
          }
        ]
      },
      {
        "matcher": "Edit|Write|NotebookEdit",
        "hooks": [
          {
            "type": "command",
            "command": "[ \"$FLAI_ROLE\" = plan ] || exit 0; out=$(flai guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"
          }
        ]
      }
    ]
  }
}

I will add a test that reads both files and checks both entries, so this cannot slip again.

### 2026-10-04T03:27:17Z alex
why is the template's settings file different from the one in `.claude/settings.json`?

### 2026-10-04T03:27:30Z agent-S-0208
Only the command each hook runs differs, and on purpose (delegation.md's project addition, ADR-0060). This repository builds flai itself, so its hooks run `"$CLAUDE_PROJECT_DIR"/scripts/flai.sh guard`: the guard built from this tree, so a story that changes the guard is held by the guard it is changing. A project made from the template has no flai source, so its hooks run the installed `flai guard` from PATH. Matchers, the FLAI_ROLE check, and the pass-only-refusals wrapper are the same in both, and the test I am adding checks that both files carry both entries.

The files still hold only the Edit|Write|NotebookEdit entry; the whole contents in my previous entry put back the Bash|mcp__flai__.* one beside it. Reply here once they are in, and I will commit them with T-0790.

### 2026-10-04T03:28:40Z alex
ok, I've made the edits

### 2026-10-04T03:34:56Z agent-S-0208
Resolved: Both settings files now carry the Bash|mcp__flai__.* entry and the planner-only Edit|Write|NotebookEdit entry; committed as T-0790 with a test that checks both files, and the adapter no longer passes --settings.
