---
id: T-0394
type: task
nature: remediation
title: A serve restart from the dashboard, retried after the connection drops, restarts serve once, tried live
status: done
parent: S-0109
owner: alex
created: 2026-09-24T07:37:31Z
updated: 2026-09-24T07:46:04Z
transitions:
  - to: ready
    at: 2026-09-24T07:37:34Z
    by: agent-S-0109
  - to: in-progress
    at: 2026-09-24T07:42:08Z
    by: agent-S-0109
  - to: done
    at: 2026-09-24T07:46:04Z
    by: agent-S-0109
stream: S-0109
tags: []
---
# T-0394 A serve restart from the dashboard, retried after the connection drops, restarts serve once, tried live

## Work

- On a scratch `flai host` built from this branch, with its own configuration and project, and a flaiover dev server for that project, ask for a serve restart from the dashboard with the retry turned on for that write (a local, uncommitted change to the dev server's copy).
- Read the host's log and the journal: serve restarted once, and the retry was answered from the record.
- Stop everything started, by exact PID. The operator's own `flai serve` is not touched.

## Done when

- What was seen is written in the story's Notes, and criterion 2 is ticked only if serve restarted once.

## Notes
