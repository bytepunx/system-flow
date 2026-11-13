---
id: TH-1139
title: Confirm ADR-0295's choices for E-0036's agent images, so S-0733 can accept it
anchor:
  path: wip/kanban/stories/S-0733-an-adr-decides-where-an-installed-flai-s-agent-images-come-from-how-each-is.md
  item: S-0733
status: resolved
participants: [agent-S-0733, orchestrator, alex]
created: 2026-11-13T07:39:08Z
updated: 2026-11-13T07:50:27Z
---

# TH-1139 Confirm ADR-0295's choices for E-0036's agent images, so S-0733 can accept it

On wip/kanban/stories/S-0733-an-adr-decides-where-an-installed-flai-s-agent-images-come-from-how-each-is.md.

## Entries

### 2026-11-13T07:39:08Z agent-S-0733
**Recommended:** answer "as recommended". I then accept ADR-0295 with `flai adr accept` and move S-0733 to review. If you change one, name it, and I amend the ADR and flai-cli.md first.

ADR-0295 is on `story/S-0733`, proposed, refining ADR-0212: `design/adrs/0295-an-installed-flai-pulls-its-agent-images-by-the-digests-a-signed-list-on-its.md`. `design/system/flai-cli.md` § Where the agent images come from (ADR-0295) describes it, with each part marked by the E-0036 story that builds it.

The choices that are yours:

1. **Pull first, build as the fallback.** Images: `ghcr.io/bytepunx/flai-agent-bedrock` and `ghcr.io/bytepunx/flai-agent-claude-code`, public on GHCR. They are pulled only by a digest from two signed lists on each `flai vX.Y.Z` release, `flai-agent-<image>_<version>.digests`, each with its `.sig`, in flaiover's format. One list per image, because `dashboard.Parse` holds one repository per list. The fallback build uses Dockerfiles embedded in flai, and `flai` and `pilot` for Linux from the release, checked against its signed `checksums.txt`. *Alternatives:* pull only, which fails with no registry; or build only, which takes minutes on every new host and hits the network failures of 2026-11-11.
2. **No consent asked** before the first automatic pull or build. The job is shown as it runs, and `flai serve agent image auto off` turns it off per project. *Alternative:* ask once per host on a thread.
3. **Strategic agents run on the harness image alone**, never the project's layer, and start once it is current. This changes today's behaviour: `containerImageHost` gives them the layer. They run flai and read files, not the test tiers. *Alternative:* keep the layer and wait for its build.
4. **A development tree.** A development flai (`dev`) is never pulled for. Its images come from the operator's build command, which for this repository is `scripts/agent-image.sh` (S-0760). They are labelled `dev` with a `dev.system-flow.commit` label. The commit is not compared, so a commit rebuilds nothing, and `make agent-images` renews them by hand. *Cost:* the flai inside the container can lag the tree until a recipe changes or you rebuild.
5. **Limits on the image job.** A step with no output for 10 minutes fails the job (`agent.image_stall`). A failure on the network is tried twice more, after 30 seconds and then 2 minutes.
6. **Where your own build command lives.** It goes in the host configuration only, `flai serve agent image command`, never the manifest, and `flai guard` refuses it to served agents. The layer's build arguments are the Dockerfile's `ARG` defaults, overridden by literal values in `agent.image.args` in the manifest.

Also yours, outside the ADR: the GHCR packages must be public for an anonymous pull. S-0734 publishes them.

### 2026-11-13T07:40:25Z orchestrator (recommendation)
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

### 2026-11-13T07:49:55Z alex
Confirmed the recommendation of 2026-11-13T07:40:25Z orchestrator.

Source: design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md § Decision

### 2026-11-13T07:50:27Z agent-S-0733
Resolved: alex confirmed all six choices as recommended; ADR-0295 accepted on story/S-0733 (2727377).
