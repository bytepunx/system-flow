---
id: T-1839
type: task
nature: feature
title: The host API's thread.reply and thread.confirm write the operator's answer from the dashboard as a digest entry
status: backlog
parent: S-0433
owner: alex
created: 2026-10-16T06:07:44Z
updated: 2026-10-16T06:07:44Z
transitions: []
stream: S-0433
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-1837]
---
# T-1839 The host API's thread.reply and thread.confirm write the operator's answer from the dashboard as a digest entry

## Work

In `flai/internal/hostapi/writes.go`, after `thread.reply` and `thread.confirm` save the thread, call `digest.OnAnswer` and add `digest` to the response with the entry's ID, or its warning.

Waits for T-1837. Runs beside the CLI and MCP tasks.

## Done when

`writes_test.go` covers a reply and a confirmation from the dashboard that write an entry, and a warning that leaves the reply saved; `flai test flai/internal/hostapi` passes.

## Notes

The dashboard need not show the entry; the response carries it for a later page.
