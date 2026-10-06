---
id: I-0084
title: Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent
class: blocker
status: open
count: 2
cost: 10m
first_reported: 2026-10-06T07:00:22Z
last_reported: 2026-10-06T07:00:22Z
updated: 2026-10-06T11:44:49Z
---

# I-0084 Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent

## Description

flai serve runs a story's agent headless, `claude -p`. When the agent ends a turn while a sub-agent it launched in the background is still running, Claude Code 2.1.290 keeps the process for ten minutes. A sub-agent that finishes in that time begins the agent's next turn with its notice. One that does not is ended with the process, mid-work, and the session is over. In the one case timed to the second, the process ended 10m00s after the turn did.

flai serve then records the run as failed and does nothing more. The story stays in progress with no agent, and every ready story it holds by overlap stays held, until someone restarts the agent: `flai serve agent restart`, or the story page's Retry button. On 2026-10-06 that was ten ready stories behind S-0220.

Ending the turn is what the start prompt's sentence on sub-agents leads to when it is followed: "Wait for a sub-agent run in the background through the harness's notice that it has finished". The other habit, holding `wait_for_events`, keeps the process alive and wastes the wait instead ([I-0083](I-0083-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md)). The agent has no way to wait for a long sub-agent that is both safe and prompt.

## Instances

### 2026-10-06T07:00:22Z
Story: S-0220.
S-0220's second session ended its turn at 06:47:45Z saying it was waiting on T-0913, whose task sub-agent ran in the background. The sub-agent was still working at 06:57:39Z, reviewing its diff and running the markdown lint; the process was gone by 06:58Z and flai serve recorded the run as failed, 'ended (an exit code nobody saw) with S-0220 in in-progress'. The sub-agent's twelve changed files stayed in the worktree, uncommitted. Nothing restarted the agent; Claude, watching the board for the operator, ran flai serve agent restart at 06:59:56Z. The cost counts the sub-agent's cut-off work and the gap; unattended, the gap has no end.

Found afterwards in the first session's transcript and flai serve's log: it ended its turn at 06:27:32Z waiting on T-0912, whose sub-agent was held in permission_prompt on TH-0163, and the process ended at 06:37:32Z, exit 0, ten minutes later to the second; flai serve recorded it as failed. The operator restarted the agent from the dashboard at 06:38:50Z.

## Remediation

Two separate things went wrong, and each has its own fix.

1. **The agent has no safe way to wait.** The prompt and `delegation.md` name one that keeps the session alive for a sub-agent of any length: a launch with `run_in_background: false`, once tested against a sub-agent that runs longer than ten minutes. This is the first change I-0083 recommends, and one story can make it for both issues.
2. **A story whose agent ended without finishing waits for a person.** flai serve could restart such an agent itself: once, or a small number of times, when the run ended with the story still in progress, no thread awaiting the operator, and no block on the story; and tell the operator when the limit is reached. Whether flai serve restarts on its own is the operator's decision, since ADR-0043 gave the restart to a person.

S-0285 made the first remediation ([ADR-0092](../adrs/0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md)). It measured the ten minutes again on 2.1.290: a turn ended at 10:34:48Z with a background sub-agent out, and the process exited at 10:44:49Z with the sub-agent cut off. It also measured the remedy: a launch with `run_in_background` false returned an 11-minute sub-agent's result as the tool's result, and three launched in one message ran together. The start prompt and `delegation.md` now name that way and say never to end the turn on a background sub-agent. The second remediation is still open.

Story S-0294 remediates this issue, created from it at 2026-10-06T11:44:49Z.
