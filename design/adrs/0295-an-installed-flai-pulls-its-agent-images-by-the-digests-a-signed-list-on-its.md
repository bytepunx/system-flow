---
id: ADR-0295
title: "An installed flai pulls its agent images by the digests a signed list on its release names, builds them on the host from the Dockerfiles it carries when the pull fails, and keeps them current through a flai serve host job that holds every start needing them, or runs the operator's own build command in its place"
status: accepted
date: 2026-11-13
supersedes: []
superseded_by: []
refines: [ADR-0212]
topics: [docker, release, security, cli]
---

# ADR-0295 An installed flai pulls its agent images by the digests a signed list on its release names, builds them on the host from the Dockerfiles it carries when the pull fails, and keeps them current through a flai serve host job that holds every start needing them, or runs the operator's own build command in its place

## Context

[ADR-0212](0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md), decision 8, built the agent images on the operator's machine with `scripts/agent-image.sh` and left the signed, published image of [ADR-0203](0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) to a later story. That story was never written, so only this repository can build the images, and a released flai has no `pilot`. A missing image is only reported: `containerImage` in `flai/internal/serve/container.go` refuses the start, and `flai check` warns `agent.image`. E-0036 is the missing work, widened to recovery and to the operator's own machinery.

On 2026-11-11 every agent image was gone from the operator's host. The rebuild took three runs: one hung for four hours with no output, one failed on connection resets from `proxy.golang.org`, and the third passed in seconds from BuildKit's cache.

What exists today:

- `flai dashboard` pulls flaiover by the index digest a signed list names: `flaiover_<version>.digests` and `.digests.sig` on the `flaiover vX.Y.Z` release, checked by `dashboard.Verify` (`flai/internal/dashboard/digests.go`) ([flai-cli.md](../system/flai-cli.md#verifying-a-dashboard-release-before-it-is-run), [release-signing.md](../system/release-signing.md#the-digest-list)). `dashboard.Parse` refuses a list whose lines name more than one repository.
- Bedrock's Dockerfile compiles `flai` and `pilot` from the tree and labels the image `dev.system-flow.flai` with its `FLAI_VERSION` argument. `scripts/agent-image.sh` passes the newest `flai/v*` tag reachable from `HEAD` (`git describe --abbrev=0`) as that argument, so the label and the image's `buildinfo.Version` both name a release, `1.86.12` on this host today, while the tree's own flai reports `dev`. Every layer inherits it; nothing in Go reads it.
- A tree build reports the version `dev`, with its commit in a separate field; a GoReleaser snapshot reports `0.0.0-<short commit>`. `buildinfo.Development` is true for both.
- This repository's manifest names `flai.command: scripts/flai.sh` ([ADR-0201](0201-when-the-manifest-names-flai-command-task-done-s-check-step-runs-that-command-s.md)) and `agent.image.layer: flai/images/system-flow/Dockerfile`. The script reads the layer's build arguments from `flai/go.mod` and `scripts/install-tools.sh`.
- `startOf` in `flai/internal/serve/agents_start.go` takes the image from `containerImageHost`, which gives the project's layer to every start, so the strategic agents run in the layer today.

## Decision

An installed flai pulls its agent images by the digests a signed list on its release names, builds them on the host from the Dockerfiles it carries when the pull fails, and keeps them current through a flai serve host job that holds every start needing them, or runs the operator's own build command in its place.

1. **Source: pull first, build as the fallback.**
   - The registry is GHCR, the packages public, so a pull needs no login: `ghcr.io/bytepunx/flai-agent-bedrock` and `ghcr.io/bytepunx/flai-agent-claude-code`, each tagged with the flai version without its `v`, for linux/amd64 and linux/arm64, and run by the image index's digest.
   - Each `flai/vX.Y.Z` release publishes one digest list per image, signed with the release key ([ADR-0070](0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md)), on the GitHub release `flai vX.Y.Z`: `flai-agent-bedrock_<version>.digests` and `flai-agent-claude-code_<version>.digests`, each with its `.sig`. They are in flaiover's format and are verified by the same code. One list per image, because that format holds one repository.
   - flai pulls only a digest a verified list names (ADR-0203), and caches the lists beside the dashboard's, as `agent-image-digests/`.
   - The fallback, when the registry cannot be reached or the pull fails: flai builds on the host from the bedrock and claude-code Dockerfiles it carries in its binary. It takes `flai` and `pilot` for Linux from the same release's archives, checked against its signed `checksums.txt`, for the engine's architecture (`docker info`, amd64 or arm64).
   - Bedrock's Dockerfile gains a target that copies the two prebuilt binaries from the build context, beside its tree build stage. The release pipeline builds the published images from that target, so a pulled image and a built one come from one recipe. GoReleaser ships `pilot` for linux/amd64 and linux/arm64, named in the same `checksums.txt`.
   - A built image is tagged `flai-agent-<inner>:<version>`, with no registry, so it is never taken for a pulled one.

2. **Tie and staleness.** Three labels tie an image to what it was built from:

   | Label | Holds | On |
   |---|---|---|
   | `dev.system-flow.flai` | The version the image's flai reports, as today | bedrock, inherited by every layer |
   | `dev.system-flow.recipe` | `sha256:` over the Dockerfile's bytes, the target, and the build arguments as sorted `NAME=value` lines, leaving out `BASE` and `FLAI_VERSION` | each image, its own |
   | `dev.system-flow.base` | The image ID of the image it was built over | each layer over another |

   - A missing image is one not on the host.
   - A pulled image is current when one of its `RepoDigests` is on the list of the release it must hold (`List.Matches`).
   - A built image is stale when its flai label is not the flai it must hold, its recipe is not the one flai would build with now, or its base is not the image now current under it. Base images are pinned in the Dockerfiles, so a moved base is a changed recipe.
   - `BASE` and `FLAI_VERSION` are left out of the recipe because the base label and the flai label carry them, more exactly than a tag does.

3. **Development tree.**
   - The flai an image must hold is the project's flai: the version that the manifest's `flai.command` prints with `version --json` when the manifest names one, else the running flai's.
   - When the manifest names a command, flai serve also takes the recipes it expects from that command's `flai image status --json`, run in the main checkout, as ADR-0201 asks that command for the check.
   - A development flai (`buildinfo.Development`: `dev`, empty, or `0.0.0-…`) is never pulled for. Its images come from the operator's build command (point 7); with none named, the job fails, naming the setting. This repository names `scripts/agent-image.sh` as its own, which builds from the tree.
   - For a development flai, an image is current when its flai label is a development version and its recipe and base labels match. The commit is not compared, so a commit to the tree rebuilds nothing.
   - `scripts/agent-image.sh` labels and stamps the tree's own version, `dev`, in place of the newest tag, and adds `dev.system-flow.commit`, which `flai image status` shows.

4. **The layer's build arguments.**
   - The Dockerfile's `ARG` defaults are the pins, kept in the Dockerfile.
   - The manifest may override them with literal values in a map, `agent.image.args`, under `agent.image` beside `layer`.
   - flai always passes `BASE` itself: the harness image it builds over.
   - flai reads no other project file: no `go.mod`, no `install-tools.sh`. The build context stays the project's root.
   - The Dockerfile's defaults and the manifest's overrides are both in the recipe, so a change to either makes the layer stale.

5. **flai serve on a missing or stale image.**
   - flai serve runs a host job, `image`, one at a time per host, since builds compete for the engine and the network. A need that arises while a job runs waits for it, or joins it when the job is for the same images.
   - The job pulls or builds, and holds every start that needs its images: a story's agent, a strategic agent, the close-out's container, and the checks' containers. When the job ends, they start, or fail with the job's cause.
   - Its record and output go to `image-jobs/<id>.job.json` and `<id>.log` under flai serve's state folder, beside `dashboard-digests/`, as the close-out job keeps `S-nnnn.job.json` and `S-nnnn.log` under `.flai-cache/verify/jobs/`.
   - The dashboard shows the job as it runs, its output, and a failure's cause, and offers a retry.
   - No consent is asked, the first time or after. The job is journalled and shown as it runs, and point 7's setting turns it off.

   | Failure | What the job does |
   |---|---|
   | A step with no output for `agent.image_stall`, a host setting, 10 minutes by default | Stops the step, and fails naming it |
   | A pull or a build whose last output names a connection reset, a timeout, a DNS failure, a TLS handshake failure, or a registry 5xx | Tries again, up to twice more, after 30 seconds and then 2 minutes; BuildKit's cache makes the steps that passed free |
   | Anything else | Fails at once |

   An upgrade of flai changes the version the images must hold, so they are stale, and the same path renews them.

6. **Strategic agents.**
   - The planner, the orchestrator, the analyzer, the curator, and the roles a project defines run on the harness image alone, never the project's layer. They start once the harness image is current, without waiting for the layer's build.
   - They run flai and read files, not the project's test tiers. This changes today's behaviour, where `containerImageHost` gives them the layer. S-0751 builds the change.

7. **Substitution.**
   - The operator names a build command of their own in the host configuration, per project: `flai serve agent image command [--project <root>] -- <program> <args>…`, stored as `agent.projects.<root>.image.command` ([ADR-0261](0261-a-story-or-strategic-agent-that-names-no-harness-falls-back-after-the-project-s.md)).
   - It is never named in the manifest: a manifest command would run on the host with the operator's rights, from whatever a branch says.
   - `flai guard` refuses both settings of this point to every served agent, as it refuses every write of the host configuration ([ADR-0220](0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md)).
   - flai runs the command in the project's main checkout in place of its own pull and build, then checks the images and their labels. A non-zero exit, or an image that is missing or carries other labels, fails the job.

     | Variable | Holds |
     |---|---|
     | `FLAI_IMAGE_VERSION` | The flai version the images must hold |
     | `FLAI_IMAGE_INNER` | The inner harness, `claude-code` |
     | `FLAI_IMAGE_BEDROCK`, `FLAI_IMAGE_HARNESS`, `FLAI_IMAGE_LAYER` | The tags the images must have; the layer's empty when the manifest names none |
     | `FLAI_IMAGE_RECIPE_BEDROCK`, `FLAI_IMAGE_RECIPE_HARNESS`, `FLAI_IMAGE_RECIPE_LAYER` | The recipe label each must carry |
     | `FLAI_IMAGE_STALE` | The images found missing or stale, of `bedrock`, `harness`, and `layer`, separated by spaces |

   - A second setting, `flai serve agent image auto off [--project <root>]`, turns the job off for the project. A start then refuses as today, naming the image and `flai image ensure`.
   - `flai serve agent harness container --image <inner>=<image>` stays, for an image of the operator's own. flai never pulls it, builds it, or judges it stale. A project's layer builds over it when it is set.

The rest of ADR-0212 and ADR-0203 stands: the layers, pilot, the mounts, the MCP server over HTTP, the guard, and secrets by name.

## Consequences

- Any project gets its images from flai, with no script of its own; a pull takes seconds where a build takes minutes.
- Which story builds each part:

  | Story | Builds |
  |---|---|
  | S-0734 | The release pipeline: both images from bedrock's prebuilt target, the two signed digest lists, and `pilot` for Linux in GoReleaser |
  | S-0735 | `flai image status`: which images a project's agents need, and whether each is on the host and current |
  | S-0740 | `flai image ensure` for the harness images: the embedded Dockerfiles, the pull by a verified digest, the fallback build, and the labels |
  | S-0745 | The layer built by flai over the harness image, with `agent.image.args` |
  | S-0751 | The `image` job in flai serve that holds the starts, its stall limit and retries, and the strategic agents on the harness image |
  | S-0757 | The dashboard's view of the job, its output, its cause, and the retry |
  | S-0759 | The operator's build command, the `auto` setting, and the guard's refusal of both |
  | S-0760 | This repository on the general path, with `scripts/agent-image.sh` as its build command, and the documents and the template |

- For a development flai, a change to `flai` or `pilot` in the tree reaches the images only when a recipe changes or the operator runs `make agent-images`. Until then a container's flai, which runs the guard and the verify inside, can be older than the tree's. `flai image status` shows the commit it was built from, so the gap is visible.
- The first start after a flai upgrade waits for a pull, or a build.
- The strategic agents lose the project's tools; a role that needs a test tier runs it on the host.
- `harness.DefaultImages`' `flai-agent-claude-code:local` gives way to the image the job resolved. The `agent.image` finding and the refusal name `flai image ensure` in place of `make agent-images`.
- `flai/internal/dashboard`'s list code is generalised past `flaiover_` names, so one parser and one verifier serve both.

Security:

- A pulled image is verified by a signed digest, as flaiover's is. The GHCR packages must stay public for an anonymous pull.
- A host build is unsigned, but it comes from Dockerfiles in a signed flai binary and release binaries checked against signed checksums, with the base images pinned in those Dockerfiles.
- The project's layer stays the project's risk, as ADR-0203 has it. The manifest gives it literal arguments only, never a command.
- An operator's command, and what it builds, are the operator's own, trusted as their checkout is. It lives where no agent writes.
- The container still never gets the Docker socket. The job is flai serve's, on the host.

## Alternatives considered

- **Pull only.** One path, and every image a signed one. But a host that cannot reach the registry, or a registry outage, would leave agents with no way back, which E-0036 exists to remove.
- **Build only.** No publication to keep. But every new host, prune, and upgrade costs minutes of building, with the network failures of 2026-11-11 each time, and no image would carry a signature.
- **Build from source on the host.** The tree's build stage, anywhere. But it needs Go and the module's source, which an installed flai lacks, and builds code no release signed.
- **Keep refusing, and document the build.** Nothing to build. But the refusal's advice can be followed only in this repository.
