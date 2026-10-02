---
title: Signing releases and verifying the peer
updated: 2026-10-02
status: active
topics: [cli, dashboard, release, security]
---

# Signing releases and verifying the peer

The finding of S-0193, for E-0015. The operator asked that the dashboard and the CLI be signed in CI with a private key, that the CLI verify a release with the public key before it installs it, and that each component use the same mechanism to refuse the other when it is not a signed release: flai closes its connection to an unsigned flaiover, and flaiover closes the connection from an unsigned flai. This document says how the two are released, upgraded, and connected today, what can be signed and how it can be verified, what a running process can and cannot prove about itself to a peer, and ends with one recommendation. Nothing here is built. The decision and the stories that follow are at the end.

Documentation was read on 2026-10-02. Nothing was tried on this host; every statement about tools comes from their documentation and every statement about this repository from its files.

## Today

### How flai is released

- `release-flai.yml` runs GoReleaser on a tag `flai/v*` with `--clean --skip=validate`, under `permissions: contents: write`, and publishes to the GitHub release of `bytepunx/system-flow` named `flai vX.Y.Z`.
- `flai/.goreleaser.yaml` builds `flai` for linux, darwin, and windows on amd64 and arm64 with `CGO_ENABLED=0`, `-trimpath`, and `-s -w`, and sets `buildinfo.Version`, `Commit`, and `Date` by `ldflags`. It archives each build as `flai_<version>_<os>_<arch>.tar.gz` (`.zip` on windows) and writes `checksums.txt`, SHA-256 of every archive. It has no `signs`, no `sboms`, and no `docker_signs` section: nothing is signed.
- `flai version` prints the three build values and the Go version; a build without `ldflags` falls back to `debug.ReadBuildInfo` for the commit (`flai/internal/buildinfo/buildinfo.go`). A build from this tree by `scripts/flai.sh` has no version.

### How flaiover is released

- `release-flaiover.yml` runs on a tag `flaiover/v*` and on every push to `main` that touches `flaiover/`, `flai/`, or itself, under `permissions: contents: read, packages: write`. It builds `flaiover/Dockerfile` with buildx for linux/amd64 and linux/arm64 and pushes to `ghcr.io/bytepunx/flaiover` with the tags `docker/metadata-action` gives it: `latest` on `main`, the version and its major on a tag, and `sha-<short>` always. It uses `docker/build-push-action` with its default `provenance`, so the registry holds a buildx provenance attestation it never asks anyone to check; there is no cosign step and no `actions/attest-build-provenance`.
- The image gets `FLAIOVER_VERSION` and `FLAIOVER_COMMIT` as build arguments from the newest `flaiover/v*` tag and the commit, kept as environment variables of the running container, and reports them as `flaiover_build_info` on `/metrics` (`flaiover/src/lib/server/metrics.ts`). The image carries no flai, no git, and nothing of any project ([ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)).
- An image built on this host with `flai dashboard --build` or `make flaiover-image` is `flaiover:local`, with the same build arguments taken from the newest tags (`flai/cmd/dashboard.go`).

### How the template is released

The template is a git repository, `bytepunx/system-flow-template`, versioned by `template/template.yaml`; `flai new` and `flai upgrade` clone it and record what they rendered in the lock file ([ADR-0015](../adrs/0015-template-lock-file.md)). It is not part of this finding: a clone over SSH or HTTPS to GitHub is authenticated by the transport, and nothing of it is executed.

### How a release is verified before it is installed

- `flai self-upgrade` (`flai/cmd/selfupgrade.go`, `flai/internal/selfupgrade/selfupgrade.go`) lists the releases through the GitHub API (`FLAI_RELEASES_API`, default `https://api.github.com`), takes the newest `flai/v*` that is neither a draft nor a pre-release, or the one `--version` names, downloads the archive for its platform and `checksums.txt` as release assets, checks the archive's SHA-256 against the line in `checksums.txt`, and replaces the binary atomically: a temporary file beside the target, mode 0755, then a rename, with the old binary moved aside on windows. A token from `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` is sent when the repository is private.
- `install.sh` does the same from a shell with `curl` and `sha256sum` or `shasum`, into `FLAI_INSTALL_DIR`, default `~/.flai/bin`.
- `flai host upgrade` runs the same upgrade for the host's flai and restarts `flai serve` and the MCP servers on it (`flai/cmd/host.go`).
- Both checks prove only that the download matches the `checksums.txt` published beside it. Whoever can write the release's assets, the owner of the repository, a holder of a token with `contents: write`, or a workflow that runs with it, can write both the archive and the checksums. The checksum stops a corrupt or truncated download, not a substituted release. The TLS connection to GitHub is the only thing that authenticates the publisher, and it authenticates GitHub, not the signer.

### How the dashboard's image is chosen and run

- `flai dashboard` resolves the image and tag from its flags, then the project's `system-flow.yaml`, then the host's configuration (`dashboard.image`, default `ghcr.io/bytepunx/flaiover`; `dashboard.tag`, default `latest`), pulls it for the daemon's platform unless it is present or `--pull` is given, and starts one container named `flaiover` for every project on the host ([ADR-0033](../adrs/0033-one-login-token-and-one-agent-credential-per-user-serve-every-project.md)). The pull is `docker pull --quiet <image>:<tag>`; the reference is a tag, never a digest, and nothing checks what the tag resolved to.
- `flai dashboard check` pulls the tag again and says whether it differs from what the running container was started from; `flai dashboard upgrade` does the same and, when it differs, starts the new image beside the old one, waits for it to answer `/_health`, and only then swaps. `flai host` restarts the container when it is recorded for restart ([ADR-0062](../adrs/0062-flai-host-restarts-the-dashboard-container-flai-dashboard-recorded-when-it-is.md)). None of them asks where the image came from.
- `flai dashboard --build` builds `flaiover:local` from this monorepo and runs that instead. It is how this repository runs its own dashboard.

### How the two authenticate each other on the channel

- flai dials the dashboard; the dashboard never dials the host ([ADR-0029](../adrs/0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). `flai serve` opens one WebSocket per registered project to `/agent` on the container's published port and speaks JSON-RPC 2.0 over it (`flai/internal/channel/channel.go`, `flaiover/src/lib/server/agent.ts`).
- Both sides hold one agent credential, a random secret `flai dashboard` generates once per user and keeps beside `flai serve`'s state, and gives the container as a file under `/run/secrets`, mode 0600, read-only ([ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md), [ADR-0033](../adrs/0033-one-login-token-and-one-agent-credential-per-user-serve-every-project.md)). The credential never crosses the connection. The handshake is `hello`, in which flai sends the protocol number, a nonce, its version (`flai`), the project, the methods it serves, and its kind; the dashboard answers with its own nonce, its version (`dashboard`), and an HMAC-SHA256 over its role and both nonces keyed with the credential; flai checks it and sends `hello.prove`, the same construction for its role, which the dashboard checks. A wrong proof is closed with `4401`, a malformed hello with `4400`, a slow one with `4408`; a second flai for a project a connected one still answers for is closed with `4409` and backs off ([ADR-0063](../adrs/0063-the-dashboard-refuses-a-second-flai-for-a-project-with-close-code-4409-while.md)).
- The dashboard compares the methods flai named with `REQUIRED_METHODS` and tells the designer which a project's flai lacks (S-0098). Nothing else is made of either version: the two numbers are shown, not checked.
- So the channel today authenticates the *installation*: the connection is between the flai and the flaiover the same `flai dashboard` set up, because only those two hold the credential. It says nothing about whether either of them is a build anyone released. A flai built from a branch, or a `flaiover:local` image, proves itself exactly as a release does.
