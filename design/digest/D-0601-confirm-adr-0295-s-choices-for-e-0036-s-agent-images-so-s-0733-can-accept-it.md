---
id: D-0601
title: "Confirm ADR-0295's choices for E-0036's agent images, so S-0733 can accept it: Recommended for the operator: answer \"as recommended\" on all six. agent-S-0733 then accepts ADR-0295 and moves S-0733 to…"
kind: answer
sources: [TH-1139, S-0733]
topics: [release, docker, security, flai, cli, go, code, all]
paths: [design/adrs/README.md, design/system/flai-cli.md, design/adrs/0295-an-installed-flai-pulls-its-agent-images-by-the-digests-a-signed-list-on-its.md]
date: 2026-11-13
status: absorbed
absorbed_by: "design/adrs/0295-an-installed-flai-pulls-its-agent-images-by-the-digests-a-signed-list-on-its.md#decision"
asked_again: []
cited: []
---

# D-0601 Confirm ADR-0295's choices for E-0036's agent images, so S-0733 can accept it: Recommended for the operator: answer "as recommended" on all six. agent-S-0733 then accepts ADR-0295 and moves S-0733 to…

The operator, alex, answered TH-1139 at 2026-11-13T07:49:55Z.

What was asked:

```text
**Recommended:** answer "as recommended". I then accept ADR-0295 with `flai adr accept` and move S-0733 to review. If you change one, name it, and I amend the ADR and flai-cli.md first.
```

The answer:

```text
**Recommended for the operator:** answer "as recommended" on all six. agent-S-0733 then accepts ADR-0295 and moves S-0733 to review.

## Why

- Choice 1 does what ADR-0203 § Decision already asked for: `ghcr.io/bytepunx/...`, signed and listed as flaiover's image is, and run only by a digest a signed release names. The fallback build comes from a signed flai binary, and from release binaries checked against signed checksums. So no unsigned source comes in.
- Choice 3 removes the layer's build from the strategic agents' start path. The planner, orchestrator, analyzer, and curator run flai and read files. None of them runs a test tier.
- Choices 5 and 6 answer the failures of 2026-11-11: the four-hour hang and the `proxy.golang.org` resets. They also keep a host command out of the manifest, where a branch could change it.
- Choice 4 has a cost: a container's flai can lag the tree. `flai image status` shows that gap, and `make agent-images` closes it.

## Why it is yours alone

- Choice 2 lets flai serve pull and build images on your host with no consent asked. That is a change to the host.
- The GHCR packages must be public for an anonymous pull. Making them public is a forge setting outside the repository, and S-0734 needs it.
- flai serve does not confirm this after `orchestration.permissions.confirm_after`.

**Alternatives:** name the choice you would change, such as asking consent once per host (2), or keeping the layer for strategic agents (3). agent-S-0733 then amends the ADR and `design/system/flai-cli.md` first.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md § Decision
```

Source: TH-1139.

S-0733's narrative (`wip/archive/agents/S-0733.md`, Decisions) records the outcome: alex confirmed all six of ADR-0295's choices as recommended, and ADR-0295 was accepted as written. The six: pull first, with the build as the fallback; no consent asked; strategic agents on the harness image alone; a development tree's images from the operator's build command, its commit not compared; a 10-minute stall limit and two network retries; the build command in the host configuration only, and the layer's build arguments as the Dockerfile's `ARG` defaults overridden by `agent.image.args`.

Source: S-0733.
