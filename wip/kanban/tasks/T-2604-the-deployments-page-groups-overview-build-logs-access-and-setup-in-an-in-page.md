---
id: T-2604
type: task
nature: feature
title: The Deployments page groups Overview, Build, Logs, Access, and Setup in an in-page rail, with a panel per component build
status: ready
parent: S-0595
owner: alex
created: 2026-11-07T22:08:12Z
updated: 2026-11-07T22:11:22Z
transitions:
  - to: ready
    at: 2026-11-07T22:10:15Z
    by: agent-S-0595
stream: S-0595
tags: [dashboard]
touches: [flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts]
after: [T-2602, T-2603]
---
# T-2604 The Deployments page groups Overview, Build, Logs, Access, and Setup in an in-page rail, with a panel per component build

## Work

ADR-0263 §5. Rebuild `flaiover/src/routes/ops/deployments/+page.svelte` around `PageSidebar` (as `routes/items/[id]/+page.svelte` uses it, hash-kept panes) with five panes:

- **Overview**: the status line (host action, cluster, deploys set up or not, the running job), main's namespace and each story's, and the stories in review not deployed, with their actions (Deploy, Teardown after confirmation, Promote), as the page has them now. With no `deploy` (`deployable` false), Deploy and Promote are replaced by a link that opens the Setup pane, never a call that flai refuses.
- **Build**: a chooser of the namespace (the running job's by default, else main's, else the first), then one `LogFeed` panel per build step of its record (`build`, or `build-<name>` per component), each with its status badge and expandable live feed; the load step after them; a sentence when the block builds nothing.
- **Logs**: the same namespace's every step (build, load, deploy, teardown) as `LogFeed` panels, and the running up/down/deploy job's own log tail.
- **Access**: the port-forwards per namespace, as **Reach it for a demo** has them now, with the one warning that a forward is reachable from the host's network.
- **Setup**: `DeploySettings` with `kluster.status`'s settings, saving back into the page's status; a link to the Klusters page for the cluster and its config file.

Counts on the rail entries where they help (deployments, running). Keep the 3-second poll of deployments while a job runs, and the notices for the action off, no block, and an old flai. UI from the design tokens only. Waits for T-2602 and T-2603, the components it uses.

## Done when

- `deployments.svelte.test.ts` covers the five panes, the Setup link in place of Deploy with no `deploy`, one panel per component build with its status, pending steps shown as pending, and the existing deploy, teardown, promote, and forward behaviour; it fails without the change.
- `flai test` on the changed paths passes.

## Notes
