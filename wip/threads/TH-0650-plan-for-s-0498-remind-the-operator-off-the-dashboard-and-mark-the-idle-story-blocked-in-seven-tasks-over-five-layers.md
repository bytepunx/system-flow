---
id: TH-0650
title: "Plan for S-0498: remind the operator off the dashboard and mark the idle story blocked, in seven tasks over five layers"
anchor:
  path: wip/kanban/stories/S-0498-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md
  item: S-0498
status: resolved
participants: [planner-S-0498, orchestrator]
created: 2026-10-26T04:17:16Z
updated: 2026-10-26T04:17:37Z
---

# TH-0650 Plan for S-0498: remind the operator off the dashboard and mark the idle story blocked, in seven tasks over five layers

On wip/kanban/stories/S-0498-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md.

## Entries

### 2026-10-26T04:17:16Z planner-S-0498
Plan by planner-S-0498. S-0498 is a draft for you to finalize. I added tags `flai` and `template`, topics `orchestration` and `cli`, 34 file touches, a forecast of 1h45m (flai's 54m raised; the reason is under `### Planning`), and the cost of delay value 17550 USD a week from the 117h input flai carried over from I-0151.

## The proposed solution

I-0151's five instances, read on their threads:

| Thread | Waited | What it needed | Would S-0502's `confirm_after` have ended it |
|--------|--------|----------------|------------------------------------------------|
| TH-0324 (S-0232) | 112.7 h | your cosign public key, two scope choices | no: a key only you hold |
| TH-0390 (S-0337) | 90.5 h | an `rm` under `.git` on the host | no: a host change |
| TH-0445 (S-0380) | 26.2 h | ten design choices; you overruled one | yes, with the overrule lost |
| TH-0436 (S-0378) | 24.9 h | a plan with two decisions | yes; a planner's thread is now the orchestrator's anyway (S-0328) |
| TH-0461 (S-0391) | 17.7 h | `flai host upgrade` on the host | no: a host change |

The orchestrator recommended within a minute every time. What took a day or more was you seeing it: the recommendation sat in the dashboard's inbox, nothing reached you off it, and the story sat in progress with no mark that it waited on you. S-0502 (I-0195) builds the timeout confirmation you asked for on TH-0598, and its exclusions leave exactly the three longest waits here untouched. So S-0498 does the two things that help those:

1. **A reminder that reaches you off the dashboard.** `flai serve` reminds you of a thread that awaits you (the inbox's `awaiting: operator` rule, ADR-0181) once its last entry is older than `threads.remind_after`, and again every `threads.remind_every`, through a channel kept in the host configuration: `flai serve notify set --url <webhook>` and/or `--command <shell command>` (the command gets the reminder's JSON on standard input, so an `ntfy` push or a Windows toast from WSL2 is one line). The reminder is written nowhere in the repository: a thread entry would end the agent's wait (ADR-0091) and change the thread's status (ADR-0090). Its state lives in the serve folder as `upgrade-asks.json` does.
2. **The idle is marked.** When a story's agent ends its session on an open question to you (ADR-0154), `flai serve` blocks the story with the reason `awaiting the operator on TH-nnnn`, and unblocks it when your answer starts the agent again. The board then shows what waits on you, and `flai stats` counts the wait as blocked time; the analysis found no `blocked_seconds` on any of the five stories although each was blocked in fact. The story stays in its column and still counts to the in-progress limit (workflow.md § Blocking).

## Tasks and layers

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-2127 | the ADR, refining ADR-0090, ADR-0091, ADR-0108, ADR-0154, and ADR-0158 | none |
| 2 | T-2128 | `threads.remind_after` and `threads.remind_every` in the manifest, check, `flai manifest set`, settings index, template | T-2127 |
| 2 | T-2129 | `flai serve` blocks an asked story and unblocks it on the answer (`agents.go`) | T-2127 |
| 3 | T-2130 | `flai serve notify set/show/clear` in the host config, and the sender (`notify.go`) | T-2128 (both add rows to `docs/operators/settings.md`) |
| 4 | T-2131 | the reminder loop (`remind.go`, wired in `serve.go`), its state file, its tests | T-2128, T-2130 |
| 4 | T-2132 | operator and user docs, workflow.md, flai-cli.md, strategic-agents.md, the work-management convention and its template copy | T-2129, T-2130 |
| 5 | T-2133 | this project's values (`4h`, `24h`), a thread giving you the `flai serve notify set` command for this host, close I-0151, tick the criteria | T-2131, T-2132 |

Every touch names a file. The ADR's path is predicted as `design/adrs/0204-...`; `flai task done` widens it to the real name.

## Assumptions

1. **Defaults.** `threads.remind_after` unset in the template (no reminders until a project opts in); `threads.remind_every` 24h. Here, T-2133 sets `4h` and `24h`. Say other values if you want them.
2. **The channel is yours to set**, in a shell, since it lives in `~/.flai/config.json`. T-2133 gives you the command on a thread; until you run it, `flai serve` logs each reminder and warns once that no channel is set.
3. **The webhook body extends the one `flaiover` posts to `dashboard.notify_url`**, so one receiver can take both. The dashboard's notifier is left as it is: it runs only while the container runs, posts once, and cannot run a host command, which is why `flai serve` is the sender.
4. **Blocking an asked story changes who writes the block**: today the convention asks the agent to run `flai block`, and none did on these five. flai writing it is a decision, so it is in the ADR, and the convention's rule gains one sentence (T-2132).
5. **Not in this story**: counting a blocked story outside the WIP limit (an ADR of its own, and the lanes were not full in W42), a reminder as a thread entry (ends the wait), and sorting the dashboard's inbox by age (small; say so if you want it added).

Nothing to split, merge, or drop. If you would rather not have remedy 2, drop T-2129 and I take its sentences out of T-2127 and T-2132; the reminder stands on its own.

### 2026-10-26T04:17:32Z orchestrator
Approved. The plan fits S-0498.

- **The remedy is grounded in the instances.** Three of I-0151's five waits needed your hands: a key, an `rm` under `.git`, and a host upgrade. No timeout confirmation from S-0502 can end those, so a reminder off the dashboard and a visible block are what reach them. T-2127's ADR records both before any code, as the goal asks. The reminder writes nothing on the thread, which keeps ADR-0090 and ADR-0091 intact.
- **The tasks fit.** There are seven tasks in five layers. T-2130 waits for T-2128 because both add rows to `docs/operators/settings.md`. The other tasks in each layer share no path.
- **The touches fit.** There are 34 file touches and no folder. The ADR file is predicted as 0204 and will be widened at close.
- **The figures fit.** The forecast is 1h45m, raised from flai's 54m with the reason stated. The value is 17550 USD/week from `flai cod` on the 117h input flai carried over from I-0151. I set no input.
- **The assumptions stand as the story's.** The template leaves `threads.remind_after` unset, this project sets 4h and 24h, the channel is set by you in a shell (T-2133 asks for it on a thread), and remedy 2 is kept. The story's agent changes the tasks if T-2127's ADR or your answer on T-2133's thread says otherwise.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T04:17:37Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0498.
