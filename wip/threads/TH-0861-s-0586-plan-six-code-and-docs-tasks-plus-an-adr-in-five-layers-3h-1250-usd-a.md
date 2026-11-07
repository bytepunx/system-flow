---
id: TH-0861
title: "S-0586 plan: six code and docs tasks plus an ADR in five layers, 3h, 1250 USD a week; three questions on what stays outside the dashboard"
anchor:
  path: wip/kanban/stories/S-0586-provide-first-time-wizard-for-kluster-setup.md
  item: S-0586
status: open
participants: [planner-S-0586, orchestrator, agent-S-0586]
created: 2026-11-07T08:33:47Z
updated: 2026-11-07T08:53:28Z
---

# TH-0861 S-0586 plan: six code and docs tasks plus an ADR in five layers, 3h, 1250 USD a week; three questions on what stays outside the dashboard

On wip/kanban/stories/S-0586-provide-first-time-wizard-for-kluster-setup.md.

## Entries

### 2026-11-07T08:33:47Z planner-S-0586
S-0586 is planned, except for three decisions that are yours. Each one decides whether criterion 3 holds in full. Tasks for questions 2 and 3 are not written yet.

## Questions

**1. Should the `settings` host action stay the one step turned on in a shell? I recommend yes.**
- The Features toggle already needs `settings`, and `settings` can only be turned on in a shell (ADR-0039).
- The setup page's writes would sit behind it too, because whoever sets the `build`, `deploy`, and `teardown` commands runs them on the host as you.
- With yes, `flai serve enable settings` stays the one shell step. A hosted flai would turn it on at provisioning.
- The alternative is a dashboard toggle for `settings`. That meets criterion 3 literally, but it reverses ADR-0039: the dashboard token then becomes equal to a shell on the host, with no consent given on the host.

**2. Should the dashboard install the host's tools? I recommend an Install button for `kluster` and `kubectl`, with Docker left as a stated prerequisite.**
- flai serve would install the two tools into a folder flai owns, behind the `kluster` host action. Docker needs root and a daemon, so the page would only report whether it is there.
- This adds one task and about 45 minutes. It also runs a download on your host from a button, so it is a host change only you can approve.
- The alternative is that the page reports what is missing and gives the install sentence. That is smaller, but it leaves a shell step.

**3. Should the dashboard run the port-forwards? I recommend Start and Stop buttons per service on the Deployments page.**
- Today that page's "Reach it for a demo" shows `kubectl port-forward` commands for you to run on the host.
- flai serve would run each forward, bound to `127.0.0.1` on the host, and the page would link to it. This adds one task and about 40 minutes.
- Limit: a browser on another machine still cannot reach a forward on the host's loopback. Reaching it from there would need a proxy through the dashboard, which I would make a separate story.
- The alternative is to leave the commands as they are, which leaves a shell step.

## Tasks and layers

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-2559 | ADR refining ADR-0254, and `design/system/kluster.md` |
| 2 | T-2560 | `flai serve enable kluster` writes and commits the block and config; `flai manifest set` writes the block's keys |
| 2 | T-2561 | flai serve's kluster pass does nothing while no `deploy` is set |
| 3 | T-2562 | `kluster.status` reports what setup lacks; new `kluster.workloads` read |
| 4 | T-2563 | The Klusters page's setup steps and management view |
| 4 | T-2564 | Deployments notices; a refused Features toggle shows its real state at once |
| 5 | T-2565 | Docs, listing any step left outside the dashboard |

## Assumptions

- **The scaffold:** the block names only `config: kluster.yaml`, and `kluster.yaml` at the repository root names the cluster and `k3d`, with no profile. This is your "empty cluster" option.
- **`deploy` becomes optional, against ADR-0254:** the other way to pass `flai check` is a fake `deploy: ["true"]`. That would bring up a cluster for every story and show "passed" with nothing deployed. With `deploy` absent, flai runs no job until you set it.
- **The project's scripts:** build, deploy, and teardown scripts are project code, so the page offers a button that makes a backlog story for an agent to write them, rather than a script editor.
- **Figures:** the forecast is 3h, not `flai forecast`'s 4m, based on S-0577, S-0579, S-0580, and S-0581. The cost of delay is 1250 USD a week, from your input.

## Not verified

- I could not read kluster's config schema. WebFetch was refused, and `kluster`, `curl`, and `gh` are not on the planner's host. So "a config with no profile gives an empty cluster" is unconfirmed. T-2559 checks it against the real command first, and asks on a thread if it can't.
- The start prompt said you wrote on a thread named `asked`. No such thread exists and S-0586 had no threads, so I planned from the story itself. If you did ask for something specific, please say it here.
- Writing touches clears the story's `no-touches` hold, so its agent can start once in-progress has room, possibly before these answers. T-2559 tells it to read this thread first.

### 2026-11-07T08:34:58Z orchestrator (recommendation)
**Recommended:** approve the plan, and answer the three questions as the planner proposes:

1. **Yes:** `settings` stays the one step turned on in a shell.
2. **Yes:** add an Install button for `kluster` and `kubectl`, behind the `kluster` host action. Docker stays a stated prerequisite.
3. **Yes:** add Start and Stop buttons for port-forwards on the Deployments page. A proxy for other machines becomes a separate story.

Once you answer, the planner writes the tasks for 2 and 3 and revises the 3h forecast.

## The plan: fits the story

- **Every criterion is covered.** T-2560 covers criterion 1. T-2562, T-2563, and T-2564 cover criterion 2. T-2565 lists any step left outside the dashboard, for criterion 3.
- **The layers follow real dependencies.** T-2559 sets the ADR first. T-2560 and T-2561 share no path.
- **Touches, forecast, and value fit.** The 41 touches cover every task's touches. The 3h forecast states its basis, S-0577 to S-0581. The value is 1250 USD a week, from your revenue input.
- **Making `deploy` optional is right.** A fake `deploy: ["true"]` would report "passed" with nothing deployed. It refines ADR-0254 through T-2559's new ADR, so nothing is reversed silently.

## Why each answer is yours

1. **Reversing ADR-0039 is your call.** You decided on 2026-09-23 that `settings` is turned on only in a shell. A dashboard toggle would make the dashboard token equal to a shell on the host. One shell step at setup does not defeat the story's goal, and a hosted flai would turn `settings` on at provisioning.
2. **Installing software on the host is yours alone.** A button that downloads and installs binaries changes the host. It adds one task and about 45 minutes. The smaller alternative, a page that reports what is missing, leaves a shell step against criterion 3.
3. **This is a scope addition.** It adds one task and about 40 minutes. It removes the last everyday shell command from the Deployments page. The forwards bind to `127.0.0.1` only, so nothing new is exposed.

## Note

S-0586 is already in ready, so its agent may start before you answer. T-2559 tells it to read this thread first.

Source: design/adrs/0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md § Decision

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md § Decision

### 2026-11-07T08:37:51Z agent-S-0586 (recommendation)
S-0586 is started. I checked kluster against its source, and it changes one assumption of the plan. I am building the parts no answer changes now.

## Finding: kluster v1.3.1 has no empty cluster

- A config with no `profile` gets kluster's default profile, `spire`: cert-manager, SPIRE, and traefik-tls (`kluster/cmd/up.go`, `stringFlag` in `kluster/cmd/config.go`).
- An empty `profile` falls back to the same default, and an unknown one is refused (`resolveProfiles` in `kluster-lib/cluster/cluster.go`).
- So the planner's "a config with no profile is an empty cluster" is false. I did not run `kluster up`: the source settles it, and a bring-up changes your host.

The block flai writes is still a no-op. It names no `deploy`, so flai runs no job until you set one.

## Question 4: which profile should the scaffold name

I recommend `spire`, written out.

- `spire` is kluster's default and its smallest documented profile.
- The Klusters page would let you change the profile (`spire`, `signet`, `authstar`) and the addons before the first Up.
- The alternative is a lighter group such as `tracing`. kluster accepts it, but it is undocumented as a base profile.

## My recommendation on all four

1. Yes: `settings` stays the one shell step (ADR-0039).
2. Report missing tools, and give the install sentence, with no Install button. `kluster setup` installs `kubectl` with `sudo` on Linux, so a button in `flai serve` would hang or fail without passwordless sudo. Docker cannot be installed by it at all.
3. Yes: Start and Stop for each port-forward on the Deployments page, bound to `127.0.0.1`.
4. `profile: spire`, changeable on the Klusters page.

Confirm this to make it the answer, or answer otherwise per number. Until then I build the block, the optional `deploy`, the serve pass, and the status read. Those need no answer.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md

### 2026-11-07T08:40:05Z orchestrator (recommendation)
**Recommended:** take agent-S-0586's four answers as they stand. I withdraw my own answer 2.

1. **Yes:** `settings` stays the one shell step (ADR-0039).
2. **No Install button.** The page reports what is missing and gives the install sentence. kluster's own install puts its binary in `/usr/local/bin`, and `kluster setup` installs `kubectl` with `sudo` on Linux. A button in `flai serve` would hang or fail without passwordless sudo, and Docker cannot be installed that way at all.
3. **Yes:** Start and Stop for each port-forward on the Deployments page, bound to `127.0.0.1`.
4. **`profile: spire`**, written out in the scaffold and changeable on the Klusters page. kluster v1.3.1 has no empty cluster, and `spire` is its default.

**What this means for criterion 3:** with answer 2, two steps stay outside the dashboard: turning on `settings`, and installing Docker, `kluster`, and `kubectl` on the host. T-2565 lists both in the docs. Both are host changes, so they are yours to accept.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/tech/kluster.md

### 2026-11-07T08:53:28Z planner-S-0586
I'm handing this thread to agent-S-0586 and ending the planner's run. S-0586 is in progress, and its agent now owns the plan (work-management.md). If I also wrote the tasks for your answers to 3 and 4, they could duplicate the ones it writes.

Two corrections to my first entry:

- **Question 2:** I withdraw the Install button. agent-S-0586 found that kluster installs to `/usr/local/bin`, and that `kluster setup` uses `sudo` for `kubectl`. A button run by flai serve would hang or fail.
- **The scaffold:** my "no profile means an empty cluster" was wrong. kluster v1.3.1 falls back to `spire`. T-2559 still says otherwise, and agent-S-0586 is correcting that as it goes.

Once you answer, agent-S-0586 does three things:

- writes the port-forward task (answer 3) and the profile step (answer 4);
- adds both to T-2565's `after`;
- revises the 3h forecast. Each new task adds about 40 minutes.
