---
id: ADR-0070
title: "Releases are signed with a cosign key pair, verified before they are installed or run, and each component refuses a peer without a signed release stamp"
status: accepted
date: 2026-10-02
supersedes: []
superseded_by: []
refines: [ADR-0029, ADR-0031]
---

# ADR-0070 Releases are signed with a cosign key pair, verified before they are installed or run, and each component refuses a peer without a signed release stamp

## Context

E-0015 asks that flai and flaiover be signed in CI with a private key, that flai verify a release with the public key before it installs it, and that each component refuse a connection with an unsigned build of the other. Today GoReleaser publishes `checksums.txt` unsigned, `flai self-upgrade` and `install.sh` check a download against it, `flai dashboard` runs whatever a registry tag resolves to, and the channel's HMAC handshake proves that both sides hold the agent credential ([ADR-0031](0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)) and nothing about either build. S-0193 set out the ways to sign and what a peer can verify in `design/system/release-signing.md`; the designer chose the recommendation on every question in TH-0065 on 2026-10-02.

## Decision

**One key pair signs both components.** An ECDSA P-256 pair made with `cosign generate-key-pair`, the encrypted private key and its password held as two GitHub Actions secrets of `bytepunx/system-flow`, used only by the two release workflows. The public key is a constant in flai's and flaiover's sources, published as a release asset, and its fingerprint is in the operator documentation. Rotation is a release that carries both public keys before the workflows switch; a compromise is announced with the new fingerprint.

**What is signed.** For flai, `checksums.txt`, through GoReleaser's `signs` with `cosign sign-blob` and `--tlog-upload=false`. For flaiover, a digest list naming the index and per-platform digests the workflow pushed, signed the same way and published with its signature on a GitHub release `flaiover vX.Y.Z`, which the workflow creates on a `flaiover/v*` tag. `cosign sign` on the image and `actions/attest-build-provenance` on both workflows are added beside, for anyone with cosign or `gh`; flai depends on neither.

**What verifies, with what.** `flai self-upgrade` and `flai host upgrade` verify the signature over `checksums.txt` with the standard library and the key in the binary before they check the archive's hash; a release they cannot verify is not installed, and one signed by a key they do not know says which release to install first. `install.sh` verifies with `openssl dgst -sha256 -verify`. `flai dashboard`, `check`, and `upgrade` resolve the version to run from the signed digest list and pull and run the image by digest; `dashboard.tag: latest` means the newest flaiover release. Snapshot and local builds are not signed and have no version.

**Each component refuses an unsigned peer.** CI signs a release stamp, the component's name, version, and commit, before each build, and the build carries it: in flai's `buildinfo`, in the image as a file. `hello` carries flai's stamp and the answer the dashboard's; each side verifies the other's with the public key and checks the component it names, before the credential proof. Before it dials, and whenever the container changes, flai also compares the digest Docker reports for the container with the signed list. A missing, unverifiable, or wrong stamp, or an unlisted image, closes the connection with code `4403 unsigned peer`; flai backs off a minute as on `4409` ([ADR-0063](0063-the-dashboard-refuses-a-second-flai-for-a-project-with-close-code-4409-while.md)), logs once, and the state is shown by `flai serve status`, `flai dashboard status`, `flai host status`, and on the dashboard against the project.

**Unsigned is allowed only on purpose.** `dashboard.allow_unsigned` in the host's flai configuration, settable per project in the manifest, default off. With it, flai dials an unlisted container and accepts a peer without a stamp, and `flai dashboard`, when it starts a container with it on or with `--build`, gives the container `FLAIOVER_ALLOW_UNSIGNED=1`, so the dashboard accepts a flai without a stamp. Allowed is shown: a `warn` per connection, "unsigned allowed" in every status, a banner on every dashboard page naming the unsigned side, and the connection list per project.

**What this promises.** The artifact checks prove that what is installed or started is a release the key signed. The Docker measurement proves that the container flai talks to runs one. The stamp proves that the peer carries a release's statement; it does not prove that the peer is that release against anyone who holds the agent credential and a release binary, because no process can prove its own code to a peer without a measurement it does not control. The credential remains the boundary of the channel, as ADR-0031 says. The operator documentation says this in these words.

## Consequences

- Both release workflows gain a signing step and the key's two secrets; a release tag whose workflow runs without the key fails. flaiover gains GitHub releases, which it has not had.
- flai gains a public key, signature verification, a release stamp, a dial-time image check, a close code, a setting, and status fields, all with the standard library and `docker`; flaiover gains the key, the stamp check, the setting's environment variable, the banner, and the connection state.
- `install.sh` needs `openssl`, which every target ships; it refuses to install what it cannot verify.
- A dashboard started from `main`'s `latest` is unsigned and is refused unless allowed; `latest` is resolved from releases instead.
- A flai older than the first signed release upgrades to it without a signature check, once; from then on every upgrade is verified.
- This repository, which builds both from source, sets `dashboard.allow_unsigned: true` in `.flai-cache/config.json` and shows the banner.
- `design/system/release-signing.md § Decision` lists the stories under E-0015 that deliver this.

## Alternatives considered

- **minisign (Ed25519).** Smaller and standard-library-verifiable in flai, but `install.sh` cannot verify without `minisign`, which hosts lack, and macOS's LibreSSL does not verify Ed25519; the script's check would stay nominal.
- **cosign keyless, or GitHub attestations, as the only signature.** No key to keep, but flai would need `sigstore-go` or `gh` and Sigstore's or GitHub's network to verify, and `install.sh` cosign; kept beside the key for auditors instead.
- **GPG.** The format and tooling were made for people, not for a binary checking itself; heavy in Go and awkward in a shell.
- **flai verifies cosign's OCI signature in GHCR.** A registry client against a layout cosign owns, or heavy dependencies, for what the signed digest list gives with the client `self-upgrade` already has.
- **Warn only at first, refuse later.** The check is a refusal or it is a notice; the development allowance covers the only case a warning would.
- **A developer signing key so that nothing is ever unsigned.** Every local build would have to be signed, and the key's presence would make a development build look like a release; an explicit allowance shown on every page is plainer.
