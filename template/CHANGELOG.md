# Changelog

## 1.0.72 - 2026-10-07

- S-0330 flai message sends a message from one open story's agent to another's, kept apart from the operator's threads (patch).

## 1.0.71 - 2026-10-07

- S-0328 Add Permission and Ability to Orchestrate to Trigger Planner (patch).
- S-0328 The orchestrator plans the backlog's stories and settles their planner's threads with `plan_backlog_stories` (patch, ADR-0119): `strategic-agents.md`'s section "As the orchestrator" says what it does with the new permission. After the epics, it starts the planner (`plan`) for each story `flai plan --candidates` lists, one at a time: a backlog story with no touches, no forecast duration, no cost of delay value, or no task. On a thread a story's planner opened, its opener `planner-S-nnnn`, it acts whatever `answer_threads` says. It approves a plan that fits the story by replying and resolving the thread, and answers the planner's questions. For a cost of delay input, it takes the figure the planner recommends unless the thread or the story gives a reason for one of its alternatives, sets the inputs alone with `item_edit` `cost_of_delay`, never a value, replies naming what it set, and resolves the thread. A question that asks for a change of scope or a decision not yet recorded still gets a recommendation, and the thread stays open. It never confirms a recommendation, and an epic's planner's threads stay under `answer_threads`. It logs each with `activity_log`. The cost of delay inputs are the operator's, or the orchestrator's under `plan_backlog_stories` for a backlog story when neither it nor its epic has any. The orchestrator's agent definition, `orchestrator.md`, says the same and gains the tool `thread_resolve`. It needs a flai that has the permission.

## 1.0.70 - 2026-10-07

- S-0279 flai check finds `wip.overlap` outside the story at close-out (patch).
- S-0279 A close-out records no `wip.overlap` and notes only an overlap naming its story (patch, ADR-0115): `work-management.md` says that a check finding outside the story is still printed as a note and recorded in an issue, which the close-out commits, but that a `wip.overlap` naming your story is printed as a note and recorded in no issue, and one between two other stories is left out. `wip.overlap` compares the claims of the stories in progress, as the pull hold does, once per pair of stories, so a story is no longer reported against another story's task or against files inside its folder touch that its tasks do not name. It needs a flai that has it; an older flai still records every `wip.overlap` at close-out.

## 1.0.69 - 2026-10-07

- S-0275 Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch (patch).
- S-0275 The agent records an issue or an ADR in one call from its story's worktree (patch): `continuous-improvement.md` says that, working a story, you record, bump, or close an issue with `flai issue new|bump|close --commit`, or the MCP tool `issue_new`, `issue_bump`, or `issue_close` with `commit`, which numbers a new issue from the whole repository, commits exactly what it wrote on `story/S-nnnn` as `docs: [S-nnnn] ...`, with each `--trailer`, and adds it to the story's touches; you follow it with no `git add`, `git commit`, or `flai touches`. Off the story's branch it refuses and writes nothing, and `--autocommit` stays the main checkout's. `decisions.md` says the same of `flai adr new --commit`, or the MCP tool `adr_new` with `commit`, and that `flai adr new` numbers from the whole repository: this checkout, main, every worktree, and every story branch. `design/adrs/README.md` names `flai adr new --commit`. It needs a flai that has `--commit` and the MCP tools `issue_close` and `adr_new`.

## 1.0.68 - 2026-10-07

- S-0274 Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together (patch).
- S-0274 The agent begins a story with `flai story start` (patch): `session-start.md` says to pull the top ready story with `flai story start S-nnnn`, or the MCP tool `story_start`, which in one call moves it to `in-progress`, opens its narrative and its branch in a worktree, primes, and answers the inbox, and to work in the worktree it names; it refuses a story that is not ready or that the board holds, and changes nothing. Starting a story primes you; `flai prime --story` is for taking up a story already in progress, and `story_start` answers the pack's first part, or only its header when that part does not fit, so you read every part it does not hold with `prime`. A session flai started for `FLAI_STORY` calls `inbox` and then starts that story with `flai story start`. `work-management.md`'s order for pulling a story begins with `flai story start` in place of moving it and opening the narrative, and `wait_for_work`'s pull is started with `story_start`: when it says the story is in progress already, another agent pulled it, and you wait again. Step 3 of `CLAUDE.md`'s "Prime your session" pulls with `flai story start`, and step 1 says that starting a story primes you. It needs a flai that has `flai story start`.

## 1.0.67 - 2026-10-07

- S-0272 A story's agent with an open question ends instead of waiting, and flai serve starts it again on the answer (patch): `work-management.md` says that when you need the designer to decide something, you ask on a thread on the story or the task and go on with any work of the story that does not wait on the answer; when nothing is left but the answer, you write the narrative's `## Current state` and `## Next steps`, naming the question and what each answer leads to, and end the session, and `flai serve` starts you again, in the same session, when the thread is answered, with the answer in your first `inbox`. When `wait_for_events` answers `end: true`, the agent does the same; its `why` names the open question. An agent run by hand, not by `flai serve`, is not started again and holds `wait_for_events` until the thread is answered. `delegation.md` says to end on an open question only once every sub-agent is back, and that `wait_for_events` is held for a thread awaiting the designer by an agent that `flai serve` does not start again. The planner, the orchestrator, and the analyzer still hold `wait_for_events` for an answer. The `end: true` answer needs a flai that has it.
- S-0269 The agent closes each task with `flai task done` (patch): `git.md` says to close each task when it is done with `flai task done T-nnnn -m "<message>"` in the story's worktree, or the MCP tool `task_done`, which commits on `story/S-nnnn`, runs `flai stream sync`, moves the task to `done`, logs it in the narrative, adds the paths the commit changed to the task's and the story's `touches`, runs `flai check --strict` scoped to the story, and answers the inbox, stopping at the first step that fails. When the sync stops on conflicts, you resolve each path it lists, `git add` it, run `git rebase --continue`, and call it again; when the check stops, you fix what it found and call it again. You then run the task's tests with `flai test` on the paths it changed and close any fix by calling it again. `work-management.md`'s task cycle says the same, and says that `flai task done` answers the inbox, moves the task, and writes its log entry. It needs a flai that has `flai task done`.
- S-0270 The agent runs the close-out, which runs `flai verify`, itself, and keeps the verifier for the review (patch): `scripts/close-out.sh` runs `flai verify S-nnnn --record-issues` for its checks: no rebase left unfinished, the branch contains the main branch, the narrative's `## Current state` and `## Next steps` are written, `flai check --strict` scoped to the story, and the test and lint tiers the branch's changes select, cheapest first; it keeps the commit and the sync check after it. `delegation.md` says that before review the story's agent runs the close-out itself, once, as one command with the longest timeout and without a pipe or a file, reads its last line and the findings of the step it stopped at, fixes what it names, and runs it again; it never hands the run to a sub-agent, and `flai verify S-nnnn --last` prints the last result again. The verifier is kept for the review of the diff against the story's acceptance criteria and the conventions, when the diff is too large to read in the agent's own context: it reads the result with `flai verify S-nnnn --last` rather than running the suite, and runs the close-out or `flai verify` only when asked. `work-management.md` names `flai verify` as what the close-out runs. `.claude/agents/verifier.md` says the same, and `code-quality.md` no longer hands the suite to a verifier. Under `accept_reviews`, `strategic-agents.md` and `.claude/agents/orchestrator.md` have the orchestrator read or run `flai verify` at the branch's head before its verifier checks the diff against each criterion. `system-flow.yaml` gains a `markdown` tier, since the close-out's markdown lint is now a tier. It needs a flai that has `flai verify`.
- S-0271 The agent writes the narrative's state with `flai stream state` (patch): `tooling.md` says to write a narrative's `## Current state` and `## Next steps` with `flai stream state S-nnnn --current "<text>" --next "<text>"`, or the MCP tool `stream_state`, never by editing the file; it replaces the two sections, or the one you give, and leaves the rest alone. `## Context`, `## Decisions`, and `## Open questions` are still edited by hand, and the `## Log` is still appended by `flai stream log`. `work-management.md` says to rewrite the two sections with it after every task transition and before any long operation, and when you end on an open question; it refuses a story that is not in progress or in review, so when your story is cancelled you log what state the work is in with `flai stream log` instead. `session-start.md` says to write them with it before a long-running or risky operation and at the end of every session. It needs a flai that has `flai stream state`.

## 1.0.66 - 2026-10-07

- S-0273 flai test runs the project's test and lint tiers for a path or a package and answers pass or the first failures as findings (patch).

## 1.0.65 - 2026-10-06

- S-0273 A new project declares its test tiers in the manifest (patch): `system-flow.yaml` carries `tests`, the project's test tiers, cheapest first, each with a `name`, a `command` (an argument list that may hold `{packages}` or `{files}` as an argument of its own), the `paths` that select it (globs from the root, as `claims.shared`'s, with `!` to take paths out), the `format` of its output (`go-test-json`, `vitest-json`, `golangci-json`, `gofmt-list`, or `plain`), and optionally `dir`, `all_only`, and `all_command`. It starts with `scripts/test.sh` as the tier `test` for every path, and `scripts/integration.sh` and `scripts/smoke.sh` as `integration` and `smoke`, which run only under `--all`; all three are `plain`. Change them with `flai manifest set tests='<JSON list>'` or by hand. A project made before this has no `tests` key, and flai gives it one `plain` tier running `scripts/test.sh` when it exists. An older flai ignores the key, so it needs a flai that has it.
- S-0273 The agent runs the tests for what it changed with `flai test` (patch): `delegation.md` says that while it works, the story's agent runs the tests for what it changed itself, with `flai test` on the paths it changed or the MCP tool `test` with them, which runs the tiers the manifest's `tests` select for those paths, cheapest first, and answers pass or the first findings. It no longer runs `go test`, vitest, golangci-lint, or gofmt by hand and reads their logs, nor hands that run to a sub-agent; the verifier keeps the whole suite, the whole lint, and `flai check`, through the close-out. It runs `flai test` on what a task sub-agent changed when it reviews the task, and a task sub-agent runs its own tests the same way. `work-management.md`'s task cycle runs the task's tests with `flai test` on the paths it changed. `tooling.md` says `flai test` is the entry point for the tests and lint of what changed, and that the `Makefile`'s test targets stay for people and CI. It needs a flai that has `flai test`.

## 1.0.64 - 2026-10-06

- S-0300 An epic's planner drafts the tasks of each story it drafts (patch): `strategic-agents.md`'s section "As the planner" says that for an epic it enriches each story it drafts as it would a story, and drafts its tasks as it would for a story with no tasks, and does the same for each story it revisits that is a draft with no tasks. The tasks carry no draft flag: they are drafts because their story is one, and the story's agent reviews them when it pulls the story. The epic's one plan thread names each story's tasks and their layers, and the summary names by ID the stories and tasks created and the stories revisited. `.claude/agents/planner.md` says the same, and fixes "an story's" in its description. The planner's prompt needs a flai that says so.

## 1.0.63 - 2026-10-06

- S-0229 The orchestrator reads its permissions, policy, and release policy again before each decision (patch): `strategic-agents.md`'s section "As the orchestrator" says the operator may change any of them while it runs, not only `answer_threads`, as they now can from the dashboard's Orchestrator page through `flai manifest set` (ADR-0101). `flai guard` already holds each call to the permissions as they are then. The dashboard's settings panels need a flai and a flaiover that have them.

## 1.0.62 - 2026-10-06

- S-0224 The analyzer files each actionable finding of its report as an issue (patch): `.claude/agents/analyzer.md` gains the MCP tools `issue_new` and `issue_bump` and says how. It lists the open issues with `flai issue list --json`, or reads `design/issues/summary.md`; when one already records the finding, under whatever title, it bumps it with `flai issue bump <id> --report <its report>` and the finding's impact; otherwise it files one with `flai issue new "<title>" --class <class> --report <its report> --json`, the class from the finding (`defect`, `efficiency`, or `impression` for a risk with no measured instance), its impact as `--time-lost-per-cycle`, `--revenue-per-week`, or `--penalty-per-week`, and `--evidence`. With `--report`, flai bumps an open issue of the same title rather than opening a second, and the issue's Remediation section links the report. The report links each issue it filed or bumped from its finding, as `[I-nnnn](../issues/<file>)`. It never runs `flai issue story` or `flai issue close`, edits no file under `design/issues`, and still authors no stories; `flai guard` lets an analyzer's session run `flai issue new`, `bump`, and `list` and call `issue_new` and `issue_bump`, and refuses it `flai issue story`, `flai issue close`, `flai story`, `epic`, and `task`, and the MCP tool `issue_story`. `design/conventions/strategic-agents.md` says the same under `As the analyzer`: file each finding with `--report`, link each issue from its finding, and leave the issues to flai. The commands need a flai that has them.

## 1.0.61 - 2026-10-06

- S-0223 The analyzer is an agent flai serve starts for a project behind the analyze host action, on the operator's word or on `analysis.schedule`, and it writes one report under `design/analysis/` (patch): the template ships `.claude/agents/analyzer.md`, the definition the analyzer's session runs as. It primes with role `analyze`, reads the metrics with `flai stats --json`, the design with `doc_search` and `doc_get`, and the issues, hands wide search to the explorer, and looks for bottlenecks in the flow, gaps between `design/system` and the code, or technical and security risks, as its focus says, or for all three. It writes one report, `design/analysis/<date>-<focus>.md`, with front matter `title`, `updated`, `status`, `focus`, `from`, and `to`, and one section per finding with its evidence, severity, and estimated impact, and adds it to `design/analysis/README.md`. It edits nothing else, writes no work item, and authors no stories, and its last line names the report. It has `Edit` and `Write`, and no tool that writes a work item. `.claude/settings.json`'s `PreToolUse` entry on `Edit`, `Write`, and `NotebookEdit` runs `flai guard` in an analyzer's session (`FLAI_ROLE=analyze`) as well as a planner's and an orchestrator's, so that its edits stay inside `design/analysis/`, and still exits at once in a story's. A new project has `design/analysis/README.md`, which says what the folder is for, how a report is named, its front matter and statuses, and has a table with a row per report; `design/README.md` and `CLAUDE.md` list the folder in their layout tables. The commands need a flai that has them.

## 1.0.60 - 2026-10-06

- S-0295 One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three (patch).

## 1.0.59 - 2026-10-06

- S-0295 The planner and the story's agent name files, not folders, in touches (patch): `strategic-agents.md`'s section "As the planner" says to predict a story's touches and its tasks' file by file, to keep a folder touch only where the story may add files there that no task can name yet, and to record each folder touch kept, and why, under the story's `### Planning` heading. `work-management.md` says to name files in touches wherever you can, to narrow a story's folder touch to the files its tasks name with `flai touches`, or leave it for flai, which replaces it in the story's claim with its tasks' touches inside it (ADR-0096), and never to ask for a path to be added to the shared paths (`claims.shared`) to escape a hold, since that list is the operator's. The narrowing and the shared paths need a flai that has them.
- S-0295 A new project lists its shared paths in the manifest (patch): `system-flow.yaml` carries `claims.shared`, the glob patterns of paths that many stories change in separate sections or new files, starting with `<design>/adrs` and `<design>/issues` under the design layout folder. An overlap of two stories' claims that lies wholly inside a pattern holds no ready story and is no `wip.overlap` (ADR-0096). `flai shared` lists, checks, adds, and removes the patterns; only the operator's own session changes them. An older flai ignores the key and holds as before, so it needs a flai that has it.

## 1.0.58 - 2026-10-06

- S-0221 The orchestrator accepts a story in review when `accept_reviews` allows it (patch): `strategic-agents.md`'s section "As the orchestrator" says how. It hands the story's worktree to its verifier, which runs the tests, the lint, and `flai check --strict`, matches each acceptance criterion to the changed files that meet it, and names the commit it verified. It runs `flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run` and reads the blockers. With none, and every criterion matched, it runs `flai accept` with `--verified` and `--evidence -`, a `Verdict:` line and one `- <n>: <files>` item per criterion on standard input, and logs the acceptance with `activity_log`, naming the story and the commit. Otherwise it leaves the story in review, opens a thread on it saying what is missing, and logs that decision. It never moves a story to `done` with `item_move`. `.claude/agents/orchestrator.md` says the same. The commands need a flai that has them.

## 1.0.57 - 2026-10-06

- S-0285 A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout (patch).

## 1.0.56 - 2026-10-06

- S-0285 A story's agent waits for its sub-agents in the foreground, and flai guard refuses it wait_for_events while one runs (patch): `delegation.md` says to launch each sub-agent in the foreground, a layer's in one message, so that each result comes back as the tool's result however long the sub-agent runs (in Claude Code, the Agent tool's `run_in_background` set to false); never to end the turn while a sub-agent runs in the background, since a headless Claude Code session ends ten minutes after its turn ends, and the sub-agent with it; and never to wait for a sub-agent with `wait_for_events`, which reports work items and threads, not sub-agents, and is for a thread awaiting the designer; `flai guard` refuses it to a story's agent while a sub-agent of its session runs and no thread on the story or one of its tasks is open. `.claude/settings.json` runs `flai guard` on `SubagentStart` and `SubagentStop` in a story's agent's session (`FLAI_STORY` set), silently and never refusing, to record the session's running sub-agents. The commands need a flai that has them.

## 1.0.55 - 2026-10-06

- S-0220 The orchestrator answers threads, or recommends an answer, as its permission allows (patch): `strategic-agents.md`'s section "As the orchestrator" says what it does with the threads awaiting the operator under `answer_threads`, whose value it reads afresh each time. With `off` it leaves them alone. With `recommend` it replies to each with `thread_reply`, `recommendation: true`, and a `source`, the ADR, design section, or convention its answer rests on. With `autonomous` it answers with a `source` when one settles the question, and posts a recommendation instead, escalating to the operator, when none does or the question asks for the operator's judgement: a decision not yet recorded, a change of scope, or money. It never resolves a thread it did not open, answers one it opened, or confirms a recommendation. `.claude/agents/orchestrator.md` says the same. The commands need a flai that has them.

## 1.0.54 - 2026-10-06

- S-0282 Acceptance criteria are ticked through flai once verified (patch): `work-management.md` says to tick each criterion as soon as you have verified it, with `flai criteria tick <story> <n>` or the MCP tool `criteria_tick` and the numbers `flai criteria list` prints, never by editing the item's file, and that the move to `review` warns of an unticked one. `delegation.md` says a task sub-agent names in its final message the criteria its task meets, a verifier says which criteria the diff meets and which it does not, and the story's agent ticks them after its review. The commands need a flai that has them.

## 1.0.53 - 2026-10-06

- S-0219 The orchestrator plans, finalizes, promotes, and orders by its permissions (patch): `strategic-agents.md`'s section "As the orchestrator" says what it does with each. With `plan_backlog_epics` it runs `flai plan --candidates` and starts the planner (`plan`) for each epic listed, one at a time. With `finalize_drafts` it runs `flai promote --drafts`, finalizes a complete draft whose criteria, touches, forecast, and value it judges consistent, and opens one thread on any other saying what is missing or inconsistent. With `promote_to_ready` it moves `flai promote --candidates`'s candidates to `ready` in their order while the ready limit has room, never a draft or a held story. With `order_ready` it runs `flai order --by <policy> --apply` after each change to the ready column, which keeps a story placed by hand within the last day. It logs every action with `activity_log` and the policy figure that justified it, and a refusal from `flai guard` or from flai ends that attempt, which it logs and does not retry until something changes. `.claude/agents/orchestrator.md`, the orchestrator's definition, says the same. The commands need a flai that has them.

## 1.0.52 - 2026-10-05

- S-0218 The orchestrator is an agent flai serve runs for a project behind the orchestrate host action (patch): the template ships `.claude/agents/orchestrator.md`, the definition the orchestrator's session runs as, which acts only within `orchestration.permissions`, takes every order, candidate, and release figure from `order_by_policy`, `promote_candidates`, and `release_evaluate`, logs each decision with `activity_log`, waits on `wait_for_events` between decisions, and writes work items and threads through flai only. `.claude/settings.json`'s `PreToolUse` entry on `Edit`, `Write`, and `NotebookEdit` runs `flai guard` in an orchestrator's session (`FLAI_ROLE=orchestrate`) as well as a planner's, and still exits at once in a story's.

## 1.0.51 - 2026-10-05

- S-0244 Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes (patch).

## 1.0.50 - 2026-10-05

- S-0244 Two stories in progress whose claims grew to overlap are told so, and the story agent waits for its sub-agents through the harness (patch): `work-management.md` says that an `overlapped` change whose `cause` is a story in progress, not an accepted one, means a write grew the two stories' claims until they overlap; coordinate with that story's agent on a thread while the work is small and narrow your `touches` if you can, else note the overlap in the narrative and go on, since a `wip.overlap` does not stop the close-out. `delegation.md` says to wait for a sub-agent run in the background through the harness's notice that it has finished, not by polling `wait_for_events`, which reports work items and threads, not sub-agents.

## 1.0.49 - 2026-10-05

- S-0268 The story agent's prompt says to batch independent edits and commands in one turn, to take a new task through ready in one command, and to prove a test fails without the shared stash (patch).

## 1.0.48 - 2026-10-05

- S-0249 flai check --strict stops a story's close-out on wip/ findings outside the story (patch).
- S-0266 The verifier reads a close-out run's exit status without piping it away, and the story agent tells it where the run is expected to stop (patch).

## 1.0.47 - 2026-10-05

- S-0249 flai check --strict stops a story's close-out on wip/ findings outside the story (patch): `scripts/close-out.sh` exports the story as `CLOSE_OUT_STORY`, and `scripts/check.sh` then runs `flai check --strict --story "$CLOSE_OUT_STORY" --record-issues`. A finding outside the story, and every `wip.overlap`, is printed as a note that does not stop the close-out, and each rule's are recorded in one open issue under `design/issues`, opened or bumped, which the close-out's commit step commits. Without `CLOSE_OUT_STORY`, as in CI and `make`, the check is unscoped and every finding counts. `work-management.md`'s close-out rule says the same. The scripts need a flai that has `--story` and `--record-issues`.

## 1.0.46 - 2026-10-04

- S-0210 The planner enriches a story with predicted touches, a forecast, and a cost of delay value (patch).

## 1.0.45 - 2026-10-04

- S-0210 The planner enriches a story with predicted touches, a forecast, and a cost of delay value (patch): `strategic-agents.md`'s section "As the planner" says to predict a story's `touches` from `flai touches suggest`, its goal and criteria, the design it links, and the code layout, keeping every touch it declares; to take its forecast from `flai forecast` and its cost of delay `value` from `flai cod`; to review each figure and adjust it with a stated reason; and to record where each touch came from and why each figure stands under a `### Planning` heading in the story's Notes, which is the planner's to rewrite. `.claude/agents/planner.md`, the planner's definition, says the same. The three commands print and write nothing, and need a flai that has them.

## 1.0.44 - 2026-10-04

- S-0255 The planner drafts a story's tasks into the backlog and revisits the children it already has (patch): `strategic-agents.md` says the planner drafts a story's tasks in layers, each with `## Work`, `## Done when`, a nature, tags, `touches`, and `after`, revisits the open ones, proposes splits, merges, and drops on the plan's thread, never cancels or rewrites a task it did not write without asking, adds a task's topic to the story, and names the items it created and revisited in its summary. `work-management.md` says the planner may draft a story's tasks, and that the agent that pulls the story reviews them before it works them.

## 1.0.43 - 2026-10-04

- S-0209 The planner drafts an epic's stories into the backlog and revisits the children it already has (patch): `strategic-agents.md`'s section "As the planner" says to create each story as a draft (`draft: true`), since flai refuses a planner's story that is not one, and to write every section of the item's template, since flai check refuses a story missing one; to summarise the plan in one thread on the item, naming on an epic the stories, their order, and the assumptions made; for an epic with stories, to enrich each one not done or cancelled again and propose in that thread each story to split, merge, add, or drop, drafting the additions only; and to end with one line naming the stories created and revisited. `.claude/agents/planner.md`, the planner's definition, says the same.

## 1.0.42 - 2026-10-04

- S-0208 The planner is an agent flai serve starts for an epic or a story, behind the plan host action (patch).

## 1.0.41 - 2026-10-04

- S-0208 The planner is an agent flai serve starts for an epic or a story, behind the plan host action (patch): the template ships `.claude/agents/planner.md`, the definition `flai plan` runs the planner's session as, which writes work items and threads through flai only. `.claude/settings.json` keeps its hook on `Bash` and flai's MCP tools and gains a second `PreToolUse` entry that runs `flai guard` before `Edit`, `Write`, and `NotebookEdit` in a planner's session alone (`FLAI_ROLE=plan`), so the guard can refuse the planner a file edit while a story's session runs no hook per edit (ADR-0082, the operator's choice on TH-0096).

## 1.0.40 - 2026-10-03

- S-0206 Planner, orchestrator, and analyzer runs are logged in an activity document, with cost, seconds, and work completed in its front matter (patch): `wip/agents/README.md` names the three activity documents flai writes, and `strategic-agents.md` says `activity_log` takes the kind, the summary, and the items touched, that flai measures the duration and cost, and that a run's final reply summarises what it did after its last reported activity (ADR-0079).

## 1.0.39 - 2026-10-03

- S-0200 An epic follows its stories: to ready and in-progress with the first, to review and done with the last (patch): `git.md`'s release rule says an epic contributes no bump of its own, its stories carry theirs, now that an epic is accepted with its last open story (ADR-0078, the designer's choice on TH-0082).

## 1.0.38 - 2026-10-03

- S-0207 The conventions and prime know the planner, orchestrator, and analyzer as roles (patch): `strategic-agents.md` says what the planner, the orchestrator, and the analyzer do, never do, log, and ask, read by the roles `plan`, `orchestrate`, and `analyze`. The baseline's `roles` name the new roles where they apply, as the designer decided on TH-0081: `work-management.md` adds `plan` and `orchestrate`, `decisions.md` `orchestrate`, `continuous-improvement.md` and `code-quality.md` `analyze`, and `telemetry.md`, read by every agent until now, lists `story`, `explore`, `verify`, and `analyze`.

## 1.0.37 - 2026-10-03

- S-0199 Work items carry planning data: draft, cost of delay inputs and value, and a forecast with who set it (patch): `work-management.md`'s definition of ready says a draft story is finalized before it is ready (ADR-0074).

## 1.0.36 - 2026-10-03

- S-0243 Review column exceeds its WIP limit while acceptance is batched (patch).

## 1.0.35 - 2026-10-03

- S-0243 Review column exceeds its WIP limit while acceptance is batched (patch): `work-management.md` says `wait_for_work` names a story to pull only while review is under its limit as well as the in-progress limit, and that on a timeout it may be waiting for a story in review to be accepted or sent back (ADR-0073, with the designer's consent on TH-0079).

## 1.0.34 - 2026-10-03

- S-0197 Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve (patch).
- S-0242 A publish refused because the remote moved drops its tags, takes the remote by rebase or merge, and works conflicts through tasks and threads (patch).
- S-0198 Give the operator the option to have all issues turned into stories (patch).

## 1.0.33 - 2026-10-03

- S-0197 Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve (patch): `git.md` says to use `flai stream sync` for the branch's git operations and never start a `git rebase` or `git merge` by hand, finishing a rebase the sync stopped on with `git add` and `git rebase --continue`, or undoing it with `git rebase --abort`; to commit each task when it is done, then sync and resolve, then run its tests and commit any fix; and before review to commit what is outstanding, sync again, and close out. `work-management.md` points to that cycle for each task, before review, and on an `overlapped` change, and its definition of done asks for the branch synced after its last commit. `delegation.md` syncs again before the verifier runs. `scripts/close-out.sh` stops while a rebase is unfinished, and when the branch does not contain the main branch, saying to run `flai stream sync`: before the tests when there is nothing to commit, after the commit otherwise.

## 1.0.32 - 2026-10-02

- S-0194 An experiment story is accepted with its results document under design/experiments (patch).

## 1.0.31 - 2026-10-01

- S-0189 Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs (patch).

## 1.0.30 - 2026-10-01

- S-0189 Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs (patch): `.claude/agents/explorer.md` runs `haiku` and `verifier.md` runs `sonnet`, cheaper than the story's agent, since neither decides anything. `delegation.md` says the story's agent runs only the tests for what it changed and leaves the whole suite, the lint, and `flai check` to one verifier before review, through the close-out script, and one more after fixing what that one found; it fixes what a verifier finds itself, never through a sub-agent, and does not repeat a verifier's passing run.

## 1.0.29 - 2026-10-01

- S-0187 The conventions say how agents write shell commands for zsh and gate each step of a chain (patch): `tooling.md` says the host shell may be zsh, so quote globs, URLs, and variables, use `[ a = b ]`, and put long sequences in a script, and that each step of a chain that must succeed is joined with `&&` or runs under `set -e`, never `;` or a pipe that hides its exit code. `scripts/close-out.sh` runs the lint, the three test tiers, `flai check --strict`, the narrative check, and the commit before review, stopping at the first failure, and `work-management.md` points to it; `make lint-md` calls the new `scripts/lint-md.sh`.

## 1.0.28 - 2026-10-01

- S-0179 flai lints the markdown it writes in the main checkout, and a thread's entries in one second share a heading (patch).

## 1.0.27 - 2026-10-01

- S-0179 flai lints the markdown it writes in the main checkout, and a thread's entries in one second share a heading (patch): `tooling.md` says flai lints what it writes in `wip/` with the project's markdownlint configuration, warning in `flai check` and refusing a write that would bring a finding (ADR-0061).

## 1.0.26 - 2026-10-01

- S-0175 Agents the claude-code adapter starts hand noisy work and verification to sub-agents that cannot act for the story (patch).

## 1.0.25 - 2026-10-01

- S-0175 Agents the claude-code adapter starts hand noisy work and verification to sub-agents that cannot act for the story (patch): `.claude/agents/explorer.md` and `verifier.md`, sub-agents whose tools only read; `.claude/settings.json`, which runs `flai guard` before a sub-agent's shell and flai calls; a baseline `delegation.md` convention; and `roles` on the conventions each sub-agent reads with `flai prime --role` (ADR-0059, ADR-0060).

## 1.0.24 - 2026-09-29

- S-0149 Provide briefs in cases where large files are part of context (patch).

## 1.0.23 - 2026-09-29

- S-0149 Provide briefs in cases where large files are part of context (patch): CLAUDE.md and session-start.md say a large file a story names only by its path written out is briefed (ADR-0050).

## 1.0.22 - 2026-09-29

- S-0148 Agents flai starts, and any agent with a story, prime with the budgeted pack and fetch what it briefs (patch).

## 1.0.21 - 2026-09-28

- S-0134 Conventions, design, tech files, and ADRs carry topics on the file and on headings, and flai check keeps them honest (patch).

## 1.0.20 - 2026-09-27

- S-0140 Agents should not leave uncommitted work in their branch when moving to review (patch).

## 1.0.19 - 2026-09-26

- S-0132 Accepting a story tells every open story that overlaps it which paths changed (patch).

## 1.0.18 - 2026-09-23

- S-0097 Agents should connect to the MCP and consistently monitor for ready tickets (patch).

## 1.0.17 - 2026-09-21

- S-0087 Moving a story to done merges it; a publish button on the board tags and pushes everything accumulated since the last one (patch).

## 1.0.16 - 2026-09-20

- S-0079 A story moved to ready starts an agent on the host, with the command the operator configured (patch).

## 1.0.15 - 2026-09-20

- S-0076 flai serves MCP itself on the host, over HTTP as well as stdio, and the dashboard's /mcp endpoint goes (patch).

## 1.0.14 - 2026-09-20

- S-0070 Cancelling an epic cancels its open stories and their tasks, and cancelling a story cancels its open tasks (patch).

## 1.0.13 - 2026-09-19

- S-0060 Operators create ADRs from the dashboard (patch).

## 1.0.12 - 2026-09-19

- S-0063 An acceptance that is not pushed stays visible on the board and in inbox, and flai push --pending pushes it (patch).

## 1.0.11 - 2026-09-19

- S-0053 Research and experiment stories can be accepted: findings land on main and are pushed, with no release (patch).

## 1.0.10 - 2026-09-19

- S-0057 Cards can be dragged within a column to change the pull order (patch).

## 1.0.9 - 2026-09-19

- S-0058 An agent learns through MCP that the designer moved work to ready (patch).

## 1.0.8 - 2026-09-18

- S-0051 Accepting from the dashboard shows uncommitted changes outside wip and lets the designer include them (patch).

## 1.0.7 - 2026-09-18

- S-0049 A story is ready without tasks; the agent that pulls it writes them (patch).

## 1.0.6 - 2026-09-18

- S-0039 flai mcp: inbox, reply, read, and move tools with the template .mcp.json (patch).

## 1.0.5 - 2026-09-18

- S-0038 Comment threads and answers as files under wip (patch).

## 1.0.4 - 2026-09-18

- S-0037 Story branches with wip on main, stream sync, and the touches flag (patch).

## 1.0.3 - 2026-09-17

- S-0035 flai dashboard runs from a private registry or a local build (patch).

## 1.0.2 - 2026-09-17

- S-0033 Four-digit work item IDs (patch).

## 1.0.1 - 2026-09-17

- S-0021 flai template push publishes template changes to a remote (patch).

## 1.0.0 - 2026-09-16

Epic E-0005 (agent conventions) complete: twelve baseline conventions under
`design/conventions`, `design/issues` with `flai issue`, `scripts/` behind
the Makefile with three test tiers, and a `CLAUDE.md` that primes every
session from the conventions. First major release of the template.

## 0.4.0 - 2026-09-16

- `design/conventions/logging.md` and `telemetry.md` are active: five log levels including fatal, structured events, `/_health`, `/_ready`, `/metrics`, OpenTelemetry, golden signals (S-0030).

## 0.3.0 - 2026-09-16

Note: under the refined release rule (incidental additive touches are a patch) this would have been 0.2.1; the version stands as cut.

- `design/conventions/logging.md` and `telemetry.md` added as drafts, indexed at 110 and 120 (S-0025).
- Markdownlint excludes test fixtures and `bin/`.

## 0.2.0 - 2026-09-16

- `design/conventions/`: ten baseline agent norms with a project-additions marker (S-0023, S-0026).
- `design/issues/` skeleton and `scripts/` with test tier stubs and Makefile delegation (S-0026).
- `CLAUDE.md` opens with a priming section and is a map; norms live in the conventions (S-0024).
- Rendered projects are markdownlint-clean; lint config aligned with the item format.

## 0.1.0 - 2026-09-15

Initial prototype. Layout, work item schema, baseline CLAUDE.md, devex files.
