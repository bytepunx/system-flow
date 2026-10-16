---
id: T-1860
type: task
nature: feature
title: flai/internal/components starts the sub-project's own vite on 127.0.0.1 and a free port, waits for it to answer, and stops it
status: backlog
parent: S-0438
owner: alex
created: 2026-10-16T21:33:43Z
updated: 2026-10-16T21:33:43Z
transitions: []
stream: S-0438
tags: [cli]
touches: [flai/internal/components/devserver.go, flai/internal/components/devserver_test.go]
after: [T-1858]
---
# T-1860 flai/internal/components starts the sub-project's own vite on 127.0.0.1 and a free port, waits for it to answer, and stops it

## Work

`devserver.go`: `Start(ctx, sub-project, config)` picks a free loopback port and runs the sub-project's `node_modules/.bin/vite` with `--config`, `--host 127.0.0.1`, `--port`, and `--strictPort`, in the sub-project's folder. It waits until the harness URL answers, with a timeout, and returns the URL and the process. `Stop` ends the process group. A missing `vite` is an error naming the fix, and an exit before it answers returns Vite's last lines.

Waits for T-1858.

## Done when

Tests with a fake `vite` script cover the start, the wait, an early exit with its lines, a timeout, and the stop leaving no process. `flai test flai/internal/components` passes.

## Notes

None.
