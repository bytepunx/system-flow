---
id: D-0174
title: A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout
kind: cause
sources: [I-0083]
topics: [cli, conventions, flai, go, template, code, all]
paths: [flai/internal/harness, flai/internal/mcpserver, flai/internal/guard, flai/cmd/guard.go, design/conventions/delegation.md, template, ".claude/settings.json", design/system/flai-cli.md, design/system/agent-context.md, docs/users, docs/operators, design/issues, design/adrs, flai/cmd/guard_test.go]
date: 2026-10-06
status: absorbed
absorbed_by: "design/system/agent-context.md#waiting-for-a-sub-agent"
asked_again: []
cited: []
---

# D-0174 A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout

I-0083 was closed on 2026-10-06.

The cause:

```text
A story's agent that flai serve starts hands its tasks and its close-out to sub-agents. Claude Code runs them in the background, and 2.1.290 does so whether or not the launch asks for it: S-0282's two verifier launches passed no `run_in_background` and came back as background launches. The agent then has to wait. In S-0218 and S-0282 it waited by holding the MCP tool `wait_for_events` for 300 to 1200 seconds.

`wait_for_events` returns on a change to a thread, a work item, or a narrative. A sub-agent writes none of these, so the wait runs to its timeout. Claude Code queues the sub-agent's finished notice while the tool call is open and delivers it only when the call returns. The agent therefore learns that a four-minute verifier has finished up to twenty minutes after it asked. It also asks for a longer timeout each time, so the cost of each wait grows through the story.

flai serve's start prompt and `delegation.md` already say not to do this, since T-0869, which closed remediation 4 of [I-0059](I-0059-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md): "Wait for a sub-agent run in the background through the harness's notice that it has finished, not by polling the flai MCP tool wait_for_events, which reports work items and threads, not sub-agents." S-0218's and S-0282's agents both had that sentence in their prompt and held the wait anyway, nine times between them.

The instruction is not always ignored. S-0219's agent, on the same prompt and the same Claude Code version as S-0282's, launched ten sub-agents in the background and ended its turn after each launch. Claude Code kept the session and began the agent's next turn within seconds of each notice. It never waited on a finished sub-agent.

That way of waiting is safe only for a short sub-agent. Every one of S-0219's finished within about four minutes. S-0220's agent waited the same way twice for a longer one, and Claude Code ended the process ten minutes after the turn ended, with the sub-agent still working ([I-0084](I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md)). Holding `wait_for_events` is slow, but it is what kept S-0282's process alive.

Three things make the wrong choice likely:

- The sentence says what not to do and names a notice, but not how to wait for one. Ending the turn works for a sub-agent that finishes within ten minutes and loses the session for a longer one, and the same prompt uses "end" for finishing the story and for stopping with a question open.
- The sentence is conditional on "a sub-agent run in the background", which reads as a choice the agent no longer makes.
- `wait_for_events` is the only blocking tool the agent has, and its description says "Hold this when idle".

The cost is wall-clock time on a story that is otherwise done, and a longer cycle time in the flow metrics. In S-0282 the waits were 57m56s of 1h52m.
```

The fix:

```text
S-0285: flai serve's start prompt, delegation.md, and the wait_for_events description say to launch sub-agents with run_in_background false, which returns the result as the tool's result however long the sub-agent runs (measured on Claude Code 2.1.290), and that wait_for_events is for a thread awaiting the designer; flai guard records a session's running sub-agents from SubagentStart and SubagentStop and refuses a story agent's wait_for_events while one runs and no thread on its story is open (ADR-0092)
```

Source: I-0083.
