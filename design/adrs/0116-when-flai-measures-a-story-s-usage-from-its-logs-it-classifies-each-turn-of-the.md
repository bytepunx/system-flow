---
id: ADR-0116
title: "When flai measures a story's usage from its logs it classifies each turn of the story's own agent, records the counts per day in usage.turns, and flai stats reports them"
status: proposed
date: 2026-10-07
supersedes: []
superseded_by: []
---

# ADR-0116 When flai measures a story's usage from its logs it classifies each turn of the story's own agent, records the counts per day in usage.turns, and flai stats reports them

## Context

E-0017 moves the story loop's mechanical work from the agent into flai: one command closes a task, verifies a story, runs the tests a change selects, ticks a criterion, writes a narrative's state, and starts a story. It is done when a story agent's run shows no turn that only commits, syncs, moves, logs, ticks, or reads a test log. Its evidence came from a one-off classification of the main agents' tool calls in the logs flai serve kept to 2026-10-04, and nothing measured the turns after it, so no story of the epic could show its saving.

What agents spend reaches `flai stats` through each story's `usage` front matter, which flai serve measures from the stream-json logs of the agents it started for the story ([ADR-0051](0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)); S-0272 added the empty wakes there the same way, and left the classification of every turn to S-0293 ([ADR-0105](0105-a-story-s-empty-wakes-the-wait-for-events-calls-of-its-agents-that-timed-out.md)). The questions are what a turn is, what classes it falls in, where the counts are kept, and how `flai stats` reports them.

## Decision

**When flai measures a story's usage from its agents' logs, it puts each turn of the story's own agent in one of five classes and records the counts per UTC day in the story's usage as `turns`; `flai stats` reports them over the window, per day, and per story.**

1. **A turn** is one assistant message of the session's own agent, by its message ID, in an event without `parent_tool_use_id`, that calls at least one tool. Claude Code repeats a message once per content block, so a turn's calls are gathered across the repeats of its ID, and its day is the UTC date of its first event. A sub-agent's turns, and a message that calls no tool, are not counted.
2. **Its class** is the first of these it matches:
   1. `test_runs`: a Bash call runs a test, lint, or format tool by hand (`go test`, `go vet`, `gofmt`, `golangci-lint`, `vitest`, `svelte-check`, `eslint`, `prettier --check`, `markdownlint`, the `make` test and lint targets, and the scripts behind them), found in its command with heredoc bodies and quoted strings removed. `flai test`, `flai verify`, and the close-out are what replace them, and are not test runs.
   2. `hand_edits`: a call edits a story's narrative, a task, or a story by hand: an Edit, Write, or MultiEdit of its file, or a Bash call that names its path and writes, with `sed -i`, `perl -i`, `tee`, a redirection, or a Python script that writes.
   3. `empty_wakes`: every call is a `wait_for_events` that was an empty wake by ADR-0105.
   4. `ceremony`: every call is the MCP `inbox` or `item_move`, or a Bash call whose every command is one a story-loop command replaces: `git add`, `commit`, `status`, `log`, `diff --stat`, `show --stat`, or `rebase --continue`, or `flai stream sync`, `stream log`, `stream open`, `move`, `touches`, or `check`. The story-loop commands themselves (`task done`, `story start`, `stream state`, `criteria tick`, `test`, `verify`) are work.
   5. `work`: every other turn.
3. **Recorded in `usage.turns`**, a list of days in order, each with `day` and the count of each class, a count of 0 omitted and a day with no turns left out. A measurement replaces the list. A task carries none; an epic's usage, summed from its stories, carries their days summed. Stories measured before the release carry none until they are measured again, as `flai serve agent usage --all --write` does.
4. **`flai stats` reports them from the items**, archived stories included, whatever the report's type: under `turns` in `--json`, the classes, the total over the window, every day of the window with each class summed over the stories, and each story with turns in the window; and a section of the text report.

## Consequences

- Each E-0017 story's saving shows in `turns.days` once its release is installed, against the baseline the logs to 2026-10-04 give when measured again; the analyzer reads the same figures.
- The classes count turns, not calls or commands, and not cost: a turn that runs three ceremony commands is one ceremony turn, and a turn that runs a test and moves a task is a test run.
- The rules match commands by their text, so they miss what a script runs on its own (`scripts/foo.sh` that runs `go test` inside), and a tool run under another name is work.
- `usage.Usage` gains `Turns`, listed in `front-matter-fields.txt`, so the release that carries it raises `flai.minimum`, as `empty_wakes` did.
- `design/system/metrics.md` and `design/system/work-hierarchy.md` define the values.

## Alternatives considered

- **`flai stats` reads the logs itself.** Measuring every logged story on this host took 73 seconds over 412 MB of logs; every `flai stats` and every dashboard chart would pay it, and another clone has no logs. ADR-0105 rejected it for the same reasons.
- **Classify calls, not turns.** The epic's done condition is about turns, and a model call is what each costs; a turn of several calls is counted once.
- **Count every `wait_for_events` as an empty wake, as the epic's evidence did.** A wait that wakes on a change does something; ADR-0105's definition is kept, so the count and `usage.empty_wakes` agree.
