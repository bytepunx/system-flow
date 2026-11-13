---
id: S-0740
type: story
nature: feature
title: flai image ensure brings back the harness images an installed flai needs, pulled by a signed digest or built on the host from the Dockerfiles flai carries
status: in-progress
parent: E-0036
owner: alex
created: 2026-11-12T15:59:12Z
updated: 2026-11-13T11:44:32Z
transitions:
  - to: ready
    at: 2026-11-13T07:29:14Z
    by: alex
  - to: in-progress
    at: 2026-11-13T11:39:55Z
    by: agent-S-0740
tags: [cli]
topics: [docker, release, security]
touches: [flai/images/embed.go, flai/internal/agentimage/pull.go, flai/internal/agentimage/pull_test.go, flai/internal/agentimage/build.go, flai/internal/agentimage/build_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/internal/dashboard/digests.go, flai/internal/dashboard/digests_test.go, flai/internal/dashboard/resolve.go, flai/cmd/image.go, flai/cmd/image_test.go, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, flai/images/embed_test.go]
after: [S-0734, S-0735]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 280
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 19
  models:
    - model: claude-opus-5-5
      input: 46
      output: 344
      cache_read: 3612303
      cache_write: 265837
      cost: 1.447
  strategic:
    - kind: orchestrator
      seconds: 31
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 35
          cache_read: 349274
          cache_write: 29200
          cost: 0.1088
forecast:
  duration: 1h
  delivery: 2026-11-13T13:32:00Z
  basis: "Its own forecast of 1h; 6th in the pull order with an in-progress limit of 3, behind S-0763, S-0760, S-0745, S-0759, S-0751 and S-0757."
  by: flai
  at: 2026-11-13T11:39:47Z
finalized:
  by: alex
  at: 2026-11-13T07:29:13Z
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
- T-3170 flai builds the harness images on the host from its embedded Dockerfiles and its release's Linux binaries checked against the signed checksums, or from a development tree's own
- T-3171 flai image ensure pulls or builds each missing or stale harness image, and flai image build builds without the pull
- T-3172 The user guide, the reference, the agent-container runbook, and flai-cli.md describe flai image ensure and build as the way back for a missing harness image
- T-3287 A docker step of a pull or a build streams its output, is stopped after the stall limit, and is tried again on a network failure

## Notes

Proposed slice 3b of E-0036. It builds only the harness images. The project's layer is S-0745, since its build arguments need a manifest change of their own.

A story's container has no Docker socket (ADR-0203), so its agent tests over a fake docker, and a real pull and build are the operator's to run on the host once the story is in review.

### Planning

- Touches:
  - `flai/images/embed.go`: design. The second criterion names it; `go:embed` cannot reach a parent folder, so the package sits beside the Dockerfiles. `flai/images` is protected for a container's agent (`flai/internal/protected/protected.go`, lines 75 to 78).
  - `flai/internal/agentimage/pull.go`, `build.go`, `ensure.go`, and their tests: new files in S-0735's package.
  - `flai/internal/dashboard/digests.go`, `digests_test.go`, `resolve.go`: layout. The signed digest list's parser and resolver live there (`Parse` at line 92, `Verify` at 125, `Resolve` at 94, `FetchList` at 182), and the pull reuses them.
  - `flai/cmd/image.go`, `image_test.go`: S-0735's command gains `ensure` and `build`.
  - `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/runbooks/agent-container.md`, `design/system/flai-cli.md`: named in the last criterion.
  - Left out from `flai touches suggest`: `docs/operators/index.md` (23%), `docs/operators/settings.md` (22%), and the dashboard's documents. No setting or page changes here.
- `after`:
  - S-0733 is not named, though this story reads the ADR's choice of pull and build, the step limit, the retry count, and what a development tree runs: it waits for S-0733 through S-0735, which names it.
  - S-0734: the pull reads the digest list's format and the release's Linux `pilot`, which that story publishes.
  - S-0735: `ensure` reads `flai/internal/agentimage`'s status and extends `flai/cmd/image.go`.
- Tasks, in three layers: T-3168 (embed) and T-3169 (pull) share no file and run together; T-3170 (build) waits for T-3168's embedded Dockerfiles; T-3171 chooses between T-3169 and T-3170; T-3172 describes it.
- Forecast 1h: `flai forecast` gives 42m (109 s per unit over 118 done large feature stories, size 23). Raised: two new mechanisms, a verified pull and a host build with a silent-step limit and network retries, each tested over a fake docker. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
