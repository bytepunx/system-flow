---
id: T-0286
type: task
nature: feature
title: Host actions are off until enabled by name on the host, and an acceptance from the board pushes nothing unless push is (I-0028)
status: done
parent: S-0078
owner: alex
created: 2026-09-20T13:06:48Z
updated: 2026-09-20T13:15:04Z
transitions:
  - to: ready
    at: 2026-09-20T13:06:49Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T13:06:49Z
    by: system-flow
  - to: done
    at: 2026-09-20T13:15:04Z
    by: system-flow
stream: S-0078
tags: []
---
# T-0286 Host actions are off until enabled by name on the host, and an acceptance from the board pushes nothing unless push is (I-0028)

## Work
The host configuration gains the actions enabled and for which projects; flai serve enable, disable, and actions manage it, per project unless all projects are named. hostapi asks a gate before a host action. accept.run passes --no-push unless the push action is enabled for the project, which closes I-0028. A disabled action answers with a typed error that says what to run to enable it. No method of the channel reads or writes the configuration, and a test says so.

## Done when
- Tests: the command line of accept.run with the action off and on; the disabled answer; enable and disable per project and for all
- make flai-test passes

## Notes
