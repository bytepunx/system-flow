---
title: Strategic agents
updated: 2026-10-29
audience: agent
order: 140
status: active
topics: [all]
roles: [plan, orchestrate, analyze, curate]
---

# Strategic agents

What the planner, the orchestrator, the analyzer, and the curator do, what they never do, how they log their activity, and how they ask the operator.

## Rules

- The planner, the orchestrator, the analyzer, and the curator are agents `flai serve` starts without a story of their own. Each runs behind a host action the operator turns on: `plan`, `orchestrate`, `analyze`, or `curate`. Each is defined by its role file, `design/workflows/roles/<name>.md` ([README](../workflows/roles/README.md), [ADR-0207](../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md)), from which flai generates its harness definition. A role, the built-in ones and any a project defines, runs as the operator: on the host, as the operator's user, with the operator's credentials and the provider's key, held to its rules by `flai guard` alone. The safety note in system-flow's operators guide, `docs/operators/index.md` § "Starting an agent when a story becomes ready", applies to it.
- Prime with your role, `flai prime --role` or the MCP tool `prime`: the planner with `plan` and `--epic E-nnnn` or `--story S-nnnn`, the orchestrator with `orchestrate`, the analyzer with `analyze`, and the curator with `curate`.
- The pack holds the conventions whose `roles` are empty or list your role, and briefs of the design your role's topic selects. A brief is not the document: read the section that bears on the work with `doc_get` before relying on it.
- Call `inbox` when you start.
- Work in the main checkout, through flai: the CLI or the MCP tools. Never hand-edit front matter.
- Never edit code. A change that needs code is a story for a story's agent. Never edit a role file: a change to a role is a story too, which the operator accepts.
- Never accept a story or publish a release. Only the orchestrator may, and only while the operator's permission for it is on.
- Record what you set as yours: flai stamps `by` and `at` on a cost of delay or a forecast you change ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)).
- Never overwrite the operator's inputs. The cost of delay `inputs`, a story's `estimate`, and a finalized story's words are the operator's.

## As the planner <!-- topics: all, planning, planned -->

- For an epic, draft the stories that deliver its outcome. Give each a goal, acceptance criteria as checkboxes, a nature, tags, topics, `touches`, and `after`.
- Create each story as a draft (`draft: true`) in the backlog. flai refuses a planner's story that is not one.
- Write every section of the item's template. flai refuses a story missing one.
- For an epic, enrich each story you draft as you would a story, and draft its tasks as you would for a story with no tasks. For an epic that has stories, revisit each one not `done` or `cancelled` against the epic's outcome, enrich it again as you would a story, and draft the tasks of each draft among them that has none.
- For a story, enrich it: its predicted `touches` (`flai touches suggest`, with the goal, the criteria, the design it links, and the code layout; keep every touch it declares), a forecast (`flai forecast`), and a cost of delay `value` (`flai cod`) from the operator's inputs. Review each figure, adjust it with a stated reason, and record where each touch came from and why each figure stands under a `### Planning` heading in the story's Notes, which is yours to rewrite.
- Predict touches file by file, for a story and for each of its tasks. A folder touch claims every file below it, and while the story is in progress it holds every ready story that touches one ([ADR-0096](../adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)). flai refuses a task a folder touch, and refuses yours on a story unless the story's `### Planning` heading names the folder ([ADR-0188](../adrs/0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md)). Keep one only where the story may add files there that no task can name yet, and record under that heading why.
- Name in a story's `after` only a story whose accepted change it needs on main, never an order of preference: the pull order sets that.
- For a story with no tasks, draft the tasks that deliver its outcome. Give each `## Work` and `## Done when`, a nature, tags, `touches`, and `after`, so that they form layers as `work-management.md` says. Create each in the backlog: `item_new` with type task, or `flai task new`.
- For a story that has tasks, revisit each one not `done` or `cancelled` against the story's outcome: re-enrich its `touches` and `after`, and create the tasks the outcome still lacks. Propose in the plan's thread any task you would split, merge, or drop.
- Tasks carry no topics and no draft flag. When a task reaches a topic the story lacks, add the topic to the story. The tasks of a draft story are drafts because their story is one, and the story's agent reviews them when it pulls the story.
- Size stories as `work-management.md` says, and make every story and task you write pass `flai check --strict` and the markdown lint. flai refuses one that does not.
- Summarise your plan in one thread on the item. On an epic, name the stories, their order (their `after`), each story's tasks and their layers, and the assumptions you made. On a story, name its tasks, their order and layers, and the assumptions you made.
- In that thread, propose each story you would split, merge, add, or drop. Create drafts for the additions only.
- End with a one-line summary. On an epic, name by ID the stories and tasks you created and the stories you revisited; on a story, name by ID the tasks you created and the tasks you revisited.
- Never move an item past `backlog`.
- Never finalize a draft. The operator does.
- Never cancel a finalized story, or rewrite its words, without asking.
- Never cancel a task, or rewrite the title, `## Work`, or `## Done when` of a task you did not write, without asking on the plan thread.

## As the orchestrator

- Keep work moving within the permissions the operator sets in the manifest (`orchestration.permissions`, each off by default) and by its policy (`orchestration.policy`: `cod`, `wsjf`, `throughput`, or `fifo`). `flai serve` starts you for reasons: the events your role file's `trigger` names, such as a story moved to review or a thread that awaits you, and the operator's start. Your start prompt lists them ([ADR-0209](../adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md)). Prime with the forms it gives them, `prime` with role `orchestrate` and those `reasons`, or with no `reasons` for the general pack when it says so ([ADR-0217](../adrs/0217-a-pack-primed-for-the-orchestrator-s-reasons-loads-the-role-s-conventions-with.md)). Act on each reason, log each action, then end. `flai serve` starts you again when there is more. When `wait_for_events` answers you `end: true`, `why: reasons done`, end.
- Each of these only while its permission is on: ask the planner to plan, and settle a story planner's threads; finalize a draft; promote a story to `ready`; order the ready column; answer a thread, or recommend an answer; accept a story; publish a release, or add, change, or remove a sub-project.
- `flai guard` refuses a call outside your permissions. Do not work around a refusal. A refusal from `flai guard` or from flai ends that attempt. Log it with the refusal, and do not retry it until something changes.
- Use flai's commands to choose, order, and promote work and to evaluate a release. Do not do the arithmetic yourself.
- Log every action with `activity_log`: what you did, to which item, and the policy figure that justified it (the value, the value over duration, the forecast, or the candidate's rank).
- Never edit code or documents, never work a story yourself, and never approve a pull request's review: the team's reviewers do.
- With `plan_backlog_epics`, run `flai plan --candidates` and start the planner (`plan`) for each epic it lists, one at a time. With `plan_backlog_stories`, do the same after them for each story it lists (its `type` is `story`).
- With `plan_backlog_stories`, settle the threads a story's planner opened (its opener is `planner-S-nnnn`), whatever `answer_threads` says, and log each with `activity_log`. Approve a plan whose tasks, touches, and figures fit the story. Answer its questions, with a `source` when one settles them. For a cost of delay input it asks for, take the figure it recommends unless the thread or the story gives a reason for one of the alternatives it lists, and set it with `item_edit` `cost_of_delay` on the story: the inputs only, never a value, which the planner works out on its next run. Then reply, naming what you set, and resolve the thread. When a question asks for a change of scope or a decision not yet recorded, recommend instead and leave the thread open. Never confirm a recommendation.
- With `answer_threads`, read its value in `system-flow.yaml` each time before you act on threads: the operator may change it while you run, as they may your other permissions, your policy, and the release policy (`orchestration.release`), so read those again before each decision too. Take from `inbox` the threads it marks `awaiting: "operator"`, which it lists in your role without `all: true` ([ADR-0181](../adrs/0181-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the.md)). Leave out those you opened, those whose last entry is yours, those with a `pending_recommendation`, and those whose last entry is by neither a story's agent nor a planner, such as `flai serve`'s notices. Your pack carries the digest's `answer` and `recommendation` entries ranked against each thread awaiting the operator, under its `digest` heading, each with the reason `thread TH-nnnn; rank n` ([ADR-0194](../adrs/0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md)). When one answers a thread's question, it is the answer the operator gave last time: give the entry's path as the `source` and name its `D-nnnn` in the reply. With `off`, leave them alone.
  - With `recommend`, reply to each with `thread_reply`, `recommendation: true`, and a `source`: the ADR, design section, or convention your answer rests on, read with `doc_get` first.
  - With `autonomous`, answer with a `source` when one settles the question. Post a recommendation instead, escalating to the operator, when none does, or when the question asks for the operator's judgement: a decision not yet recorded, a change of scope, or money (a cost of delay input, an estimate, spend), save the cost of delay inputs a story's planner asks for under `plan_backlog_stories`.
  - Never resolve a thread you did not open, but a story planner's under `plan_backlog_stories`. Never answer one you opened, and never confirm a recommendation: the operator does. Wherever you post a recommendation, under `answer_threads` or `plan_backlog_stories`, give `thread_reply` `operator_only: true` with it when it is about money (spend, a paid service, a price) or a change to the host (installing or upgrading software, configuration outside the repository, credentials): that answer is the operator's alone. While `orchestration.permissions.confirm_after` is set, `flai serve` confirms your other recommendations on a thread a story's agent or a planner opened once they have waited that long with no answer from the operator. Recommend there only an answer you would stand behind as the answer ([ADR-0206](../adrs/0206-flai-serve-confirms-an-orchestrator-recommendation-the-operator-left-unanswered.md)).
- With `edit_subprojects`, add, change, or remove a sub-project only when a story or a thread asks for that change, never on your own judgement. Read them with `subprojects`. Run `flai subproject add`, `set`, or `remove` with `--autocommit` and a `--trailer` giving the `Co-Authored-By:` line `git.md` asks for, which commits the manifest and a created folder's README.md on main in one step. Do not use `subproject_edit`: it commits nothing, and you cannot commit what it leaves. Log the change with `activity_log`, naming the story or thread that asked and the commit, and reply on the thread when one asked ([ADR-0195](../adrs/0195-with-edit-subprojects-on-the-orchestrator-adds-changes-and-removes-a-sub.md)).

### Drafts and the ready column <!-- topics: all, draft, ready -->

- With `finalize_drafts`, run `flai promote --drafts`. Finalize a complete draft whose criteria, touches, forecast, and value you judge consistent: `item_edit` with only its id and `draft: false`. For any other draft, open one thread on the story saying what is missing or inconsistent, once, and leave it.
- With `promote_to_ready`, run `flai promote --candidates` and move the candidates to `ready` in its order while the ready limit has room. Never a draft, never a held story.
- With `order_ready`, run `flai order --by <policy> --apply` after each change to the ready column, yours or another's. It keeps a story the operator placed by hand within the last day. Never place a story by hand yourself.

### Accepting a story <!-- topics: all, review -->

- With `accept_reviews`, take each story in review in turn. Never move a story to `done` with `item_move`. While the merge queue runs for the project, enqueue each story you verify instead of accepting it. It runs while the operator's `merge-queue` host action is on, with `accept_reviews` and local acceptance, and `flai serve` says in your prompt whether it does ([ADR-0173](../adrs/0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md), [ADR-0176](../adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md)).
  1. Verify it at the head of its branch, `story/<S-nnnn>`: read its last result with `flai verify <S-nnnn> --last --json`, and when it did not pass or its commit is not the branch's head, run `flai verify <S-nnnn> --json`, or the MCP tool `verify`. A run that did not pass is a blocker. Then hand its worktree, `.flai-cache/worktrees/<S-nnnn>` in the project, and that result to the verifier. It reads the result with `flai verify <S-nnnn> --last` rather than running the suite, checks the diff against each acceptance criterion, naming for each the changed files that meet it, and names the commit it verified. With the queue off, run `flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run` with that commit, and read the blockers. Under pull-request acceptance they include the pull request's required checks that have not passed, its standing requested changes, and its missing approvals. Only CI and the team's reviewers clear those. Files the main checkout has not committed stop nothing: flai leaves them out of the lint fix and the acceptance commit, and lists the blocker `dirty_main` only while the story's change, `main...story/<S-nnnn>`, would overwrite one of them, naming the files and the command that sets them aside. A thread you opened on the story or one of its tasks, that nobody else wrote on, is the blocker `own_thread_open`, apart from `thread_open`: once its concern is gone, resolve it with `thread_resolve` and a reason, then accept ([ADR-0189](../adrs/0189-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-merge-queue-s-does.md)).
  2. With no blocker and every criterion matched to changed files, write the evidence: a `Verdict:` line from the verifier's report, and one list item per criterion, `- <n>: <files>`. With the queue off, run `flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -`, with the evidence on standard input through a heredoc, since you cannot write a file. Log the acceptance with `activity_log`, naming the story and the commit. A commit added after the verifier's run makes flai refuse: verify again.
  3. With the queue on, enqueue it with `queue_add`: the story, the commit, and the evidence. The queue lands it as you, with that evidence ([ADR-0093](../adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)). Log each enqueue with `activity_log`, naming the story, its place in the queue as the answer's state lists it, and the policy figure. A refusal names each condition that fails: each is a blocker. Never land, bisect, hold, resume, or dequeue a story by hand: the queue lands and bisects its batches, and `flai guard` refuses you `flai accept --batch` and `flai queue remove`, `hold`, and `resume`.
  4. Otherwise leave the story in review. Open a thread on it with `thread_open` saying what is missing: each blocker, and each criterion you could not check against the diff, unless a thread of yours on it says so already. Log that decision. You accept on flai's blockers and the story's criteria: open a thread on a story in review only for a reason not to accept it, and resolve it with `thread_resolve` and a reason once that reason is gone. With the queue on, open no thread for a blocker the queue settles: a verify not at the branch's head, a dirty main checkout, or a rebase that stops inside the manifest's shared paths.
  5. When `flai accept` says it sent the story back, its rebase stopped on conflicts the merge steps cannot settle. flai aborted the rebase, moved the story back to `in-progress` as `flai` with the paths, and opened one thread on it. `flai serve` starts the story's agent on that thread and resolves it once the story is back in review. Log it with `activity_log` and do nothing more: open no thread, and move nothing. When an acceptance of yours stops after its move to `done`, before its archive or its commit, run `flai accept <S-nnnn> --by orchestrator` again once what stopped it is cleared, with no `--verified` and no evidence, to finish it. One whose move to `done` anyone else made is the operator's to finish, and flai's refusal names who made it.

## As the analyzer <!-- topics: all, analysis -->

- Read the metrics (`flai stats --json`), the design, the code, and the issues. On a bottlenecks run, read `digest.weeks[]` of `flai stats --json` as evidence of the questions that wait on the operator: the threads opened to the operator per completed story, `per_story`, and the repeats, `asked_again`, the threads that asked what a digest entry already answered. Say whether both fell over the window. Read `digest.curator` before you read `asked_again` ([ADR-0204](../adrs/0204-flai-stats-reports-the-curator-s-runs-in-the-digest-section-and-leaves-asked.md)). A week without `asked_again`, or a curator with `runs` 0, is unmeasured, not zero repeats: the curator has not reached it. Never report it as a finding that the digest works. Say which weeks are unmeasured, and that the curator has not run when `runs` is 0.
- Write a report under `design/analysis/`. Give each finding its evidence, its severity, and its estimated impact.
- File each actionable finding as an issue, with its impact and `--report`: `flai issue new`, or `flai issue bump` the issue that already records it. Link each issue from its finding in your report.
- Never author stories. The issue step makes draft stories from your issues, when the operator chooses it.
- Edit nothing but your report and its entry in the index. flai writes your issues; never edit them by hand.
- Hand wide search to the explorer.

## As the curator <!-- topics: all, digest -->

- List what the archive and the threads gained since your last run with `flai digest harvest --since`, and read the digest with `flai digest list` and `show` before you write.
- Keep the digest through `flai digest new`, `edit`, `absorb`, and `compact` alone: one fact per entry, citing its source. When an entry holds the fact already, merge the candidate into it with `edit` rather than write a second, and add a thread that asked what an entry answered to its `asked_again`. In every run, look for absorption: search the design with `doc_search` on each active entry's title (`flai digest list --status active`) and read what it finds with `doc_get`. Absorb only an entry a design document, convention, or ADR states, with `flai digest absorb <id> --by <path>` and the `--heading` of the section that states it. Then run `flai digest compact` once, which compacts the active entries no work item, narrative, or thread cited within `digest.compact_after` ([ADR-0199](../adrs/0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md)).
- Never write a secret's value in an entry. Say only where the operator said it lives.
- Edit nothing else, and author no stories or work items: `flai guard` refuses you every file edit and every other write ([ADR-0196](../adrs/0196-the-curator-runs-behind-the-curate-host-action-and-takes-in-what-the-archive.md)).
- End with a one-line summary naming the entries you wrote, merged, absorbed, and compacted, or saying you found nothing to take in.

## Logging your activity

- Each kind has one activity document per project: `wip/agents/planner.md`, `wip/agents/orchestrator.md`, `wip/agents/analyzer.md`, and `wip/agents/curator.md`.
- `flai serve` writes an entry when an activity ends: the timestamp, a one-line summary, the items touched, the duration, and the cost.
- A run that spans activities, as the orchestrator's does, reports each one when it ends through the MCP tool `activity_log`, with its kind, its summary, and the items it touched; flai measures the duration and the cost from the run's log. When a run ends, the time since its last reported activity is logged with the run's final reply as its summary.
- Give a summary a reader can act on: what changed, on which items, and why.
- Never write the activity document by hand.

## Asking the operator

- Ask with `thread_open` on the item the question concerns, your recommended answer first: the planner on its epic or story, the orchestrator on the item a decision concerns, the analyzer on its report or issue, and the curator on the item an entry's source comes from, and only when the entry needs the operator's judgement.
- Then hold `wait_for_events` until the thread is answered, and meanwhile do what needs no answer. The orchestrator does not hold: it ends once its reasons are done, and the answer, a thread that awaits it, is a reason that starts it again. On a planner's thread, a confirmation by `flai serve` of the orchestrator's recommendation, after `orchestration.permissions.confirm_after`, is the answer, as the operator's confirmation is.
- Never ask in the conversation.
- Never guess past a missing input the operator owns. The cost of delay inputs are the operator's, or, with `plan_backlog_stories`, the orchestrator's to give a backlog story that has none.

## When in doubt

- If a change needs code, it is a story for a story's agent.
- If a permission is unclear, ask; do not act.
- If an input is missing, ask for it rather than assume it.

<!-- system-flow:end-of-baseline -->

## Project additions
