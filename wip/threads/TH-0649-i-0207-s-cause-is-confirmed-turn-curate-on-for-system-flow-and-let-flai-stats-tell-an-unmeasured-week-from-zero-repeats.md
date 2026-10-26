---
id: TH-0649
title: "I-0207's cause is confirmed: turn curate on for system-flow, and let flai stats tell an unmeasured week from zero repeats"
anchor:
  path: wip/kanban/tasks/T-2120-confirm-i-0207-s-cause-on-the-host-and-propose-the-remedy-to-the-operator-on-a-thread.md
  item: T-2120
status: open
participants: [agent-S-0504, orchestrator]
created: 2026-10-26T04:12:27Z
updated: 2026-10-26T04:12:36Z
---

# TH-0649 I-0207's cause is confirmed: turn curate on for system-flow, and let flai stats tell an unmeasured week from zero repeats

On wip/kanban/tasks/T-2120-confirm-i-0207-s-cause-on-the-host-and-propose-the-remedy-to-the-operator-on-a-thread.md.

## Entries

### 2026-10-26T04:12:27Z agent-S-0504
**Cause confirmed on the host.** `~/.flai/config.json` has `host_actions.curate` set for `/home/alex/git/bytepunx/signet` only. `wip/agents/curator.md` does not exist on main. So no curator has run here. No entry lists a thread in `asked_again`, and every week reads `asked_again 0`.

**Second cause, in flai.** `metrics.digestSection` reads only threads and entries. Nothing in the `digest` section says whether a curator ever ran. A project with no curator therefore reads 0 repeats, and the analyzer and the Digest chart take that 0 as a measurement. Every project hits this until its curator runs, not only this one.

**Recommendation: fix both.**

1. **You, on the host:** run `flai serve enable curate` in the system-flow checkout, then `flai curate` to start the first run now, rather than wait for the next acceptance. Only you can turn on a host action. An agent `flai serve` started cannot start the curator.
2. **This story, in code:** an ADR and `metrics.md` § Digest add the curator's runs to the section:
   - `digest.weeks[].curator_runs`: the curator runs that ended in the week, from `wip/agents/curator.md`.
   - `digest.curator`: `runs` and `last_run`.
   - `digest.weeks[].asked_again` is left out of a week no curator run has reached, meaning no run ended on or after its Monday. This matches how `per_story` is left out of a week with no completed story.
   - `flai stats` prints `asked again –` for such a week, and a `curator  runs 3 · last <time>` line.
   - The Digest chart and the docs follow.

**Alternative:** turn curate on and change no code. That clears I-0207 here but leaves the same false 0 on every other project.

I am building item 2 on this recommendation now, on the story's branch. If you answer otherwise, I change the tasks. I close I-0207 only once the curator has run here and `flai stats` shows that run. Please reply after you run `flai curate`, or say if you would rather wait for an acceptance to trigger it.

### 2026-10-26T04:12:36Z orchestrator (recommendation)
**Recommended for the operator:** take agent-S-0504's recommendation and fix both causes. Run these on the host, in the system-flow checkout:

```sh
flai serve enable curate
flai curate
```

Then reply here once the curator has run.

- **Why it needs you.** A host action is the operator's to turn on, and an agent that `flai serve` started cannot start the curator. The question is a host change, so I escalate it rather than answer it.
- **Why `flai curate` now.** ADR-0199 § Decision makes `curate` start a run only on an acceptance seen while the action is on, or on `digest.schedule`. An acceptance made while `curate` is off starts nothing, then or later. Without a manual run, the first run waits for the next acceptance after you enable it.
- **Why the code half too.** Without it, every project whose curator has not yet run reads `asked_again 0` as if measured. The ADR the story writes keeps that 0 from being read as a measurement.

The story builds item 2 meanwhile. Only T-2122 and T-2126 wait on your reply.

Source: design/adrs/0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md § Decision
