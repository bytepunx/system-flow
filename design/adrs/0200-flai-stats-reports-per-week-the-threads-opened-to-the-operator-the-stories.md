---
id: ADR-0200
title: "flai stats reports per week the threads opened to the operator, the stories completed, the threads per completed story, and those that asked what a digest entry already answered"
status: accepted
date: 2026-10-24
supersedes: []
superseded_by: []
refines: [ADR-0081, ADR-0091]
topics: [cli, dashboard, analysis]
---

# ADR-0200 flai stats reports per week the threads opened to the operator, the stories completed, the threads per completed story, and those that asked what a digest entry already answered

## Context

E-0025 builds a digest of the archive, one entry per fact under `design/digest` ([ADR-0191](0191-the-digest-is-one-markdown-file-per-entry-under-the-design-folder-s-digest.md)), that primes a story's agent with what past stories decided and what the operator answered ([ADR-0194](0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md)), so that the same question is not asked twice. I-0151 asks how anyone would know whether it worked. The curator records in an entry's `asked_again` each thread that asked what the entry already answered ([ADR-0196](0196-the-curator-runs-behind-the-curate-host-action-and-takes-in-what-the-archive.md), [ADR-0199](0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md)), but nothing counts those threads, or the threads the operator is asked at all, against the work done.

[ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md) has `flai stats` report what the dashboard charts, so that the dashboard works out no figures of its own, and its waiting section counts the stories completed per week. [ADR-0091](0091-a-recommendation-ends-no-thread-wait-and-flai-stats-reports-the-waits-the.md) tells the waits by who ended them, but not how many questions were asked. `design/system/metrics.md`, the contract between `flai stats` and the dashboard, changes only with an ADR. S-0435, the experiment, compares the two weeks before the digest with the two weeks after, so the measure must count weeks before any entry existed, from the threads alone.

## Decision

**`flai stats` reports under `digest`, per ISO week of the window, the threads opened to the operator, the stories completed, the threads per completed story, and how many of those threads asked what a digest entry already answered, with the digest's entries counted by status.** It reports them as `metrics.md` § Digest defines them, refining ADR-0081 and ADR-0091.

- **A thread opened to the operator.** A thread under `wip/threads`, open or resolved, whose first entry's author is neither the operator, the manifest's `owner`, nor flai itself, `flai` or `flai serve`, which write notices such as the upgrade thread rather than questions. The planner, the orchestrator, a story's agent, and anyone else count. A thread with no entry has no opener and is not counted.
- **Its week.** The ISO week that holds its first entry's time, from the window's start to now, as every weekly series is.
- **The stories completed.** The stories completed in the week, cancelled ones left out, as `waiting.weeks[].items` counts them in a report of stories. The section counts stories whatever the report's type, since the questions are asked per story.
- **The rate.** `per_story`, the threads opened to the operator over the stories completed, not rounded. Absent in a week with no completed story: such a week has no rate, which is a gap, not 0.
- **A repeat.** One of the week's threads opened to the operator that some entry lists in `asked_again` and whose `date` is a day before the day of the thread's first entry, in UTC. An entry's `date` is a day, so an entry of the same day, which may have been written after the thread, never makes a repeat. An entry of any status counts: a compacted or absorbed entry had answered the question when it was asked. A thread several entries list is one repeat.
- **The entries.** `digest.entries` counts the digest's entries now by status, `active`, `compacted`, and `absorbed`, over the whole digest, not the window.
- **Without a digest.** A project with no digest folder reports its weeks with no repeat and every entry count 0, so the weeks before the digest are measured too. A digest flai cannot read is warned of on stderr, and `digest` is absent.

## Consequences

- The dashboard's Digest chart reads `flai stats --json` alone, and matches it. A flai older than this release sends no `digest`, and the chart says so.
- The rate counts every question to the operator, not only those the digest could have answered. A fall in it after the digest is evidence, not proof: S-0435 says how much it can bear.
- A repeat is counted only once the curator has listed the thread in an entry's `asked_again`. Until its next run, this week's repeats read low.
- A thread opened by the operator, or by flai's own notices, is never counted, however it reads. A question an agent asks in a conversation with another story's agent is not a thread and is not counted either.
