---
title: "Runbook: update"
updated: 2026-10-12
status: active
---

# Update

flai and the dashboard are released separately (`flai/vX.Y.Z` and `flaiover/vX.Y.Z` tags) and are updated separately. Neither update touches a project's files. Bringing a project to a newer template is in [migrate.md](migrate.md#a-project-to-a-newer-template).

Before either, read the release notes on the [releases page](https://github.com/bytepunx/system-flow/releases) for anything marked breaking, and the operators guide's sections headed "Upgrading from", which say what to do when a release changed how things are kept.

## flai

1. See whether there is one:

   ```bash
   flai self-upgrade --check   # the installed version and the newest
   flai host check             # the same, for the flai the host runs
   ```

2. Update, and restart everything that runs flai, in one step:

   ```bash
   flai host upgrade
   ```

   The host's flai verifies the release's signature with the release key built into it, downloads the archive and checks it against the signed `checksums.txt`, installs it over the flai it runs, then stops `flai serve` and every MCP server it runs and starts itself again on the new binary. From the dashboard, the Upgrade button does the same once the `host` host action is on.
3. Or, with no host running, or to choose the release:

   ```bash
   flai self-upgrade                   # the newest
   flai self-upgrade --version 1.16.1  # a given release
   flai host stop && flai host start   # a running host keeps the old binary until it restarts
   ```

   If the release is refused, nothing was installed and nothing restarted. The refusal names the release and why: it has no signature, its signature does not verify, or it is signed by a key this flai does not know. For an unknown key it names the release to install first: install that one with `--version`, then upgrade again. What each refusal means is in [What an upgrade verifies](../../users/flai.md#what-an-upgrade-verifies), and the key's fingerprint under [Current key](release-key.md#current-key).

4. Check it:

   ```bash
   flai version
   flai host status   # every process on the same version
   ```

5. Agents already running keep the flai they started with: an MCP server that `.mcp.json` started over stdio runs until its session ends. Start the session again to give it the new one.

### To go back, or to a chosen release

1. List the published releases, newest first:

   ```bash
   flai host versions          # through the host
   flai self-upgrade --list    # where no host runs
   ```

   Each list marks the installed release, the newest, any below the `flai.minimum` of a project `flai serve` serves, and, `not installable`, any the flai that would install it cannot verify. Those are refused, so you cannot go back past flai 1.42.0, the first signed release.
2. Install one:

   ```bash
   flai host upgrade --version 1.16.1     # through the host
   flai self-upgrade --version 1.16.1     # where no host runs; then restart the host as in step 3
   ```

   A release that is not published is refused, naming the published ones, before anything is downloaded. A release below a project's `flai.minimum` is installed with a warning: a flai below it will not read that project. Through the host, an earlier release restarts the host, `flai serve`, and the MCP servers exactly as installing the newest does.

From the dashboard, the flai host area of the Updates page has **Versions**, which lists the same releases and deploys a chosen one once the `host` host action is on ([flaiover guide](../../users/flaiover.md#host)).

A `flai` that runs from inside a system-flow project, such as a checkout's own `bin/flai`, is not replaced: `flai self-upgrade` installs the release into `~/.flai/bin` (or `FLAI_INSTALL_DIR`) and says to put that folder first on your `PATH`.

## flaiover

1. See whether there is one:

   ```bash
   flai dashboard check
   ```

2. Update:

   ```bash
   flai dashboard upgrade
   ```

   It pulls the configured image and tag and, when it differs from what runs, starts it beside the running container, waits for it to answer healthy, and only then replaces the running one. If the new one never answers healthy, the running container is left as it was and the command says why. If it answers healthy but then fails to start in the running one's place, the command starts the previous image again and fails with an error ending `<name> runs the previous image <ref> again`. Only if that fails too does the error end `the dashboard is down: run flai dashboard to start it`; run `flai dashboard`. From the dashboard, its own Upgrade button does the same once the `dashboard` host action is on.
3. Check it: `flai dashboard status`, and the dashboard's `/metrics` reports the release in `flaiover_build_info`. Browser sessions survive: the login token did not change.

### To go back, or to a chosen release

1. List the published releases, newest first. Each is a `flaiover/vX.Y.Z` tag, and its image tag is the bare `X.Y.Z`:

   ```bash
   flai dashboard versions   # marks the running release, the configured tag, and the newest
   ```

2. Deploy one:

   ```bash
   flai dashboard upgrade --published --tag 0.27.4
   ```

   It swaps the container as `flai dashboard upgrade` does, only once the new one answers healthy. With `--published` a tag that is not a published release is refused, naming the published ones, before anything is pulled. Without it the tag is used as given, such as a mirror's.

From the dashboard, the Dashboard area of the Updates page has **Versions**, which lists the same releases and deploys a chosen one once the `dashboard` host action is on ([flaiover guide](../../users/flaiover.md#host)).

A chosen release applies to the container, not the configuration. It keeps running through restarts (`flai dashboard restart`, or flai host's watch) until the next `flai dashboard upgrade` without a tag, the page's Upgrade, which goes to the configured `dashboard.tag`, or until `flai dashboard` starts a container after `flai dashboard stop` ([ADR-0118](../../../design/adrs/0118-a-dashboard-release-chosen-on-the-updates-page-keeps-running-through-restarts.md)).

`latest` follows the main branch. To stay on a release, pin its tag, then upgrade to it:

```bash
flai config set dashboard.tag 0.27.4   # or dashboard.tag in system-flow.yaml, for one project
flai dashboard upgrade
```

### Released and unreleased images

| Tag | Built from | Signed |
|-----|------------|--------|
| `X.Y.Z`, and its major `X` | a `flaiover/vX.Y.Z` tag | yes: the image carries a cosign signature by the release key, and its digests are on the release's signed digest list |
| `latest`, `sha-<short>` | every push to main | no: no signature, on no digest list |

`latest` and `sha-*` are builds of main, not releases. Today `flai dashboard` runs them as it runs a release. Once `flai dashboard` verifies images against the signed list, it will run such a build only as a development image, which the operator allows with the `dashboard.allow_unsigned` setting. A released `X.Y.Z` needs no allowance. How the list and the check work is in [release-signing.md](../../../design/system/release-signing.md#the-digest-list).

### Check a release image by hand

1. From the release `flaiover vX.Y.Z` on the [releases page](https://github.com/bytepunx/system-flow/releases), download `flaiover_X.Y.Z.digests` and `flaiover_X.Y.Z.digests.sig`, and `cosign.pub` from the same release or from your own copy. Compare the key's fingerprint with the one under [Current key](release-key.md#current-key) first.
2. Verify the list, with cosign:

   ```bash
   cosign verify-blob --key cosign.pub --signature flaiover_X.Y.Z.digests.sig --insecure-ignore-tlog flaiover_X.Y.Z.digests
   ```

   or with openssl:

   ```bash
   openssl base64 -d -A -in flaiover_X.Y.Z.digests.sig -out flaiover_X.Y.Z.digests.der
   openssl dgst -sha256 -verify cosign.pub -signature flaiover_X.Y.Z.digests.der flaiover_X.Y.Z.digests   # Verified OK
   ```

3. Verify the image's signature in the registry, by the `index` line's digest:

   ```bash
   cosign verify --key cosign.pub --insecure-ignore-tlog=true ghcr.io/bytepunx/flaiover@sha256:…
   ```

4. Check that the image you run is on the list: `docker image inspect --format '{{json .RepoDigests}}' ghcr.io/bytepunx/flaiover:X.Y.Z` names a digest the list has.
