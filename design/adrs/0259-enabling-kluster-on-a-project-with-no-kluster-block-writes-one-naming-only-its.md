---
id: ADR-0259
title: "Enabling kluster on a project with no kluster: block writes one naming only its config and a kluster.yaml with the spire profile, deploy becomes optional, every setting and the port-forwards are worked from the dashboard, and settings stays the one host action turned on in a shell"
status: accepted
date: 2026-11-07
supersedes: []
superseded_by: []
refines: [ADR-0254]
topics: [kluster, deployment, manifest, dashboard]
---

# ADR-0259 Enabling kluster on a project with no kluster: block writes one naming only its config and a kluster.yaml with the spire profile, deploy becomes optional, every setting and the port-forwards are worked from the dashboard, and settings stays the one host action turned on in a shell

## Context

[ADR-0254](0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md) has a project's kluster run from the manifest's `kluster:` block. It makes `config` and `deploy` required (§3), and refuses the `kluster` host action for a project with no block (§7). So turning kluster on from the dashboard's Features page failed with a sentence telling the operator to edit `system-flow.yaml`, and the toggle looked on until a refresh. S-0586 asks for the opposite: turning kluster on gives the project a block that does nothing yet, or brings up an empty cluster, and the Klusters page holds everything setup, management, and control need, with no command or file edit outside the dashboard, since a hosted flai cannot assume a shell.

Three facts shaped the answer:

- **kluster has no empty cluster.** In kluster v1.3.1, `kluster up` takes the profile from `--profile`, then the config's `profile`, then the flag's default, `spire` (`kluster/cmd/up.go` and `stringFlag` in `kluster/cmd/config.go`, at kluster's commit 6d7a87c). An empty `profile` falls back to `spire`, and an unknown one is refused (`resolveProfiles` in `kluster-lib/cluster/cluster.go`). No profile installs nothing. The plan's assumption that a config with no profile gives an empty cluster was wrong.
- **The tools cannot be installed from flai serve.** kluster's `install.sh` writes `/usr/local/bin`, and `kluster setup` installs `kubectl` with `sudo` on Linux. flai serve has no terminal to answer `sudo`. Docker needs root and a daemon, and `kluster setup` never installs it.
- **The settings host action is turned on only in a shell** ([ADR-0039](0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md)). The block's `build`, `deploy`, and `teardown` are commands flai runs on the host as the operator, so whoever sets them holds a shell there.

The operator answered the plan's four questions on TH-0861 on 2026-11-07 (D-0511): `settings` stays the one shell step; no Install button; Start and Stop for each port-forward, bound to `0.0.0.0` so that a hosted flai can reach them; and `profile: spire` written out.

## Decision

**Enabling kluster on a project with no `kluster:` block writes one naming only its config and a `kluster.yaml` with the `spire` profile, `deploy` becomes optional, every setting and the port-forwards are worked from the dashboard, and `settings` stays the one host action turned on in a shell.**

1. **Enabling writes the block.** `flai serve enable kluster`, which the Features toggle and the Klusters page run through `settings.action`, gives a project with no block a `kluster:` block naming only `config: kluster.yaml`. It writes `kluster.yaml` at the repository root when that file is absent, since the layout allows no new top-level folder, and keeps one that is there. It commits those two files alone, with the dashboard's trailer, unless the project sets `dashboard.autocommit: false`. The manifest is edited a line at a time, so every other line stays as it was. This replaces ADR-0254 §7's refusal. `--all-projects` writes no block.
2. **The scaffold's config.** `kluster.yaml` names the cluster, as the manifest's name made a DNS label with `-local` after it, the provider `k3d`, and the profile `spire`, written out. `spire` is kluster's own default and its smallest documented profile: cert-manager, SPIRE, and traefik-tls. The Klusters page changes the profile, among `spire`, `signet`, and `authstar`, and the addons, before the first Up.
3. **`deploy` is optional.** This replaces ADR-0254 §3's "required". A block with no `deploy` is a project whose cluster can come up and go down and whose deploys are not set up yet. flai serve deploys and promotes nothing for it, reading the block at every look, so a `deploy` set later is picked up at the next look. It still tears down a namespace already deployed. `flai kluster deploy` and `promote`, and the dashboard's Deploy and Promote, refuse it. `flai check` no longer reports a missing `deploy`.
4. **The block is written from the dashboard, under `settings`.** Each key, `config`, `namespace`, `build`, `deploy`, `teardown`, and `timeout`, is a setting `flai manifest set` writes and `settings.manifest` reaches. The commands are argument lists, one argument a line on the page. Those writes stay under the `settings` host action, as ADR-0039 has them, because the commands run on the host as the operator.
5. **The config file is saved through `file.get` and `file.save`.** [ADR-0147](0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md)'s reach now includes the file `kluster.config` names, and no other path outside the design and docs folders. That save is not held by `settings`: the file names a profile and addons, not a command, and bringing the cluster up needs the `kluster` host action anyway.
6. **`settings` stays the one shell step.** ADR-0039 stands: no dashboard turns `settings` on. A hosted flai turns it on when it is provisioned.
7. **No Install button.** The Klusters page reports whether `kluster`, `docker`, and `kubectl` are on the host's `PATH`, and for a missing one what flai needs it for, the sentence that installs it, and a link to its documentation. flai installs none of them.
8. **The port-forwards are run by flai serve.** The Deployments page starts and stops each service port's forward, which flai serve runs on the host with flai's kubeconfig, bound to `0.0.0.0`, behind the `kluster` host action. A running forward is reachable by anything that reaches the host's network, with no sign-in of its own. The forwards stop when flai serve stops, and a namespace's when it is torn down.

## Consequences

- Turning kluster on from the dashboard works for a project that has never had a cluster, and the Klusters page leads the operator from there to a running cluster in five steps: the tools, turning kluster on, the cluster's config, the deploy commands, and Up. A button there makes a backlog story for an agent to write the project's build, deploy, and teardown scripts, since they are the project's code.
- Two steps stay outside the dashboard, both host changes the operator makes once: turning `settings` on, and installing Docker, `kluster`, and `kubectl`. The documentation names both.
- `kluster.status` gains what setup still lacks: the block's settings, whether `settings` is on, the config file's name, provider, profile, and addons, `flai check`'s findings on the block, whether a `deploy` is set, and the host's tools. A new read, `kluster.workloads`, lists the cluster's nodes and each namespace's pods.
- The Deployments page and the Features toggle name no command. A toggle flai refuses goes back at once to the state flai holds, with flai's sentence beside its row.
- A forward on `0.0.0.0` exposes a story's services to the host's network. On a shared network the operator stops a forward when the demo ends.
- A Playwright spec, `flaiover/e2e/klusters.e2e.ts`, walks the Klusters page in Chromium from no block to a running cluster, with the host's answers faked.

## Alternatives considered

- **A dashboard toggle for `settings`.** It would meet "no command outside the dashboard" literally, but it reverses ADR-0039: the dashboard token alone would be worth a shell on the host, with no consent given there.
- **An Install button run by flai serve.** `kluster setup` needs `sudo` for `kubectl` and the release binary goes to `/usr/local/bin`, so the button would hang or fail without passwordless sudo, and it cannot install Docker at all.
- **Forwards bound to `127.0.0.1`.** Safer on a shared network, but a hosted flai, and a browser on any other machine, could not reach them.
- **A placeholder such as `deploy: ["true"]`.** It would pass `flai check` with `deploy` still required, but flai serve would bring the cluster up for every story in review and report it deployed with nothing deployed.
- **A config with no profile.** It reads as an empty cluster but installs `spire` all the same, so writing the profile out says what the cluster gets.
- **A new host method to write the config file.** `file.save` already has the hash, the check, and the commit, and the manifest names the path.
