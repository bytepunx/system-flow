---
title: Pushing an acceptance made from the board
updated: 2026-09-19
status: active
---

# Pushing an acceptance made from the board

The finding of S-0052. An acceptance made from the flaiover board is merged, archived, committed, and tagged in the clone, and not pushed: `flai dashboard` gives the container the host's git identity and global excludes and no credential ([ADR-0018](../adrs/0018-dashboard-token.md), [operators](../../docs/operators/index.md)). The board says "accepted locally" and prints the `git push` command. The operator was surprised by this twice, and between 2026-09-18 and 2026-09-19 an agent session pushed `main` and the release tags from the host after every one of ten board acceptances. This document weighs the ways the push could happen, says what was tried, and ends with one recommendation. Nothing here is built; the decision and the stories that follow are at the end.

> Superseded in practice, 2026-09-20: the push key this finding led to (ADR-0026, S-0062) was retired by S-0077 ([ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)). The container now holds no git, no ssh, and no file of the project, so nothing can be pushed from it; an acceptance is pushed from the host (`flai push --pending`, and by `flai serve` on the operator's word with S-0078). What follows is the finding as it was made.

## What is at stake

A push of `main` with a tag `flai/vX.Y.Z` or `flaiover/vX.Y.Z` starts the release workflows, which publish binaries and the dashboard image that `flai dashboard --pull` and `install.sh` hand to every user. Whatever can push tags can publish a release.

The dashboard is published on every interface of the host by default and guarded by one project token (ADR-0018); S-0043 shaped it to be reached through a hub. Two kinds of person matter:

- **A holder of the dashboard token, using the API.** They can do what the API allows: move items, accept a story that is in review, save Markdown under `design/`, `docs/`, and `wip/` (the editor refuses every other path and every other file type, `flai/internal/docedit`), and talk to agents through threads. Since S-0076 the token no longer opens MCP: that is flai's on the host, with a token of its own (ADR-0030). They cannot run git. If acceptances were pushed automatically, this person could publish a release of any story an agent has put in review, and their Markdown edits would reach the remote with it. They cannot force push, delete tags, or change workflow files, because the API never asks git to.
- **Someone who has compromised the container** (a flaw in flaiover or one of its dependencies). They can do whatever the container's user can do with what the container was given. A credential in the container is theirs.

One thing found on the way changes how much the second case is worth protecting against, and it is true today, with no credential in the container. The clone is mounted read-write, `.git` included, because acceptance needs it. A process in the container can write `.git/hooks/pre-push`, and the next `git push` the operator or an agent runs on the host executes it as the operator. **Tried** in a scratch project: a hook written from inside the container ran on the host as the host user at the next push. So a compromised container can already reach the operator's credentials, with a delay and with the operator's unknowing help. Keeping credentials out of the container limits the token holder; against a compromised container it buys time, not safety. Follow-up 3 below is about this.

## How it was tried

A scratch project with the flai layout, and a scratch remote in a container on its own docker network: a bare repository served over SSH (OpenSSH, throwaway ed25519 keys) and over HTTP (busybox `httpd` with `git-http-backend` and a throwaway password). Pushes were made from containers of the published dashboard image (`ghcr.io/bytepunx/flaiover`, 0.17.0) run the way `flai dashboard` runs it: as the host user, with the project mounted at its host path. Host: WSL2, Docker Engine. Everything was generated for the purpose and deleted afterwards. The operator's SSH agent, keys, `gh` token, and this repository's remote were not used or looked at. TLS was not part of the test: the HTTP remote stands in for HTTPS, and only the credential handling is claimed.

Not tried, and marked so below: anything on macOS or Windows, anything against GitHub itself. Statements about GitHub and Docker Desktop are from their documentation, cited.

## Pushing from inside the container

The image has `git` with the HTTP transport and the credential helpers, and **no SSH client** (`ssh`, `ssh-add`, and `ssh-keyscan` are absent; **tried**: `git push ssh://…` fails with "cannot run ssh"). With no credential an HTTP push fails with "could not read Username … terminal prompts disabled" (**tried**). So the two SSH mechanisms need the image changed (`openssh-client`, about 0.7 MB on Alpine); the token needs nothing in the image.

### An HTTPS token given to a credential helper

- **What `flai dashboard` would do.** Mount a token file read-only (say `/run/flaiover/git-token`) and set a credential helper that reads it, for the push only: `git -c credential.helper='!f() { echo username=x-access-token; echo password=$(cat /run/flaiover/git-token); }; f' push …`. **Tried**: the push succeeded from the unmodified image as the host user, and nothing was written into the clone's config. Passing the token in the environment also works (**tried**) and is worse: every process in the container and `docker inspect` on the host can read it.
- **Hosts.** A file mount and HTTPS work the same on Linux, WSL2, and Docker Desktop on macOS and Windows. **Not tried** off WSL2; nothing in it is platform specific.
- **What it lets someone do.** A token holder using the API: publish by accepting, as above. A compromised container: read the token and do to the repository whatever the token allows. On the unprotected scratch remote that was everything (**tried**: force push over `main`, create a tag, delete a tag, all accepted). The limits have to come from the provider.
- **How narrowly it can be scoped.** A GitHub fine-grained personal access token can be limited to one repository with the Contents permission set to write, which is what pushing commits and tags needs; changing files under `.github/workflows` needs the separate Workflows permission, which would not be given ([permissions for fine-grained tokens](https://docs.github.com/en/rest/authentication/permissions-required-for-fine-grained-personal-access-tokens)). Rulesets on the repository can block force pushes and restrict deletion, creation, and update of matching branches and tags ([available rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets)). Restricting tag creation would also stop the legitimate release tags unless the token's owner is on the bypass list, so the realistic protection is: no force push, no deletion, no workflow changes. A stolen token could still publish a release from a commit of the attacker's making. **Not tried** against GitHub.
- **Rotation and revocation.** The token has an expiry chosen at creation and is revoked in the account's settings; replacing the file and restarting the dashboard rotates it.
- **A failed push** (expired token, no network) is reported as it is today: accepted locally, with the command. Pushing later is idempotent.

### A deploy key mounted read-only

- **What `flai dashboard` would do.** Add an SSH client to the image; mount the private key and a `known_hosts` file read-only; set `GIT_SSH_COMMAND="ssh -i /run/flaiover/deploy_key -o IdentitiesOnly=yes -o UserKnownHostsFile=/run/flaiover/known_hosts -o StrictHostKeyChecking=yes"`. **Tried** with a lab image that adds `openssh-client-default`: the push succeeded as the host user.
- **Hosts.** OpenSSH refuses to run for a user ID with no passwd entry: "No user exists for uid 4321" (**tried**). `flai dashboard` runs the container as the host user; it worked here only because the host user is 1000, which the image happens to name `node`. Any other host UID needs a passwd entry made at start. The key file must be mode 600 and owned by the container's user; true with a bind mount on Linux and WSL2, **not tried** on Docker Desktop for macOS or Windows.
- **What it lets someone do.** As the token: the API publishes by accepting; a compromised container has the private key. A deploy key grants access to a single repository, read-only unless write is given, and "a deploy key with write access lets a deployment push to the repository" ([managing deploy keys](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/managing-deploy-keys)). That page sets no limit on workflow files for a write key, and rulesets' bypass lists name roles, teams, and apps, not deploy keys, so whether rules bind a deploy key the way they bind a token was **not established**.
- **Rotation and revocation.** "Deploy keys are credentials that don't have an expiry date" and "are usually not protected by a passphrase, making the key easily accessible if the server is compromised" (same page). Revoked by deleting the key from the repository's settings.

### The host's SSH agent socket, forwarded

- **What `flai dashboard` would do.** Add an SSH client to the image; mount `$SSH_AUTH_SOCK` and set `SSH_AUTH_SOCK` in the container; mount `known_hosts`. **Tried** with a scratch agent holding a throwaway key: the container listed the key and pushed; no key material was in the container.
- **Hosts.** Linux and WSL2 when an agent is running in that environment. Docker Desktop forwards the agent through a fixed path, `/run/host-services/ssh-auth.sock`, on Mac and Linux only, not Windows ([Docker Desktop networking how-tos](https://docs.docker.com/desktop/features/networking/networking-how-tos/)); **not tried**. The same passwd-entry problem as the deploy key. A Unix socket path is limited to about 108 bytes, which the first attempt here exceeded.
- **What it lets someone do.** This is the widest of the three. While the container runs, a process in it can sign with **every key in the operator's agent**: every repository of their account, and any server those keys open. It cannot copy the keys, and it cannot be scoped to one repository.
- **Revocation and absence.** Removing keys from the agent cuts the container off at once (**tried**: "Permission denied (publickey)"), and so does the agent ending (**tried**), which is what happens when the operator logs out or the host restarts. A dashboard that stays up longer than the operator's session then fails every push until it is restarted from a session with an agent. It is unattended only while someone is logged in.

## Pushing with the container holding nothing

### A host-side command: push what is pending

- **What it is.** `flai push --pending` (name open), run on the host by whoever is there: the operator, an agent session, or a timer. When `main` is ahead of its remote-tracking branch and the commits ahead include an acceptance, it pushes `main` and the tags that point at unpushed commits, with the host's own credentials, the ones the operator already pushes with. **Tried** as a shell prototype against the scratch remote: after an acceptance made inside a credential-free container (`pushed: false`, `push_error` set), it pushed three commits and the release tag; a second run said "nothing pending"; the container saw zero commits ahead immediately, because it shares the clone.
- **Hosts.** Wherever the operator can `git push` today.
- **What it lets someone do.** The token holder gains nothing at the moment of acceptance. What they accepted is published when someone on the host next runs the command, so a person or an agent stands between the board and the release, if they look.
- **When nobody is there.** Nothing is pushed until somebody or something runs it. A systemd user timer or cron entry makes it unattended, and then it is an unattended push again, with the operator's full credentials rather than a scoped one; that is the operator's choice to make on their own host, outside flai.
- **A failed push** is an ordinary failed command on the host; the board keeps saying "not pushed".

### The agent session pushes

What happens now, by convention and by hand: an agent told of an acceptance fetches, pushes `main` and the tags, and restarts the dashboard. **Tried** ten times on this repository on 2026-09-18 and 2026-09-19, never failing, once after the operator had committed through GitHub in between (merged, not rebased). It depends on an agent running and being told; an agent holding `wait_for_events` sees the acceptance as a change but nothing tells it that the acceptance is unpushed. With the command above and an `unpushed` fact in `inbox`, this becomes a rule an agent can follow without being asked.

### A git hook in the clone

A `post-commit` hook runs where the commit is made. **Tried**: for an acceptance made from the container, the hook ran inside the container (the container's hostname, the host user's UID, no agent socket), where there is no credential, so it cannot push. There is no hook for tag creation short of `reference-transaction`, which runs in the same place. A hook cannot move the push to the host. It lost on that alone; the hook finding above is a second reason to want fewer hooks, not more.

### Not pushing, and making that hard to miss

- **What it is.** The board and the item page show a standing state, "accepted, not pushed: 3 commits, tags flaiover/v0.17.1", with the command and with `flai push --pending` once it exists, until it is no longer true. Today the notice appears once, after the acceptance, and is gone on the next page load.
- **Can a container with no credential know?** **Tried**: offline, from the clone, it can count the commits `main` is ahead of `origin/main`, list them, and list the tags that point at commits not yet on `origin/main`. It cannot ask the remote (the same "could not read Username" failure), so it cannot see a push made from another clone until someone fetches here; after a push from this clone it is right at once.
- **What it lets someone do.** Nothing new.

## The operator's answers

Asked in the session on 2026-09-19, after the research, each question with a recommended answer first. The operator's answers, as given:

| Question | Recommended | The operator's answer |
|----------|-------------|-----------------------|
| Who can reach the dashboard today, and who will through a hub or tunnel? | (none) | "Me, over a public tunnel": only the operator holds the token, and the address will be reachable from the internet, so the token is the only lock |
| A pushed release tag publishes binaries and the image. Is a push nobody on the host approved acceptable? | No: push from the host side | "Yes, every acceptance" |
| If the container were given a credential, which kind? | None | First a question back: "is there a way to side-load/mount the ssh key from the host machine?" Answered: yes, a key file mounts read-only and that is the deploy key mechanism tried above; mounting the operator's own key puts a copyable key for their whole account in the container, a dedicated key for this one repository does not, and a passphrase makes either unusable unattended. Asked again with a dedicated deploy key recommended: "My own SSH key, mounted" |
| When a push fails or has not happened yet, what should the board do? | Standing state | "Standing state (Recommended)": the acceptance stands, and the board and the item page keep showing "accepted, not pushed" with the command until it is no longer true; agents see it in `inbox` |

## Recommendation

The operator wants to drive the project from the board, remotely, and an acceptance that is not published until someone reaches a shell defeats that. Given their answers, the host-side command alone is not the answer for them: it still needs someone on the host. So the recommendation is that **`flai dashboard` can give the container one credential to push acceptances with, as an explicit opt-in that is off by default, and that the board shows a standing "accepted, not pushed" state whether or not a credential is given.**

For the credential the recommendation was a **dedicated deploy key**: a key made for this project's dashboard, added to this one repository with write access, mounted read-only. It fits the repository's SSH remote as it is, and a stolen copy opens one repository. Why the others lost:

- **The operator's own key, mounted.** The same mechanism and the simplest to set up, but the private key is a readable file in a container whose only lock is a token on a public address, and it opens every repository of the account and any server that trusts it. A passphrase, which such a key should have, makes it unusable unattended.
- **The forwarded agent.** Never hands the key over, but lends every key in the agent for as long as the container runs, cannot be scoped, fails whenever the operator is logged out, and does not work on Windows with Docker Desktop.
- **A fine-grained HTTPS token.** The best scoped of all (one repository, contents only, no workflow permission, an expiry) and needs no change to the image, but this repository's remote is SSH, so the push would have to be redirected to an HTTPS URL, and it is one more credential kind for the operator to keep. It is the right second option, and the right first one for a project whose remote is HTTPS.
- **A git hook.** Runs in the container, where there is no credential.
- **The host-side command alone, the agent session alone, or not pushing.** All keep the container empty-handed and all need someone or something on the host. They stay the default for every project that does not opt in, which is why the standing state and the pending-push command are still worth building.

Whatever the credential, a stolen one can publish a release from a commit of the thief's making. Repository rules (no force push, no deletion of branches or tags) and no workflow permission narrow what else it can do; they do not remove that.

## Decision

The operator decided on 2026-09-19, recorded in [ADR-0026](../adrs/0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md):

- Every acceptance made from the board is pushed by the container, when the operator has opted in by giving it a key.
- The key is **the operator's own SSH key, mounted read-only**, not the dedicated key that was recommended. The operator was told what that exposes before choosing. It is their key and their repository, and the decision is theirs. The mechanism is the same for any key file, so a dedicated key remains a matter of which path the operator names, and the documentation will say which is safer.
- The board shows a standing "accepted, not pushed" state; no retry button.

Nothing was built in S-0052. What followed:

1. **S-0062** (built; see [operators](../../docs/operators/index.md#pushing-what-the-board-accepts)) The container pushes an acceptance when the operator has given it a key: the opt-in, the SSH client in the image, the passwd entry, pinned host keys, refusing a passphrase-protected key with the reason, a warning at start that names what the container holds, and the operators' documentation with the safer choice first.
2. **S-0063** A standing "accepted, not pushed" state on the board and the item page, the same fact in `inbox` for agents, and `flai push --pending` on the host for projects that give the container nothing and for pushes that failed.
3. **S-0064** A container cannot leave git hooks or configuration that run on the host: the finding under "What is at stake", which matters most for the default set-up, where the container holds nothing and the hook is the only way to the operator's credentials.
