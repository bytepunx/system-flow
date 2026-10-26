---
id: ADR-0205
title: "flai serve reminds the operator of a thread that awaits them past threads.remind_after through a channel kept on the host, writing nothing on the thread, and blocks a story whose agent ended asking until the answer starts it again"
status: accepted
date: 2026-10-26
supersedes: []
superseded_by: []
refines: [ADR-0090, ADR-0091, ADR-0108, ADR-0154, ADR-0158]
topics: [cli, orchestration]
---

# ADR-0205 flai serve reminds the operator of a thread that awaits them past threads.remind_after through a channel kept on the host, writing nothing on the thread, and blocks a story whose agent ended asking until the answer starts it again

## Context

[I-0151](../issues/I-0151-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md) found five threads in W42 on which the orchestrator recommended within a minute and the operator confirmed 17 to 113 hours later: TH-0324 (S-0232), TH-0390 (S-0337), TH-0445 (S-0380), TH-0436 (S-0378), and TH-0461 (S-0391). Meanwhile each story sat in progress, unmarked, holding what named it in `after` or overlapped its claim. Three of the five needed the operator's hands on the host (a cosign key, an `rm` under `.git`, `flai host upgrade`), which no setting may confirm for them; S-0502 takes the design-choice waits with `orchestration.permissions.confirm_after`.

What took a day was the operator seeing the recommendation. It sat in the dashboard's inbox, and nothing reached them off it. [ADR-0090](0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md) and [ADR-0091](0091-a-recommendation-ends-no-thread-wait-and-flai-stats-reports-the-waits-the.md) keep the wait open until the operator confirms. [ADR-0154](0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md) lets a story's agent end its session on an open question, and `flai serve` starts it again on the answer. Nothing marks the story while it waits: the convention asks the agent to run `flai block`, and none of the five did, so `flai stats` counted no blocked time on any of them. [ADR-0158](0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md) already gives `flai serve` a per-minute look per project that asks the operator and keeps its state in the serve folder.

## Decision

**`flai serve` reminds the operator, through a channel kept in the host configuration, of each thread that has awaited them longer than `threads.remind_after`, and again every `threads.remind_every`, and writes nothing on the thread; and it blocks a story whose agent ended its session asking on a thread, with the reason `awaiting the operator on TH-nnnn`, until the answer starts the agent again.** This refines ADR-0090, ADR-0091, ADR-0108, ADR-0154, and ADR-0158; the rest of each stands.

1. **The manifest's `threads` block.** `threads.remind_after` is a Go duration, unset by default: with it unset nothing is reminded. `threads.remind_every` is a Go duration, `24h` when unset. `flai check` reports a value that is not a duration, a negative one, or a `remind_every` of zero as `manifest.threads`. Both are set with `flai manifest set`.
2. **What is reminded.** A thread the inbox marks awaiting the operator ([ADR-0181](0181-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the.md): a pending recommendation, or a last entry by neither the project's owner nor its story's agent), unresolved, on a document or on a work item that is not done or cancelled, whose last entry is older than `remind_after`. It is reminded again every `remind_every` while it still awaits them. A new entry on the thread starts the count over from that entry.
3. **How it reaches them.** Through a channel in the host configuration (`~/.flai/config.json`, [ADR-0008](0008-cli-config-location.md)), set in a shell with `flai serve notify set --url <webhook>` and/or `--command <shell command>`, shown with `flai serve notify show`, removed with `flai serve notify clear`, as `flai serve agent provider` keeps the providers ([ADR-0129](0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)). The reminder is JSON: `project`, `thread`, `title`, `item`, `awaiting_seconds`, `since`, `recommendation` with its `author`, `first_line`, and `source` when one is pending, and `href`, the dashboard link `/items/<id>?thread=TH-nnnn`. The URL gets a POST, one attempt, a 5 s timeout, a non-2xx being a failure; the command runs through the host's shell with the reminder on standard input and a 30 s timeout. With no channel set, `flai serve` logs the reminder and warns once per project per run that no channel is set. No token and no file content is sent. The body extends the one the dashboard posts to `dashboard.notify_url`, so one receiver can take both.
4. **Never on the thread.** A reminder is written nowhere in the repository: a thread entry would end the agent's wait (ADR-0091) and change what the thread awaits (ADR-0090). What was reminded, and when, is kept per project in `reminders.json` in the serve folder, beside `upgrade-asks.json`: per thread, the time of the entry the reminder was for and when it was last sent. A restart of `flai serve` repeats nothing. A thread that no longer awaits the operator is dropped from the file at the next look.
5. **The idle is marked.** When `flai serve` judges a story's run as ended asking with the thread it asked on, it appends a blocked interval to the story with the reason `awaiting the operator on TH-nnnn`, unless the story is blocked already. When the answer starts the agent again, it closes that interval before the start; an interval whose reason names no such thread, one an agent or the operator wrote, is left open. A run that ended awaiting only another story's reply blocks nothing. The story stays in its column ([workflow.md § Blocking](../system/workflow.md)) and still counts to the in-progress limit. [ADR-0108](0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md)'s rule stands: a blocked story is not restarted on its own; the asked run is started again on its answer alone, as before.

## Consequences

- The operator learns of a pending recommendation, or a question, `remind_after` after it was written, off the dashboard: on a phone through a push service, or as a desktop notification on the host.
- The board shows which stories wait on the operator and on which thread, and `flai stats` counts those waits as blocked time.
- The agent need no longer run `flai block` for a question it ends its session on; the work-management convention says so.
- Nothing changes for a project that leaves `threads.remind_after` unset, apart from the block on an asked story.
- The channel is the operator's to set, in a shell: no agent writes the host configuration.

## Alternatives considered

- **A reminder as a thread entry.** Rejected: it ends the agent's wait and changes what the thread awaits.
- **The dashboard's `dashboard.notify_url` notifier as the sender.** Rejected: it runs only while the container runs, posts once per entry, and cannot run a host command.
- **Confirming a recommendation after a timeout.** S-0502's, which three of these instances could not use, since they needed the operator's hands.
- **Leaving a blocked story out of the in-progress limit.** Not decided here: it is a decision of its own, and the lanes were not full in W42.
