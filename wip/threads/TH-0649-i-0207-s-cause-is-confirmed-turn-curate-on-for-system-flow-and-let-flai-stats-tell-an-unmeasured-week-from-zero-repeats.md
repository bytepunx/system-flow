---
id: TH-0649
title: "I-0207's cause is confirmed: turn curate on for system-flow, and let flai stats tell an unmeasured week from zero repeats"
anchor:
  path: wip/kanban/tasks/T-2120-confirm-i-0207-s-cause-on-the-host-and-propose-the-remedy-to-the-operator-on-a-thread.md
  item: T-2120
status: resolved
participants: [agent-S-0504, orchestrator, alex]
created: 2026-10-26T04:12:27Z
updated: 2026-10-28T12:14:56Z
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

### 2026-10-28T12:02:08Z alex
Confirmed the recommendation of 2026-10-26T04:12:36Z orchestrator.

Source: design/adrs/0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md § Decision

### 2026-10-28T12:03:34Z agent-S-0504
Thank you. `curate` is now on for system-flow. **The curator has not yet completed a run here, so please run `flai curate` once** in the system-flow checkout and reply here.

- **What happened.** Accepting S-0498 and S-0501 started the curator at 2026-10-26T16:03:52Z. Claude Code refused it at once: `429 usage_limit_reached`, "You've hit your weekly limit · resets Oct 28, 2am". `serve/agents.json` records the run `failed (exit 1)`. It used no tokens, so it logged nothing in `wip/agents/curator.md`. Nothing started it again after the reset. I recorded this as I-0212.
- **Why I can't start it.** The MCP tool `curate` refuses an agent that flai serve started.
- **The code half works on this repository.** This branch's `flai stats` reads W40 to W44 as `asked again –` and `curator  runs 0 · never run`. Before, it read 0 in every week.
- **Once you reply:** I check that `wip/agents/curator.md` has its first entry and that `flai stats` shows `curator  runs 1` with the current week measured. Then I close I-0207, run the close-out, and move S-0504 to review. If you would rather let the next acceptance start the curator, say so: I will then close I-0207 on the code half and leave the first run to I-0212.

### 2026-10-28T12:13:40Z alex
let the next acceptance trigger the curator

### 2026-10-28T12:14:56Z agent-S-0504
Resolved: The operator turned curate on and chose to let the next acceptance start the curator's first run. S-0504 closed I-0207 on ADR-0204's measure, and I-0212 records the run that failed on the usage limit.
