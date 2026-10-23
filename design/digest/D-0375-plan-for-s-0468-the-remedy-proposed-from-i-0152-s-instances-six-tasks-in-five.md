---
id: D-0375
title: "Plan for S-0468: the remedy proposed from I-0152's instances, six tasks in five layers: Recommended for the operator: take the planner's four remedies, with the refusal for a task's folder touch…"
kind: answer
sources: [TH-0599]
topics: [planning, flai, cli, go, template, conventions, code, all]
paths: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/planning_test.go, flai/internal/itemedit/folder_test.go, flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, ".claude/agents/planner.md", template/root/.claude/agents/planner.md, template/template.yaml, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md, design/issues/summary.md, flai/cmd/edit.go, flai/cmd/check.go, flai/cmd/check_stats_test.go, design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0375 Plan for S-0468: the remedy proposed from I-0152's instances, six tasks in five layers: Recommended for the operator: take the planner's four remedies, with the refusal for a task's folder touch…

The operator, alex, answered TH-0599 at 2026-10-23T04:57:33Z.

What was asked:

```text
**The remedy proposed.** I-0152's instance is E-0016's planned stories: folder touches (`flai/cmd`, `flai/internal/*`, `flaiover/src/routes`, `docs/users`) and `after` hubs (S-0218, S-0223). What stands already asks for files but enforces nothing: the convention and `planPrompt` say name files (S-0295), ADR-0096 narrows a story's folder to the files its tasks name, ADR-0134 shares. ADR-0096's narrowing works only when tasks name files, and nothing makes them. So the tasks build, through flai's write paths rather than words alone:
```

The answer:

```text
Recommended for the operator: take the planner's four remedies, with the refusal for a task's folder touch rather than a warning.

- **Why a refusal.** The convention and `planPrompt` already ask for file touches (S-0295), and I-0152 shows they were not followed: folder touches and `after` hubs held up to ten ready stories a day. ADR-0096's narrowing only works once tasks name files. A refusal at the write is testable, and it says what to do.
- **Why the planner alone, for stories.** Your words on a story stay yours. A story's agent's folder touches are narrowed by its tasks and reported by `touches.folder` while in progress, so they hold only as long as they need to.
- **The hub rule** (`wip.after-hub` at three dependents) and the convention line that `after` names only a needed change are advisory, and cheap.
- **The one cost.** Agents that today write a folder touch on a task get refused and must name files. That is the intended behaviour, but expect a few refusals in the first stories after it lands.

S-0468 is held on overlap by S-0466 (`flai/internal/harness/harness.go`), so this can be answered before it starts. Its value is 6,600 USD a week.

Source: design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md
```

Source: TH-0599.
