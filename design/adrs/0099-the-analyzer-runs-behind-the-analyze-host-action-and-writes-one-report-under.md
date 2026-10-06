---
id: ADR-0099
title: "The analyzer runs behind the analyze host action and writes one report under design/analysis, which flai guard holds it to"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0082]
topics: [cli, template, analysis]
---

# ADR-0099 The analyzer runs behind the analyze host action and writes one report under design/analysis, which flai guard holds it to

## Context

The analyzer reads the metrics, the design, the code, and the issues, and says what it finds: bottlenecks in the flow, gaps between `design/system` and the code, and technical and security risks (E-0016). S-0207 gave it a pack (`flai prime --role analyze`, ADR-0075) and S-0206 an activity document (ADR-0079), but nothing started it. The planner (ADR-0082) and the orchestrator (ADR-0087) set the shape a strategic agent takes: a host action, a role in `FLAI_ROLE`, a definition in `.claude/agents/`, and guard rules of its own. The analyzer fits neither of them as it stands. Unlike the planner it has no item; unlike the orchestrator it has an end, its report. And unlike both it must write a file, which the guard refused the planner and the orchestrator outright. The convention says it never authors stories: what it finds is for the operator to read and act on. Its findings must also stay apart from the living design, which says how the system is, not what an agent thought was wrong with it.

## Decision

The analyzer is a session `flai serve` starts for the whole project behind the host action `analyze`, on the operator's word or on `analysis.schedule`; it writes one report under `design/analysis/` and edits nothing else, which `flai guard` enforces, and its activity entry names the report.

1. **Behind a host action.** `analyze` is off by default. `flai analyze [--focus bottlenecks|intent|risk]`, the host API's `analyze.run` (`{focus?}`, which runs `flai analyze --json`), and the MCP tool `analyze` (for the operator's own agent; refused to an agent `flai serve` started) start it on demand. `flai analyze` starts the run itself, as `flai plan` does, and the serving flai settles it; with no `flai serve` running it says so, and the analyzer still runs.
2. **On a schedule.** `analysis.schedule` in `system-flow.yaml`, a five-field cron expression in UTC or `daily`, parsed as `planning.schedule` is, starts it with no focus each time it comes round while the action is on. The run's trigger is `schedule <spec>`. Times missed between two looks come round once, a schedule set or changed comes round first at its next time from then, and times that passed while `flai serve` was down or the action was off do not come round. A start refused when its time comes, as while an analyzer runs, is logged and not tried again until the schedule next comes round.
3. **A focus, or all of them.** The focus is `bottlenecks`, `intent`, or `risk`. A run asked for none looks for all three, and its focus is recorded as `all`.
4. **One run per project.** A second start while one runs is refused, naming the run. So are a start while the action is off, a focus none of the three, and a start with nothing to run it: no harness and no command.
5. **Its agent.** `analysis.agent`, merged over the project's `agent` as `planning.agent` is (ADR-0082), gives its harness, model, config, and roles. `flai check` reports an `analysis.agent` or `analysis.schedule` that is not valid as `manifest.analysis`.
6. **How it runs.** In the project's main checkout, as `analyzer`, with `FLAI_ROLE=analyze`, no `FLAI_ITEM` or `FLAI_STORY`, and a prompt of its own. On `claude-code` it runs as the project's `.claude/agents/analyzer.md`, selected with `--agent analyzer`, its session named `<key> analyze` and the focus. The operator's command gets `FLAI_ROLE`, and `FLAI_FOCUS` when there is a focus, since it has no prompt to be told it in.
7. **One report.** It writes `design/analysis/<date>-<focus>.md`, the date its UTC day, with front matter `title`, `updated`, `status`, `focus`, and the window its metrics cover, `from` and `to`, and one section per finding with its evidence, severity, and estimated impact. It adds a row for it to `design/analysis/README.md`. The statuses are the documentation standard's: `draft` while the analyzer writes it, `active` once the run has ended, and `deprecated` when a later report replaces it. `flai check` validates each report's name and front matter (`analysis.report`), and warns of a report the README does not list and of a row with no report (`analysis.index`).
8. **Held by the guard.** With `FLAI_ROLE=analyze`, `flai guard` checks the session's own calls. It passes the reads a sub-agent has, `flai stats` among them, and the MCP tools `inbox`, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events`. It passes `Edit`, `Write`, and `NotebookEdit` only on a file inside the design folder's `analysis/`, the path resolved from the project's root with `..` and symbolic links followed. It refuses `item_new`, `item_edit`, and `item_move`, since the analyzer authors no stories, and every other flai tool and command that writes, `issue new` and `bump` among them, and git's writes. The `Edit|Write|NotebookEdit` entry in `.claude/settings.json` runs the guard when `FLAI_ROLE` is `analyze` as well as `plan` or `orchestrate`.
9. **Recorded as the other runs are.** The run is kept in `serve/agents.json` under `analyzer`, the newest alone, with its focus, its trigger (`asked`, or the schedule's), and, once it has ended, its report: the newest Markdown file in the folder, its README aside, changed since the run started. Its log is `<key>-analyzer-<start>.log`. It does not count against the in-progress limit and holds no story back. It ends `failed` on a failed exit and on an exit that left no report, and `worked` otherwise. flai serve then logs its activity in `wip/agents/analyzer.md`: the run's last line as the summary, with the report added when the line does not name it, or `no report written under design/analysis`, its trigger, and the run's seconds and cost. It charges no item.

## Consequences

- The operator asks for an analysis from a shell or their own agent, or sets a schedule, and pays for an analyzer only then. A schedule spends without asking each time, so it is unset until the operator sets it.
- Every report sits in one folder, named for its day and focus, listed in one index, and checked by `flai check`, so the dashboard's documents page shows them as it shows any design document.
- The analyzer is the first strategic agent the guard lets edit a file, and only inside its folder. A project that keeps its old settings has an analyzer whose file edits the guard does not see; its definition still lists `Edit` and `Write`, so nothing fences it then. As for the planner, other shell commands pass, as a sub-agent's do (ADR-0060), so a write through the shell is not seen.
- A project whose `.claude/agents/` has no `analyzer.md` cannot start the analyzer on `claude-code`; the refusal names `flai upgrade`.
- The convention's "As the analyzer" says it files its actionable findings as issues. Filing issues is S-0224's: until then the guard refuses it `flai issue new` and `bump`, and its findings stay in the report. Its activities charge no item until S-0227 charges them to the issues it files.
- flai checks that a report's status is one of the three; it does not change it. The analyzer writes it.
- No dashboard page has an Analyze button yet; `analyze.run` is there for one (S-0228).
- One run at a time keeps the report the end of a run found unambiguous: the newest file changed since it started is that run's.

## Alternatives considered

- **The report under `design/system` or `wip/`.** `design/system` says how the system is, and a finding is a claim about it until someone acts on it. `wip/` is work state, and a report is not work. A folder of its own, as `design/experiments` is for results, keeps both clean.
- **The analyzer authoring stories.** Stories are the operator's and the planner's; the convention keeps the analyzer to findings, and S-0224 routes its actionable ones through issues, which the operator chooses to turn into stories.
- **No `Edit` or `Write`, and the report written through flai**, as the planner writes items. flai has no command that writes a new report, and the guard would still have had to hold a path. Holding Claude Code's file tools to one folder is the narrower change.
- **A long-running analyzer, like the orchestrator.** Its work has an end, a report; a session kept waiting would spend for nothing between reports.
- **One run per focus at once.** Three sessions would triple the spend, and a run's report could no longer be told by being the newest file changed since it started.
- **A host cron job running `flai analyze`.** It would live outside the project, unseen by `flai check` and the operator's manifest, and would not see the host action. `flai serve` already reads `planning.schedule` the same way.
- **An exit with no report counted as `worked`.** The run would claim success with nothing to read; `failed` with the reason says what went wrong.
