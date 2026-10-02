---
title: Signing releases and verifying the peer
updated: 2026-10-02
status: active
topics: [cli, dashboard, release, security]
---

# Signing releases and verifying the peer

The finding of S-0193, for E-0015. The operator asked that the dashboard and the CLI be signed in CI with a private key, that the CLI verify a release with the public key before it installs it, and that each component use the same mechanism to refuse the other when it is not a signed release: flai closes its connection to an unsigned flaiover, and flaiover closes the connection from an unsigned flai. This document says how the two are released, upgraded, and connected today, what can be signed and how it can be verified, what a running process can and cannot prove about itself to a peer, and ends with one recommendation. Nothing here is built. The decision and the stories that follow are at the end.

Every statement about this repository comes from its files, read on 2026-10-02. Nothing was tried on this host, and no page was fetched in this session: what this document says about GoReleaser, cosign, minisign, GPG, Sigstore, and GitHub's attestations is from their documentation as the author knew it, and the story that adopts a tool checks it against the tool's current documentation first. Where a claim matters to the choice and could be wrong, it is marked **to check**.

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

## Signing releases

What the epic asks of a release is two things: a signature made in CI with a private key only CI holds, and a check with the public key wherever the release is consumed, before it is installed or run. The consumers are `flai self-upgrade` and `flai host upgrade` (Go, inside flai), `install.sh` (a shell with `curl` and whatever else a fresh machine has), and `flai dashboard` with its `check` and `upgrade` (Go, with `docker`). The dashboard's own code never consumes a release; it is the thing released.

Two properties decide most of the choice. First, **what the consumer needs to verify**: a public key in the binary and the standard library, or a library, or a tool on the path, or the network. flai is a single static binary installed by a shell script on machines that have nothing else; a verification that needs cosign installed, or the Sigstore trust root fetched, is a weaker promise than one that needs nothing. Second, **where the private key lives**, because the key, not the format, is what a compromise is about.

### What is signed

For flai, the right thing to sign is `checksums.txt`, not each archive: it already names every archive with its SHA-256, GoReleaser already publishes it, and `self-upgrade` and `install.sh` already download it and check the archive against it. One signature over that file turns the check they make today from "the archive matches a file published beside it" into "the archive matches a file the release key signed". GoReleaser's `signs` section does exactly this with `artifacts: checksum`, running any command with the checksum file as its input and uploading what it writes as a release asset. The archives themselves are then covered through their hashes and need no signature of their own.

For flaiover, the artifact is an image in GHCR and the thing that identifies it is its digest, `sha256:…`, which Docker checks on pull against the content it receives. Signing the digest signs the image. There are two ways to publish that signature, under [How the image's signature is published](#how-the-images-signature-is-published).

### The ways to sign

#### Ed25519 in the minisign format

[minisign](https://jedisct1.github.io/minisign/) is a small signing tool with a fixed format: Ed25519 keys, a public key that is one base64 line carrying a key ID, and a signature file of four lines, an untrusted comment, the base64 signature, a trusted comment, and a global signature over both. The private key is a file protected by a password. GoReleaser runs it through `signs` with `cmd: minisign`.

- **Verifying in flai**: `crypto/ed25519` in the standard library, after decoding the signature line and checking the key ID. No dependency; `aead.dev/minisign` is a small pure-Go library if the format's edge cases are not worth owning. **To check**: whether the default signature is over the file or its BLAKE2b-512 prehash (minisign signs the prehash with `-H`, and newer versions by default), because a prehashed signature needs `golang.org/x/crypto/blake2b`, a dependency, though a well-known one.
- **Verifying in `install.sh`**: needs `minisign` on the path, which almost no machine has, or OpenSSL 3 with `pkeyutl -rawin` for Ed25519, which macOS does not ship (its `openssl` is LibreSSL; **to check** whether its version verifies Ed25519). So the script would most likely have to say "verified only where minisign is installed", which is the current state with a different wording.
- **Verifying in flaiover** (for [Verifying the peer](#verifying-the-peer)): Node's `crypto.verify` with an Ed25519 key, standard library.
- **Key**: a password-protected file. In CI it is a GitHub Actions secret written to disk for the step, with the password in another secret.

#### cosign with a key pair

[cosign](https://docs.sigstore.dev/cosign/) is Sigstore's signing tool. With `cosign generate-key-pair` it makes an ECDSA P-256 key pair, the private key encrypted with a password (`cosign.key`), the public key as PEM (`cosign.pub`). `cosign sign-blob --key cosign.key --output-signature checksums.txt.sig checksums.txt` writes the signature, which is the base64 of a DER-encoded ECDSA signature over the SHA-256 of the file. By default it also records the signature in Rekor, Sigstore's public transparency log; `--tlog-upload=false` turns that off. The key can be given as `--key env://COSIGN_PRIVATE_KEY` with the password in `COSIGN_PASSWORD`, which is how CI uses it; `cosign generate-key-pair github://owner/repo` stores a fresh pair straight into the repository's Actions secrets. GoReleaser's documented `signs` example is cosign, with `artifacts: checksum`, and it also signs images with `docker_signs`.

- **Verifying in flai**: `encoding/pem`, `crypto/x509.ParsePKIXPublicKey`, `crypto/sha256`, and `crypto/ecdsa.VerifyASN1`, all standard library, with the public key as a PEM string in the source. No dependency, no network beyond the download itself. The Rekor entry, when one is made, is not needed to verify and is not consulted.
- **Verifying in `install.sh`**: `openssl dgst -sha256 -verify cosign.pub -signature <(base64 -d checksums.txt.sig) checksums.txt`, which every OpenSSL and LibreSSL does. That is the one line that makes the script's check real on a fresh machine, and the reason this option beats minisign for this repository. Anyone with cosign can also run `cosign verify-blob --key cosign.pub --signature … checksums.txt`.
- **Verifying in flaiover**: Node's `crypto.verify('sha256', data, pem, der)`, standard library.
- **Key**: an encrypted PEM and its password, two GitHub Actions secrets. Rotation means a new pair and a release of flai that carries both public keys; see [Keys](#keys).

#### cosign keyless

The same tool, with no key of ours: the workflow authenticates to Sigstore's Fulcio with its GitHub Actions OIDC token and receives a short-lived certificate that names the workflow (`https://github.com/bytepunx/system-flow/.github/workflows/release-flai.yml@refs/tags/flai/v1.2.3`) and the issuer (`https://token.actions.githubusercontent.com`); the signature and the certificate are recorded in Rekor; verification checks the certificate chain to Sigstore's root, the identity and issuer against what the verifier expects, and the Rekor entry's inclusion. GoReleaser's cosign example covers this form too.

- **What it buys**: no private key to keep, rotate, or lose. The signer is the workflow, and a compromise of the repository that could run the workflow is the only way to sign.
- **What it costs the consumer**: the verifier needs Sigstore's trust root, fetched over TUF and refreshed, and the Rekor entry; in Go that is `sigstore-go`, a large dependency tree for a CLI that is otherwise standard library, and verification needs the network, which the download needs anyway, but to a second party, Sigstore's public instance, whose availability becomes flai's. `install.sh` would need cosign installed. This is heavy for what flai is.
- **Where it fits**: as a second signature anyone can check with public tooling, beside a key the binaries verify themselves. Not as the only one.

#### GitHub artifact attestations

`actions/attest-build-provenance` signs a SLSA provenance statement for an artifact with the workflow's identity, in Sigstore's bundle format, and stores it with the repository, where `gh attestation verify <file> --owner bytepunx` checks it. It is keyless signing with GitHub as the trust root (public repositories use Sigstore's public instance, private ones GitHub's own, **to check** on which plans). It describes *how the artifact was built*, not only that the key holder approved it, which is the stronger supply-chain statement.

- **Consumer**: `gh` or `sigstore-go`, and the network to GitHub's attestation API. The same objection as keyless cosign for flai and for `install.sh`, with GitHub as the only external party. It can attest the image too, by its digest, and the image's attestation can be verified from the image alone with `gh attestation verify oci://…`.
- **Where it fits**: a cheap addition to both release workflows for anyone with `gh`, and the right thing to offer auditors; not the mechanism flai verifies with.

#### GPG

GoReleaser's default `signs` command is `gpg --detach-sign`. The format is OpenPGP; the key lives in a keyring.

- **Consumer**: in Go, `golang.org/x/crypto/openpgp` is frozen and deprecated and the maintained fork is ProtonMail's, both large for one signature; in a shell, `gpg` is widely installed but needs the public key imported into a keyring first, and users are used to being told to trust a key ID. The format and tooling were made for people and email, not for a binary checking itself. Not recommended.

### Comparison

| | Ed25519, minisign | cosign, key pair | cosign, keyless | GitHub attestations | GPG |
|--|--|--|--|--|--|
| Signs in CI with | `minisign`, secret key file | `cosign`, secret key in env | `cosign`, OIDC identity | `actions/attest-build-provenance` | `gpg`, secret key in keyring |
| Private key | ours, a file | ours, a file | none, Sigstore issues per run | none, GitHub issues per run | ours, a keyring |
| flai verifies with | stdlib (`crypto/ed25519`), maybe blake2b | stdlib (`crypto/ecdsa`) | `sigstore-go`, TUF root, Rekor | `sigstore-go` or `gh`, GitHub API | OpenPGP library |
| `install.sh` verifies with | `minisign`, rarely present | `openssl dgst -verify`, present | `cosign`, rarely present | `gh`, often present | `gpg`, often present, needs import |
| Network to verify | none | none | Sigstore | GitHub | none |
| Rotation | new key, release carries both | new key, release carries both | nothing to rotate | nothing to rotate | new key, release carries both |
| Compromise means | whoever holds the file and password signs | same | whoever can run the workflow signs | same | same as a file |
| Also signs the image | no, digest list only | `cosign sign` in GHCR, or digest list | `cosign sign` in GHCR | yes, by digest | no |

### How the image's signature is published

**A signed digest list as a release asset.** `release-flaiover.yml` already knows the digest it pushed (`docker/build-push-action` outputs it). On a `flaiover/v*` tag the workflow writes one file, `flaiover_<version>.digests`, with the digest for each platform and the tag it was pushed under, signs it with the same key and tool as `checksums.txt`, and uploads both to a GitHub release `flaiover vX.Y.Z`, which does not exist today and would be created by the workflow. `flai dashboard` then resolves the version it is to run through the same release client `self-upgrade` already has, downloads the list and its signature, verifies with the key in the binary, and pulls and runs `ghcr.io/bytepunx/flaiover@sha256:…` by digest, so that Docker's own content check completes the chain. One mechanism, one key, one client, no new dependency; the dashboard's version becomes something flai chooses from a signed list rather than a tag it trusts.

**cosign on the image.** `cosign sign --key cosign.key ghcr.io/bytepunx/flaiover@sha256:…` stores the signature in GHCR as an OCI artifact under a tag derived from the digest, and anyone verifies with `cosign verify --key cosign.pub ghcr.io/bytepunx/flaiover:1.2.3`. For flai to verify it without cosign, it would speak the registry's HTTP API (a token exchange with GHCR, a manifest fetch, a blob fetch) and check the ECDSA signature over the payload in the manifest's annotation: a few hundred lines of standard library against a layout cosign owns, or `go-containerregistry` plus `sigstore` as dependencies. It is the right thing to publish for the ecosystem and the wrong thing for flai to depend on.

**Notary (Docker Content Trust) and notation** sign tags in a registry-side trust server; GHCR does not run one. Not an option here.

A consequence either way: today `latest` is pushed on every merge to `main`, and `sha-*` tags with it. Those builds are not releases and would carry no signature. Once `flai dashboard` runs only what a signed list names, `dashboard.tag: latest` has to mean "the newest flaiover release", resolved from the releases, not the registry's `latest` tag; and an unreleased `main` build is run only as a development image, the way [Verifying the peer](#verifying-the-peer) describes.

### Keys

- **One key for both components**, held as two GitHub Actions secrets (the encrypted private key and its password) in the `bytepunx/system-flow` repository, used only by `release-flai.yml` and `release-flaiover.yml`. An Actions *environment* with required reviewers is the one setting that stops a workflow run from a branch using it: both release workflows run on tags, and a tag can be pushed by anyone with write access. **To check**: environment secrets and their reviewer rule on the repository's plan.
- **The public key is in both sources**, as a PEM constant in `flai/internal/buildinfo` (or a sibling package) and in `flaiover/src/lib/server`, with the key's SHA-256 fingerprint in `docs/operators` so an operator can compare, and published as a release asset too. It is not fetched from anywhere at verification time: a key fetched over the same channel as the artifact proves nothing.
- **Rotation**: generate the new pair, release a flai that carries both public keys, then switch the workflows to the new key; a flai older than that release cannot verify a newer release and `self-upgrade` tells the operator to install the bridging release first, by its version. A signature file is tried against every key the binary knows; nothing in the format needs a key ID.
- **Compromise**: rotate as above, and the release that carries the new key is signed with the old one, which an attacker also holds. So a compromise is announced in `CHANGELOG.md` and the release notes with the fingerprint of the new key, and `install.sh` and the documentation tell the operator to compare it. There is no stronger recovery without a second, offline root key; that is deliberate for a project of this size, and said here so that nobody expects otherwise.
- **Snapshot and local builds** (`scripts/flai.sh`, `goreleaser --snapshot`, `flai dashboard --build`) are not signed. They are development builds and are treated as such by the peer check.

### Recommendation

Sign `checksums.txt` and the image digest list with **cosign and a key pair** (ECDSA P-256), through GoReleaser's `signs` for flai and a step in `release-flaiover.yml` for the image, with `--tlog-upload=false` at first so that nothing depends on Rekor, and verify in flai with the standard library and in `install.sh` with `openssl`. Publish the image's digest list as a release asset of a `flaiover vX.Y.Z` release and make `flai dashboard` run by digest from it. Add `actions/attest-build-provenance` to both workflows as a second, keyless statement for anyone with `gh`, costing one step each and nothing in flai. Keep the template out: it is a clone.

The reasons, in order: the consumer verifies with nothing but what it already has (stdlib, `openssl`, `docker`), on every platform flai is built for, without a second service; the private key is the operator's, which is what the epic asked; cosign is the mainstream tool for exactly this and is GoReleaser's documented example, so the CI side is a few lines; the same key, format, and client cover both components, which is what lets [Verifying the peer](#verifying-the-peer) be one mechanism; and the keyless statement is kept for what it is good at, provenance for auditors, without making flai depend on it. minisign is the runner-up and would be chosen over cosign only if `install.sh` were allowed to skip verification where `minisign` is absent, which it should not be.

## Verifying the peer

The epic asks for the same mechanism to make each component refuse the other when it is not a signed release: flai closes its connection to an unsigned flaiover, flaiover closes the connection from an unsigned flai. This section says what that can and cannot mean, and what reaches as far as the facts allow.

### What a running process can prove about itself

A signature is a statement about bytes at rest: this archive, this image, was approved by the key holder. It is checked by whoever receives the bytes before running them, which is what [Signing releases](#signing-releases) arranges for `self-upgrade`, `install.sh`, and `flai dashboard`. A running process has no such thing to show a peer. Whatever it sends over the connection, a hash of its own executable, a version, a signed statement it was built with, is data it chose to send; a modified build can send the same data, copied from a real one. Proving to a remote party which code is running needs a measurement made by something the code does not control, such as a hardware root of trust, and nothing in flai's or flaiover's environment offers one. The design below is therefore honest about three levels:

| Level | Who checks | What it proves | What it does not |
|-------|-----------|----------------|------------------|
| Artifact verified before it runs | `self-upgrade`, `install.sh`, `flai dashboard` | the bytes installed or started are a signed release | nothing about what runs later |
| Peer's artifact measured by a trusted substrate | flai, through `docker inspect` | the container flai talks to runs an image the release key signed | nothing where flai cannot ask Docker |
| Release stamp exchanged on the channel | both, in `hello` | the peer carries a statement the release key signed, naming its component and version | that the peer *is* that release, against anyone who holds a release binary and the credential |

The threat each level stops: the first stops a substituted download or image, which is what signing is for. The second stops a dashboard container that was started from, or replaced by, an image nobody signed, including the `latest` tag moving to an unreleased `main` build. The third stops the honest mistakes, a flai built from a branch pointed at a production dashboard, a `flaiover:local` left running, a release of one talking to a snapshot of the other, and makes them visible by name. What none of them stops is an attacker who already holds the agent credential and can run code on the host or in the container: they can replay a real stamp. The credential is the boundary there, as [ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md) already says, and signing does not move it.

### The agent credential, as it is

The HMAC handshake already proves to each side that the other was set up by the same `flai dashboard`, because only those two hold the credential, which never crosses the wire. It is the thing that makes the channel the operator's and not anyone's, and it stays. It does not say whether either side is a release, and that is the gap the two designs below fill.

### flai measures the container before it dials

`flai serve` dials the container `flai dashboard` started on this host; it knows the container's name (`flaiover`) and can ask Docker what it runs: `docker inspect flaiover --format '{{.Image}}'` gives the image ID, and `docker image inspect` of that gives `RepoDigests`, the registry digests the image was pulled under. Before it dials, and again whenever the connection is opened anew, flai compares that digest with the signed digest list of [Signing releases](#how-the-images-signature-is-published), fetched and verified as `flai dashboard` did when it chose the image, and cached beside `flai serve`'s state. A digest on the list is a signed release; flai dials. A digest not on any list is unsigned; flai does not dial, logs one `error` event with the digest and the image, records the state so that `flai serve status`, `flai dashboard status`, and `flai host status` say "the dashboard runs an unsigned image", and tries again when the container changes. This is a real check: Docker, not the peer, says what the container runs, and flai already trusts Docker to run the container at all.

Its limits: it needs Docker where `flai serve` runs, so a dashboard reached over the network on another machine cannot be measured and falls back to the stamp; `RepoDigests` is empty for a locally built image, which is the development case and is right to refuse by default; and the list must name every digest a pull can leave behind, the index digest and each platform's manifest digest, because which one `RepoDigests` records depends on how the image was pulled (**to check**, in the story that builds it). Since `flai dashboard` will already run the image by digest from the signed list, this check at dial time is a second look at a decision flai itself made, and it catches what happened in between: a container restarted from a different image, a tag that moved, an operator's `docker run`.

### A release stamp in `hello`

The stamp is a small statement CI signs with the release key before the build, and the build carries: the component name, the version, and the commit, as one line, with its signature. For flai the workflow signs the line and hands the signature to GoReleaser through the environment, which `ldflags` write into `buildinfo` beside `Version`; for flaiover the workflow signs the line and passes it as a build argument that the image keeps as a file. The statement cannot cover the finished artifact's hash, because the artifact does not exist yet when it must be embedded; it covers what the release *is*, and the artifact checks of the first level cover what it *contains*.

`hello` grows one field each way: flai sends `release: {statement, signature}` with the fields it sends today, and the dashboard answers with its own. Each side verifies the other's signature with the public key it carries, checks that the statement names the component it expects (`flaiover` from the dashboard, `flai` from flai) and that the version is one the protocol number already admits, and then goes on with the credential proof. A missing, unverifiable, or wrong-component stamp closes the connection with a close code of its own, `4403 unsigned peer`, from whichever side found it; flai backs off a minute on it as it does on `4409` ([ADR-0063](../adrs/0063-the-dashboard-refuses-a-second-flai-for-a-project-with-close-code-4409-while.md)), logs once, and shows the state everywhere the measured check does. The dashboard shows a project whose flai was refused with the reason, where it shows a flai that lacks a required method today, so that the designer sees "unsigned flai" and not a connection that never comes.

The stamp is what lets flaiover say anything about flai at all, since the container can measure nothing on the host; and it is what covers flai's side when Docker is out of reach. It is checked by both for symmetry and because it costs nothing more once the key is in both sources. Its promise is stated above and should be written the same way in the operator documentation: it names the peer's release and refuses builds that are not one; it does not defend against someone who holds the credential.

### The development case

This repository runs flai from `scripts/flai.sh`, with no version and no stamp, and its dashboard from `flai dashboard --build`, an image with no digest on any list. Every contributor and every project on the template that builds its own flai is in the same position, so unsigned must be allowed on purpose, in one place, and visible.

- **The allowance is the operator's, on the host, once.** A setting in the host's flai configuration, `dashboard.allow_unsigned: true`, with the project manifest able to set it for one project as it sets `dashboard.image`. With it, flai dials a container whose image is not on a list, and accepts a peer with no stamp; and `flai dashboard`, when it starts a container with the allowance on or with `--build`, gives the container `FLAIOVER_ALLOW_UNSIGNED=1`, so that the dashboard in turn accepts a flai with no stamp. The container is told by the thing that started it, which is the one act that makes a build a development build; nothing in the image decides for itself.
- **Allowed is not silent.** With the allowance on, `flai serve` logs a `warn` once per connection, `flai dashboard status` and `flai serve status` say "unsigned allowed", the dashboard shows a banner on every page naming which side is unsigned (its own build, the project's flai, or both), and the connection list says it per project. A release build of flai talking to a release image with the allowance on shows nothing: the allowance permits, it does not downgrade.
- **The default is to refuse**, for a release build of either side. An operator who installed flai with `install.sh` and runs the published image is never asked, and never shows the banner; the one who builds either side is told what they built.
- **Release builds are never allowed to be unsigned.** A `flai/v*` tag whose workflow runs without the key fails the release, and a `flaiover/v*` tag likewise; the stamp is not optional for a build that claims a version. `goreleaser --snapshot` and `flai dashboard --build` give a version of `0.0.0-<commit>` or none, and that is what "development build" means in both places.

### Recommendation

Do all three levels, in this order of worth: verify every artifact before it is installed or started, as [Signing releases](#signing-releases) recommends; have flai measure the container's image digest through Docker against the signed list before it dials and refuse an unlisted one; and exchange a signed release stamp in `hello`, checked by both, closed with `4403` when it fails. Allow unsigned only by the operator's setting on the host, passed to the container by `flai dashboard`, and show it wherever the connection is shown. Write the threat model into `docs/operators` in the words of the table above, so that the connection check is read as what it is: a check that the peer is a release, with the credential as the thing that keeps the channel the operator's.

## Decision

The designer took the recommendation on every question in TH-0065 on 2026-10-02, recorded as [ADR-0070](../adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md): one cosign key pair signs flai's `checksums.txt` and flaiover's digest list; flai, `install.sh`, and `flai dashboard` verify before they install or run, with the standard library, `openssl`, and Docker's digest; flai measures the container's image before it dials; both sides exchange a signed release stamp in `hello` and close with `4403` when it fails; and unsigned is allowed only by `dashboard.allow_unsigned` on the host, passed to the container by `flai dashboard`, and shown wherever the connection is. The stories under E-0015 that deliver it follow.
