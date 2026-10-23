---
id: S-0435
type: story
nature: experiment
title: "Experiment: does priming with the digest cut the threads opened to the operator per completed story on this repository, two weeks before against two weeks after"
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:41Z
updated: 2026-10-23T17:53:06Z
transitions: []
tags: [cli]
topics: [analysis, planning]
touches: [design/experiments/README.md, design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md, design/issues/summary.md]
after: [S-0434]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 450
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 110
          output: 1916
          cache_read: 16155300
          cache_write: 22978
          cost: 3.5585
cost_of_delay:
  value: 60.95
  by: planner-S-0435
  at: 2026-10-23T17:37:52Z
forecast:
  duration: 45m
  delivery: 2026-10-24T14:00:00Z
  basis: "Its own forecast of 45m; 26th in the pull order with an in-progress limit of 5, behind S-0429, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0430, S-0431, S-0432, S-0433 and S-0434."
  by: flai
  at: 2026-10-23T17:53:06Z
finalized:
  by: orchestrator
  at: 2026-10-23T17:38:15Z
---
# S-0435 Experiment: does priming with the digest cut the threads opened to the operator per completed story on this repository, two weeks before against two weeks after

## Goal

Find out, on this repository, whether the digest does what E-0025 built it for. The hypothesis is that once the digest primes (S-0430's acceptance), fewer threads per completed story are opened to the operator, and fewer of them ask what an entry already answered. The measure is S-0434's, read from `flai stats --json` over the two weeks before S-0430 was accepted and the two weeks after, with the results document under `design/experiments/`.

## Acceptance criteria

- [ ] The results document, `design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md` from `design/experiments/template.md`, states the hypothesis, the two windows by date, and the method: `flai stats --json` with each window, the `digest.weeks[]` figures, and the waits `waiting.weeks[].threads` reports, as flai gave them.
- [ ] It reports for each window the threads opened to the operator, the stories completed, the threads per completed story, the threads that asked what an entry already answered, and the thread waits' total and mean, in one table.
- [ ] It names what else changed between the windows that could explain a difference (other epics accepted, the orchestrator's `answer_threads` setting, the stories' mix by nature), and says whether the result supports the hypothesis, refutes it, or cannot tell, and why.
- [ ] It lists three threads of the after window that an entry answered and the entry each cited, or says that none did, and the entries the curator absorbed or compacted in the window.
- [ ] It ends with a recommendation: keep the digest as it is, change it (naming what, as an issue filed with `flai issue new` and linked), or drop it; and `design/experiments/README.md` lists the document.

## Tasks
- T-1847 The results document states the hypothesis and the two windows and gives the measured figures for each, as flai stats reports them
- T-1848 The results document weighs what else changed, cites threads an entry answered, says whether the hypothesis holds, and recommends keeping, changing, or dropping the digest

## Notes

- Drafted by the planner for E-0025: the epic's last story.
- It cannot start before 14 days have passed since S-0430's acceptance, since the after window must be whole. The operator promotes it no earlier than that; T-1847 blocks the story if it starts sooner.
- An experiment is accepted like any other story and releases nothing (ADR-0066).

### Planning

Planned again on 2026-10-23 by planner-S-0435, with S-0428 accepted (ADR-0191 and `flai digest list --json` exist) and E-0025's cost of delay input set. The two tasks of 2026-10-16 and their layers stand; their touches and `after` were checked and left as they are. The plan is on the thread this run opened on S-0435.

- **Touches.** Three, all declared, all files, each a task's: T-1847 writes the results document; T-1848 writes it too and adds `design/experiments/README.md` and `design/issues/summary.md`, which `flai issue new` regenerates when the recommendation files a change.
  - **From design:** `design/experiments/template.md` and the folder's README, as ADR-0066 and the documentation standard ask of an experiment story; `flai check` validates the document's front matter (S-0194).
  - **Co-change, kept out:** everything `flai touches suggest` lists (issue files at 7% and below, `flai-cli.md`, `docs/users/flai.md`) co-changes with `design/issues/summary.md`, which almost every close-out touches; none of it is this story's to change.
  - **No folder touch.** The issue file a recommendation may add goes under `design/issues`, which is in `claims.shared` and so holds nothing; its name is flai's to allocate, and the close-out commits it.
- **Forecast.** 45m, delivery 2026-11-07T08:25Z. flai's figure was 14m (median 100 s per unit over 297 done stories, size 8), played out to 2026-10-24T13:07Z. Raised: the two experiments done here, S-0176 and S-0427, took 1h46m and 1h31m of agent time, and both built and ran a trial; this one reads its figures from `flai stats`, the manifest's history, the threads, and the entries, and argues them in one document, about half of that. The delivery is not flai's play-out, which ignores the after window: it is S-0430's forecast delivery (2026-10-24T07:40Z) plus 14 days plus the duration, and moves with S-0430's acceptance.
- **Cost of delay.** 60.95 USD a week, flai's figure: S-0435's share of E-0025's 600 USD a week (`time_lost_per_cycle: 4h`, the answer the operator confirmed on TH-0476), 45m of 7h23m forecast over the epic's 7 open stories without inputs. Not adjusted: the story's weight in the epic is its duration.
- **Assumption on the windows.** `flai stats` takes `--since` alone, a window that ends now, and reports `digest.weeks[]` and `waiting.weeks[]` per ISO week. The story's agent runs it once with a `--since` long enough to cover both windows (28d from S-0430's acceptance plus the days since) and takes for each window the weeks whose `start` falls inside it; a week that straddles S-0430's acceptance belongs to neither and is named in the method. This is proposed as an edit to T-1847's Work on the plan thread.
