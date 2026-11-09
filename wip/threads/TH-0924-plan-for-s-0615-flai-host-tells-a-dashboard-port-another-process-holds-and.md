---
id: TH-0924
title: "Plan for S-0615: flai host tells a dashboard port another process holds and stops restarting it"
anchor:
  path: wip/kanban/stories/S-0615-flai-host-cannot-see-a-dashboard-container-docker-lost-track-of-which-keeps-the.md
  item: S-0615
status: open
participants: [planner-S-0615, orchestrator]
created: 2026-11-09T09:54:23Z
updated: 2026-11-09T09:55:05Z
---

# TH-0924 Plan for S-0615: flai host tells a dashboard port another process holds and stops restarting it

On wip/kanban/stories/S-0615-flai-host-cannot-see-a-dashboard-container-docker-lost-track-of-which-keeps-the.md.

## Entries

### 2026-11-09T09:54:23Z planner-S-0615
S-0615 is planned: six tasks in four layers, its touches, a forecast of 2h, and a cost of delay value of 21600 USD a week. The proposed fix and the reasons are in the story's Notes under `### Planning`.

## Proposed fix

- The cause in the code: `dashboardState` in `flai/cmd/dashboard_watch.go` answers not answering whenever `docker port` gives no address. So a container whose port map is empty is restarted at every back-off, for ever.
- A new watch state, `port-held`: docker runs the recorded container, it publishes nothing, and the recorded host port answers anyway. The watch does not restart on it. It logs it once each time the state is entered, naming the holder.
- The holder is the `docker-proxy` whose `/proc/<pid>/cmdline` carries `-host-port <port>`, which any user can read. Where none is found, flai says only that another process holds the port.
- `flai dashboard status` and `flai host status` say what holds the port and how to free it.
- A new ADR refines ADR-0062.

## Tasks

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-2714 ADR and `design/system/flaiover-dashboard.md` § Runtime | nothing |
| 1 | T-2715 holder lookup and port probe, `flai/cmd/dashboard_porthold.go` | nothing |
| 2 | T-2716 the `port-held` state in `flai/internal/host/watch.go` and `flai/cmd/dashboard_watch.go`, with the test that reproduces I-0243 | T-2715 |
| 3 | T-2717 status lines in `flai/cmd/dashboard.go` and `flai/cmd/host.go` | T-2716 |
| 3 | T-2719 `flai issue close I-0243` | T-2714, T-2716 |
| 4 | T-2718 `docs/operators/index.md` and `docs/users/flai.md` | T-2714, T-2717 |

## Assumptions

- The watch should not restart while another process holds the port, and flai should not kill that process. The holder runs as root, and killing it is the operator's call.
- The holder lookup is for Linux only. On macOS and Windows flai says only that the port is held.
- The story keeps one folder touch, `design/adrs`, because the ADR's number is not known yet. The folder is in `claims.shared`, so it holds no ready story.
- The story has no `after`.

## Question: the cost of delay input

I recommend lowering `time_lost_per_cycle` from 144h to 4h. That would make the value 600 USD a week instead of 21600.

flai set the 144h from I-0243. It is how long the defect lasted, six days of restarts. It is not work lost: browsers reached the orphan container throughout. The work lost was the operator's time to find the orphan.

Alternatives:

- Keep 144h. The value stays at 21600 USD a week, which puts this story above most of the backlog.
- Name another figure.

The input is yours, so I have left it unchanged. The plan does not wait on your answer.

### 2026-11-09T09:55:05Z orchestrator (recommendation)
**Recommended:** lower S-0615's `time_lost_per_cycle` from 144h to 4h, as the planner proposes. flai then works the value out at 600 USD a week instead of 21600.

**Why:**

- 144h is how long I-0243 lasted, six days of restarts. It is not work lost: browsers reached the orphan container throughout.
- The work lost was your time to find the orphan and free the port. A few hours fits that.
- It changes no order today. The policy is `throughput`, so the ready column is ordered by forecast (2h). It does keep `flai stats` from counting this defect as the costliest item in the backlog.

**Why I recommend rather than set it:** S-0615 already has an input, the one flai set from I-0243. `plan_backlog_stories` lets me set only a backlog story's first input. Once you confirm or name another figure, I set it under `set_cost_of_delay`, the input only.

**Alternatives:** keep 144h (21600 USD a week), or name another figure.

**The plan:** it fits the story. Six tasks in four layers, file touches but the shared `design/adrs`, a 2h forecast, and no `after`. Not killing the docker-proxy is right: it runs as root, and freeing it is yours. I am finalizing the draft now, since `flai promote --drafts` lists it complete and the policy does not use its value.

Source: design/system/strategic-agents.md § Cost of delay

Source: design/system/strategic-agents.md § Cost of delay
