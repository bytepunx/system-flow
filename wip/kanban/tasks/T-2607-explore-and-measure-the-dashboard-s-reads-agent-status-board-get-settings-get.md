---
id: T-2607
type: task
nature: research
title: "Explore and measure the dashboard's reads: agent.status, board.get, settings.get, threads, messages, and item.get, and whether flai serve answers channel requests one at a time"
status: in-progress
parent: S-0596
owner: alex
created: 2026-11-07T22:09:51Z
updated: 2026-11-07T22:12:47Z
transitions:
  - to: ready
    at: 2026-11-07T22:12:44Z
    by: agent-S-0596
  - to: in-progress
    at: 2026-11-07T22:12:47Z
    by: agent-S-0596
stream: S-0596
tags: [cli]
---
# T-2607 Explore and measure the dashboard's reads: agent.status, board.get, settings.get, threads, messages, and item.get, and whether flai serve answers channel requests one at a time

## Work

Read the code behind the slowest reads in `~/.flai/serve/serve.log` since 2026-11-04 (agent.status p50 796 ms over 84,417 calls, about 30% of serve's wall time; board.get p50 800 ms; settings.get p50 825 ms; inbox.designer and threads.list about 600 ms in threads.read). Time each with `flai hostapi --timing` on this repository. Say for each where the file I/O goes, which reads run one after another where they could run at once, what is read again on every request that could be kept, and whether flai serve answers the dashboard's requests on one channel one at a time. Waits for no task: it writes nothing and returns its findings to the story's agent.

## Done when

- [ ] The findings name each method's cost, its I/O pattern by path and line, and each serial read that could run in parallel or be kept, with a measurement.

## Notes
