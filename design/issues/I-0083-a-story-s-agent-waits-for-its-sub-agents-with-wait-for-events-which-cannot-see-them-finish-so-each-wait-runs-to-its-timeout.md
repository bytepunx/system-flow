---
id: I-0083
title: A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout
class: efficiency
status: closed
count: 2
cost: 39m
first_reported: 2026-10-06T06:02:33Z
last_reported: 2026-10-06T06:02:38Z
updated: 2026-10-06T10:54:40Z
---

# I-0083 A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout

## Description

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

## Instances

### 2026-10-06T06:02:33Z
Story: S-0282.
The story's agent launched eight sub-agents in the background and after each launch or layer held wait_for_events for 300 to 1200 seconds. All six waits timed out with nothing changed. Each sub-agent's finished notice was queued during the wait and delivered when it returned: 3m32s, 6m17s, 5m15s, 10m15s, 15m50s, and 16m47s after the sub-agent ended, 57m56s of a run of 1h52m. The two close-out verifiers ran about 4 and 3 minutes and each cost 20.

### 2026-10-06T06:02:38Z
Story: S-0218.
Found in the transcript of S-0218's agent, 2026-10-05T07:09Z to 09:31Z, whose prompt already carried the instruction T-0869 added. Three waits on background sub-agents ran to their timeouts of 600, 900, and 900 seconds and returned 1m08s, 8m31s, and 9m25s after the sub-agent's notice was queued: 19m04s, the last of them on the close-out verifier.

## Remediation

Two changes, both needed. The wording T-0869 added did not hold in two stories, so better wording is not enough without a check in flai.

1. **Say how to wait.** The start prompt in `flai/internal/harness`, `delegation.md` here and in the template, and the `wait_for_events` tool description say how the agent waits for a sub-agent when it has nothing else to do, in a way that keeps the session alive however long the sub-agent runs. Ending the turn is not that way: it works while the sub-agent finishes within ten minutes, as S-0219's did, and loses the session and the sub-agent's work past that, as in S-0220 (I-0084). The candidate is a launch with `run_in_background: false`, which returns the result as the tool's result; test it against the Claude Code that flai serve runs, 2.1.290 when this was recorded, with a sub-agent that runs longer than ten minutes, and with a layer of several launched in one message. They say when `wait_for_events` is right: only for a thread awaiting the designer. The sentence no longer depends on the agent having chosen the background.
2. **Make flai catch it.** `flai guard` already runs as a `PreToolUse` hook on the story agent's `mcp__flai__` calls, so flai sees the `wait_for_events` call when it is made. Either of these keeps a mistaken wait short:
   - A Claude Code hook on a sub-agent's stop (`SubagentStop`) tells flai, and a `wait_for_events` held by that session returns at once, saying which sub-agent finished. Confirm first what that hook's input carries.
   - `flai guard` refuses a `wait_for_events` call made while a sub-agent of the session is running and no thread on the story awaits the designer, and its refusal says how to wait instead.

The measure, on a story flai serve works after the change is installed: no wait on a sub-agent outlasts the sub-agent by more than a few seconds.

Story S-0285 remediates this issue, created from it at 2026-10-06T06:03:25Z.
Closed 2026-10-06T10:54:40Z: S-0285: flai serve's start prompt, delegation.md, and the wait_for_events description say to launch sub-agents with run_in_background false, which returns the result as the tool's result however long the sub-agent runs (measured on Claude Code 2.1.290), and that wait_for_events is for a thread awaiting the designer; flai guard records a session's running sub-agents from SubagentStart and SubagentStop and refuses a story agent's wait_for_events while one runs and no thread on its story is open (ADR-0092)
