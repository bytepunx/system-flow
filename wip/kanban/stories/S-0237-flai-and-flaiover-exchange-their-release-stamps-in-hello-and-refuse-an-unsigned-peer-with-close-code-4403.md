---
id: S-0237
type: story
nature: feature
title: flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-08T04:25:41Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [flai/internal/channel, flaiover/src/lib/server/agent.ts, docs/operators, docs/users/flaiover.md, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/HostFlaiBanner.svelte, flaiover/src/lib/components/HostFlaiBanner.svelte.test.ts, flai/cmd/serve.go, flai/cmd/host.go, design/system/dashboard-host-channel.md, design/system/flaiover-dashboard.md, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flaiover/src/routes/api/projects/+server.ts, flaiover/src/routes/api/projects/projects.test.ts, flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/settings.ts, flai/internal/buildinfo/stamp.go, flai/internal/buildinfo/stamp_test.go, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts, flaiover/src/routes/api/agent/+server.ts, flaiover/src/routes/api/agent/agent.test.ts, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts, flai/cmd/serve_test.go, flai/cmd/host_test.go, docs/users/flai.md, design/system/flai-cli.md, design/system/release-signing.md]
after: [S-0235]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 381
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 433
          cache_read: 3118479
          cache_write: 11596
          cost: 0.7713
cost_of_delay:
  value: 3.95
  by: planner-E-0015
  at: 2026-10-07T22:13:45Z
forecast:
  duration: 1h30m
  delivery: 2026-10-08T09:07:00Z
  basis: "Its own forecast of 1h30m; 12th in the pull order with an in-progress limit of 3, behind S-0232, S-0316, S-0324, S-0318, S-0320, S-0319, S-0309, S-0312, S-0326, S-0287, S-0233, S-0234, S-0235 and S-0236."
  by: flai
  at: 2026-10-08T04:25:41Z
---
# S-0237 flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403

## Goal

The channel's handshake carries each side's release stamp; each verifies the other's with the public key and the component it names before the credential proof, and a missing, unverifiable, or wrong stamp closes the connection with `4403 unsigned peer` (ADR-0070). The refusal is visible on both sides.

## Acceptance criteria
- [ ] `hello` carries flai's `release` (statement and signature) and the answer carries the dashboard's; each side verifies the other's signature with the embedded key, checks the component (`flaiover` from the dashboard, `flai` from flai), and otherwise closes with `4403`; the protocol number is bumped and an older peer's hello without the field is treated as unsigned.
- [ ] flai backs off a minute on `4403` as on `4409`, logs one `error` with the peer's version and the reason, and `flai serve status` and `flai host status` say "the dashboard refused this flai: unsigned" or "flai refused the dashboard: unsigned".
- [ ] The dashboard shows a project whose flai it refused with the reason where it shows a flai that lacks a required method, and its connection list says per project whether the flai is signed.
- [ ] Tests on both sides cover a signed peer, an unsigned one, a wrong component, and an older peer without the field.
- [ ] The operator documentation states the promise and its limit as `release-signing.md § Verifying the peer` does, and `docs/users/flaiover.md` describes what the connection list shows.

## Tasks
- T-1247 flai sends its release stamp in hello, verifies the dashboard's before the proof, and closes or backs off a minute on 4403
- T-1248 The dashboard verifies flai's release stamp in hello before the proof, answers with its own, and closes an unsigned flai with 4403
- T-1249 flai serve status and flai host status say which side refused the other as unsigned
- T-1250 The dashboard shows a project whose flai it refused as unsigned where it shows a flai that lacks a required method
- T-1251 The dashboard's connection list says per project whether its flai is signed or was refused
- T-1252 The operator, user, and design documentation describe the stamp exchange, 4403, its promise and limit, and the connection list

## Notes

### Planning

Tasks, in three layers:

| Layer | Tasks | Why |
|-------|-------|-----|
| 1 | T-1247 (flai handshake), T-1248 (dashboard handshake) | No path in common; they meet only in the wire format criterion 1 fixes, and both use S-0235's stamp and S-0232's key, on main first |
| 2 | T-1249 (status) after T-1247; T-1250 (banner) and T-1251 (connection list) after T-1248 | Each shows the refusal its layer-1 task records; T-1250 and T-1251 share no path |
| 3 | T-1252 (docs) after T-1249, T-1250, and T-1251 | Documents the exchange and the words as built |

Touches:

- Declared, all kept: `flai/internal/channel`, `flaiover/src/lib/server/agent.ts` and its test, `docs/operators`, `docs/users/flaiover.md`, `flai/cmd/serve.go`, `flai/cmd/host.go`, `flai/internal/serve/serve.go` and its test, `flaiover/src/lib/hostflai.svelte.ts`, `flaiover/src/lib/components/HostFlaiBanner.svelte` and its test, `flaiover/src/routes/api/projects/+server.ts` and its test, `flaiover/src/lib/components/ProjectSwitcher.svelte`, `flaiover/src/lib/settings.ts`, `design/system/dashboard-host-channel.md`, `design/system/flaiover-dashboard.md`.
- Layout, added by this run: `flai/internal/buildinfo/stamp.go` and its test, and `flaiover/src/lib/server/release.ts` and its test, where S-0235 verifies a stamp and each side gains a peer's check; `flaiover/src/routes/api/agent/+server.ts` and its test, which turn a flai's missing methods into the error the banner shows; `flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts`; `flai/cmd/serve_test.go`; `flai/cmd/host_test.go`, new, since no test of `flai host status` exists.
- Design, added by this run: `docs/users/flai.md` and `design/system/flai-cli.md`, which describe `flai serve status`; `design/system/release-signing.md`, which records the stamp exchange as built.
- Co-change: `flai touches suggest` listed `design/system/flai-cli.md` (37%) and `docs/users/flai.md` (34%), added above. `flai/internal/hostapi/writes.go` (13%) and `flai/cmd/serve_actions.go` (8%) were not added: no host API method or host action changes.
- Folder touches kept, as declared: `flai/internal/channel`, where T-1247 names `channel.go`, `channel_test.go`, and `channeltest/dashboard.go`; `docs/operators`, where T-1252 names `index.md`. The story's agent may add a runbook page on a refusal. In the claim, both narrow to the files their tasks name (ADR-0096).

Forecast: 1h30m, kept. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 1h1m after this run's touches, from 104 s per unit of size over 33 done large-band feature stories, times size 35.
- S-0235 has five tasks, two of them code, and was planned at 1h. This story has four code tasks with tests on both sides, a protocol 1 compatibility case, and docs, so 1h30m stands over flai's 1h1m.

Cost of delay: 3.95 USD a week, as `flai cod` gives it: this story's 1h30m share of the 9h30m forecast over E-0015's eight open stories, of the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure is closed only when the chain is done, so a share by work fits.
