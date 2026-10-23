---
id: D-0254
title: "S-0340 plan: six tasks in three layers; three questions: don't cut a PR, we can leave it without an issue filing for now"
kind: recommendation
sources: [TH-0370]
topics: [ci, testing, release, flai, cli, go, code, all]
paths: [scripts/install-test.sh, scripts/smoke.sh, scripts/flai-snapshot.sh, scripts/release-server.sh, install.sh, ".github/workflows/system-flow-check.yml", design/tech/ci.md, design/system/devex.md, docs/users/flai.md, flai/.goreleaser.yaml, flai/internal/releaseserver/releaseserver.go, flai/internal/releaseserver/releaseserver_test.go, flai/internal/releaseserver/serve/main.go, scripts/install-published-test.sh, ".github/workflows/install-published.yml", Makefile, scripts/README.md, design/system/flai-cli.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0254 S-0340 plan: six tasks in three layers; three questions: don't cut a PR, we can leave it without an issue filing for now

The operator, alex, answered TH-0370 at 2026-10-08T08:19:40Z.

What was asked:

```text
The plan for S-0340: six tasks in three layers. Touches, forecast, and cost of delay are set; the reasons are in the story's Notes under `### Planning`.
```

The recommendation it overruled:

```text
The plan is approved: six tasks in three layers, with touches that share no path within a layer. Putting the server in Go under `flai/` is right, since only `flai/**` selects smoke for the story's own close-out. The 1h30m forecast fits. Questions 1 to 3 are a decision not yet recorded and two scope changes, so they are the operator's. My recommendations:

1. **CI failure of the GitHub check: take the alternative.** CI fails visibly, and the script records the issue through `flai issue` when it runs on a host. ADR-0067 § Decision brings work to the remote only by publishing. A bot pull request that carries an issue file is a second way onto `main`, with machinery of its own to keep up. GitHub already notifies you of a failed scheduled run.
2. **S-0291: do not cancel it outright.** S-0340 takes smoke off GitHub, but real users' `install.sh` and `flai self-upgrade` still reach GitHub and still drop the 1 MB listing, which TH-0365 measured at 2 in 5. S-0340's own T-1352 check of the published install would flake the same way without S-0291's retries. I recommend keeping S-0291's T-1338 and T-1339 (the retries) and T-1341 (the docs), and dropping T-1340, since S-0340 rewrites `install-test.sh`. Let S-0340's T-1357 close I-0086. The two stories overlap on `install.sh`, `scripts/install-test.sh` and I-0086, so flai serializes them: S-0291 (1h) is ahead in the pull order.
3. **The hold on S-0232: your call, and it depends on your key.** S-0232 has waited on TH-0324 for your cosign key since 2026-10-07T22:24Z. If you will provide it soon, keep criterion 1 and test GoReleaser's real archive. If not, the alternative, a plain `go build` archive, frees S-0340 now.
4. **The cost of delay input:** set `time_lost_per_cycle` to `5h30m` on S-0340, as the planner recommends. flai refuses me inputs on a story in ready. Note that its 2h for failed smoke runs is the same time S-0291's recommended 2h counts. If you keep both stories, count it once.

Source: design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md § Decision
```

The operator's answer:

```text
1. don't cut a PR, we can leave it without an issue filing for now
```

Source: TH-0370.
