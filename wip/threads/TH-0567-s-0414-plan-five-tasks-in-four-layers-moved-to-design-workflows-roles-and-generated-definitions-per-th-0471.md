---
id: TH-0567
title: "S-0414 plan: five tasks in four layers, moved to design/workflows/roles and generated definitions per TH-0471"
anchor:
  path: wip/kanban/stories/S-0414-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-that-flai-reads-and-flai-check-validates-and-the-template-ships-the-planner-the-orchestrator-and-the-analyzer-as-baseline-role-files.md
  item: S-0414
status: answered
participants: [planner-S-0414, orchestrator, alex]
created: 2026-10-21T11:11:14Z
updated: 2026-10-21T12:44:16Z
---

# TH-0567 S-0414 plan: five tasks in four layers, moved to design/workflows/roles and generated definitions per TH-0471

On wip/kanban/stories/S-0414-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-that-flai-reads-and-flai-check-validates-and-the-template-ships-the-planner-the-orchestrator-and-the-analyzer-as-baseline-role-files.md.

## Entries

### 2026-10-21T11:11:14Z planner-S-0414
S-0414 is replanned on your TH-0471 answers. It keeps its five tasks and four layers, and now has a forecast and a cost of delay.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-1720 | none | The ADR: the role file, `design/workflows/`, the definition fields, and the named conditions |
| 2 | T-1721 | T-1720 | `flai/internal/roles`: the loader and the validation |
| 2 | T-1722 | T-1720 | The three baseline role files in both places, and the upgrade test |
| 3 | T-1723 | T-1721, T-1722 | `roles.*` in `flai check`, and `flai role list` and `flai role show` |
| 4 | T-1724 | T-1723 | The design, the convention, the docs, and the layout tables |

T-1721 and T-1722 run together; they share no path.

## What changed

- **Folder (answer 2).** The title, goal, criteria, and every touch now say `design/workflows/roles/`. The ADR also records `design/workflows/`, since `repository-layout.md` adds no `design/` folder without one. Four layout tables were added to T-1724: `design/README.md`, `template/root/design/README.md.tmpl`, `CLAUDE.md`, and `template/root/CLAUDE.md.tmpl`.
- **Definitions (answer 6).** Criteria 1, 3, and 6 now say flai generates a role's harness definition from its role file. In this story the role files carry the definition fields and `roles.definition` checks them. Nothing is generated yet.
- **Forecast.** 1h17m, flai's figure kept: size 43 at 107 s a unit. Delivery 2026-10-21T15:40Z.
- **Cost of delay.** 66.76 USD a week. This is S-0414's share, by forecast duration, of the 600 USD a week your 4h gives E-0024. I worked it out by hand: `flai cod` refuses while E-0024 has no inputs.

## Questions, recommendation first

1. **Task words.** The tasks predate your answers. I added a Notes entry to each and left their Work and Done when alone, as I may not rewrite them unasked.
   - **Recommended:** let me rewrite them to the new folder and generated definitions: T-1720's first Work bullet, T-1722's title and first line, T-1723's `roles.definition` bullet and its Done when, and T-1724's layout bullet.
   - Alternative: leave the Notes entries to the story's agent.
2. **Where generation lands.** Answer 6 is bigger than this story.
   - **Recommended:** S-0419, which already moves the built-in prompts into role files. It would generate the definitions and check that they match today's `.claude/agents` files. `subAgents` in `flai/internal/harness/adapters.go` already passes definitions to Claude Code as `--agents` JSON. A session flai serve starts could take the generated definition that way, with no write under `.claude/`. A symlink would be needed only for sessions started by hand, and only if Claude Code follows one.
   - Alternative: a new story right after S-0414.
3. **E-0024's input.** Your 4h is not on the epic yet.
   - **Recommended:** run `flai edit E-0024 --time-lost-per-cycle 4h`. `flai cod` then gives every E-0024 story its share. For S-0414 that matches the 66.76 above.

## Assumptions

- The role files stay data only. Nothing changes in behaviour until S-0415 to S-0419 read them.
- S-0416, S-0419, S-0421, and S-0423 still name `design/roles` or `design/formulas`. They need E-0024 planned again; that is outside this run.

### 2026-10-21T11:11:30Z orchestrator
Approved under plan_backlog_stories: the five tasks in four layers, the touches moved to `design/workflows/roles/`, and the 1h17m forecast.

**Question 1:** yes, rewrite the tasks to match the new folder and the generated definitions:

- T-1720's first Work bullet;
- T-1722's title and first line;
- T-1723's `roles.definition` bullet and its Done when;
- T-1724's layout bullet.

This applies the operator's recorded answers on TH-0471 (2 and 6) and adds no scope. You are asking on the plan thread, as the convention requires before you rewrite a task you did not write.

Questions 2 and 3 follow as a recommendation for the operator. The thread stays open for them.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T11:11:33Z orchestrator (recommendation)
Recommended for the operator: take the planner's answers to questions 2 and 3.

- **Question 2: generation lands in S-0419.** That story already moves the built-in prompts into role files, and it would generate the definitions and check them against today's `.claude/agents`. flai serve can pass a generated definition through `--agents`, so it needs no write under `.claude/`. A symlink is then needed only for sessions started by hand. This changes S-0419's scope, so it is yours to decide. The alternative is a new story right after S-0414.
- **Question 3: write your TH-0471 figure on the epic.** Run `flai edit E-0024 --time-lost-per-cycle 4h`. `flai cod` then values every E-0024 story from it, in place of the hand-worked 66.76 on S-0414. Setting an epic's input is yours: my permission covers a backlog story's inputs only, and a story's input would stop the epic from taking yours.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-21T11:12:52Z planner-S-0414
The four tasks are rewritten as approved. T-1720's first Work bullet now names `design/workflows/roles/` and the ADR for `design/workflows/`. T-1722's title and first line name the new folder. T-1723's `roles.definition` checks the role file's own definition fields, and its Done when names `design/workflows/roles/`. T-1724's layout bullet covers the tree, the contract, and the four layout tables. Layers and `after` are unchanged. S-0414's Planning notes record the rewrite. Questions 2 and 3 are still the operator's to answer.

### 2026-10-21T12:44:16Z alex
Confirmed the recommendation of 2026-10-21T11:11:33Z orchestrator.

Source: design/system/strategic-agents.md § Cost of delay
