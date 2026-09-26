---
title: flai CLI
updated: 2026-09-26
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

```bash
flai board 2>/dev/null                 # output only
LOG_FORMAT=json flai check 2>events.jsonl
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
| `host_actions` | `flai serve enable <action>`, `flai serve disable <action>` | For each host action (`push`, `agent`, `dashboard`, `checks`, `settings`, `host`), the projects it is on for: main checkout paths, or `*` for every project. Absent means none. `flai serve actions` says what each lets a dashboard do |
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
flai move S-0001 review         # needs at least one task
flai move S-0001 done --by alex # needs every task closed and every criterion checked
flai move S-0001 in-progress --reason "tests missing"     # from review
flai move S-0002 cancelled --reason "superseded by S-0005"
flai move E-0003 cancelled --reason "a different route" --dry-run   # what would go with it
```

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
```

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
flai stream sync S-0037        # rebase the branch onto main; run at every task transition
flai stream open S-0037 --no-branch
```

Each story is worked on its own branch, checked out in a worktree under `.flai-cache/worktrees/`. Code, design, and docs changes land there; `wip/` is always written in the main checkout, so the board and the dashboard stay current whatever branches exist. `flai stream sync` rebases the branch onto the main branch, stashing uncommitted work around it; conflicts stop inside the worktree and are listed, resolve them, `git rebase --continue`, and sync again. `flai accept` rebases, fast-forwards the branch into main, removes the worktree and branch, then tags and pushes.

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
flai edit S-0085 --harness claude-code --model claude-sonnet-5 --agent-config effort=high
```

`flai edit` changes what an item says about itself: title, nature, tags, touches, parent, a story's `after` and agent, and the body below its heading, any of them together. What is the item's state stays with its own commands: the status with `flai move`, blocking with `flai block`. A closed or archived item is refused.

A title lives in several places, and a retitle keeps them in step: the front matter, the heading, the file's name, the line in the parent's list, the story's narrative, and links to the old file name under design, docs, and wip (from a story's worktree only under wip, because design and docs there are another branch's). With `--hash`, the one `--show` printed, a change someone made meanwhile is a conflict (exit 3) and nothing is written. `flai check` runs with the change in place: what the change introduces refuses it, every file is put back, and the findings are printed (exit 4). What is simply not allowed, a nature there is not, an epic as a task's parent, is said as a `rule:`. `--autocommit` commits every file the edit touched in one commit; nothing is pushed.

Agents connected over MCP are told of an edit someone else made, as a change of kind `edited` that names what changed. That comes from a small log under `.flai-cache`, outside git, like an agent's read marker: hand edits of a file are not reported, as before.

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
| `inbox` | (Since S-0085 `changes` also reports `edited`: someone changed an item's title, fields, or body with `flai edit` or from the dashboard, and `to` names what.) Threads awaiting the agent (`awaiting: you` when the last entry is not the agent's; `story` filters, `all` includes the rest), `ready`: the stories ready to pull, in pull order, with `can_pull` from the in-progress limit and `held` with why on a story an open story's claim holds or that waits for a story it names in `after`, and `changes`: what others did to work items since this agent last looked (moved, blocked, unblocked, pull order changed), each reported once `unpushed`, on every call while it is true: an acceptance made in this clone and not pushed (items, commits ahead, tags), which the agent pushes from the host with `git fetch` and `flai push --pending` |
| `board` | The board as `flai board --json` prints it, a held ready story with `held` and why; `all` adds epics and tasks |
| `thread_get`, `thread_open`, `thread_reply`, `thread_resolve` | Read, start, answer, and close threads as the agent (`FLAI_AGENT`) |
| `item_get`, `item_move` | Read an item with its children, a story's agent and the project's default, and the hash of its file; transition it with the workflow rules. Moving a story or epic to done is refused: acceptance is yours |
| `item_new`, `item_edit` | Create an epic, a story (with an `agent` over the project's default), or a task; change an item's own words, as `flai edit` does: `agent` replaces a story's agent whole and `clear_agent` removes it, `after` replaces the stories a story waits for and an empty list removes them, and the `hash` from `item_get` refuses a change made meanwhile. Neither commits: the agent commits with its work |
| `doc_get` | A markdown document under the design, docs, or wip folders; nothing else in the repository is served |
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
```

`flai adr new` takes the next number from the files in `design/adrs` (one more than the highest; gaps are not filled), names the file `NNNN-slug.md`, writes `id`, `title`, `status` (`proposed` unless you say `--status accepted`), `date`, `supersedes`, `superseded_by`, and `refines`, adds the row to `design/adrs/README.md`, and sets `superseded_by` on each ADR it supersedes, which is the one edit allowed to an accepted ADR. The body is your template's sections, or standard input with `--body-stdin`. `flai check` runs with everything in place: if it reports anything the ADR introduces, every file is put back, the findings are printed, and the exit code is 4. `--autocommit` makes one `docs: ADR-NNNN <title>` commit of what was written, unless the project sets `dashboard.autocommit: false`; nothing is pushed.

`flai check` warns with `adr.index` when an ADR file has no row in the index or a row has no file. An accepted ADR is immutable: `flai doc` refuses its body, and so does the dashboard's editor.

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
flai stats --json                   # per-item values, weekly throughput, burn-up and cumulative flow series, aging
```

The table shows completed and cancelled counts, throughput per week, current WIP, cycle, lead, and queue time distributions (p50, p85, max, mean), flow efficiency, time-in-state share, aging work against the cycle time p85, and throughput by week. Definitions are in [design/system/metrics.md](../../design/system/metrics.md); the dashboard uses the same numbers.

## Prime a session

```bash
flai prime          # paths of design/conventions in read order, README first
flai prime --cat    # the same files' contents, each under a header
flai prime --json
```

Agents read these before any change; a shell hook or a wrapper can pipe `flai prime --cat` into the session. `flai check` validates the folder: every file needs `title`, `updated`, `audience: agent`, a unique `order`, and `status`; exactly one baseline marker followed by a `## Project additions` section; under 120 lines; and the README must list each file exactly once.

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
flai push --pending                    # tag whatever has accumulated and push it, from the host
```

Acceptance is one step, and for a story it is the only way to reach done: `flai move S-0031 done` from review, a card dropped on done in the dashboard, and `flai accept S-0031` all run the same flow with the same flags. It rebases the story branch and fast-forwards it into the main branch, moves the item to done (the same rules as `flai move`), archives it with its children and narrative, and commits. It computes no release, creates no tag, and pushes nothing: that is a deliberate step of its own, not tied to any one item, done by whichever of the two commands below you reach for. `--dry-run` prints anything that would block acceptance and any uncommitted files outside `wip/`, and stops without refusing; `--trailer` appends lines such as co-author attribution to the commit message. The working tree must be clean outside `wip/` so the acceptance commit holds only acceptance, unless you pass `--yes`, which includes those files in it. From the dashboard the same choice is a checkbox in the confirmation.

Acceptance checks what could fail midway before it changes anything: without a git committer identity it refuses and the story stays in review; an experiment story (ADR-0025) is refused outright and stays on its branch, since accepting it onto main is not what an experiment is for. `flai board` says when something is accepted and not pushed (`accepted, not pushed: S-0031 (3 commit(s) ahead of origin/main)`), `flai board --json` carries it as `unpushed`, and one command on the host finishes the job — computing and tagging whatever release has accumulated since each component's last tag, then pushing branch and tags together, before anything else is decided:

```bash
flai push --pending             # tag whatever has accumulated and push it, branch and tags together
flai push --pending --dry-run   # say what would be tagged and pushed
flai push --pending --publish   # also publish the template when those commits moved its version
```

It only acts when the commits ahead of the remote include an acceptance, or there is a release to tag from one already pushed; ordinary commits are yours to push with git. It never forces: when the remote has commits this clone lacks it refuses and tells you to fetch and merge first. It answers from what this clone knows, so a push made from another clone is not seen until you fetch. `flai release --pending` does the same computing, applying, and tagging on its own, ahead of a push, for seeing or forcing it separately; the acceptance commit is pushed either way, released or not.

A story that is `done` but was never accepted (an older flai, a hand edit) is flagged by `flai check` as `story.unaccepted`, and `flai accept` completes it.

The release follows the git convention. Components are the `projects` in `system-flow.yaml`. The component the item delivers to gets the delivery-type bump: epic major, feature story minor, remediation or improvement patch. It is found from the item's tags (a project name or one of its `tags` aliases), then its epic's tags, and only among the components the item's commits touched: when the tags name several, the one with the most touched files delivers, the earlier tag breaking a tie, and a tag naming a component no commit touched never delivers, so that component gets no release at all. `--deliver` overrides all of that. With no tag deciding, the only touched component delivers, and several ask you for a tag or `--deliver`. Every other component the item's commits touched gets a patch. Code components get an annotated tag `<name>/vX.Y.Z`; a `template` component gets its `template.yaml` version and `CHANGELOG.md` bumped instead. An item whose commits touch no component, such as design or docs work, releases nothing. A research story releases nothing either, whatever it touched: its findings are merged and pushed like any acceptance, and if its commits changed a component's files the plan says that component lands on main without a release. The acceptance commit is pushed whether or not a release was cut.

## Upgrade to a newer template

```bash
flai upgrade --dry-run      # what would change
flai upgrade                # interactive: keep, replace, or diff each conflict
flai upgrade --keep-all     # scripts and CI: never overwrite a project edit
flai upgrade --relock       # a project assembled by hand: record the current files at this version
```

`flai new` writes `system-flow.lock.yaml`, a hash of every file the template rendered. On upgrade each template path is classified: **add** when the project lacks it, **merge** for files with the baseline marker such as `CLAUDE.md` and the conventions (template text above the marker, yours below), **replace** when your copy still matches the lock, **unchanged** when identical, otherwise a **conflict** that you decide. Without a lock every difference is a conflict, which is what `--relock` fixes. A dirty git tree is refused unless `--force`, so an upgrade is one reviewable diff, and the manifest's template version is updated only when no conflict is left undecided. Kept conflicts stay divergent and come back next time; replace them or add your rule below a marker instead.

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

`flai serve` does what a dashboard asks only among the methods flai offers, and what touches your credentials is off until you turn it on. These *host actions* are yours to enable, by name, in a shell on the host; `push` lets an acceptance made from the board be pushed and published:

```bash
flai serve actions          # what there is, what each means, and where each is on
flai serve enable push      # for this project; --all-projects for every project
flai serve disable push
flai serve journal          # every host action asked for, and what became of it
```

A second one, `agent`, starts an agent for each story that becomes ready: the harness and model the story names (see [Who works a story](#who-works-a-story-its-agent)), with the program and permissions you set for that harness on the host (`flai serve agent harness`), or a command you wrote for stories that name none (`flai serve agent set -- <program> [args...]`). Turn it on with `flai serve enable agent`. A story in ready or in progress whose agent dropped or failed gets a new one with `flai serve agent restart <story>`, or **Retry** on its page. For a story in ready while the in-progress limit is full, the new agent is queued and starts as soon as there is room. The new agent is told how the last one ended and goes on from the story's narrative. flai serve does not start a story that is held (see [Touches](#touches)), and says why on the board; a restart of one is queued until it is clear. A story in ready gets its agent at once, whatever flai serve's own rules say, held or past the limit, with `flai serve agent start <story>`, or **Start agent** on its page. A third, `host`, lets the dashboard have `flai host` restart `flai serve` or the MCP servers, or upgrade flai and restart on it. A fourth, `settings`, lets the dashboard's Settings page change all of this for you ([the operator guide](../operators/index.md#the-settings-host-action-changing-the-hosts-settings-from-the-dashboard-s-0105) says what that gives the dashboard token). What enabling each means, for who can publish a release and who can start a process on your machine, is in the operator guide; read it first.

The dashboard needs its login token for everything but health and readiness: one token per user, for every project the dashboard serves. `flai dashboard` creates it on first run, as `dashboard.token` in the `serve` folder beside flai's config file, and prints a login link; open the link (or paste the token on the login page) and the browser keeps a session cookie. Tools send it as `Authorization: Bearer`. Details and the exposure table are in the operator guide.

The container runs detached as `flaiover`, one for every project on the host, published on every interface of the host (`--bind`, or `dashboard.bind`, restricts it) on the configured port. Nothing of the project is mounted into it: it is given its port, the login token, and a credential for the host's flai, and everything it shows and changes it asks of `flai serve` on your machine, which `flai dashboard` has `flai host` start alongside it. Commits and acceptances made from the board are therefore made on your machine, as you, wherever the repository lies. Image, tag, port, and bind address come from flags, then the `dashboard` section of `system-flow.yaml`, then `~/.flai/config.json`. If Docker is not installed the command says so with an install pointer. The container holds no git credential: an acceptance made from the board is pushed from the host, by `flai serve` when the `push` host action is enabled, or with `flai push --pending` in a shell. The old `--push-key` flag and `dashboard.push_key` setting are retired and ignored.

The image lives on GHCR and is private while the repository is. When the pull is refused, `flai dashboard` logs Docker into the registry with `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` and retries once. The token needs the `read:packages` scope; `gh auth refresh -h github.com -s read:packages` adds it. Inside the monorepo, `--build` sidesteps the registry by building the image from `flaiover/` as `flaiover:local`.

The image is published for `linux/amd64` and `linux/arm64`. `flai dashboard`, `flai dashboard upgrade`, and `flai dashboard check` pull the image for the platform the Docker daemon runs, such as `linux/arm64` on Apple silicon. If the image already present was built for another platform (for example an amd64 image pulled before arm64 was published), `flai dashboard` pulls again without being asked. A tag published without your platform fails with an error that names the platform; pick a newer tag, or use `--build` in the monorepo.
