---
id: S-0414
type: story
nature: feature
title: A strategic role is a markdown file under design/workflows/roles that flai reads and flai check validates, and the template ships the planner, the orchestrator, and the analyzer as baseline role files
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:39:25Z
updated: 2026-10-28T12:40:13Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, template, conventions]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/workflows/roles/analyzer.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/workflows/roles/README.md, design/system/repository-layout.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/root.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/README.md, template/template.yaml, design/system/flai-cli.md, design/system/project-manifest.md, design/README.md, template/root/design/README.md.tmpl, CLAUDE.md, template/root/CLAUDE.md.tmpl]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 417
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 584
          cache_read: 4071625
          cache_write: 35288
          cost: 0.9069
draft: true
cost_of_delay:
  value: 66.76
  by: planner-S-0414
  at: 2026-10-21T11:10:38Z
forecast:
  duration: 1h17m
  delivery: 2026-10-28T17:21:00Z
  basis: "Its own forecast of 1h17m; 9th in the pull order with an in-progress limit of 5, behind S-0503, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411 and S-0412."
  by: flai
  at: 2026-10-28T12:40:13Z
---
# S-0414 A strategic role is a markdown file under design/workflows/roles that flai reads and flai check validates, and the template ships the planner, the orchestrator, and the analyzer as baseline role files

## Goal

A project can say what a strategic role is in a file, not in Go. A role file under `design/workflows/roles/` gives the role's name, the verb it primes and runs as, the fields flai generates its harness definition from, its host action and trigger, the item it takes, its topics, the permissions it may hold, the MCP tools, commands, and files the guard lets it use, and whether it may work a story; its body is its harness prompt, with the baseline marker below which a project adds its own. flai reads the files, `flai role list` and `flai role show` print them, and `flai check` refuses a file that does not parse or that grants its role a call the operator's permissions do not allow. The template ships the planner, the orchestrator, and the analyzer as baseline role files that state today's rules, so the later stories of E-0024 can read them in place of the Go that holds them now.

## Acceptance criteria

- [ ] An ADR records the role file: where it lives (`design/workflows/roles/<name>.md`, in a `design/workflows/` folder the ADR adds to the layout, as TH-0471 settled), its front matter fields and their meaning, which of them flai generates the role's harness definition from and how that definition reaches Claude Code without a hand-written `.claude/agents` file (TH-0471, answer 6), its body as the harness prompt, the baseline marker, and how a rule that is not a plain allow-list (a story created only as a draft, a move only to backlog, an edit only inside one folder, a thread reply as the orchestrator's `answer_threads` allows) is named as a condition flai implements and the file references.
- [ ] `flai/internal/roles` reads every role file of a project, defaults and all, and `flai role list [--json]` and `flai role show <name> [--json]` print each role with its source file, action, trigger, item kind, topics, permissions, and calls.
- [ ] `flai check` reports `roles.file` for a role file that does not parse, misses a field, or names an unknown condition, trigger, or item kind; `roles.grant` for a role that grants itself, outside any permission it declares, a call that only the operator's permissions may allow (an acceptance, a publish, a move past backlog, finalizing a draft, setting cost of delay inputs); and `roles.definition` for a role whose harness definition fields are missing or list a tool its rules refuse.
- [ ] `design/workflows/roles/planner.md`, `orchestrator.md`, and `analyzer.md` in the template and in this repository state the three roles' rules as `design/system/strategic-agents.md` and ADR-0082, ADR-0087, ADR-0099, and ADR-0143 describe them, with a `README.md` for the folder, and `flai check --strict` reports nothing on them here.
- [ ] `flai upgrade` brings a baseline role file a project lacks and leaves what a project wrote below the marker of one it has, shown by a test.
- [ ] `design/system/strategic-agents.md`, `design/system/repository-layout.md`, `design/README.md` and `CLAUDE.md` (template and here), the convention `strategic-agents.md` (template and here), `docs/users/flai.md` with `flai-reference.md` regenerated, `docs/operators/index.md`, and `template/CHANGELOG.md` say what a role file holds, where it lives, and that flai generates its harness definition; the convention says a defined role runs as the operator, as the epic's Notes ask.

## Tasks
- T-1720 An ADR records the role file: where it lives, its front matter, its body as the prompt, the marker, and the named conditions it may reference
- T-1721 flai/internal/roles reads and validates a project's role files into roles with their fields, conditions, and definitions
- T-1722 The template and this repository ship design/workflows/roles with planner.md, orchestrator.md, and analyzer.md stating today's rules, and flai upgrade brings them
- T-1723 flai check reports roles.file, roles.grant, and roles.definition, and flai role list and flai role show print a project's roles
- T-1724 The design, the convention, and the user and operator docs say what a role file holds, where it lives, and that a defined role runs as the operator

## Notes

- Drafted by the planner for E-0024. First story of the epic: the format every later story reads.
- The three role files are data only in this story. flai keeps deciding from its Go until S-0415 (prime), S-0416 (guard), and S-0417 to S-0419 (serve) read them, so nothing changes in behaviour here.
- `.claude/agents/*.md` stay the harness's definitions (ADR-0130). The role file names its definition; flai never writes under `.claude/`.

### Planning

- **TH-0471's answers (2026-10-21).** The plan follows them, and they change this story in three places.
  - Answer 2 moves the folder from `design/roles/` to `design/workflows/roles/`. The title, the goal, and criteria 1, 4, and 6 now say so, and every `design/roles` touch became its `design/workflows/roles` path. `repository-layout.md` adds no `design/` subfolder without an ADR, so T-1720's ADR also records `design/workflows/`.
  - Answer 6 has flai generate a role's harness definition from its role file. That overrides the third note above. Criteria 1, 3, and 6 now say so. Generating the definition, and the symlink into `.claude/agents` where Claude Code follows one, is proposed on TH-0567 for a later story. In this story the role files carry the definition fields and `roles.definition` checks them, so behaviour stays unchanged.
  - Answer 1's input, `time_lost_per_cycle: 4h` on E-0024, is not yet on the epic.
  - T-1720, T-1722, T-1723, and T-1724 were rewritten to these answers, as the orchestrator approved on TH-0567. Their layers and `after` are unchanged.
- **Touches.** All 37 are files but `design/adrs`.
  - **Declared, kept:** the 24 earlier touches that are not role files, from the five tasks.
  - **Declared, moved:** the eight role files and their two READMEs, now under `design/workflows/roles/` (answer 2).
  - **From design:** the ADR's index; `strategic-agents.md`, `repository-layout.md`, `project-manifest.md`, the convention in both copies, `docs/users/flai.md`, `docs/operators/index.md`, and `template/CHANGELOG.md` for criterion 6.
  - **From layout, added:** `design/README.md`, `template/root/design/README.md.tmpl`, `CLAUDE.md`, and `template/root/CLAUDE.md.tmpl`. Each lists the `design/` folders and gets a row for `design/workflows/` (T-1724).
  - **Co-change, added:** `docs/operators/settings.md` (16% of the seed commits). `make flai-reference` rewrites its flag index when `flai role` lands (T-1723). `design/system/flai-cli.md` stays; it is shared, so it holds nothing.
  - **Co-change, not added:** `flaiover-dashboard.md`, `docs/users/flaiover.md`, `workflow.md`, `work-management.md`, and `flai/internal/harness`. No page, workflow rule, or start path changes here.
  - **`upgrade.go` kept** although its merge is already path-generic (`Compute`, `merge`), in case a role file needs a case the merge lacks.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. The folder is in `claims.shared`, so it holds no story.
- **Forecast.** flai's 1h17m is kept: size 43 at 107 s a unit. The four new layout rows are one line each, so they overstate the size a little. Settling how the definition is generated adds about as much to the ADR. The delivery is flai's play-out, 2026-10-21T15:40Z.
- **Cost of delay.** 66.76 USD a week, worked out by hand from the operator's answer 1, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week.
  - S-0414 takes the share its 1h17m is of the 11h32m forecast for E-0024's twelve open stories: 77 / 692 × 600.
  - Once the input is on E-0024, `flai cod S-0414` gives the same figure while the forecasts stand.
