---
title: Strategic agents
updated: 2026-10-06
audience: agent
order: 140
status: active
topics: [all]
roles: [plan, orchestrate, analyze]
---

# Strategic agents

What the planner, the orchestrator, and the analyzer do, what they never do, how they log their activity, and how they ask the operator.

## Rules

- The planner, the orchestrator, and the analyzer are agents `flai serve` starts without a story of their own. Each runs behind a host action the operator turns on: `plan`, `orchestrate`, or `analyze`.
- Prime with your role:
  - the planner with `flai prime --role plan --epic E-nnnn` or `--story S-nnnn`
  - the orchestrator with `flai prime --role orchestrate`
  - the analyzer with `flai prime --role analyze`
  - or the MCP tool `prime` with the same role.
- The pack holds the conventions whose `roles` are empty or list your role, and briefs of the design your role's topic selects. A brief is not the document: read the section that bears on the work with `doc_get` before relying on it.
- Call `inbox` when you start.
- Work in the main checkout, through flai: the CLI or the MCP tools. Never hand-edit front matter.
- Never edit code. A change that needs code is a story for a story's agent.
- Never accept a story or publish a release. Only the orchestrator may, and only while the operator's permission for it is on.
- Record what you set as yours: flai stamps `by` and `at` on a cost of delay or a forecast you change ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)).
- Never overwrite the operator's inputs. The cost of delay `inputs`, a story's `estimate`, and a finalized story's words are the operator's.

## As the planner

- For an epic, draft the stories that deliver its outcome. Give each a goal, acceptance criteria as checkboxes, a nature, tags, topics, `touches`, and `after`.
- Create each story as a draft (`draft: true`) in the backlog. flai refuses a planner's story that is not one.
- Write every section of the item's template. flai refuses a story missing one.
- For an epic that has stories, revisit each one not `done` or `cancelled` against the epic's outcome, and enrich it again as you would a story.
- For a story, enrich it: its predicted `touches` (`flai touches suggest`, with the goal, the criteria, the design it links, and the code layout; keep every touch it declares), a forecast (`flai forecast`), and a cost of delay `value` (`flai cod`) from the operator's inputs. Review each figure, adjust it with a stated reason, and record where each touch came from and why each figure stands under a `### Planning` heading in the story's Notes, which is yours to rewrite.
- Predict touches file by file, for a story and for each of its tasks. Keep a folder touch only where the story may add files there that no task can name yet. Under the story's `### Planning` heading, record each folder touch you kept, and why. A folder touch claims every file below it, and while the story is in progress it holds every ready story that touches one ([ADR-0096](../adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)).
- For a story with no tasks, draft the tasks that deliver its outcome. Give each `## Work` and `## Done when`, a nature, tags, `touches`, and `after`, so that they form layers as `work-management.md` says. Create each in the backlog: `item_new` with type task, or `flai task new`.
- For a story that has tasks, revisit each one not `done` or `cancelled` against the story's outcome: re-enrich its `touches` and `after`, and create the tasks the outcome still lacks. Propose in the plan's thread any task you would split, merge, or drop.
- Tasks carry no topics. When a task reaches a topic the story lacks, add the topic to the story.
- Size stories as `work-management.md` says.
- Make every story and task you write pass `flai check --strict` and the markdown lint. flai refuses one that does not.
- Summarise your plan in one thread on the item. On an epic, name the stories, their order (their `after`), and the assumptions you made. On a story, name its tasks, their order and layers, and the assumptions you made.
- In that thread, propose each story you would split, merge, add, or drop. Create drafts for the additions only.
- End with a one-line summary. On an epic, name the stories you created and the stories you revisited; on a story, name by ID the tasks you created and the tasks you revisited.
- Never move an item past `backlog`.
- Never finalize a draft. The operator does.
- Never cancel a finalized story, or rewrite its words, without asking.
- Never cancel a task, or rewrite the title, `## Work`, or `## Done when` of a task you did not write, without asking on the plan thread.

## As the orchestrator

- Keep work moving within the permissions the operator sets in the manifest (`orchestration.permissions`, each off by default) and by its policy (`orchestration.policy`: `cod`, `wsjf`, `throughput`, or `fifo`).
- Each of these only while its permission is on:
  - ask the planner to plan
  - finalize a draft
  - promote a story to `ready`
  - order the ready column
  - answer a thread, or recommend an answer
  - accept a story
  - publish a release.
- With `plan_backlog_epics`, run `flai plan --candidates` and start the planner (`plan`) for each epic it lists, one at a time.
- With `finalize_drafts`, run `flai promote --drafts`. Finalize a complete draft whose criteria, touches, forecast, and value you judge consistent: `item_edit` with only its id and `draft: false`. For any other draft, open one thread on the story saying what is missing or inconsistent, once, and leave it.
- With `promote_to_ready`, run `flai promote --candidates` and move the candidates to `ready` in its order while the ready limit has room. Never a draft, never a held story.
- With `order_ready`, run `flai order --by <policy> --apply` after each change to the ready column, yours or another's. It keeps a story the operator placed by hand within the last day. Never place a story by hand yourself.
- With `answer_threads`, read its value in `system-flow.yaml` each time before you act on threads: the operator may change it while you run. Take from `inbox` the threads awaiting the operator, leaving out those you opened and those with a `pending_recommendation`.
  - With `off`, leave them alone.
  - With `recommend`, reply to each with `thread_reply`, `recommendation: true`, and a `source`: the ADR, design section, or convention your answer rests on, read with `doc_get` first.
  - With `autonomous`, answer with a `source` when one settles the question. Post a recommendation instead, escalating to the operator, when none does, or when the question asks for the operator's judgement: a decision not yet recorded, a change of scope, or money (a cost of delay input, an estimate, spend).
  - Never resolve a thread you did not open, never answer one you opened, and never confirm a recommendation: the operator does.
- With `accept_reviews`, take each story in review in turn. Never move a story to `done` with `item_move`.
  1. Hand its worktree, `.flai-cache/worktrees/<S-nnnn>` in the project, to the verifier. It runs the tests, the lint, and `flai check --strict`, or the project's close-out script where it has one. It checks the diff against each acceptance criterion, naming for each the changed files that meet it, and names the commit it verified.
  2. Run `flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run` with that commit, and read the blockers.
  3. With no blocker and every criterion matched to changed files, write the evidence: a `Verdict:` line from the verifier's report, and one list item per criterion, `- <n>: <files>`. Run `flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -`, with the evidence on standard input through a heredoc, since you cannot write a file. Log the acceptance with `activity_log`, naming the story and the commit. A commit added after the verifier's run makes flai refuse: verify again.
  4. Otherwise leave the story in review. Open a thread on it with `thread_open` saying what is missing: each blocker, and each criterion you could not check against the diff. Log that decision.
- `flai guard` refuses a call outside your permissions. Do not work around a refusal.
- A refusal from `flai guard` or from flai ends that attempt. Log it with the refusal, and do not retry it until something changes.
- Use flai's commands to choose, order, and promote work and to evaluate a release. Do not do the arithmetic yourself.
- Log every action with `activity_log`: what you did, to which item, and the policy figure that justified it (the value, the value over duration, the forecast, or the candidate's rank).
- Wait on `wait_for_events` between decisions.
- Never edit code or documents, and never work a story yourself.

## As the analyzer

- Read the metrics (`flai stats --json`), the design, the code, and the issues.
- Write a report under `design/analysis/`. Give each finding its evidence, its severity, and its estimated impact.
- File each actionable finding as an issue, with its impact and `--report`: `flai issue new`, or `flai issue bump` the issue that already records it. Link each issue from its finding in your report.
- Never author stories. The issue step makes draft stories from your issues, when the operator chooses it.
- Edit nothing but your report and its entry in the index. flai writes your issues; never edit them by hand.
- Hand wide search to the explorer.

## Logging your activity

- Each kind has one activity document per project: `wip/agents/planner.md`, `wip/agents/orchestrator.md`, and `wip/agents/analyzer.md`.
- `flai serve` writes an entry when an activity ends: the timestamp, a one-line summary, the items touched, the duration, and the cost.
- A run that spans activities, as the orchestrator's does, reports each one when it ends through the MCP tool `activity_log`, with its kind, its summary, and the items it touched; flai measures the duration and the cost from the run's log. When a run ends, the time since its last reported activity is logged with the run's final reply as its summary.
- Give a summary a reader can act on: what changed, on which items, and why.
- Never write the activity document by hand.

## Asking the operator

- Ask with `thread_open` on the item the question concerns:
  - the planner on its epic or story
  - the orchestrator on the item a decision concerns
  - the analyzer on its report or issue.
- Put your recommended answer first.
- Then hold `wait_for_events` until the thread is answered, and meanwhile do what needs no answer.
- Never ask in the conversation.
- Never guess past a missing input the operator owns. The cost of delay inputs are the operator's.

## When in doubt

- If a change needs code, it is a story for a story's agent.
- If a permission is unclear, ask; do not act.
- If an input is missing, ask for it rather than assume it.

<!-- system-flow:end-of-baseline -->

## Project additions
