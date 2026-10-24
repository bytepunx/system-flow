---
id: T-1828
type: task
nature: feature
title: The curator's definition lists its tools and rules, and settings.json runs flai guard on its file edits, in the template and here
status: done
parent: S-0431
owner: alex
created: 2026-10-16T06:06:08Z
updated: 2026-10-24T13:21:34Z
transitions:
  - to: ready
    at: 2026-10-23T23:42:58Z
    by: agent-S-0431
  - to: in-progress
    at: 2026-10-23T23:42:59Z
    by: agent-S-0431
  - to: done
    at: 2026-10-24T13:21:33Z
    by: agent-S-0431
stream: S-0431
tags: [cli, template]
touches: [".claude/agents/curator.md", template/root/.claude/agents/curator.md, ".claude/settings.json", template/root/.claude/settings.json, flai/cmd/guard_test.go, flai/internal/harness/harness_test.go, flai/internal/metrics/strategic_test.go]
after: [T-1823]
usage:
  source: log
  seconds: 5994
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 189
      output: 68763
      cache_read: 13478569
      cache_write: 444363
      cost: 6.7501
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
