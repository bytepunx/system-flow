---
id: TH-0843
title: "S-0567 live planner run: set the container login and the planner's harness, then plan S-0559 from the dashboard"
anchor:
  path: wip/kanban/stories/S-0567-flai-serve-starts-the-planner-in-a-container-on-harness-container-and-it-plans.md
  item: S-0567
status: resolved
participants: [agent-S-0567, orchestrator, alex]
created: 2026-11-07T03:49:04Z
updated: 2026-11-07T06:08:55Z
---

# TH-0843 S-0567 live planner run: set the container login and the planner's harness, then plan S-0559 from the dashboard

On wip/kanban/stories/S-0567-flai-serve-starts-the-planner-in-a-container-on-harness-container-and-it-plans.md.

## Entries

### 2026-11-07T03:49:04Z agent-S-0567
The preconditions an agent can check are in place. Three steps are yours, because the guard refuses an agent host-configuration writes (ADR-0220) and the manifest is yours to set.

### What I checked

S-0567's Notes, `### Live run`, hold the detail.

| Check | Found |
|-------|-------|
| Installed flai | 1.83.0. It holds S-0509 (released from 1.81.0). |
| Image | The agent images were built 2026-11-02T18:51Z, before S-0509's pilot change. I rebuilt them with `make agent-images` from main at `e41d6124e`. `flai-agent-sf-claude-code:local` is now `sha256:521d859c7bca`, with Claude Code 2.1.312 and flai 1.84.0. |
| Container harness | program `docker`, default bridge; login not set |
| Planner's agent | `planning.agent` names a model and no harness |
| Story to plan | S-0559, the first story `flai plan --candidates` lists with no plan |

### What I ask you to do, in this order

1. Set the login folder the container mounts read-only. Your claude-code harness runs `rmhclod`, so your login may live in `~/.claude-rmh` rather than `~/.claude`. Use the folder that holds the login `rmhclod` uses:

   ```sh
   flai serve agent harness container --login ~/.claude-rmh
   ```

2. Put the planner on the container harness:

   ```sh
   flai manifest set planning.agent='{"harness":"container","model":"claude-opus-5-5"}'
   ```

3. From the dashboard, plan S-0559. Do not use `flai plan` in a shell: that starts the planner as you, not through `flai serve`.

### After the run

I recommend you set `planning.agent` back to the host harness at once:

```sh
flai manifest set planning.agent='{"model":"claude-opus-5-5"}'
```

While it names `container`, every planner run on this project starts in a container, including replans the `plan` host action starts on its own. Several planners are running now.

### Reply here

Reply once the plan is asked for, or name another story to plan. I then check the run against criteria 1 to 3 (T-2502) and record what it showed (T-2503).

### 2026-11-07T03:54:28Z agent-S-0567
A correction to the last line above: T-2502 and T-2503 are cancelled as duplicates (TH-0838). Once you reply, T-2499 and T-2500 check the run against criteria 1 to 3, and T-2501 records what it showed. The three steps you are asked for are unchanged.

### 2026-11-07T03:55:55Z orchestrator (recommendation)
**Recommended:** do agent-S-0567's three steps in its order, then set `planning.agent` back to the host harness as soon as the planner's run on S-0559 ends.

1. `flai serve agent harness container --login <folder>`. Use the folder that holds the `.credentials.json` your login uses. ADR-0247 mounts that folder read-only at `/run/claude-login`, and only `.credentials.json` is read. This is one host setting, and S-0565 needs it too. TH-0836 step 3 names `~/.claude` and this thread names `~/.claude-rmh`. Pick the one that holds the credentials, and it serves both proofs.
2. `flai manifest set planning.agent='{"harness":"container","model":"claude-opus-5-5"}'`.
3. Plan S-0559 from the dashboard.

**Why only you:** each step is a credential path, host configuration, or a manifest key, and the guard refuses all of them to agents (ADR-0220).

**What I do meanwhile:**
- I start no planner with `plan` while `planning.agent` names `container`, so no planner of mine lands in the container beside the proof run.
- I skip S-0559 in `flai plan --candidates` while this thread is open, so the proof run is not used up on the host harness.

For the operator alone: this is a host change.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0247-a-strategic-agent-s-container-mounts-the-main-checkout-read-write-at-its-host.md

### 2026-11-07T05:38:10Z alex
Confirmed the recommendation of 2026-11-07T03:55:55Z orchestrator.

Source: design/adrs/0247-a-strategic-agent-s-container-mounts-the-main-checkout-read-write-at-its-host.md

### 2026-11-07T05:44:32Z agent-S-0567
The live run worked. planner-S-0559 ran on the container harness under pilot, with Claude Code 2.1.312 from the image where the host has 2.1.319. It planned S-0559 through flai alone (T-2539, T-2540, TH-0847), and `wip/agents/planner.md` gained its entry, 239 s and 2.1455 USD from the container's log. Inside the image, `flai guard` refuses the planner what the host refuses it, word for word. Criteria 1 to 3 are ticked. S-0567's Notes, `### Live run`, hold the detail, and I-0274 to I-0276 record what the run could not show.

One step of yours is still open. `system-flow.yaml` still names `harness: container` under `planning.agent`, so every planner run here starts in a container. As recommended, set it back:

```sh
flai manifest set planning.agent='{"model":"claude-opus-5-5"}'
```

I go on with T-2501 now: the design section, then the close-out.

### 2026-11-07T05:45:34Z orchestrator (recommendation)
**Recommended, alex:** set the planner back to the host harness now, as your confirmed answer says (D-0493):

```sh
flai manifest set planning.agent='{"model":"claude-opus-5-5"}'
```

**Why now:** planner-S-0559's container run has ended, and agent-S-0567 has ticked criteria 1 to 3. `system-flow.yaml` still has `harness: container` under `planning.agent`. So every planner run starts in a container, including the replans the `plan` host action starts by itself.

**What I do meanwhile:** while `planning.agent` names `container`, I start no planner with `plan`, so none of mine runs in the container. You do not need to answer my hold on S-0559 any more: S-0559 is planned (T-2539, T-2540, TH-0847).

**Why only you:** this is a manifest key and a host setting, and the guard refuses both to agents (ADR-0220).

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0493-s-0567-live-planner-run-set-the-container-login-and-the-planner-s-harness-then.md

### 2026-11-07T06:08:55Z alex
Resolved: S-0567 was accepted
