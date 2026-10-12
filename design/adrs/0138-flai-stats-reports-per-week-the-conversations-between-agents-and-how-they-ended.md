---
id: ADR-0138
title: "flai stats reports per week the conversations between agents and how they ended, the conflicts a trial merge found against those a rebase met at sync or acceptance, from a conflict log under .flai-cache, and the stories started on a share with the hold each would have waited"
status: accepted
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0081, ADR-0113]
topics: [cli, dashboard]
---

# ADR-0138 flai stats reports per week the conversations between agents and how they ended, the conflicts a trial merge found against those a rebase met at sync or acceptance, from a conflict log under .flai-cache, and the stories started on a share with the hold each would have waited

## Context

E-0018 sets out to run more stories at once with less conflict work at the end. The agents of two open stories talk in a conversation per pair ([ADR-0120](0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md)), flai tells them there of a trial-merge conflict or a claim grown into another's ([ADR-0121](0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md)), and a holding story's agent may share the overlapping paths so that the held story starts ([ADR-0134](0134-flai-serve-asks-the-agent-of-a-story-in-progress-about-each-ready-story-it.md)). S-0337 measures whether it works. [ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md) has `flai stats` report what the dashboard charts, so that the dashboard works out no figures of its own, and [ADR-0113](0113-flai-stats-reports-held-stories-per-day-held-time-by-reason-per-week-and-the.md) added the time held per week. `design/system/metrics.md`, the contract between `flai stats` and the dashboard, changes only with an ADR.

Three things stand in the way. A conversation has no kind, can mix notices with talk, and its `updated` moves with every entry, so neither its `created` nor its `updated` says when it agreed or ended. A rebase that stops on conflicts at `flai stream sync` or `flai accept` is recorded nowhere: once it is settled, git keeps no trace of it. And the hold replay behind `held_seconds` ignores shares, so the hold a share cut short is not reported.

## Decision

**`flai stats` reports, under `coordination` and per ISO week of the window, the conversations between stories and how they ended, the conflicts a trial merge found while both stories were open against those a rebase met at sync or at acceptance, read from a new conflict log, and the stories started on a share with the hold each would have waited.** It reports them as `metrics.md` § Coordination defines them, refining ADR-0081 and ADR-0113.

- **Conversations**, from the files under `wip/messages`, each event in the ISO week of the time in its entry's heading, the time it happened:
  - `opened`: the conversation's first entry.
  - `agreed`: the first of a share in its `shares[]`, at the share's `at`, and an entry closing it because the two branches merge cleanly at a sync. Each conversation counts once.
  - `escalated`: its first entry written by `flai message escalate`. Each conversation counts once.
  - `closed_with_story`: the entry closing it because a story was accepted, archived, or no longer open. A conversation read as closed with no such entry counts at the first done or cancelled transition of either story.
- **Conflicts**, from `.flai-cache/conflicts.jsonl` in the main checkout's cache, beside `overlaps.jsonl`. It holds one JSON line per event, with `at`, `story`, `other` (on a trial-merge line only), `where` (`trial`, `sync`, or `accept`), and `paths`.
  - A `trial` line is written only when the trial merge at `flai stream sync` sends or changes a conflict message. The conversation already skips a message naming the same paths as the newest, so a pair counts once per distinct conflict, not once per sync.
  - A `sync` line is written where `storygit.Rebase` stops on conflicts, which `flai stream sync` and `flai task done` reach. An `accept` line is written where `flai accept`'s rebase stops, before it refuses.
  - Generated files the sync settles itself ([ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)) are left out of `paths`, and a stop settled by them alone writes no line.
  - The log is append-only and never trimmed, so a window of any length can be counted from it. A failed write is logged as a warning and the command goes on: the log measures, it does not gate. `overlaps.jsonl` differs on both counts: it drops a failed write silently and trims to 500 lines.
  - A missing log means no conflicts. An unreadable one is warned of, `logged` is false, and the weeks carry no `conflicts`.
- **Shares**, from the board's transitions and the conversations' `shares[]`. A story started on a share is one whose first move to in-progress falls while a share names it as `held`: the share's `at` is at or before the move, and its holder has not yet left in-progress. It counts in that move's week, with the hold it would have waited, from the move to the holder's first move to review, done, or cancelled after it, or to now while the holder is still in progress. This is the counterpart of ADR-0113's `held_seconds`, the hold that ran.
- **Not counted**: a narrowed claim, and agreement in free text. Nothing records either against the conversation: a narrowed claim is an edit of a story's touches, and a reply that agrees is words. Counting them would be a guess.

## Consequences

- The Coordination chart reads `flai stats --json` alone, and matches it.
- `agreed` is a floor. A pair that agreed in a reply, or by narrowing a claim, shows as opened and later closed with its story.
- Conflicts are counted from the release that writes the log. Weeks before it show none found and none met, which is no evidence either way.
- The log grows by a line per conflict event and is never trimmed. A write that fails loses that event from the stats; the warning says so.
- A story's saved hold counts whole in the week it started, not split over the weeks the hold would have spanned, so it grows while its holder is still in progress and settles once the holder moves on.

## Alternatives considered

- **Count conversations by `created` and `updated`.** `updated` moves with every entry, and a conversation can mix notices with agreement, so the week it agreed or ended is lost.
- **Read conflicts from `overlaps.jsonl`.** It holds notices of overlapping claims, not conflicts, drops failed writes silently, and trims to 500 lines, so a long window would undercount.
- **A trial-merge line at every sync that finds a conflict.** A pair that syncs ten times over one conflict would count ten.
- **Count agreement from free text.** It would guess from words, and two readers could count it differently.
- **Split the saved hold over the weeks it spans, as `claims.weeks[].held_seconds` is split.** The saving is known only once the holder moves on, and it belongs to the story counted in `started`; counting it in that week keeps the two together.
