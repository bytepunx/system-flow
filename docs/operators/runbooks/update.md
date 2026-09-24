---
title: "Runbook: update"
updated: 2026-09-24
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

To go back, `flai self-upgrade --version <the previous release>`, then restart the host as in step 3.

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

`latest` follows the main branch. To stay on a release, pin its tag, then upgrade to it; the same moves back to an earlier one:

```bash
flai config set dashboard.tag 0.27.4   # or dashboard.tag in system-flow.yaml, for one project
flai dashboard upgrade
```

`flai dashboard upgrade --tag 0.27.4` does it once without changing the setting; the next `flai dashboard` that starts the container uses the setting again.
