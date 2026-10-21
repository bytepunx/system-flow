---
id: I-0163
title: The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project
class: defect
status: closed
count: 1
cost: 20m
first_reported: 2026-10-18T02:25:28Z
last_reported: 2026-10-18T02:25:28Z
updated: 2026-10-21T12:12:03Z
---

# I-0163 The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project

## Description
The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project

## Instances

### 2026-10-18T02:25:28Z
Reported by the operator from the Components page: 'Could not read the components: system-flow.yaml has no sub-project "sf"; name one of flai, flaiover, template'. The dashboard's agent sends every call as { project: <project key>, ...params } (flaiover/src/lib/server/agent.ts:594) and flai's channel refuses a call whose project is not its key (flai/internal/channel/channel.go:850), while componentsParams (flai/internal/hostapi/components.go:188) reads that same field as the sub-project. components.list and components.status, sent with no sub-project, take the key sf as the sub-project and are refused; components.start and components.stop, sent with project: flaiover, overwrite the key and the channel refuses them as another project's. Reproduced: flai hostapi components.status '{"project":"sf"}' gives the error, '{}' answers. The Go tests call the handlers directly and the dashboard's tests use a fake source, so no test sent a call through the channel.

## Impact
The Components page is the whole of E-0026's operator-facing outcome and fails on every call made through the dashboard; the operator reported it unusable on 2026-10-17

- penalty_per_week: 1000

## Remediation

Story S-0475 remediates this issue, created from it at 2026-10-21T11:59:09Z.
Closed 2026-10-21T12:12:03Z: Fixed by S-0448, accepted 2026-10-18: the components methods take their sub-project as sub_project (382fab56f), and the dashboard's agent puts the channel's project key after a call's params (d94c114d4). Released in flai v1.63.1 and flaiover v0.52.5. S-0475 confirmed it: flai hostapi components.status '{"project":"sf"}' now answers the status, and the regression tests TestComponentsAnswerThroughTheChannelWithTheProjectKeyInTheirParams, TestNoMethodTakesAParameterNamedProject, and the S-0448 cases in flaiover's agent.test.ts and components.test.ts pass.
