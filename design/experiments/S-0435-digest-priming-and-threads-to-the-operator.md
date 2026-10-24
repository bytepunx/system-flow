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
