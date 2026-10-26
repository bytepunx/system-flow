---
id: T-2120
type: task
nature: research
title: Confirm I-0207's cause on the host and propose the remedy to the operator on a thread
status: done
parent: S-0504
owner: alex
created: 2026-10-26T04:08:51Z
updated: 2026-10-26T04:12:31Z
transitions:
  - to: ready
    at: 2026-10-26T04:12:15Z
    by: agent-S-0504
  - to: in-progress
    at: 2026-10-26T04:12:16Z
    by: agent-S-0504
  - to: done
    at: 2026-10-26T04:12:31Z
    by: agent-S-0504
stream: S-0504
tags: [digest, curator]
usage:
  source: log
  seconds: 15
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1672
      cache_read: 259819
      cache_write: 10540
      cost: 0.1514
---
# T-2120 Confirm I-0207's cause on the host and propose the remedy to the operator on a thread

## Work

Confirm the cause from the host, not from the issue's words: `flai serve status` for whether the `curate` action is on for this project, `~/.flai/serve/journal.jsonl` for any `curate` action on it (read it, do not change it), and `ls wip/agents/curator.md` in the main checkout, which does not exist while the curator has never run. Read `design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md` § What else changed between the windows and § Recommendation, and `design/system/metrics.md` § Digest.

Then open one thread on S-0504 with `thread_open` proposing the remedy, the recommendation first, as the story's goal asks before anything is built:

1. The operator turns the action on for system-flow with `flai serve enable curate` in a shell on the host, or on the Settings page while `settings` is on. The agent cannot: a host action is the operator's.
2. `flai stats` tells an unmeasured repeat count from a measured zero: `digest.weeks[]` gains `curator_runs`, the curator's runs ended in the week, read from `wip/agents/curator.md`, and leaves `asked_again` out of a week no curator run has ended since its start, as `per_story` is left out of a week with no completed story; `digest` gains `curator` with `runs` over the window and `last_run`. An ADR records it, since `metrics.md` changes only with one.
3. The alternative to name: leave `flai stats` as it is and only turn the action on. Say why it is weaker: the measure would still read 0 on any project whose curator has not run, and the analyzer and the dashboard would still read that 0 as a result.

Record the answer in the narrative's `## Decisions`. Hold `wait_for_events` for the answer; when it changes the remedy, change the tasks of this story with `flai edit` and say why in `## Decisions`. The tasks after this one wait on it because the remedy must be agreed before an ADR is written.

## Done when

- The narrative's `## Decisions` records whether `curate` was on for system-flow, what the journal showed, and whether `wip/agents/curator.md` existed, each with the command that showed it.
- A thread on S-0504 holds the proposal with the recommendation first and the alternative, and the operator has answered it.
- The tasks of this story match the answer.

## Notes

Drafted by planner-S-0504. No file of the repository changes in this task, so it has no touches.
