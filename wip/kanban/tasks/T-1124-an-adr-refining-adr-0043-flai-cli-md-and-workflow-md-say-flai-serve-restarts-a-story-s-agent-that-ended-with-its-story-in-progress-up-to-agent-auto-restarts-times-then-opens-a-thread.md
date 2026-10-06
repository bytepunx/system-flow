---
id: T-1124
type: task
nature: remediation
title: An ADR refining ADR-0043, flai-cli.md, and workflow.md say flai serve restarts a story's agent that ended with its story in progress, up to agent.auto_restarts times, then opens a thread
status: backlog
parent: S-0294
owner: alex
created: 2026-10-06T23:07:37Z
updated: 2026-10-06T23:09:13Z
transitions: []
stream: S-0294
tags: [docs, serve]
touches: [design/adrs, design/system/flai-cli.md, design/system/workflow.md]
---
# T-1124 An ADR refining ADR-0043, flai-cli.md, and workflow.md say flai serve restarts a story's agent that ended with its story in progress, up to agent.auto_restarts times, then opens a thread

## Work

ADR-0043 gave a story's restart to a person. I-0084's second remediation hands it to flai serve, within a limit, on the operator's answer on TH-0208.

Write the ADR with `flai adr new`, refining ADR-0043, and record the decision the operator took on TH-0208:

- **When.** flai serve's launcher restarts a story's agent once its run ends, if all of these hold:
  - the story is still in progress;
  - the story is not blocked;
  - the run did not end asking a question (`OutcomeAsked`, which the restart on an answer already covers);
  - the operator did not stop the run (ADR-0058);
  - the `agent` host action is on.
- **How.** It starts the agent as `flai serve agent restart` does: a new session, the same agent name, the run recorded in `serve/agents.json`, and the new agent told how the last one ended.
- **The limit.** The host setting `agent.auto_restarts` sets how many times, 2 by default; 0 turns it off. The count is per story. It resets when the operator restarts the agent by hand, or when the story leaves in-progress.
- **At the limit.** flai serve opens one thread on the story. It says how many times the agent ended and how the last run ended, and names `flai serve agent restart` and Retry. The story is then left as before.

Then bring the living design into line:

- `design/system/flai-cli.md`'s `flai serve` section: the restart, the setting, and the thread.
- `design/system/workflow.md`'s paragraph that says a story whose agent dropped or failed gets a new one "on the operator's word": add the automatic restart and the limit.

## Done when

- The ADR is in `design/adrs`, proposed, refining ADR-0043, and links I-0084 and TH-0208.
- `flai-cli.md` and `workflow.md` describe the automatic restart, its conditions, `agent.auto_restarts`, and the thread at the limit.
- `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner. The operator confirmed TH-0208's recommendation on 2026-10-06. The ADR file's name is not known until `flai adr new` numbers it, so the task touches `design/adrs`, a shared path.
