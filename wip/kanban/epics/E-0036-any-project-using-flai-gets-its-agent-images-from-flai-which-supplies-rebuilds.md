---
id: E-0036
type: epic
nature: improvement
title: Any project using flai gets its agent images from flai, which supplies, rebuilds, and keeps them current, and an operator can customize them or substitute their own machinery
status: in-progress
owner: alex
created: 2026-11-12T04:19:48Z
updated: 2026-11-13T07:28:50Z
transitions:
  - to: ready
    at: 2026-11-13T07:28:41Z
    by: alex
  - to: in-progress
    at: 2026-11-13T07:28:50Z
    by: agent-S-0733
tags: [cli]
topics: [release]
usage:
  source: sum
  seconds: 4894
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 12
      hand_edits: 1
      work: 129
  models:
    - model: claude-opus-5-5
      input: 660
      output: 133754
      cache_read: 48963539
      cache_write: 1646047
      cost: 22.5856
  strategic:
    - kind: planner
      seconds: 5945
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 907
          output: 108243
          cache_read: 177333136
          cache_write: 652148
          cost: 46.0268
    - kind: orchestrator
      seconds: 380
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 101
          output: 791
          cache_read: 3713714
          cache_write: 422922
          cost: 1.1853
        - model: claude-sonnet-5-5
          input: 8
          output: 49
          cache_read: 76668
          cache_write: 44576
          cost: 0.0931
---
# E-0036 Any project using flai gets its agent images from flai, which supplies, rebuilds, and keeps them current, and an operator can customize them or substitute their own machinery

## Outcome

A project made from the template runs its agents on the `container` harness with no image machinery of its own. flai supplies the harness images that match the installed flai, builds the project's layer over them when the manifest names one, and brings back an image that is missing or stale by itself, showing the work in the dashboard, instead of refusing the start with an error that tells the operator to run a script. The operator can still customize what the images hold, name an image of their own, or put their own build machinery in place of flai's.

## Stories
- S-0733 An ADR decides where an installed flai's agent images come from, how each is tied to the flai that runs it, and how flai brings back a missing one or hands the build to the operator's own machinery
- S-0734 Each flai release builds, signs, and publishes the bedrock and claude-code agent images with a signed digest list, and ships pilot for Linux
- S-0735 flai image status says which agent images a project's agents need, and whether each is on the host and current for the flai that runs it
- S-0740 flai image ensure brings back the harness images an installed flai needs, pulled by a signed digest or built on the host from the Dockerfiles flai carries
- S-0745 flai image ensure builds the project's agent.image.layer over the harness image itself, with build arguments the manifest gives in any project
- S-0751 flai serve brings a missing or stale agent image back as a host job that holds the start, in place of refusing it
- S-0757 The dashboard shows each agent image job as it runs, its output, and a failure's cause, and offers the operator a retry
- S-0759 An operator names a build command of their own that flai runs in place of its image build, or turns the automatic image job off and manages the images outside flai
- S-0760 This repository builds its agent images through flai's general path, with scripts/agent-image.sh as its own build command, and the documents and the template say how any project gets its images

## Notes

### Why

On 2026-11-11 every agent image was gone from the operator's host, and starting the planner failed with "the image flai-agent-sf-claude-code:local ... is not on this host ... build it with make agent-images". Two gaps are behind that message.

1. A missing image is only reported. `containerImage` in `flai/internal/serve/container.go` refuses the start, and `flai check` reports `agent.image` (`flai/internal/check/image.go`). Nothing rebuilds, whatever removed the images (a prune, a Docker Desktop reset, a new host).
2. Only this repository can build the images. `scripts/agent-image.sh` and `make agent-images` build `flai` and `pilot` from this tree with the Dockerfiles under `flai/images/`. The template ships neither, a released flai has no `pilot` (ADR-0212, consequences), and nothing publishes the images. In any other project the error's advice cannot be followed.

The design did not intend the second gap. [ADR-0203](../../../design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) names a baseline image flai builds, signs, and publishes, verified by digest. [ADR-0212](../../../design/adrs/0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md), decision 8, built the images locally within E-0029 and left the published image "to a later story". E-0029 was accepted without that story being written. This epic is that work, widened to recovery and to substitution.

### Principles the operator set (2026-11-11)

- Defining, creating, and managing the images is flai's, not each project's.
- An operator can customize the images or substitute their own machinery for flai's.
- A missing image, whatever the cause, has an appropriate way back that is not only an error in the dashboard.

### Proposed slices, for the planner to draft and correct

1. An ADR that decides where the harness images come from for an installed flai, refining ADR-0212's decision 8. The candidates: images published and signed per flai release and pulled by digest, as `flai dashboard` already does for flaiover; or Dockerfiles carried in the flai binary and built on the host, which needs `flai` and `pilot` for Linux from the release; or both, pull first and build as the fallback. It also decides how an image is tied to the flai that runs it (the `dev.system-flow.flai` label exists already), and what a development tree such as this repository does, where the images must hold the tree's flai.
2. The release pipeline builds, signs, and publishes the bedrock and claude-code images with each flai release, if the ADR chooses publication.
3. A flai command that owns the images on a host, such as `flai agent image status`, `ensure`, and `build`: it says which images the project's agents need and whether each is present and current, fetches or builds the harness images, and builds the manifest's `agent.image.layer` over them itself. The layer's build arguments are today read by `scripts/agent-image.sh` from this repository's own files (`flai/go.mod`, `scripts/install-tools.sh`), so the manifest or the Dockerfile needs a general way to give them.
4. `flai serve` ensures the images in place of refusing. When a start finds an image missing or stale, serve runs the ensure as a host job, one at a time, holds the start until it ends, and shows it in the dashboard with its progress. When the build fails, the dashboard shows the cause and the build's output and offers a retry. The container never gets the Docker socket (ADR-0203), so this is serve's work on the host, not an agent's. An upgrade of flai makes the images stale, and the same path renews them.
5. Customization and substitution. What exists stays: `agent.image.layer` for a project's tools, and `flai serve agent harness container --image <inner>=<image>` for an image of the operator's own. Added: a build command the operator names, which flai runs in place of its own machinery for ensure, and a setting that turns the automatic ensure off for an operator who manages images outside flai. Where that command is named matters: a command in the manifest runs on the host with the operator's rights, so it is the host's configuration or a protected path (`flai/internal/protected` already lists `scripts/agent-image.sh`).
6. This repository moves onto the general path: `scripts/agent-image.sh` becomes its substitute build, from the tree, and `make agent-images` stays as its entry. The refusal text, the `agent.image` finding, `docs/operators/` (the settings table, the agent-container runbook, the line saying nothing publishes the images), `design/system/agent-adapters.md`, `design/system/project-manifest.md`, and `design/tech/docker.md` say what is true after each story, and the template's documents say how a new project gets its images.

### What the rebuild of 2026-11-11 showed

The rebuild with `make agent-images` took three runs. The first hung for four hours in the claude-code layer's `apk add` and `npm install` step with no output, until it was killed. The second built that layer in under two minutes and then failed in the project's layer, where `go install` of goreleaser got connection resets from `proxy.golang.org`. The third, `scripts/agent-image.sh --layer`, passed in seconds from the download cache. So the ensure in slice 4 needs a limit on how long a build step may run without progress, a retry of a step that fails on the network, and the build's output where the operator can read it while it runs.

### Open questions for planning

- Pull, build, or both (slice 1). Publication fits ADR-0203 and the dashboard's existing path, and makes recovery a pull of seconds rather than a build of minutes; a host with no access to the registry needs the build.
- Whether the automatic ensure in slice 4 needs the operator's consent the first time, since it uses the network and several minutes, or only reports that it runs.
- Whether the strategic agents, which never get the project's layer, wait for the layer's build or start as soon as the harness image is there.
- 2026-11-13T07:28:41Z: moved to ready: follows S-0733, which moved to ready
- 2026-11-13T07:28:50Z: moved to in-progress: follows S-0733, which moved to in-progress
