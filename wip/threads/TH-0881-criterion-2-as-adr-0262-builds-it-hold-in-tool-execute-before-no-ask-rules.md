---
id: TH-0881
title: "Criterion 2 as ADR-0262 builds it: hold in tool.execute.before, no ask rules"
anchor:
  path: wip/kanban/stories/S-0408-an-opencode-agent-s-protected-write-is-asked-of-the-story-s-owner-through-flai-s-hold-and-ask-from-the-plugin-s-permission-asked-or-refused-by-deny-rules-under-auto-as-an-adr-settles.md
  item: S-0408
status: resolved
participants: [agent-S-0408, orchestrator, alex]
created: 2026-11-07T21:51:06Z
updated: 2026-11-10T09:32:31Z
---

# TH-0881 Criterion 2 as ADR-0262 builds it: hold in tool.execute.before, no ask rules

On wip/kanban/stories/S-0408-an-opencode-agent-s-protected-write-is-asked-of-the-story-s-owner-through-flai-s-hold-and-ask-from-the-plugin-s-permission-asked-or-refused-by-deny-rules-under-auto-as-an-adr-settles.md.

## Entries

### 2026-11-07T21:51:06Z agent-S-0408
**Recommended:** let me reword criterion 2 of S-0408 to what ADR-0262 builds, and tick it once the integration test shows it.

Criterion 2 was drafted before S-0406 read OpenCode. As written, its hold-and-ask branch cannot work under `opencode run` (agent-adapters.md § Read from the source at v1.19.5, finding 2):

- `opencode run` answers every `permission.asked` itself, at once, and does not await the plugin's event hook. A plugin that waits for the owner loses the race.
- An `ask` rule on a protected path is refused by `opencode run` even after the owner allowed.

ADR-0262 does this instead:

- The plugin's `tool.execute.before`, which OpenCode awaits, hands the call to `flai guard`. The guard holds it through `ask.Asker.Hold` for at most four minutes, on the story's permission thread as ADR-0124 says, and answers allow (exit 0) or deny (exit 2).
- With the plugin present, the adapter writes no rule on a protected path. Without it, it writes the deny rules. It never passes `--auto`.

Proposed wording:

> Under hold-and-ask: the plugin hands each call to `flai guard` on `tool.execute.before`, which holds a protected write through `ask.Asker.Hold` within the bound, opening or reusing the story's permission thread as ADR-0124 says, and passes or refuses it with the owner's answer; the configuration writes no `ask` rule, and denies the protected paths only in a project without the plugin.

Reply `yes` to reword it, or say what to keep. I am building it this way meanwhile.

### 2026-11-07T21:52:31Z orchestrator (recommendation)
**Recommended:** say yes. agent-S-0408 rewords criterion 2 to the proposed wording and ticks it once the integration test shows the hold.

**Why:**

- Criterion 2's hold-and-ask branch cannot be met as written. Finding 2 in `design/system/agent-adapters.md` § Read from the source at v1.19.5 (S-0406), on main, says three things. `opencode run` rejects every `permission.asked` at once. It does not await the plugin's `event` hook. The `permission.ask` hook is never triggered at v1.19.5.
- The same finding names `tool.execute.before` as the place flai can hold a call. OpenCode awaits it, and a throw refuses the call. ADR-0262 builds on that.
- The outcome stays the same: the owner is asked before a protected write, through ADR-0124's thread and retry. The goal already says "as an ADR settles". Only the mechanism in the criterion changes.
- The new wording fits ADR-0262 §1, §2, and §5: no `ask` rule, no `--auto`, and deny rules only without the plugin. Criterion 3 (`Ask: true`) still holds.

**Why I recommend rather than answer:** this changes the words of a finalized story's criterion. ADR-0262 was written in this story and is only on `story/S-0408`, so you have not reviewed it yet. That makes it your call.

**Alternative:** keep criterion 2 as written. I don't recommend it. Its hold-and-ask branch cannot be met, and its `--auto` branch is the one ADR-0262 rejects.

Source: design/system/agent-adapters.md § Read from the source at v1.19.5 (S-0406)

### 2026-11-10T09:32:31Z alex
Resolved.
