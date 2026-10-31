---
id: T-2352
type: task
nature: improvement
title: The planner's, orchestrator's, analyzer's, and curator's .claude/agents files, here and in the template, become what the ADR says, written by the operator
status: backlog
parent: S-0553
owner: alex
created: 2026-10-30T20:47:39Z
updated: 2026-10-30T20:47:39Z
transitions: []
stream: S-0553
tags: [template]
touches: [".claude/agents/planner.md", ".claude/agents/orchestrator.md", ".claude/agents/analyzer.md", ".claude/agents/curator.md", template/root/.claude/agents/planner.md, template/root/.claude/agents/orchestrator.md, template/root/.claude/agents/analyzer.md, template/root/.claude/agents/curator.md]
after: [T-2348, T-2351]
---
# T-2352 The planner's, orchestrator's, analyzer's, and curator's .claude/agents files, here and in the template, become what the ADR says, written by the operator

## Work

Make the four strategic files under `.claude/agents`, here and in `template/root`, what the ADR says. If Claude Code follows the link, they become links to the definitions `flai role definition --write` writes. Otherwise they are written from that command's output. `.claude/` is protected (ADR-0106), so ask the operator for the write through the permission prompt or a thread, with the whole contents or the link commands. Leave the explorer and the verifier alone.

Third layer: it waits for the ADR and for the command whose output it writes. It is a task of its own so the rest of the story can close while the write waits on the operator (D-0411).

## Done when

`flai check --strict` reports no `roles.definition` here, and a project made from the template has the four files as the ADR says. A session started by hand lists the four agents.

## Notes

Drafted by the planner of E-0024 on 2026-10-30, as S-0536 and S-0537 put their protected writes in tasks of their own.
