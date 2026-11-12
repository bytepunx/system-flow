---
id: S-0740
type: story
nature: feature
title: flai image ensure brings back the harness images an installed flai needs, pulled by a signed digest or built on the host from the Dockerfiles flai carries
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T15:59:12Z
updated: 2026-11-12T15:59:57Z
transitions: []
tags: [cli]
topics: [docker, release, security]
touches: [flai/images/embed.go, flai/internal/agentimage/pull.go, flai/internal/agentimage/pull_test.go, flai/internal/agentimage/build.go, flai/internal/agentimage/build_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/internal/dashboard/digests.go, flai/internal/dashboard/digests_test.go, flai/internal/dashboard/resolve.go, flai/cmd/image.go, flai/cmd/image_test.go, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md]
after: [S-0733, S-0734, S-0735]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0740 flai image ensure brings back the harness images an installed flai needs, pulled by a signed digest or built on the host from the Dockerfiles flai carries

## Goal

When an agent image is missing today, the only way back is `make agent-images` in this repository, which builds `flai` and `pilot` from the tree. A project made from the template has no such target, and its installed flai cannot make the images.

After this story, `flai image ensure` makes each missing or stale harness image (bedrock and claude-code) for the flai that runs it, as S-0733's ADR decides. It pulls the image by the index digest the release's signed digest list names, verified as `flai dashboard` verifies flaiover's (S-0734 publishes it). Where it cannot pull, it builds the image on the host from the Dockerfiles carried in the flai binary, with that release's `flai` and `pilot` for Linux. `flai image build` forces the build. Each build step is limited in how long it may run with no output, and a step that fails on the network is tried again.

## Acceptance criteria

- [ ] `flai image ensure` pulls each missing or stale harness image by the digest that the signed list of the running flai's release names. It refuses a list that does not verify against the release key, and a digest the list does not name, with `flai dashboard`'s refusal kinds. It tags the pulled image with the name the start runs
- [ ] When the pull cannot run (no registry, no list, or a refusal the ADR lets fall back), it builds the image from the Dockerfiles embedded in the binary (`flai/images/embed.go`). It uses `flai` and `pilot` for Linux from the same release, checked against its signed `checksums.txt`, and says which way it took
- [ ] `flai image build` builds without trying the pull. In a development tree, as S-0733's ADR names one, both commands build from the tree's own Dockerfiles and source instead of a release's
- [ ] Each image it makes carries the label S-0735's status reads, and `flai image status` reports it current afterwards
- [ ] A build step that writes no output for the limit the ADR sets is stopped and reported, and a step that fails on the network (a reset, a timeout, a 5xx) is tried again up to the ADR's count. Each case is shown by a test over a fake docker
- [ ] The build's output streams to standard output as it runs, and to a log file under `.flai-cache` that the command names
- [ ] `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/runbooks/agent-container.md` § Build, and `design/system/flai-cli.md` describe both commands, and the runbook gives `flai image ensure` as the way back before `make agent-images`

## Tasks
- T-3168 The flai binary carries the bedrock and claude-code Dockerfiles through a go:embed package at flai/images
- T-3169 The agent images are pulled by the digest their release's signed list names, verified as flaiover's list is
- T-3170 flai builds the harness images on the host from its embedded Dockerfiles and its release's Linux binaries, with a limit on a silent step and a retry on the network
- T-3171 flai image ensure pulls or builds each missing or stale harness image, and flai image build builds without the pull
- T-3172 The user guide, the reference, the agent-container runbook, and flai-cli.md describe flai image ensure and build as the way back for a missing harness image

## Notes

Proposed slice 3b of E-0036. It builds only the harness images. The project's layer is the next story in the plan's order, since its build arguments need a manifest change of their own.
