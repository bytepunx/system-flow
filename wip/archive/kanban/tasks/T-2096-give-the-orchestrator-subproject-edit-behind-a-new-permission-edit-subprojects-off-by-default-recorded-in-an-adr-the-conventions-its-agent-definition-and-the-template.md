---
id: T-2096
type: task
nature: improvement
title: Give the orchestrator subproject_edit behind a new permission edit_subprojects, off by default, recorded in an ADR, the conventions, its agent definition, and the template
status: done
parent: S-0491
owner: alex
created: 2026-10-23T20:40:36Z
updated: 2026-10-23T22:17:48Z
transitions:
  - to: ready
    at: 2026-10-23T21:50:49Z
    by: agent-S-0491
  - to: in-progress
    at: 2026-10-23T21:50:49Z
    by: agent-S-0491
  - to: done
    at: 2026-10-23T22:17:48Z
    by: agent-S-0491
stream: S-0491
tags: [cli, template]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, design/adrs/README.md, design/system/strategic-agents.md, design/conventions/strategic-agents.md, ".claude/agents/orchestrator.md", docs/operators/settings.md, docs/users/flai.md, template/root/design/conventions/strategic-agents.md, template/root/.claude/agents/orchestrator.md, flai/internal/harness/harness.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings_test.go, flai/internal/guard/guard_test.go, flai/internal/guard/subproject_test.go, flai/internal/mcpserver/subproject.go, design/adrs/0195-with-edit-subprojects-on-the-orchestrator-adds-changes-and-removes-a-sub.md]
after: [T-2095]
usage:
  source: log
  seconds: 1619
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 157
      output: 64560
      cache_read: 9665390
      cache_write: 376126
      cost: 5.2902
---
# T-2096 Give the orchestrator subproject_edit behind a new permission edit_subprojects, off by default, recorded in an ADR, the conventions, its agent definition, and the template

## Work

The story's third criterion names the orchestrator among the agents that manage sub-projects, and the orchestrator acts only within `orchestration.permissions` (`strategic-agents.md`). Add the permission `edit_subprojects`, off by default:

- `manifest.go`: the constant `PermitEditSubprojects` in `PermissionNames`, in their order, so `Permissions.UnmarshalYAML` and `Allows` know it; `settings.go`: its entry in the settings catalog, so `flai manifest set orchestration.permissions.edit_subprojects=true` writes it and the dashboard's orchestrator panel shows it as it shows the others.
- `guard.go` § `orchestrate`: `subproject_edit` and the shell commands `flai subproject add`, `set`, and `remove` need `PermitEditSubprojects`; refused otherwise with `off(needs)`, naming the permission. Tests in `orchestrate_test.go` show on and off.
- Record the decision with `flai adr new`: the orchestrator edits a project's sub-projects through `subproject_edit` only behind `orchestration.permissions.edit_subprojects`, commits the main checkout's change itself with `flai subproject ... --autocommit` where it would rather use the shell, and never creates a sub-project on its own judgement without a story or thread that asks for it. The ADR's path is not known until it is made, so this task names `design/adrs/README.md`, the index `flai adr new` rewrites, and the story keeps the folder in its touches.
- Say it in `design/system/strategic-agents.md` § What it does with each permission, in `design/conventions/strategic-agents.md` § As the orchestrator ("Each of these only while its permission is on" gains the line, and the `With ...` list gains `edit_subprojects`), in `.claude/agents/orchestrator.md` (its tools gain `mcp__flai__subprojects` and `mcp__flai__subproject_edit`, and a step says when it uses them), and in the template's copies of both; the permission's row goes into `docs/operators/settings.md` § Project manifest, which `TestSettingsIndexIsComplete` demands, and the permission list in `docs/users/flai.md` where the other permissions are described.

Writes under `.claude/` are asked of the operator through `permission_prompt`: make them yourself, not from a sub-agent.

It waits for T-2095, whose tool and guard rules it extends; it runs beside the dashboard task, which touches other files.

## Done when

- `flai manifest set orchestration.permissions.edit_subprojects=true` writes the key, and `flai check` accepts it.
- `flai guard` lets the orchestrator call `subproject_edit` and run `flai subproject add` with the permission on, and refuses both naming `orchestration.permissions.edit_subprojects` with it off, shown in `orchestrate_test.go`.
- An ADR under `design/adrs/` records the decision and `design/adrs/README.md` lists it.
- The convention, the design, the orchestrator's agent definition, `docs/operators/settings.md`, `docs/users/flai.md`, and the template's copies name the permission, and `go test ./cmd/ -run TestSettingsIndexIsComplete` and `scripts/flai-test.sh` pass for the paths changed.

## Notes

Precedent: `plan_backlog_stories` (ADR-0119), added to `PermissionNames`, the catalog, the guard, the convention, the agent definition, and the template in one story.
