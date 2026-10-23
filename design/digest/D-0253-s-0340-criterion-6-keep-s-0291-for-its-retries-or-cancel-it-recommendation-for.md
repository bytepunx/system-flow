---
id: D-0253
title: "S-0340 criterion 6: keep S-0291 for its retries, or cancel it: Recommendation for the operator: keep S-0291 for its retries (T-1338, T-1339, T-1341) and drop T-1340, as the agent recommends"
kind: answer
sources: [TH-0372]
topics: [ci, testing, release, flai, cli, go, code, all]
paths: [scripts/install-test.sh, scripts/smoke.sh, scripts/flai-snapshot.sh, scripts/release-server.sh, install.sh, ".github/workflows/system-flow-check.yml", design/tech/ci.md, design/system/devex.md, docs/users/flai.md, flai/.goreleaser.yaml, flai/internal/releaseserver/releaseserver.go, flai/internal/releaseserver/releaseserver_test.go, flai/internal/releaseserver/serve/main.go, scripts/install-published-test.sh, ".github/workflows/install-published.yml", Makefile, scripts/README.md, design/system/flai-cli.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0253 S-0340 criterion 6: keep S-0291 for its retries, or cancel it: Recommendation for the operator: keep S-0291 for its retries (T-1338, T-1339, T-1341) and drop T-1340, as the agent recommends

The operator, alex, answered TH-0372 at 2026-10-08T08:18:49Z.

What was asked:

```text
Criterion 6 of S-0340 asks you whether S-0291 is cancelled or closed as taken over.
```

The answer:

```text
Recommendation for the operator: keep S-0291 for its retries (T-1338, T-1339, T-1341) and drop T-1340, as the agent recommends. I recommended the same on TH-0370. Users' `install.sh` and `flai self-upgrade` still reach GitHub, and S-0340's new published-install check would flake without the retries. S-0340 closes I-0086. Criterion 6 leaves this to you. One answer here also settles TH-0370's question 2.

Source: wip/kanban/stories/S-0340-the-smoke-tier-installs-and-self-upgrades-from-a-local-release-server-built-from-the-tree-so-no-close-out-depends-on-github.md
```

Source: TH-0372.
