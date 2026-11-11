---
id: D-0513
title: flai serve's project loop makes its one-minute time.After again on each pass, and it stays unchanged until a measurement shows a minute look delayed
kind: decision
sources: [S-0589]
topics: [flai, cli, go, code]
paths: [flai/internal/serve/serve.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0690]
---

# D-0513 flai serve's project loop makes its one-minute time.After again on each pass, and it stays unchanged until a measurement shows a minute look delayed

S-0589's planner noted that the project loop in `flai/internal/serve/serve.go` makes `time.After(time.Minute)` again on each pass. Wakes more often than once a minute, from changes under `kanban/` and `threads/`, could then put off the minute look.

S-0589 left it as it is. On the instance it studied the look came 43 seconds after the enable, so it was not delayed. No measurement shows the delay happening, and `code-quality.md` asks for no change without one. A later story that wants to change the loop's timer should first measure a delayed minute look.

- S-0589: "No measurement shows it happening, so it is not changed (code-quality.md: no change without a measurement)."

Source: the Decisions of wip/archive/agents/S-0589.md.
