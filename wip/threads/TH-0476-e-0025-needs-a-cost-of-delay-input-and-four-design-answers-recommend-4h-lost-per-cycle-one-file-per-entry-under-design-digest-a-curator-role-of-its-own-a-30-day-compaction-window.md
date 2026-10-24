---
id: TH-0476
title: "E-0025 needs a cost of delay input and four design answers: recommend 4h lost per cycle, one file per entry under design/digest, a curator role of its own, a 30-day compaction window"
anchor:
  path: wip/kanban/epics/E-0025-a-digest-of-the-archive-kept-by-an-agent-and-decayed-as-the-design-absorbs-it-primes-a-story-s-agent-with-what-past-stories-decided-and-what-the-operator-answered-so-the-same-question-is-not-asked-twice.md
  item: E-0025
status: resolved
participants: [planner-E-0025, orchestrator, alex]
created: 2026-10-16T05:57:18Z
updated: 2026-10-24T23:50:56Z
---

# TH-0476 E-0025 needs a cost of delay input and four design answers: recommend 4h lost per cycle, one file per entry under design/digest, a curator role of its own, a 30-day compaction window

On wip/kanban/epics/E-0025-a-digest-of-the-archive-kept-by-an-agent-and-decayed-as-the-design-absorbs-it-primes-a-story-s-agent-with-what-past-stories-decided-and-what-the-operator-answered-so-the-same-question-is-not-asked-twice.md.

## Entries

### 2026-10-16T05:57:18Z planner-E-0025
I am drafting E-0025's stories now. Five inputs are yours. Each has my recommendation first. I draft on the recommended answers and change the stories if you answer otherwise.

## 1. Cost of delay input (blocks the stories' values)

`flai cod E-0025` refuses: the epic has no inputs.

- **Recommended:** `time_lost_per_cycle: 4h` on E-0025. At `planning.hour_rate` 150 and the default 168h cycle, that is 600 USD a week, shared among the stories by duration. It counts the operator's time on questions asked again and the rework of a story that guessed against an earlier answer. I-0151's 117h a cycle is mostly agents idling, not work, so I take only a share of it.
- Alternatives: `2h` (300 USD a week), or `8h` (1200 USD a week).

## 2. Where the digest lives and how an entry is written (the epic's first open question)

- **Recommended:** one markdown file per entry, `design/digest/D-nnnn-<slug>.md`, named in the manifest as `design/analysis` is. The title is the fact in one sentence. Front matter holds `id`, `kind` (`decision`, `answer`, `recommendation`, `cause`), `sources` (S-, TH-, I- IDs), `topics`, `paths`, `date`, `status` (`active`, `absorbed`, `compacted`), `absorbed_by`, and `asked_again` (the threads that asked what it already answered). `doc_search` and the dashboard's documents page then find and render each entry as they do an ADR, and prime selects an entry by its topics and paths as it selects an ADR.
- Alternatives: one file per topic, each entry a section (an entry on two topics has to pick one); or front matter on the archived items (scatters the digest across `wip/archive`, which `doc_search` does not index).

## 3. Who keeps the digest (the epic's second open question)

- **Recommended:** a role of its own, the **curator**: role `curate`, host action `curate`, activity document `wip/agents/curator.md`, definition `.claude/agents/curator.md`, settings `digest.agent` and `digest.schedule`. It is built on today's analyzer machinery, as the epic's Notes allow, without waiting for E-0024. It writes entries only through `flai digest` commands, and the guard holds it to them.
- Alternatives: the analyzer keeps it, with a `digest` focus (cheaper, but its activity mixes with analyses and its guard grows a second folder); or the curator waits for E-0024 and is a role file after S-0419 (less code, but the epic then waits behind E-0024's twelve stories).

## 4. The compaction window (the epic's third open question)

- **Recommended:** `digest.compact_after: 30d`. An entry is needed when a narrative, thread, or work item cites its ID. An active entry nobody has cited in 30 days is compacted to its one-line title, which is still primed, as a brief is. Prime records nothing, so the measure stays in the repository.
- Alternatives: `60d`; or count an entry as needed whenever a pack included it (needs a use log under `.flai-cache`, lost on another host).

## 5. Who counts as the operator for an answer that becomes an entry at once

- **Recommended:** the manifest's `owner` (alex), as `permission_prompt` takes it (ADR-0097). Their reply to a thread they did not open, or their confirmation of a recommendation, becomes an `answer` entry. A reply that overrules a pending recommendation becomes a `recommendation` entry that says what was overruled and why.
- Alternative: any author that is not an agent run (`planner-*`, `orchestrator`, `analyzer`, `agent-S-*`, `flai`).

Meanwhile I plan everything that needs no answer. Once question 1 is answered, I write the stories' values.

### 2026-10-19T04:51:24Z planner-E-0025
An update. E-0025's eight drafts, S-0428 to S-0435, and their 42 tasks, T-1807 to T-1848, are complete except for their cost-of-delay values, which wait on question 1. The drafts already follow my recommended answers to questions 2 to 5. The plan, its order, and its layers are on TH-0480.

Since I opened this thread:

- **S-0426 was accepted.** Its finding, in design/system/agent-memory.md, supports a digest primed by topics and paths.
- **E-0021 and E-0026 are done.** Nothing in E-0025 waits on them.
- **No E-0024 role story has started.** So the curator still builds on the analyzer machinery (question 3).

The quickest answer is "yes to all five". I would then set `time_lost_per_cycle: 4h` on E-0025, write each story's value from `flai cod`, and leave the drafts for you to finalize. If any answer differs, name only that one, and I will change the stories it reaches.

### 2026-10-20T19:02:52Z planner-E-0025
**Question 3, revised recommendation.** You moved E-0024 to ready on 2026-10-20, and that changes what I would choose.

- **Recommended now:** the curator waits for E-0024 and becomes a role file. S-0431 goes after S-0419, and its triggers in S-0432 use S-0418's (on an event and on a schedule).
- **Why:** E-0024's role stories are ahead of E-0025's in the pull order. flai plays S-0419 out to 2026-10-21T03:45Z, before S-0428 even starts, so waiting delays nothing. Building the curator on the analyzer machinery first would mean S-0419 or a later story has to move it onto the role path again.
- **What it changes:** S-0431 becomes a `design/roles/curator.md` role file with its permitted `flai digest` calls, its pack topics, and its prompt. The new host action, guard rules, prompt code, and serve wiring drop out. S-0432's schedule and after-acceptance triggers become the role file's trigger. Both forecasts fall, from 1h23m and 55m to roughly 45m and 30m, which I would re-forecast with `flai forecast`. S-0430, S-0433, and S-0434 are unchanged.
- **The other two answers still stand:** a role of its own built now on the analyzer code (my first recommendation), or the analyzer with a `digest` focus.

The quickest answer is still one line, now "yes to all five, with question 3 as revised": `time_lost_per_cycle: 4h` on E-0025, one file per entry under `design/digest`, the curator as an E-0024 role file, `digest.compact_after: 30d`, and the manifest owner as the operator. I change the drafts only once you answer.

### 2026-10-21T05:57:49Z planner-E-0025
**Question 3, back to the first recommendation.** You moved E-0024 back to the backlog on 2026-10-21. My revision of 2026-10-20 rested on E-0024 being ready and ahead of E-0025, so it no longer holds.

- **Recommended again:** the curator is a role of its own, built on today's analyzer machinery, as my first entry describes, without waiting for E-0024. The drafts S-0431 and S-0432 already follow this, so nothing changes.
- **Why:** waiting for S-0419 would now hold the curator, and with it S-0432 to S-0435, behind an epic with no place in the pull order.
- If E-0024 goes ahead of E-0025 later, the role-file option comes back, and I would propose it again then.

The quickest answer is again "yes to all five", as in my second entry.

### 2026-10-21T11:04:41Z orchestrator (recommendation)
Recommended: "yes to all five", as planner-E-0025's entry of 2026-10-21T05:57Z puts it.

1. Set `time_lost_per_cycle: 4h` on E-0025. At 150 USD an hour over a 168h cycle, that is 600 USD a week, shared across S-0428 to S-0435 by forecast. This is the same figure you chose for E-0024 on TH-0471. Without inputs on the epic, `flai cod` refuses all eight stories, and none of them can be promoted.
2. Keep one file per entry under `design/digest`. Your TH-0471 answer moves roles and formulas to `design/workflows`, but it does not touch the digest.
3. Build the curator now as a role of its own on the analyzer machinery. E-0024 is back in the backlog, so waiting for S-0419 would hold S-0431 to S-0435. If E-0024 lands first, S-0419 moves the curator onto the role path (TH-0473).
4. and 5. Take the planner's recommendations as written into the drafts.

Confirming this sets no input by itself: the epic's input is yours to set, with `flai edit E-0025` and its cost of delay flags.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-21T12:59:17Z alex
Confirmed the recommendation of 2026-10-21T11:04:41Z orchestrator.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-24T23:50:56Z orchestrator
Resolved: S-0435 was accepted
