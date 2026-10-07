---
title: "Runbook: update"
updated: 2026-10-07
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

   The host downloads the release (checked against its checksum), installs it over the flai it runs, then stops `flai serve` and every MCP server it runs and starts itself again on the new binary. From the dashboard, the Upgrade button does the same once the `host` host action is on.
3. Or, with no host running, or to choose the release:

   ```bash
   flai self-upgrade                   # the newest
   flai self-upgrade --version 1.16.1  # a given release
   flai host stop && flai host start   # a running host keeps the old binary until it restarts
   ```

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

   Each list marks the installed release, the newest, and any below the `flai.minimum` of a project `flai serve` serves.
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

   It pulls the configured image and tag and, when it differs from what runs, starts it beside the running container, waits for it to answer healthy, and only then replaces the running one. If the new one never answers healthy, the running container is left as it was and the command says why. From the dashboard, its own Upgrade button does the same once the `dashboard` host action is on.
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
