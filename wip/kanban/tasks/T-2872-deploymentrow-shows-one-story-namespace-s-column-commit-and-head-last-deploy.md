---
id: T-2872
type: task
nature: feature
title: DeploymentRow shows one story namespace's column, commit and head, last deploy and checks run, and who asked, with its actions and its examples
status: done
parent: S-0663
owner: alex
created: 2026-11-10T08:25:31Z
updated: 2026-11-13T15:00:19Z
transitions:
  - to: ready
    at: 2026-11-13T14:51:14Z
    by: agent-S-0663
  - to: in-progress
    at: 2026-11-13T14:51:15Z
    by: agent-S-0663
  - to: done
    at: 2026-11-13T15:00:19Z
    by: agent-S-0663
stream: S-0663
tags: [dashboard]
touches: [flaiover/src/lib/components/DeploymentRow.svelte, flaiover/src/lib/components/DeploymentRow.svelte.test.ts, design/components/flaiover/src/lib/components/DeploymentRow.examples.json]
after: [T-2871]
usage:
  source: log
  seconds: 544
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 84
      output: 427
      cache_read: 4342069
      cache_write: 132347
      cost: 1.6743
---
# T-2872 DeploymentRow shows one story namespace's column, commit and head, last deploy and checks run, and who asked, with its actions and its examples

## Work

Draw one story namespace's row out of the Overview pane's `deployed` snippet in `flaiover/src/routes/ops/deployments/+page.svelte` into a component, `flaiover/src/lib/components/DeploymentRow.svelte`, so it can be tested and shown on the Components page with examples.

- It shows the story's ID and title, linking its item page, and its column; the namespace; the commit, and whether it is the head, with the head when it differs; the last deploy and the last checks run, each with its outcome, when, who asked, and a link that opens Logs at its step; and a badge per step as today.
- Its actions are props the page handles: **Deploy** for a story in progress or in review deployed at a commit other than its head, or whose deploy failed or was interrupted; **Run checks** when the namespace is deployed at the head; **Teardown**, after confirmation. With no `deploy` set, **Deploy** gives way to **Set up deploys**, as the page has it. While a job runs the actions wait.
- Use the design system's tokens and the components the page already uses. Its props take the types T-2871 added to `repo.ts`.
- Write its examples, `design/components/flaiover/src/lib/components/DeploymentRow.examples.json`, as `design/components/README.md` says: a story in progress deployed at its head and checked, one behind its head, one whose deploy failed, one in review, and one asked for by its agent.
- Test it in `DeploymentRow.svelte.test.ts`.

It waits for T-2871, whose types its props take.

## Done when

- The component's tests pass for each example's case and each action.
- `flai check` reports no `components.examples` finding for it.
- `flai test` passes on the changed paths, `vitest` and `svelte-check` among them.

## Notes
