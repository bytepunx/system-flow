---
id: T-2690
type: task
nature: improvement
title: agent.status takes a finished parameter, and without it sends only the open and live runs
status: backlog
parent: S-0611
owner: alex
created: 2026-11-07T23:55:49Z
updated: 2026-11-07T23:55:49Z
transitions: []
stream: S-0611
tags: [cli]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/contract_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/system/flaiover-dashboard.md]
---
# T-2690 agent.status takes a finished parameter, and without it sends only the open and live runs

## Work

Criterion 2, the host side. `agent.status` (`flai/internal/hostapi/hostapi.go`, about line 351) takes no parameter, and `host().Agent` (`flai/cmd/serve_actions.go`, about lines 47 to 84) answers every story's newest run (`stories`, `last`), every item's newest planner run (`plans`), and the past runs of each role (`past_orchestrators`, `past_analyzers`, `roles.<name>.runs`). That is 437 finished runs and 288 planner runs, 392 KB.

- Give `agent.status` a boolean parameter, `finished`. With it, the answer is as today. Without it, the answer keeps:
  - each run that is live, and each story run whose activity is `waiting` (an agent that ended asking, a held story, or one run elsewhere), as `serve.Activity` and `runActivity` tell them (`flai/internal/serve/agents.go`, about lines 2657 to 2768);
  - the newest run, `worked` or `failed`, of each story or planned item that is not done, cancelled, or archived, so the board's `failed` dot and a story's agent panel keep theirs;
  - `running`, `waiting`, `enabled`, `command`, `budget`, and each role's state and run under way.
  It leaves out the past runs, and the planner runs and story runs of items that are done, cancelled, or archived.
- Change `Host.Agent` (`flai/internal/hostapi/writes.go`, about line 554) to take the parameter, and every test that builds a `Host`.
- Document the parameter in the method's comment above `agent.status` and in `design/system/flaiover-dashboard.md` § API, beside the runs the pages list.
- `agent.stream` is unchanged.

Waits for no task. Changes no file another task of this story changes; T-2691 reads the parameter this one adds.

## Done when

- [ ] A test with a live run, a waiting run, a failed run on an open story, and finished runs on a done story and a done planned item shows the answer without `finished` keeps the first three and drops the rest, and with `finished: true` keeps all of them.
- [ ] The method's comment and `design/system/flaiover-dashboard.md` § API name `finished` and what it adds.
- [ ] `flai test` passes on the changed paths.

## Notes

- Planned by planner-S-0611. S-0603 changes `AgentStates` and `serve_actions.go` for cause 12; settle any overlap in the two stories' conversation.
