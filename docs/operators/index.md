---
title: Operators guide
updated: 2026-09-23
status: draft
---

# system-flow for operators

## Running the dashboard

`flai dashboard` makes sure one container, `flaiover`, is running and serves this project. If it is already running, for this project or another, this registers the project with it and starts nothing more; if not, it pulls (or, inside the monorepo with `--build`, builds from `flaiover/` as `flaiover:local`) and starts it, published on every interface at port `4242` by default (`--bind 127.0.0.1`, or `dashboard.bind`, restricts it to this host), and prints a login link; see Authentication below. The container is given that port and two secrets, each one file mounted read-only: the login token and the credential the host's flai proves itself with, both now per user, shared by every project on this host, not generated per project ([ADR-0033](../../design/adrs/0033-one-login-token-and-one-agent-credential-per-user-serve-every-project.md)). **No file of any project is mounted into it** ([ADR-0031](../../design/adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)). Everything the dashboard shows and everything it changes it asks of `flai serve` on this host, which `flai dashboard` registers the project with and starts; see the next section. `flai dashboard status` lists every project it currently serves; `logs` and `stop` manage the container (`stop` only actually stops it once no project is left registered — see below). Change the image, tag, port, or bind address in `~/.flai/config.json` or per project in `system-flow.yaml` under `dashboard`; they matter only the first time, when they decide what the shared container is started with.

Run in a folder that is not a project, such as `~/git`, `flai dashboard` starts the same container and `flai serve`, registers every project it finds in the folder and up to three levels below it, and names the folder for import, so that the folder's other git repositories are offered on the board. `flai serve import remove .` undoes that naming ([ADR-0036](../../design/adrs/0036-a-folder-that-is-not-a-project-is-served-whole-by-flai-mcp-and-flai-dashboard.md)). Its other subcommands (`status`, `stop`, `token`, `restart`, `upgrade`, `check`, `logs`) work in such a folder too. There, `stop` stops the container only once no project is left registered. `flai serve` itself, started in a folder that is not a project (in the foreground, or with `flai serve start`), does the same for as long as it runs: it serves the projects below the folder, without writing them to its registry, and offers the folder's other git repositories for import, as though the folder were named (S-0102).

The container runs as you (`--user`), for one reason: the two secrets are files only you can read. There is nothing of yours within its reach to write.

Without `flai`, the container alone is (paths below are wherever your `flai serve` keeps its state — `flai dashboard status --json` or `flai serve status --json` names it):

```bash
docker run --detach --rm --name flaiover \
  --publish 0.0.0.0:4242:3000 \
  --mount type=bind,source="$HOME/.flai/serve/dashboard.token",target=/run/secrets/flaiover_token,readonly \
  --env FLAIOVER_TOKEN_FILE=/run/secrets/flaiover_token \
  --mount type=bind,source="$HOME/.flai/serve/dashboard.agent-key",target=/run/secrets/flaiover_agent_key,readonly \
  --env FLAIOVER_AGENT_KEY_FILE=/run/secrets/flaiover_agent_key \
  --user "$(id -u):$(id -g)" ghcr.io/bytepunx/flaiover:latest
```

It shows nothing until a `flai serve` that knows a project and holds the same credential connects to it, which is what `flai dashboard` arranges; started by hand, every page says that no flai is connected.

### Upgrading from a release that mounted the repository

Earlier releases mounted the repository read-write into the container, with `.git/hooks`, `.git/config`, and `.git/info` read-only over it, your git identity and global excludes passed in, and optionally an SSH key to push with. All of that is gone.

- **Restart the dashboard** with the new flai: `flai dashboard stop`, then `flai dashboard`. A container started by an older flai keeps its mounts until then, and `flai dashboard status` says so.
- **If you configured a push key** (`dashboard.push_key`, `--push-key`, `dashboard.push_known_hosts`): it is ignored, and `flai dashboard` says so at every start until you clear it with `flai config set dashboard.push_key ""` (and `dashboard.push_known_hosts` likewise). The container has no git or ssh to push with and no repository to push from. An acceptance from the board is committed on the host; it is tagged and pushed from the host with your own credentials: `flai push --pending`, or from the board once you enable the push action, both below. If the key was a deploy key made for this purpose, delete it from the repository's settings and from `~/.ssh`; nothing uses it any more. `.flai-cache/dashboard.known_hosts` and `.flai-cache/dashboard.passwd` can be deleted.
- **Nothing to do for git settings.** Commits from the board are made on the host with your git configuration of the moment; the advice to restart the dashboard after changing git settings no longer applies.
- **A Windows or otherwise unusual host path** no longer matters: git runs where `flai serve` runs, so a story with a branch can be accepted from the board anywhere.

### Upgrading from one container per project

Before this release, `flai dashboard` started a container per project, each with its own login token and its own agent credential under that project's `.flai-cache`. Since S-0080 ([ADR-0033](../../design/adrs/0033-one-login-token-and-one-agent-credential-per-user-serve-every-project.md)) one container, `flaiover`, serves every project registered with the same `flai serve`, with one shared login token and one shared agent credential kept beside `flai serve`'s own state.

- **Move each project over in turn.** In each project, `flai dashboard stop` then `flai dashboard`. The first one you do this in starts the shared container; every one after it just registers with the container the first one started, and starts nothing. Until you do this for a project, it keeps running its own old container and needs its own old login link.
- **Old logins stop working project by project, as you move each one over.** Once a project has moved to the shared container, its old per-project login link and token no longer apply: log in again with the link `flai dashboard` prints, which now opens every project on the host, not just the one it was run in.
- **The old per-project secret files become inert, not migrated.** A project's `.flai-cache/dashboard.token` and `.flai-cache/dashboard.agent-key` from before this release are simply no longer read once that project has moved over; nothing deletes or migrates them. Once every project on the host has moved, they can be deleted (they are already git-ignored, so leaving them costs nothing but disk).
- **The container's name changes.** It was `flaiover-<project>`; it is now `flaiover` for every project. `docker ps`, any firewall rule, reverse proxy, or systemd unit that names the old container by its project-specific name needs updating to the new fixed name.
- **One login for every project on the host** is the point: once every project has moved over, you and anyone you share a link with sign in once, and the project switcher in the header (once more than one project is connected) chooses between them.

### Who acts when the dashboard writes

The dashboard has one login token, shared by every project it serves, and so one holder. What it does for a given project is recorded as that project's manifest's `owner` (`system-flow.yaml`), or `designer` when there is none: thread entries, moves made on the board, and acceptances, which run `flai accept <id> --by <owner>`. All of it is done by `flai serve` on the host, as you: commits carry your git identity and the dashboard's `Co-Authored-By` trailer.

### Accepting computes no release; pushing does, before it pushes (S-0087, S-0094)

A story accepted from the board (or `flai accept`, or `flai move <story> done`) is merged, archived, and committed in your clone by `flai serve`, and that is all: nothing is tagged or pushed there. Releasing what has accumulated happens at the push, covering everything merged since each component's last tag, not one release per story: three small stories against the same component become one release, not three. This was once a step you had to remember to take deliberately, separate from pushing (S-0087); since S-0094 pushing does it for you, first, so a release is never a step someone forgot.

**Pushing, which now also releases.** `flai push --pending`, on the host, tags everything accumulated and unreleased for each component since its last tag — the highest delivery type among what is pending, not a sum — bumps and commits the version files, then pushes the branch and every tag together with your own credentials. `--dry-run` shows the release it would tag, and what it would push, and changes nothing. It never forces, and refuses when the remote has moved until you fetch and merge; run it again after it fails partway (a tag made, the push not) and it does not redo what already succeeded. `--publish` also publishes the template when its version moved. Tags go three to a push, the branch last, because GitHub starts no tag-triggered workflow when one push carries more than three. The board, the story's page, `flai board`, and the agents' MCP `inbox` all keep saying "accepted, not pushed" until it is pushed; an agent session that is running does this itself when `inbox` reports it. From the board, once you enable the push action (below), the standing notice's **Push now** button does the same thing and reports the tags it just made.

**Publishing, separately, when you want to see or force it ahead of a push.** On the host, `flai release --pending` does the same computing, applying, and tagging `flai push --pending` now does, on its own — useful for seeing the release before deciding to push, or for cutting one without pushing yet. `--dry-run` shows what it would do and changes nothing; run again after a partial failure. From the board, once you enable the push action, the done column shows a card as **published** or **waiting to publish**, and a banner at its top lists what publishing now would release — which components, which bump, which stories — with its own **Publish** button; off, the banner says what enables it and changes nothing; the dashboard cannot enable it itself.

### The push host action

Both publishing and pushing ordinary merged commits are a *host action*: something `flai serve` does on your machine, as you and with your credentials, because the dashboard asked. It is off until you enable it by name, in a shell on the host:

```bash
flai serve actions          # what there is, what each means, where each is on
flai serve enable push      # for the project in the working directory; --all-projects for every project
flai serve journal          # every host action asked for, refused ones included
flai serve disable push     # --all-projects turns it off everywhere
```

Enabling and disabling take effect at once, with no restart. The setting is `host_actions` in your flai configuration (`~/.flai/config.json`); `flai config set` does not reach it, and nothing the dashboard can ask for reads or changes it or the journal.

**Understand what enabling it means.** The dashboard token becomes the power to publish: whoever holds it can publish any release accumulated so far, and push what any story an agent has put in review adds once accepted, with your git credentials and with nobody at the keyboard. A compromised dashboard container could ask for the same. If that is more than you want a token to be worth, leave it off and publish and push by hand, or from a timer of your own.

With it on, the done column's Publish button and the standing "accepted, not pushed" notice's **Push now** button both work; the confirmation for each says what it will do before you click it. Both are refused (HTTP 403) and journalled while it is off, and say what enables it. Neither ever forces a push, and when the remote has moved they are refused with the reason and the command to run by hand.

The journal is `journal.jsonl` beside `flai serve`'s state (the `serve` folder next to your flai configuration), mode 0600, one line per request: when, the action, the method, the project, for whom (the manifest's `owner`), the request, and the outcome (`done`, `failed`, or `disabled`) with what was pushed and published or why not.

> **Correction, 2026-09-20.** From flai 1.5.3 to 1.6.1, an acceptance made from the board was pushed and published with your credentials although nothing had enabled it (I-0028). Acceptance had moved from the container, which held no credential, to `flai serve` on the host, which has yours, and the push was not turned off on the way. This page said during that time that nothing was pushed unasked and that the token was not the power to publish; both were wrong. If you accepted from the board with one of those releases, check what was pushed against what you meant to release. The release after 1.6.1 restores the default described above. Since S-0087, acceptance itself does not push at all, whatever this setting is: only publishing and the standing push notice do.

### Starting an agent when a story becomes ready

`flai serve` can start an agent for each story that becomes ready, so that a move to ready on the board starts work with no one at a keyboard. Each story names its agent: a harness, a model, and options (`flai agent set` gives new stories a default; `flai edit` or the dashboard changes one story's). flai starts a harness it knows through its adapter, with the program and the arguments you set for it on this host. For a story that names no harness, it runs a command you wrote, if you set one. It is a host action, off until you enable it ([ADR-0038](../../design/adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md)).

```bash
flai serve enable agent        # for this project; --all-projects for every project
flai serve agent harness       # each harness: the program and what its agent may do here
flai serve agent harness claude-code --program /opt/claude/bin/claude
flai serve agent harness claude-code -- --permission-mode acceptEdits --allowedTools Bash,mcp__flai
flai serve agent harness claude-code --reset
flai serve agent set --name builder -- /home/me/bin/start-agent {story} --model {model}
flai serve agent show
flai serve journal             # every start and every failure
flai serve disable agent
```

**Understand what enabling it means.** Whoever can move a story to ready starts an agent on your machine, as you: you at the board, an agent with `flai move`, and anyone who holds the dashboard token. So can whoever can edit a story, because a story chooses its harness and its model. A compromised dashboard container could do both. What a story cannot choose is what runs or what the agent may do.

- **Harnesses.** `claude-code` runs `claude -p` headless, in the project's directory, with a prompt that tells it to work that story, and only that story, to review the way the project's conventions say. It gets the story's `--model` and flai's MCP server (this flai's `flai mcp`), and no other MCP server unless you add one with `--mcp-config` in its arguments. The only options a story may give it are `effort` (low, medium, high, xhigh, max), `max_budget_usd`, and `fallback_model`, and each is checked. Anything else refuses the start, and the story's card says why. By default an agent may edit files in the project, run any shell command (git, flai, your tests), and use flai's tools: `--permission-mode acceptEdits --allowedTools Bash,mcp__flai`. `flai serve agent harness claude-code -- ...` replaces those arguments, and `--program` names the binary. Both are stored in your flai configuration (`agent.harnesses` in `~/.flai/config.json`), which `flai config set` does not reach.
- **Your command.** A story that names no harness, or names `command`, is started with the argument list you set with `flai serve agent set -- ...`. It is run as it stands, never through a shell. `{story}`, `{root}`, `{model}`, and `{harness}` in an argument are replaced, and the story's options are in `FLAI_AGENT_CONFIG` as a JSON object. With no command set, such a story is not started, and the board says so.
- **The environment** carries `FLAI_AGENT` (your `--name`, default `agent`, and the story: `agent-S-0104`), `FLAI_STORY`, `FLAI_SESSION`, and `FLAI_STARTED_BY=flai-serve`, with `FLAI_MODEL` and `FLAI_HARNESS` for your command. Each agent has a name of its own, so its inbox and its place in `wait_for_work` are its own. The story's ID is flai's own reading of the project's files, checked to be a story's ID. Nothing the dashboard sends becomes part of a command.
- **When one starts.** One agent per ready story, in your pull order, while the in-progress limit leaves room. An agent started for a story that is still in ready counts as a story in progress. A story is started once each time it enters ready, from the board, the CLI, or an agent: one whose agent failed waits until it is moved to ready again. Nothing starts while someone else is attending the project.
- **Attending** is judged from files, because flai asks nobody. An agent connected over MCP rewrites its read marker under `.flai-cache/mcp/` at every look, at least every five minutes while it waits. An agent at work writes its story's narrative under `wip/agents/`. If the newest of those is younger than `--attended-minutes` (6 by default), someone is attending and will see the story in their inbox. The marker and narrative of an agent flai started do not count, and a log entry you add to a narrative from the dashboard does, for those minutes.
- **What does not start one.** Stories that were already ready when `flai serve` began: a restart of `flai serve` never starts a session. An agent that outlived a restart is still watched, and is settled at the next look after it ends. Only projects `flai serve` serves start agents, that is, ones whose dashboard was started with `flai dashboard`.
- **How it ends.** An agent that ends with its story in review or done has worked. One that ends with its story anywhere else has failed: the story's card and page say where it left the story and whether it was blocked, and how it exited. An agent that asks you something opens a thread on its story and waits for the answer, which you give on the story's page or in the inbox.
- **Where to look.** Each story's card on the board has a dot while its agent runs: green at work, yellow waiting for you (it asked on the story, or the story is blocked), red when it failed. The story's page says which harness and model, since when, the question it waits on, or why it failed and where its log is. `flai serve journal` has every start and failure, and each agent's output is in a log under `serve/agents/` beside `flai serve`'s state. Stopping an agent is yours to do on the host (its PID is in the journal); the dashboard cannot stop, start, or configure one.

### The dashboard host action (S-0081)

`flai serve` can restart, upgrade, or stop the dashboard container itself, on your machine, with Docker, because the dashboard's own page asked. It is a host action, off until you enable it:

```bash
flai serve enable dashboard      # for the project in the working directory; --all-projects for every project
flai serve journal                # every restart, upgrade, and stop, refused ones included
flai serve disable dashboard
```

**Understand what enabling it means.** Whoever holds the dashboard token can restart the container, pull and switch it to whatever image your configuration currently names, or stop it (for this project, or for every project it serves, if it was the last one registered). None of it reaches past Docker on this host: no git credential, no project file.

**An upgrade never touches the running container until the new one has proven itself.** `flai dashboard upgrade` pulls the configured image and, if it differs from what is running, starts it as a second, temporary container of its own and waits for it to answer healthy; only then does it stop the running container and start the new image in its place. If the new image never answers healthy, the temporary container is removed and the one you already had keeps running, unchanged — the same command run by hand (`flai dashboard upgrade`, `flai dashboard check` to look without changing anything, `flai dashboard restart` to cycle the process without a pull) behaves identically, on or off the board.

### The checks host action (S-0082)

`flai serve` can run checks for a story in review, in the story's worktree, on your machine, because the review page asked. It is a host action, off until you enable it:

```bash
flai serve checks set --name flai -- scripts/flai-test.sh
flai serve checks set --name flaiover -- bash -c "cd flaiover && pnpm run check && pnpm run test:unit -- run && pnpm run lint"
flai serve checks show
flai serve checks timeout 20      # minutes for one run of every command together; 15 by default
flai serve enable checks          # for the project in the working directory; --all-projects for every project
flai serve journal                # every run and cancel, refused ones included
flai serve checks clear flaiover  # one name, or every name with no argument
```

Naming nothing here leaves the manifest's own `checks:` in effect instead — the project's own default, committed and visible to everyone. Naming any command here, on this host, uses this list instead, for this host only: your own real choice of where to name them, never a merge of both.

**Understand what enabling it means.** Whoever holds the dashboard token can run every named check, in order, in a story's worktree, on your machine, as you, and cancel a run early. Each command is an argument list, run as it stands, never through a shell; `{story}` and `{root}` are replaced in an argument, nothing else is interpreted. The first command to fail stops the rest; a run is bounded by the timeout above, and cancelling it (terminate, then, after a short grace period, kill) reaches the whole process tree the command started, not only its direct child. One run per story at a time: a second while one is active is refused, not queued. The outcome and duration stay with the story (`flai checks status <id>`) until you accept it.

It can reach its port, the network, and the two secrets. It cannot read or write any file of the project or of the host: no work tree, no `.git`, no `.flai-cache` beyond those two files, no flai configuration. So it cannot leave a git hook or setting that would run on your machine, change a tracked file or a branch, or plant an ignored file your tools execute, which were the routes open or guarded while the repository was mounted (I-0022, ADR-0027, superseded).

What a compromised container could still do is what the dashboard itself does: ask `flai serve` for the named methods it offers (reads of the three folders, moves, saves of Markdown under them, an acceptance of a story in review, a push if you enabled that host action, and, by moving a story to ready, a start of your agent command if you enabled that one), each of which flai checks and performs itself, and present the login token. Treat a dashboard you expose beyond your own network accordingly, and see the next section for how to turn the connection off.

| Setting | Where | Default | Effect |
|---------|-------|---------|--------|
| `dashboard.notify_url` | `system-flow.yaml`, per project | unset | A webhook for the designer's inbox. When it is an `http` or `https` URL, the dashboard's server POSTs `{ "project": "<name>", "entry": { "key", "kind", "title", "href", "at" } }` as JSON for each inbox entry that appears after the server started: `kind` is `thread`, `question`, `review`, `blocked`, or `overlap`, and `href` is a path in the dashboard. One attempt per entry, a five second timeout, no retry; a failure is a warning in the log naming the host only, since the URL may carry a secret. The project token, file contents, and anything else are never sent. Entries that existed at start are not posted, so a restart does not replay the inbox. Read once at start: restart the dashboard after changing it |
| `dashboard.autocommit` | `system-flow.yaml`, per project | `true` | Documents saved from the dashboard's editor are committed on the main checkout by `flai serve`, one path per commit, authored by your git identity, with a `Co-Authored-By: flaiover` trailer. `false` leaves them uncommitted: agents on story branches then do not see the edit until someone commits it, and an acceptance from the board lists it as an uncommitted change. Nothing is pushed ([ADR-0023](../../design/adrs/0023-documents-are-saved-through-flai.md)) |
| `worktrees.relative_paths` | `~/.flai/config.json`, per user (`flai config set`) | `false` | With git 2.48 or newer, `flai stream open` links new worktrees with relative paths, for clones that are moved or shared. The dashboard no longer has any use for it. Sets `extensions.relativeWorktrees` on the clone, after which git older than 2.48 refuses the repository. Never enabled automatically. How to turn it back off is in [the flai guide](../users/flai.md) |

## The connection from flai on the host

`flai dashboard` starts a second thing beside the container: `flai serve`, one process per user on the host, which opens a WebSocket to the dashboard's `/agent` endpoint and keeps it open. Over it the dashboard asks flai for named things and flai answers. The direction is deliberate: the container is given no socket, pipe, or address of the host's, and cannot start a conversation with it. It is the only way the dashboard reaches the project: everything it reads and everything it changes goes over it, and the container holds no file of the project ([ADR-0031](../../design/adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)). The image holds no `flai`, no `git`, and no `ssh`; commits made from the board are made by `flai serve` on the host as you, with your git configuration, and carry the dashboard's trailer.

- **The credential.** `flai dashboard` creates one agent credential per user, kept beside `flai serve`'s own state, not any one project's `.flai-cache` (S-0080, [ADR-0033](../../design/adrs/0033-one-login-token-and-one-agent-credential-per-user-serve-every-project.md)), and mounts it read-only into the shared container at `/run/secrets/flaiover_agent_key`. It is not the login token and opens nothing but `/agent`. It never travels: `flai serve` sends a nonce, the dashboard answers with its own nonce and an HMAC of both under the credential, and `flai serve` answers with the mirror HMAC. A flai that dials a port something else is listening on gives nothing away and says so in `flai serve status`; a client that cannot prove the credential is dropped before it is served anything. Sharing the credential across projects proves who `flai serve` is, not which project a given connection may answer for: every request still names its project, and flai refuses one that does not match what that connection's `hello` declared.
- **Who can reach `/agent`.** Anyone who can reach the dashboard's port can open the socket, and gets nothing without the credential. A request with an `Origin` header is refused outright, so no page in a browser can try. One connection per project per credential: a newer proven connection for the same project replaces the older; a different project registered with the same `flai serve` opens its own connection alongside it.
- **What can be asked.** Only the methods `flai serve` offers (`flai hostapi` lists them), each checked against the project that connection's `hello` named. There is no method that takes a command line or a file path of the caller's choosing, and the dashboard cannot add one. What a holder of the dashboard token can do is what the dashboard's pages can do, as before; what a compromised container can do is call those same methods and nothing more: it has no file of any project, no `git`, and no credential. Whatever it sends is validated as data (an item ID, a state, a Markdown path under the project's three folders, text) and reaches a command only as `--flag=value` or after `--`. It cannot choose who a change is recorded as.
- **When it is not there.** Pings every 4 seconds in both directions; a `flai serve` that is frozen or gone is marked so within ten seconds, a restarted container is reconnected within a second, and requests in flight fail with a clear error. Tried on Linux under WSL2 with Docker Engine. On macOS and Windows with Docker Desktop nothing differs in principle, since flai dials the published port on the host as a browser does, but it has **not been tried** there; on Windows `flai serve` runs in its own process group and `flai serve stop` ends it without a graceful signal.
- **Where its files are.** A folder named `serve` beside flai's config file: `projects.json` (what it serves), `state.json` (rewritten every second while it runs), `serve.log`, and, since S-0080, the login token and the agent credential themselves — nothing in the folder is secret except those two files (mode 0600).
- **Turning it off.** `flai dashboard --no-serve` starts the container without registering the project. The dashboard then has no flai to ask for that project: its board, items, and threads answer 503, every page says why, and `/_ready` reports not ready (its `host_flai` check), which matters if you probe it. Use it when you run `flai serve` yourself in a terminal.
- **What flai serve reads.** Each registered project's design, docs, and wip folders and its manifest, looked at every 300 ms for changes (a stat of each file; a few hundred files cost well under a millisecond of CPU a look), and, when the dashboard asks, Markdown files under those three folders and nothing else: a path that climbs out, a link that leads out, another folder, or another file type is refused by flai, whatever the dashboard sends. The search index is built from the same files and kept in `flai serve`'s memory.
- **Stopping the container.** The image's server now closes every connection five seconds after SIGTERM (`SHUTDOWN_TIMEOUT`, in seconds), so a container asked to stop does stop, even with event streams open; before, a server that had stopped listening could linger while its container showed as up.

### Importing repositories from the board (S-0098)

Name, on the host, a folder that holds git repositories which are not system-flow projects yet, and the board offers each of them for import ([ADR-0035](../../design/adrs/0035-repositories-under-folders-the-operator-names-can-be-imported-from-the-board.md)):

```bash
flai serve import add ~/git     # offer the repositories in it
flai serve import list          # the folders named, and what is found in them now
flai serve import remove ~/git
```

- **What is offered.** Every git repository in the folder, or up to three folders below it, with no `system-flow.yaml`. Nothing inside a repository, a hidden folder, or build output is looked at. `flai serve` looks again every 30 seconds, and `flai serve status` lists what it offers. A repository is offered to the dashboards the projects `flai serve` already serves connect to, so with no project served, nothing is offered.
- **What naming a folder allows.** Anyone who can use the board can import any repository offered. An import writes the standard's files into the repository, runs the repository's own tests on this host as you (the checks `flai serve checks set` names, if any, else what the repository has: its Makefile's `test` target, `go test`, its package manager's test script, `cargo test`, pytest), and commits the import when they pass. Naming the folder is your say that this may happen to what is in it: no host action needs enabling as well. Name only folders whose repositories you would run tests from.
- **What an import does.** It is `flai import --yes --commit` in the repository. A repository with uncommitted changes is refused, so the commit holds the import alone. When a test fails nothing is committed and the files stay for you to fix. Either way the repository is registered with `flai serve` and served as a project from then on, under a key made from its folder's name. Every import is in the journal (`flai serve journal`, action `import`).

## MCP over HTTP

MCP is served by flai on the host, not by the dashboard ([ADR-0030](../../design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)). An agent on the host needs nothing from you: `.mcp.json` starts `flai mcp` on stdio. For an agent that cannot start a process there, the same server runs over HTTP, one per project.

`flai serve` keeps it running for every project it serves ([ADR-0034](../../design/adrs/0034-flai-serve-keeps-each-served-project-s-http-mcp-server-running.md)): it starts one when it starts serving a project, starts it again if it stops, and stops it when the project is no longer served or `flai serve` stops, even when `flai serve` was killed. One you started yourself with `flai mcp start` is used as it is and never stopped by `flai serve`. `flai serve status` shows where each project's server listens. The first project gets `127.0.0.1:4243`, the next the next free port, and each keeps its port from then on. For a project no `flai serve` serves, start it by hand:

```bash
flai mcp start      # in the background; flai mcp http runs it in the foreground
flai mcp status     # where it listens, and what an agent's configuration looks like
flai mcp token      # its bearer token; --rotate replaces it and restarts the server
flai mcp stop
```

| | |
|-|-|
| Address | `http://127.0.0.1:4243/mcp` unless `--addr` says otherwise. The address is remembered per project (`.flai-cache/mcp-http.addr`), so an agent's configuration survives a restart; a second project on the same machine needs another port, which `flai serve` picks by itself (the next free one from 4243) |
| Transport | MCP Streamable HTTP, answers as `application/json`. Clients on revisions up to 2025-11-25 get a session (`Mcp-Session-Id`), ended by `DELETE` or after `--idle` without a request (default 30 minutes), at most `--max-sessions` at once (default 16), after which `initialize` answers 503. Clients on 2026-07-28, which has no sessions, are served request by request at the same address |
| Authentication | `Authorization: Bearer <token>` only, from `.flai-cache/mcp.token` (mode 0600, git-ignored, created when first needed). It is not the dashboard's token, and the dashboard's token does not open it. A request with an `Origin` header, which is what a browser sends, is refused with 403 |
| The agent's name | The `X-Flai-Agent` header, else the client's own name, made safe for a file name. It is who thread entries and transitions are attributed to, and whose cursor `inbox` keeps |
| Long requests | `wait_for_events` holds its request open until something changes, for up to five minutes. A proxy or tunnel in front must allow an idle response that long, or agents see their wait cut short. At most 64 requests are in flight at once; more answer 503 |
| State and log | `.flai-cache/mcp-http.json` while it runs, `.flai-cache/mcp-http.log` for its events (`mcp server started`, `mcp session requested` with the agent and how many sessions are open, `mcp server stopped`) |
| What it serves | The main checkout, read when it starts: restart it after changing `system-flow.yaml` |
| Stopping | `flai mcp stop` ends held `wait_for_events` calls as waits that ran out, so a waiting agent gets an answer, not a broken connection, and simply starts a session again when the server is back. On a server `flai serve` keeps, `flai mcp stop` lasts only until its next look (within 15 seconds), when it is started again; to stop it for good, stop serving the project (`flai dashboard stop`) or stop `flai serve` |

It listens on this machine only by default. `--addr 0.0.0.0:4243` or another interface is allowed and logged as a warning: the token travels in every request and flai does not encrypt it, so beyond the machine put an SSH forward, a tunnel, or a proxy that terminates TLS in front and give agents that address.

**The dashboard's `/mcp` is gone.** Until flaiover 0.22 the dashboard served MCP itself (ADR-0024, superseded). Any request to `/mcp` now answers 410 with a message naming `flai mcp start`, so an agent configured for the old address is told what to do. `FLAIOVER_MCP_MAX_SESSIONS` and `FLAIOVER_MCP_IDLE_MINUTES` no longer mean anything.

### The tunnel expectation

The dashboard speaks plain HTTP and the token travels in every request. On the machine itself or a network you trust that is acceptable. Anywhere else, put a tunnel or a reverse proxy that terminates TLS in front of it and use the `https` address; never publish the dashboard's port to the internet as it is. To keep it to the host and let only the tunnel reach it, bind it to loopback: `dashboard.bind: 127.0.0.1` in `system-flow.yaml`, or `flai dashboard --bind 127.0.0.1`. The same holds for `flai mcp` over HTTP, which has its own port and token.

### How the dashboard knows whether it is served over HTTPS

It matters for one thing: the browser's session cookie is marked `Secure` exactly when the page is served over HTTPS, and a browser drops a `Secure` cookie it receives over plain HTTP from anywhere but `localhost`.

- **Directly, over plain HTTP** (by `localhost`, a LAN address, a VPN address): nothing to configure. The dashboard sees that the connection is plain and sets an ordinary cookie.
- **Behind a proxy or tunnel that terminates TLS**: it must tell the dashboard the original scheme in `X-Forwarded-Proto`. cloudflared, Caddy, and Traefik send it unasked; nginx needs `proxy_set_header X-Forwarded-Proto $scheme;`. The dashboard takes the first value of the header, accepts only `http` or `https`, and otherwise goes by the connection. Without the header the cookie is not `Secure` although the page is HTTPS, which works and is weaker than it should be.
- A client that sends `X-Forwarded-Proto: https` by itself over plain HTTP affects only its own session: its cookie is marked `Secure` and its browser drops it.

Until flaiover 0.22.3 the dashboard took every request for HTTPS, so logging in from any address but `localhost` over plain HTTP came straight back to the login page with no reason given (S-0083). The login page now checks that the session was kept and says so when it was not. Bearer tokens were never affected.

### Project identity

Every `/api/*` response of the dashboard, and every answer of `flai mcp` over HTTP, names the project it came from: the headers `X-Flai-Project-Key` and `X-Flai-Project-Name` (URI-encoded), and `project: { name, key }` in JSON object bodies. Both come from `name` and `key` in `system-flow.yaml`; `flai check` warns when `key` is missing. The dashboard's refusals for lack of a token carry neither; flai's MCP server sends the headers on refusals too, since its caller chose the project by choosing the address.

## Authentication

Every request except `/_health` and `/_ready` needs the token (ADR-0018), one per user since S-0080, shared by every project the dashboard serves (ADR-0033).

| Step | How |
|------|-----|
| Set | `flai dashboard` generates it once, on first use: 32 random bytes, base64url. `flai dashboard token` prints it; `--rotate` replaces it and restarts the shared container, ending every session for every project it serves |
| Stored | Beside `flai serve`'s own state (the `serve` folder next to your flai configuration), mode 0600, shared by every project registered with that `flai serve`. That folder must be outside version control; `flai check` still refuses a repository where any project's old `.flai-cache/` is not git-ignored |
| Handed to the container | Bind-mounted read-only at `/run/secrets/flaiover_token` with `FLAIOVER_TOKEN_FILE` pointing at it. Never an environment variable, so `docker inspect` does not show it |
| Browser | Open the login link (`http://host:4242/login#token=...`). The fragment never leaves the browser; the page exchanges it for an HttpOnly, SameSite=Lax cookie (Secure over HTTPS, which the dashboard knows as described under the tunnel expectation) and rewrites history. Pasting the token on `/login` also works |
| Agents and tools | `Authorization: Bearer <token>` on every request |
| Metrics | `/metrics` needs the token unless `FLAIOVER_METRICS_PUBLIC=true` |
| Development | `scripts/flaiover-dev.sh` points the dev server at the same file; `FLAIOVER_AUTH=off` disables authentication outside production only |

Exposure, plainly:

| Risk | Answer |
|------|--------|
| Sniffed on the network | flaiover speaks plain HTTP. Beyond a trusted LAN, terminate TLS in a tunnel or proxy; the cookie is then marked Secure |
| Token in logs | Request logs carry method, route, status, and duration only; never header or cookie values |
| Token in the image or repository | Never baked in, never committed; verified by `flai check` |
| Guessing | 256 bits of randomness |
| Stolen cookie | HttpOnly blocks script access; rotation invalidates every session |
| Docker socket access on the host | Can read the mounted file; that is host ownership, not something the token scheme addresses |

## Access to the image

`ghcr.io/bytepunx/flaiover` is a private package while the repository is private, so a plain `docker pull` is refused. Three ways in:

| Path | What to do |
|------|------------|
| Let `flai` log in | Have `GITHUB_TOKEN` or `GH_TOKEN` set, or be logged in with `gh`, and give the token the `read:packages` scope: `gh auth refresh -h github.com -s read:packages`. On a refused pull `flai dashboard` runs `docker login ghcr.io` with that token and retries once. The token never appears in arguments or logs |
| Log in yourself | `gh auth token \| docker login ghcr.io --username <github-user> --password-stdin`, then `flai dashboard` or the `docker run` above |
| Build locally | In the monorepo, `flai dashboard --build` (or `make flaiover-image` then `flai dashboard --image flaiover --tag local`) builds `flaiover:local` from `flaiover/Dockerfile` and runs it; no registry access needed |

Making the package public on GHCR removes the need for a token entirely; that is an organisation setting, not something `flai` changes.

## The image

`ghcr.io/bytepunx/flaiover` is built from `flaiover/Dockerfile` at the repository root: the SvelteKit build on `node:24-alpine` and nothing else (no `flai`, no `git`, no `ssh`; the label `dev.system-flow.flai` names the flai of the commit it was built and tested with), listening on `3000`, running as an unprivileged user by default and working under any `--user`. Tags: `latest` (main), `X.Y.Z` and `X` from `flaiover/vX.Y.Z` release tags, and `sha-<commit>`. While the repository is private the package is too: `docker login ghcr.io` with a token that has `read:packages` before `flai dashboard` can pull it. The container writes nothing that outlives it: it has no volume, and its only home is `/tmp`.

Build locally with `make flaiover-image` (tag `flaiover:local`) and run it with `flai dashboard --image flaiover --tag local`.

## Telemetry

flaiover follows the logging and telemetry conventions.

| Signal | Where | Notes |
|--------|-------|-------|
| Logs | stdout, one JSON event per line | `ts`, `level`, `service`, `component`, `msg`, and fields; one line per request with `trace_id` (or `request_id` when tracing is off), `method`, `route`, `path`, `status`, `duration_ms`; `LOG_LEVEL` (`debug`, `info`, `warn`, `error`), `LOG_FORMAT=text` for key-value text outside production |
| Liveness | `GET /_health` | Always 200 while the process runs; touches nothing |
| Readiness | `GET /_ready` | Checks that flai on the host is connected (`host_flai`) and, through it, the project's manifest and item listing, with a two second timeout each; 503 names the failing check |
| Metrics | `GET /metrics` | Prometheus format: `flaiover_http_requests_total{method,route,status}`, `flaiover_http_request_duration_seconds` (histogram), `flaiover_http_requests_in_flight`, `flaiover_build_info{version,commit}` (the flaiover release tag and commit baked into the image), plus Node process metrics prefixed `flaiover_`. Route labels are SvelteKit route ids, never paths with IDs |
| Traces | OTLP/HTTP | Exported only when `OTEL_EXPORTER_OTLP_ENDPOINT` is set; one server span per request named `<method> <path>` with W3C context taken from the incoming headers; standard `OTEL_*` variables apply (`OTEL_SERVICE_NAME`, `OTEL_TRACES_SAMPLER`, `OTEL_TRACES_SAMPLER_ARG`) |

Local stack: `PROJECT=$PWD docker compose -f flaiover/compose.yaml up --build` runs the dashboard on <http://localhost:4242> next to `grafana/otel-lgtm` (collector, Tempo, Prometheus, Loki, Grafana on <http://localhost:3001>), with traces exported to it by default.

## Security posture

The dashboard authenticates every request with the project token (above), and through `flai serve` it can change the project: moves, saves, acceptances. By default it is published on every interface of the host so a team can reach it over a private network or VPN; beyond a trusted LAN put a TLS-terminating tunnel or proxy in front of it, because the token travels in clear over plain HTTP. To keep it to the machine it runs on, set `dashboard.bind: 127.0.0.1` in `system-flow.yaml` or config, or pass `--bind 127.0.0.1`. The container holds no file of the project and no git credential. Whether the token is the power to publish is yours to decide: it is not unless you enable the push action on the host (above), and it is once you do. Keep the dashboard off public addresses or behind a tunnel you trust all the same, and rotate the token (`flai dashboard token --rotate`) when in doubt.

## Requirements

- Docker Engine 24 or newer on `PATH`
- `git` on `PATH` for template cloning
