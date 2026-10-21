---
id: S-0414
type: story
nature: feature
title: A strategic role is a markdown file under design/roles that flai reads and flai check validates, and the template ships the planner, the orchestrator, and the analyzer as baseline role files
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:39:25Z
updated: 2026-10-21T00:41:56Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, template, conventions]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/roles/analyzer.md, design/roles/orchestrator.md, design/roles/planner.md, design/roles/README.md, design/system/repository-layout.md, design/system/strategic-agents.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/root.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/roles/analyzer.md, template/root/design/roles/orchestrator.md, template/root/design/roles/planner.md, template/root/design/roles/README.md, template/template.yaml, design/system/flai-cli.md, design/system/project-manifest.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h1m
  delivery: 2026-10-21T04:54:00Z
  basis: "Its own forecast of 1h1m; 11th in the pull order with an in-progress limit of 5, behind S-0445, S-0453, S-0455, S-0454, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411 and S-0412."
  by: flai
  at: 2026-10-21T00:41:56Z
---
# S-0414 A strategic role is a markdown file under design/roles that flai reads and flai check validates, and the template ships the planner, the orchestrator, and the analyzer as baseline role files

## Goal

A project can say what a strategic role is in a file, not in Go. A role file under `design/roles/` gives the role's name, the verb it primes and runs as, its harness definition, its host action and trigger, the item it takes, its topics, the permissions it may hold, the MCP tools, commands, and files the guard lets it use, and whether it may work a story; its body is its harness prompt, with the baseline marker below which a project adds its own. flai reads the files, `flai role list` and `flai role show` print them, and `flai check` refuses a file that does not parse or that grants its role a call the operator's permissions do not allow. The template ships the planner, the orchestrator, and the analyzer as baseline role files that state today's rules, so the later stories of E-0024 can read them in place of the Go that holds them now.

## Acceptance criteria

- [ ] An ADR records the role file: where it lives (`design/roles/<name>.md`, or the folder TH-0471 settles), its front matter fields and their meaning, its body as the harness prompt, the baseline marker, and how a rule that is not a plain allow-list (a story created only as a draft, a move only to backlog, an edit only inside one folder, a thread reply as the orchestrator's `answer_threads` allows) is named as a condition flai implements and the file references.
- [ ] `flai/internal/roles` reads every role file of a project, defaults and all, and `flai role list [--json]` and `flai role show <name> [--json]` print each role with its source file, action, trigger, item kind, topics, permissions, and calls.
- [ ] `flai check` reports `roles.file` for a role file that does not parse, misses a field, or names an unknown condition, trigger, or item kind; `roles.grant` for a role that grants itself, outside any permission it declares, a call that only the operator's permissions may allow (an acceptance, a publish, a move past backlog, finalizing a draft, setting cost of delay inputs); and `roles.definition` for a role whose harness definition is missing or lists a tool its rules refuse.
- [ ] `design/roles/planner.md`, `orchestrator.md`, and `analyzer.md` in the template and in this repository state the three roles' rules as `design/system/strategic-agents.md` and ADR-0082, ADR-0087, ADR-0099, and ADR-0143 describe them, with a `README.md` for the folder, and `flai check --strict` reports nothing on them here.
- [ ] `flai upgrade` brings a baseline role file a project lacks and leaves what a project wrote below the marker of one it has, shown by a test.
- [ ] `design/system/strategic-agents.md`, `design/system/repository-layout.md`, the convention `strategic-agents.md` (template and here), `docs/users/flai.md` with `flai-reference.md` regenerated, `docs/operators/index.md`, and `template/CHANGELOG.md` say what a role file holds and where it lives; the convention says a defined role runs as the operator, as the epic's Notes ask.

## Tasks
- T-1720 An ADR records the role file: where it lives, its front matter, its body as the prompt, the marker, and the named conditions it may reference
- T-1721 flai/internal/roles reads and validates a project's role files into roles with their fields, conditions, and definitions
- T-1722 The template and this repository ship design/roles with planner.md, orchestrator.md, and analyzer.md stating today's rules, and flai upgrade brings them
- T-1723 flai check reports roles.file, roles.grant, and roles.definition, and flai role list and flai role show print a project's roles
- T-1724 The design, the convention, and the user and operator docs say what a role file holds, where it lives, and that a defined role runs as the operator

## Notes

- Drafted by the planner for E-0024. First story of the epic: the format every later story reads.
- The three role files are data only in this story. flai keeps deciding from its Go until S-0415 (prime), S-0416 (guard), and S-0417 to S-0419 (serve) read them, so nothing changes in behaviour here.
- `.claude/agents/*.md` stay the harness's definitions (ADR-0130). The role file names its definition; flai never writes under `.claude/`.

### Planning

- **Touches.** The story declared none. Every touch comes from its five tasks, file by file.
  - **From design:** the ADR's index (`design/adrs/README.md`); `strategic-agents.md`, `repository-layout.md`, `project-manifest.md`, the convention in both copies, `docs/users/flai.md`, `docs/operators/index.md`, and `template/CHANGELOG.md` for criterion 6.
  - **From layout:** the new package `flai/internal/roles` (`roles.go`, `conditions.go`, and their tests), `flai/cmd/role.go` and `root.go`, `flai/internal/check/check.go` with a new `roles_test.go`, `flai/internal/upgrade/upgrade.go` for the marker merge, the six role files and two READMEs, and `template/template.yaml` for the version.
  - **Co-change, added:** `design/system/flai-cli.md` (44% of the seed commits), which lists every command; it is a shared path, so it holds nothing.
  - **Co-change, not added:** `docs/operators/settings.md`, `flaiover-dashboard.md`, `docs/users/flaiover.md`, `workflow.md`, and `work-management.md`. No setting, page, or workflow rule changes here.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. The folder is in `claims.shared`, so it holds no story.
- **Forecast.** Kept at flai's 1h1m. The duplicated template paths inflate the size about as much as the ADR and three roles' rules add. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
