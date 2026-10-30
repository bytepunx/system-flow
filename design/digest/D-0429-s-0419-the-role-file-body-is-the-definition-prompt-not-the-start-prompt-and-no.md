---
id: D-0429
title: "S-0419: the role file body is the definition prompt, not the start prompt, and no continuous launcher exists; how should criteria 1, 2, and T-1759 read: if we go with A now, does that prevent us from…"
kind: recommendation
sources: [TH-0760]
topics: [cli, planning, orchestration, template, conventions, flai, go, code, all]
paths: [design/conventions/strategic-agents.md, design/workflows/roles/analyzer.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/system/strategic-agents.md, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/serve/agents.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analyze.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/plan.go, flai/internal/serve/replan.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, flai/internal/check/activity_test.go, flai/internal/check/check_test.go, flai/internal/check/epic_lags_test.go, flai/internal/metrics/testdata/good/design/workflows/roles/analyzer.md, flai/internal/metrics/testdata/good/design/workflows/roles/curator.md, flai/internal/metrics/testdata/good/design/workflows/roles/orchestrator.md, flai/internal/metrics/testdata/good/design/workflows/roles/planner.md]
date: 2026-10-30
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0429 S-0419: the role file body is the definition prompt, not the start prompt, and no continuous launcher exists; how should criteria 1, 2, and T-1759 read: if we go with A now, does that prevent us from…

The operator, alex, answered TH-0760 at 2026-10-30T12:56:51Z.

What was asked:

```text
Two of the plan's premises do not hold in the code today (worktree read 2026-10-30).
```

The recommendation it overruled:

```text
Recommended for the operator: take option A.

- **ADR-0207 § 4 already settles premise 1.** The role file's body is the definition's prompt, the system prompt. The start prompt stays flai's: `planPrompt`, `orchestratePrompt`, and `analyzePrompt`. Criterion 1 and T-1753, as written, contradict the accepted ADR. A follows it.
- **Premise 2 holds too.** ADR-0209 made the orchestrator run for a reason and then end. T-1759's "continuous role" has no launcher to build on, so drop `role_continuous.go` from it.
- **Criterion 1 would read:** "The three roles' definition prompts carry their role file's body, its project additions included, from the project's file or the embedded baseline, and `flai/internal/harness/harness_test.go` (`TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, `TheAnalyzer`) passes with the same start prompt text."
- **Criterion 2:** keep it, with "the same start path a defined role takes" in place of "the defined-role path".
- **What stays out of S-0419 (D-0351):** the follow-up story keeps the rest of definition generation. That is checking the generated definitions against `.claude/agents/*.md` and settling the symlink for sessions started by hand. In S-0419, the additions reach the prompt only through the `--agents` JSON, and nothing is written under `.claude/`.
- **Until that story lands, D-0411 still applies:** a change to a baseline prompt is made in both `.claude/agents/<name>.md` and the role file.

B would put start-prompt text into the bodies, which goes against ADR-0207 § 4. C splits off the launcher work with no gain, since A keeps the start prompts in Go and the risk is small.

Rewording criteria 1 and 2 and T-1753 and T-1759 is yours.

Source: design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md § 4. The body and the harness definition
```

The operator's answer:

```text
if we go with A now, does that prevent us from eventually moving role definitions out of Go implementation so that the 3 existing strategic agents become defined solely by data files?
```

The operator, alex, answered TH-0760 again at 2026-10-30T14:05:53Z.

The answer:

```text
Recommended for the operator: take option A for S-0419. It does not block roles defined only by data files.

**Your question: does A prevent it?** No. agent-S-0419's answer above holds against ADR-0207 § 4.

- After A, each built-in role's definition is read from its role file.
- All three roles start through the one start path a defined role takes.
- What stays in Go after A is the start prompt and the replan and schedule rules. Moving the start prompt means adding a templated start-prompt section to the role file. That refines ADR-0207 § 4 in a new ADR and leaves the body alone. B would put that text into the body now, against § 4 as accepted.

**What to do on S-0419:**

1. Criterion 1 reads: "The three roles' definition prompts carry their role file's body, its project additions included, from the project's file or the embedded baseline, and `flai/internal/harness/harness_test.go` (`TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, `TheAnalyzer`) passes with the same start prompt text."
2. Criterion 2 keeps its words, with "the same start path a defined role takes" in place of "the defined-role path".
3. T-1753 drops the move of start-prompt text into the bodies.
4. T-1759 drops `role_continuous.go`. ADR-0209 makes the orchestrator run for a reason and end, so there is no continuous launcher to build on.

**Follow-up, a scope addition for you:** add a backlog story under E-0024, after the definition-generation story (D-0351). It covers two things: a templated start-prompt section in the role files, under a new ADR refining ADR-0207 § 4, and the planner's replan rules and the analysis schedule as `event` and `schedule` triggers. Once it lands, the three built-in roles are defined by their files alone. E-0024's planner can draft it.

Until the generation story lands, D-0411 still holds: a change to a baseline prompt is made in both `.claude/agents/<name>.md` and the role file.

Source: design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md § 4. The body and the harness definition
```

Source: TH-0760.
