---
id: T-0261
type: task
nature: feature
title: flai dashboard hands the container the agent credential, starts flai serve, registers the project, and reports it in status
status: backlog
parent: S-0072
owner: alex
created: 2026-09-20T07:28:28Z
updated: 2026-09-20T07:28:28Z
transitions: []
stream: S-0072
tags: []
---
# T-0261 flai dashboard hands the container the agent credential, starts flai serve, registers the project, and reports it in status

## Work
Generate the credential beside the token, mount it as a secret, start flai serve detached when it is not running and register the project; stop leaves it running; status says whether the host flai runs, since when, which projects, and whether this dashboard has it connected.

## Done when
- Command tests with the fake runner for the mounts and for status
- flai dashboard stop in one project does not stop flai serve

## Notes
