---
id: TH-0202
title: "S-0298 plan: nine tasks in five layers, and two questions on scope (MCP deploy, dashboard pin)"
anchor:
  path: wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md
  item: S-0298
status: resolved
participants: [planner-S-0298, alex]
created: 2026-10-06T21:47:10Z
updated: 2026-10-06T22:52:52Z
---

# TH-0202 S-0298 plan: nine tasks in five layers, and two questions on scope (MCP deploy, dashboard pin)

On wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md.

## Entries

### 2026-10-06T21:47:10Z planner-S-0298
Plan for S-0298, by planner-S-0298. The touches, forecast (1h1m, delivery 2026-10-07T09:29Z) and cost of delay (37.50 USD a week, from your 15m per cycle) are on the story. Where each touch came from is under `### Planning` in its Notes.

Two questions for you. My recommended answer comes first in each.

1. **Over MCP, should agents only list versions, or also deploy one?** Recommended: list only. The new MCP tool `versions` (T-1041) reads the published flai and dashboard releases. Deploying stays with you, through the CLI and the dashboard's `host` and `dashboard` host actions. An upgrade from MCP would restart the host and the MCP server that took the call, and it would let an agent change the flai every other agent runs on. The alternative is an MCP `upgrade` tool gated by a host action. T-1041 would add it.
2. **When you pick a dashboard version on the Updates page, should it stay pinned, or apply once?** Recommended: once, like `flai dashboard upgrade --tag`. The page would say that the configured tag (`dashboard.tag`, today `latest`) is used again the next time the container starts, and how to pin with `flai config set dashboard.tag`. Pinning from the page would write the host's configuration, which is the `settings` host action's job (ADR-0039). The alternative: the page also offers "keep this version", which sets `dashboard.tag` behind the `settings` host action. T-1044 and T-1046 follow your answer either way.

Tasks, layers and order:

- Layer 1:
  - T-1038 an ADR, plus the design in `flai-cli.md` and `flaiover-dashboard.md`. It fixes the names and decides that a dashboard request may name a release only when it is published.
  - T-1039 `selfupgrade.List`: the published `flai/v` or `flaiover/v` releases, newest first, across pages.
- Layer 2:
  - T-1040 `flai self-upgrade --list`, `flai dashboard versions`, and `flai dashboard upgrade --published`. It waits for T-1038 and T-1039.
  - T-1041 the MCP tool `versions`. It waits for T-1038 and T-1039.
  - T-1042 the dashboard's `/api/host` and `/api/dashboard` routes, plus the channel method list. It waits for T-1038 only, because its tests stand in for the channel.
- Layer 3:
  - T-1043 `flai host versions`, and `flai host upgrade --version`, which restarts the host as the newest does. It waits for T-1040.
  - T-1044 the Versions control on the Updates page's Dashboard and flai host areas. It waits for T-1042.
- Layer 4: T-1045 the host API's reads `host.versions` and `dashboard.versions`, and `host.upgrade {version}` and `dashboard.upgrade {tag}`. It waits for T-1043.
- Layer 5: T-1046 the runbook, the guides, and the regenerated reference. It waits for T-1041, T-1044 and T-1045.

Assumptions:

- Listing reads GitHub releases, the same way `self-upgrade` already does. Dashboard image tags match `flaiover/vX.Y.Z` release versions, as the update runbook's `dashboard.tag 0.27.4` example shows.
- The dashboard never names an image, only a version. flai on the host accepts a version only when it is a semver and a published release, so a request from the page cannot install anything you have not released. Your own shell's `--tag` (a mirror, say) stays as it is today.
- "Works like updating to the latest" means the same paths: the same host action gating, the same health-checked swap for the dashboard, and the same host restart with its processes for flai.
- A flai below a served project's `flai.minimum` is marked in the list and warned about. It is not refused.
- The story keeps your declared folder touches `flai/cmd` and `flaiover/src`. The tasks name files inside them, so the story's claim narrows to those files while it is in progress (ADR-0096). The `design/adrs` folder touch stays because the ADR's number is not known yet.
- `flai/internal/guard/guard.go` is left out: `versions` is not added to `guard.MCPReads`. That keeps S-0298 clear of the guard file that S-0229 and S-0299 hold.

No task to split, merge or drop.

### 2026-10-06T21:51:45Z alex
1. list only, 2. apply once, no pinning

### 2026-10-06T21:56:45Z planner-S-0298
Thank you. Both answers are now in the tasks:

- T-1038: the ADR records that MCP lists only and that a dashboard version picked on the page applies once.
- T-1041: the tool `versions` reads and deploys nothing.
- T-1044: the page says a dashboard version picked there lasts until the next start, and offers no pin.
- T-1046: the docs say the same.

The tasks, layers and figures are otherwise unchanged.

### 2026-10-06T22:52:52Z alex
Resolved.
