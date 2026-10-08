---
id: TH-0370
title: "S-0340 plan: six tasks in three layers; three questions"
anchor:
  path: wip/kanban/stories/S-0340-the-smoke-tier-installs-and-self-upgrades-from-a-local-release-server-built-from-the-tree-so-no-close-out-depends-on-github.md
  item: S-0340
status: open
participants: [planner-S-0340, orchestrator]
created: 2026-10-08T08:07:28Z
updated: 2026-10-08T08:07:52Z
---

# TH-0370 S-0340 plan: six tasks in three layers; three questions

On wip/kanban/stories/S-0340-the-smoke-tier-installs-and-self-upgrades-from-a-local-release-server-built-from-the-tree-so-no-close-out-depends-on-github.md.

## Entries

### 2026-10-08T08:07:28Z planner-S-0340
The plan for S-0340: six tasks in three layers. Touches, forecast, and cost of delay are set; the reasons are in the story's Notes under `### Planning`.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1348 `flai-snapshot.sh` builds a host-only release with a plain version | `scripts/flai-snapshot.sh`, `flai/.goreleaser.yaml` | none |
| 1 | T-1350 A local release server answers GitHub's three release endpoints | `flai/internal/releaseserver/{releaseserver,releaseserver_test}.go`, `.../serve/main.go`, `scripts/release-server.sh` | none |
| 1 | T-1352 The GitHub check is its own script and workflow | `scripts/install-published-test.sh`, `.github/workflows/install-published.yml`, `Makefile` | none |
| 2 | T-1354 `install-test.sh` installs from the local server, with GitHub unreachable | `scripts/install-test.sh`, `scripts/smoke.sh`, `install.sh`, `system-flow-check.yml` | T-1348, T-1350 |
| 3 | T-1356 Five documents say what smoke installs from and what the GitHub check covers | `docs/users/flai.md`, `design/tech/ci.md`, `design/system/devex.md`, `design/system/flai-cli.md`, `scripts/README.md` | T-1354, T-1352 |
| 3 | T-1357 Close I-0086, and settle S-0291 | I-0086, `design/issues/summary.md` | T-1354, T-1352 |

T-1354 runs what T-1348 and T-1350 build. T-1356 and T-1357 describe and close what layers 1 and 2 built. The tasks of each layer share no path.

## Assumptions

- **Snapshot versions.** `selfupgrade.List` accepts only plain `X.Y.Z` tags. `install.sh` names the archive from the tag. A snapshot's `0.0.0-<commit>` and `flai_snapshot_*` fit neither. So T-1348 adds a local mode with a plain version. That touches `flai/.goreleaser.yaml`, which the story did not declare.
- **The server is Go.** It lives under `flai/internal/releaseserver`, not in a shell script. It needs no interpreter. Smoke and integration are selected by `flai/**` only, so a change limited to `scripts/` would never run smoke in the story's own close-out. Its `main` is planned at `flai/internal/releaseserver/serve/main.go`, not under `flai/cmd`, which is the cobra package.
- **No change to `selfupgrade.go` or `install.sh` is expected.** `FLAI_RELEASES_API`, the hidden `--api`, `FLAI_API`, and `FLAI_REPO` already exist. Both stay in the touches as declared.
- **CI needs GoReleaser.** `system-flow-check.yml` gets GoReleaser through `goreleaser-action` with `install-only` (T-1354).
- **`self-upgrade --dir`** answers "already the latest" today and downloads nothing, because `install.sh` installed the same release just before. T-1354 keeps that, with a second, newer release as an option.

## Questions

1. **Recording a CI failure of the GitHub check (criterion 4).** `flai issue` writes a file in the runner's checkout, which is then lost.
   - Recommended: on failure, the workflow runs `flai issue bump` or `new` and opens a pull request with the issue file, so it reaches `main` when you merge it.
   - Alternative: CI only fails visibly. The script records through `flai issue` when it runs on a host, such as `make install-published-test`.
   - CI should not push to `main` directly. That would bypass publishing (ADR-0067).
2. **S-0291 (ready) overlaps this story almost entirely.** Recommended: cancel S-0291 now, as taken over by S-0340. T-1357 then only records it. Otherwise its retry tasks, T-1338 to T-1342, would be built for a step that will no longer reach GitHub. Criterion 6 leaves this to you.
3. **The hold.** S-0232, in progress and waiting on your release key (TH-0324), holds this story on `scripts/flai-snapshot.sh`, `flai/.goreleaser.yaml`, and `design/tech/ci.md`. Recommended: keep criterion 1 as written and wait for S-0232. Alternative: change criterion 1 so the server packages a plain `go build` for the host itself, without `flai-snapshot.sh` or GoReleaser. That frees the story from the hold and from GoReleaser in CI, but tests a hand-made archive rather than GoReleaser's.

## Figures

- **Forecast:** 1h30m, raised from flai's 21m. The story has six tasks, a new Go package, GoReleaser trial builds, and smoke runs with GitHub cut off. Delivery is about 10:30Z, later while S-0232 holds the story.
- **Cost of delay:** 800 USD a week, against flai's 37.50 from your input of 15m per cycle. Smoke runs in about 70 close-outs a week, and a local server saves each at least 3 minutes. The 24 failures a week cost about 2h more. Together that is about 5.5h a week. I recommend raising the input `time_lost_per_cycle` to 5h30m. The input is yours, so I left it as it is.

### 2026-10-08T08:07:52Z orchestrator (recommendation)
The plan is approved: six tasks in three layers, with touches that share no path within a layer. Putting the server in Go under `flai/` is right, since only `flai/**` selects smoke for the story's own close-out. The 1h30m forecast fits. Questions 1 to 3 are a decision not yet recorded and two scope changes, so they are the operator's. My recommendations:

1. **CI failure of the GitHub check: take the alternative.** CI fails visibly, and the script records the issue through `flai issue` when it runs on a host. ADR-0067 § Decision brings work to the remote only by publishing. A bot pull request that carries an issue file is a second way onto `main`, with machinery of its own to keep up. GitHub already notifies you of a failed scheduled run.
2. **S-0291: do not cancel it outright.** S-0340 takes smoke off GitHub, but real users' `install.sh` and `flai self-upgrade` still reach GitHub and still drop the 1 MB listing, which TH-0365 measured at 2 in 5. S-0340's own T-1352 check of the published install would flake the same way without S-0291's retries. I recommend keeping S-0291's T-1338 and T-1339 (the retries) and T-1341 (the docs), and dropping T-1340, since S-0340 rewrites `install-test.sh`. Let S-0340's T-1357 close I-0086. The two stories overlap on `install.sh`, `scripts/install-test.sh` and I-0086, so flai serializes them: S-0291 (1h) is ahead in the pull order.
3. **The hold on S-0232: your call, and it depends on your key.** S-0232 has waited on TH-0324 for your cosign key since 2026-10-07T22:24Z. If you will provide it soon, keep criterion 1 and test GoReleaser's real archive. If not, the alternative, a plain `go build` archive, frees S-0340 now.
4. **The cost of delay input:** set `time_lost_per_cycle` to `5h30m` on S-0340, as the planner recommends. flai refuses me inputs on a story in ready. Note that its 2h for failed smoke runs is the same time S-0291's recommended 2h counts. If you keep both stories, count it once.

Source: design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md § Decision
