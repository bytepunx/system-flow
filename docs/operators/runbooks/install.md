---
title: "Runbook: install"
updated: 2026-09-24
status: active
---

# Install

## Before you start

| Need | For | Check |
|------|-----|-------|
| Linux or macOS on amd64 or arm64 | `install.sh`; on Windows use the zip from the releases page, and note that `flai host` there has not been tried | `uname -sm` |
| `curl`, `tar`, and `sha256sum` or `shasum` | `install.sh` | `command -v curl tar sha256sum shasum` |
| `git` | every project command | `git --version` |
| Docker Engine 24 or newer, on `PATH` | the dashboard only | `docker version` |
| A GitHub token, while the repository is private: `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth login` | the download, and the dashboard image (scope `read:packages`) | `gh auth status` |

Nothing here needs `sudo`. Everything flai keeps is under `~/.flai` and each project's `.flai-cache/`.

## flai

1. Install the newest release into `~/.flai/bin`:

   ```bash
   curl -fsSL https://raw.githubusercontent.com/bytepunx/system-flow/main/install.sh | sh
   ```

   `FLAI_VERSION=1.16.1` pins a release; `FLAI_INSTALL_DIR=/some/dir` installs elsewhere, into a folder you can already write to. The script verifies the archive's SHA-256 against the release's `checksums.txt` and stops if it differs.
2. Put the folder on your `PATH`, with the line the script prints, in your shell profile. If it says an earlier `PATH` entry takes precedence, another `flai` comes first: remove it or reorder.
3. Check it:

   ```bash
   flai version
   flai config path     # creates ~/.flai/config.json with the defaults
   flai config get
   ```

4. Change what you need in the configuration ([settings](../settings.md#flai-configuration)). Most installs change nothing. A fork of the template: `flai config set template.repo <url>`.
5. Start or join a project:

   ```bash
   flai new my-project          # a new repository from the template
   cd existing-repo && flai import   # or bring an existing one under the standard
   flai check --strict          # in a project: the repository conforms
   ```

   Both are in [the flai guide](../../users/flai.md#create-a-project). A project's `.mcp.json` runs `flai mcp` from your `PATH`, so an agent in it reaches the flai you just installed.

If the download is refused, the repository is private and no token was found: `gh auth login`, or set `GITHUB_TOKEN`, and run it again.

## flaiover

1. Give Docker access to the image, while it is private: `gh auth refresh -h github.com -s read:packages`. `flai dashboard` logs Docker in with that token when a pull is refused ([Access to the image](../index.md#access-to-the-image)). Inside this monorepo, `flai dashboard --build` builds the image instead.
2. Decide where it listens. On a machine only you use, keep it to this machine: `flai config set dashboard.bind 127.0.0.1`. Read [Security posture](../index.md#security-posture) before you publish it on any other address.
3. In a project, start it:

   ```bash
   flai dashboard
   ```

   This pulls the image, starts the one `flaiover` container for this host, starts `flai host` (which runs `flai serve` and the project's MCP server), registers the project, and prints a login link. Run it in each project you want on the board, or once in a folder above them to serve every project below it.
4. Open the login link. The header says `host flai: connected` once `flai serve` has reached the dashboard.
5. Check it:

   ```bash
   flai dashboard status   # the container, and a host flai line
   flai host status        # the host, flai serve, and each MCP server: running
   flai serve status       # the projects served and the dashboards connected
   ```

6. Leave the host actions off unless you want them from the board. Each is turned on by name, in a shell: `flai serve actions` lists them, and [the operators guide](../index.md#the-push-host-action) says what each gives whoever holds the token.

| When | Do |
|------|----|
| The port is taken | `flai dashboard --port 8080`, or `flai config set dashboard.port 8080` |
| The pull is refused | The token lacks `read:packages`: step 1 again, or `flai dashboard --build` in the monorepo |
| `flai host` will not start: another program holds `127.0.0.1:4241` | Set `FLAI_HOST_ADDR` to another loopback address for every flai you run |
| `flai dashboard` is refused because a host runs with another configuration | One host runs per machine, for one configuration; use the same `FLAI_CONFIG`, or stop the other host with `flai host stop` |
| Pages say no flai is connected | `flai serve status` says why; `flai host restart serve` |
