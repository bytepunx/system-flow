---
id: T-0337
type: task
nature: feature
title: A claim kept with the repository stops flai serve from dialing a project another live process already serves
status: done
parent: S-0084
owner: alex
created: 2026-09-23T01:08:32Z
updated: 2026-09-23T01:09:52Z
transitions:
  - to: ready
    at: 2026-09-23T01:08:59Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T01:09:06Z
    by: system-flow
  - to: done
    at: 2026-09-23T01:09:52Z
    by: system-flow
stream: S-0084
tags: []
---

# T-0337 A claim kept with the repository stops flai serve from dialing a project another live process already serves

## Work
`flai serve`'s own "one per user" guarantee (ADR-0029) is enforced only per configuration folder (`Dir.ReadStatus`, keyed on `DirFor(configPath)`): two flai processes under different configurations (an install under `~/.flai`, a tree's own `.flai-cache`) never see each other, so both dial the same project's dashboard and fight over its one connection (I-0029: 457 reconnects in 45 minutes).

`internal/serve/claim.go` (new): `Claim{PID, Version, Config, Started, Updated}` kept at `<root>/.flai-cache/serve-claim.json` — WITH the repository, not any one configuration — read and written with the same freshness-plus-liveness pattern `Dir.ReadStatus`/`Alive` already use (`ReadClaim`, `WriteClaim`, `ReleaseClaim`). `serve.Run`'s `reconcile()` now reads a project's claim before dialing it: a live claim held by another PID is left alone (logged once, not every tick); otherwise this process writes its own claim first, then dials. A held claim is refreshed every tick so a live holder is never mistaken for gone; a released or dropped root's claim is removed so a clean shutdown or restart hands it over at once rather than waiting out the staleness window. Not a distributed lock — a race between two processes claiming in the same instant is accepted, the same as `flai serve`'s own existing single-instance check already accepts, since the loser's next tick (a second later) finds the winner's fresher claim.

`serve.Options` gained `ConfigPath`, so a claim names which configuration holds it; `cmd/serve.go` passes it.

## Done when
Acceptance criterion 1. Go tests: a claim round-trips, a stale one and a dead one's PID both count as not held, release only undoes this process's own claim; a repository claimed by another live process is never dialed (checked against a real `channeltest` fake dashboard: the connection channel stays empty); a dead holder's claim is taken over within a few ticks. `go test ./...`, `go vet`, `gofmt -l .` clean.

## Notes
A genuinely two-OS-process version of "never dialed" is T-0340's own real, hands-on verification (two actual flai builds); a single test binary cannot fork itself a second PID to race honestly, so the unit test here simulates the other holder by writing its claim directly (the same shape the pre-existing `TestASecondServeRefusesWhileTheFirstRuns` already uses for "another live process," `PID: os.Getppid()`).
