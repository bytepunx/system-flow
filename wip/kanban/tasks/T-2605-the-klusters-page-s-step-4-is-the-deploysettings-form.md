---
id: T-2605
type: task
nature: feature
title: The Klusters page's step 4 is the DeploySettings form
status: ready
parent: S-0595
owner: alex
created: 2026-11-07T22:08:21Z
updated: 2026-11-07T22:11:36Z
transitions:
  - to: ready
    at: 2026-11-07T22:10:19Z
    by: agent-S-0595
stream: S-0595
tags: [dashboard]
touches: [flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts]
after: [T-2603]
---
# T-2605 The Klusters page's step 4 is the DeploySettings form

## Work

ADR-0263 §5. Replace the Klusters page's step 4, "Deploying stories", with `DeploySettings`, removing the code it now holds (`DEPLOY_KEYS`, the deploy fields, `writeWithAgent` if nothing else uses it), so the setup lives in one form. Its Done or To do still follows `deployable`. Add a link from the step to Deployments' Setup pane. Waits for T-2603, the component. It touches no path T-2604 touches, so the two run together.

## Done when

- `klusters.svelte.test.ts` passes with its step 4 tests moved to what the page still owns (the step's state and the component's presence); the deploy fields' behaviour is tested in `DeploySettings.svelte.test.ts`.
- `flai test` on the changed paths passes.

## Notes
