---
title: Priming with the digest, against the threads opened to the operator per completed story
updated: 2026-10-24
status: active
story: S-0435
---

# Priming with the digest, against the threads opened to the operator per completed story

The results of S-0435, an `experiment` story run on 2026-10-24, the last story of E-0025. E-0025 built a digest of the archive under `design/digest`, one entry per fact. Since S-0430, a story's pack loads the entries that bear on it ([ADR-0194](../adrs/0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md)), so that an agent does not ask the operator what a past story already settled. S-0434 gave `flai stats` the measure ([ADR-0200](../adrs/0200-flai-stats-reports-per-week-the-threads-opened-to-the-operator-the-stories.md)). This experiment reads that measure on this repository, before S-0430's acceptance and after it.

This is a preliminary reading. The story was planned for two whole weeks after the acceptance, which are whole from 2026-11-09. The story started one day after the acceptance and was blocked on TH-0642's answer. The operator unblocked it at 2026-10-24T23:43:59Z, so it measures the 26 hours 45 minutes that had passed.

## Hypothesis

Once the digest primes a story's agent, fewer threads are opened to the operator per completed story, and fewer of them ask what an entry already answered. The digest holds what past stories decided and what the operator answered. An agent that reads those answers in its pack has less to ask.

## Success measure

The hypothesis holds if both of the following are true:

- `digest.weeks[].per_story` in the after window is below the before window's rate. The gap must be wider than the spread between the before window's own weeks, so that it is more than the week-to-week noise.
- The after window's `asked_again` is lower than the before window's.

The thread waits that `waiting.weeks[].threads` reports are read beside the rate. Fewer questions should mean less time waiting on the operator.

## What was done

Nothing was built. The figures are those that `flai stats --json` gives, from the flai of this tree (`scripts/flai.sh`, with S-0434). They are read against the stories, threads, digest entries, and host journal behind them.

### The windows

S-0430 was accepted at 2026-10-23T21:00:06Z (`flai show S-0430`), in ISO week 2026-W43. `flai stats` takes `--since` alone, a window that ends now, and reports `digest` and `waiting` per ISO week.

| Window | Dates (UTC) | Weeks | Length |
|--------|-------------|-------|--------|
| Before | 2026-10-05 to 2026-10-18 | W41, W42 | 14 days |
| Left out | 2026-10-19 to 2026-10-23T21:00:06Z | W43 up to the acceptance | 4 days 21 hours |
| After | 2026-10-23T21:00:30Z to 2026-10-24T23:45:30Z | the rest of W43 | 26 hours 45 minutes |

- **Before.** The two whole ISO weeks before the week that holds the acceptance. The planner's rule takes a window's weeks as those whose start falls inside it. That rule would keep W42 alone, since W41 starts on 2026-10-05, four days before the 14 days that end at the acceptance. W41 is taken too, so that the window holds two whole weeks. Each week is reported on its own as well, so W42 can be read alone.
- **Left out.** W43 holds the acceptance, so as a whole week it belongs to neither window.
- **After.** No week of the planned after window, W44 and W45, had begun. The after window is instead the stretch from the acceptance to now, taken with a `--since` that starts it. flai reports it as a W43 bucket that holds only that stretch. It starts 24 seconds after the acceptance, so S-0430's own completion is left out.

### Commands

```bash
scripts/flai.sh stats --json --since 30d      # at 2026-10-24T23:45:08Z: W39 (part) to W43
scripts/flai.sh stats --json --since 26h45m   # at 2026-10-24T23:45:30Z: window_start 2026-10-23T21:00:30Z
scripts/flai.sh digest list --json            # the 382 entries, their status, asked_again, and cited
```

The figures for each window come from `digest.weeks[]` (`threads_to_operator`, `stories_completed`, `per_story`, `asked_again`), from `digest.entries`, and from `waiting.weeks[].threads` (`total_seconds`, `mean_seconds`, the mean being per completed story). The before window's totals are sums of its two weeks, and its rate is threads over stories.

## Results

### The measure

| Window | Threads to the operator | Stories completed | Threads per completed story | Asked again | Thread waits, total | Thread waits, mean per story |
|--------|------------------------:|------------------:|----------------------------:|------------:|--------------------:|-----------------------------:|
| W41 (2026-10-05) | 242 | 106 | 2.283 | 0 | 158,193 s (43h56m) | 1,492 s (24m52s) |
| W42 (2026-10-12) | 98 | 80 | 1.225 | 0 | 1,085,231 s (301h27m) | 13,565 s (3h46m) |
| **Before** (W41 and W42) | **340** | **186** | **1.828** | **0** | **1,243,424 s (345h24m)** | **6,685 s (1h51m)** |
| W43 whole, left out | 111 | 60 | 1.850 | 0 | 266,989 s (74h10m) | 4,450 s (1h14m) |
| **After** (26h45m) | **14** | **10** | **1.400** | **0** | **56,740 s (15h46m)** | **5,674 s (1h35m)** |

`digest.entries` was the same in both runs: 382 active, 0 compacted, 0 absorbed. Every entry was seeded by S-0429 on 2026-10-23 (commit `1981f3001`). The entries' `date` runs from 2026-09-16 to 2026-10-23.

Taking the after window out of W43 leaves 97 threads over 50 stories, 1.94 per story, for the days just before the acceptance. That is a subtraction of the two runs, not a figure flai gives.

### What else changed between the windows

Several other things changed between the two windows, and each could explain a difference in the rate without the digest.

- **The work.** The before window completed 186 stories: 64 features, 73 improvements, 44 remediations, 4 research stories, and 1 experiment. It also accepted 8 epics: E-0015 to E-0021 and E-0026. The after window completed 10 stories: 4 features, 2 improvements, and 4 remediations. It accepted no epic. Three more epics, E-0022, E-0027, and E-0028, were accepted in the left-out days of W43.
- **What the threads asked.** The threads' titles sort them roughly. This is a heuristic on the title, not a figure flai gives.

  | Kind of thread | Before | After |
  |----------------|-------:|------:|
  | A plan to approve (`… plan: …`) | 137 | 2 |
  | A cost of delay input or value | 61 | 0 |
  | A permission prompt (`Allow Write …`, `Allow Edit …`) | 14 | 8 |
  | A story only the operator may accept (ADR-0106) | 5 | 2 |
  | Anything else | 123 | 2 |

  Eight of the after window's 14 threads are permission prompts, opened as S-0431 and S-0432 wrote `.claude/agents/curator.md`, `.claude/settings.json`, and their template copies, and as S-0491 wrote `.claude/agents/orchestrator.md`. Two more ask the operator to accept S-0491 and S-0432, which change those protected files. No digest entry can answer a permission prompt or an acceptance. The other four are two plans (TH-0627, TH-0637), the refused release publish (TH-0641), and this story's own question of when to measure (TH-0642).
- **When the priming arrived.** A story's agent primes through the `flai` MCP server, which runs the host's installed flai. flai 1.68.0, the first release with S-0430, was tagged at 2026-10-23T23:39:07Z. The host's auto-upgrade failed six times from 23:40:49Z, and first installed it, as 1.68.1, at 2026-10-24T01:09:37Z (`~/.flai/serve/journal.jsonl`). Eight of the after window's 14 threads, TH-0627 to TH-0635, were opened before that time. An agent that primed through MCP in those four hours got no digest entries.
- **The orchestrator's permissions.** `orchestration.permissions.answer_threads` has been `autonomous` since 2026-10-07 (commit `9c28ced8a`), before both windows. The manifest's other changes were `analysis.agent` on 2026-10-15, `planning.agent` on 2026-10-21, a sub-project's install command (S-0483) on 2026-10-21, and the go-test tier's flags (S-0490) on 2026-10-24. None of them changes who opens a thread or who answers one.
- **The curator.** The curator has never run on this project. The host journal shows the `curate` host action turned on only for another project, at 2026-10-24T22:41:18Z. So no thread is in any entry's `asked_again`, and the repeat count is 0 in both windows by construction, not by measurement. This is recorded as [I-0207](../issues/I-0207-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md).

### Threads an entry answered

None of the after window's threads was answered by an entry.

- None of the 14 threads cites a `D-nnnn` in any entry. None is listed in any entry's `asked_again`.
- One decision in the window did rely on an entry instead of a thread. S-0431's narrative cites D-0363, the operator's answer on TH-0476, for the curator's design (`wip/archive/agents/S-0431.md`, `## Decisions`). That is a question not asked. The measure does not count it, because it counts only threads.
- The curator absorbed no entry and compacted none in the window: `digest.entries` is 382 active, 0 compacted, and 0 absorbed. S-0432 tried `flai digest cited` once in its worktree, which wrote `cited` into D-0001, D-0363, and D-0382. It then put the three files back from main (`wip/archive/agents/S-0432.md`, `## Decisions`).

### Verdict

The result cannot tell whether the digest works. The rate fell from 1.828 threads per completed story before to 1.400 after, a fall of 23%. That fall does not meet the success measure, for four reasons.

1. **The fall is inside the noise.** The before window's own weeks ran from 1.225 (W42) to 2.283 (W41). An after rate of 1.400 lies between them.
2. **The after window is too small.** It is 26 hours 45 minutes, with 10 stories and 14 threads. One thread more or less moves the rate by 0.1.
3. **The after threads are of other kinds.** Ten of the 14 are permission prompts and acceptance asks, which no entry can answer. In the before window, plan approvals and cost of delay inputs made 198 of 340 threads. The after window had 2 plans and no cost of delay input. Its rate tracks what the work asked, not what the agents knew.
4. **Most of the threads came without priming.** Eight of the 14 threads came before the host's flai carried S-0430's priming. The repeat count reads 0 because the curator never ran, not because no question was repeated.

The thread waits point the same way, and they are weaker still. The mean wait per completed story fell from 1h51m to 1h35m. But W42's three-hour mean is what lifts the before window, and the after window holds 10 stories.

## Recommendation

**Adopt: keep the digest as it is, and measure it again.** Nothing here argues for changing the digest or for dropping it. The one entry relied on in the window, D-0363 in S-0431, saved a question as the digest intends. The measurement, not the digest, needs two changes before the measure can say anything:

- **Turn the curator on for this project.** Without it, `asked_again` is 0 whatever is asked, and the repeat half of the hypothesis cannot be tested. This is [I-0207](../issues/I-0207-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md). The `curate` host action is the operator's to turn on.
- **Run the experiment again on the planned window.** Measure once W44 and W45 are whole, from 2026-11-09, with the same commands and the same before window. Count the plan and cost of delay threads apart from the rest, since those are what a digest entry can answer. Of all the kinds in the table above, they are the ones whose rate the digest should move. I-0207 carries the re-run; the operator can make it a story from the issue at acceptance.
