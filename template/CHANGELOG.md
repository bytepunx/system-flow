# Changelog

## 1.1.7 - 2026-10-29

- S-0537 The conventions mark the sections a pack primed for the orchestrator's reasons keeps (patch, ADR-0217): in `strategic-agents.md`, the section "As the orchestrator" gains two sections, "Drafts and the ready column" (`all, draft, ready`) and "Accepting a story" (`all, review`), by moving its rules, whose words do not change, and its list of permissions is written on one line; the rules for every action, planning, threads, and sub-projects stay under the section's own heading. "As the planner" carries the heading topics `all, planning, planned`, "As the analyzer" `all, analysis`, and "As the curator" `all, digest`, and in `work-management.md` "Rules" carries `all, review, draft, planned`. `all` keeps each section in every pack that has it today, so no other pack changes; a pack primed with reasons keeps a section tagged with other words only when one of them is a reason of the run or `orchestration`. It needs a flai that has it, the release that carries S-0537: an older flai prints these sections in every pack, a reason pack among none, and its `flai check` warns `doc.topic` on each of the reason and role words in these headings that no sub-project, story, or epic of the project uses.
- S-0537 The orchestrator primes for the reasons its start prompt names (patch, ADR-0217): step 1 of `design/workflows/roles/orchestrator.md` calls `prime` with role `orchestrate` and the forms the prompt gives its reasons, such as `thread:TH-nnnn` or `review:S-nnnn`, as `reasons`, and with none for the general pack when the prompt says so. In `strategic-agents.md`, the first rule of "As the orchestrator" says the same. Such a pack holds the role's conventions with the sections the reasons do not need left out, what each reason loads, and the digest's answers ranked against the reasons' threads, and no briefs or catalog. It needs a flai that has it, the release that carries S-0537: an older `flai serve` gives the reasons no forms, so the orchestrator primes with the general pack, and an older `prime` does not take `reasons`.

## 1.1.6 - 2026-10-29

- S-0536 The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted (patch).
- S-0536 The orchestrator runs for a reason and ends (patch, ADR-0209): `design/workflows/roles/orchestrator.md` names under `trigger` the events it acts on, `asked`, `{event: thread, when: awaits-orchestrator}`, `{event: moved, when: to-review}`, `{event: moved, when: to-backlog}`, `{event: edited, when: draft-changed}`, `{event: planned}`, `{event: release-due}`, `{event: ready-empty}`, `{event: reordered}`, and `{schedule: orchestration.schedule}`, in place of `continuous`, which `flai check` now reports on any role as `roles.file`. Its prompt says to act on each reason the start prompt lists, log each action, then end, and that `wait_for_events` answers `end: true`, `why: reasons done`, once they are done. In `strategic-agents.md`, the section "As the orchestrator" says the same and no longer says to wait on `wait_for_events` between decisions, and "Asking the operator" says the orchestrator ends rather than hold for an answer. `design/workflows/roles/README.md` drops `continuous` from the trigger vocabulary and lists `orchestration.schedule` among the schedules a trigger may name. It needs a flai that has it, the release that carries S-0536: an older flai's `flai check` reports `orchestration.schedule` in a trigger as `roles.file`, and its `flai serve` starts the orchestrator again whenever a run ends.

## 1.1.5 - 2026-10-29

- S-0417 flai serve starts a role a project defines on the operator's word, behind the host action its file names, with an activity document under wip/agents and the same journal entries as the built-in roles (patch).
- S-0417 A role a project defines in a role file under `design/workflows/roles/` runs on the operator's word (patch, ADR-0207): `flai role run <name> [<item>] [--input name=value]...`, the MCP tool `role_run`, and the host API's `role.run` start it behind the host action its file names, which `flai serve enable <action>` turns on, on an item its file takes or on none, in the main checkout, as `<name>` or `<name>-<item>`, with `FLAI_ROLE` set to its file's `role`. Its agent is `roles.<name>.agent` in `system-flow.yaml` over `agent`, which `flai manifest set roles.<name>.agent=<json>` writes; on `claude-code` it runs as a definition flai generates from its file, and nothing is written under `.claude/`. flai serve records the run in `serve/agents.json` under `roles`, its output in `serve/agents/<key>-<name>-<time>.log`, and what it did and cost in its own activity document, `wip/agents/<name>.md`, which `wip/agents/index.md` lists and `flai check` validates; `activity_log` and the new `flai activity log <kind> "<summary>"` take its name as the kind. No template file changes. It needs a flai that has it, the release that carries S-0417: an older flai has no `flai role run` or `flai activity log`, ignores `roles.<name>.agent`, and neither lists nor checks a defined role's activity document.

- S-0416 flai guard reads a role's permitted calls from its role file and refuses the rest, and the planner's, the orchestrator's, and the analyzer's rules move into their files with their guard tests unchanged (patch).

## 1.1.3 - 2026-10-29

- S-0535 flai emits the events the orchestrator's trigger names: planned, ready-empty, release-due, a thread coming to await an agent, and the conditions awaits-orchestrator, draft-changed, and to-review (patch).
- S-0535 The roles README lists the trigger events and conditions (patch, ADR-0211): `design/workflows/roles/README.md` gains a section, "Triggers", with a table of every event a role file's `trigger` may name, `moved`, `thread`, `planned`, `ready-empty`, and `release-due` among them, and a table of every trigger condition, with the event it qualifies and when it holds, `to-review`, `to-backlog`, `draft-changed`, and `awaits-orchestrator` among them. No role file changes. It needs a flai that has it, the release that carries S-0535: an older flai's `flai check` reports a role file that names a new event or condition as `roles.file`.

## 1.1.3 - 2026-10-28

- S-0416 `flai guard` holds every strategic role to its role file (patch, ADR-0207, ADR-0210): the `Edit|MultiEdit|Write|NotebookEdit` entry of `.claude/settings.json` runs `flai guard` for any `FLAI_ROLE`, not only `plan`, `orchestrate`, `analyze`, and `curate`, behind `[ -n "$FLAI_ROLE" ] || [ -n "$FLAI_STORY" ] || exit 0;`, so the file edits of a role a project defines are held too. `flai guard` reads each role's rules from its file under `design/workflows/roles/`: the reads every agent has pass; each MCP tool, flai command, and folder to edit the file lists passes under its condition, and a call under a `permissions` entry only while the operator has its setting on, `orchestration.permissions.<p>` for the orchestrator and `roles.<name>.permissions.<p>` in `system-flow.yaml` for any other role; everything else is refused, naming the role, the rule, and the file. A built-in role with no file is held to flai's own copy; a defined role with no file, or a role whose file does not parse, is refused every write. A defined role writes as itself alone, as the planner does. The planner may no longer set or clear a story's cost of delay inputs. It needs a flai that has it, the release that carries S-0416: an older `flai guard` holds only the four built-in roles, from its own code, and passes a defined role's calls.

## 1.1.2 - 2026-10-28

- S-0502 A thread with a pending orchestrator recommendation waits on the operator however long it takes, with no setting to confirm it after a threshold (patch).

## 1.1.1 - 2026-10-28

- S-0502 The orchestrator marks a recommendation on money or a change to the host for the operator alone, and `flai serve` confirms its other recommendations after `orchestration.permissions.confirm_after` (patch, ADR-0206): in `strategic-agents.md`, the section "As the orchestrator" says that wherever it posts a recommendation, under `answer_threads` or `plan_backlog_stories`, it gives `thread_reply` `operator_only: true` with one about money (spend, a paid service, a price) or a change to the host (installing or upgrading software, configuration outside the repository, credentials), that while `confirm_after` is set `flai serve` confirms its other recommendations on a thread a story's agent or a planner opened once they have waited that long with no answer from the operator, and that it recommends there only an answer it would stand behind as the answer. The section "Asking the operator" says that on a planner's thread such a confirmation is the answer, as the operator's confirmation is. The orchestrator's definition, `.claude/agents/orchestrator.md`, says the same. It needs a flai that has it, the release that carries S-0502: an older flai does not know `orchestration.permissions.confirm_after`, so its `flai check` reports it as a permission it does not know and its `flai serve` confirms nothing, and its `thread_reply` has no `operator_only`.

## 1.1.0 - 2026-10-28

- S-0414 A strategic role is a markdown file under design/workflows/roles that flai reads and flai check validates, and the template ships the planner, the orchestrator, and the analyzer as baseline role files (minor).

## 1.0.122 - 2026-10-28

- S-0414 A strategic role is a file under `design/workflows/roles/` (patch, ADR-0207): the template ships a new design folder, `design/workflows/`, for what flai runs as data, with `roles/` holding `planner.md`, `orchestrator.md`, `analyzer.md`, and `curator.md` and a `README.md` listing them. A role file's front matter names the role, its host action, its triggers, the items a run takes, its topics, the calls `flai guard` allows it, plain or under a condition flai implements, the calls only the operator's settings allow under `permissions`, and the fields flai generates its harness definition from; its body is the definition's prompt, with the baseline marker, below which a project adds its own and `flai upgrade` keeps it. `design/README.md` and `CLAUDE.md` list the folder. In `strategic-agents.md`, the first rule says each role is defined by its role file, from which flai generates its harness definition, and that a role, a project's own included, runs as the operator, on the host with the operator's credentials, so the safety note on starting an agent applies to it; the rule on code says never to edit a role file, since a change to a role is a story the operator accepts. `flai check` reports `roles.file`, `roles.grant`, and `roles.definition`, each an error, on a file that does not parse or names what flai does not know, grants a call outside the permission that gates it, or has a harness definition its rules refuse. Nothing reads the files to decide yet: `.claude/agents/` keeps the definitions, and `flai serve` starts only the four built-in roles. It needs a flai that has it, the release that carries S-0414: an older flai has no `flai role`, and its `flai check` does not read the role files.

## 1.0.121 - 2026-10-28

- S-0504 The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked (patch).

## 1.0.120 - 2026-10-28

- S-0504 The analyzer reads the curator's runs before the digest's repeats (patch, ADR-0204): in `strategic-agents.md`, the section "As the analyzer" says to read `digest.curator` of `flai stats --json` before `asked_again`. A week without `asked_again`, or a curator with `runs` 0, is unmeasured, not zero repeats, and is never reported as a finding that the digest works. It needs a flai that has it, the release that carries S-0504: an older flai sends `asked_again` 0 in every week and no `curator`, so the analyzer cannot tell an unmeasured week from one without repeats.

## 1.0.119 - 2026-10-28

- S-0498 A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress (patch).

## 1.0.118 - 2026-10-26

- S-0498 The manifest's `threads` block says when `flai serve` reminds the operator of a thread that awaits them (patch, ADR-0205): the template's `system-flow.yaml` gains a commented `threads:` block with `remind_after`, how long a thread may await the operator before they are reminded of it, unset by default so that nothing is reminded, and `remind_every`, how often they are reminded again while it still awaits them, `24h` when unset, both Go durations that `flai manifest set` writes. It needs a flai that has it, the release that carries S-0498: an older flai does not know `threads.remind_after` or `threads.remind_every`: its `flai manifest set` refuses both keys, and written by hand they are ignored, by its `flai check` too, and its `flai serve` reminds no one.
- S-0498 `flai serve` blocks a story whose agent ended its session asking the operator, until the answer starts it again (patch, ADR-0205): in `work-management.md`, the rule on asking the designer says that while the agent waits on its question, `flai serve` blocks the story with the reason `awaiting the operator on TH-nnnn` and closes that interval before it starts the agent again on the answer, so the agent need not run `flai block` for the question, and that a block it wrote itself with `flai block` stays open until it runs `flai unblock`. It needs a flai that has it, the release that carries S-0498: an older `flai serve` blocks nothing, and an agent that ends on a question there leaves its story unmarked.

## 1.0.117 - 2026-10-25

- S-0434 flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it (patch).

## 1.0.116 - 2026-10-24

- S-0434 The analyzer reads the digest's measure on a bottlenecks run (patch, ADR-0200): in `strategic-agents.md`, the section "As the analyzer" says that on a bottlenecks run it reads `digest.weeks[]` of `flai stats --json` as evidence of the questions that wait on the operator, the threads opened to the operator per completed story, `per_story`, and the repeats, `asked_again`, and says whether both fell over the window. It needs a flai that has it, the release that carries S-0434: an older flai's `flai stats --json` has no `digest`, so the analyzer finds no such figures to read.

## 1.0.115 - 2026-10-24

- S-0432 The curator runs after each acceptance and on a schedule, marks an entry the design now records absorbed, and compacts an entry no one has cited within the window (patch).

## 1.0.114 - 2026-10-24

- S-0432 The curator runs after acceptances and on `digest.schedule`, marks absorbed the entries a document now records, and compacts the entries no one cited within `digest.compact_after` (patch, ADR-0199): the curator's definition, `.claude/agents/curator.md`, says `flai serve` starts it after a story's acceptance, when `digest.schedule` comes round, or when the operator asks, and gains two steps: in every run it searches the design with `doc_search` on each active entry's title and marks absorbed with `flai digest absorb <id> --by <path>`, and `--heading`, only an entry a design document, convention, or ADR states, then runs `flai digest compact` once; its one-line summary names the entries it absorbed and compacted too. In `strategic-agents.md`, the section "As the curator" says the curator keeps the digest through `flai digest new`, `edit`, `absorb`, and `compact`, looks for absorption in every run, runs `flai digest compact` once after its absorptions, and ends naming the entries it wrote, merged, absorbed, and compacted. The template's `system-flow.yaml` gains a commented `digest:` block with `schedule` and `compact_after`, which `flai manifest set` writes. It needs a flai that has it, the release that carries S-0432: an older flai has no `flai digest absorb`, `cited`, or `compact`, and its `flai serve` starts the curator only when the operator asks. It does not know `digest.schedule` or `digest.compact_after`: its `flai manifest set` refuses both keys, and written by hand they are ignored, by its `flai check` too. Its guard refuses the curator `flai digest absorb`, `cited`, and `compact`.

## 1.0.113 - 2026-10-24

- S-0431 The curator, a strategic agent behind the curate host action, takes in what the archive and the threads gained since its last run as digest entries written through flai alone (patch).

## 1.0.112 - 2026-10-24

- S-0431 The curator, a fourth strategic agent, keeps the digest behind the `curate` host action (patch, ADR-0196): the template ships `.claude/agents/curator.md`, the curator's definition, which reads what `flai digest harvest` lists since its last run, writes an entry with `flai digest new` for each fact the digest lacks, merges one into the entry that holds it with `flai digest edit`, records `asked_again`, never writes a secret's value, edits no file, authors no work item, and ends with a one-line summary naming the entries it wrote and merged. `.claude/settings.json` runs `flai guard` on `Edit|MultiEdit|Write|NotebookEdit` when `FLAI_ROLE` is `curate`, as for `plan`, `orchestrate`, and `analyze`. In `strategic-agents.md`, `roles` gains `curate`, the opening, the rules, and "Asking the operator" name the curator, a new section "As the curator" says the same as its definition, and "Logging your activity" names `wip/agents/curator.md`. The conventions README's row for `strategic-agents.md` and `wip/agents/README.md` name the curator and its activity document. It needs a flai that has it, the release that carries S-0431: an older flai has no `flai curate`, no `curate` host action, no MCP tool `curate`, and no `--role curate` for `flai prime`, and its `flai check` warns `conventions.roles` on `curate`. Its guard does not know the role: it passes a curator's calls unchecked, file edits and work item writes included, save the shared-paths rules and `flai forge push`, which hold every session `flai serve` starts, and `release_publish`, which only the orchestrator may call.

## 1.0.111 - 2026-10-23

- S-0430 flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's (patch).
- S-0491 Add sub-project management (patch).
- S-0491 Sub-projects are added, changed, and removed through `flai subproject`, and the orchestrator may do it with `edit_subprojects` (patch, ADR-0195): in `strategic-agents.md`, the section "As the orchestrator" says that with `orchestration.permissions.edit_subprojects` it adds, changes, or removes a sub-project only when a story or a thread asks, reads them with `subprojects`, runs `flai subproject add`, `set`, or `remove` with `--autocommit` and a `Co-Authored-By:` trailer, never `subproject_edit`, which commits nothing, and logs the change with `activity_log`. The orchestrator's agent definition, `orchestrator.md`, gains the tool `subprojects`, the permission in its step 3, and the same rule in a new step 13. The template's `system-flow.yaml` gains a comment above `projects` saying what an entry holds and that `flai subproject add` writes one, and `CLAUDE.md`'s row for `system-flow.yaml` says `flai subproject` adds, changes, and removes sub-projects. It needs a flai that has it, the release that carries S-0491: an older flai has no `flai subproject`, no MCP tools `subprojects` and `subproject_edit`, and no `manifest.projects` check. It does not know `edit_subprojects`: its `flai manifest set` refuses the key, and a permission written by hand stays off, its `flai check` reports it as a `manifest.orchestration` error, and its `flai manifest set` refuses every other change while it is there.

## 1.0.110 - 2026-10-23

- S-0430 flai prime gives a story's agent the digest entries its topics or claim select, and the orchestrator the answers the operator gave to the questions its threads ask again (patch, ADR-0194): in `work-management.md`, the rule on asking the designer says to read first the digest entries the pack loads under its `digest` heading, what past stories decided and the operator answered; to ask only what no entry answers; and, when relying on an entry, to cite its `D-nnnn` in the narrative's `## Decisions` or on the thread and check it against the design before changing what it describes, since citing an entry is what keeps it from being compacted. In `strategic-agents.md`, the section "As the orchestrator" says under `answer_threads` that its pack carries the digest's `answer` and `recommendation` entries ranked against each thread awaiting the operator, with the reason `thread TH-nnnn; rank n`, and that when one answers a thread's question it is the answer the operator gave last time: it gives the entry's path as the `source` and names its `D-nnnn` in the reply. It needs a flai that has it, the release that carries S-0430: an older flai's prime loads no digest entry.

## 1.0.109 - 2026-10-23

- S-0428 The digest, what past stories learned, is one markdown file per fact in `design/digest/`, written through `flai digest` (patch, ADR-0191): the template ships `design/digest/README.md`, which says what an entry is, that it is named `D-nnnn-<slug>.md`, what its front matter holds, that the folder has no index since `flai digest list` is one, that an entry says where a secret lives and never its value, that only flai and the digest's keeper write there, and that nothing is deleted. `design/README.md` and `CLAUDE.md` gain the folder's row beside `analysis/`. Writing, listing, and checking entries need a flai that has them, the release that carries S-0428: `flai digest new`, `edit`, `list`, and `show`, and `flai check`'s `digest.entry` and `digest.secret`; an older flai has no `flai digest`, and its `flai check` warns `doc.updated` on every entry, since it checks the folder as ordinary design documents.

## 1.0.108 - 2026-10-22

- S-0476 A component's example may carry `answers`, so the component harness answers a framed component's requests to the app's own server from them (patch, ADR-0187): `design/components/README.md` says each example may carry the `answers` the harness gives the component's requests to the app's own server, as fixtures not checked against the app's routes, and that a story that changes what a route answers edits the examples' answers that copy it. Checking answers and answering from them need a flai that has them, the release that carries S-0476; an older flai's `flai check` reports `answers` as a key it does not name.
- S-0467 The orchestrator accepts on a dirty main, leaves a story sent back to its agent, finishes an acceptance of its own that stopped, and resolves a thread of its own once its concern is gone (patch, ADR-0189): in `strategic-agents.md`, the section "As the orchestrator" says under `accept_reviews` that files the main checkout has not committed stop nothing, and the blocker `dirty_main` holds a story only while its change would overwrite one of them; that a thread it opened that nobody else wrote on is the blocker `own_thread_open`, which it resolves with `thread_resolve` and a reason once its concern is gone, then accepts; that it opens a thread on a story in review only for a reason not to accept it, and resolves that thread once the reason is gone; that when `flai accept` says it sent the story back over a rebase it cannot settle, it logs that and does nothing more, since flai moved the story back to `in-progress` and `flai serve` starts its agent on flai's thread; and that it runs `flai accept <S-nnnn> --by orchestrator` again to finish an acceptance of its own that stopped after its move to `done`, while one anyone else made stays the operator's. The orchestrator's agent definition, `orchestrator.md`, says the same in its step 11. It needs a flai that has it, the release that carries S-0467: an older flai refuses the orchestrator's acceptance on any uncommitted file outside `wip/`, lists no `dirty_main` or `own_thread_open`, leaves a stopped rebase in the story's worktree, and refuses the orchestrator the finishing of any stopped acceptance.
- S-0468 A task names files in its touches, the planner keeps a folder on a story only with its Planning note, and `after` names only a change a story needs (patch, ADR-0188): in `strategic-agents.md`, the section "As the planner" says flai refuses a task a folder touch and refuses the planner's folder touch on a story unless the story's `### Planning` heading names the folder, and that `after` on a story names only a story whose accepted change it needs on main, never an order of preference, which the pull order sets. In `work-management.md`, the touches rule says flai refuses a task's folder touch, and the close-out leaves out a `touches.folder` or a `wip.after-hub` on another story, both advisory, keeping the story's own, advisory. The planner's agent definition, `planner.md`, says the same folder and `after` rules in its step 3. It needs a flai that has it, the release that carries S-0468: an older flai neither refuses a folder touch nor reports `touches.folder` and `wip.after-hub`.

## 1.0.107 - 2026-10-21

- S-0464 The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories (patch).
- S-0464 The merge queue's batch scopes `scripts/check.sh` to the stories it verifies, so a finding on main outside every one of them no longer fails it (patch): with `MERGE_QUEUE_STORIES` set, as the merge queue's batch and each bisection step set it to the stories they verify, comma-separated, and `CLOSE_OUT_STORY` unset, `check.sh` runs `flai check --strict --story "$MERGE_QUEUE_STORIES"`, so a finding outside every one of them is a note, and records no issue, since the batch's checkout is never committed. With `CLOSE_OUT_STORY` set it runs as before, and with neither, as in CI and make, every finding counts. `scripts/README.md` says so. It needs a flai that has it, the release that carries S-0464: `flai check --story` with several stories, and the batch's `MERGE_QUEUE_STORIES`; an older flai's batch sets neither variable, and the check is unscoped as before.

## 1.0.106 - 2026-10-21

- S-0453 A thread awaits an agent only when the agent takes part in it or it is on the agent's own story, so a story's agent no longer counts the threads awaiting the operator as its own (patch).
- S-0453 The orchestrator answers the threads `inbox` marks `awaiting: "operator"` (patch, ADR-0181): in `strategic-agents.md`, the section "As the orchestrator" says that under `answer_threads` it takes from `inbox` the threads it marks `awaiting: "operator"`, which `inbox` lists in the orchestrator's role without `all: true`, and leaves out those it opened, those whose last entry is its own, those with a `pending_recommendation`, and those whose last entry is by neither a story's agent nor a planner, such as `flai serve`'s notices. The orchestrator's agent definition, `orchestrator.md`, names them the same way in its step 9. It needs a flai that has it, the release that carries S-0453: an older flai marks no thread `awaiting: "operator"`, so under it the orchestrator finds none to answer.

## 1.0.105 - 2026-10-20

- S-0403 With accept_reviews on and the queue running, the orchestrator enqueues each story in review it has verified instead of accepting it, never lands or bisects by hand, and opens no gate thread on a story the queue will verify (patch).

## 1.0.104 - 2026-10-20

- S-0403 With accept_reviews on and the queue running, the orchestrator enqueues each story in review it has verified instead of accepting it, never lands or bisects by hand, and opens no gate thread on a story the queue will verify (patch).
- S-0403 While the merge queue runs, the orchestrator enqueues each story it has verified instead of accepting it, and the queue lands it as the orchestrator with its evidence (patch, ADR-0093, ADR-0173, ADR-0176): in `strategic-agents.md`, the section "As the orchestrator" says that under `accept_reviews`, while the merge queue runs for the project, with the `merge-queue` host action on, `accept_reviews` on, and acceptance local, it verifies each story as before, writes the evidence, and enqueues the story with the MCP tool `queue_add`, giving the commit its verifier checked and the evidence, in place of `flai accept --by orchestrator`. It logs each enqueue with `activity_log`, naming the story, its place in the queue, and the policy figure. It never lands, bisects, holds, resumes, or dequeues a story by hand, and opens no thread for a blocker the queue settles: a verify not at the branch's head, a dirty main checkout, or a rebase that stops inside the shared paths. `flai serve` says in its prompt whether the queue runs. The orchestrator's agent definition, `orchestrator.md`, gains the tool `queue_add`. It needs a flai that has it, the release that carries S-0403: `flai queue add --by orchestrator --verified --evidence` and `queue_add`, the refusal of `flai accept --by orchestrator` while the queue runs, and the guard's rules for the queue.

## 1.0.103 - 2026-10-19

- S-0405 With the merge queue on, a story's close-out and flai verify leave the test tiers the manifest marks for the queue to the batch, so a close-out main outruns no longer runs integration and smoke twice (patch).

## 1.0.102 - 2026-10-19

- S-0405 While the merge queue can batch, a story's close-out and flai verify leave the test tiers the manifest marks queue to the batch, which runs them (patch).
- S-0405 While the merge queue can batch, a story's close-out and flai verify leave the test tiers the manifest marks `queue: true` to the batch, which runs them (patch, ADR-0174): in `work-management.md`, the close-out rule says that while the merge queue can batch for the project, with the `merge-queue` host action on for it, `orchestration.permissions.accept_reviews` on, and acceptance local, `flai verify` lists each such tier as `queued`, left to the merge queue, and the close-out passes without it, and that `flai verify S-nnnn --all` runs it too. `scripts/close-out.sh` says in its header which tiers it leaves to the queue. A tier in `system-flow.yaml`'s `tests` may carry `queue: true`; the template's manifest marks none. The merge queue's batch runs every tier its paths select, those marked `queue` included. It needs a flai that has it, the release that carries S-0405; an older flai ignores the key and runs every tier, and its `flai manifest set tests=<JSON>` refuses a tier that carries it.

## 1.0.101 - 2026-10-19

- S-0399 flai serve builds a batch of the enqueued stories in a worktree of its own, rebased on main in order, and verifies it once with the tiers flai verify runs, behind the merge-queue host action (patch).
- S-0399 The template's `system-flow.yaml` shows the merge queue's batch size, commented out (patch, ADR-0173): a `merge_queue` block with `batch_max`, the most enqueued stories one batch verifies together, 4 when unset. A value that is not a whole number of one or more stops every command that reads the manifest, naming the key, and `flai manifest set merge_queue.batch_max=N` writes it. The queue runs in `flai serve` with no model, behind the `merge-queue` host action, off by default, while `orchestration.permissions.accept_reviews` is on, and under local acceptance only. `flai queue add`, `remove`, `hold`, and `resume` enqueue a story in review, take it out, and hold and resume the queue. Landing a verified batch is a later story's, so stories are still accepted by hand. The key, the commands, and the host action need a flai that has them, the release that carries S-0399; an older flai ignores the block.

## 1.0.100 - 2026-10-18

- S-0446 flai check finds `components.examples` outside the story at close-out (patch).
- S-0446 A close-out records no issue for the note `components.examples` gives when a story's worktree has no sub-project install (patch, ADR-0171): in `work-management.md`, the close-out rule says that the check scoped to the story also leaves out the note `components.examples` gives when it could not check a component's examples for want of `node` or a sub-project's svelte, unless the story changes a file under the folder it names, since only installing the tool clears it and the main checkout's check still reports it. A sub-project whose svelte is not installed now gets one such note on its examples folder rather than a warning on each examples file. It needs a flai that has it, the release that carries S-0446.

## 1.0.99 - 2026-10-18

- S-0397 Under pull-request acceptance, flai accept merges the story's pull request through the remote by the project's merge method, then fetches main and moves to done, archives, and commits as today (patch).
- S-0397 Under pull-request acceptance, flai accept merges the story's pull request through the remote by the project's merge method, then fetches main and moves to done, archives, and commits as today (patch, ADR-0170): `git.md`'s rule on accepting an item now covers both modes. Under local acceptance flai merges the story branch on the host. Under pull-request acceptance flai merges the story's pull request on the forge by `acceptance.merge_method`, through `flai serve` when the token is not at hand, deletes the story branch on `origin`, and brings main up to `origin`'s. Either way nothing is tagged and nothing else is pushed: the acceptance commit reaches the remote only when published. The merge needs a flai that has it, the release that carries S-0397; an older flai merges on the host under either mode.

## 1.0.98 - 2026-10-17

- S-0395 Under pull-request acceptance, flai accept and the orchestrator's acceptance refuse a story while a required check has not passed or a required review is outstanding, naming each (patch).
- S-0395 Under pull-request acceptance, flai accept and the orchestrator's acceptance refuse a story while a required check has not passed or a required review is outstanding, naming each (patch, ADR-0169): in `strategic-agents.md`, the orchestrator's dry run step says that under pull-request acceptance the blockers include the pull request's required checks that have not passed, its standing requested changes, and its missing approvals, which only CI and the team's reviewers clear. The orchestrator gains a rule that it never approves a pull request's review. The blockers need a flai that has them, the release that carries S-0395.

## 1.0.97 - 2026-10-17

- S-0393 Under pull-request acceptance, a story's move to review has flai serve push its branch and open or update one pull request from the project's pull request template (patch).
- S-0393 Under pull-request acceptance, flai serve pushes a story's branch when it moves to review and opens or updates its one pull request from the project's pull request template (patch, ADR-0166): `git.md` names the one exception to never force pushing, a story's own branch that `flai stream sync` rebased, which `flai serve`, never an agent, pushes with `--force-with-lease` on the head it last pushed. Its rule on pull requests now says that under pull-request acceptance flai opens and updates the story's pull request from `.github/pull_request_template.md`, and that an agent never opens one itself. The template's `.github/pull_request_template.md` marks the places flai fills, each between `<!-- flai:<part> -->` and `<!-- /flai:<part> -->`: `story` for the story's ID, `criteria` for its acceptance criteria as ticked, `definition-of-done`, whose checklist flai keeps as the project wrote it, and `verify` for the story's last `flai verify` result and its commit. It gains the sections `Acceptance criteria` and `Last flai verify` for them. A project that removes a mark gets that part appended under its own heading. The push and the pull request need a flai that has them, the release that carries S-0393, with the `forge` host action on.

## 1.0.96 - 2026-10-17

- S-0392 A project chooses local or pull-request acceptance in system-flow.yaml, and flai reaches its GitHub remote through a forge contract with a token the host names (patch).
- S-0392 A project chooses local or pull-request acceptance in its manifest, and the host names the variable that holds the forge's token (patch, ADR-0165): the template's `system-flow.yaml` shows, commented out, an `acceptance` block with `mode` (`local`, the default, or `pull-request`), `forge` (`github`, the default and only one), and `merge_method` (`merge`, the default, `squash`, or `rebase`). `flai manifest set` writes each key, and a value outside a key's values stops every command that reads the manifest rather than falling back to `local`. The forge's token is the host's, read from `GITHUB_TOKEN` unless `flai serve forge github --token-env <NAME>` names another variable, and never in the manifest; `flai forge check` says whether it can read and push to the repository `origin` names. Nothing opens or merges a pull request yet: `flai accept` merges on the host under either mode until the later stories of E-0021. The key and the commands need a flai that has them, the release that carries S-0392; an older flai ignores the block.

## 1.0.95 - 2026-10-17

- S-0443 A component's examples are kept under design/components, chosen on the Components page, and kept in step with its props by flai check (patch).
- S-0443 A component's examples files under `design/components/` are the second exception to "Everything is markdown" (patch, ADR-0164): `documentation.md` says that everything is markdown with YAML front matter except two kinds of JSON file, the design token file and the component examples files, one `*.examples.json` per component under `design/components/`. A new project has `design/components/README.md`, which says the folder holds supporting files for the project's components, today one examples file per component at `<sub-project>/<path>`, the sub-project's `name` and the component's path in it with `.svelte` replaced by `.examples.json`, holding `examples`, each with a `name`, its `props`, and an optional `theme` and `width`. It says the dashboard's Components page writes them through `file.save`, a story that changes a component keeps its examples in step, and `flai check` reports `components.examples`. `design/README.md` and `CLAUDE.md` list the folder in their layout tables. Saving an example from the page and the check need a flai that has them, the release that carries S-0443; an older flai neither checks the files nor lets the page create one.

## 1.0.94 - 2026-10-15

- S-0386 A story's agent working on a project's UI is primed with its design system, and a baseline convention tells it to build from the tokens and propose changes on its branch (patch).

## 1.0.93 - 2026-10-15

- S-0381 A project's design tokens are a DTCG file the manifest names, beside its guidance, and flai validates them, renders them as CSS, reports their contrast, and serves them to the dashboard (patch).
- S-0381 The design token file is the one exception to "Everything is markdown" (patch, ADR-0146): `documentation.md` says that everything is markdown with YAML front matter except the design token file, JSON in the W3C Design Tokens Community Group (DTCG) format at the path the manifest's `design.tokens` names, `design/system/design-system.tokens.json` by default. Its guidance, `design/system/design-system.md`, is markdown like any other design document. The manifest key and `flai design init`, which writes both files, need a flai that has them, the release that carries S-0381.

## 1.0.92 - 2026-10-14

- S-0374 An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again (patch).
- S-0374 A close-out that main outran syncs and verifies again, at most three rounds (patch, ADR-0156): `work-management.md` says that when the main branch moved during the close-out's tiers, such as by an acceptance, the close-out syncs and runs every step again on its own. The template's `scripts/close-out.sh`, when its last check `flai verify "$story" --sync-only` stops at the `sync` step, runs `flai stream sync "$story"` and, when the sync ends clean, `flai verify "$story" --record-issues`, the commit step, and the last check again, at most three rounds in all. It prints each new round as `close-out: main moved during the run; flai stream sync S-nnnn (round 2 of 3)`, and its last line names the round it ended on when it ran more than one, as in `close-out: S-nnnn passed every step in round 2 of 3; ready to move to review`; a run of one round ends as before. After the third round's last check fails, it stops at the sync check as before; a sync that stops on a conflict or a refusal stops it at `flai stream sync`, and a verify that fails in a later round at `flai verify`. `flai verify` reuses across that clean rebase the tiers that the main branch's new paths do not select, so a later round runs again only the tiers main's change can reach; that needs a flai that has it, the release that carries S-0374, after 1.49.4. With an older flai the rounds still run, and each runs every tier.

## 1.0.91 - 2026-10-14

- S-0372 A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other (patch).
- S-0372 A reply that asks nothing is sent settled, and wakes no one (patch, ADR-0154): `work-management.md` says to answer a message with `flai message reply --settled`, or `message_reply` with `settled`, when the answer asks nothing, such as agreement, so the conversation awaits neither story and the other story's agent is not woken; a settled entry to your story owes no reply. A `flai task done` notice that does not break your work is answered settled, and one that does gets a plain reply saying how. The notice leaves out the paths sync regenerates, and while your story is blocked `flai serve` starts you again only for the answer to your thread, not for a message. It needs a flai that has `--settled`; an older flai refuses the flag, and its `message_reply` has no `settled`.

## 1.0.90 - 2026-10-12

- S-0239 dashboard.allow_unsigned lets a development build connect, shown on every page and in every status (patch).
- S-0239 The operators guide says how to allow unsigned builds (patch, ADR-0070): `docs/operators/index.md` gains a section on the dashboard and `dashboard.allow_unsigned`, set in the host's flai configuration, in `system-flow.yaml` for one project, or for one start with `flai dashboard --allow-unsigned`. With it, flai dials a dashboard image on no signed list and accepts a dashboard with no release stamp, and the dashboard, started with `FLAIOVER_ALLOW_UNSIGNED=1`, accepts a flai with none; a stamp that does not verify is still refused, and the credential is still proven. It is shown as `unsigned allowed: <sides>` in a banner on every dashboard page and in the three statuses, and as a `warn` per connection. It needs a flai that has it; an older flai ignores `dashboard.allow_unsigned` and refuses `--allow-unsigned`.

## 1.0.89 - 2026-10-12

- S-0367 Perform lint fix on every markdown file a story touches as part of validation (patch).
- S-0367 A test tier may carry `fix_command` (patch, ADR-0139): an argument list, with the placeholders `command` takes, that rewrites in place the files the tier selects. flai runs it before `flai task done` commits, after each clean `flai stream sync` over the files the branch changes, committing what it rewrote on the branch, and before `flai accept` commits the main checkout; its exit status is ignored. The template's `system-flow.yaml` adds none: the comment on its `markdown` tier says to add `fix_command: [scripts/lint-md.sh, --fix, "{files}"]` once `scripts/lint-md.sh` takes the paths. It needs a flai that has it; an older flai ignores `fix_command`.

## 1.0.88 - 2026-10-12

- S-0349 A story's agent names a provider from the host's providers map, and a project's manifest may override the entry's api, base_url, and models (patch).
- S-0349 A story's agent names a provider, and a project may override the provider's api, base_url, and models (patch, ADR-0129): the template's `system-flow.yaml` shows, commented out, an `agent` block with `provider` beside `harness` and `model`, and a `providers` entry with `api`, `base_url`, and `models`. A provider is an entry on the host, set with `flai serve agent provider <name>`; its `key_env`, the name of the variable that holds the key, is the host's alone, and `flai check` refuses it in the manifest. No agent is started through a provider yet: S-0350 has `flai serve` read the variable `key_env` names and set the harness's environment. It needs a flai that has it; an older flai ignores `provider` and `providers`, and drops a story's `provider` when it writes the story.

## 1.0.87 - 2026-10-12

- S-0366 flai check finds `threads.archived` outside the story at close-out (patch).
- S-0366 A close-out records no `threads.archived` (patch, ADR-0137): `work-management.md` says that a `threads.archived` outside the story, an open thread on an archived item, is left out of a close-out's check, as an `item.archive`, a `board.wip-limit`, and an `issues.no-story` are, since it names that thread, not the story, and only the operator or the thread's opener clears it, with `flai thread resolve` in the main checkout. It needs a flai that has it; an older flai still records every `threads.archived` at close-out.

## 1.0.86 - 2026-10-12

- S-0365 flai check finds `issues.no-story` outside the story at close-out (patch).
- S-0365 A close-out records no `issues.no-story` (patch, ADR-0136): `work-management.md` says that an `issues.no-story` outside the story, an open issue no open story links, is left out of a close-out's check, as an `item.archive` and a `board.wip-limit` are, since it names that issue, not the story, and only the operator's choice of a story for that issue clears it. It needs a flai that has it; an older flai still records every `issues.no-story` at close-out.

## 1.0.85 - 2026-10-08

- S-0342 A test tier the manifest marks as covered by another is skipped when both are selected, so verify does not run the short Go tests and then the full suite (patch).
- S-0342 A test tier may name the tiers it covers (patch): a tier in the manifest's `tests` may carry `covers: [<tier>...]`, the tiers it runs everything of, and when `flai verify` or `flai test --all` selects it with one it covers, the covered tier is skipped and reported `skipped: covered by <tier>`. The template's `system-flow.yaml` covers nothing, since its `test` and `integration` scripts run different things; a comment on its `integration` tier says to add `covers: [test]` when that script runs the behaviour tests too. It needs a flai that has it; an older flai ignores `covers` and runs both tiers.

## 1.0.84 - 2026-10-08

- S-0347 The close-out's last check that the branch contains main fails when flai commits wip on main during its run (patch).
- S-0347 A close-out passes over flai's `wip/` commits on the main branch during its run (patch, ADR-0135): `work-management.md` says that a story's branch at review may lack commits on the main branch that change only `wip/` paths it does not change, such as flai's replans, and that the close-out's last check is `flai verify S-nnnn --sync-only`. Both of the close-out's checks that the branch contains the main branch pass over such commits; any other commit the branch lacks, such as an acceptance, still stops it, and needs `flai stream sync` and another run. The template's `scripts/close-out.sh` runs `flai verify "$story" --sync-only` for its last check, in place of `git merge-base --is-ancestor`. It needs a flai that has `--sync-only`; an older flai refuses the flag, so the close-out stops at its last check.

## 1.0.83 - 2026-10-08

- S-0334 A story held on overlap is asked about by message, and starts when the holding story's agent shares the paths or narrows its claim (patch).
- S-0334 A story held on overlap is asked about by message, and starts when the holding story's agent shares the paths or narrows its claim (patch, ADR-0134): `work-management.md` says that when a message by `flai` asks about a ready story your story holds on overlap, you answer it before you go on with its paths: narrow your `touches` with `flai touches` if your story will not change them, share them with `flai message share <MS-nnnn> --paths <path> "<split>"`, or `message_share`, saying who changes what, if the two stories can change them apart, so the ready story can start, or reply saying why the hold stands. When `flai serve` starts your story on a share, its prompt names the conversation, the shared paths, and the split: you keep to your side of the split and say so in the conversation when it no longer fits. It needs a flai that has it; an older flai never asks about a hold, has no `flai message share`, and ignores a conversation's `shares`, so the overlap holds as before.

## 1.0.82 - 2026-10-08

- S-0348 flai check finds `board.wip-limit` outside the story at close-out (patch).
- S-0348 A close-out records no `board.wip-limit` (patch, ADR-0133): `work-management.md` says that a `board.wip-limit`, a column of the board over its limit, is left out of a close-out's check, as an `item.archive` is, since the board counts the main checkout's stories, which no story branch changes, and only the operator's acceptances, cancellations, or a raised limit clear it. It needs a flai that has it; an older flai still records every `board.wip-limit` at close-out.

## 1.0.81 - 2026-10-08

- S-0322 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit (patch).
- S-0322 Closing a task commits only the task's own paths, so the tasks of one layer close apart (patch, ADR-0128): `git.md` says that `flai task done`, or `task_done`, commits the changed paths the task's `touches` cover and those no other open task covers. It leaves a path only another open task covers for that task's close, lists it, and goes on, and the close that leaves none syncs the branch. `-m` is needed only when there is something to commit. It needs a flai that has it; an older flai commits every uncommitted file in the worktree and needs `-m` on every close.

## 1.0.80 - 2026-10-08

- S-0287 flai guard lets a task sub-agent run flai adr new but refuses flai adr topics (patch).
- S-0287 A task sub-agent records an ADR with flai and leaves the commit to the story's agent (patch, ADR-0127): `delegation.md` § As a task sub-agent says to record a decision with `flai adr new` or `adr_new`, and set its topics with `flai adr topics`, never by copying what flai writes or guessing a number, and never with `--commit`, `--autocommit`, or `adr_new`'s `commit`, which `flai guard` refuses. It needs a flai that has it; an older flai's guard still refuses a sub-agent every `flai adr` command and `adr_new`.

## 1.0.79 - 2026-10-08

- S-0323 flai check finds `narrative.state` outside the story at close-out (patch).
- S-0323 A close-out records no `narrative.state` finding on another open story's narrative (patch, ADR-0125): `work-management.md` says that a `narrative.state`, a `## Current state` or `## Next steps` not yet written, on the narrative of another story not yet done or cancelled is left out of a close-out's check, as a markdown finding there is, since only that story's agent writes them. It needs a flai that has it; an older flai still records every `narrative.state` outside the story at close-out.

## 1.0.78 - 2026-10-08

- S-0318 flai check finds `markdown.MD038` outside the story at close-out (patch).
- S-0318 A close-out records no markdown finding on another open story's narrative (patch, ADR-0123): `work-management.md` says that a markdown finding on the narrative of another story not yet done or cancelled is left out of a close-out's check, as an `item.archive` is, and recorded in no issue, since that story's own close-out checks its narrative and stops on it. A markdown finding anywhere else outside the story is still recorded. It needs a flai that has it; an older flai still records every markdown finding outside the story at close-out.

## 1.0.77 - 2026-10-07

- S-0280 flai check finds `item.archive` outside the story at close-out (patch).
- S-0280 A close-out records no `item.archive` (patch, ADR-0122): `work-management.md` says that an `item.archive`, a done or cancelled item not yet archived, is left out of a close-out's check, as a `wip.overlap` between two other stories is, since it never names your story and only `flai archive` in the main checkout clears it. It needs a flai that has it; an older flai still records every `item.archive` at close-out.

## 1.0.76 - 2026-10-07

- S-0332 flai tells two stories' agents of a trial-merge conflict or a grown overlap with a message between them, and asks the operator only when they do not agree (patch).
- S-0332 flai tells two stories' agents of a trial-merge conflict or a grown overlap with a message between them, and asks the operator only when they do not agree (patch, ADR-0121): `work-management.md` says that when a write grows the two stories' claims to overlap, flai has messaged the other story from the one whose claim grew, in the two stories' conversation, `about` the paths, and you agree there with `message_reply`, or `flai message reply`, which of you changes them first before you change them. When `flai stream sync` finds that your branch conflicts with another open story's, it tells the two stories in their conversation, not on a thread, and names the conversation: you answer it and settle who changes what, one of you narrowing its change or naming the other in `after:` and waiting for it, and the next sync that finds the two merging cleanly closes the conversation. You ask the operator only when the two of you do not agree, with `flai message escalate <MS-nnnn> "<reason>"`, or `message_escalate`, which opens a thread on your story and records it in the conversation, which stays open. It needs a flai that has it; an older flai opens a conflict thread at sync and has no `flai message escalate`.

## 1.0.75 - 2026-10-07

- S-0333 Closing a task messages every open story whose claim covers a path the task changed, before either story is accepted (patch).
- S-0333 Closing a task messages every open story whose claim covers a path the task changed, before either story is accepted (patch, ADR-0120): `work-management.md` says that a message telling you a task of another open story changed paths your story claims comes from `flai task done` as it closed that task, naming the task, the commit and its subject, and the paths in `about`. You read the change with the `git show` it names, check whether it breaks your story's work, and reply on the conversation either way, since it awaits your story until you do: when the change breaks your work, or your work depends on the shape it changed, you say how and note it in the narrative; otherwise you say it does not. You adjust to it early where you can, and you do not sync to get it: `flai stream sync` brings it only once that story is accepted, and the `overlapped` change follows then. It needs a flai that has it; an older flai sends no message when it closes a task.

## 1.0.74 - 2026-10-07

- S-0335 A story's agent that waits only on another agent's reply ends, and flai serve starts it again when the reply or a new message to its story comes (patch).
- S-0335 A story's agent that waits only on another story's agent's reply ends, and flai serve starts it again when the reply or a new message to its story comes (patch, ADR-0120): `work-management.md` says that when nothing is left but the other story's agent's reply, you write the narrative's `## Current state` and `## Next steps` with `flai stream state`, naming the conversation and what each reply leads to, and end the session, as for the designer's answer; `flai serve` starts you again, in the same session, when that agent replies or a new message to your story comes, and your first `inbox` lists it under `messages`. When `wait_for_events` answers `end: true`, its `why` names the open question or the conversation awaiting another story's agent's reply. An agent run by hand, not by `flai serve`, is not started again and holds `wait_for_events` until the reply comes. It needs a flai that has it; an older flai keeps such an agent waiting.

## 1.0.73 - 2026-10-07

- S-0331 A story's agent sends and answers messages through MCP, inbox lists those to its story, and wait_for_events wakes on one (patch).
- S-0331 A story's agent messages the agent of another open story, not a thread, before changing a path both claim (patch, ADR-0120): `work-management.md` says that when an `overlapped` change's `cause` is a story in progress, you message that story's agent with `message_send`, or `flai message send`, naming the paths in `about`, before you change them, and narrow your `touches` if you can. `inbox` lists your story's open conversations under `messages`, and `wait_for_events` wakes on a message to your story as a change of kind `message`. You answer a message awaiting your story with `message_reply`, or `flai message reply`, before you go on with the paths it is about, and ask the operator on a thread only when the two of you do not agree; a question for the designer is still a thread. `delegation.md` says a sub-agent does not send messages, which `flai guard` refuses it. It needs a flai that has the message tools.

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
