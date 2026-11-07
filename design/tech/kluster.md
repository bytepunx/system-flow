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
| Decision | [ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md), refined by [ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md); the design is [kluster.md](../system/kluster.md) |

## Why

A project such as authstar already describes its local cluster and shared services as a kluster config. flai runs that config rather than describing the cluster again, so the cluster a story is checked against is the one the project's own people run ([ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md)). flai shells out to `kluster` as it does to `git` and `docker` ([ADR-0010](../adrs/0010-shell-out-to-git-and-docker.md)), so a kluster upgrade on the host needs no flai release and flai takes no Kubernetes or Helm libraries into its module.

## The config file

What kluster v1.3.1 reads from its config, as S-0586 found it in kluster's source (`kluster/cmd/up.go`, `kluster/cmd/config.go`, and `kluster-lib/cluster/cluster.go`, at 6d7a87c). A flag given on the command line wins over the file, and the file over the flag's default.

| Key | What it is |
|-----|------------|
| `name` | The cluster's name. Required: from the file or `--name`, else `up` refuses |
| `provider` | `k3d`, the default, or `kind`. flai accepts `k3d` only |
| `profile` | The shared services `up` installs: `spire`, the default, `signet`, which adds Signet over `spire`, or `authstar`. An unknown profile is refused. A file with no profile, or an empty one, gets `spire`: kluster has no empty cluster |
| `addons` | Opt-in addons beside the profile, such as `argocd`, or the groups `observability` and `tracing` |

So the scaffold flai writes when the `kluster` host action is turned on names `profile: spire` out ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)), and the dashboard's Klusters page changes it before the first `up`.

`kluster setup` checks for Docker and `kubectl` and installs `kubectl` when it is missing. On Linux it installs it with `sudo` into `/usr/local/bin`, so flai serve, which cannot answer a password prompt, does not run it: the Klusters page names it for the operator to run on the host. It never installs Docker.

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
