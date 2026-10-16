---
id: T-1828
type: task
nature: feature
title: The curator's definition lists its tools and rules, and settings.json runs flai guard on its file edits, in the template and here
status: backlog
parent: S-0431
owner: alex
created: 2026-10-16T06:06:08Z
updated: 2026-10-16T06:06:08Z
transitions: []
stream: S-0431
tags: [cli, template]
touches: [".claude/agents/curator.md", template/root/.claude/agents/curator.md, ".claude/settings.json", template/root/.claude/settings.json]
after: [T-1826]
---
# T-1828 The curator's definition lists its tools and rules, and settings.json runs flai guard on its file edits, in the template and here

## Work

- Write `template/root/.claude/agents/curator.md` on the pattern of `analyzer.md`: prime with role `curate`, keep the digest through `flai digest` alone, never write a secret's value, edit no file, and its tools: `Read`, `Grep`, `Glob`, `Bash`, `Agent`, and flai's MCP tools `prime`, `inbox`, `board`, `item_get`, `doc_get`, `doc_search`, `thread_get`, `who_touches`, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events`; no `Edit`, `Write`, or `NotebookEdit`.
- Copy it to `.claude/agents/curator.md` here.
- In `template/root/.claude/settings.json` and `.claude/settings.json`, add `curate` to the roles for which the `Edit|MultiEdit|Write|NotebookEdit` entry runs `flai guard`.

Waits for T-1826, so the definition and the guard say the same.

## Done when

The four files say what the guard enforces, a template test or `flai upgrade`'s test sees the new definition, and `flai check --strict` reports nothing.

## Notes

Writes under `.claude/` here go through `permission_prompt`, which asks the operator on a thread; the story's agent makes them itself, not a sub-agent (ADR-0102), and waits for the answer.
