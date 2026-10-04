---
id: T-0818
type: task
nature: feature
title: "The planner enriches a story with the three commands and records its reasoning under ### Planning"
status: done
parent: S-0210
owner: alex
created: 2026-10-04T19:46:00Z
updated: 2026-10-04T20:10:15Z
transitions:
  - to: ready
    at: 2026-10-04T19:59:36Z
    by: agent-S-0210
  - to: in-progress
    at: 2026-10-04T19:59:36Z
    by: agent-S-0210
  - to: done
    at: 2026-10-04T20:10:15Z
    by: agent-S-0210
stream: S-0210
tags: [flai]
touches: [flai/internal/guard, flai/internal/harness, template/root/.claude/agents/planner.md, ".claude/agents/planner.md", design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml]
after: [T-0815, T-0816]
usage:
  source: log
  seconds: 639
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 85
      output: 32364
      cache_read: 4208590
      cache_write: 120410
      cost: 2.2001
---
# T-0818 The planner enriches a story with the three commands and records its reasoning under ### Planning

## Work

Let the planner and every sub-agent run the three commands: `forecast`, `cod`, and `touches suggest` join the guard's reads. Tell the planner to use them: `planPrompt` and `planner.md` say to run `flai touches suggest` and keep any touch the story already declares, run `flai forecast` and `flai cod`, review each figure, adjust it with a stated reason, write `touches`, `forecast`, and the cost of delay `value` through flai, and record which touch came from where and why each figure stands under a `### Planning` heading in the story's Notes, which it owns and rewrites, leaving the rest of Notes as it was. The convention's "As the planner" names the commands; the template gets the same, a changelog entry, and a version bump.

Waits for T-0815 and T-0816: the prompt names their commands and flags.

## Done when

- [ ] The guard passes `flai forecast`, `flai cod`, and `flai touches suggest` for a sub-agent and the planner, and still refuses `flai touches S-nnnn path` to a sub-agent; tests pin it
- [ ] The planner's prompt test pins the three commands and the `### Planning` heading
- [ ] `template/root/.claude/agents/planner.md` and `.claude/agents/planner.md` match, the convention and its template copy match, and the template's version and changelog record the change

## Notes
