---
id: T-1927
type: task
nature: remediation
title: Record in an ADR and the design that a thread awaits an agent only when the agent takes part in it or it is on the agent's own story, and is awaiting operator or other otherwise
status: backlog
parent: S-0453
owner: alex
created: 2026-10-20T22:20:10Z
updated: 2026-10-20T22:20:10Z
transitions: []
stream: S-0453
tags: [flai]
touches: [design/adrs/README.md, design/system/agent-narrative.md, design/system/flai-cli.md]
---
# T-1927 Record in an ADR and the design that a thread awaits an agent only when the agent takes part in it or it is on the agent's own story, and is awaiting operator or other otherwise

## Work

Write the ADR with `flai adr new`, refining ADR-0179 and ADR-0090, and accept it. It records the rule from criterion 1:

- A thread awaits an agent (`awaiting: "you"`) only when the agent is among its `participants`, or the thread is on the agent's own story or one of its tasks, and its last entry is not the agent's.
- ADR-0090's rule stands: an opener whose recommendation awaits confirmation is not awaited.
- Any other open thread is `awaiting: "operator"` when it awaits the operator, else `"other"`.
- The operator is the manifest's `owner` or the thread's story's `owner`, the two names ADR-0097 takes an answer from. A thread awaits the operator when a recommendation is pending or its last entry is by neither. This is the rule `inbox.designer` in `flai/internal/hostapi/people.go` already applies, so the dashboard is unchanged.
- A reply already adds its author to `participants` (`threads.ReplyWith`), so an agent that replied on a thread takes part in it.
- The orchestrator's `inbox` lists the threads awaiting the operator by default, as its role, without `all: true`.
- The agent's own story is the one `recordingStory` names: `FLAI_STORY`, else `agent-S-nnnn` in the agent's name, else the story branch checked out.

Update `design/system/agent-narrative.md` § What an agent is told through MCP and `design/system/flai-cli.md` § Commands where they say what `inbox` and `wait_for_work` report as awaiting the agent.

It waits for nothing: the code tasks build what it records.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`.
- The two design sections state the rule and link the ADR.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes

The ADR's own file is named by `flai adr new`; `flai task done` widens the touches with it.
