---
title: kluster
updated: 2026-11-07
status: active
topics: [kluster, cli]
---

# kluster

| | |
|-|-|
| What | `kluster`, the command of [bytepunx/kluster](https://github.com/bytepunx/kluster): it makes a local Kubernetes cluster with k3d, embedded as a Go library, and installs a profile of shared services into it |
| Version | v1.3.1 tested (`kluster version v1.3.1` on the dev host). flai pins none: it runs the `kluster` on `PATH` |
| Install | `curl -fsSL https://raw.githubusercontent.com/bytepunx/kluster/main/install.sh \| bash`, which installs the release binary to `/usr/local/bin`, or `go install github.com/bytepunx/kluster/kluster@latest`. It needs Docker. Only a host that runs a project with a `kluster:` block needs it; system-flow itself runs no cluster |
| Used in | `flai/internal/kluster`, behind `flai kluster up`, `down`, and `status` (S-0577). `flai kluster deploy`, `teardown`, and `promote` (S-0578) and `flai serve`'s jobs behind the `kluster` host action (S-0579) build on it |
| Calls | `kluster --config <config> --provider <provider>` with `up`, `down --name <name> --yes`, `status`, and `kubeconfig --name <name> --output <file>`, run in the project's main checkout |
| Decision | [ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md); the design is [kluster.md](../system/kluster.md) |

## Why

A project such as authstar already describes its local cluster and shared services as a kluster config. flai runs that config rather than describing the cluster again, so the cluster a story is checked against is the one the project's own people run ([ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md)). flai shells out to `kluster` as it does to `git` and `docker` ([ADR-0010](../adrs/0010-shell-out-to-git-and-docker.md)), so a kluster upgrade on the host needs no flai release and flai takes no Kubernetes or Helm libraries into its module.

## What flai works around

kluster lacks what flai would read most easily; each workaround goes once kluster has the change ([ADR-0254 § Proposals for kluster's repository](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md#proposals-for-klusters-repository)).

| kluster lacks | flai does instead |
|---------------|-------------------|
| A JSON status | Parses `kluster status`'s table by its `NAME`, `RUNNING`, and `AGE` headers, passing over a column it does not know |
| An image import | Loads each image with `docker save` into `ctr images import` in each k3d node (S-0578) |
| A prompt-free default for `down` | Always passes `--yes`, and asks its own confirmation on a terminal |

## Considered

- `k3d` directly: the host has no `k3d` command, since kluster drives k3d as a library, and k3d alone installs none of the profile's shared services.
- kind, kluster's CI provider: left out until a project needs it; `flai check` reports a config whose provider is not `k3d`.
- A Go dependency on kluster's library: ties flai's release to kluster's and brings the k3d and Helm SDKs into flai's module.
