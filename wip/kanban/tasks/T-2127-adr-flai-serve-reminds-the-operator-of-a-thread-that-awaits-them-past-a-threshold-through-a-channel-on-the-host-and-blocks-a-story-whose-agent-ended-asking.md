---
id: T-2127
type: task
nature: improvement
title: "ADR: flai serve reminds the operator of a thread that awaits them past a threshold through a channel on the host, and blocks a story whose agent ended asking"
status: done
parent: S-0498
owner: alex
created: 2026-10-26T04:13:26Z
updated: 2026-10-26T04:19:30Z
transitions:
  - to: ready
    at: 2026-10-26T04:18:41Z
    by: agent-S-0498
  - to: in-progress
    at: 2026-10-26T04:18:42Z
    by: agent-S-0498
  - to: done
    at: 2026-10-26T04:19:30Z
    by: agent-S-0498
stream: S-0498
tags: [flai]
touches: [design/adrs/README.md, design/adrs/0205-flai-serve-reminds-the-operator-of-a-thread-that-awaits-them-past-threads.md]
usage:
  source: log
  seconds: 48
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 47
      cache_read: 705989
      cache_write: 10237
      cost: 0.2423
---
# T-2127 ADR: flai serve reminds the operator of a thread that awaits them past a threshold through a channel on the host, and blocks a story whose agent ended asking

## Work

Write the ADR that records the two decisions this story builds, with `flai adr new` (or the MCP tool `adr_new`), and accept it. It refines ADR-0090 and ADR-0091 (a recommendation ends no wait; the operator confirms it), ADR-0154 (an agent ends its session on an open question), ADR-0108 (a blocked story is not restarted on its own), and ADR-0158 (the pattern of a periodic look in `flai serve` that asks the operator). Topics `cli` and `orchestration`.

Context: I-0151's five instances, from TH-0324, TH-0390, TH-0445, TH-0436, and TH-0461. The orchestrator recommended within a minute each time; the operator confirmed 17 to 113 hours later, while the story sat in progress and held what named it in `after` or overlapped its claim. Three of the five needed the operator's hands on the host (a key, an `rm` under `.git`, `flai host upgrade`), which no setting may confirm for them; S-0502 takes the rest with `orchestration.permissions.confirm_after`.

Decision, in numbered points:

1. **The manifest's `threads` block.** `threads.remind_after`, a duration, unset by default, and `threads.remind_every`, a duration, default `24h`. With `remind_after` unset nothing is reminded.
2. **What is reminded.** A thread the operator's inbox marks awaiting the operator (ADR-0181: a pending recommendation, or a last entry by neither the project's owner nor its story's agent), on an item that is not done or cancelled, whose last entry is older than `remind_after`. A reminder repeats every `remind_every` while the thread still awaits them. A new entry on the thread starts the count over.
3. **How it reaches them.** Through a channel kept in the host configuration (`~/.flai/config.json`, ADR-0008), set in a shell with `flai serve notify set --url <webhook>` and/or `--command <shell command>`, as `flai serve agent provider` keeps the providers (ADR-0129). The body is JSON: the project, the thread's ID, title, item, how long it has awaited them, the pending recommendation's first line and source when there is one, and the dashboard link `/items/<id>?thread=TH-nnnn`. The URL gets a POST, one attempt, a short timeout; the command gets it on standard input. With no channel set, `flai serve` logs the reminder and warns once per project that no channel is set. The host's token and no file content are sent.
4. **Never on the thread.** A reminder is written nowhere in the repository: a thread entry would end the agent's wait (ADR-0091) and change the thread's status (ADR-0090). What was reminded, and when, lives in the serve folder, `reminders.json` per project, as `upgrade-asks.json` does, so a restart repeats nothing.
5. **The idle is marked.** When `flai serve` judges a story's run as ended asking, with the thread it asked on, it appends a blocked interval to the story with the reason `awaiting the operator on TH-nnnn`, unless the story is blocked already; when the answer starts the agent again, it closes the interval. The board shows the story blocked with its reason, and `flai stats` counts the wait as blocked time. The story stays in its column (workflow.md § Blocking) and still counts to the in-progress limit.

Consequences and alternatives: a reminder in the thread itself (rejected: it ends the wait); the dashboard's `dashboard.notify_url` notifier (rejected as the sender: it runs only while the container runs, posts once, and cannot run a host command); confirming after a timeout (S-0502's, which these instances could not use); not counting a blocked story to the WIP limit (rejected: an ADR of its own, and the lanes were not full in W42).

Add the ADR to `design/adrs/README.md` as `flai adr new` does. Predicted path `design/adrs/0204-...`: `flai task done` widens the touch to the real name.

## Done when

- The ADR exists under `design/adrs/` with status accepted, refining ADR-0090, ADR-0091, ADR-0108, ADR-0154, and ADR-0158, and is in the ADR index.
- Its decision names the two settings and their defaults, the awaiting rule it reminds on, the channel and the JSON body, the serve-folder state, and the blocked interval and its reason.
- `flai check --strict` and the markdown lint pass on it.

## Notes

Layer 1: waits for nothing. The code tasks wait for it, so the decision is recorded before anything is built on it; if the operator changes a point on the plan's thread, change the ADR first.
