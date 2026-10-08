---
id: T-1204
type: task
nature: feature
title: An overlap on shared paths does not hold, and the held reason names only the paths still held
status: done
parent: S-0334
owner: alex
created: 2026-10-07T20:16:13Z
updated: 2026-10-08T09:48:40Z
transitions:
  - to: ready
    at: 2026-10-08T09:45:12Z
    by: agent-S-0334
  - to: in-progress
    at: 2026-10-08T09:45:12Z
    by: agent-S-0334
  - to: done
    at: 2026-10-08T09:48:40Z
    by: agent-S-0334
stream: S-0334
tags: [flai]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/share.go, flai/internal/workitem/share_test.go]
after: [T-1202]
usage:
  source: log
  seconds: 208
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 7044
      cache_read: 1153844
      cache_write: 34096
      cost: 0.5775
---
# T-1204 An overlap on shared paths does not hold, and the held reason names only the paths still held

## Work

Make the hold honour shares, as ADR-0134 says. It waits for T-1202, the ADR. It comes before T-1203, since the messages package writes the `workitem.Share` type this task defines and reads (messages imports workitem, not the other way).

- `workitem.Share`: the holding story, the held story, the paths, the split, who shared, and when, with the YAML keys a conversation's `shares` front matter takes.
- `(r *Repo) Shares()` reads the shares in force from the conversations under `<wip>/messages`: stored open, both stories in ready, in progress, or in review, and neither moved to backlog, done, or cancelled since the share was made.
- `(h *Holds) WithShares(shares)`; `Repo.Holds` passes `r.Shares()` in, so `flai board`, `inbox`, `wait_for_work`, `flai story start`, and the launcher agree. An overlapping pair does not hold when the narrower path lies wholly inside a path the holding story shares with the held one.
- The `held` reason names only the paths still held, and every story that still holds.
- `(h *Holds) OverlapsBy(story)`: for a ready story held on overlap alone, each in-progress story that holds it with the overlapping paths still held, for flai serve's ask (T-1205).

## Done when

- Tests cover a share that clears a hold, one that clears part of it, a share with a third story that clears nothing, a share that ended, and `OverlapsBy`.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
