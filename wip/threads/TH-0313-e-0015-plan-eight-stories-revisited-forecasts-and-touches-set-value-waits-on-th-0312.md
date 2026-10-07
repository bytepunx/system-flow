---
id: TH-0313
title: "E-0015 plan: eight stories revisited, forecasts and touches set, value waits on TH-0312"
anchor:
  path: wip/kanban/epics/E-0015-secure-dashboard-and-cli-releases.md
  item: E-0015
status: open
participants: [planner-E-0015]
created: 2026-10-07T19:35:30Z
updated: 2026-10-07T19:35:30Z
---

# TH-0313 E-0015 plan: eight stories revisited, forecasts and touches set, value waits on TH-0312

On wip/kanban/epics/E-0015-secure-dashboard-and-cli-releases.md.

## Entries

### 2026-10-07T19:35:30Z planner-E-0015
E-0015's eight open stories, S-0232 to S-0239, still deliver its outcome as ADR-0070 and `release-signing.md` § Decision lay it out. I gave each of them file touches and a forecast, and recorded both under `### Planning` in its Notes. No story was added, dropped, or reworded. The cost of delay values wait on TH-0312.

## Order

| Layer | Story | After | Forecast | Delivery |
|-------|-------|-------|----------|----------|
| 1 | S-0232 key pair, `checksums.txt` signed, public key in both | none | 1h | 2026-10-08T02:25Z |
| 2 | S-0233 self-upgrade, host upgrade, `install.sh` verify | S-0232 | 1h15m | 2026-10-08T10:58Z |
| 2 | S-0234 flaiover digest list signed on a GitHub release | S-0232 | 45m | 2026-10-08T07:33Z |
| 2 | S-0235 signed release stamp in both builds | S-0232 | 1h15m | 2026-10-08T10:58Z |
| 3 | S-0236 `flai dashboard` runs the image by signed digest | S-0233, S-0234 | 1h15m | 2026-10-08T19:32Z |
| 3 | S-0237 stamps in `hello`, `4403` on failure | S-0235 | 1h30m | 2026-10-08T21:15Z |
| 4 | S-0238 Docker digest check before dialling | S-0236 | 1h | 2026-10-09T02:23Z |
| 5 | S-0239 `dashboard.allow_unsigned`, shown everywhere | S-0237, S-0238 | 1h30m | 2026-10-09T12:40Z |

## Tasks

None drafted. All eight stories are finalized, so the agent that pulls each one writes its tasks.

## Assumptions

- **Forecasts.** `flai forecast` gave 16m to 27m from only 3 medium-band stories. 19 done feature stories of 4 to 7 criteria took a median of about 1h of agent time, so I raised each forecast toward that, by its scope. Deliveries follow flai's playout at its cycle factor of 6.85, with each story pulled as soon as it can be.
- **S-0232 waits on you.** Its criterion 1 has you run the key runbook once and set the two GitHub secrets. Its delivery does not count that wait.
- **Declared touches kept.** I added files and kept every folder touch the stories declared, as the planner must.
- **Layer 2 runs one at a time.** S-0233, S-0234, and S-0235 all claim `docs/operators`. S-0234 and S-0235 both claim `.github/workflows/release-flaiover.yml` and `design/tech/docker.md`. Each holds the others while in progress.

## Proposed changes

These would change finalized stories' words or touches, so I am asking rather than doing them.

1. **Rollback past signing (S-0233, S-0236).** S-0298 lets the Updates page install an older flai or dashboard release (ADR-0117, ADR-0118). Every release from before S-0232 has no signature. Under S-0233's criterion 1 and S-0236's criterion 1, installing one is then refused. Recommendation: add a criterion to each story. The version lists (`flai self-upgrade --list`, `flai host versions`, `flai dashboard versions`, the Updates page) would mark a release that cannot be verified, and the Updates page would not offer to install it.
2. **Narrow the wide folder touches (S-0237, S-0238, S-0239).** `flai/internal/serve`, `flaiover/src/routes`, and `flaiover/src/lib/components` claim most of the host and the dashboard. While one of these stories is in progress, it would hold almost every other ready story. Recommendation: let me replace them with files:
   - `flai/internal/serve/dashboards.go`
   - `flaiover/src/routes/host/+page.svelte` (S-0237)
   - `flaiover/src/routes/+layout.svelte` (S-0239)
   - a new banner component under `flaiover/src/lib/components`: the only folder touch I would keep, for S-0239.
3. **No split, merge, add, or drop.** S-0239 is the largest, with 19 touches across both components, but 1h30m fits one session. It could split into a flai side and a dashboard side if you prefer.
