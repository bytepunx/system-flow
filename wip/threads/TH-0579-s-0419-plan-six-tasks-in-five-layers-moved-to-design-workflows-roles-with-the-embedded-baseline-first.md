---
id: TH-0579
title: "S-0419 plan: six tasks in five layers, moved to design/workflows/roles, with the embedded baseline first"
anchor:
  path: wip/kanban/stories/S-0419-the-planner-the-orchestrator-and-the-analyzer-start-through-the-path-a-defined-role-takes-with-their-prompts-read-from-their-role-files-and-behave-exactly-as-today.md
  item: S-0419
status: resolved
participants: [planner-S-0419, orchestrator, alex]
created: 2026-10-21T12:16:34Z
updated: 2026-10-21T16:01:49Z
---

# TH-0579 S-0419 plan: six tasks in five layers, moved to design/workflows/roles, with the embedded baseline first

On wip/kanban/stories/S-0419-the-planner-the-orchestrator-and-the-analyzer-start-through-the-path-a-defined-role-takes-with-their-prompts-read-from-their-role-files-and-behave-exactly-as-today.md.

## Entries

### 2026-10-21T12:16:34Z planner-S-0419
S-0419 is replanned on your TH-0471 answers. It keeps its six tasks. The embedded baseline now comes before the prompt move. The story has a forecast and a cost of delay value.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-1754 | none | Baseline role files embedded in flai, the fallback when a project has none, and the `roles.baseline` warning |
| 2 | T-1753 | T-1754 | `planPrompt`, `orchestratePrompt`, and `analyzePrompt` read from the role files' bodies |
| 3 | T-1756 | T-1753 | The planner and the replanner start through the defined-role path, and the marker test |
| 4 | T-1757 | T-1756 | The analyzer and the analysis schedule start through the path |
| 4 | T-1759 | T-1756 | The orchestrator runs as a continuous role, with its stop, start, and hold |
| 5 | T-1761 | T-1754, T-1757, T-1759 | The design, the convention, and the changelog |

T-1757 and T-1759 run together; they share no path.

## What changed

- **Order.** T-1754 used to wait for T-1753. Now T-1753 waits for T-1754. The harness's prompt tests start roles in a temporary root with no role files. Without the embedded fallback they cannot pass unchanged.
- **Folder (answer 2).** The goal, criterion 4, and every role-file touch now say `design/workflows/roles/`.
- **Touches added.**
  - `flai/internal/roles/baseline/*.md`, on T-1754 and T-1753: `go:embed` cannot reach `template/` from the `flai` module, so flai embeds a copy.
  - `flai/internal/serve/agents.go` on T-1756: built-in runs keep `plans`, `orchestrator`, and `analyzer`, not `roles`.
  - `flai/internal/serve/serve.go` on T-1759: the loop's `orch.look` and S-0418's role loop must not both start the orchestrator.
- **Forecast.** 1h, raised from flai's 45m. The work is three launchers rewired with their tests unchanged, in five serial layers. Delivery 2026-10-21T22:37Z.
- **Cost of delay.** 48.45 USD a week. This is S-0419's share, 60 of 743 minutes, of the 600 USD a week your 4h gives E-0024. I worked it out by hand, because `flai cod` refuses while E-0024 has no inputs.

## Proposals, recommendation first

1. **Rewrite three task lines I did not write.**
   - **Recommended:** let me make these changes. T-1754's title and Work say `design/workflows/roles`, and its Work drops "Waits for the prompt task". T-1753's Work says "Second layer: waits for the embedded baseline" in place of "First layer". T-1756's Work says the run's record goes through `agents.go`.
   - Alternative: leave them to the story's agent, which reads this thread.
2. **Generated definitions for the built-in roles (TH-0471, answer 6; TH-0567, question 2).**
   - **Recommended:** a new draft story after S-0419, not inside it. `TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, and `TheAnalyzer` write their fixtures under `.claude/agents`. They also assert the refusal that names `flai upgrade` when a definition is missing. Criterion 1 keeps them unchanged, and generating the definitions would change both. The new story would generate the three from their role files, pass them in `--agents` as S-0417 does for defined roles, and check that they match today's `.claude/agents` files. It would also settle the symlink for sessions started by hand.
   - Alternative: widen S-0419 and let criterion 1 change those tests.
   - I have not drafted the story. That is for the epic's planner, or for you.
3. **E-0024's input.** Your 4h is not on the epic yet. I recommend `flai edit E-0024 --time-lost-per-cycle 4h`, as on TH-0567.

## Assumptions

- S-0414's role body format can hold the prompts' branches: epic or story, merge queue running or not, and focus or none. If it holds only plain substitution, T-1753 keeps the branches in Go, and the bodies hold the text of each branch.
- The defined-role path from S-0417 and S-0418 lets a role take its definition from `.claude/agents/<name>.md`, as the three built-in roles do today.
- The analysis scheduler stays. It starts the analyzer through the path, and the analyzer's role file names no schedule of its own, so the analyzer is not started twice.
- `flai/cmd` keeps its `serve` calls. No command, flag, or manifest key changes.

### 2026-10-21T12:19:30Z orchestrator
Approved under plan_backlog_stories: six tasks in five layers with the embedded baseline first, since the prompt tests need the fallback, plus the three added touches and the 1h forecast.

**Question 1:** yes, rewrite the three task lines as you recommend: T-1754's folder and its dropped wait, T-1753's layer, and T-1756's record through `agents.go`. They only follow the new order and the TH-0471 folder.

Questions 2 and 3 follow as a recommendation for the operator, so the thread stays open.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:19:34Z orchestrator (recommendation)
Recommended for the operator:

- **Question 2:** take the planner's answer. Generate the built-in roles' definitions in a new story after S-0419, not inside it. S-0419's criterion 1 keeps the three claude-code start tests unchanged, and those tests write fixtures under `.claude/agents`, so generating definitions there would break that promise. The new story would generate the three definitions from their role files and pass them in `--agents`, as S-0417 now does for defined roles. It would check them against today's files and settle the symlink for sessions started by hand. This replaces my recommendation on TH-0567, question 2, which put generation in S-0419. A new story is a scope addition, so it is yours; the epic's planner can draft it.
- **Question 3:** as on TH-0567, run `flai edit E-0024 --time-lost-per-cycle 4h`. Then `flai cod` values all of E-0024's stories in place of the planners' hand-worked figures. S-0419's would be 48.45 USD a week.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:21:02Z planner-S-0419
The three tasks are rewritten as approved.

- T-1754's title and Work name `design/workflows/roles/` and the embedded copy under `flai/internal/roles/baseline/`. Its Work now says it is the first layer, and why the prompt task waits for it.
- T-1753's Work says it is the second layer and also keeps the embedded copy in step.
- T-1756's Work routes a built-in role's record through `agents.go`, and its Done when adds `agents_test.go`, unchanged.

Layers, `after`, and touches are as in the table above. S-0419's Planning notes record the rewrite. Questions 2 and 3 are still the operator's.

### 2026-10-21T12:33:31Z alex
Confirmed the recommendation of 2026-10-21T12:19:34Z orchestrator.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:33:55Z orchestrator
Noted: the operator confirmed both answers. Two steps remain, and both need the operator, since neither is within my permissions:

1. E-0024 still has no cost of delay input. Run `flai edit E-0024 --time-lost-per-cycle 4h`, and `flai cod` will then value every E-0024 story.
2. The new story, generating the built-in roles' definitions after S-0419, needs E-0024's planner: `flai plan E-0024`, or Plan on the epic's page. `flai plan --candidates` does not list an epic that already has stories, so I cannot start that planner myself.

S-0419's own plan is settled. I am leaving this thread open until the new story exists.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-21T12:40:26Z alex
are you able to record the two items that just came up as issues? you should be allowed to update the cost of delay and you should be allowed to invoke the planner on epics/stories even if they’ve already been planned before

### 2026-10-21T16:01:44Z orchestrator
Yes. Both are recorded as issues, committed on main:

- **I-0187** (cabd78500): the orchestrator cannot set a cost of delay input you decided on an epic, or on a story past the backlog. I tried E-0024's 4h just now, and flai guard refused: it allows me story inputs only.
- **I-0188** (04ce72d24): `flai plan --candidates` lists no epic or story planned before, so I cannot start the planner on E-0024 again for the new story after S-0419.

Until those change, two steps stay with you:

1. `flai edit E-0024 --time-lost-per-cycle 4h`, and `flai edit E-0025 --time-lost-per-cycle 4h` (TH-0476).
2. `flai plan E-0024`, to draft the story that generates the built-in roles' definitions.

S-0419's own plan is settled, so I am resolving this thread.

### 2026-10-21T16:01:49Z orchestrator
Resolved: S-0419's plan is settled. The operator's two follow-ups are recorded as I-0187 and I-0188, and the E-0024 input and replan steps are named for the operator.
