---
id: T-2564
type: task
nature: improvement
title: The Deployments page and the Features toggle say a block's state without naming a command, and a refused toggle shows its real state at once
status: done
parent: S-0586
owner: alex
created: 2026-11-07T08:31:50Z
updated: 2026-11-07T09:32:23Z
transitions:
  - to: ready
    at: 2026-11-07T09:21:30Z
    by: agent-S-0586
  - to: in-progress
    at: 2026-11-07T09:21:31Z
    by: agent-S-0586
  - to: done
    at: 2026-11-07T09:32:23Z
    by: agent-S-0586
stream: S-0586
tags: [kluster, flaiover]
touches: [flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/routes/api/kluster/config/+server.ts, flaiover/src/routes/api/kluster/config/config.test.ts]
after: [T-2562]
usage:
  source: log
  seconds: 652
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 22160
      cache_read: 8058371
      cache_write: 135751
      cost: 2.8388
---
# T-2564 The Deployments page and the Features toggle say a block's state without naming a command, and a refused toggle shows its real state at once

## Work

Waits for T-2562, whose `deployable` field the Deployments page reads. It shares no path with T-2563, so the two can run together.

- **Deployments.** In `flaiover/src/routes/ops/deployments/+page.svelte`, the notices for the host action off and for no block link the Klusters page's setup in place of naming `flai serve enable kluster` or the block to add. With no `deploy` set, **Deploy** and **Promote** are replaced by a sentence saying deploys are not set up yet, linking the same setup, rather than offering a button flai refuses.
- **Features.** The operator's report on S-0586 says the toggle looked on until a refresh, after flai refused it. In `flaiover/src/lib/components/SettingsPanel.svelte`, a toggle flai refuses goes back to the state flai holds at once, with flai's sentence beside that row, not at the bottom of the page. Confirm the fault first with a test that fails before the change. This applies to every host action, not to `kluster` alone.

Leave the Deployments page's port-forward commands until the plan thread settles them.

## Done when

- `flaiover/src/routes/ops/deployments/deployments.svelte.test.ts` covers the three notices with no `flai` command in their text.
- `flaiover/src/lib/components/SettingsPanel.svelte.test.ts` has a test that failed before the change: a refused toggle shows its real state and flai's sentence beside its row, with no refresh.
- `flai test` passes on the paths this task changes.

## Notes
