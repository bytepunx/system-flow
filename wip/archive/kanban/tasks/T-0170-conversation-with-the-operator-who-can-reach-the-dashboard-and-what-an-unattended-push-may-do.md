---
id: T-0170
type: task
nature: research
title: "Conversation with the operator: who can reach the dashboard, and what an unattended push may do"
status: done
parent: S-0052
owner: alex
created: 2026-09-19T01:55:25Z
updated: 2026-09-19T07:56:34Z
transitions:
  - to: ready
    at: 2026-09-19T07:36:30Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T07:36:31Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:56:34Z
    by: system-flow
stream: S-0052
tags: []
---

# T-0170 Conversation with the operator: who can reach the dashboard, and what an unattended push may do

## Work
Ask the operator, through a thread on S-0052 (`thread_open`) or in the session, and record the answers in the finding. Who can reach the dashboard today, and who will once S-0043 makes it reachable through a hub or tunnel? Is a push that nobody on the host approved acceptable, given that release tags publish images? Should every acceptance push, or only when asked from the confirmation? Is a credential stored on the host for this acceptable, and which kind? What should happen when the push fails: leave it accepted locally, as now? Put the recommended answer first with each question. Do the research tasks first so the questions carry what was found.

## Done when
- The questions were asked once, with recommendations, and the operator's answers are quoted in the finding with the date
- Any answer still missing is listed as open, and the story is blocked on it with `flai block --reason` if the recommendation depends on it

## Notes
