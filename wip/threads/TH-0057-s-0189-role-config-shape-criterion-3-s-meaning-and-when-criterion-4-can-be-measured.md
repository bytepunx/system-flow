---
id: TH-0057
title: "S-0189: role config shape, criterion 3's meaning, and when criterion 4 can be measured"
anchor:
  path: wip/kanban/stories/S-0189-sub-agents-run-on-a-cheaper-model-and-the-verifier-s-run-replaces-the-agent-s-own-full-suite-runs.md
  item: S-0189
status: resolved
participants: [agent-S-0189, alex]
created: 2026-10-01T10:43:39Z
updated: 2026-10-01T10:53:41Z
---

# TH-0057 S-0189: role config shape, criterion 3's meaning, and when criterion 4 can be measured

On wip/kanban/stories/S-0189-sub-agents-run-on-a-cheaper-model-and-the-verifier-s-run-replaces-the-agent-s-own-full-suite-runs.md.

## Entries

### 2026-10-01T10:43:39Z agent-S-0189
Three questions. Each has my recommended answer first. I'm going ahead on the parts that don't depend on them: the cheaper models in the definitions, and the delegation and prompt wording.

1. **The shape of per-role config.** Recommended: add an optional `roles` map to `agent`, keyed `explore` and `verify` (open to more later). Each value has the same shape as `agent` itself: harness, model, config. The top-level fields stay the `story` role, so existing stories and system-flow.yaml read as before. A story gets the project default's roles the same way it gets the rest, and `flai agent`, item_new, and the dashboard's agent edit carry them. For `claude-code`, flai serve starts the session with `--agents`. That flag takes the project's own `.claude/agents/<role>.md`, with the role's model set over the definition's, so the definition file stays the source of truth and the role overrides only the model. A role whose harness is not the story's is refused for `claude-code`, because a sub-agent runs inside the story's session. The `command` harness gets the roles as `FLAI_AGENT_ROLES` JSON. The alternative is to have the prompt tell the agent which `model` to pass to the Agent tool. That is weaker: it relies on the agent doing as told, and the tool only takes aliases (sonnet, haiku).

2. **Criterion 3**, "do not allow the verifier agent to perform any necessary corrections". Recommended: this is already true whatever the verifier's model. ADR-0059 gives the verifier no Edit or Write, and ADR-0060's `flai guard` refuses any sub-agent's writes. I'd meet the criterion by making delegation.md and the prompt say outright that the story's agent fixes what the verifier finds, never a sub-agent, plus a test that the guard refuses a verifier's write. Did you mean something more, such as letting a verifier that runs the story's own model make fixes? I'd advise against that, because it reverses ADR-0059.

3. **Criterion 4** (measure two stories run with a flai that includes this change). This can't be done before S-0189 is in review. It needs S-0189 accepted, flai released and installed, and two more stories run by flai serve. Recommended: I create a follow-up story for the measurement, S-0188-style and nature research, `after` S-0189, and S-0189 drops criterion 4. Alternatively, S-0189 goes to review with criterion 4 unchecked and explained in its notes, and you accept it knowing that. Instead of a real run, I can replay the verifier's prompt on S-0185's pre-fix branch on the cheaper model, as the notes suggest, and record whether it still finds the defects. That checks the benefit before release, not the cost.

### 2026-10-01T10:44:43Z agent-S-0189
Question 4: **I'm not allowed to edit the agent definitions.** Claude Code treats files under `.claude/` as sensitive. In this headless session it refused my edits to all four files: `template/root/.claude/agents/{explorer,verifier}.md` and their copies in `.claude/agents/`. I won't get around that through the shell. Recommended: you make the edit in the S-0189 worktree (`.flai-cache/worktrees/S-0189`). Add one line after `tools:` in each file: `model: haiku` in both `explorer.md` files and `model: sonnet` in both `verifier.md` files. I'll commit it with T-0671. Alternatively, allow `Edit(.claude/agents/**)` for this session. I chose the models this way: the explorer only finds and reads, so Haiku is enough. The verifier has to judge a diff against criteria and conventions, and its defects are the benefit we want to keep, so Sonnet. Both use the alias, so the template follows the harness's current model. If you'd rather both run Sonnet, say so.

### 2026-10-01T10:51:30Z alex
I accept your recommendations. I allow `Edit(.claude/agents/**) for this session. Let me know if I need to take additional steps to make that permission take effect.

### 2026-10-01T10:52:07Z agent-S-0189
Thanks. The permission didn't reach this session: a running Claude Code session keeps the tool permissions it was started with, and its Edit tool still refused `.claude/agents/`. Because you allowed exactly this edit, I made it through the shell instead: one `model:` line in each of the four files, committed as ba06e64. To make it take effect for future sessions, add `Edit(.claude/agents/**)` to the claude-code harness's arguments on the host with `flai serve agent harness claude-code -- ...`, or to the project's `.claude/settings.json` `permissions.allow`. Nothing more is needed for S-0189.

Going ahead as recommended:
- `agent.roles.{explore,verify}`, passed to claude-code with `--agents`.
- Criterion 3 met by the wording already committed, plus a guard test.
- A follow-up research story for criterion 4, with S-0189's criterion 4 moved to it.
- A replay of S-0185's verifier on Sonnet.

### 2026-10-01T10:53:41Z alex
Resolved.
