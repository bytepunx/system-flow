---
title: flai CLI
updated: 2026-10-01
status: active
---

# flai

The system-flow command line tool. This guide says how the commands fit together; [flai-reference.md](flai-reference.md) lists every command and flag, generated from `flai --help`. The design is in [design/system/flai-cli.md](../../design/system/flai-cli.md).

## Install

One line on macOS or Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/bytepunx/system-flow/main/install.sh | sh
```

The script detects the OS and architecture, resolves the newest `flai/v*` release, downloads the archive and `checksums.txt`, verifies the SHA-256, and installs `flai` into `$HOME/.flai/bin`, a directory you already own, alongside flai's own config and cache under `~/.flai`. No `sudo` is ever used. It ends by printing where it installed and, if `$HOME/.flai/bin` is not already on your `PATH`, a line to add to your shell profile.

| Variable | Effect |
|----------|--------|
| `FLAI_INSTALL_DIR` | Install somewhere else, for example `/usr/local/bin` (needs write access there already; the script does not use `sudo`) |
| `FLAI_VERSION` | Pin a release, for example `1.0.3` |
| `GITHUB_TOKEN` or `GH_TOKEN` | Authenticate against the GitHub API. While the repository is private one of these, or a `gh auth login` session, is required; the script borrows `gh auth token` when it can |

### Upgrade

```bash
flai self-upgrade --check      # installed and latest versions
flai self-upgrade              # replace this binary with the latest release
flai self-upgrade --version 1.0.3
flai self-upgrade --dir /some/other/directory
```

`self-upgrade` performs the same steps as the script from inside the binary: resolve, download, verify, and replace the running executable atomically. It uses the same token sources. Without `--version` it does nothing when the installed version is already the latest. With no `--dir` it replaces whichever binary is running, so once installed under `~/.flai/bin` an agent or a scheduled job can run `flai self-upgrade` itself, with no `sudo` and no extra configuration. The exception is a flai that runs from inside a system-flow project, such as a checkout's own `bin/flai`. It installs the release where `install.sh` would (`$FLAI_INSTALL_DIR`, else `~/.flai/bin`), leaves the checkout alone, and says to put that folder on your PATH ahead of the checkout's. A build there is replaced by the next build, and that folder exists on no other machine.

`FLAI_RELEASES_API` points it at another source of releases that answers as GitHub's releases API does, such as a mirror. `flai host check` and `flai host upgrade`, and the dashboard's Check for upgrade and Upgrade buttons, which ask the host, follow it too.

### Other ways

With Go 1.26 or newer, `go install github.com/bytepunx/system-flow/flai@latest`. Or download an archive from the [releases page](https://github.com/bytepunx/system-flow/releases): releases named `flai vX.Y.Z` carry `flai_X.Y.Z_<os>_<arch>.tar.gz` (a zip on Windows) for Linux, macOS, and Windows on amd64 and arm64, plus `checksums.txt`. On Windows unpack the zip and put `flai.exe` on your `PATH`.

```bash
flai version
```

## Global flags

| Flag | Effect |
|------|--------|
| `--config <path>` | Use this config file. Falls back to `$FLAI_CONFIG`, then `~/.flai/config.json`. |
| `--json` | Structured output for scripts and agents. |
| `--yes`, `-y` | Answer yes to confirmations. |
| `--verbose`, `-v` | Debug-level log events on stderr. |

## Logging

Command output goes to stdout. Everything else is a structured log event on stderr, one per line, with `ts`, `level`, `component`, `msg`, and named fields. On a terminal events are key-value text; when stderr is redirected they are JSON. Levels are `DEBUG`, `INFO`, `WARN`, `ERROR`, and `FATAL`; a failed command logs one `FATAL` event with an `err` field and exits non-zero.

| Control | Effect |
|---------|--------|
| `--verbose` | Debug level, overrides `LOG_LEVEL` |
| `LOG_LEVEL` | `debug`, `info` (default), `warn`, `error`, `fatal` |
| `LOG_FORMAT` | `text` or `json`; default text on a terminal, json otherwise |
| `FLAI_SLOW_REQUEST` | How long a request to `flai serve` or `flai mcp` may take before its `request answered` event is logged at `INFO` rather than `DEBUG`; a Go duration, default `500ms` |

```bash
flai board 2>/dev/null                 # output only
LOG_FORMAT=json flai check 2>events.jsonl
```

### Request timing (S-0152)

`flai serve` and `flai mcp` log one `request answered` event, component `perf`, for every request they answer: a dashboard's method over the channel, or an MCP request, named by its tool for a tool call. It is measured inside flai, from when the request was read to when its answer is ready to write, so a request's duration in the dashboard's own log, less this one, is the transport.

| Field | Holds |
|-------|-------|
| `transport` | `channel` (the dashboard) or `mcp` |
| `method` | The method, or the MCP tool |
| `project` | The project's key, over the channel |
| `duration_ms` | The time flai took to answer |
| `bytes` | The answer's size: the JSON result over the channel, a tool's text over MCP |
| `phases` | Where the time went, longest first, as `name=milliseconds`, with `xN` for a step taken N times; phases nest and overlap, so they need not add up |
| `err` | The error answered, if any |

A request that takes `FLAI_SLOW_REQUEST` or longer is logged at `INFO`; the rest at `DEBUG`. `wait_for_work` and `wait_for_events` wait by design and are logged at `DEBUG` however long they take.

Phases are named for what the time went to: `repo.open`, `repo.list` (every work item, the archive included when the method asks for it), `repo.get`, `board.load`, `board.view`, `release.pending`, `pending.detect`, `threads.read`, `threads.view`, `narratives.read`, `check.overlap` (the designer's inbox's overlapping touches), `docs.walk`, `doc.read`, `adrs.read`, `search.index`, `search.query`, `agent.state`, `agent.stream`, `manifest.load`, `changes.read`, `item.show`, `cancel.preview`, `accept.preview`, `stream.diff`, `stats.compute`, `push.preview`, `encode`; and `exec.<program>.<command>` for each process flai started to answer, such as `exec.git.log`, or `exec.flai.move` for a write, which runs the flai command. Since S-0159 no read starts flai.

To time one method with no transport at all, ask it of `flai hostapi` with `--timing`; its event goes to stderr at `INFO` whatever it took:

```bash
flai hostapi --timing board.get '{"all":true}' >/dev/null
```

To see where `flai serve` spends its CPU, start it with `FLAI_PPROF_ADDR` set to a loopback address and read a profile with `go tool pprof`:

```bash
FLAI_PPROF_ADDR=127.0.0.1:6060 flai serve
go tool pprof 'http://127.0.0.1:6060/debug/pprof/profile?seconds=30'
```

## Configuration

`flai` keeps its settings in one JSON file per user: `~/.flai/config.json`, or the path in `--config`, or else `FLAI_CONFIG`. The first command that needs it creates it with these defaults:

```json
{
  "template": {
    "repo": "https://github.com/bytepunx/system-flow-template",
    "ref": "main"
  },
  "dashboard": {
    "image": "ghcr.io/bytepunx/flaiover",
    "tag": "latest",
    "port": 4242,
    "bind": "0.0.0.0",
    "push_key": "",
    "push_known_hosts": ""
  },
  "cache_dir": "~/.flai/cache",
  "author": "<your username>",
  "worktrees": {
    "relative_paths": false
  }
}
```

A key flai does not know is refused when the file is read, so a typo is an error rather than a silent default. `flai host` and `flai serve` keep their state, tokens, and logs in the folders `host` and `serve` beside this file, so a different `FLAI_CONFIG` gives them a separate home too.

Read and change it by dotted key:

```bash
flai config get
flai config get template.repo
flai config set template.repo git@github.com:me/system-flow-template.git
flai config set template.ref my-branch
flai config set dashboard.port 8080
flai config set worktrees.relative_paths true   # opt in to relative worktree links, git 2.48 or newer
flai config path
```

| Key | Default | Meaning |
|-----|---------|---------|
| `template.repo` | `https://github.com/bytepunx/system-flow-template` | Where `flai new`, `flai import`, and `flai upgrade` take the template from when `--template` is not given: a git URL or a local directory |
| `template.ref` | `main` | The template's branch, tag, or commit, when `--ref` is not given |
| `dashboard.image` | `ghcr.io/bytepunx/flaiover` | The image `flai dashboard` runs |
| `dashboard.tag` | `latest` | The image's tag |
| `dashboard.port` | `4242` | The host port the dashboard is published on |
| `dashboard.bind` | `0.0.0.0` | The host address it is published on; `127.0.0.1` keeps it to this machine |
| `dashboard.push_key`, `dashboard.push_known_hosts` | empty | Retired and ignored ([ADR-0031](../../design/adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)). They still load so an old file works, and `flai dashboard` names them when set; clear them with `flai config set dashboard.push_key ""` |
| `cache_dir` | `~/.flai/cache`, or `FLAI_CACHE_DIR` when the file is created | Where git templates are cloned |
| `author` | your login name | The default owner of new items and the default `--by` of transitions and acceptances |
| `worktrees.relative_paths` | `false` | Create story worktrees with relative links (git 2.48 or newer); see [Relative worktree links](#relative-worktree-links-opt-in) |

Point `template.repo` at any fork and `template.ref` at any branch, tag, or commit to use your own template. A local directory path also works, which is how the system-flow repo develops against its own `./template`. Git templates are cloned under `cache_dir`; set `FLAI_CACHE_DIR` before the first run to choose where that default lands (for example inside a repository or a CI workspace). For the dashboard, flags come first, then the `dashboard` section of the project's `system-flow.yaml`, then these keys.

The same file holds what `flai serve` may do on this host. These keys are not reachable with `flai config set`, and nothing a dashboard can ask for writes them unless you enable the `settings` host action; each has its own command, run in a shell on the host:

| Key | Set with | Meaning |
|-----|----------|---------|
| `host_actions` | `flai serve enable <action>`, `flai serve disable <action>` | For each host action (`push`, `auto-publish`, `agent`, `dashboard`, `checks`, `settings`, `host`), the projects it is on for: main checkout paths, or `*` for every project. Absent means none. `flai serve actions` says what each lets a dashboard do |
| `agent.command` | `flai serve agent set -- <program> [args...]`, `flai serve agent clear` | What is started for a ready story that names no harness: an argument list, never run through a shell, with `{story}`, `{root}`, `{model}`, and `{harness}` replaced |
| `agent.name` | `flai serve agent set --name` | The `FLAI_AGENT` prefix of the agents it starts, `agent` when empty; the story is appended, as in `agent-S-0104` |
| `agent.attended_minutes` | none | Retired ([ADR-0043](../../design/adrs/0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md)): a configuration that sets it still loads, and nothing reads it. `flai serve agent set --attended-minutes` is accepted and does nothing |
| `agent.harnesses.<name>.program`, `agent.harnesses.<name>.args` | `flai serve agent harness <name> --program <path> -- [args...]`, `--reset` | The program a harness is on this host, and the arguments that replace its adapter's defaults and say what the agent may do |
| `checks.commands` | `flai serve checks set --name <name> -- <program> [args...]`, `flai serve checks clear [name]` | The named commands a story in review is checked with, in order, in its worktree; empty means the manifest's `checks:` |
| `checks.timeout_minutes` | `flai serve checks timeout <minutes>` | The bound on one run of every check together; 15 when unset |
| `import_roots` | `flai serve import add <folder>`, `flai serve import remove <folder>` | The folders whose git repositories the board offers to import, and whose repositories with a `system-flow.yaml` `flai serve` serves |

`flai serve agent show`, `flai serve checks show`, `flai serve import list`, and `flai serve actions` print them. The [operator guide](../operators/index.md) says what enabling each host action gives a dashboard.

Environment variables flai reads:

| Variable | Effect |
|----------|--------|
| `FLAI_CONFIG` | The config file, when `--config` is not given |
| `FLAI_CACHE_DIR` | The `cache_dir` written when the config file is created |
| `FLAI_AGENT`, `FLAI_SESSION` | Who writes narrative entries, thread entries, and transitions, and in which session; `FLAI_AGENT` is also the MCP server's agent unless `flai mcp --agent` names one |
| `LOG_LEVEL`, `LOG_FORMAT` | See [Logging](#logging) |
| `FLAI_HOST_ADDR` | Where `flai host` listens, `127.0.0.1:4241` by default |
| `FLAI_RELEASES_API` | Another source of releases for `flai self-upgrade` and `flai host check` and `upgrade` |
| `FLAI_INSTALL_DIR` | Where `flai self-upgrade` from a checkout installs (`~/.flai/bin` by default) |
| `GITHUB_TOKEN`, `GH_TOKEN` | GitHub API and registry authentication for `flai self-upgrade` and `flai dashboard`; `gh auth token` is borrowed when neither is set |

`flai host` sets `FLAI_HOST_URL` and `FLAI_HOST_TOKEN` for the processes it runs; do not set them yourself.

## Version

```bash
flai version
flai version --json
```

Prints the version, commit, build date, and Go version.

## Create a project

```bash
flai new my-project
```

In a terminal this asks for the project name, key, description, and owner, with defaults pre-filled. In scripts use `--defaults` and `--var`:

```bash
flai new my-project --defaults --var "description=Billing platform" --var owner=core
```

| Flag | Effect |
|------|--------|
| `--template <url or dir>` | Use this template instead of the configured one. A local directory works. |
| `--ref <branch, tag, or commit>` | Template version to use. |
| `--var name=value` | Set a variable. Repeatable. |
| `--layout key=name` | Rename a documentation folder, for example `--layout design=architecture`. |
| `--defaults` | Never prompt. Use defaults for anything not given with `--var`. |
| `--force` | Overwrite files that already exist. Otherwise they are kept and reported. |
| `--no-git` | Do not run `git init`. |

The result has a `system-flow.yaml` recording the template and version, a `CLAUDE.md` for agents, and the `design`, `docs`, and `wip` folders ready to use.

## Convert an existing repository

```bash
cd existing-repo
flai import --dry-run     # the proposal, nothing changes
flai import               # interactive: folder names, moves, where each markdown file goes
flai import --yes         # accept every default, leave loose markdown in place
flai import --yes --commit  # then run its tests, and commit the import if they pass
```

`import` scans the tree and proposes: the three documentation folders (reusing `docs/`, `design/`, or `wip/` if they exist, or names you choose with `--layout`), whole-folder moves for `adr/`, `adrs/`, `architecture/`, `doc/`, and `documentation/`, a list of loose markdown files to place, and the code sub-projects it found by their build files (`go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`). Applying it creates the structure, renders every template file that does not already exist, performs the moves with `git mv` when the file is tracked, writes `system-flow.yaml` with the sub-projects, and runs `flai check`. Existing files are never overwritten; a conflicting move is reported and the source left in place. A repository that already has `system-flow.yaml` is refused unless `--force`.

`--commit` makes the import end in a commit. The repository must be a git repository with no uncommitted changes, so that the commit holds the import and nothing of yours. After importing, it runs the repository's tests: the checks `flai serve checks set` names on this host, if any, else what the repository has (its own Makefile's `test` target, `go test ./...`, the package manager's test script, `cargo test`, or pytest). When they pass, or none are found, it commits exactly what the import wrote and moved. When one fails it commits nothing, shows what failed, and exits with code 5; the imported files stay in the working tree for you to fix and commit. The dashboard does the same when you import a repository from the board (see the dashboard guide).

An imported project shows in the dashboard's project switcher without another step when `flai host` runs: once `system-flow.yaml` is written, with or without `--commit`, `import` registers the project with `flai serve`, as `flai dashboard` would, and says the address the dashboard shows it at. When no host runs it registers nothing and says so: `flai dashboard` in the project serves it, starting the dashboard and the host if they are not running. `--json` has the same as `serve`: `served`, `dashboard`, and, when it is not served, `reason` and `next`.

The template's `repo_url` is offered as the repository's `origin` remote made a web address (`git@github.com:owner/repo.git` becomes `https://github.com/owner/repo`). `import` and `new` refuse a `repo_url` that is not an http or https address with a host and a path, at the prompt or given with `--var`, before anything is changed, and say why: `https://github.com:owner/repo` is refused.

## Manage the template source

```bash
flai template show                       # where the template comes from and what it asks for
flai template update                     # re-fetch a git template into the cache
flai template use git@github.com:me/system-flow-template.git --ref my-branch
flai template use ./template             # a local directory, no fetching
```

Git templates are cloned into `~/.flai/cache/templates`. Private repositories work with whatever git credentials you already have.

## Work items

All of these run inside a conforming repository (anywhere below `system-flow.yaml`).

```bash
flai epic new "Billing v2" --nature feature
flai story new "Invoice PDF export" --epic E-0001
flai task new "Render invoice template" --story S-0001 --tag pdf
flai show S-0001
```

Items are created from the template's item bodies with the next free ID and linked into their parent's Stories or Tasks list. Natures: `feature`, `improvement`, `remediation`, `research`, `experiment`. `--epic` is optional for a story: not every story fits an active epic, and an epic made only to hold one is not wanted (S-0092); `--story` is still required for a task.

To create an item with its body already written, give the body on standard input. This is what the dashboard's "new" form does, and it is one step that happens or does not:

```bash
flai story new --print-body > story.md          # the sections your project's template gives a story
$EDITOR story.md                                # write the goal, criteria as - [ ] lines, notes
flai story new "Invoice PDF export" --epic E-0001 --body-stdin --autocommit < story.md
```

The heading (`# S-0007 Invoice PDF export`) and the front matter are flai's; the body is everything below the heading. `flai check` runs with the new item in place: if it reports anything the item introduces, nothing is created, the parent is left as it was, the findings are printed, and the exit code is 4. `--autocommit` commits the new item and its parent on their own (`chore: [S-0007] create story: ...`) unless the project sets `dashboard.autocommit: false`; `--trailer` adds trailer lines. Nothing is pushed.

### Who works a story: its agent

```bash
flai agent                                                   # the project's default, if any
flai agent set --harness claude-code --model claude-opus-5-5 --config effort=high
flai agent set --model claude-sonnet-5 --unset effort        # change only what you give
flai story new "Invoice PDF export" --model claude-haiku-4-5 --agent-config max_turns=20
flai edit S-0007 --model claude-opus-5-5 --agent-config max_turns=   # key= removes a key
flai edit S-0007 --clear-agent                               # the story has none
flai agent clear
```

A story's agent says what will work it: the harness that runs the agent (such as `claude-code`), the model, and options for that harness as `key=value`, which flai passes on without reading them ([ADR-0037](../../design/adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)). The project's default lives in `system-flow.yaml` under `agent`. Every story created while it is set gets a copy in its front matter, with whatever `--harness`, `--model`, or `--agent-config` you give on `flai story new` laid over it. The copy belongs to the story. Changing the default later affects only stories made after it; to move an open story to another model, edit it. `flai show` and `flai edit --show` print a story's agent. Epics and tasks have none.

A project with no default and no story with an agent has no `agent` keys at all, and an older flai reads it as before. Once a story has one, use a flai that knows about agents (`flai self-upgrade`).

### Moving work

```bash
flai move S-0001 ready          # needs acceptance criteria; tasks are not required
flai move S-0001 in-progress    # warns if the WIP limit is exceeded
flai move T-0001 in-progress
flai move T-0001 done           # tasks may skip review
flai move S-0001 review         # needs at least one task and nothing uncommitted in the story's worktree
flai move S-0001 done --by alex # needs every task closed and every criterion checked
flai move S-0001 in-progress --reason "tests missing"     # from review
flai move S-0002 cancelled --reason "superseded by S-0005"
flai move E-0003 cancelled --reason "a different route" --dry-run   # what would go with it
flai move S-0003 backlog        # back from ready, or reopened from cancelled
flai move S-0004 ready          # back from in-progress; it goes last in the ready order
```

An item moves back one column from ready, in-progress, review, or cancelled, and never out of done ([ADR-0055](../../design/adrs/0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)). A reopened item is refused while its parent is cancelled; move the parent back first. What its cancellation cancelled under it stays cancelled. Until it closes again it has no completed time, lead time, or cycle time in `flai stats`.

### Tokens and cost

An item may carry `usage`: the tokens each model read and wrote on it, what they cost in US dollars, and how long agents worked ([ADR-0051](../../design/adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)). You never type it. `flai serve` measures a story, and each of its tasks, from the logs of the agents it started for the story (see [flai serve agent usage](#flai-serve-flai-on-the-host-for-the-dashboards)). Whenever an item enters done, whether by `flai move`, the MCP `item_move`, or `flai accept`, each item above it that was not measured gets the sum of its children's, up to its epic: an epic always sums its stories. `flai show` prints it:

```text
  usage: 20.1M tokens · $8.13 · 18m3s of agent work · measured from its agents' logs
    claude-opus-5-5  input 256 · output 89.3K · cache read 19.7M · cache write 327.6K · $8.1258
```

`(estimated)` after the cost means some of it was not reported by the harness: a task's share of its story's session, or a run that ended without its totals, priced at the rate the logs report for the model. `flai show --json` returns it under `item.usage`, and `flai stats` charts it ([Flow metrics](#flow-metrics)). A flai older than the one that brought `usage` refuses an item that carries it: upgrade the flai on your host first.

### Cancelling

Cancelling an epic cancels every story under it that is still open and their open tasks; cancelling a story cancels its open tasks. Items that are done or already cancelled are left alone. `flai move` lists what will be cancelled first and, on a terminal, asks before doing it (`--yes` skips the question, `--dry-run` only lists). Each cancelled item records its own transition and a note that names the cause, such as `E-0003 cancelled: a different route`, so an archived task still says why it ended.

```text
Cancelling E-0003 also cancels 3 items:
  story S-0009  in-progress Berths
    task  T-0031  backlog     Dredge
  story S-0010  review      Cranes  (in review: its work stays on its branch, unmerged)
```

A story in review cannot be cancelled directly: accept it or send it back. It is cancelled only with its epic, and the list says so; accept it first if you want the work. Cancelled is final. Nothing of git is touched: the command names each cancelled story's narrative, branch, and worktree and leaves them for you to keep or remove (`flai archive` moves the narratives). If a tree was cancelled by an older flai or edited by hand, `flai check` reports each open item under a cancelled parent as `item.parent-cancelled`.

Every move appends to the item's `transitions` with a timestamp and who made it (`--by`, default the config author). Reasons land under the item's Notes.

A story is ready once its goal and acceptance criteria are written. The agent that pulls it moves it to `in-progress` and then writes its tasks; a story with no tasks is refused at `review`, and `flai check` reports `story.tasks` for a `review` or `done` story that has none.

### Blocking

```bash
flai block T-0001 --reason "waiting on API keys"
flai unblock T-0001
```

Blocked items keep their column; the interval is recorded so blocked time shows in the charts.

### The board

```bash
flai board          # stories by column with age, nature, blocked flag, WIP counts
flai board --all    # epics and tasks too
flai board --json
flai board limit in-progress 3   # the column's WIP limit; 0 for none
```

`flai board limit` sets the WIP limit of `ready`, `in-progress`, or `review` in `wip/kanban/board.md`, the one place flai, `flai serve`, `flai check`, and the dashboard read it from; the other columns have none. A move past a limit warns and is made. `flai serve` starts an agent for a ready story only while `in-progress` has room, so raising that limit lets it start the next one. The board's lane menu in the dashboard does the same.

The `backlog` and `ready` columns are listed in pull order: the stories named in `order` in `wip/kanban/board.md`, in that order, then the rest by ID. Agents pull the first ready story, so this is how you say what comes next.

```bash
flai order S-0061 --top                # first in its column
flai order S-0059 --before S-0061      # just above another story of the same column
flai order S-0047 --after S-0053
flai order S-0056 --bottom
```

Only ready and backlog stories can be placed, and only relative to a story in the same column; `flai move` changes the column. A story you move to `ready` joins the end of the ready stories. Backlog stories you have never placed stay out of the list and come last, by ID. The dashboard's board does the same thing when you drag a card up or down within a column.

### Story branches

```bash
flai stream open S-0037        # narrative, plus branch story/S-0037 in .flai-cache/worktrees/S-0037
flai stream sync S-0037        # rebase the branch onto main, then check it against the other open branches; run at every task transition
flai stream open S-0037 --no-branch
```

Each story is worked on its own branch, checked out in a worktree under `.flai-cache/worktrees/`. Code, design, and docs changes land there; `wip/` is always written in the main checkout, so the board and the dashboard stay current whatever branches exist. `flai stream sync` rebases the branch onto the main branch, stashing uncommitted work around it; conflicts stop inside the worktree and are listed, resolve them, `git rebase --continue`, and sync again. `flai accept` rebases, fast-forwards the branch into main, removes the worktree and branch, then tags and pushes.

After a clean rebase, sync checks the branch against the other stories in progress or in review, so that two stories that change the same lines find out while both are still open, not when the second is accepted:

```text
story/S-0131 is rebased onto main
story/S-0131 conflicts with story/S-0130 (in progress) in flai/cmd/edit.go; see TH-0024
story/S-0131 merges cleanly with story/S-0129 (in review)
story/S-0131 changed 1 path outside S-0131's touches: flai/internal/threads/threads.go
widen them so that stories that overlap wait: flai touches S-0131 flai/cmd docs flai/internal/threads/threads.go
```

- **Conflicts.** Sync merges the two branches in git's object store only (`git merge-tree --write-tree`, git 2.38 or newer; an older git skips it with a warning), so nothing changes in either worktree. For each pair that conflicts, flai opens one thread on the story that synced, titled `S-0130 and S-0131 conflict when merged`, listing the paths. It shows in both stories' agents' MCP `inbox` and in the designer's inbox on the dashboard. Settle it between the two stories: one narrows its change, or names the other in `after:` and waits. A later sync with the same paths adds nothing, new paths add an entry, and flai resolves the thread once the two merge cleanly or the other story is no longer open.
- **Outside the touches.** Sync lists the files the branch changed since main that the story's touches, and its open tasks', do not cover (see [Touches](#touches)), and prints the `flai touches` command that widens them. Touches that are too narrow let a story that overlaps start beside it.
- Neither check fails the sync. `--json` adds `branches` (each with `story`, `status`, `branch`, `clean`, `conflicts`, and `thread`), `outside_touches`, and `trial_merge_skipped` when git is too old.

#### Relative worktree links (opt-in)

`worktrees.relative_paths` is off by default, and flai never turns it on for you, whatever git you have. Set it when you move the clone around or share it between machines; the dashboard, which once needed it for clones it could not mount at their own path, no longer mounts anything. With it on and git 2.48 or newer, `flai stream open` creates the worktree with `git worktree add --relative-paths`, so the links work wherever the repository lies. With it on and an older git, flai warns, naming your git version and the setting, and creates an ordinary worktree.

Turning it on changes the clone, not just the worktree: the first relative worktree sets `extensions.relativeWorktrees` in `.git/config`, and any git older than 2.48 then refuses the whole repository with `unknown repository extension found: relativeworktrees`. That includes other tools on your machine that bundle their own git. `flai check` warns with `git.relative-worktrees` when it finds the extension and an older git on your `PATH`.

To go back, with git 2.48 or newer:

```bash
git worktree repair --no-relative-paths .flai-cache/worktrees/S-0001   # once per worktree; the path is required
git config --unset extensions.relativeWorktrees
flai config set worktrees.relative_paths false
```

With only an older git: delete the `relativeWorktrees = true` line from `.git/config`, then run `git worktree repair .flai-cache/worktrees/S-0001` for each worktree before anything prunes them, and set the key to false.

### Touches

```bash
flai story new --epic E-0006 "Title" --touches flaiover/src/lib,docs/users
flai touches T-0121 flai/internal/workitem
flai touches T-0121 --clear
```

`touches` is the list of paths or components a story or task is changing. `flai check` warns (`wip.overlap`) when two in-progress items cover the same path, the board prints it under each card, and the dashboard shows a "being worked on" notice on those documents.

It is also a claim that decides what starts ([ADR-0046](../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)). A story's claim is its `touches` and those of its tasks that are not done or cancelled; a sub-project's name or tag (`cli`, `flai`) means its path. A ready story is *held* while its claim overlaps the claim of a story in progress or in review: the same path, or one inside the other (`flai/cmd` and `flai/cmd/serve`, not `flai` and `flaiover`). A story with no touches may change anything, so it is held while any story is open, and while it is open itself it holds every ready story. Declare touches when you create a story; the agent that pulls it may widen them.

A held story is not started by `flai serve` and not offered by `wait_for_work`. The next ready story that is not held goes ahead of it, and it keeps its place and goes first once it is clear. The board says why:

```text
ready
  S-0130 Serve the board faster                         feature          2m HELD
         touches flai/cmd/serve
         held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0128 (in progress) touches; starts when S-0128 is accepted, cancelled, or sent back
```

You can still start it yourself: `flai move S-0130 in-progress` and `flai serve agent start S-0130` warn and go ahead.

### Waiting for another story

```bash
flai edit S-0131 --after S-0129,S-0130     # S-0131 starts once both are done
flai edit S-0131 --clear-after
```

Some stories depend on another for a reason that is not about files: they build on its API, or on what its research found. Name those stories in `after` and the story is held in ready until every one of them is done, however little they touch in common:

```text
         held (after): waits for S-0129 (in progress) and S-0130 (ready); starts when S-0129 and S-0130 are done
```

A story named that is cancelled keeps the hold, and the reason says so: drop it from `after` if the story no longer needs it. `flai check` refuses an `after` that names a story that does not exist, the story itself, or a cycle (`S-0131` waits for `S-0132`, which waits for `S-0131`). The dashboard's story editor and the MCP `item_edit` tool set it too.

A flai older than the one that brought `after` refuses to read a story that carries it, so upgrade the flai on your host (`flai self-upgrade`) before you use it.

### A story's topics

```bash
flai story new "Structured request logs" --epic E-0004 --topics logging
flai epic new "Release automation" --topics release
flai edit S-0131 --topics logging,release    # replaces them; --clear-topics removes them
flai show S-0131                              # the story's topics and where each came from
```

`topics` say what a story or epic is about beyond the components it reaches, in words such as `logging` or `release` ([ADR-0047](../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md)). They choose what an agent working the story is primed with, as the topics on documents do (below); `tags` still say which component a story delivers to and so which release it cuts. Each topic is one word of letters, digits, dot, dash, or underscore. A task carries none.

You rarely need to write many. flai works out a story's topics from its own and its epic's, from the name, tags, and kind of every sub-project in `system-flow.yaml` that one of its tags names or its touches (and its open tasks' touches) reach, `code` when one of those is not the template, and `all`. `flai show` lists them:

```text
  topics:
    logging    own
    flai       flai by tag cli and touches flai/internal/logx (S-0131)
    cli        flai by tag cli and touches flai/internal/logx (S-0131)
    go         flai by tag cli and touches flai/internal/logx (S-0131)
    code       flai is code
    all        every story
```

`flai show --json` returns the same under `topics`, with every source: its `kind` (`own`, `epic`, `tag`, `claim`, `code`, `all`), the `item` that says it, the sub-project (`project`) it reached, and the tag or touch (`via`) that reached it. A touch outside every sub-project, such as `design/system`, adds nothing. The dashboard's story and epic editors, and the MCP `item_new` and `item_edit` tools, set topics too. As with `after`, a flai older than the one that brought `topics` refuses to read an item that carries them: upgrade the flai on your host first.

### Threads

```bash
flai thread new --on design/system/flaiover-dashboard.md --heading "Workbench (E-0006)" "Title" "Question"
flai thread new --on S-0038 "Title" "Question"        # anchored to an item
flai thread list                                       # unresolved threads
flai thread list --on S-0038 --all
flai thread show TH-0001
flai thread reply TH-0001 "Answer"
flai thread resolve TH-0001 --reason "settled in ADR-0021"
```

A thread is one file under `wip/threads/`, anchored to a document, a heading in it, or a work item, with dated entries by author. The author is `--by`, else `FLAI_AGENT`, else the config author. A reply from anyone but the opener marks the thread `answered`; the opener's follow-up makes it `open` again; `resolve` closes it. Unresolved threads on a story or its tasks are mirrored into the story narrative under `## Open questions`, so an agent sees them without the dashboard. `flai check` validates threads: the anchor must exist, a named heading must still be in the document, and open threads on archived items are flagged.

### Changing an item after it was made

```bash
flai edit S-0085 --show                         # the fields, the body below the heading, and the hash
flai edit S-0085 --title "A better name" --autocommit
flai edit S-0085 --nature improvement --tag dashboard --tag cli --touches flaiover/src
flai edit S-0085 --parent E-0004                # an open epic for a story, an open story for a task
flai edit S-0085 --body-stdin --hash <hash> < body.md
flai edit S-0085 --clear-tags --clear-touches
flai edit S-0085 --after S-0084                 # hold it until S-0084 is done; --clear-after lets it go
flai edit S-0085 --topics logging               # what it is about; --clear-topics removes them
flai edit S-0085 --harness claude-code --model claude-sonnet-5 --agent-config effort=high
```

`flai edit` changes what an item says about itself: title, nature, tags, touches, parent, a story's or epic's `topics`, a story's `after` and agent, and the body below its heading, any of them together. What is the item's state stays with its own commands: the status with `flai move`, blocking with `flai block`. A closed or archived item is refused.

A title lives in several places, and a retitle keeps them in step: the front matter, the heading, the file's name, the line in the parent's list, the story's narrative, and links to the old file name under design, docs, and wip (from a story's worktree only under wip, because design and docs there are another branch's). With `--hash`, the one `--show` printed, a change someone made meanwhile is a conflict (exit 3) and nothing is written. `flai check` runs with the change in place: what the change introduces refuses it, every file is put back, and the findings are printed (exit 4). What is simply not allowed, a nature there is not, an epic as a task's parent, is said as a `rule:`. `--autocommit` commits every file the edit touched in one commit; nothing is pushed.

Agents connected over MCP are told of an edit someone else made, as a change of kind `edited` that names what changed. That comes from a small log under `.flai-cache`, outside git, like an agent's read marker: hand edits of a file are not reported, as before.

They are also told when an accepted story changed paths their own story claims, as a change of kind `overlapped` on their story. `cause` is the accepted story and `to` lists the paths. The agent syncs its story and runs its tests again before it goes on. See acceptance, below.

### Serving agents over MCP

```bash
flai mcp          # an MCP server on stdio; agents start it, you do not
flai mcp start    # the same server over HTTP, for an agent that cannot start a process here
```

`flai mcp` gives an agent session a typed, low-latency view of the repository. Register it once per project in `.mcp.json` (new projects get this from the template):

```json
{ "mcpServers": { "flai": { "command": "flai", "args": ["mcp"] } } }
```

The server works as the agent `FLAI_AGENT` names, or as `--agent` names when given: `flai mcp --agent agent-S-0104`. `flai serve` starts the agents it starts that way, so a `FLAI_AGENT` that an agent's own settings put into the server's environment (Claude Code's `env` in `.claude/settings.json` does that) cannot make every agent one name.

Started in a folder that is not itself a project, such as `~/git`, `flai mcp` serves every system-flow project in that folder and up to three levels below it ([ADR-0036](../../design/adrs/0036-a-folder-that-is-not-a-project-is-served-whole-by-flai-mcp-and-flai-dashboard.md)). One agent started there works across all of them. `inbox`, `wait_for_work`, and `wait_for_events` cover every project and say which one each thing is in, and every other tool takes `project`: a key `inbox` lists, or the project's folder. A project created or imported below the folder joins within a few seconds. Started in a folder with no project below it at all, the server still starts and says there is none yet. Over HTTP (`flai mcp start` and the rest) it still serves one project, so run those in the project.

| Tool | What it does |
|------|--------------|
| `inbox` | (Since S-0085 `changes` also reports `edited`: someone changed an item's title, fields, or body with `flai edit` or from the dashboard, and `to` names what. Since S-0132 it reports `overlapped`: a story was accepted, `cause`, that changed paths this story claims, `to`.) Threads awaiting the agent (`awaiting: you` when the last entry is not the agent's; `story` filters, `all` includes the rest), `ready`: the stories ready to pull, in pull order, with `can_pull` from the in-progress limit and `held` with why on a story an open story's claim holds, and `changes`: what others did to work items since this agent last looked (moved, blocked, unblocked, pull order changed), each reported once `unpushed`, on every call while it is true: an acceptance made in this clone and not pushed (items, commits ahead, tags), which the agent pushes from the host with `git fetch` and `flai push --pending` |
| `board` | The board as `flai board --json` prints it, a held ready story with `held` and why; `all` adds epics and tasks |
| `thread_get`, `thread_open`, `thread_reply`, `thread_resolve` | Read, start, answer, and close threads as the agent (`FLAI_AGENT`) |
| `item_get`, `item_move` | Read an item with its children, a story's agent and the project's default, and the hash of its file; transition it with the workflow rules. Moving a story or epic to done is refused: acceptance is yours |
| `item_new`, `item_edit` | Create an epic, a story (with an `agent` over the project's default), or a task; change an item's own words, as `flai edit` does: `agent` replaces a story's agent whole and `clear_agent` removes it, `after` replaces the stories a story waits for and an empty list removes them, and the `hash` from `item_get` refuses a change made meanwhile. Neither commits: the agent commits with its work |
| `doc_get` | A markdown document under the design, docs, or wip folders; nothing else in the repository is served. With `heading`, only that section and the sections below it, with its heading path and line ([Read design on demand](#read-design-on-demand)) |
| `doc_search` | The sections of the design and docs folders, the conventions among them, that rank highest against `query`: at most 20 (`limit` for fewer), each with its path, the document's title, its heading path, line, first lines, and size ([Read design on demand](#read-design-on-demand)) |
| `prime` | A story's context pack, as `flai prime --story <id> --json` prints it, fitted to `budget` (default the project's `prime.budget`, else 80 KB): its topics, the conventions with the sections those topics leave out taken out, what the story names whole (a large document it names only by a path written out as a brief), briefs of the design and tech files and the ADRs its topics and one link step select, ranked sections to fill the budget, each with its reason and size, and a catalog of the rest, to read with `doc_get` and a `heading` when needed; with `role` (`explore` or `verify`), the smaller pack for a sub-agent ([Sub-agents](#sub-agents)) ([Prime a session](#prime-a-session)) |
| `who_touches` | In-progress and in-review items whose `touches` cover a path |
| `wait_for_work` | What to do when you have nothing to work on. Answers at once with `resume` and your own story if one is still in progress, `thread` and the threads awaiting you that were written to since it last answered, or `pull` and the first ready story that is not held when the in-progress limit leaves room. Otherwise it waits until one of those is true, up to the timeout, and then says whether it was waiting for room, for a held story to be clear (`held`: each ready story says why), or for a story to be ready: call it again. Hold it whenever you are idle, and you pull the next story as soon as there is one |
| `wait_for_events` | Returns at once when something changed since this agent last looked, otherwise blocks until a thread, item, or narrative changes, or the timeout passes. Returns `events` in the same shape as `changes`, and the changed paths |

An agent that cannot start a process on the host reaches the same server over HTTP. flai serves it itself, one server per project, on the host ([ADR-0030](../../design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)); the dashboard is not involved and need not run. When `flai serve` serves the project (which `flai dashboard` arranges), `flai host` keeps that server running for you ([ADR-0040](../../design/adrs/0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)). `flai host status` or `flai mcp status` says where it listens. Otherwise start it yourself:

```bash
flai mcp start      # in the background, at http://127.0.0.1:4243/mcp unless --addr says otherwise
flai mcp status     # where it listens, and an agent's configuration to copy
flai mcp token      # its bearer token (.flai-cache/mcp.token); --rotate replaces it
flai mcp stop
flai mcp http       # the same server in the foreground; Ctrl-C stops it
```

The agent's configuration names that address and sends the token as a bearer; `X-Flai-Agent` is the name it works under, as `FLAI_AGENT` is locally (without it, the client's own name is used):

```json
{
  "mcpServers": {
    "flai": {
      "type": "http",
      "url": "http://127.0.0.1:4243/mcp",
      "headers": { "Authorization": "Bearer ${FLAI_MCP_TOKEN}", "X-Flai-Agent": "claude@laptop" }
    }
  }
}
```

The tools and their behaviour are the same, because it is the same server. Keep the token out of the file: Claude Code expands `${VAR}` in `.mcp.json`, as above; for another client, check how it takes a secret. This token is the MCP server's own, not the dashboard's. The server listens on the host only; from another machine, reach it through an SSH forward or a tunnel that terminates TLS, as the operator guide describes. It serves the main checkout, and a second project on the same host needs another port (`--addr`, remembered per project). Until flaiover 0.22 the dashboard served MCP at its `/mcp`; that address now answers 410 and says to use this instead.

"Since this agent last looked" is a marker per agent name (`--agent`, else `FLAI_AGENT`) under `.flai-cache/mcp/`, outside git. It only decides which changes are news; ready work and open threads are listed on every call, so nothing depends on it. One look reports at most 50 changes, the newest, and says in `changes_omitted` (`events_omitted` for `wait_for_events`) how many older ones it left out; those are not reported later. The first look under a new name covers the last 24 hours and tells of stories and epics only, not task transitions: to an agent that has just arrived, a day of task moves is history, and `board` and `item_get` show how things stand. An agent that ends its turn between your messages calls `inbox` when it starts again and hears what you did in between; one that stays running holds `wait_for_events` and hears within a second. `flai move` records `FLAI_AGENT` as who moved an item when it is set, so an agent is not told about its own moves; a move on the dashboard's board is recorded as the project's `owner`.

Design and docs files are also exposed as resources (`flai://design/...`, `flai://docs/...`). Everything the server writes goes through the same code as the CLI, so files stay the record and the dashboard shows agent replies as they land.

### Agent narratives

```bash
export FLAI_AGENT=claude FLAI_SESSION=abc123
flai stream open S-0001
flai stream log S-0001 "T-0001 done, starting T-0002"
```

`wip/agents/index.md` is regenerated after every open, log, move, and archive.

### Archiving

```bash
flai archive --dry-run
flai archive               # everything done or cancelled that is safe to move
flai archive S-0001         # one story with its tasks and narrative
```

### Item IDs

New items and issues get four-digit IDs (`S-0034`, `I-0012`), the same width as ADRs. Commands accept an ID with any zero padding or none, so `flai show` and `flai issue bump` find the item whether you type the number with two, three, or four digits. A repository created with three-digit IDs keeps working as it is; to widen it in one step:

```bash
flai migrate ids --dry-run   # list every rename and rewrite
flai migrate ids             # git mv the files, rewrite references, refresh the index
```

The migration touches `design/`, `docs/`, `wip/`, the root markdown and yaml files, and each project's root markdown files. Fixtures under `testdata/`, `node_modules`, and hidden folders are left alone. Commit the result as a `chore:` on its own.

## Record a decision

```bash
flai adr new "The dashboard authenticates every request with a project token" --refines 16
flai adr new --print-body > decision.md         # the sections of your design/adrs/0000-template.md
flai adr new "Replace the SPA with server rendering" --status accepted --supersedes 7 --body-stdin --autocommit < decision.md
flai adr accept 27                               # a proposed ADR becomes accepted, dated today
flai adr topics 27 cli template                  # the stories it is for, on an ADR of any status
```

`flai adr new` takes the next number from the files in `design/adrs` (one more than the highest; gaps are not filled), names the file `NNNN-slug.md`, writes `id`, `title`, `status` (`proposed` unless you say `--status accepted`), `date`, `supersedes`, `superseded_by`, and `refines`, adds the row to `design/adrs/README.md`, and sets `superseded_by` on each ADR it supersedes, one of the two edits allowed to an accepted ADR (the other is `flai adr topics`). The body is your template's sections, or standard input with `--body-stdin`. `flai check` runs with everything in place: if it reports anything the ADR introduces, every file is put back, the findings are printed, and the exit code is 4. `--autocommit` makes one `docs: ADR-NNNN <title>` commit of what was written, unless the project sets `dashboard.autocommit: false`; nothing is pushed.

`flai adr topics` sets which stories an ADR is for: it writes `topics: [cli, template]` into the front matter, in place of the key when it is there and at the end otherwise, and changes nothing else. It works on an ADR of any status, since `topics`, with `superseded_by`, is the one key an accepted ADR may gain. Topics are separate arguments or comma separated. Each is a word, and `flai check` must know it: `all` (every story), `code`, a sub-project's name, tag, or kind in `system-flow.yaml`, or a topic a story or epic declares. A topic it does not know is refused, the file is put back, and the exit code is 4. `--autocommit` and `--trailer` work as for `flai adr new`.

`flai check` warns with `adr.index` when an ADR file has no row in the index or a row has no file. An accepted ADR is immutable but for its `topics`: `flai doc save` refuses any other change, and the dashboard's editor shows it read-only.

### Topics on documents

Conventions, `design/system` and `design/tech` files, and ADRs say which stories they are for with `topics: [...]` in their front matter, and any heading can narrow that with a comment at the end of its line ([ADR-0047](../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md)):

```markdown
## Go <!-- topics: cli, go -->
```

A heading's topics cover everything down to the next heading at its level or higher; a heading without them takes its parent's, and the top headings take the file's. `all` is every story, and a convention without `topics` is read as `[all]`. A convention's `roles` (`explore`, `verify`) say which sub-agents read it as well ([Sub-agents](#sub-agents)); `flai check` warns with `conventions.roles` on any other role. `flai check` warns with `doc.topics` on a `design/system` or `design/tech` file without topics, and with `doc.topic` on a topic that is not `all`, `code`, a sub-project's name, tag, or kind, or one that a story or epic declares ([A story's topics](#a-storys-topics)). `flai prime --story` selects conventions, design, tech files, and ADRs by them ([Prime a session](#prime-a-session)). A design, tech, or ADR file without topics comes into a pack only when something names or links it, or it ranks.

## Read design on demand

A context pack briefs most of what it selects ([ADR-0049](../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)). Find the sections that bear on a question, then fetch one by its heading rather than the whole file:

```bash
flai doc search context pack budget                 # the 20 best sections of design and docs
flai doc search --limit 5 --json worktree           # fewer, as data
flai doc show design/system/flai-cli.md --heading "Configuration"
flai doc show design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md --heading "Decision › On demand"
```

`flai doc search` ranks every section of the design and docs folders, the conventions among them, by BM25 against the words given, with the index the pack's ranked step uses: each section runs from its heading to the next heading of any level. Folder READMEs and the ADR template are left out, and so are words such as "the" and single letters. It prints at most 20, best first: the path and heading path, the size in bytes, the line, and the first lines. `--json` gives `query` and `hits`, each with `path`, `title`, `heading`, `line`, `lines`, `size`, and `score`. The MCP tool `doc_search` returns the same.

`flai doc show <path> --heading "<heading>"` prints that section and the sections below it, down to the next heading at its level or higher. A heading is its text (`Configuration`), a heading path joined by `›` or `>` that matches at the end (`Commands › flai prime`), or its anchor slug (`#flai-prime`), compared without case or punctuation. A heading that names more than one section is refused with each one's heading path and line, and one that names none is refused with every heading the document has. `--json` gives `path`, `heading`, `line`, `size`, and `text`, and no hash: a section is for reading, and `flai doc save` saves whole documents. The MCP `doc_get` takes the same `heading` and returns the section as `body`, with `heading` and `line`. Without a heading both return the whole document, as before.

## Edit a document through flai

```bash
flai doc show design/system/overview.md --json      # content, its sha256, and the edit mode
flai doc save design/system/overview.md --hash "$h" --message "clarify the overview" < new.md
flai doc save wip/kanban/stories/S-0001-x.md --hash "$h" --no-commit < new.md
```

This is the save path of the dashboard's editor, usable from a script too. `show` reports the mode: `full` for design and docs files, `body` for work items, narratives, and `board.md`, whose front matter flai owns, and `none`, with the reason, for generated files, threads, issues, and the archive. `save` reads the new content on standard input and needs the hash `show` gave for the content you started from. If the file has changed since, it stops as a conflict (exit code 3) and prints a diff of the current file against yours; with `--json` the current content and its hash come too, and saving over it means saving again with that hash. If the content changes front matter flai owns, or the check finds anything your change introduced, the save is refused (exit code 4), the findings are printed, and the file is left as it was. A saved file is committed on its own, `docs:` for design and docs and `chore:` for wip with the item's ID in brackets, authored by your git identity; `--message` sets the subject, `--trailer` adds lines, `--no-commit` or `dashboard.autocommit: false` in `system-flow.yaml` skips the commit. Nothing is pushed. For design and docs files the `updated` date is set to today unless you changed it.

## Check the repository

```bash
flai check            # errors exit 1
flai check --strict   # warnings exit 1 too, use this in CI
flai check ../other-repo --json
```

Every finding is one line, `path:line: level: rule: message`, so editors and CI annotate it. Rules cover the manifest and layout, every work item in `kanban/` and `archive/` (front matter, IDs and file names, parents and children, state history, acceptance criteria, required sections), narratives and their index, the board's WIP limits and pull order, and front matter on `design/` and `docs/` files including ADR numbering. `README.md` files are exempt from front matter.

## Flow metrics

```bash
flai stats                          # stories completed in the last 30 days
flai stats --since 90d --by nature  # grouped
flai stats --type task
flai stats --json                   # per-item values, weekly throughput, burn-up and cumulative flow series, aging, usage
flai stats --since 7d --bucket hour --json   # spend over time by the hour (day by default, or week)
```

The table shows completed and cancelled counts, throughput per week, current WIP, cycle, lead, and queue time distributions (p50, p85, max, mean), flow efficiency, time-in-state share, aging work against the cycle time p85, and throughput for each week of the window, weeks with nothing done included. When items done in the window carry `usage` ([Tokens and cost](#tokens-and-cost)), a `usage` line adds what agents spent on them: tokens, cost, and agent time; then what that comes to per item, with the agent time an item took, per minute of agent work, and per dollar; then per model its tokens, cost, tokens per minute of agent work, and how many items it worked on. `--json` has every item of the type under `items`, done in the window or not, and burn-up and cumulative flow for each day of the window, from the day that holds its start, or the first item's creation if later (S-0166). It has each item's usage under `items[].usage`, and under `usage` the totals, the models, and the items done in order with what had been done and spent by then (`done`, and per model `by_model`). Since S-0163 `usage.spend` lays out what was spent on epics, on stories, and on tasks over time, whatever `--type` is: one point per `--bucket` (`hour`, `day`, or `week`; a day unless you say, and an hour only over a window of 31 days or less), from the first bucket in which an item with usage was done to now, each with the items done in it, their tokens, cost, and agent seconds, the tokens, cost, and agent minutes per item (`minutes_per_item`, since S-0169), the tokens per agent minute and per dollar, the running means per bucket, and the same per model. An item counts in the bucket in which it entered done. Rates are per minute of agent work (`tokens_per_minute`); `tokens_per_hour` is still there for scripts written before. Definitions are in [design/system/metrics.md](../../design/system/metrics.md); the dashboard uses the same numbers.

## Prime a session

```bash
flai prime          # paths of design/conventions in read order, README first
flai prime --cat    # the same files' contents, each under a header
flai prime --json
flai prime --story S-0137         # S-0137's context pack, fitted to 80 KB
flai prime --story S-0137 --budget 120KB
flai prime --story S-0137 --json
flai prime --story S-0137 --role verify   # the pack for a verifier sub-agent
```

Agents read these before any change. An agent with a story primes with `flai prime --story <id>`, or the MCP tool `prime`, which returns the same pack as `--story --json` and takes `budget` too: the prompt `flai serve` gives the agents it starts says so, as do the MCP server's instructions, `CLAUDE.md`, and `session-start.md`. They also tell the agent that a brief is not the document, and to read the section that bears on the story with `doc_get` and its `heading` before relying on it or changing what it describes. Without a story an agent primes with `flai prime --cat`; a shell hook or a wrapper can pipe either into the session. A sub-agent the story's agent starts primes with `--role` ([Sub-agents](#sub-agents)).

`--story` prints what an agent working that story needs: its context pack, fitted to a size budget ([ADR-0047](../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), [ADR-0049](../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)). The budget is `--budget`, else `prime.budget` in `system-flow.yaml`, else 80 KB: bytes, or a number with `KB` or `MB`. It counts everything printed, the header included.

The header names the story and its topics, with where each came from (as `flai show` gives them). It gives the pack's size in bytes and lines against the budget, and one line per thing printed with its size. When the conventions alone exceed the budget, it says so, and the pack is the conventions and a catalog. When the conventions and what the story names exceed it, it says that, and nothing is ranked. When the briefs take it over, it says that too: every brief is kept and nothing is ranked (TH-0032). Until a project narrows its conventions' topics, a pack of a code story is usually over 80 KB this way.

The pack, in order:

1. Each convention as `--cat` prints it, with every section whose topics include neither `all` nor one of the story's left out ([Topics on documents](#topics-on-documents)). The front matter, the baseline marker, and the `## Project additions` heading always stay, and so does the heading above a section that is kept. Then the open issues. Nothing here is cut for the budget.
2. What the story, its epic, and its tasks name, whole: a link, a path written out, or an ID (`ADR-0046`). A link to `file.md#heading` brings that section and the ones below it. Each is headed with its path, its heading path when it is a section, and a `reason:` line. Nothing here is cut either; a story that names more than the budget holds is a story to split. One exception ([ADR-0050](../../design/adrs/0050-a-document-a-story-names-only-by-its-path-written-out-is-briefed-when-it-is.md)): a document named only by its path written out, such as a task's `design/system/flai-cli.md` as a file to update, is briefed instead of loaded when it is larger than an eighth of the budget (10 KB at 80 KB). Link a document, or name an ADR by its ID, when the agent should read it before it starts; write out the path when it only says where the work lands.
3. `briefs`: first the documents named by a path written out and too large to load (reason `named in T-0525`), each with a line telling the agent that the story names it and that it reads the file, or the sections it changes, before relying on it or changing it; then each `design/system` and `design/tech` file whose topics, or a heading's, select it, as one line with its path, title, whole size, and reason, then its first paragraph and heading outline. Headings the topics chose are marked `selected` when they are not all of them, and those loaded below are marked `loaded`. A brief is not the document: when one bears on the story, read the section with `flai doc show <path> --heading "<heading>"` or the MCP `doc_get` with its `heading`, or the whole file, before relying on it or changing what it describes ([Read design on demand](#read-design-on-demand)).
4. `decisions`: each ADR reached by its topics or one step from what is named or briefed, as one line with its ID, title, path, reason, and decision sentence: the first sentence under its `## Decision`.
5. Ranked sections: the sections that best match the story's title, goal, and acceptance criteria among those not loaded, best first, each cut at its own heading, while the budget has room. An ADR loads whole when it fits and as its best section when it does not.
6. The catalog, one line per document neither loaded nor briefed, then `left out`, one line per convention section left out with its topics, such as `code-quality.md § Rules › Go (go)`.

| Reason | Why it is there |
|--------|-----------------|
| `linked from S-0137` | The story, its epic, or one of its tasks names it. |
| `named in T-0525` | The story, its epic, or one of its tasks writes out its path, and it is too large to load whole, so it is briefed. |
| `topics: cli` | Its topics, or a heading's, include `all` or one of the story's. |
| `linked from design/system/x.md § Heading` | A named or briefed section links or names this ADR. Links are followed one step only. |
| `refined by ADR-0047`, `refines ADR-0047` | A named or briefed ADR refines this one, or this one refines it. |
| `supersedes ADR-0042` | It replaces a superseded ADR that something above chose. The superseded one is not printed. |
| `rank 3` | Third among the ranked sections that loaded. |

Nothing is printed twice. The first reason wins; a loaded item lists the later ones after it (`reason: linked from S-0137; also topics: cli`), and a brief counts them (`topics: cli, and 2 more`). `--json` returns the same as data: `budget`, `exceeded` (`conventions`, `named`, or `briefs`: the part that takes the pack over), and `size`; per convention, the sections kept and left out with their heading paths, lines, and topics, and its `size`; `items`, each with `path`, `id`, `title`, `heading`, `label`, `step` (`named`, `briefed`, or `ranked`), `reason`, `also`, `size`, `whole` (a brief's document size), and `text` (what is loaded, or the brief); and `catalog` with `not_loaded` and `in_part` (each entry with its `outline`). An archived story gets the pack it would get today, which is how you replay one. An ID that is unknown or not a story is refused with a message naming it. `flai check` validates the folder: every file needs `title`, `updated`, `audience: agent`, a unique `order`, and `status`; exactly one baseline marker followed by a `## Project additions` section; under 120 lines; and the README must list each file exactly once. It also warns with `adr.decision` about an ADR whose `## Decision` does not open with a sentence, because that sentence is the ADR's brief.

## Sub-agents

An agent `flai serve` starts with `claude-code` is told to keep its own context for decisions and edits and to hand noisy work to sub-agents: search across many files to the explorer, test, lint, and `flai check` runs and long logs to the verifier, and, before it moves its story to review, a check of its diff against the story's criteria and the conventions to a fresh verifier ([ADR-0059](../../design/adrs/0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md)). The convention `design/conventions/delegation.md` says the same to any agent. The template defines both sub-agents for Claude Code:

| File | What it is |
|------|------------|
| `.claude/agents/explorer.md` | Finds and reads: `Read`, `Grep`, `Glob`, and flai's read tools. No shell. |
| `.claude/agents/verifier.md` | The explorer's tools and `Bash`, to run the project's tests, lint, and checks. Told not to edit. |
| `.claude/settings.json` | Runs `flai guard` before every shell command and flai tool call. |

A sub-agent may read and run checks. It may not move, create, or edit a work item, write to a thread, read the inbox, or wait for events or work: neither definition has those tools, and `flai guard` refuses them, and the flai and git commands that write, to any sub-agent, the built-in ones included ([ADR-0060](../../design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)). The story's agent's own calls pass. A refused call tells the sub-agent to say what it needs in its final message instead. The guard looks at every word of a command line, so `env`, `sudo`, `timeout`, `xargs`, `find -exec`, and `bash -c` do not hide a command from it; it is not a shell, and a command hidden on purpose, in a variable or with a backslash in its name, gets past it.

```bash
echo '{"tool_name":"Bash","tool_input":{"command":"flai move S-0001 review"},"agent_type":"verifier","agent_id":"a1"}' | flai guard
# a sub-agent (verifier) cannot run "flai move S-0001 review": ... ; exit status 2
```

A sub-agent primes with `flai prime --story <id> --role explore` or `--role verify`, or the MCP tool `prime` with `role`. Its pack is the conventions whose `roles` list the role, the story's goal and acceptance criteria, and briefs, never bodies, of what the story names and what its topics and links select, as many as fit half the budget; the header counts those left out, and `doc_search` finds them.

A sub-agent that needs the designer puts the question in its final message, and the story's agent asks it with `thread_open`. A sub-agent's calls reach flai under the story's agent's name; in the agent's log under `flai serve`, each of its events carries `parent_tool_use_id`.

## Record recurring friction

```bash
flai issue new "golangci-lint on the host is v1 but the config is v2" --class efficiency --cost 5m
flai issue bump I-0001 --cost 8m --note "reinstalled again in S-0008"
flai issue close I-0001 --reason "scripts/install-tools.sh pins v2"
flai issue list [--all]
flai issue summary
```

Issues live in `design/issues/`, one file per recurring problem with a class (`defect`, `blocker`, `efficiency`, `impression`), a count, an average cost per occurrence, and one dated instance per occurrence. `bump` increments the count, updates the average, and appends the instance. Every command regenerates `summary.md`, the table of open issues most expensive first, and `flai prime` lists it after the conventions when anything is open. `flai check` validates the files and warns when the summary is stale.

## Accept and release

```bash
flai release S-0031 --dry-run          # what a release would look like now
flai accept S-0031 --by alex           # move to done, archive, commit — no release, no tag, no push
flai accept E-0002 --by alex           # an epic: the same
flai push --pending                    # push what was accepted, from the host; releases nothing by default
flai release --pending                 # publish: tag whatever has accumulated and push it
```

Acceptance is one step, and for a story it is the only way to reach done: `flai move S-0031 done` from review, a card dropped on done in the dashboard, and `flai accept S-0031` all run the same flow with the same flags. It rebases the story branch and fast-forwards it into the main branch, moves the item to done (the same rules as `flai move`), archives it with its children and narrative, and commits. It computes no release, creates no tag, and pushes nothing: that is a deliberate step of its own, not tied to any one item, done by whichever of the two commands below you reach for. `--dry-run` prints anything that would block acceptance and any uncommitted files outside `wip/`, and stops without refusing; `--trailer` appends lines such as co-author attribution to the commit message. The working tree must be clean outside `wip/` so the acceptance commit holds only acceptance, unless you pass `--yes`, which includes those files in it. From the dashboard the same choice is a checkbox in the confirmation.

Acceptance then tells the stories still in progress or in review what it changed under them. For each one whose `touches`, with those of its open tasks, cover a path the merge brought into the main branch, it records which paths those are. A story with no touches is told of every path. The command prints `told S-0040 it overlaps: flai/cmd/accept.go`, `--json` lists them in `overlaps`, and the story's agent sees it in its MCP `inbox` as an `overlapped` change. Acceptance from the dashboard does the same.

Acceptance checks what could fail midway before it changes anything: without a git committer identity it refuses and the story stays in review; an experiment story (ADR-0025) is refused outright and stays on its branch, since accepting it onto main is not what an experiment is for. `flai board` says when something is accepted and not pushed (`accepted, not pushed: S-0031 (3 commit(s) ahead of origin/main)`), `flai board --json` carries it as `unpushed`, and one command on the host pushes it. It releases nothing unless the `auto-publish` host action is on for the project (`flai serve enable auto-publish`, off by default, S-0144): then it first computes and tags whatever release has accumulated since each component's last tag, and pushes branch and tags together:

```bash
flai push --pending             # push the branch and any tags already made; with auto-publish on, tag the pending release first
flai push --pending --dry-run   # say what would be pushed, and tagged when auto-publish is on
flai push --pending --publish   # also publish the template when those commits moved its version
```

It only acts when the commits ahead of the remote include an acceptance or a release tag, or, with `auto-publish` on, there is a release to tag from one already pushed; ordinary commits are yours to push with git. It never forces: when the remote has commits this clone lacks it refuses and tells you to fetch and merge first. It answers from what this clone knows, so a push made from another clone is not seen until you fetch. `flai release --pending` is how what has accumulated is published when you choose, so that several acceptances release together: it computes, applies, tags, and pushes, sending the tags on their own when the acceptances were pushed already. The acceptance commit is pushed either way, released or not.

What is pending is worked out from this clone's release tags, and flai never fetches. So before it plans, `flai release --pending` asks the remote for its tags. When the remote has a newer `<name>/vX.Y.Z` than this clone, as after publishing from another clone and pulling without tags, it plans nothing, says which tags are missing, and refuses to publish until you fetch them:

```bash
git fetch --tags origin
```

When the remote cannot be reached, `--dry-run` still shows the plan with a warning that it was not checked, and publishing waits until the remote can be reached, since it pushes there anyway. `flai push --pending` with `auto-publish` on does the same. A clone with no remote publishes locally as before.

An accepted item that no plan can cover is named with the reason (`left out:`) instead of being skipped silently: one touching two components with no tag saying which it delivers to, for example. Tag it, or its epic, and it is planned next time.

A story that is `done` but was never accepted (an older flai, a hand edit) is flagged by `flai check` as `story.unaccepted`, and `flai accept` completes it.

The release follows the git convention. Components are the `projects` in `system-flow.yaml`. The component the item delivers to gets the delivery-type bump: epic major, feature story minor, remediation or improvement patch. It is found from the item's tags (a project name or one of its `tags` aliases), then its epic's tags, and only among the components the item's commits touched: when the tags name several, the one with the most touched files delivers, the earlier tag breaking a tie, and a tag naming a component no commit touched never delivers, so that component gets no release at all. `--deliver` overrides all of that. With no tag deciding, the only touched component delivers, and several ask you for a tag or `--deliver`. `flai check` warns about that earlier (`story.component-tag`): an open story whose touches reach two or more components while no tag of its own or its epic's names one of them, with the `flai edit --tag` that fixes it. Every other component the item's commits touched gets a patch. Code components get an annotated tag `<name>/vX.Y.Z`; a `template` component gets its `template.yaml` version and `CHANGELOG.md` bumped instead. An item whose commits touch no component, such as design or docs work, releases nothing. A research story releases nothing either, whatever it touched: its findings are merged and pushed like any acceptance, and if its commits changed a component's files the plan says that component lands on main without a release. The acceptance commit is pushed whether or not a release was cut.

## Upgrade to a newer template

```bash
flai upgrade --dry-run      # what would change
flai upgrade                # interactive: keep, replace, or diff each conflict
flai upgrade --keep-all     # scripts and CI: never overwrite a project edit
flai upgrade --relock       # a project assembled by hand: record the current files at this version
```

`flai new` writes `system-flow.lock.yaml`, a hash of every file the template rendered. On upgrade each template path is classified: **add** when the project lacks it, **merge** for files with the baseline marker such as `CLAUDE.md` and the conventions (template text above the marker, yours below; `topics` you set on a convention's front matter stay yours, and a convention whose topics you left as the template gave them takes the new template's), **replace** when your copy still matches the lock, **unchanged** when identical, otherwise a **conflict** that you decide. Without a lock every difference is a conflict, which is what `--relock` fixes. A dirty git tree is refused unless `--force`, so an upgrade is one reviewable diff, and the manifest's template version is updated only when no conflict is left undecided. Kept conflicts stay divergent and come back next time; replace them or add your rule below a marker instead.

## Publish a template

```bash
flai template push ./template --dry-run
flai template push ./template            # commit and push to publish.repo / publish.ref from template.yaml
flai template push ./template --tag      # also tag v<version>
```

A template developed inside another repository is published by cloning its remote branch, replacing the contents with the local template, committing with the template version, and pushing. The branch is created from the default branch if it does not exist; `--tag` refuses a tag that already exists. Git errors are shown as git reports them; nothing is retried, and nothing is force-pushed without `--force`. `flai accept` runs this with `--tag` whenever an accepted item releases a template component, so bumping the template and publishing it are one step.

## Run the dashboard

```bash
flai dashboard                 # pull the image if needed, run it, print the URL
flai dashboard --port 8080 --pull
flai dashboard --attach        # follow the logs; Ctrl-C leaves the container running
flai dashboard --build         # build flaiover:local from this monorepo and run that
flai dashboard --bind 127.0.0.1  # this host only (default: every interface)
flai dashboard status
flai dashboard logs [-f]
flai dashboard stop
flai dashboard token           # print the token and login link
flai dashboard token --rotate  # new token; a running dashboard restarts
flai dashboard --no-serve      # do not register with flai serve, or start flai host
```

### flai host: the process that runs flai serve and the MCP servers

`flai dashboard` also starts `flai host`, one process of yours per machine. It runs `flai serve` and each served project's MCP server over HTTP as its children. It starts them again if they end, and stops them all when it stops ([ADR-0040](../../design/adrs/0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)). Its state, token, and log are in a folder named `host` beside flai's config file.

```bash
flai host status            # the host and each process it runs: state, PID, version, restarts
flai host                   # run it in the foreground to watch it; Ctrl-C stops it and everything it runs
flai host start             # in the background (flai dashboard does this for you)
flai host restart serve     # serve, mcp, or all; start <process> and stop <process> too
flai host check             # is a newer flai published?
flai host upgrade           # install it and restart the host and everything it runs on it
flai host stop              # the host and everything it runs
```

One host runs per machine: a second one is refused and told which runs. The host listens on `127.0.0.1:4241` (`FLAI_HOST_ADDR` moves it), where `flai serve` and these commands reach it with the token it wrote. [The operator guide](../operators/index.md#flai-host-the-process-that-runs-the-others-s-0106) has the rest.

### flai serve: flai on the host, for the dashboards

`flai dashboard` registers the project with `flai serve`, a small process of yours on the host that `flai host` runs. `flai serve` opens a connection to the project's dashboard and keeps it open, and the dashboard asks flai for what it needs over that connection; the dashboard never connects to the host. It is how the dashboard reads everything it shows (the board, work items, threads, documents, decisions, the inbox, activity, and search) and how it hears that a file changed; and it is how the dashboard changes anything: every move, save, and acceptance made in the browser is done by `flai serve`, by running the same flai command you would have run, as the project's owner. Without `flai serve` the dashboard says on every page that it has no flai to ask. One `flai serve` serves every project you start a dashboard for.

```bash
flai serve status     # does it run, which projects, which dashboards have it connected
flai serve            # run it in the foreground, outside flai host, to watch it; it keeps no MCP server then
flai serve start      # have flai host run it, starting the host if needed (flai dashboard does this for you)
flai serve stop       # the host stops it, and keeps it stopped until flai serve start
```

Which projects `flai serve` serves is yours to change without starting or stopping the dashboard (S-0121):

```bash
flai serve project add              # serve the project in this folder on the dashboard
flai serve project add ~/git/blog   # or the one in that folder
flai serve project list             # every project served and how it is, and what is not served
flai serve project remove blog      # by key or folder; the dashboard's switcher drops it within a second
```

`add` registers the project as `flai dashboard` does, for the dashboard address your configuration names, and starts `flai host` when it is not running. It refuses a folder with no `system-flow.yaml`, a manifest with no `key`, and a key another served project has, and says which. `remove` touches none of the project's files and leaves the dashboard running for your other projects. A project below a folder named for import, or below the folder `flai serve` was started in, is served from there whether registered or not, so `remove` also puts its folder on `flai serve`'s list of removed projects, which it does not serve from below a folder (S-0123). `add` on a folder in that list takes it off: below a folder `flai serve` serves from, it is served from there again without being registered; elsewhere it is registered as usual. `flai serve import remove` stops every project below an import folder at once. `list` (and `--json`) shows each served project with its key, folder, dashboard address, and whether it is connected, the last error if not, or why it cannot be served: its folder is gone, it has no `system-flow.yaml`, or the manifest has no key. `flai serve` says that once in its log rather than trying it every second, and `flai serve status` shows it too. `list` also shows the projects served because they are below a folder named for import, or below the folder `flai serve` was started in; then the repositories offered for import on the board, and the projects below those folders that it does not serve, each with why. Do not confuse it with `flai serve import`, which names the folders whose repositories the board offers to import.

`flai hostapi` shows what the dashboard can ask, and answers one question on the terminal: `flai hostapi` lists the methods, `flai hostapi board.get '{"all":true}'` prints what the board page is given.

A git repository with a `system-flow.yaml` below a folder named with `flai serve import add` is served too, registered or not, for as long as it is there (a repository imported on the command line before a host ran, say). `flai serve status` lists those it serves under `served from the folders named for import`, and under `not served` each project there it does not serve, with why: you removed it with `flai serve project remove`, its `system-flow.yaml` does not load or has no key, another project has its key, no dashboard is known yet, or `flai serve` is not running.

It needs no root and no configuration. `flai dashboard stop`, like `flai serve project remove`, takes the project out of it and leaves it running for your other projects; `flai serve stop` stops it until `flai serve start`, and `flai host stop` stops it with everything else. `flai dashboard status` has a `host flai` line: connected and since when, or why not. The dashboard shows the same at the right of its header, and "host flai: not connected" there means `flai serve` is not running or cannot reach the dashboard: `flai serve status` says which. Its list of projects, its list of removed projects (`removed.json`), its state, and its log (`serve.log`) are in a folder named `serve` beside flai's config file, `~/.flai/serve` unless `FLAI_CONFIG` points elsewhere.

`flai serve` does what a dashboard asks only among the methods flai offers, and what touches your credentials is off until you turn it on. These *host actions* are yours to enable, by name, in a shell on the host; `push` lets an acceptance made from the board be pushed, and published when you press Publish:

```bash
flai serve actions          # what there is, what each means, and where each is on
flai serve enable push      # for this project; --all-projects for every project
flai serve disable push
flai serve journal          # every host action asked for, and what became of it
```

`auto-publish`, off by default, makes every push of accepted work release first, from the board or `flai push --pending` in any shell that uses this configuration. Off, acceptances wait, unreleased, to be published together from Publish or `flai release --pending` (S-0144).

A second one, `agent`, starts an agent for each story that becomes ready: the harness and model the story names (see [Who works a story](#who-works-a-story-its-agent)), with the program and permissions you set for that harness on the host (`flai serve agent harness`), or a command you wrote for stories that name none (`flai serve agent set -- <program> [args...]`). Turn it on with `flai serve enable agent`. A story in ready or in progress whose agent dropped or failed gets a new one with `flai serve agent restart <story>`, or **Retry** on its page. For a story in ready while the in-progress limit is full, the new agent is queued and starts as soon as there is room. The new agent is told how the last one ended and goes on from the story's narrative. flai serve does not start a story that is held (see [Touches](#touches)), and says why on the board; a restart of one is queued until it is clear. A story in ready gets its agent at once, whatever flai serve's own rules say, held or past the limit, with `flai serve agent start <story>`, or **Start agent** on its page. A story in review whose agent left changes uncommitted in its worktree, which blocks its acceptance, gets an agent to commit them, and do nothing else, with `flai serve agent commit <story>`, or **Have an agent commit them** in the acceptance confirmation or on its review page. A story's agent that runs, or waits for an answer, is stopped with `flai serve agent stop <story>`, or **Stop** on the dashboard's activity page: its process and everything it started end, the story stays where it is with its worktree as the agent left it, and it gets no agent until it is retried or moved back to ready. What a story's agent is saying and doing, read from the log flai serve gives it, is printed by `flai serve agent stream <story>` (`--follow` until it ends), and shown on each card of the dashboard's activity page. From the same logs flai serve measures what the story's agents spent, when an agent ends and when a task of a story whose agent runs enters done, and writes it as the story's and its tasks' `usage` (see [Tokens and cost](#tokens-and-cost)). `flai serve agent usage <story>` prints what the logs say, `--write` records it, and `--all --write` fills in every story flai serve kept a log for, such as those worked before your flai measured them. Run it with the configuration of the flai serve that started the agents (the installed flai, or `--config`), since the logs are in its serve folder. Only Claude Code's logs are read. A third, `host`, lets the dashboard have `flai host` restart `flai serve` or the MCP servers, or upgrade flai and restart on it. A fourth, `settings`, lets the dashboard's Settings page change all of this for you ([the operator guide](../operators/index.md#the-settings-host-action-changing-the-hosts-settings-from-the-dashboard-s-0105) says what that gives the dashboard token). What enabling each means, for who can publish a release and who can start a process on your machine, is in the operator guide; read it first.

The dashboard needs its login token for everything but health and readiness: one token per user, for every project the dashboard serves. `flai dashboard` creates it on first run, as `dashboard.token` in the `serve` folder beside flai's config file, and prints a login link; open the link (or paste the token on the login page) and the browser keeps a session cookie. Tools send it as `Authorization: Bearer`. Details and the exposure table are in the operator guide.

The container runs detached as `flaiover`, one for every project on the host, published on every interface of the host (`--bind`, or `dashboard.bind`, restricts it) on the configured port. Nothing of the project is mounted into it: it is given its port, the login token, and a credential for the host's flai, and everything it shows and changes it asks of `flai serve` on your machine, which `flai dashboard` has `flai host` start alongside it. Commits and acceptances made from the board are therefore made on your machine, as you, wherever the repository lies. Image, tag, port, and bind address come from flags, then the `dashboard` section of `system-flow.yaml`, then `~/.flai/config.json`. If Docker is not installed the command says so with an install pointer. The container holds no git credential: an acceptance made from the board is pushed from the host, by `flai serve` when the `push` host action is enabled, or with `flai push --pending` in a shell. The old `--push-key` flag and `dashboard.push_key` setting are retired and ignored.

The image lives on GHCR and is private while the repository is. When the pull is refused, `flai dashboard` logs Docker into the registry with `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` and retries once. The token needs the `read:packages` scope; `gh auth refresh -h github.com -s read:packages` adds it. Inside the monorepo, `--build` sidesteps the registry by building the image from `flaiover/` as `flaiover:local`.

The image is published for `linux/amd64` and `linux/arm64`. `flai dashboard`, `flai dashboard upgrade`, and `flai dashboard check` pull the image for the platform the Docker daemon runs, such as `linux/arm64` on Apple silicon. If the image already present was built for another platform (for example an amd64 image pulled before arm64 was published), `flai dashboard` pulls again without being asked. A tag published without your platform fails with an error that names the platform; pick a newer tag, or use `--build` in the monorepo.
