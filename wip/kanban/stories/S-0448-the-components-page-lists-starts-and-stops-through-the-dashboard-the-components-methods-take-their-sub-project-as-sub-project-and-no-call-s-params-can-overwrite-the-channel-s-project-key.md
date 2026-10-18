---
id: S-0448
type: story
nature: remediation
title: "The Components page lists, starts, and stops through the dashboard: the components methods take their sub-project as sub_project, and no call's params can overwrite the channel's project key"
status: review
owner: alex
created: 2026-10-18T02:25:46Z
updated: 2026-10-18T02:51:51Z
transitions:
  - to: ready
    at: 2026-10-18T02:26:04Z
    by: alex
  - to: in-progress
    at: 2026-10-18T02:26:13Z
    by: agent-S-0448
  - to: review
    at: 2026-10-18T02:51:51Z
    by: agent-S-0448
tags: [cli, dashboard]
topics: [cli, dashboard]
touches: [flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/components, design/system/component-library.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, flai/internal/hostapi/params_test.go, design/system/dashboard-host-channel.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1553
  priced_by: harness
  turns:
    - day: 2026-10-18
      hand_edits: 2
      work: 45
  models:
    - model: claude-opus-5-5
      input: 182
      output: 57604
      cache_read: 9236803
      cache_write: 322659
      cost: 5.1541
  strategic:
    - kind: orchestrator
      seconds: 147
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 48
          cache_read: 3619644
          cache_write: 5889
          cost: 0.9088
        - model: claude-sonnet-5-5
          input: 12
          output: 80
          cache_read: 152351
          cache_write: 57378
          cost: 0.1717
cost_of_delay:
  inputs:
    penalty_per_week: 1000
    by: alex
    at: 2026-10-18T02:25:46Z
---
# S-0448 The Components page lists, starts, and stops through the dashboard: the components methods take their sub-project as sub_project, and no call's params can overwrite the channel's project key

## Goal

Make the Components page work through the dashboard, which it has never done ([I-0163](../../../design/issues/I-0163-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md)). The dashboard's agent sends every call over the channel as `{ project: <project key>, ...params }` (`flaiover/src/lib/server/agent.ts`, `ask`), and flai's channel refuses a call whose `project` is not the key it serves (`flai/internal/channel/channel.go`, `call`). The four components methods read that same `project` field as the sub-project (`componentsParams` in `flai/internal/hostapi/components.go`). So `components.list` and `components.status`, sent with no sub-project, take the key (`sf` here) as the sub-project and refuse it: "system-flow.yaml has no sub-project "sf"; name one of flai, flaiover, template". `components.start` and `components.stop`, sent with `project: "flaiover"`, overwrite the key, and the channel refuses them as another project's. The Go tests call the handlers directly and the dashboard's tests use a fake source, so no test sent a call through the channel.

Reproduce in the main checkout:

```bash
flai hostapi components.status '{"project":"sf"}'
flai hostapi components.status '{}'
```

The first is refused, the second answers.

## Acceptance criteria

- [x] `components.list`, `components.status`, `components.start`, and `components.stop` take the sub-project as `sub_project`, and none reads `project` from its params. A call that carries the channel's `project` key and no `sub_project` lists, or reports, every supported sub-project; one with `sub_project` acts on that sub-project alone. The answers' own `project` field, the sub-project's name in `ComponentServer` and the catalogue, is unchanged.
- [x] The dashboard's `Repo` sends `sub_project` for `components.start` and `components.stop`, and the HTTP routes under `/api/components` keep their `project` body field, so the page and its URL are unchanged.
- [x] The dashboard's agent puts the channel's `project` key after a call's params, so that no method's params can overwrite it, and a test fails if a call's params carry a `project` of their own.
- [x] A Go test sends each of the four methods through the channel's `call` with the project key in its params, as the dashboard does, and fails on today's code: list and status answer every supported sub-project, start and stop with `sub_project` reach their handler. A test in flai fails if any channel method's params type has a field named `project`, so that no later method repeats this.
- [ ] With the fixed flai and flaiover running on this host, the Components page lists flaiover's components, starts its dev server, frames a component, and stops the server, and the story's notes say so with the versions used.
- [x] `design/system/component-library.md` and `design/system/flaiover-dashboard.md` name `sub_project` as the methods' parameter and say why `project` is the channel's alone, and `docs/users/flai-reference.md` is regenerated if it names the parameter.

## Tasks
- T-1896 The components methods take their sub-project as sub_project, and Go tests send them through the channel with the project key and fail on any params type with a project field
- T-1897 The dashboard's agent puts the channel's project key after a call's params, and Repo sends sub_project for components.start and components.stop
- T-1898 The design names sub_project as the components methods' parameter and says the project key is the channel's alone

## Notes

- Opened by the operator's request on 2026-10-17 after they reported the Components page unusable; they asked for it in ready with a high cost of delay. The input, a penalty of 1000 a week, is the highest this repository has given a story or an epic, because E-0026's whole operator-facing outcome is the page and it fails on every call.
- E-0026, the component library epic, is done and archived, so this story has no parent.
- Compatibility, for the release notes: a fixed flai under today's dashboard answers list and status, since it ignores `project`, but start and stop stay refused by the channel until the dashboard is upgraded too; a fixed dashboard under an older flai is refused as today. Both components are released together, and the host's flai is upgraded before the page is checked.
- `flai components list --project` and `flai components serve --project` on the command line do not use the channel and keep their flag.
- Criterion 5 is not verified here: it needs this story's flai and flaiover released and the host's flai serve and dashboard upgraded to them, which only the operator does, and an agent must not start a second flai serve. What was checked, on 2026-10-18 in the main checkout: the installed flai 1.57.0 refuses `flai hostapi components.status '{"project":"sf"}'` with `system-flow.yaml has no sub-project "sf"`. A flai built from story/S-0448 answers it with flaiover's dev server, answers `components.status` with `"sub_project":"flaiover"`, and answers `components.list` with the key, listing flai, flaiover with 54 components, and template. The page itself, starting the dev server, framing a component, and stopping it, is the operator's check once both are released.
