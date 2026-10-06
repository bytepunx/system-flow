---
id: TH-0208
title: "S-0294: may flai serve restart a story's agent on its own, and how many times?"
anchor:
  path: wip/kanban/stories/S-0294-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md
  item: S-0294
status: open
participants: [planner-S-0294, alex]
created: 2026-10-06T22:51:08Z
updated: 2026-10-06T23:09:31Z
---

# TH-0208 S-0294: may flai serve restart a story's agent on its own, and how many times?

On wip/kanban/stories/S-0294-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md.

## Entries

### 2026-10-06T22:51:08Z planner-S-0294
S-0294 has one part left: the second remediation from I-0084. When a story's agent ends with the story still in progress, flai serve marks the run failed and does nothing more. ADR-0043 gave restarting to a person, so whether flai serve may restart on its own is your call.

**My recommendation: yes, up to 2 automatic restarts per story, then a thread to you.** In detail:

- **When it restarts.** A run ended, and all of these hold:
  - the story is still in-progress;
  - the story is not blocked;
  - no unresolved thread on the story or its tasks awaits you (the restart on an answer already covers that case);
  - the operator did not stop the agent (ADR-0058);
  - the `agent` host action is on.
- **How.** It restarts the same way `flai serve agent restart` does: a new session, the same name, and the run recorded in `serve/agents.json`.
- **The limit.** Set by a new manifest key, `agent.auto_restart`. It defaults to 2, and 0 turns automatic restarts off. The count resets when you restart the agent by hand or when the story leaves in-progress.
- **At the limit.** flai serve opens a thread on the story saying that the agent ended N+1 times and how to restart it. It then leaves the story as it does today.
- **The record.** An ADR that refines ADR-0043, and `flai-cli.md`'s `flai serve` section, would record this.

Reply `yes` to take the recommendation, or give your own terms: a different limit, no automatic restart (the story would then only add the thread at the first failure), or other conditions. I'm drafting the touches, forecast, and tasks in the meantime.

### 2026-10-06T22:53:07Z planner-S-0294
A correction to where the limit would live. The code shows that flai serve's agent settings are this host's, not the project's: `agent.command`, `agent.name`, and the harnesses sit in `~/.flai/config.json`, and `flai serve agent set` writes them. The `agent` host action that gates restarting is also the host's.

So I now recommend a host setting, **`agent.auto_restarts`**, rather than a manifest key. You would set it with `flai serve agent set --auto-restarts N`, and `docs/operators/settings.md` would list it. Everything else in my first entry stands.

Reply `yes` to take this, or say `manifest` if you would rather each project carry its own limit.

### 2026-10-06T23:00:25Z alex
yes

### 2026-10-06T23:08:56Z planner-S-0294
This is the plan for S-0294. It is drafted on my recommendation above, and I'll adjust it to your answer.

**The tasks, in three layers:**

- **Layer 1, with no prerequisites:**
  - T-1124: an ADR that refines ADR-0043. It records the automatic restart, its conditions, the limit, and the thread at the limit. `flai-cli.md` and `workflow.md` are brought into line.
  - T-1125: the host setting `agent.auto_restarts`, 2 when unset and 0 for off. It is set with `flai serve agent set --auto-restarts`.
- **Layer 2, after T-1124 and T-1125:**
  - T-1126: the launcher restarts a run that ended with its story in progress, up to the limit, and then opens one thread to you. Its tests in `agents_test.go` start with one that reproduces I-0084 and fails on main. Then one per condition: blocked, stopped, asked, limit 0, the reset on a manual restart, and the single thread.
- **Layer 3, after T-1126:**
  - T-1127: the operators' and users' guides and the reference. I-0084 is closed, naming S-0285 and the new ADR.

**The figures:**

- **Forecast.** 40m, against flai's 20m: an ADR, a setting, and a launcher change tested across six conditions. S-0285 took 42m.
- **Cost of delay.** 50 USD a week, flai's figure from the inputs, unchanged.

**Assumptions:**

- The answer is `yes`. With `manifest`, T-1125 moves to `flai/internal/manifest`. With "no automatic restart", T-1126 only opens the thread when the run fails, and T-1125 is dropped.
- A run the operator stopped, and one that ended asking a question, are never restarted automatically. The restart on an answer is unchanged.
- The count is per story. It resets on a manual restart and when the story leaves in-progress.
- The dashboard needs no change: the run's `why` already shows on the story page, and the Settings page does not list host agent settings.
- S-0272 also changes `flai/internal/serve/restart.go` and the serve docs, so whichever starts second is held until the other is in review.

### 2026-10-06T23:09:31Z planner-S-0294
Thank you. I've recorded your `yes`. T-1124, T-1125, T-1126 and T-1127 now carry the confirmed terms:

- the host setting `agent.auto_restarts`, 2 when unset and 0 for off;
- a restart only for a run that ended with its story in progress, unblocked, not stopped, and asking nothing;
- one thread to you at the limit.

S-0294's `### Planning` notes say the same. The layers and figures are as in my plan above. One correction to the plan: S-0272 overlaps on `restart.go` and the serve docs, not on `agents.go`.
