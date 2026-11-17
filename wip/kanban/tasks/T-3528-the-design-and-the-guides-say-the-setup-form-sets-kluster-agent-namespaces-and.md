---
id: T-3528
type: task
nature: remediation
title: The design and the guides say the setup form sets kluster.agent_namespaces, and I-0356 is closed
status: backlog
parent: S-0830
owner: alex
created: 2026-11-17T08:29:07Z
updated: 2026-11-17T08:29:07Z
transitions: []
stream: S-0830
tags: [docs, kluster]
touches: [design/system/flaiover-dashboard.md, design/system/kluster.md, docs/operators/settings.md, docs/users/flaiover.md]
after: [T-3526]
---
# T-3528 The design and the guides say the setup form sets kluster.agent_namespaces, and I-0356 is closed

## Work

It waits for T-3526, because it describes the field as T-3526 builds it: its label, its input, and how it saves.

Four documents say the form leaves the setting out, or list the form's fields without it:

- `design/system/kluster.md`, the paragraph on `kluster.agent_namespaces` under the limit: drop "The dashboard's setup form does not show it yet" and say the setup form sets it.
- `docs/operators/settings.md`, the `kluster.agent_namespaces` row: drop "the dashboard's setup form does not show it yet" and name the setup form on the Klusters and Deployments pages beside `flai manifest set`.
- `design/system/flaiover-dashboard.md`, the `/ops/klusters` row's step 4: add `agent_namespaces`, a whole number, to the fields `DeploySettings` shows from `settings`.
- `docs/users/flaiover.md`, step 4, **Deploying stories**: add the limit to the fields the form holds, and say what 0 does.

Then close the issue: `flai issue close I-0356 --reason "S-0830: DeploySettings.svelte names kluster.agent_namespaces in PLAIN_KEYS and shows a number setting, so the setup form shows and saves it"`.

## Done when

- None of the four documents says the setup form does not show `kluster.agent_namespaces`.
- `flaiover-dashboard.md` and `flaiover.md` list it among the form's fields.
- I-0356 is closed with a reason naming what fixed it.
- `flai test` passes on the changed paths, the markdown lint and `design-routes` included.

## Notes
