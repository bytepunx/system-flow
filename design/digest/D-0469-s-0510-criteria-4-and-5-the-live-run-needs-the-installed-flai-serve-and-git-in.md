---
id: D-0469
title: "S-0510 criteria 4 and 5: the live run needs the installed flai serve and git in the container, so move it to S-0565: Recommended: say yes, as you answered for S-0508 on TH-0812 (D-0466) and for…"
kind: answer
sources: [TH-0814]
topics: [harness, docker, manifest, flai, cli, go, template, conventions, code, all]
paths: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/image.go, flai/internal/check/image_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, scripts/agent-image.sh, Makefile, scripts/README.md, flai/images/system-flow/Dockerfile, system-flow.yaml, design/system/project-manifest.md, docs/operators/settings.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/summary.md, design/issues, flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/store_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/cmd/check.go, flai/cmd/check_stats_test.go, flai/internal/check/story.go, flai/internal/verify/story.go, design/issues/I-0255-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md, scripts/env.sh, scripts/flai-build.sh]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0469 S-0510 criteria 4 and 5: the live run needs the installed flai serve and git in the container, so move it to S-0565: Recommended: say yes, as you answered for S-0508 on TH-0812 (D-0466) and for…

The operator, alex, answered TH-0814 at 2026-11-02T23:25:31Z.

What was asked:

```text
**Recommended:** answer "as recommended". Criterion 5 moves to S-0565, and criterion 4's close-out part is proven there too. S-0510 proves criterion 4 with a container run by hand of the layered image on its own worktree.
```

The answer:

```text
**Recommended:** say yes, as you answered for S-0508 on TH-0812 (D-0466) and for S-0423 on TH-0803 (D-0461).

- S-0510's agent removes criterion 5 from S-0510, and narrows criterion 4 to what a hand run of the layered image on its own worktree shows: each tier `flai test` runs inside it, recorded in the story's notes.
- It adds criterion 5's words and criterion 4's close-out half to S-0565, which is still a draft, and adds `S-0510` and `S-0517` to S-0565's `after`. S-0565's Notes already say its proof needs the layer (S-0510) and git in the container (S-0517).
- S-0510's Notes record the move.

**Why:** the same case as D-0466. The layer exists only on `story/S-0510`, and this host's `flai serve` is the installed release. The close-out also needs I-0215's fix, S-0517, still in progress. So S-0510 could not meet the definition of done with criterion 5 kept.

**Why I recommend rather than answer:** this takes a criterion out of a finalized story and widens another. That changes scope, so it is your call. Finalizing S-0565 stays yours.

**Alternative:** keep criterion 5. I don't recommend it: S-0510 would wait in review until S-0517 lands and S-0510 is published.

Source: design/digest/D-0466-s-0508-the-live-container-run-criterion-6-t-2169-s-first-half-cannot-run-before.md
```

Source: TH-0814.
