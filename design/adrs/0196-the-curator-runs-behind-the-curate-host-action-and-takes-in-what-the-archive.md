---
id: ADR-0196
title: "The curator runs behind the curate host action and takes in what the archive and the threads gained since its last run as digest entries written through flai digest alone, which flai guard holds it to"
status: accepted
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0075, ADR-0099]
topics: [cli, template, analysis, planning]
---

# ADR-0196 The curator runs behind the curate host action and takes in what the archive and the threads gained since its last run as digest entries written through flai digest alone, which flai guard holds it to

## Context

E-0025 keeps a digest of what past stories decided and what the operator answered, one fact per file under the design folder's `digest/` (ADR-0191), and a story's pack loads the entries its topics or claim select (ADR-0194). `flai digest seed` writes an entry for each operator answer and each closed issue's cause, mechanically. It cannot judge which decisions in a done story's narrative are worth keeping, nor tell that a new answer repeats an entry, nor merge two raw entries that state one fact. Only a model can. TH-0476's answer to question 3 settled that the agent that does this is a role of its own, built now on the analyzer's machinery (ADR-0099) and its pack by role (ADR-0075), rather than waiting for E-0024's role files. `flai digest harvest`, `new`, `edit`, `list`, and `show` exist (S-0428 to S-0430); nothing starts an agent that uses them.

## Decision

The curator is a session `flai serve` starts for the whole project behind the host action `curate`, on the operator's word; it takes in what the archive and the threads gained since its last run as digest entries written through `flai digest new` and `edit` alone, which `flai guard` enforces, and its activity entry names the entries it wrote.

1. **Behind a host action.** `curate` is off by default. `flai curate [--since <date>] [--json]`, the host API's `curate.run` (`{since?}`, which runs `flai curate --json`), and the MCP tool `curate` (for the operator's own agent; refused to an agent `flai serve` started) start it on demand. `flai curate` starts the run itself, as `flai analyze` does, and the serving flai settles it. There is no schedule yet.
2. **One run per project.** A second start while one runs is refused, naming the run. So are a start while the action is off, a `--since` that is neither a `YYYY-MM-DD` day nor an RFC 3339 time, and a start with nothing to run it: no harness and no command.
3. **Its agent.** `digest.agent`, merged over the project's `agent` as `analysis.agent` is (ADR-0099), gives its harness, model, config, and roles. `flai check` reports a `digest.agent` that is not valid as `manifest.digest`.
4. **How it runs.** In the project's main checkout, as `curator`, with `FLAI_ROLE=curate`, no `FLAI_ITEM` or `FLAI_STORY`, and a prompt of its own. On `claude-code` it runs as the project's `.claude/agents/curator.md`, selected with `--agent curator`, its session named `<key> curate`. The operator's command gets `FLAI_ROLE`, and `FLAI_SINCE` when the run has a since, since it has no prompt to be told it in.
5. **Its pack.** `flai prime --role curate`, with no item: its topics are `digest` and `planning`, and it holds the conventions whose `roles` are empty or list `curate`, `strategic-agents.md` among them.
6. **What it reads and writes.** It reads `flai digest harvest --since` the run's since: the `--since` it was given, else the time of the last entry in `wip/agents/curator.md`, else everything on its first run. For each candidate it writes an entry with `flai digest new` when no entry holds the fact, or merges it into the entry that does with `flai digest edit`, and records in that entry's `asked_again` each thread that asked what the entry already answered. It never writes a secret's value, only where the operator said it lives. It edits no file and authors no work item. It asks the operator on a thread only when an entry would need their judgement. It ends with a one-line summary naming the entries it wrote and merged.
7. **Held by the guard.** With `FLAI_ROLE=curate`, `flai guard` checks the session's own calls. It passes the reads a sub-agent has; `flai digest new`, `edit`, `harvest`, `list`, and `show`; and the MCP tools `inbox`, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events`. It refuses `Edit`, `Write`, `MultiEdit`, and `NotebookEdit` on every file, `flai digest seed`, every work item write, every other flai tool and command that writes, `curate`, `analyze`, and `plan` among them, and git's writes, each refusal saying that the curator writes entries only with `flai digest new` and `edit`. It holds the curator to its own name, `curator`, as ADR-0143 holds the analyzer to its. The `Edit|MultiEdit|Write|NotebookEdit` entry in `.claude/settings.json` runs the guard when `FLAI_ROLE` is `curate` as well as `plan`, `orchestrate`, or `analyze`.
8. **Recorded as the other runs are.** The run is kept in `serve/agents.json` under `curator`, the newest alone, marked with its role so that it is never taken for the orchestrator's run, with its since and its trigger (`asked`). Its log is `<key>-curator-<start>.log`. It does not count against the in-progress limit and holds no story back. It ends `worked` on a clean exit, whether it wrote or merged an entry or found nothing to take in, and `failed` on a failed exit, with why. flai serve then logs its activity in `wip/agents/curator.md`: the run's last line as the summary, its items the entries whose files the run wrote or changed, its trigger, and the run's seconds and cost. It charges no item: its cost stays in its kind's project strategic total, which `flai stats` reports as it reports the analyzer's.

## Consequences

- The operator asks for a curation from a shell or their own agent and pays for a curator only then. Starting it after acceptances or on a schedule, absorbing entries the design took in, and compacting old ones are S-0432's.
- Every entry the curator writes passes through `flai digest`, which numbers it, names its file, and checks its front matter, so a curator run cannot leave a malformed or misnumbered entry, and nothing in the digest is deleted (ADR-0191).
- The curator's entries land in the main checkout, uncommitted, as the planner's items do; the operator's next commit, or an acceptance, takes them.
- A project whose `.claude/agents/` has no `curator.md` cannot start the curator on `claude-code`; the refusal names `flai upgrade`. A project that keeps its old settings has a curator whose file edits the guard does not see; its definition lists no `Edit` or `Write`, so it has none to make.
- As for the analyzer, other shell commands pass, as a sub-agent's do (ADR-0060), so a write through the shell is not seen.
- When E-0024's role files land, S-0419 moves the curator onto the role path as a fourth baseline role (TH-0473).

## Alternatives considered

- **Waiting for E-0024's role path.** E-0024 is back in the backlog; waiting would hold S-0431 to S-0435 (TH-0476).
- **The analyzer curating the digest.** Its guard lets it edit one folder and its run ends in one report; a curator writes many entries through flai and none by hand, so one role would need two fences and two settlements.
- **`Edit` and `Write` held to the digest folder, as the analyzer's are to its own.** flai already numbers, names, and checks entries; a hand-written entry could skip the next number or break the front matter that `flai check` refuses.
- **`flai digest seed` for the curator.** Seed writes an entry per raw answer and cause, which is what the curator exists to judge and merge; it stays the operator's.
- **An exit that wrote no entry counted as `failed`.** A run that found every candidate already held did its work; the exit code says whether it failed.
