---
id: D-0611
title: flai serve's checks wait for the image job before their timeout starts, outside flai serve checks timeout, so a long image build does not time the checks out
kind: decision
sources: [S-0751]
topics: [docker, code, go]
paths: [flai/internal/serve/checks.go, flai/internal/serve/imagejob.go]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0611 flai serve's checks wait for the image job before their timeout starts, outside flai serve checks timeout, so a long image build does not time the checks out

S-0751 made a missing or stale agent image a host job that holds every start needing it, the checks' container among them (ADR-0295 §5). It placed the checks' wait for that job ahead of their time limit. From S-0751's Decisions (`wip/archive/agents/S-0751.md`):

```text
The checks' wait runs outside the checks' timeout, so a long build does not time them out.
```

In `flai/internal/serve/checks.go` the wait for the image job comes before the `context.WithTimeout` of `flai serve checks timeout`, and a cancel ends the wait as it ends a check. An image build can take many minutes, so a story that changes the checks' run should keep the wait outside that timeout. Otherwise a checks run that only waited for its image would fail as timed out. `design/system/flai-cli.md` says the checks wait for the job (**The container harness**), but not where the timeout starts.

Source: S-0751.
