---
id: TH-0760
title: "S-0419: the role file body is the definition prompt, not the start prompt, and no continuous launcher exists; how should criteria 1, 2, and T-1759 read?"
anchor:
  path: wip/kanban/stories/S-0419-the-planner-the-orchestrator-and-the-analyzer-start-through-the-path-a-defined-role-takes-with-their-prompts-read-from-their-role-files-and-behave-exactly-as-today.md
  item: S-0419
status: resolved
participants: [agent-S-0419, orchestrator, alex, planner-E-0024]
created: 2026-10-30T00:15:32Z
updated: 2026-10-31T06:23:24Z
---

# TH-0760 S-0419: the role file body is the definition prompt, not the start prompt, and no continuous launcher exists; how should criteria 1, 2, and T-1759 read?

On wip/kanban/stories/S-0419-the-planner-the-orchestrator-and-the-analyzer-start-through-the-path-a-defined-role-takes-with-their-prompts-read-from-their-role-files-and-behave-exactly-as-today.md.

## Entries

### 2026-10-30T00:15:32Z agent-S-0419
Two of the plan's premises do not hold in the code today (worktree read 2026-10-30).

**1. The role file body is not what `planPrompt` produces.**

- The body of `design/workflows/roles/planner.md` (and the orchestrator's and analyzer's) is the agent *definition* prompt, the system prompt. It is the same text as `.claude/agents/planner.md`, and D-0411 keeps the two in step.
- `planPrompt`, `orchestratePrompt`, and `analyzePrompt` (`flai/internal/harness/harness.go:576`, `:692`, `:786`) are the *start* prompt, the user message. It is per-run Go text: epic or story branches, reasons with prime forms, the merge queue, the focus.
- T-1753 says to "bring the bodies ... to the exact text the functions produce now". That would replace the definition text in the role files with start-prompt text and a template syntax for the per-run parts. The bodies would then no longer match `.claude/agents/*.md`, which D-0351 keeps as the built-ins' definitions in this story.

**2. There is no continuous-role launcher.**

- S-0418 shipped `event` and `schedule` triggers only. `roles.go:570` refuses `continuous`, and `role_continuous.go` does not exist.
- ADR-0209 made the orchestrator run for reasons and end. Its triggers are already read from `orchestrator.md` (`orchestrate.go:270`, `orchestratorTriggers`).
- T-1759 asks to run the orchestrator "as the continuous role", which has nothing to build on.

**Recommended: A.**

- **A. Keep the body as the definition, and make the start path shared.**
  - The three built-in launches load their role file: the project's file, or the embedded baseline when it is missing.
  - They start through the same harness request a defined role builds: `l.role` is generalised, and the built-in start prompts stay in Go.
  - The harness appends the text below the role file's marker to the definition it passes in `--agents`. So a project's additions reach the planner, orchestrator, and analyzer without a flai release. The harness tests stay unchanged.
  - Records stay under `plans`, `orchestrator`, and `analyzer`. The replanner, the analysis schedule, and the orchestrator's reasons look and hold keep deciding *when*.
  - T-1759 drops `role_continuous.go`.
  - Criterion 1 would read: "the three roles' definition prompts carry their role file's body, its additions included, from the project's file or the embedded baseline, and the harness tests pass with the same prompt text."
  - Criterion 2 stands, saying "the same start path" in place of "the defined-role path".
- **B. As planned.** The start prompts move into the role bodies as templates over the request's fields. The bodies stop matching `.claude/agents/*.md` until the generation story lands. This is larger, and it is the riskier change to all three launchers.
- **C. Do A's additions, baseline, and `roles.baseline` here, and leave the launcher unification to a story of its own.**

Rewording the criteria is yours, so I have not edited S-0419. Meanwhile I am working T-1754 (the embedded baseline compared whole, and the `roles.baseline` warning), which no answer changes.

### 2026-10-30T00:17:02Z orchestrator (recommendation)
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

### 2026-10-30T12:56:51Z alex
if we go with A now, does that prevent us from eventually moving role definitions out of Go implementation so that the 3 existing strategic agents become defined solely by data files?

### 2026-10-30T12:57:19Z agent-S-0419
No. A does not prevent it, and it is the first step toward it.

What A moves to data now:

- Each built-in role's definition comes from its file, the project's or flai's embedded copy. A project's additions below the marker reach it.
- The three roles start through the same start path a defined role takes, so one path serves all of them.
- The orchestrator's triggers already come from its file.

What stays in Go after A, and how each part could move later:

1. **The start prompts** (`planPrompt`, `orchestratePrompt`, `analyzePrompt`).
   - They have per-run branches: an epic or a story, reasons with their prime forms, the merge queue on or off, the focus.
   - To move them, a role file gains a start-prompt section beside the body, templated over the run's fields (item, kind, reasons, inputs, focus).
   - ADR-0207 § 4 makes the body the definition, so this would be a new ADR refining it, not a rewrite of the body. A defined role's start prompt (`definedPrompt`) would move the same way.
2. **When each role runs.**
   - The planner's replan rules and the analysis schedule are Go today. Both fit the `event` and `schedule` triggers S-0418 built for defined roles.
   - The orchestrator's reasons look already reads its triggers from its file.
3. **The records and hold.** `plans`, `orchestrator`, `analyzer`, the past runs, and the orchestrator's stop and hold are bookkeeping. They can key off the role's name once nothing else tells the roles apart.

Under A the launchers already load the role file and go through the shared start path. Each of those later steps would then change only where its text or rule comes from, not how the run starts.

B would do step 1 now, inside the body. That goes against ADR-0207 § 4 and against the `.claude/agents` copies that D-0411 keeps in step until the generation story.

Recommendation: take A now. Add a backlog story under E-0024 for the start-prompt section and the triggers (steps 1 and 2), after the definition-generation story.

Meanwhile I am doing the part every option needs, which A and C share: the harness passes each built-in role's definition from its role file with its additions. On your word for A, I will rewrite criteria 1 and 2 and the remaining tasks as the orchestrator's recommendation above words them.

Source: design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md § 4. The body and the harness definition

### 2026-10-30T12:57:49Z orchestrator (recommendation)
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

### 2026-10-30T13:34:05Z agent-S-0419
Status: S-0419's work is done, and its close-out passed (`flai verify S-0419 --last`). Integration and smoke are left to the merge queue.

What was built:

- `startOf`, the start every run goes through (a defined role's included), gives the planner, orchestrator, analyzer, and curator their role file: the project's, else flai's baseline.
- The harness makes that file's body, with the project's additions, the prompt of the definition it passes in `--agents`.
- The start prompts, records, and when are unchanged, and the harness tests pass unchanged.
- `flai check` advises `roles.baseline` for a missing file.

T-1757 and T-1759 are cancelled: `startOf` already covers those starts.

Criteria 3, 4, and 5 are ticked. Criteria 1 and 2 are not, because as written they describe option B. The code meets the orchestrator's rewording above:

1. "The three roles' definition prompts carry their role file's body, its project additions included, from the project's file or the embedded baseline, and `flai/internal/harness/harness_test.go` (`TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, `TheAnalyzer`) passes with the same start prompt text."
2. The same as now, with "the same start path a defined role takes" in place of "the defined-role path".

Reply "A" and I will reword both criteria that way, tick them, and move S-0419 to review. Reply otherwise and I will rework to fit.

### 2026-10-30T14:05:53Z alex
Confirmed the recommendation of 2026-10-30T12:57:49Z orchestrator.

Source: design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md § 4. The body and the harness definition

### 2026-10-30T14:07:18Z orchestrator
The operator's confirmation settles TH-0760: option A.

- S-0419's criteria 1 and 2 already read as recommended and are ticked. agent-S-0419 can move the story to review once its worktree is committed.
- The follow-up story under E-0024 is still to add. It covers a templated start-prompt section under a new ADR refining ADR-0207 § 4, and the replan rules and analysis schedule as `event` and `schedule` triggers, after the definition-generation story (D-0351). E-0024 is not on `flai plan --candidates`, so I cannot start its planner for it. agent-S-0419: add it with `flai story new --epic E-0024` as a backlog draft, and the orchestrator will plan it as a story candidate.

Source: design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md

### 2026-10-30T20:38:30Z alex
Are you able to write the story?

### 2026-10-30T20:39:36Z orchestrator
Not by hand: the orchestrator cannot create work items. I have started E-0024's planner instead, and it can (planner-E-0024, started 2026-10-30T20:39:30Z, trigger TH-0760).

- **Why the planner:** `flai plan --candidates` lists E-0024 as planned before, because you wrote here after its last run. Your entry asks for a story to be added, the follow-up you confirmed at 14:05:53Z. That is a change to the plan, so I may start its planner for it (ADR-0222).
- **What it should draft:** a backlog story under E-0024, after the definition-generation story (D-0351). It covers a templated start-prompt section in the role files, under a new ADR refining ADR-0207 § 4, and the planner's replan rules and the analysis schedule as `event` and `schedule` triggers.
- **What happens next:** the planner drafts the story and its tasks, then opens a plan thread. I settle that thread under plan_backlog_stories. The story stays a draft until `flai promote --drafts` lists it as complete.

S-0419 still waits on TH-0765. TH-0760 has to be resolved before the merge queue takes S-0419. Its opener, agent-S-0419, or you can resolve it.

Source: design/adrs/0222-with-plan-backlog-epics-or-plan-backlog-stories-on-the-orchestrator-may-ask-the.md

Source: design/adrs/0222-with-plan-backlog-epics-or-plan-backlog-stories-on-the-orchestrator-may-ask-the.md

### 2026-10-30T20:40:04Z alex
Resolved.

### 2026-10-30T20:50:07Z planner-E-0024
Yes. The story is written: S-0554, a draft under E-0024. It moves the start prompts and the planner's, analyzer's, and curator's triggers into the role files, under a new ADR refining ADR-0207 § 4.

It waits for S-0553, also a new draft. S-0553 is the definition-generation story D-0351 named, which no one had created yet. The plan, both stories' tasks, and two proposals for you are on TH-0768.

### 2026-10-30T20:52:08Z orchestrator
The follow-up you asked for now exists. The plan is on TH-0768.

- **S-0553** generates the built-in definitions (D-0351). It is finalized and waits for S-0419.
- **S-0554** moves the start prompts and triggers into the role files. It is a draft and waits for S-0553. I have recommended "no split" on TH-0768. I finalize S-0554 once you answer there.

**This thread is still not resolved.** You wrote "Resolved." at 20:40:04Z, but its status is still `answered`. The merge queue refuses S-0419 while it stays open:

```text
thread_open: threads TH-0760, TH-0765 are not resolved on S-0419 or its tasks
```

I cannot resolve it, because I did not open it. To clear it, resolve it from the dashboard or run:

```bash
flai thread resolve TH-0760 --reason "option A confirmed; follow-up is S-0553 and S-0554"
```

Once it is resolved, I verify S-0419 again at its new head, `8a38ec0e`, and enqueue it. Then I resolve TH-0765.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-31T06:23:24Z alex
Resolved: S-0419 was accepted
