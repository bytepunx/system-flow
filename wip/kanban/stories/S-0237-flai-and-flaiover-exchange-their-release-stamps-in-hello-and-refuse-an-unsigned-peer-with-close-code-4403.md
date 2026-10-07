---
id: S-0237
type: story
nature: feature
title: flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-07T22:48:25Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [flai/internal/channel, flaiover/src/lib/server/agent.ts, docs/operators, docs/users/flaiover.md, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/HostFlaiBanner.svelte, flaiover/src/lib/components/HostFlaiBanner.svelte.test.ts, flai/cmd/serve.go, flai/cmd/host.go, design/system/dashboard-host-channel.md, design/system/flaiover-dashboard.md, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flaiover/src/routes/api/projects/+server.ts, flaiover/src/routes/api/projects/projects.test.ts, flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/settings.ts]
after: [S-0235]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 3.95
  by: planner-E-0015
  at: 2026-10-07T22:13:45Z
forecast:
  duration: 1h30m
  delivery: 2026-10-08T02:57:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 3, behind S-0232, S-0332, S-0233, S-0234, S-0235 and S-0236."
  by: flai
  at: 2026-10-07T22:48:25Z
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

## Notes

### Planning

Touches:

- Declared: `flai/internal/channel`, `flai/internal/serve`, `flaiover/src/lib/server/agent.ts`, `flaiover/src/routes`, `docs/operators`, `docs/users/flaiover.md`.
- Layout: `flaiover/src/lib/server/agent.test.ts`; `flaiover/src/lib/hostflai.svelte.ts` and `flaiover/src/lib/components/HostFlaiBanner.svelte` with its test, which show a flai the dashboard cannot use; `flai/cmd/serve.go` and `flai/cmd/host.go`, which print `flai serve status` and `flai host status`.
- Design: `design/system/dashboard-host-channel.md` and `design/system/flaiover-dashboard.md` § The channel to flai on the host, which describe `hello`.
- Folder touches kept, as declared: `flai/internal/channel` (`channel.go`, its test, and `channeltest`), where `hello` and the close codes are; `docs/operators`, whose runbooks may gain a page on the refusal.
- Narrowed on TH-0313: `flai/internal/serve` to `flai/internal/serve/serve.go`, which builds the channel client, and its test; `flaiover/src/routes` to `flaiover/src/routes/api/projects/+server.ts` and its test, the per-project connection list, with `flaiover/src/lib/components/ProjectSwitcher.svelte` and `flaiover/src/lib/settings.ts`, which show it.

Forecast: 1h30m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 22m from 116 s per unit of size, over only 3 medium-band feature stories.
- Done feature stories of this size took a median of about 1h of agent time. This one changes the protocol on both sides, with status, UI, and four test cases on each, so 1h30m.
- The first delivery was played out after S-0235 at flai's cycle factor of 6.85.

Cost of delay: 3.95 USD a week, as `flai cod` gives it: this story's 1h30m share of the 9h30m forecast over E-0015's eight open stories, of the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure is closed only when the chain is done, so a share by work fits.
