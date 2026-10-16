---
id: TH-0460
title: "S-0391 plan: the host runs flai 1.48.0, older than the self-upgrade check, so it needs one flai host upgrade by you; then docs, verify, close I-0149"
anchor:
  path: wip/kanban/stories/S-0391-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md
  item: S-0391
status: resolved
participants: [planner-S-0391, orchestrator, agent-S-0391]
created: 2026-10-15T05:06:10Z
updated: 2026-10-15T06:23:34Z
---

# TH-0460 S-0391 plan: the host runs flai 1.48.0, older than the self-upgrade check, so it needs one flai host upgrade by you; then docs, verify, close I-0149

On wip/kanban/stories/S-0391-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md.

## Entries

### 2026-10-15T05:06:10Z planner-S-0391
**Recommendation:** run `flai host upgrade` on this host now. Then turn on `flai serve enable auto-upgrade` in a shell, so the host never lags a release again. Please answer here: did you upgrade, and is `auto-upgrade` on or off?

## Why I-0149 happens

- The host's flai is 1.48.0 (`~/.flai/bin/flai version`). `host_actions` in `~/.flai/config.json` has no `auto-upgrade`.
- ADR-0158's minute check, which upgrades the host or asks you on a thread, shipped in flai 1.50.1 (S-0387).
- ADR-0160's close through the installed flai shipped in 1.51.0 (S-0390).
- A host older than 1.50.1 runs neither of them. It never notices a release and never asks you. No new code can reach a binary that is already installed.
- 1.51.2 is published. The 1.52.0 publish commit is on main, made at 05:02Z.

So the remedy is your upgrade, made once. From 1.50.1 on, ADR-0158 keeps the host current, or asks you on a thread for each release.

## Proposed tasks

S-0391 moved to in-progress at 05:04Z, while I was planning it. I created no tasks, so as not to duplicate the story agent's. agent-S-0391 can take these or plan its own:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | Docs: a host older than flai 1.50.1 does not check for releases and is upgraded by hand once. Bring update.md steps 5 and 6 in line with ADR-0160. | `docs/operators/index.md`, `docs/operators/runbooks/update.md`, `design/system/flai-cli.md` | none |
| 1 | Verify the host has caught up: `flai host status` and `flai version` at or above the newest tag, no `flai_outdated` in the inbox, and the MCP `task_done` closing a task with no refusal. Record the evidence in the story's Notes. | none in the repository | none (it waits on your upgrade) |
| 2 | Close I-0149 with `flai issue close --reason`, naming the upgrade and the docs. | the I-0149 issue file, `design/issues/summary.md` | both layer-1 tasks |

## Assumptions

- No code changes. `flai/internal/serve/upgrade_test.go` already reproduces the check on a newer host. A test of a pre-1.50.1 binary does not fit, which criterion 1 allows ("where one fits").
- Upgrading the host is yours. Agents never install or upgrade tools on your machine (`safety.md`).
- Story touches: the five files above, all files and no folders. Topics `release`, `host`, and `mcp` were added.

## Figures

- **Forecast 20m:** flai's 11m, raised for verifying through MCP. Delivery is 05:30Z, assuming a prompt upgrade.
- **Cost of delay 5 USD a week:** flai cod's figure, kept. The inputs are yours, and they likely understate the cost: while the host lags, every story's agent is refused at each task close.

The reasons are under `### Planning` in the story's Notes.

## Proposals

I propose no split, merge, or drop.

### 2026-10-15T05:06:22Z orchestrator (recommendation)
**Recommendation:** run `flai host upgrade` on the host now. Then decide on `auto-upgrade`:

- **On:** `flai serve enable auto-upgrade`, in a shell. The host then follows each published release by itself.
- **Off:** leave it. From 1.50.1 on, ADR-0158 asks you on one thread per release instead.

Answer here which you chose. S-0391's agent then verifies the host and closes I-0149.

## Why this needs you

- **It changes your machine.** The upgrade, and turning on a host action that upgrades the host by itself, are yours. Agents never install or upgrade tools on your machine (`safety.md`).
- **No code can reach the old binary.** Every inbox today reports "this flai is 1.48.0, older than flai 1.52.0". A 1.48.0 host has neither ADR-0158's minute check nor ADR-0160's close through the installed flai, so each story's `task_done` is refused until you upgrade.

## The plan fits

- Three tasks: docs, verify, and close I-0149. The first two run together, and the close comes after both.
- Every touch is a file, and no code changes. The criterion's "a test where one fits" is met by the existing `upgrade_test.go`.
- The figures are flai's. Forecast 20m, raised from 11m with a stated reason. Cost of delay 5 USD a week, which the planner notes likely understates the cost.

S-0391 is already in progress, so its agent takes these tasks or writes its own.

Source: design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md § Decision

### 2026-10-15T05:09:08Z planner-S-0391
agent-S-0391 has taken the plan over on TH-0461. It asks the same question: run `flai host upgrade`, then choose on `auto-upgrade`. Its tasks go further than mine:

- It adds `docs/operators/runbooks/migrate.md` to the docs.
- It adds a code change, with a new ADR refining ADR-0158: the upgrade thread names the newest *published* release, not the newest tag.

The story's own tasks will widen its touches to cover those paths. My three proposed tasks are superseded wherever TH-0461 differs. One answer from you, on either thread, settles both.

### 2026-10-15T05:22:12Z agent-S-0391
Agreed on the cause and on the upgrade being yours. I took the plan over on TH-0461; answer there or here, either works.

- T-1599 is done: the ADR-0158 thread now names a release only once it is published (ADR-0161, with tests).
- T-1600 is done: `update.md`, `migrate.md`, and `index.md` now cover a host older than 1.50.1. They also say which sessions ADR-0160 helps, which covers your point on steps 5 and 6.
- T-1601 waits for your `flai host upgrade`. Then I verify the host, close I-0149, and close the story out.

### 2026-10-15T06:23:34Z orchestrator
Resolved: The planner's plan is superseded by S-0391's own tasks (T-1599 and T-1600 done, T-1601 waiting). The one open ask, the operator's flai host upgrade and the auto-upgrade choice, continues on TH-0461 with my recommendation.
