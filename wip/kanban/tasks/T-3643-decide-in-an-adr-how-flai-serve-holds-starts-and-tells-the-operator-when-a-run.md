---
id: T-3643
type: task
nature: research
title: Decide in an ADR how flai serve holds starts and tells the operator when a run ends on the harness's failed login
status: backlog
parent: S-0855
owner: alex
created: 2026-11-27T23:59:32Z
updated: 2026-11-27T23:59:32Z
transitions: []
stream: S-0855
tags: [flai]
touches: [design/adrs/0323-a-run-that-ends-on-the-harness-s-failed-login-holds-every-start-on-its-account.md, design/adrs/README.md]
---
# T-3643 Decide in an ADR how flai serve holds starts and tells the operator when a run ends on the harness's failed login

## Work

Propose the fix from I-0412's instances, as the story's goal asks, before any code. Read a failed run's log on the host, one of `~/.flai/serve/agents/sf-orchestrator-20261117T17…Z.log`, and record the fields of its `result` line: `is_error`, `api_error`, `api_error_status`, the text, and `terminal_reason`. Then read ADR-0209 §5, ADR-0214, ADR-0108, and ADR-0270.

Write the ADR with the next free number. The file named in touches is tentative: 0322 is claimed by S-0852, and other planners may take more. The recommended decision, to be confirmed or changed against the log:

- A claude-code result that is an error on authentication is a limit of a kind of its own, `login`. That is status 401, or `api_error` `authentication_error`, or the text `Failed to authenticate`.
- It sets the account's budget state as a usage limit does: it holds every strategic and story start on that account, with no reset time, so the 5m backoff applies and doubles to 1h. Only a probe goes through, and the probe's end on no limit clears the state.
- It is announced as the budget state is, with wording that says the login failed and how to renew it (`claude /login` on the host).
- flai serve opens one thread to the operator in each served project the state holds, and resolves it when the state clears.
- It is not an orchestrator failure that `orchestrateRetry` retries.

Name the alternatives: a generic bound on consecutive failed runs of any role, or a check of the login before each start. Say why each was or was not chosen. Add the ADR's row to `design/adrs/README.md`.

This task waits for nothing.

## Done when

- The ADR is written and indexed in `design/adrs/README.md`, naming the result fields read from the real log.
- Its decision says what is detected, what is held and until when, what tells the operator, and what clears it.
- `flai check --strict` passes on the story.

## Notes
