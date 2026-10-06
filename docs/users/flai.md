---
title: flai CLI
updated: 2026-10-06
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

Phases are named for what the time went to: `repo.open`, `repo.list` (every work item, the archive included when the method asks for it), `repo.get`, `board.load`, `board.view`, `release.pending`, `threads.read`, `threads.view`, `narratives.read`, `check.overlap` (the designer's inbox's overlapping touches), `docs.walk`, `doc.read`, `adrs.read`, `search.index`, `search.query`, `agent.state`, `agent.stream`, `manifest.load`, `changes.read`, `item.show`, `cancel.preview`, `accept.preview`, `stream.diff`, `stats.compute`, `order.by`, `promote.candidates`, `release.evaluate`, `encode`; and `exec.<program>.<command>` for each process flai started to answer, such as `exec.git.log`, or `exec.flai.move` for a write, which runs the flai command. Since S-0159 no read starts flai.

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
| `host_actions` | `flai serve enable <action>`, `flai serve disable <action>` | For each host action (`push`, `auto-publish`, `agent`, `dashboard`, `checks`, `settings`, `host`, `plan`, `orchestrate`, `auto-approve`), the projects it is on for: main checkout paths, or `*` for every project. Absent means none. `flai serve actions` says what each lets a dashboard do; `auto-publish` and `auto-approve` are shell only, and the dashboard neither shows nor changes them |
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

## License

```bash
flai license
flai license --json
```

Prints the Bytepunx Shield License, the terms flai, flaiover, and the system-flow repository are distributed under. The text is built into the binary, so what prints is the license of the very flai that prints it, and a copy travels with every copy of flai, as the license requires. It is the same text as `LICENSE.md` at the root of the repository and the dashboard's Host › License page. `--json` prints `{"name", "text"}`.

In short: use is free for 501(c)(3) nonprofits, for educators and students at state-funded educational institutions, for security research on the software, for start-ups under the revenue and funding limits the license states, and for individuals whose work with it grosses under one million US dollars a year; any other commercial use needs a paid license from the licensor (<alex@robsonandmilligan.com>); and nobody, paid or free, may offer system-flow or a fork of it as a hosted service or license it to others. The license text is what binds; this paragraph only points at it.

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
| `--var name=value` | Set a variable. Repeatable. A required variable given empty is refused before anything is written. |
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
flai agent set --role-model explore=haiku --role-model verify=sonnet   # what its sub-agents run
flai edit S-0007 --role-model verify=claude-opus-5-5 --role-config verify.x= --unset-role explore
flai agent clear
```

A story's agent says what will work it: the harness that runs the agent (such as `claude-code`), the model, and options for that harness as `key=value`, which flai passes on without reading them ([ADR-0037](../../design/adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)). The project's default lives in `system-flow.yaml` under `agent`. Every story created while it is set gets a copy in its front matter, with whatever `--harness`, `--model`, or `--agent-config` you give on `flai story new` laid over it. The copy belongs to the story. Changing the default later affects only stories made after it; to move an open story to another model, edit it. `flai show` and `flai edit --show` print a story's agent. Epics and tasks have none.

A project with no default and no story with an agent has no `agent` keys at all, and an older flai reads it as before. Once a story has one, use a flai that knows about agents (`flai self-upgrade`).

### Moving work

```bash
flai move S-0001 ready          # needs acceptance criteria; tasks are not required
flai move S-0001 in-progress    # warns if the WIP limit is exceeded
flai move T-0001 in-progress    # warns while a task of its after is open
flai move T-0001 done           # tasks may skip review
flai move S-0001 review         # needs at least one task and nothing uncommitted in the story's worktree
flai move S-0001 done --by alex # needs every task closed and every criterion checked
flai move S-0001 in-progress --reason "tests missing"     # from review
flai move S-0002 cancelled --reason "superseded by S-0005"
flai move E-0003 cancelled --reason "a different route" --dry-run   # what would go with it
flai move S-0003 backlog        # back from ready, or reopened from cancelled
flai move S-0004 ready          # back from in-progress; it goes last in the ready order
flai move S-0005 ready --yes    # a draft: finalizes it as it moves
```

An item moves back one column from ready, in-progress, review, or cancelled, and never out of done ([ADR-0055](../../design/adrs/0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)). A reopened item is refused while its parent is cancelled; move the parent back first. What its cancellation cancelled under it stays cancelled. Until it closes again it has no completed time, lead time, or cycle time in `flai stats`.

An epic follows its stories (S-0200, [ADR-0076](../../design/adrs/0076-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review.md)). You do not move it by hand: it goes to ready with its first ready story, to in-progress with its first started one, and to review with its last open one, and it goes back only when no other story holds it. Cancelled stories do not count. Each step is recorded on the epic with the story's actor and time and a note such as `follows S-0001, which moved to ready`. `flai move` says so under the item's line, and `--json` returns it as `followed`:

```text
S-0001 → in-progress
  E-0001 → in-progress, following S-0001
```

An epic reaches done only when you accept its last open story (see [Accept and release](#accept-and-release)). You can still move an epic yourself; a story moving forward never pulls it back. The full rule is in [the workflow](../../design/system/workflow.md#an-epic-follows-its-stories). Agents see the epic's move in their MCP `inbox` as `E-0001 <title> moved to in-progress by <who>, following S-0001`, and `item_move` returns it as `followed`.

### Drafts, cost of delay, and forecasts

Stories and epics carry planning data, each part saying who set it and when ([ADR-0074](../../design/adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)).

A draft is a story an agent wrote, such as one the planner drafted or one made from an issue, that you have not read yet. It cannot go to ready until you finalize it:

```bash
flai story new "Export to CSV" --epic E-0001 --draft   # make a story as a draft
flai edit S-0005 --draft                               # mark a backlog story as a draft
flai edit S-0005 --no-draft                            # finalize it where it is
flai move S-0005 ready --yes                           # or finalize it as it moves to ready
```

`flai move S-0005 ready` without `--yes` is refused with "finalize it first". Agents cannot finalize: over MCP, `item_move` refuses a draft to ready and `item_edit` refuses `draft: false`, so an agent tells you in a thread or its narrative that a story is ready to be finalized. The orchestrator is the one exception, and only while you give it `finalize_drafts` ([Running the orchestrator](#running-the-orchestrator)). `flai check` warns (`story.draft`) on a draft that reached ready some other way.

Finalizing records who did it and when, in the story's `finalized` block, so the story still says it was a draft once the flag is gone ([ADR-0077](../../design/adrs/0077-a-story-that-was-a-draft-records-who-finalized-it-and-when-and-the-dashboard.md)). `flai edit --no-draft` records you (`--by`, else `FLAI_AGENT`, else your config author); a finalizing move records the move's `--by` and time. `flai show` prints `finalized by alex at 2026-10-03T10:00:00Z`. `flai edit --draft` makes the story a draft again and removes the block. The dashboard's **Finalize** button, on a draft story's page and in its edit form, does the same as `flai edit --no-draft` ([The item page](flaiover.md#the-item-page)).

The cost of delay says what each week of waiting for a story or an epic costs. Its inputs are yours to give, each optional: the revenue it brings each week once done, the penalty each week it is not done, and the work lost each cycle it is not done. Its value, the cost of delay per week, is the planner's, and you can set it too:

```bash
flai edit E-0001 --revenue-per-week 1200 --penalty-per-week 300 --time-lost-per-cycle 6h
flai edit E-0001 --penalty-per-week ""                 # an empty value removes one input
flai edit S-0005 --cost-of-delay-value 400
flai edit S-0005 --clear-cost-of-delay                 # remove it all
```

You can give the inputs when you make a story or an epic, with the same flags. The inputs are recorded as set by its owner (`--owner`, else your config author) at the time it is made. Without them the item has no cost of delay. The dashboard's new-item form passes its cost of delay panel this way, as you:

```bash
flai story new "Export to CSV" --epic E-0001 --revenue-per-week 1200 --time-lost-per-cycle 2h
flai epic new "Billing v2" --penalty-per-week 300
```

An amount that is not a number, or a duration that is not a Go duration, is refused and nothing is made. The value is not given at creation: it is the planner's, or yours with `flai edit`.

A story made from an issue is a draft and carries the cost of delay inputs the issue gives, set by flai ([Record recurring friction](#record-recurring-friction)).

A forecast says how long a story is expected to take in agent time, when it is expected done, and what that rests on. It sits beside `estimate`, which stays yours:

```bash
flai edit S-0005 --forecast-duration 6h --forecast-delivery 2026-10-09T17:00:00Z --forecast-basis "three tasks like S-0185's"
flai edit S-0005 --clear-forecast
```

Amounts are plain numbers in the project's currency. Durations are Go durations such as `6h` or `90m`, and the delivery is a UTC timestamp. Each edit changes only what you give, and the forecast it changes records you (`--by`, else `FLAI_AGENT`, else your config author) and the time. A cost of delay records who set its inputs and who set its value apart, each when it changes ([ADR-0080](../../design/adrs/0080-a-cost-of-delay-stamps-its-inputs-and-its-value-apart.md)), so changing an input keeps the planner's name on the value. The value is stale when the inputs changed after it: the planner should work it out again. `flai show` prints all three, with `stale: the inputs changed after the value` when it is. The project manifest sets the currency, what an hour of work costs, and the cycle time lost is counted over: `planning.currency` (default `USD`), `planning.hour_rate` (unset), and `planning.cycle` (default `168h`), and the duration a forecast falls back on, `planning.default_duration` (default `1h`), listed in the [settings index](../operators/settings.md#project-manifest).

flai works the touches, the forecast, and the cost of delay value out for you. These commands only print; the planner, or you, records what they print with `flai touches` and the `flai edit` flags above. The model behind each is in [Enriching a story](../../design/system/strategic-agents.md#enriching-a-story-s-0210).

```bash
flai touches suggest S-0005                               # files often changed with its touches
flai touches suggest S-0005 flai/cmd --min 3 --limit 10   # start from more paths, keep fewer files
flai forecast S-0005                                      # how long, when done, and why
flai cod E-0001                                           # an epic's cost of delay from its inputs
flai cod S-0005 --json                                    # a story's, or its share of its epic's
```

`flai touches suggest` takes the main branch's commits that changed the story's touches, its tasks', or the paths you give, and counts the other files each changed. It leaves out the wip folder and flai's own commits, such as acceptances and releases. It lists the files changed in at least `--min` of those commits (default 2), most first, at most `--limit` of them (default 20, 0 for all). Each line gives the file, the commits, and their share. On S-0212 with `--limit 3`:

```text
S-0212 from flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts: 150 of 778 commits changed them
design/system/flai-cli.md  67   45%
docs/users/flai.md         61   41%
docs/operators/index.md    53   35%
```

`flai forecast` sizes the story by its acceptance criteria and touches, and takes the median agent time per unit of size of the done stories most like it: the same nature, model, and size band, falling back to fewer of them until three stories match. For the delivery it plays out the board's pull order and in-progress limit from now. With fewer than three done stories to go on, it gives `planning.default_duration` and says so:

```text
S-0210 forecast 1h6m, delivery 2026-10-04T21:27:00Z
Median 151 s per unit of size over 12 done feature stories on claude-opus-5-5 in the large band, times size 26 (5 criteria, 21 touches); in progress since it started at 2026-10-04T19:35:00Z, so delivery counts from then.
```

`flai cod` adds the revenue and the penalty per week to the hours lost per cycle, priced at `planning.hour_rate`, times the cycles in a week. A story without inputs gets a share of its epic's value, by its forecast duration among the epic's open stories without inputs. It refuses an epic without inputs, and time lost without an hour rate: those are yours to give.

```text
E-0016 cost of delay 1500.00 USD a week
10h of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week: 1500.00 USD a week.
```

A flai older than the one that brought these fields reads past them with a warning and does not act on them, so it would let a draft go to ready. Publishing that release raises `flai.minimum`, so upgrade the flai on your host (`flai self-upgrade`) first.

### Tokens and cost

An item may carry `usage`: the tokens each model read and wrote on it, what they cost in US dollars, and how long agents worked ([ADR-0051](../../design/adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)). You never type it. `flai serve` measures a story, and each of its tasks, from the logs of the agents it started for the story (see [flai serve agent usage](#flai-serve-flai-on-the-host-for-the-dashboards)). A task gets the calls of the sub-agents started for it, known by the task ID their description, or else their prompt, names, and an even share of the story's agent's own calls while it was in progress, split among the tasks in progress at the time ([ADR-0071](../../design/adrs/0071-a-task-s-usage-is-the-calls-of-the-sub-agents-started-for-it-and-an-even-share.md)). Whenever an item enters done, whether by `flai move`, the MCP `item_move`, or `flai accept`, each item above it that was not measured gets the sum of its children's, up to its epic: an epic always sums its stories. `flai show` prints it:

```text
  usage: 20.1M tokens · $8.13 · 18m3s of agent work · measured from its agents' logs
    claude-opus-5-5  input 256 · output 89.3K · cache read 19.7M · cache write 327.6K · $8.1258
```

`(estimated)` after the cost means some of it was not reported by the harness: a task's share of its story's session, or a run that ended without its totals, priced at the rate the logs report for the model. `flai show --json` returns it under `item.usage`, and `flai stats` charts it ([Flow metrics](#flow-metrics)). A flai older than the one that brought `usage` refuses an item that carries it: upgrade the flai on your host first.

Since S-0225 what the planner spent planning an item is charged to it under `usage.strategic`, apart from what its agents spent ([ADR-0083](../../design/adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)). `flai show` prints it on a line of its own per kind of agent, with its models below, and without an agents' line when no agent has worked the item yet. An item with a forecast duration, or else an estimate, also gets its expected cost below its planning: that duration priced at the project's cost per agent hour, always an estimate. It is absent until a story has been measured from its logs. `flai show --json` returns it as `expected_cost`, beside `item`, with `cost`, `from` (`forecast` or `estimate`), and `estimated`.

```text
  forecast: duration 3h · set by planner at 2026-10-04T09:30:00Z
  expected cost: $12.12 (estimated, from the forecast of 3h at $4.04 per agent hour)
  planner, strategic: 812.0K tokens · $0.81 (estimated) · 6m52s, apart from the agents' usage
    claude-opus-5-5  input 12.0K · output 0 · cache read 800.0K · cache write 0 · $0.8100
```

### Cancelling

Cancelling an epic cancels every story under it that is still open and their open tasks; cancelling a story cancels its open tasks. Items that are done or already cancelled are left alone. `flai move` lists what will be cancelled first and, on a terminal, asks before doing it (`--yes` skips the question, `--dry-run` only lists). Each cancelled item records its own transition and a note that names the cause, such as `E-0003 cancelled: a different route`, so an archived task still says why it ended.

```text
Cancelling E-0003 also cancels 3 items:
  story S-0009  in-progress Berths
    task  T-0031  backlog     Dredge
  story S-0010  review      Cranes  (in review: its work stays on its branch, unmerged)
```

A story in review cannot be cancelled directly: accept it or send it back. It is cancelled only with its epic, and the list says so; accept it first if you want the work. An epic in review cannot be cancelled directly either: move it back to in-progress first. Cancelling a story moves its epic forward when the story was what held it back, to review at most: when its other stories are all done, accept the epic with `flai accept E-nnnn`. The list, `--json`, and `--dry-run` say so. Cancelled is final. Nothing of git is touched: the command names each cancelled story's narrative, branch, and worktree and leaves them for you to keep or remove (`flai archive` moves the narratives). If a tree was cancelled by an older flai or edited by hand, `flai check` reports each open item under a cancelled parent as `item.parent-cancelled`.

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

`flai board limit` sets the WIP limit of `ready`, `in-progress`, or `review` in `wip/kanban/board.md`, the one place flai, `flai serve`, `flai check`, and the dashboard read it from; the other columns have none. A move past a limit warns and is made. `flai serve` starts an agent for a ready story only while `in-progress` has room and `review` is under its limit, so raising the in-progress limit lets it start the next one. A full review holds the pull until you accept or send back a story in it; stories already in progress go on (S-0243, [ADR-0073](../../design/adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md)). The board's lane menu in the dashboard does the same.

The `backlog` and `ready` columns are listed in pull order: the stories named in `order` in `wip/kanban/board.md`, in that order, then the rest by ID. Agents pull the first ready story, so this is how you say what comes next.

```bash
flai order S-0061 --top                # first in its column
flai order S-0059 --before S-0061      # just above another story of the same column
flai order S-0047 --after S-0053
flai order S-0056 --bottom
flai order S-0056 --top --placed-by alex   # record the placement as alex's
```

Only ready and backlog stories can be placed, and only relative to a story in the same column; `flai move` changes the column. A story you move to `ready` joins the end of the ready stories. Backlog stories you have never placed stay out of the list and come last, by ID. The dashboard's board does the same thing when you drag a card up or down within a column.

Each placement is recorded in `board.md`, under `placed`, with who made it and when ([ADR-0088](../../design/adrs/0088-board-md-records-who-placed-a-story-by-hand-and-when-and-a-policy-s-order-keeps.md)). Who is `--placed-by`, else `FLAI_AGENT`, else your config author; a drag on the dashboard is recorded as `flaiover`. The record goes when the story leaves its column. A policy's order keeps a story placed by anyone but the orchestrator in the last day where it was put ([Ordering by a policy](#ordering-by-a-policy)), so a story you drag holds against the orchestrator for a day.

#### Ordering by a policy

flai can work out the ready column's order from the stories' planning data ([Drafts, cost of delay, and forecasts](#drafts-cost-of-delay-and-forecasts)), and say which backlog stories could go to ready:

```bash
flai order --by wsjf                   # the ready column by a policy, with each story's figure; writes nothing
flai order --by cod --apply            # the same, written to board.md
flai order --by cod --keep-placed 2h   # keep only the stories placed by hand in the last two hours
flai promote --candidates              # backlog stories that could go to ready, and why the others cannot
flai promote --candidates --limit 3
flai promote --drafts                  # draft stories, and what each lacks to be finalized
flai plan --candidates                 # epics the planner should plan, and why
```

These, and `flai release --evaluate` ([Accept and release](#accept-and-release)), are the orchestrator's arithmetic. They live in flai so that you and the dashboard get the same answer the orchestrator does, with no agent running; the orchestrator itself makes only the judgement calls.

`flai order --by` takes one of four policies, the names `orchestration.policy` takes in `system-flow.yaml` ([settings](../operators/settings.md)):

| Policy | Orders by | First |
|--------|-----------|-------|
| `cod` | The cost of delay value per week | Highest |
| `wsjf` | That value divided by the forecast duration in hours | Highest |
| `throughput` | The forecast duration | Shortest |
| `fifo` | When the story was created | Oldest |

A story without the figure its policy needs, such as a story with no forecast under `throughput`, goes after the stories that have it, in its current order, and the listing says what it lacks. Stories whose figures tie keep their current order. Without `--apply` nothing is written. With it, the order becomes `board.md`'s ready order, as if you had dragged each story into place, but records no placement; nothing is committed. `--by` orders the whole column, so it takes no story and none of `--before`, `--after`, `--top`, `--bottom`, or `--placed-by`.

A ready story placed by hand, by anyone but the orchestrator, within `--keep-placed` keeps its position, and the policy orders the other stories around it. The window is a day unless you give one; `--keep-placed 0` keeps none. The listing marks each kept story `(kept: placed by alex at 2026-10-06T09:00:00Z)`.

`flai promote --candidates` lists the backlog stories that could go to ready, ordered by the project's `orchestration.policy` (`fifo` when it is not set), each with its figure. A backlog story is a candidate when:

- it is not a draft;
- it meets the definition of ready: a goal, acceptance criteria with a checkbox, and an epic that is not cancelled;
- it would not be held if it were ready ([Touches](#touches), [Waiting for another story](#waiting-for-another-story)): it declares touches whenever a story is in progress, its claim overlaps no claim of a story in progress outside the shared paths, and every story it names in `after:` is done;
- it has a forecast duration and a cost of delay value.

Every other backlog story is listed with each reason it is not a candidate. `--limit` caps the candidates listed; those beyond it are left out. It writes nothing: move a candidate to ready with `flai move`.

`flai promote --drafts` lists each draft story in the backlog, in the same policy order, as `complete` or `incomplete` with each thing it lacks. A draft is complete when it has every section of the project's story template, a goal, acceptance criteria with a checkbox, at least one touch, a forecast duration and delivery, a cost of delay value, and an epic that is open, or none. Whether its criteria, touches, and forecast describe the same work is for whoever finalizes it. It writes nothing.

`flai plan --candidates` lists, in ID order, the epics the planner should plan, each with why: an epic in the backlog with no stories, for the planner to draft them, and an epic not done or cancelled whose stories are all done or cancelled with at least one done, for the planner to draft what its outcome still lacks. It lists apart an epic whose planner runs now, or whose planner asked you a question you have not answered yet. It starts nothing.

Each has a `--json` form. The dashboard reads the same answers through `flai serve` (`order.by`, `promote.candidates`), and agents through the MCP tools `order_by_policy` and `promote_candidates` ([Serving agents over MCP](#serving-agents-over-mcp)). `flai promote --drafts` and `flai plan --candidates` have no MCP tool or `flai serve` read; the orchestrator runs them in its shell. A sub-agent, the planner, and the orchestrator may run `flai order --by` without `--apply`, `flai promote --candidates` and `--drafts`, and `flai plan --candidates`, as reads.

### Story branches

```bash
flai stream open S-0037        # narrative, plus branch story/S-0037 in .flai-cache/worktrees/S-0037
flai stream sync S-0037        # rebase the branch onto main, then check it against the other open branches; run after committing each task
flai stream open S-0037 --no-branch
```

Each story is worked on its own branch, checked out in a worktree under `.flai-cache/worktrees/`. Code, design, and docs changes land there; `wip/` is always written in the main checkout, so the board and the dashboard stay current whatever branches exist. `flai stream sync` rebases the branch onto the main branch, and is the only way a story's agent rebases it: agents never start a `git rebase` or `git merge` by hand. `flai accept` rebases, fast-forwards the branch into main, and removes the worktree and branch.

A story's agent works through its tasks one at a time, and keeps the branch close to main as it goes:

1. When a task is done, it commits the task's changes, with their docs and work item updates, on `story/S-0037`.
2. It runs `flai stream sync`, and resolves each conflict sync lists.
3. It runs the tests for what the task changed, and commits any fix they need.

Before it moves the story to review, it commits whatever is outstanding, syncs again, and closes out with `scripts/close-out.sh`, which refuses a branch that does not yet contain the main branch and says to sync.

Every run of `scripts/close-out.sh` ends with one line naming the story, the outcome, and the step it stopped at, such as `close-out: S-0037 passed every step; ready to move to review` or `close-out: S-0037 stopped at markdown lint (exit 1)`, so a run that stops never needs running again to learn why. The agent's verifier runs it once and reads that line, and the agent tells the verifier any step it already knows will stop.

Sync never stashes, so it never has your work in hand when something goes wrong:

- **Uncommitted changes.** A worktree with uncommitted changes is refused before anything is touched. Sync names each path; commit them on the story branch (or stash them yourself) and sync again. A worktree where a rebase is already in progress is refused too.
- **Conflicts.** When the rebase stops on conflicts, it stays stopped in the worktree. Sync prints each conflicting path on a line of its own, then how to continue (in the worktree, resolve each path, `git add` it, run `git rebase --continue`, then sync again) and how to abort (`git rebase --abort`, which puts the branch back as it was before the sync). It exits non-zero; `--json` reports `ok: false` with `worktree`, `uncommitted`, `conflicts`, `rebase_in_progress`, `continue`, and `abort`.

```text
story/S-0037 was not synced: the rebase onto main stopped on conflicts in 2 paths:
  flai/cmd/edit.go
  docs/users/flai.md
To continue: in .flai-cache/worktrees/S-0037, resolve each conflicting path, git add it, and run git rebase --continue; then run flai stream sync S-0037 again
To abort: in .flai-cache/worktrees/S-0037, run git rebase --abort, which puts story/S-0037 back as it was before the sync
```

A story begun on another host reaches your clone with its narrative and tasks, but not its branch or worktree. `flai stream open` on it keeps the narrative and records your host, agent, and session in it. It then checks out `story/<id>`: your clone's branch if you have one, otherwise the remote's, fetched from `origin`, otherwise a new branch from main. It says which (`reopened wip/agents/S-0037.md`, `branch story/S-0037 (fetched from origin) checked out at …`). Only what the other host pushed comes with it. A story that already has its worktree here is refused, as before.

After a clean rebase, sync checks the branch against the other stories in progress or in review, so that two stories that change the same lines find out while both are still open, not when the second is accepted:

```text
story/S-0131 is rebased onto main
story/S-0131 conflicts with story/S-0130 (in progress) in flai/cmd/edit.go; see TH-0024
story/S-0131 merges cleanly with story/S-0129 (in review)
story/S-0131 changed 1 path outside S-0131's touches: flai/internal/threads/threads.go
widen them so that stories that overlap wait: flai touches S-0131 flai/cmd docs flai/internal/threads/threads.go
```

- **Conflicts.** Sync merges the two branches in git's object store only (`git merge-tree --write-tree`, git 2.38 or newer; an older git skips it with a warning), so nothing changes in either worktree. For each pair that conflicts, flai opens one thread on the story that synced, titled `S-0130 and S-0131 conflict when merged`, listing the paths. It shows in both stories' agents' MCP `inbox` and in the designer's inbox on the dashboard. Settle it between the two stories: one narrows its change, or names the other in `after:` and waits. A later sync with the same paths adds nothing, new paths add an entry, and flai resolves the thread once the two merge cleanly or the other story is no longer open.
- **Outside the touches.** Sync lists the files the branch changed since main that the story's claim does not cover: its touches, each folder among them narrowed to the files its tasks name inside it, and its open tasks' touches (see [Touches](#touches)). It prints the `flai touches` command that widens them; when a task changed a file it did not name, widen that task's touches. Touches that are too narrow let a story that overlaps start beside it.
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

`touches` is the list of paths or components a story or task is changing. Write each path from the repository's root; one that starts with a dot, such as `.claude/agents` or `.github/workflows`, is taken like any other. flai drops a trailing slash and a repeat, and refuses an entry with a comma, a leading dash, a leading slash, a `..` segment, or a control character, writing nothing, whichever way you set it: `flai touches`, `--touches` on `new` and `flai edit`, MCP's `item_new` and `item_edit`, or the dashboard. `flai check` warns (`wip.overlap`) when two in-progress items cover the same path outside the [shared paths](#shared-paths), the board prints it under each card, and the dashboard shows a "being worked on" notice on those documents.

It is also a claim that decides what starts ([ADR-0046](../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md), refined by [ADR-0096](../../design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)). A sub-project's name or tag (`cli`, `flai`) means its path. A story's claim is worked out from its touches and its tasks':

- A story touch that is a folder holding touches of the story's tasks is replaced by those touches. Done tasks count, because the branch changed their files; cancelled tasks do not. So `flai/internal/mcpserver` on the story narrows to `flai/internal/mcpserver/permission.go` once a task names that file.
- A story touch that no task names inside stays whole.
- A touch of an open task outside every story touch is added.

A ready story is *held* while its claim overlaps the claim of a story in progress: the same path, or one inside the other (`flai/cmd` and `flai/cmd/serve`, not `flai` and `flaiover`). An overlap that lies wholly inside a shared path does not hold ([Shared paths](#shared-paths)). A story in review holds nothing: its branch is finished and synced, and when it is accepted, the notice tells each overlapping story what changed. A story with no touches may change anything, so it is held while any story is in progress, and while it is in progress itself it holds every ready story. Declare touches when you create a story, a file where you can name one and a folder only where files no task can name yet may be added; the agent that pulls it may widen them, or narrow a folder to the files its tasks name. `flai touches` leaves an edit notice when it changes them, as `flai edit` does, so agents connected over MCP hear of it and `flai serve` may plan the story again ([Running the planner](#running-the-planner)).

A claim grows after its story starts, as its tasks are written. When `flai task new`, `flai edit --touches`, `flai touches`, or MCP's `item_new` and `item_edit` add paths to the claim of a story in progress that another story in progress claims, flai says so (S-0244): `overlaps S-0198 <title> (in progress) on docs/users/flai.md: both stories are told; coordinate with its agent before you change them`. `--json` and the MCP result list them in `overlaps`, and both stories' agents are told as an `overlapped` change. An overlap inside the shared paths is not told. The write stands, and `flai check` still only warns.

A held story is not started by `flai serve` and not offered by `wait_for_work`. The next ready story that is not held goes ahead of it, and it keeps its place and goes first once it is clear. The board says why:

```text
ready
  S-0130 Serve the board faster                         feature          2m HELD
         touches flai/cmd/serve
         held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0128 (in progress) touches; starts when S-0128 moves to review, is cancelled, or is sent back
```

You can still start it yourself: `flai move S-0130 in-progress` and `flai serve agent start S-0130` (**Start agent** on its page) warn and go ahead. An agent you start this way is told that you started it and what it went past, the hold, a full in-progress limit, or a full review, and the journal says the same.

An agent that ended asking you a question is started again when you answer, in its own session, while its story is in progress or in review. When the orchestrator recommends an answer, the agent waits on: it starts again when you confirm the recommendation (`flai thread confirm`) or answer otherwise. If you send the story back to ready before answering, it waits for its hold and the limit like any other ready story once you answer.

### Shared paths

```bash
flai shared list                                  # the patterns, one a line
flai shared check docs/users/flai.md flai/cmd     # whether each lies inside a pattern, and which
flai shared check S-0295                          # each entry of a story's claim, as its hold reads it
flai shared add 'docs/users/*.md' --autocommit    # commits system-flow.yaml on its own
flai shared remove design/adrs
```

A few paths are changed by almost every story, each in a section of its own or in a new file, such as the users' guide or the ADR folder. An overlap there is seldom a conflict, so it should not hold a story. `claims.shared` in `system-flow.yaml` lists them as glob patterns ([ADR-0096](../../design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)):

```yaml
claims:
  shared:
    - docs/users/flai.md
    - design/adrs
```

A pattern is a path from the repository's root, separated by `/`. `*` matches any characters within one segment, `**` zero or more whole segments, and `?` one character. A path with none of them covers itself and everything below it, so `design/adrs` and `design/adrs/**` mean the same. A project made from the template starts with `design/adrs` and `design/issues`, under its design folder.

When two claims overlap, flai takes the narrower entry of the pair, the deeper one. The pair does not hold when that entry lies wholly inside a pattern. A file lies inside a pattern that matches it. A folder lies inside only when all of it does: `docs/users` lies inside `docs/users/**`, but not inside `docs/users/*.md`, since the folder may hold other files. Such an overlap holds no ready story, is not a `wip.overlap`, and is not told as a grown claim. The trial merge at `flai stream sync` and the notice at acceptance still report it, so a real conflict on a shared path is caught, after both stories have started. flai reads the list again each time it needs it, so a change takes effect at once in a running `flai mcp` or `flai serve`.

`flai shared check` changes nothing and exits 0 whether or not an entry is shared. It reads a component's name or tag as its path. An entry whose last segment has an extension, such as `flai.md`, is a file; any other is a folder. Given a story's ID, it reports each entry of the story's claim, which tells you what a pattern would free before you add it:

```text
S-0295 docs/users/flai.md: shared, inside docs/users/flai.md
S-0295 flai/internal/workitem/hold.go: not shared
```

`flai shared add` and `remove` rewrite only `claims.shared` and keep the file's other keys and comments. A pattern that is not valid (empty, absolute, with a `..` segment, or a malformed glob), one to add that is listed already, or one to remove that is not listed is refused with the reason, and nothing is written. Each prints what it changed. Nothing is committed unless you give `--autocommit` and the project leaves `dashboard.autocommit` on; `--trailer` adds lines to the commit. `flai check` reports a pattern that is not valid as `manifest.claims` on its line, and the pattern frees nothing until it is fixed.

With `--json`, `list` prints an array of patterns; `add` and `remove` print `{"added": [...], "shared": [...]}` or `{"removed": [...], "shared": [...]}`, `shared` being the list after the change; and `check` prints an array of `{entry, path, shared, pattern, story}`, with `pattern` only when the entry is shared and `story` only for a story's claim.

The dashboard's settings page ([flaiover.md](flaiover.md#settings)) and the MCP tools `shared_paths` and `shared_paths_edit` ([Serving agents over MCP](#serving-agents-over-mcp)) read and change the same list.

The list decides what holds, so an agent that could add to it could free its own story. Only your own session changes it. `flai guard` refuses `flai shared add` and `remove`, and `shared_paths_edit`, to every sub-agent and to every session `flai serve` starts: a story's agent, the planner, the orchestrator, and the analyzer. The refusal tells the agent to ask you on a thread. The dashboard changes it while you have the `settings` host action on. Listing and checking are open to every agent.

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

A flai older than the one that brought `after` refuses to read a story that carries it, so upgrade the flai on your host (`flai self-upgrade`) before you use it. `flai story new --after` and the MCP `item_new` tool set it when the story is made.

### Planning a story's tasks

```bash
flai task new "Store the plan" --story S-0176
flai task new "Show the plan" --story S-0176 --after T-0677          # starts once T-0677 is done
flai edit T-0679 --after T-0677,T-0678                                 # replaces them
flai edit T-0679 --clear-after
```

A task's `after` names the tasks of the same story it waits for: the plan the story's agent writes with the tasks. Tasks that wait for nothing undone, and whose `touches` do not overlap, can be worked at the same time.

`flai show` prints a story's plan below its children: each task's state, and the layers of tasks that can run at once.

```text
  plan:
    T-0677  done
    T-0678  done
    T-0679  in-progress  after T-0677
    T-0680  ready        after T-0677
    T-0681  waiting      for T-0679, T-0680
  layers:
    1  T-0677, T-0678
    2  T-0679, T-0680
    3  T-0681
```

A task is `ready` to start when it has not started and every task of its `after` is done or cancelled, `waiting` while one of them is open (`for` names those), `in-progress` while it is in progress or in review, then `done` or `cancelled`. A layer is the tasks with the same longest chain of `after` steps before them: layer 1 waits for none, and each layer can run at once when the ones before it are done. A cancelled task is in no layer and holds no one up; neither is a task on a cycle, or one that waits on a cycle, which `flai check` reports. An `after` entry that names no task of the story is left out, for `flai check` to report. `flai show --json`, the MCP `item_get` tool, and the dashboard's `item.get` give the same as `plan` beside `item` and `children`, for a story with tasks: `tasks`, each with `id`, `state`, and, when they are not empty, `after` and `waiting_for`; and `layers`, lists of task IDs in ID order, the first waiting for none.

`flai board` counts a story's tasks under its card, while the story is open, and `--json` gives every story with tasks `tasks`: `ready`, `waiting`, `in_progress`, `done`, and `layers`, the number of layers.

```text
         tasks 1 ready, 1 waiting, 1 in progress, 2 done; 3 layers
```

flai does not hold a task the way it holds a story. `flai move T-0681 in-progress` while a task of its `after` is open warns, and moves it:

```text
T-0681 is waiting (after): waits for T-0679 (in progress) and T-0680 (ready); ready to start when T-0679 and T-0680 are done or cancelled
```

`flai task new --after` and `flai edit --after` run `flai check` with the change in place, and a finding refuses the change and leaves nothing written (exit 4). `flai check` reports (`task.after`) an entry that names no task, a task of another story, the task itself, or a story, and every cycle among a story's tasks (`T-0677` waits for `T-0678`, which waits for `T-0677`), once. The MCP `item_new` and `item_edit` tools set it too.

A flai older than the one that brought a task's `after` reports every task that carries one as an error, so upgrade the flai on your host first. A release with it raises `flai.minimum`, so an older flai stops before it reads any item.

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
flai thread reply TH-0001 --source "design/system/overview.md#Principles" "Answer"   # cite what it rests on
flai thread reply TH-0001 --recommend --source design/system/overview.md "Answer"
flai thread confirm TH-0001                            # make the recommendation the answer
flai thread resolve TH-0001 --reason "settled in ADR-0021"
```

A thread is one file under `wip/threads/`, anchored to a document, a heading in it, or a work item, with dated entries by author. The author is `--by`, else `FLAI_AGENT`, else the config author. A reply from anyone but the opener marks the thread `answered`; the opener's follow-up makes it `open` again; `resolve` closes it. Unresolved threads on a story or its tasks are mirrored into the story narrative under `## Open questions`, so an agent sees them without the dashboard. `flai check` validates threads: the anchor must exist, a named heading must still be in the document, and open threads on archived items are flagged.

A reply can cite a source with `--source <path>` or `--source <path>#<heading>`: a file in the repository, or an item ID, and a heading in it. Both must exist, or the reply is refused. The source is the entry's last line, `Source: <path> § <heading>`. A reply with `--recommend` is a recommendation (ADR-0090). Its heading ends `(recommendation)`, and the thread keeps its status, so it still awaits you. `flai thread show` marks it and says it awaits you. `flai thread confirm` makes it the answer: it adds an entry by you that names the recommendation and cites its source, and marks the thread `answered`. To answer otherwise, reply as usual. Confirm is refused when no recommendation is pending, and to the one who made it. The orchestrator replies this way when `orchestration.permissions.answer_threads` lets it.

### Changing an item after it was made

```bash
flai edit S-0085 --show                         # the fields, the body below the heading, and the hash
flai edit S-0085 --title "A better name" --autocommit
flai edit S-0085 --nature improvement --tag dashboard --tag cli --touches flaiover/src
flai edit S-0085 --parent E-0004                # an open epic for a story, an open story for a task
flai edit S-0085 --body-stdin --hash <hash> < body.md
flai edit S-0085 --clear-tags --clear-touches
flai edit S-0085 --after S-0084                 # hold it until S-0084 is done; --clear-after lets it go
flai edit T-0679 --after T-0677                 # a task's plan: the tasks of its story it waits for
flai edit S-0085 --topics logging               # what it is about; --clear-topics removes them
flai edit S-0085 --harness claude-code --model claude-sonnet-5 --agent-config effort=high
flai edit S-0085 --no-draft --revenue-per-week 800   # finalize a draft and give a cost of delay input
```

`flai edit` changes what an item says about itself: title, nature, tags, touches, parent, a story's or epic's `topics`, a story's or a task's `after`, a story's agent, a story's draft flag and forecast, a story's or epic's cost of delay (see [Drafts, cost of delay, and forecasts](#drafts-cost-of-delay-and-forecasts)), and the body below its heading, any of them together. What is the item's state stays with its own commands: the status with `flai move`, blocking with `flai block`. A closed or archived item is refused.

A story given another epic with `--parent` moves the epic it joined as though the story had just entered it from backlog, and the epic it left as though the story had been cancelled out of it: either goes forward, never back, and never to done. `flai edit` prints each move, indented, as `E-0004 → in-progress, following S-0085`, and `--json` lists them in `followed`.

A title lives in several places, and a retitle keeps them in step: the front matter, the heading, the file's name, the line in the parent's list, the story's narrative, and links to the old file name under design, docs, and wip (from a story's worktree only under wip, because design and docs there are another branch's). With `--hash`, the one `--show` printed, a change someone made meanwhile is a conflict (exit 3) and nothing is written. `flai check` runs with the change in place: what the change introduces refuses it, every file is put back, and the findings are printed (exit 4). What is simply not allowed, a nature there is not, an epic as a task's parent, is said as a `rule:`. `--autocommit` commits every file the edit touched in one commit; nothing is pushed.

Agents connected over MCP are told of an edit someone else made, as a change of kind `edited` that names what changed. That comes from a small log under `.flai-cache`, outside git, like an agent's read marker: hand edits of a file are not reported, as before.

They are also told when an accepted story changed paths their own story claims, as a change of kind `overlapped` on their story. `cause` is the accepted story and `to` lists the paths. The agent syncs its story and runs its tests again before it goes on. See acceptance, below. Since S-0244 an `overlapped` change whose `cause` is a story in progress, not an accepted one, means a write grew the two stories' claims to overlap on `to` ([Touches](#touches)): the agent agrees with that story's agent on a thread before it changes those paths, and narrows its touches if it can.

### Ticking acceptance criteria

```bash
flai criteria list S-0282                     # the criteria, numbered, ticked or not
flai criteria tick S-0282 1 3                 # tick the first and the third
flai criteria tick S-0282 1,3 --autocommit    # the same, committed
flai criteria untick S-0282 2
```

A story's acceptance criteria are the checkboxes under its `## Acceptance criteria` heading. `flai criteria` names them by number, from 1, in the order they appear, as `list` prints them. `tick` and `untick` change those boxes and nothing else in the file, with the checks `flai edit` makes: a closed or archived item is refused (exit 4), a change made meanwhile is a conflict with `--hash` (exit 3), and a number with no box is refused and nothing is written. Agents connected over MCP are told the criteria changed.

The story's agent ticks a criterion once it has verified it, never by editing the file, and leaves one it cannot verify unticked with the reason in the story's notes. A task sub-agent says which criteria its task meets; the story's agent ticks them after its review. `flai move` to review warns when a criterion is unticked, and the move to done refuses one ([ADR-0089](../../design/adrs/0089-acceptance-criteria-are-ticked-through-flai-by-number-by-the-story-s-agent-once.md)). The MCP tool `criteria_tick` and the dashboard's `POST /api/items/<id>/criteria` make the same change.

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
| `inbox` | (Since S-0181 `flai_outdated`, on every call while it is true: the flai serving the agent is older than the newest flai release in the project's history, with `running`, `newest`, and the `upgrade` command.) (Since S-0085 `changes` also reports `edited`: someone changed an item's title, fields, or body with `flai edit` or from the dashboard, and `to` names what. Since S-0132 it reports `overlapped`: a story was accepted, `cause`, that changed paths this story claims, `to`. Since S-0244 `cause` may be a story in progress whose claim and this story's grew to overlap on `to`: coordinate with its agent.) Threads awaiting the agent (`awaiting: you` when the last entry is not the agent's, but `other` on a thread it opened while a recommendation on it awaits your confirmation; `story` filters, `all` includes the rest), `ready`: the stories ready to pull, in pull order, with `can_pull`, false while the in-progress limit is full or review is at or over its limit, and `pull_hold` saying which, and `held` with why on a story the claim of a story in progress holds, and `changes`: what others did to work items since this agent last looked (moved, blocked, unblocked, pull order changed), each reported once. `unpublished`: the IDs of accepted items no release has covered yet, as information; publishing them is the operator's, or an agent's the operator asks (S-0195, ADR-0067) |
| `board` | The board as `flai board --json` prints it, a held ready story with `held` and why, a story with tasks with their counts in `tasks`; `all` adds epics and tasks |
| `thread_get`, `thread_open`, `thread_reply`, `thread_resolve` | Read, start, answer, and close threads as the agent (`FLAI_AGENT`) |
| `item_get`, `item_move` | Read an item with its children, a story's agent and the project's default, a story's task `plan` when it has tasks (see [Planning a story's tasks](#planning-a-storys-tasks)), and the hash of its file, and its `draft`, `cost_of_delay`, and `forecast` with the `currency`; transition it with the workflow rules. Moving a story or epic to done is refused: acceptance is yours. So is moving a draft to ready: finalizing is yours |
| `item_new`, `item_edit` | Create an epic, a story (with an `agent` over the project's default), or a task; change an item's own words, as `flai edit` does: `agent` replaces a story's agent whole and `clear_agent` removes it, `after` sets what an item waits for, a story's stories or a task's tasks of the same story (a creation that sets it, or that gives a `body`, is checked as `flai story new --body-stdin` is: a finding the item introduces, such as a missing section or a markdown lint rule, refuses it and leaves nothing), on an edit replacing them, and an empty list removes them, and the `hash` from `item_get` refuses a change made meanwhile. `draft` makes a story a draft, and `draft: false` is refused. `cost_of_delay` and `forecast` set the keys given; `clear_cost_of_delay` and `clear_forecast` remove them, and to remove one amount the agent clears the cost of delay and gives the keys to keep. Neither commits: the agent commits with its work |
| `criteria_tick` | Tick and untick a story's acceptance criteria by number, `tick` and `untick` naming them as `flai criteria list` does, with the `hash` from `item_get`; returns the criteria after. Nothing is committed. Sub-agents, the planner, and the orchestrator are refused it ([Ticking acceptance criteria](#ticking-acceptance-criteria)) |
| `issue_story` | Make a backlog story from an open issue, as `flai issue story` does: a draft carrying the issue's cost of delay inputs, named in the issue's Remediation section. `id` is the issue, `epic` puts the story under an epic, and `story` reads the issue from that story's worktree. Returns the story's `id`, `title`, `nature`, `path`, `draft`, and the `issue`. A closed issue, or one an open story already links, is refused, naming that story. Nothing is committed ([Record recurring friction](#record-recurring-friction)) |
| `doc_get` | A markdown document under the design, docs, or wip folders; nothing else in the repository is served. With `heading`, only that section and the sections below it, with its heading path and line ([Read design on demand](#read-design-on-demand)) |
| `doc_search` | The sections of the design and docs folders, the conventions among them, that rank highest against `query`: at most 20 (`limit` for fewer), each with its path, the document's title, its heading path, line, first lines, and size ([Read design on demand](#read-design-on-demand)) |
| `prime` | A story's context pack, as `flai prime --story <id> --json` prints it, fitted to `budget` (default the project's `prime.budget`, else 80 KB): its topics, the conventions with the sections those topics leave out taken out, what the story names whole (a large document it names only by a path written out as a brief), briefs of the design and tech files and the ADRs its topics and one link step select, ranked sections to fill the budget, each with its reason and size, and a catalog of the rest, to read with `doc_get` and a `heading` when needed; with `role` (`explore` or `verify`), the smaller pack for a sub-agent ([Sub-agents](#sub-agents)); with `role` `plan` and `epic` or `story`, or `role` `orchestrate` or `analyze`, a strategic agent's pack ([The planner, the orchestrator, and the analyzer](#the-planner-the-orchestrator-and-the-analyzer)) ([Prime a session](#prime-a-session)) |
| `who_touches` | In-progress and in-review items whose `touches` cover a path |
| `shared_paths` | The shared paths, `claims.shared` in `system-flow.yaml`, as `shared`, and as `invalid` each pattern that is not valid, with the reason. Given `paths` (paths or touches entries) or `story` (a story's ID, whose claim is checked entry by entry), `entries` says of each whether it lies inside a pattern and which, as `flai shared check` does. A read, open to every agent ([Shared paths](#shared-paths)) |
| `shared_paths_edit` | Remove the patterns in `remove`, then add those in `add`, as `flai shared remove` and `add` do, and return `added`, `removed`, and `shared`, the list after the change. Nothing is committed. Yours alone: `flai guard` refuses it to every sub-agent and every session `flai serve` starts, and tells the agent to ask you on a thread ([Shared paths](#shared-paths)) |
| `order_by_policy` | The ready column's order by `policy` (`cod`, `wsjf`, `throughput`, or `fifo`; the project's `orchestration.policy` when left out), as `flai order --by <policy> --json` prints it, with each story's figure. It never writes the order ([Ordering by a policy](#ordering-by-a-policy)) |
| `promote_candidates` | The backlog stories that could go to ready and why each other one cannot, as `flai promote --candidates --json` prints them; `limit` caps the candidates. It writes nothing ([Ordering by a policy](#ordering-by-a-policy)) |
| `release_evaluate` | Whether the release policy is met, with its figures, as `flai release --evaluate --json` prints it. It releases nothing ([Whether a release is due](#whether-a-release-is-due)) |
| `release_publish` | For the orchestrator alone, while you give it `publish`: publish what is accepted and not yet released, as the board's Publish does, when the release policy allows it, with a one-sentence `reason`. It refuses, changing nothing, when the policy is not met, when `whole_epics` holds the batch back, under `judgement` without a reason, and while the `push` host action is off. Returns the versions, tags, and items released ([When the orchestrator publishes](#when-the-orchestrator-publishes)) |
| `agent_start`, `agent_restart` | Start a story's agent on the host, as `flai serve agent start` and `restart` do, so that your own agent can give a story begun on another host an agent here ([ADR-0064](../../design/adrs/0064-a-story-in-ready-or-in-progress-with-no-agent-run-on-this-host-is-started-here.md)). Only while the operator has turned on the `agent` host action for the project, as for the dashboard's buttons; otherwise, and whenever flai would refuse the command, the tool's error says why. Each call is journalled with the agent that made it. A sub-agent cannot call them, and nor can an agent `flai serve` started: starting agents is your word, not a story's agent's |
| `plan` | Start the planner for an epic or a story on the host, as `flai plan` does ([Running the planner](#running-the-planner)). Only while the operator has turned on the `plan` host action for the project; otherwise, and whenever flai would refuse the command, the tool's error says why. Returns the run: its agent, PID, log, and session. Each call is journalled with the agent that made it. A sub-agent, the planner, and an agent `flai serve` started cannot call it: planning is the operator's to ask for. The orchestrator alone may, for an epic in the backlog, while you give it `plan_backlog_epics` ([Running the orchestrator](#running-the-orchestrator)) |
| `activity_log` | For the planner, the orchestrator, and the analyzer: log an activity that just ended, with `kind`, a one-line `summary`, and the `items` it touched. flai measures its seconds and cost from the agent's run log and appends it to `wip/agents/<kind>.md`. Charges its cost under `usage.strategic`: a planner's to the item its run planned, an orchestrator's split evenly between the epics, stories, and tasks named, each summed up to its epic; what no item takes is left in the kind's project strategic total. Returns the entry, the document's totals, and what the cost was charged to (`charged_to`, `charge`). A sub-agent cannot call it. Nothing is committed ([What they did: activity documents](#what-they-did-activity-documents)) |
| `permission_prompt` | Not for the agent to call: Claude Code calls it, in a session `flai serve` starts, when a call would ask a person. It allows an `Edit`, `Write`, `MultiEdit`, or `NotebookEdit` of a file in a `.claude/` folder in an in-progress story's worktree once the story's owner or the project's owner answers `allow` on a thread, or at once under `auto-approve`, and refuses everything else ([Writes under .claude/](#writes-under-claude)) |
| `wait_for_work` | What to do when you have nothing to work on. Answers at once with `resume` and your own story if one is still in progress, `thread` and the threads awaiting you that were written to since it last answered, or `pull` and the first ready story that is not held when the in-progress limit leaves room and review is under its limit. Otherwise it waits until one of those is true, up to `timeout_seconds` (300, 5 minutes, by default; at most 1800, 30 minutes), and then says whether it was waiting for room (`room`), for a story in review to be accepted or sent back (`review`), for a held story to be clear (`held`: each ready story says why), or for a story to be ready: call it again. Hold it whenever you are idle, and you pull the next story as soon as there is one |
| `wait_for_events` | Returns at once when something changed since this agent last looked, otherwise blocks until a thread, item, or narrative changes, or `timeout_seconds` passes (60 by default; at most 1800, 30 minutes). Returns `events` in the same shape as `changes`, and the changed paths. Hold it to wait for the designer's answer on a thread, or, as the orchestrator does, between decisions. It does not see a sub-agent finish: a story's agent waits for one by launching it with the Agent tool's `run_in_background` set to false, and `flai guard` refuses a story's agent this call while a sub-agent of its session runs and no thread on its story is open ([ADR-0092](../../design/adrs/0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md)) |

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

"Since this agent last looked" is a marker per agent name (`--agent`, else `FLAI_AGENT`) under `.flai-cache/mcp/`, outside git. It only decides which changes are news; ready work and open threads are listed on every call, so nothing depends on it. One look reports at most 50 changes, the newest, and says in `changes_omitted` (`events_omitted` for `wait_for_events`) how many older ones it left out; those are not reported later. The first look under a new name covers the last 24 hours and tells of stories and epics only, not task transitions: to an agent that has just arrived, a day of task moves is history, and `board` and `item_get` show how things stand. An agent that ends its turn between your messages calls `inbox` when it starts again and hears what you did in between; one that stays running holds `wait_for_events` and hears within a second. A wait is held for the `timeout_seconds` asked, up to 30 minutes (S-0244): every return is a model turn that reads the agent's whole context again, about 0.05 USD at 150k cached tokens, so a long wait is a cheap one. While it waits, the server sends a progress notification every minute, because Claude Code drops a tool call that sends nothing for 30 minutes over stdio and 5 over HTTP. `flai move` records `FLAI_AGENT` as who moved an item when it is set, so an agent is not told about its own moves; a move on the dashboard's board is recorded as the project's `owner`.

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

A heading's topics cover everything down to the next heading at its level or higher; a heading without them takes its parent's, and the top headings take the file's. `all` is every story, and a convention without `topics` is read as `[all]`. A convention's `roles` list every agent that reads it: `story` for the agent working a story; `explore` and `verify` for its sub-agents ([Sub-agents](#sub-agents)); and `plan`, `orchestrate`, and `analyze` for the planner, the orchestrator, and the analyzer ([The planner, the orchestrator, and the analyzer](#the-planner-the-orchestrator-and-the-analyzer)). A convention without `roles` is read by every agent. `flai prime --story` prints the conventions whose roles are empty or list `story`, `flai prime --role` those whose roles are empty or list the role, and `flai prime --cat` prints them all. `flai check` accepts these six roles and warns with `conventions.roles` on any other. `flai check` warns with `doc.topics` on a `design/system` or `design/tech` file without topics, and with `doc.topic` on a topic that is not `all`, `code`, a sub-project's name, tag, or kind, or one that a story or epic declares ([A story's topics](#a-storys-topics)). `flai prime --story` selects conventions, design, tech files, and ADRs by them ([Prime a session](#prime-a-session)). A design, tech, or ADR file without topics comes into a pack only when something names or links it, or it ranks.

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
flai check --strict   # warnings exit 1 too, save review over its limit and an epic behind its stories; use this in CI
flai check ../other-repo --json
flai check --strict --story S-0249                   # findings outside S-0249 are notes that do not fail the run
flai check --strict --story S-0249 --record-issues   # and each rule's are recorded in an issue
```

Every finding is one line, `path:line: level: rule: message`, so editors and CI annotate it. Rules cover the manifest and layout, every work item in `kanban/` and `archive/` (front matter, IDs and file names, parents and children, state history, acceptance criteria, required sections), narratives and their index, the board's WIP limits and pull order, and front matter on `design/` and `docs/` files including ADR numbering. `README.md` files are exempt from front matter.

`--strict` fails on every warning but two. One is `board.wip-limit` for review over its limit: only you clear it, by accepting or sending back a story (S-0243, [ADR-0073](../../design/adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md)). The other is `epic.lags-stories`, an open epic that its stories put further on than its status, because it was moved before S-0200 or moved back by hand: the message names the `flai move`s that catch it up, or `flai accept E-nnnn` when its stories are all done, and only you move an epic ([ADR-0076](../../design/adrs/0076-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review.md)). Neither must stop an agent's close-out of another story. Both are still printed, the summary line ends `(N that --strict passes over: only the operator clears them, by accepting or by moving an epic)`, and `--json` counts them in `advisory` as well as `warnings`. Ready or in-progress over its limit still fails `--strict`.

`--story S-nnnn` scopes the run to one story, as a story's close-out does (S-0249, [ADR-0085](../../design/adrs/0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)). A finding is inside the story when it is on the story's or one of its tasks' files, its narrative, a thread anchored on the story or one of its tasks, a path its branch changes against the main branch, or a path uncommitted in its worktree. Every other finding is outside it, and so is every `wip.overlap`, which the pull hold and the other story's agent clear: another story's unmerged branch, a thread answered on an archived story, or an error on the main branch no longer stops the close-out of a story that did not cause it. A finding outside keeps its level and is printed with `(outside S-nnnn)`, but neither an error nor `--strict` fails on it. The summary line counts them, and `--json` marks each `outside` and counts them in `outside`.

`--record-issues`, with `--story`, records each rule's findings outside the story in `design/issues`, so you can see where they occur. Each rule has one open issue, titled ``flai check finds `<rule>` outside the story at close-out``: the first finding opens it with class `efficiency`, and each later story that meets the rule bumps its count, with an instance naming the story and the findings. Running the check again for the same story and findings records nothing new. The issue is written in the checkout the run reads, which at close-out is the story's worktree, and `summary.md` is regenerated:

```text
$ flai check --strict --story S-0249 --record-issues
wip/archive/kanban/stories/S-0173-<slug>.md:4: warning: story.unaccepted: S-0173 is done but its branch story/S-0173 was never merged; merge or delete it (outside S-0249)
1139 items checked, 0 errors, 1 warnings; 1 outside S-0249, notes the run passes over
recorded story.unaccepted outside S-0249 in I-0070 (opened)
```

`scripts/close-out.sh` passes both flags to the check through `scripts/check.sh`, and commits the issues it records. Run without `--story`, as CI and `make` do, the check counts every finding.

When the project has a markdownlint configuration at its root (`.markdownlint.yaml`, `.yml`, `.json`, or `.jsonc`, or a `.markdownlint-cli2` file's `config`), `flai check` also lints every markdown file under the wip folder with it and warns on each finding, `markdown.MD024` and the like, with markdownlint's own message: work items, threads, and narratives are written in the main checkout, where a story's own lint never runs, and would otherwise reach CI unlinted. flai checks the markdownlint rules what it writes can break (headings, blank lines, trailing spaces, lists, emphasis, code spans, fences, bare URLs); your CI's markdownlint still checks the rest. Without a configuration nothing is linted.

The same lint guards what flai writes there. `flai story new`, `flai task new`, `flai epic new`, `flai edit`, `flai thread new`, `reply`, and `resolve`, `flai stream log`, and the MCP tools that do the same refuse a body or entry that would bring a finding, name the rule and the line, and write nothing. A title loses a trailing `.`, `,`, `;`, `:`, or `!`, which a heading may not end with. A reply and a resolution by one author in the same second share one entry heading.

Whether or not the project has a markdownlint configuration, `flai check` reports an error, `markdown.conflict-marker`, on each line of a markdown file under the design, docs, and wip folders, or at the root, that opens or closes a merge conflict: one that begins with `<<<<<<<`, `|||||||`, or `>>>>>>>` followed by a space or the line's end. Such a line means a conflict was added to git unresolved. Resolve it, keeping what both sides meant, and remove the markers. A line of `=======` alone is not reported, since markdown uses it to underline a heading.

A front-matter field flai does not know, on a work item, a thread, or an issue, is an error: `item.unknown-field`, `threads.unknown-field`, or `issues.unknown-field`, on the field's line. A newer flai wrote it, or it is misspelled. Everything else reads past it: the board, `flai serve`, and the MCP tools list the item, log a warning naming the file and the field, and keep the field when they write the file. When the flai on your host is older than the project, upgrade it ([Keeping the host's flai current](../operators/index.md#keeping-the-hosts-flai-current)).

## Flow metrics

```bash
flai stats                          # stories completed in the last 30 days
flai stats --since 90d --by nature  # grouped
flai stats --type task
flai stats --json                   # per-item values, weekly throughput, burn-up and cumulative flow series, aging, usage, forecasts, cost of delay, waiting, claims
flai stats --since 7d --bucket hour --json   # spend over time by the hour (day by default, or week)
```

The table shows completed and cancelled counts, throughput per week, current WIP, cycle, lead, and queue time distributions (p50, p85, max, mean), flow efficiency, time-in-state share, aging work against the cycle time p85, and throughput for each week of the window, weeks with nothing done included. When items done in the window carry `usage` ([Tokens and cost](#tokens-and-cost)), a `usage` line adds what agents spent on them: tokens, cost, and agent time; then what that comes to per item, with the agent time an item took, per minute of agent work, and per dollar; then per model its tokens, cost, tokens per minute of agent work, and how many items it worked on. `--json` has every item of the type under `items`, done in the window or not, and burn-up and cumulative flow for each day of the window, from the day that holds its start, or the first item's creation if later (S-0166). It has each item's usage under `items[].usage`, and under `usage` the totals, the models, and the items done in order with what had been done and spent by then (`done`, and per model `by_model`). Since S-0163 `usage.spend` lays out what was spent on epics, on stories, and on tasks over time, whatever `--type` is: one point per `--bucket` (`hour`, `day`, or `week`; a day unless you say, and an hour only over a window of 31 days or less), from the first bucket in which an item with usage was done to now, each with the items done in it, their tokens, cost, and agent seconds, the tokens, cost, and agent minutes per item (`minutes_per_item`, since S-0169), the tokens per agent minute and per dollar, the running means per bucket, and the same per model. An item counts in the bucket in which it entered done. Rates are per minute of agent work (`tokens_per_minute`); `tokens_per_hour` is still there for scripts written before. Definitions are in [design/system/metrics.md](../../design/system/metrics.md); the dashboard uses the same numbers.

Since S-0205 `flai stats` also reports what the planner's figures, cost of delay, waiting, and claims come to ([Planning, waiting, and claims](../../design/system/metrics.md#planning-waiting-and-claims-s-0205)). The table adds, each only when there is something to show: the absolute error of forecasts, delivery dates, and estimates (p50 and p85); the cost of delay outstanding in each column now and incurred over the window; the time the agents of the items done in the window waited on threads and in review, in all and per item; the time those stories were held in ready, the items in progress now against the board's limit, and how far the stories' touches were from the files their commits changed; and what the planner, the orchestrator, and the analyzer spent over the window beside the items done. `--json` has them under `forecasts`, `cost_of_delay`, `waiting`, `claims`, and `strategic_days`, with each item's values under `items`. The touches drift reads git: a story's commits are those on the main branch and the `story/` branches whose subject names it in brackets. Where git cannot be read, `claims.drift` is left out and `flai stats` logs a warning saying so; the rest is reported as usual.

Since S-0225 what the planner spent planning an item is charged to that item and the items above it, and `flai stats` reports it apart from what agents spent, never in the agents' totals or per-model figures ([Strategic usage](../../design/system/metrics.md#strategic-usage)). The table adds a `strategic usage, apart` line with its tokens, cost, and time over the items done in the window, and its cost per kind of agent, and a line with the project's cost per agent hour: the agents' cost over their hours across every story measured from its logs, archived ones included. `--json` has each item's under `items[].usage.strategic`, the totals under `usage.strategic`, each type's and each bucket's under `usage.spend.<type>.strategic`, and the rate under `usage.cost_per_agent_hour`. An item with a forecast duration, or else an estimate, has under `items[].expected_cost` what it is expected to cost at that rate, marked estimated, so it shows before any agent works it; it is absent until a story has been measured.

Since S-0226 an orchestrator activity is charged too, split evenly between the work items it named, and what no item carries is its kind's project strategic total ([ADR-0095](../../design/adrs/0095-an-orchestrator-activity-s-usage-is-charged-evenly-to-the-work-items-it-named.md)). For each strategic agent, `flai stats` prints beside its document's totals what the items carry and the project total, cost and seconds each, which add up to the totals: the project total holds the activities that named no item, the planner's before S-0225 or with no planned item, and the analyzer's for now. `--json` has them under `strategic[].items` and `strategic[].project` ([Strategic agents](../../design/system/metrics.md#strategic-agents-s-0206)).

```text
strategic agents (all time):
  planner: 2 activities, 0.5213 USD, 823 s (on items 0.4213 USD, 723 s; project 0.1000 USD, 100 s), last 2026-10-03T18:00:00Z
```

## Prime a session

```bash
flai prime          # paths of design/conventions in read order, README first
flai prime --cat    # the same files' contents, each under a header
flai prime --json
flai prime --story S-0137         # S-0137's context pack, fitted to 80 KB
flai prime --story S-0137 --budget 120KB
flai prime --story S-0137 --json
flai prime --story S-0137 --role verify   # the pack for a verifier sub-agent
flai prime --role plan --epic E-0016      # the planner's pack for an epic
flai prime --role plan --story S-0207     # the planner's pack for a story
flai prime --role orchestrate             # the orchestrator's pack
flai prime --role analyze --json          # the analyzer's pack, as data
```

Agents read these before any change. An agent with a story primes with `flai prime --story <id>`, or the MCP tool `prime`, which returns the same pack as `--story --json` and takes `budget` too: the prompt `flai serve` gives the agents it starts says so, as do the MCP server's instructions, `CLAUDE.md`, and `session-start.md`. They also tell the agent that a brief is not the document, and to read the section that bears on the story with `doc_get` and its `heading` before relying on it or changing what it describes. Without a story an agent primes with `flai prime --cat`; a shell hook or a wrapper can pipe either into the session. A sub-agent the story's agent starts primes with `--role` and `--story` ([Sub-agents](#sub-agents)). The planner, the orchestrator, and the analyzer prime with `--role` alone ([The planner, the orchestrator, and the analyzer](#the-planner-the-orchestrator-and-the-analyzer)).

`--story` prints what an agent working that story needs: its context pack, fitted to a size budget ([ADR-0047](../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), [ADR-0049](../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)). The budget is `--budget`, else `prime.budget` in `system-flow.yaml`, else 80 KB: bytes, or a number with `KB` or `MB`. It counts everything printed, the header included.

The header names the story and its topics, with where each came from (as `flai show` gives them). It gives the pack's size in bytes and lines against the budget, and one line per thing printed with its size. When the conventions alone exceed the budget, it says so, and the pack is the conventions and a catalog. When the conventions and what the story names exceed it, it says that, and nothing is ranked. When the briefs take it over, it says that too: every brief is kept and nothing is ranked (TH-0032). Until a project narrows its conventions' topics, a pack of a code story is usually over 80 KB this way.

The pack, in order:

1. Each convention the story's agent reads, those whose `roles` are empty or list `story`, as `--cat` prints it, with every section whose topics include neither `all` nor one of the story's left out ([Topics on documents](#topics-on-documents)). The front matter, the baseline marker, and the `## Project additions` heading always stay, and so does the heading above a section that is kept. Then the open issues. Nothing here is cut for the budget.
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

### The planner, the orchestrator, and the analyzer

These agents work above a story ([ADR-0075](../../design/adrs/0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md)), and each primes with its role:

| Agent | Command | Primed for |
|-------|---------|------------|
| The planner | `flai prime --role plan --epic E-nnnn` or `--story S-nnnn` | The epic or story it plans |
| The orchestrator | `flai prime --role orchestrate` | The whole project |
| The analyzer | `flai prime --role analyze` | The whole project |

The MCP tool `prime` takes the same `role`, with `epic` or `story` for the planner. The pack is shaped like a story's agent's and fits the same budget:

1. Its topics: the role's, `planning`, `orchestration`, or `analysis`, and for the planner its item's as well.
2. The conventions whose `roles` are empty or list the role, with the sections those topics leave out taken out, and no README. Then the open issues.
3. For the planner, what its item names, whole: a story's with its epic and tasks, an epic's alone. The planner reads an epic's stories as items. Then the sections ranked highest against the item.
4. Briefs of the design, tech, and ADRs the topics select, and of the ADRs one step reaches; then the catalog.

The header names the role and, for the planner, the item. `--json` gives the planner's epic or story as `item`, and has no `story`. The orchestrator and the analyzer refuse `--epic` and `--story`. The planner needs one of them, not both, and refuses a task, or an ID of the other type than its flag. `--epic` without `--role plan` is refused.

#### Running the planner

The planner plans one epic or one story and ends ([ADR-0082](../../design/adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)). For an epic with no stories it drafts the stories that deliver its outcome. For a story it adds the paths it will touch, a forecast, and a cost of delay value worked out from your inputs. For an epic with stories it revisits each one not done or cancelled and enriches it again (S-0209). Every story it creates is a draft in the backlog, for you to read and finalize; flai refuses it a story that is not one.

It works a story's figures out with `flai touches suggest`, `flai forecast`, and `flai cod` ([Drafts, cost of delay, and forecasts](#drafts-cost-of-delay-and-forecasts)), keeps every touch the story declares, and changes a figure only with a reason (S-0210). It records where each touch came from and why each figure stands under a `### Planning` heading in the story's Notes. That heading is the planner's, rewritten on each run; the rest of the Notes is yours.

On an epic it opens one thread, the plan: the stories, their order (their `after`), and the assumptions it made. When the epic already had stories, the same thread proposes each one it would split, merge, add, or drop. It drafts the additions only; it never cancels a finalized story or rewrites its words without asking, so the rest is yours to decide there.

For a story it also plans the tasks (S-0255). With none, it drafts the tasks that deliver the story's outcome, each with its work, what done means, the paths it touches, and the tasks it waits for ([Planning a story's tasks](#planning-a-storys-tasks)), in the backlog. With tasks, it revisits each one not done or cancelled and adds what the outcome still lacks. It opens one thread on the story with the plan: the tasks, their order, and the assumptions it made. There it proposes any task it would split, merge, or drop; it does not cancel a task, or rewrite one it did not write, without asking there. The agent that pulls the story reviews the planner's tasks before it works them, and changes what it would plan differently.

You start it, behind a host action that is off until you turn it on:

```bash
flai serve enable plan     # for this project; --all-projects for every project
flai plan E-0016           # plan an epic now
flai plan S-0208 --json    # or a story, and print the run as data
flai serve journal         # every planner started, or that could not be
```

**Plan** on an epic's or a story's page in the dashboard does the same, and so does the MCP tool `plan` from your own agent. `flai plan` needs `flai serve` running to record how the run ended and log what it did; without one it says so, and the planner still runs.

The planner runs in the project's main checkout, as `planner-<id>`, with the planner's agent: `planning.agent` in `system-flow.yaml`, laid over the project's `agent`, so you name only what differs, such as a cheaper model:

```yaml
planning:
  agent:
    model: claude-sonnet-5
    config:
      effort: medium
```

With `claude-code` it runs as the project's `.claude/agents/planner.md`, which the template ships; `flai upgrade` adds it to an older project. A command you set with `flai serve agent set` gets `FLAI_ROLE=plan` and `FLAI_ITEM` instead of a story.

It writes items and threads through flai and nothing else. `flai guard` holds it to that: it may create stories, epics, and tasks, edit them, set their touches, open and answer threads, record issues, and move an item to backlog. A story it creates must be a draft: `flai story new` needs `--draft` and `item_new` needs `draft: true`. It may not move an item further, finalize a draft, accept, publish, commit, or use Claude Code's `Edit` and `Write` ([Sub-agents](#sub-agents) has how the guard reads a command line). When it needs an input that is yours, such as a story's cost of delay inputs, it asks in a thread on the item and waits for your answer.

flai refuses, and says why, while the `plan` action is off, for a task, for an item done, cancelled, or archived, while a planner already runs for the item, and when the planner's agent names no harness and no command is set. A planner does not count against the in-progress limit and holds no story back. When it ends, `flai serve` records how: `worked`, `asked` (it ended with its question to you open; the answer does not start it again, so ask for another run), or `failed`. It logs what the run did and cost in `wip/agents/planner.md`. The entry's summary is the planner's last line, which on an epic names the stories it created and revisited, and on a story the tasks. Its items are the planned item, then the items under it created during the run, then those changed during it: an epic's stories and their tasks, or a story's tasks. Its output is in `serve/agents/<key>-planner-<time>.log` beside flai serve's state.

While `plan` is on, `flai serve` also plans again on its own, so that forecasts keep up with the board ([ADR-0084](../../design/adrs/0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md)):

| When | What happens |
|------|--------------|
| You or an agent edit the goal, criteria, or touches of a backlog or ready story that has a forecast or a cost of delay value, or change its cost of delay inputs after its value | The planner runs for that story. A story never planned is not planned on an edit: planning it first is yours to ask. The planner's own edits start nothing |
| A story is accepted or cancelled, or the pull order changes | What `planning.replan` says. By default flai plays the board out again with no agent and moves the delivery of each ready and backlog forecast that changed, keeping its duration, in one commit, `chore: replan forecasts after …` |
| `planning.schedule` comes round | The planner runs for every ready story |

```yaml
planning:
  replan: deterministic      # the default; or never, or agent
  schedule: "0 6 * * 1-5"    # five-field cron in UTC, or daily for 00:00 UTC; unset, no schedule
```

`never` does nothing when work ahead completes or the order changes, and `agent` does what `deterministic` does, then runs the planner for each story whose delivery moved. You set both keys by hand; [project-manifest.md](../../design/system/project-manifest.md) has them, and `flai check` reports a value flai cannot read.

flai serve sees the edits made with `flai edit`, `flai touches`, the dashboard, and the MCP tool `item_edit`; a hand edit of a file is not seen. It runs one planner at a time per project. A story queued more than once runs once, with every trigger, and one whose planner already runs, such as one you asked for, waits its turn. A queued story that can no longer be planned, such as one accepted meanwhile, is dropped. flai serve acts only on what happens while it runs with `plan` on: what passed while it was down, or while `plan` was off, is not acted on.

What it costs: `deterministic` uses no agent, but each forecast it moves takes about a second, because flai checks the project before and after each write, as `flai edit` does. `agent` and a schedule start planner sessions you pay for without asking each time, so both are off until you set them.

To see it, the dashboard's Settings page shows, for the project, whether edits start the planner, the replan policy, and the schedule with its next run ([Settings](flaiover.md#settings)). Each planner run's entry in `wip/agents/planner.md` has a `- Trigger:` line saying what started it: `asked` when you asked, `orchestrator` when the orchestrator did, otherwise such as `edited goal, touches by alex`, `accepted S-0210`, `cancelled S-0213`, `reordered`, or `schedule daily`, joined by semicolons when several came together. [Planning again](../../design/system/strategic-agents.md#planning-again) has the whole of it.

#### Running the orchestrator

The orchestrator keeps the project's work moving ([ADR-0087](../../design/adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md)). It asks the planner to plan an epic, finalizes drafts, promotes stories to ready, orders the ready column, answers threads, accepts stories, and publishes releases, but each only while you give it the permission. It takes every figure from flai: the epics to plan from `flai plan --candidates`, the drafts complete enough to finalize from `flai promote --drafts`, the order from `flai order --by` ([Ordering by a policy](#ordering-by-a-policy)), the stories that could go to ready from `flai promote --candidates`, and whether a release is due from `flai release --evaluate` ([Whether a release is due](#whether-a-release-is-due)). It never edits a file and never works a story.

It runs behind a host action that is off until you turn it on:

```bash
flai serve enable orchestrate    # for this project; --all-projects for every project
flai serve disable orchestrate   # stop it
flai serve journal               # each start, end, failure, and stop
```

The Settings page of the dashboard turns it on and off too, while `settings` is on. Within a minute of turning it on, `flai serve` starts one orchestrator for the project, in its main checkout, as `orchestrator`. It does not end by itself: after each decision it waits for the next change. When a run ends all the same, `flai serve` starts another, a minute later if the run failed. Within a minute of turning it off, `flai serve` stops it. A start flai cannot make, because the orchestrator's agent names no harness and no command is set, or the project has no `.claude/agents/orchestrator.md` (`flai upgrade` adds it), is in the journal once.

Its agent is `orchestration.agent` in `system-flow.yaml`, laid over the project's `agent` as `planning.agent` is for the planner. With `claude-code` it runs as the project's `.claude/agents/orchestrator.md`. A command you set with `flai serve agent set` gets `FLAI_ROLE=orchestrate` and no story. The run does not count against the in-progress limit. Its output is in `serve/agents/<key>-orchestrator-<time>.log` beside flai serve's state.

What it may do is yours to give, one permission at a time, under `orchestration.permissions`. Each is off until you set it:

```yaml
orchestration:
  policy: wsjf                 # how it orders the ready column
  permissions:
    promote_to_ready: true
    order_ready: true
    answer_threads: recommend  # off, recommend, or autonomous
  agent:
    model: claude-sonnet-5
```

| Permission | It may |
|------------|--------|
| `plan_backlog_epics` | Ask for the planner, one epic at a time, on each epic `flai plan --candidates` lists, with the MCP tool `plan`. The `plan` host action must be on too |
| `finalize_drafts` | Finalize a draft that `flai promote --drafts` finds complete and it judges consistent, with `item_edit` giving only `draft: false`, or `flai edit --no-draft`, and change nothing else with it. On any other draft it opens one thread saying what is missing |
| `promote_to_ready` | Move the stories `flai promote --candidates` lists to `ready`, in its order, while the ready column is under its WIP limit |
| `order_ready` | Order the ready column by `orchestration.policy` with `flai order --by <policy> --apply`, after each change to it. A story you placed by hand in the last day keeps its place |
| `answer_threads` | Reply on threads: `recommend` replies with a recommendation for you to decide on, `autonomous` with an answer of its own. `off`, the default, leaves threads to you |
| `accept_reviews` | Accept a story in review with `flai accept --by orchestrator`, once its verifier passed and nothing blocks it ([When the orchestrator accepts](#when-the-orchestrator-accepts)) |
| `publish` | Publish what is accepted and not yet released, through the MCP tool `release_publish` alone, when the release policy allows it. The `push` host action must be on too ([When the orchestrator publishes](#when-the-orchestrator-publishes)) |

Without any, it reads the board and the inbox, opens threads, records issues, and logs. `flai guard` holds it to its permissions, reading them from `system-flow.yaml` at each call, so a change applies at its next call with no restart. A call that a permission would allow is refused while that permission is off, and the refusal names it: `it needs orchestration.permissions.publish, which is off`. Anything else that writes is refused whatever you give it: editing files, committing, moving a story anywhere but `ready`, or to `done` as it accepts it, changing anything of an item but its draft flag, placing a story by hand in the pull order, and every other flai command that writes. Either way it is told to ask you on a thread rather than work around the refusal.

flai holds it to the four permissions above itself too, so a call the guard does not see is held all the same (S-0219). In the orchestrator's session, `flai plan` and the MCP tool `plan` refuse an epic `flai plan --candidates` does not list; `flai edit --no-draft` and `item_edit` refuse a draft that is not complete, naming what it lacks; `flai move` and `item_move` refuse a story that is not a candidate, with the candidates' reasons, and any story while ready is at its WIP limit; and `flai order --by --apply` is refused without `order_ready`. flai holds it to `accept_reviews` the same way ([When the orchestrator accepts](#when-the-orchestrator-accepts)). A refusal ends that attempt: it logs it and does not try again until something changes. A planner it starts records `orchestrator` as what started it, in its run and its entry in `wip/agents/planner.md`. It never places a story by hand, and its policy order keeps a story you placed in the last day where you put it ([ADR-0088](../../design/adrs/0088-board-md-records-who-placed-a-story-by-hand-and-when-and-a-policy-s-order-keeps.md)). `flai check` reports a permission it does not know, or an `answer_threads` that is none of its three values.

Where to see what it did:

| What | Where |
|------|-------|
| Each decision: what it did, on which items, why, and the policy figure behind it, such as a story's cost of delay value, its value over its duration, its forecast, or a candidate's rank | `wip/agents/orchestrator.md`, under `## Log`, one entry per decision, with its seconds and cost. A run that ends logs the time since the last decision as one more entry; a run you stopped says `stopped: orchestrate turned off` |
| Each call the guard refused it: when, the call, and the permission it needs, or `none` | `wip/agents/orchestrator.md`, under `## Refusals` |
| Its runs, starts, failures, and stops | `flai serve journal`, and the dashboard's Activity page, which shows the run and what it is saying |

[The orchestrator](../../design/system/strategic-agents.md#the-orchestrator) has the whole of it.

#### What they did: activity documents

Each of these agents has one activity document in the project: `wip/agents/planner.md`, `wip/agents/orchestrator.md`, and `wip/agents/analyzer.md` ([ADR-0079](../../design/adrs/0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md)). flai writes them; do not edit them by hand. A document appears with its agent's first activity, in the main checkout, and is committed with the work around it.

Its front matter holds the totals: `kind`, `accrued_cost` in US dollars, `accrued_seconds`, `tasks_completed` (the number of activities logged), and `last_run`, when the newest one ended. Under `## Log` is one entry per activity, newest last: when it ended, a one-line summary, for a planner run `flai serve` started what triggered it, the items it touched, its seconds, and its cost. The orchestrator's document also has `## Refusals`, after the log, with each call `flai guard` refused it ([Running the orchestrator](#running-the-orchestrator)); a refusal is not an activity and adds nothing to the totals. A cost marked `estimated` is the run's cost apportioned to the activity, as a task's is ([Tokens and cost](#tokens-and-cost)).

flai measures each activity from the log `flai serve` keeps of the agent's run. An agent that does several things in one run, such as the orchestrator, reports each when it ends with the MCP tool `activity_log`, giving its kind, a summary, and the items it touched. When a run ends, flai logs the time since the last entry as one activity, so an agent that does one thing and ends need not call the tool. `wip/agents/index.md` lists the documents under a heading of their own, `flai check` validates them, and `flai stats` reports their totals, what of them the items carry and the project total, and, with `--json`, their entries ([Flow metrics](#flow-metrics)).

## Sub-agents

An agent `flai serve` starts with `claude-code` is told to keep its own context for decisions and edits and to hand noisy work to sub-agents: search across many files to the explorer, test, lint, and `flai check` runs and long logs to the verifier, and, before it moves its story to review, a check of its diff against the story's criteria and the conventions to a fresh verifier ([ADR-0059](../../design/adrs/0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md)). While it works it runs only the tests for what it changed. The whole suite, the lint, and `flai check` are the verifier's: one run before review, and one more after the agent fixes what that one found. The agent makes the fixes itself, never a sub-agent. Since S-0176 it is also told to plan its story's tasks as it writes them, giving each its `touches` and the tasks it waits for (`--after`), and to hand each task to a task sub-agent, running a layer of tasks that wait for nothing undone and share no path at once only when the tasks are long; it reviews, commits, and moves each task itself, and a task sub-agent edits only what its task touches. The convention `design/conventions/delegation.md` says the same to any agent. The template defines both sub-agents for Claude Code:

| File | What it is |
|------|------------|
| `.claude/agents/explorer.md` | Finds and reads: `Read`, `Grep`, `Glob`, and flai's read tools. No shell. Runs `haiku`. |
| `.claude/agents/verifier.md` | The explorer's tools and `Bash`, to run the project's tests, lint, and checks. Told not to edit. Runs `sonnet`. |
| `.claude/agents/planner.md` | The planner, which `flai plan` runs as its session's own agent, not a sub-agent ([Running the planner](#running-the-planner)). Writes through flai only: no `Edit` or `Write`. |
| `.claude/settings.json` | Runs `flai guard` before every shell command and flai tool call, and, in a planner's or the orchestrator's session alone, before `Edit`, `Write`, and `NotebookEdit`. In a story's agent's session (`FLAI_STORY` set) it also runs `flai guard` as each sub-agent starts and stops (`SubagentStart`, `SubagentStop`), so that the guard knows which of the session's sub-agents still run. |

A story that changes one of these files ships it from its own branch: its agent edits the file with `Edit` or `Write`, and you allow the write on a thread ([Writes under .claude/](#writes-under-claude)).

Both sub-agents run a model cheaper than the story's agent's, since they read and run checks and decide nothing. Set `model` in a definition to change it, or to `inherit` to run it on the story's agent's model. To change it for one story, or for every new story, without editing the definitions, give the agent a role: `flai agent set --role-model verify=claude-sonnet-5-5`, or `flai story new` and `flai edit` with `--role-model`. `flai serve` then starts the story's session with that role's sub-agent on that model.

A sub-agent may read and run checks. It may not move, create, or edit a work item, write to a thread, read the inbox, or wait for events or work: neither definition has those tools, and `flai guard` refuses them, and the flai and git commands that write, to any sub-agent, the built-in ones included ([ADR-0060](../../design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)). The story's agent's own calls pass, save the wait below. A refused call tells the sub-agent to say what it needs in its final message instead. The guard looks at every word of a command line, so `env`, `sudo`, `timeout`, `xargs`, `find -exec`, and `bash -c` do not hide a command from it; it is not a shell, and a command hidden on purpose, in a variable or with a backslash in its name, gets past it.

```bash
echo '{"tool_name":"Bash","tool_input":{"command":"flai move S-0001 review"},"agent_type":"verifier","agent_id":"a1"}' | flai guard
# a sub-agent (verifier) cannot run "flai move S-0001 review": ... ; exit status 2
```

The story's agent waits for a sub-agent by launching it with the Agent tool's `run_in_background` set to false, a layer's sub-agents in one message, so that each result comes back as the tool's result however long it runs. It never ends its turn while a sub-agent runs in the background: Claude Code ends a headless session ten minutes after its turn ends, and the sub-agent with it. Nor does it wait for one with `wait_for_events`, which reports work items, threads, and narratives, which a sub-agent does not write, so the wait would run to its timeout. Since S-0285 `flai guard` refuses it that call while a sub-agent of its session runs and no unresolved thread is on its story or one of the story's tasks, naming the sub-agents and saying how to wait ([ADR-0092](../../design/adrs/0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md)). It knows which run from the `SubagentStart` and `SubagentStop` hooks, which it records in `.flai-cache/guard/<session_id>.json` in the main checkout and which never refuse. A wait with a thread on the story open passes and returns on the designer's answer. The planner, the orchestrator, and a session you run yourself are not held to this, since it needs `FLAI_STORY` and no `FLAI_ROLE`; and what the guard cannot read lets the call through.

```bash
export FLAI_STORY=S-0001
echo '{"hook_event_name":"SubagentStart","session_id":"s1","agent_id":"a1","agent_type":"verifier"}' | flai guard
echo '{"hook_event_name":"PreToolUse","session_id":"s1","tool_name":"mcp__flai__wait_for_events"}' | flai guard
# the story's agent cannot call wait_for_events while its sub-agents run (verifier a1) and no thread on S-0001 awaits the designer: ... ; exit status 2
```

A sub-agent primes with `flai prime --story <id> --role explore` or `--role verify`, or the MCP tool `prime` with `role`. Its pack is the conventions whose `roles` are empty or list the role, the story's goal and acceptance criteria, and briefs, never bodies, of what the story names and what its topics and links select, as many as fit half the budget; the header counts those left out, and `doc_search` finds them.

A sub-agent that needs the designer puts the question in its final message, and the story's agent asks it with `thread_open`. A sub-agent's calls reach flai under the story's agent's name; in the agent's log under `flai serve`, each of its events carries `parent_tool_use_id`.

## Record recurring friction

```bash
flai issue new "golangci-lint on the host is v1 but the config is v2" --class efficiency --cost 5m
flai issue bump I-0001 --cost 8m --note "reinstalled again in S-0008" [--story S-0008]
flai issue close I-0001 --reason "scripts/install-tools.sh pins v2"
flai issue list [--all] [--story S-0008]
flai issue story I-0001 [--epic E-0002] [--story S-0008] [--owner olive] [--autocommit]
flai issue summary
```

Issues live in `design/issues/`, one file per recurring problem with a class (`defect`, `blocker`, `efficiency`, `impression`), a count, an average cost per occurrence, and one dated instance per occurrence. `bump` increments the count, updates the average, and appends the instance. Every command regenerates `summary.md`, the table of open issues most expensive first, and `flai prime` lists it after the conventions when anything is open. `flai check` validates the files and warns when the summary is stale. It also warns about an open issue that no open story links once it is older than `issues.story_after` in `system-flow.yaml`: 168h, seven days, unless the project sets another duration, or `0` to turn the warning off. A story links an issue when its body names the issue's ID; `flai issue story <id>` makes one that does.

Each instance names the story it was recorded for. `new` and `bump` take it from `--story`; without the flag, from `FLAI_STORY`, else from `FLAI_AGENT` when it has the form `agent-S-nnnn`, else from the story branch, `story/S-nnnn`, checked out where the command runs. Outside any story, the instance names none. `list --story S-0008` lists only the issues with an instance recorded for that story. It reads them from the story's worktree when it has one: the issues a story records are committed on its branch, so they are there, and not in the main checkout, until the story is accepted. Without a worktree, or without `--story`, issues are read from where the command runs. Because each story records its issues on its own branch, `new` numbers an issue one past the highest it finds on the main branch, in the main checkout, in every story's worktree, even before it is committed, and on every `story/S-nnnn` branch, so two stories worked at once never take the same number. `list --json` gives every issue, closed ones too, with `stories`, the stories its instances name, and `story`, the open story that links it, or empty when none does.

`flai issue story I-0001` makes a backlog story from an open issue, the way `flai story new --body-stdin` makes one. The story takes the issue's title. Its nature is `remediation` for a defect or a blocker and `improvement` otherwise. Its goal links the issue and carries the solution the issue's Remediation recommends. That link is what ties the issue to the story. The story is a [draft](#drafts-cost-of-delay-and-forecasts): you finalize it before it can go to ready.

When the issue gives them, the story carries cost of delay inputs, set by flai at the time it is made. `time_lost_per_cycle` is the issue's cost times its count, divided by the planning cycles (`planning.cycle`, default `168h`) since it was first reported, at least one, rounded to the minute. An issue can also give inputs in an `## Impact` section, one per line:

```markdown
## Impact

- revenue_per_week: 1200
- penalty_per_week: 300
- time_lost_per_cycle: 4h
```

Other lines in the section are evidence and are not read. A `time_lost_per_cycle` there takes precedence over the one worked out from cost and count. A value that is not an amount of zero or more, or not a duration longer than zero, is left out. The story's Notes say how each input was set and which were left out.

The issue's Remediation section then gets a line naming the story by ID: `Story S-0009 remediates this issue, created from it at 2026-10-03T12:00:00Z.` It names the ID, not the file, because the story's file moves to `wip/archive/` when it is accepted.

`flai check` runs with the story in place; if it reports anything the story introduces, the story is removed and the findings are printed (exit 4). A closed issue, or one an open story already links, is refused, and the refusal names that story. `--story S-0008` reads the issue from that story's worktree, as `list --story` does, so a story in review can turn the issues it recorded into stories; the new story is made in `wip/` as always, and the issue is edited in that worktree for that story to commit. `--epic` puts the story under an epic, and `--owner` names its owner in place of your configured author. `--autocommit` commits the story, its epic, and the issue it names on their own, with each `--trailer` line, unless the project sets `dashboard.autocommit: false`; nothing is pushed. `--json` prints `id`, `title`, `nature`, `path`, `issue`, `draft`, `cost_of_delay` when it has one, and `committed`, with `commit` or `commit_error` when there is one. The MCP `issue_story` tool makes the same story and leaves it, and the issue, uncommitted. The dashboard's review page makes these stories with `--autocommit` when you accept a story with issues checked.

## Accept and release

```bash
flai release S-0031 --dry-run          # what a release would look like now
flai accept S-0031 --by alex           # move to done, archive, commit — no release, no tag, no push
flai accept E-0002 --by alex           # an epic whose stories are all done: the same
git fetch                              # flai never fetches by itself
flai release --pending                 # publish: tag whatever has accumulated and push it
```

Acceptance is one step, and for a story it is the only way to reach done: `flai move S-0031 done` from review, a card dropped on done in the dashboard, and `flai accept S-0031` all run the same flow with the same flags. It rebases the story branch and fast-forwards it into the main branch, moves the item to done (the same rules as `flai move`), archives it with its children and narrative, and commits. It computes no release, creates no tag, and pushes nothing: that is a deliberate step of its own, not tied to any one item, publishing, below. `--dry-run` prints anything that would block acceptance and any uncommitted files outside `wip/`, and stops without refusing; `--trailer` appends lines such as co-author attribution to the commit message. The working tree must be clean outside `wip/` so the acceptance commit holds only acceptance, unless you pass `--yes`, which includes those files in it. From the dashboard the same choice is a checkbox in the confirmation.

Accepting an epic's last open story accepts the epic too (S-0200, [the workflow](../../design/system/workflow.md#an-epic-follows-its-stories)): the epic moves to done after the story and is archived with it and with its cancelled stories, in the one acceptance commit, `chore: [S-0031] accept and archive, with E-0002`. `--dry-run` says `would also move E-0002 <title> from review to done, following S-0031, and archive it`, and anything that would stop the epic is a blocker before anything is merged. `--json` has the epic's move in `epic`. An epic accepted this way counts toward the next publish as one you accept yourself.

Acceptance then tells the stories still in progress or in review what it changed under them. For each one whose claim ([Touches](#touches)) covers a path the merge brought into the main branch, shared paths included, it records which paths those are. A story with no touches is told of every path. The command prints `told S-0040 it overlaps: flai/cmd/accept.go`, `--json` lists them in `overlaps`, and the story's agent sees it in its MCP `inbox` as an `overlapped` change. Acceptance from the dashboard does the same.

Acceptance checks what could fail midway before it changes anything: without a git committer identity it refuses and the story stays in review; an experiment story is refused until its results document is committed on its branch, and the refusal names the document to write, `design/experiments/<S-nnnn>-<slug>.md` (ADR-0066). Start it from `design/experiments/template.md`: front matter `title`, `updated`, `status`, and `story`, and the sections Hypothesis, Success measure, What was done, Results, and Recommendation, which says adopt, adapt, or drop; `flai check` validates it. With it, an experiment is accepted like a research story, and publishing gives no component a bump on its account, whatever it touched.

After it rebases the story branch and before it fast-forwards the main branch, acceptance reads every file the branch adds or changes, markdown or not, and refuses, merging nothing, when any carries a merge conflict marker: a line that begins with `<<<<<<<`, `|||||||`, or `>>>>>>>` followed by a space or the line's end, the lines `flai check` reports as `markdown.conflict-marker`. Such a line means a conflict met during a sync was added to git unresolved. The refusal names each one as `path:line`, for example `story/S-0031 carries merge conflict markers at design/system/workflow.md:7, design/system/workflow.md:11`, and the story stays in review with its branch and worktree. Resolve each conflict in the worktree, keeping what both sides meant, remove the markers, commit, and accept again. `flai accept --dry-run`, and the dashboard's Accept dialog that shows it, lists the same markers as a blocker, read by the same check, so you see them before you accept (S-0276).

Accepted work stays local until it is published. Publishing is the one way it reaches the remote ([ADR-0067](../../design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)): `git fetch`, then `flai release --pending`, or **Publish** on the dashboard's board, which runs the same on the host. Several acceptances release together: it computes, applies, tags, and pushes the branch and the tags, never forcing, sending the tags on their own for acceptances already pushed. Agents publish only when you ask them to. `flai board` lists what is accepted and not yet published (`accepted, not yet published: S-0031 (publishing is the operator's: git fetch, then flai release --pending)`), and `flai board --json` and the agents' `inbox` carry the IDs as `unpublished`.

What is pending is worked out from this clone's release tags and branch, and flai never fetches. So before it plans, `flai release --pending` asks the remote for its tags and for the head of the branch this clone tracks. When the remote has a newer `<name>/vX.Y.Z` than this clone, as after publishing from another clone and pulling without tags, or has commits on that branch this clone lacks, it plans, commits, and tags nothing, says what is missing, and refuses to publish (exit 3), naming what to run. Run it, then publish again:

```bash
git fetch --tags origin                     # release tags the remote has and this clone lacks
git fetch origin && git merge origin/main   # commits on the remote branch this clone lacks
git merge origin/main                       # the same, already fetched; or rebase onto origin/main instead
```

The remote can still move between that check and the push. Each push goes as one `git push --atomic`, so when the remote refuses the branch it refuses that push's tags too. flai then asks the remote again, and when its branch has moved it deletes here the release tags that did not reach the remote, because they no longer tag what will be published, names them, and exits 3. Tags an earlier push in the same publish already sent (more than three tags go in batches of three, the branch with the last) are kept, here and on the remote, and named; then merge rather than rebase, so the commits they tag stay in the history. `--json` carries `remote_moved`, `deleted_tags`, and `kept_tags`. The publish commit that bumped the template's version stays: a rebase replays it and a merge keeps it, and the rerun does not bump again. To recover:

```bash
git fetch origin
git rebase origin/main                      # or git merge origin/main; merge when tags were kept
flai check --strict                         # and your tests: verify what you now have
flai release --pending                      # tags again and pushes
```

When the rebase or merge conflicts, the conflicts are worked through threads: one with you on each accepted story whose changes conflict, naming the paths, resolved as the thread settles. A push that fails for any other reason (credentials, a hook, an unreachable remote) keeps its tags, and running `flai release --pending` again finishes the push.

When the remote cannot be reached, `--dry-run` still shows the plan with a warning that it was not checked, and publishing waits until the remote can be reached, since it pushes there anyway. A clone with no remote publishes locally as before.

`flai push --pending` and the `auto-publish` host action are kept as the operator's own shell tools, outside the workflow: no agent is told to run them, and the dashboard neither offers a push nor shows `auto-publish`. `flai push --pending` pushes the branch and any tags already made when the commits ahead include an acceptance or a release tag; ordinary commits are yours to push with git. It never forces, and refuses when the remote has commits this clone lacks. It releases nothing unless `auto-publish` is on for the project (`flai serve enable auto-publish`, off by default, S-0144): then it first tags whatever has accumulated, as `flai release --pending` would, and refuses as it does.

```bash
flai push --pending             # push the branch and any tags already made; with auto-publish on, tag the pending release first
flai push --pending --dry-run   # say what would be pushed, and tagged when auto-publish is on
flai push --pending --publish   # also publish the template when those commits moved its version
```

An accepted item that no plan can cover is named with the reason (`left out:`) instead of being skipped silently: one touching two components with no tag saying which it delivers to, for example. Tag it, or its epic, and it is planned next time.

A story that is `done` but was never accepted (an older flai, a hand edit) is flagged by `flai check` as `story.unaccepted`, and `flai accept` completes it.

### When the orchestrator accepts

With `orchestration.permissions.accept_reviews` on ([Running the orchestrator](#running-the-orchestrator)), the orchestrator accepts stories in review itself ([ADR-0093](../../design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)). It runs the same flow as you: merge, done, archive, commit. It never publishes. For each story in review it has its verifier run the tests, the lint, and `flai check --strict` in the story's worktree and match each acceptance criterion to the changed files that meet it. It then previews the acceptance and accepts only when nothing blocks it. Otherwise it leaves the story in review and opens a thread on it saying what is missing.

```bash
flai accept S-0031 --by orchestrator --verified 4f1c2a9 --dry-run
flai accept S-0031 --by orchestrator --verified 4f1c2a9 --evidence evidence.md   # --evidence - reads standard input
```

`--verified` is the commit the verifier passed, which must be the story branch's head. `--evidence` is required, except with `--dry-run`. It is markdown without headings: one `Verdict:` line from the verifier's report, and one item per acceptance criterion naming the changed files that meet it. Quote the files in backticks to add a note after them; otherwise separate them with commas.

```markdown
Verdict: pass, tests and lint clean
- 1: `flai/cmd/accept.go`, `flai/internal/preview/accept.go`
- 2: `docs/users/flai.md` (the new section)
```

The preview adds these blockers to an acceptance by the orchestrator, and `flai accept` refuses on any of them before it merges anything. `--json` lists them in `orchestrator_blockers` as `{code, message}`, and each is in `blockers` too.

| Code | Blocks when |
|------|-------------|
| `unverified` | `--verified` is missing, names no commit, or is not the story branch's head. A commit added after the verifier's run blocks until it verifies again |
| `criterion_unticked` | An acceptance criterion is unticked |
| `outside_touches` | The branch changes a file under none of the story's touches. The `wip` folder, `design/issues`, and an experiment's results never count |
| `thread_open` | A thread on the story or one of its tasks is not resolved |
| `criterion_unevidenced` | The evidence has no item for a criterion, or its item names no file the branch changes |
| `not_story` | The item is an epic, which is yours to accept |

`--by orchestrator` is refused while `accept_reviews` is off, and the refusal names it. `--verified` and `--evidence` are refused with any other `--by`. In the orchestrator's session (`FLAI_ROLE=orchestrate`), an acceptance with any other `--by`, or none, is refused by flai and by `flai guard`. The MCP tool `item_move` still refuses it a move to done, and names `flai accept --by orchestrator`. `flai move S-0031 done` takes the same flags as `flai accept`.

What it leaves for you to read:

| What | Where |
|------|-------|
| Who accepted | The story's done transition, `by: orchestrator` |
| The commit verified, when, and the evidence | The story's `## Notes`, under `### Accepted by the orchestrator`, with `- Verified: <commit>` and `- At: <time>`, in the acceptance commit and so in the archive |
| The same, for a script | `--json`: `by`, `verified`, and `evidence` (`verdict`, `criteria` as `{n, files}`, `text`). The text output adds a line naming the commit and where the evidence is |
| The decision | `wip/agents/orchestrator.md`, under `## Log` |
| On the dashboard | The story's page says who accepted it and links the evidence ([flaiover.md](flaiover.md#reviewing-a-story)) |

[Accepting a story](../../design/system/strategic-agents.md#accepting-a-story-s-0221) has the whole of it.

The release follows the git convention. Components are the `projects` in `system-flow.yaml`. The component the item delivers to gets the delivery-type bump: feature story minor, remediation or improvement patch. An epic, accepted with its last story or by hand, gets none of its own: its stories carry theirs ([ADR-0078](../../design/adrs/0078-an-accepted-epic-contributes-no-release-bump-of-its-own-its-stories-carry-theirs.md)), and `flai release <epic>` says so. It is found from the item's tags (a project name or one of its `tags` aliases), then its epic's tags, and only among the components the item's commits touched: when the tags name several, the one with the most touched files delivers, the earlier tag breaking a tie, and a tag naming a component no commit touched never delivers, so that component gets no release at all. `--deliver` overrides all of that. With no tag deciding, the only touched component delivers, and several ask you for a tag or `--deliver`. `flai check` warns about that earlier (`story.component-tag`): an open story whose touches reach two or more components while no tag of its own or its epic's names one of them, with the `flai edit --tag` that fixes it. Every other component the item's commits touched gets a patch. Code components get an annotated tag `<name>/vX.Y.Z`; a `template` component gets its `template.yaml` version and `CHANGELOG.md` bumped instead. An item whose commits touch no component, such as design or docs work, releases nothing. A research or experiment story releases nothing either, whatever it touched: its findings or results are merged like any acceptance, and if its commits changed a component's files the plan says that component lands on main without a release. Publishing pushes the acceptance commit whether or not a release was cut.

### Whether a release is due

```bash
flai release --evaluate          # is the release policy met, and on what figures
flai release --evaluate --json
```

`flai release --evaluate` says whether the project's release policy, `orchestration.release` in `system-flow.yaml` ([settings](../operators/settings.md)), is met, why, and the figures it rests on. Like `flai order --by` ([Ordering by a policy](#ordering-by-a-policy)), it is the orchestrator's arithmetic, kept in flai so that you and the dashboard get the same answer with no agent running. It weighs the stories accepted and not yet released, the ones `flai release --pending --dry-run` would publish: how many there are, and their cost of delay values per week added up.

| Policy | Met when |
|--------|----------|
| `threshold` | The summed value of the stories waiting is at or over the policy's `value`, or their count is at or over its `count`. A story with no value counts toward the count, adds nothing to the sum, and is named. With nothing waiting it is not met. |
| `theme` | Every story of the policy's `epic`, or of its `tag`, cancelled ones aside, is accepted, and at least one of them is not yet released. The ones not yet accepted are named. |
| `judgement` | Never, by itself: whether to release is the orchestrator's call, or yours. This is the default. |

It tags, bumps, commits, and pushes nothing; publishing is still `flai release --pending`. So it takes no item and refuses `--apply`, `--pending`, `--dry-run`, and `--deliver`. The dashboard reads the same answer through `flai serve` (`release.evaluate`), and agents through the MCP tool `release_evaluate`. A sub-agent and the planner may run it as a read.

With `orchestration.release.whole_epics` set, no policy is met while a story waiting belongs to an epic in neither review nor done. Those stories are listed under `held by epic:` with their epics and the epics' status, and in `held_by_epic` with `--json`. Under `judgement` they are what the orchestrator must not publish. A story with no epic is never held.

### When the orchestrator publishes

With `orchestration.permissions.publish` on ([Running the orchestrator](#running-the-orchestrator)), the orchestrator publishes for you, by the release policy ([ADR-0094](../../design/adrs/0094-with-publish-on-the-orchestrator-publishes-only-through-release-publish-by-its.md)). Turning the permission on is your asking: agents still publish only when you ask. After each acceptance, its own or another's, it evaluates the policy. Under `threshold` or `theme` it publishes when the policy is met. Under `judgement` it publishes when it judges the unreleased work coherent and complete, and says why. Under every policy it never publishes a batch `whole_epics` holds back. Each release publishes everything accepted and not yet released, as the board's Publish does.

It publishes only through the MCP tool `release_publish`, which runs the board's `flai release --pending` on the host, never forced. `flai guard` refuses it `flai release` other than `--evaluate`, `flai push`, `git push`, and `git tag`, whatever its permissions. `release_publish` refuses, changing nothing:

| When | It says |
|------|---------|
| Nothing is accepted and not yet released | That there is nothing to publish |
| `whole_epics` holds the batch back | Each story held, its epic, and the epic's status |
| Under `threshold` or `theme`, the policy is not met | The policy and its figures |
| Under `judgement`, it gives no reason | That the call needs its reason |
| The `push` host action is off | That `publish.run` needs it, and `flai serve enable push` |
| The remote has release tags this clone lacks, or its branch moved | flai's own message, beginning `conflict:`, with what to run |

What it leaves for you to read:

| What | Where |
|------|-------|
| The run | `flai serve journal`: method `mcp.release_publish`, by `orchestrator`, with the policy and its reason before the outcome |
| Each release | `wip/agents/orchestrator.md`, under `## Log`: the policy, its figures, the versions and tags, and the items bundled |
| Each decision not to publish | The same log, with the evaluation's figures |
| Each refusal | The same log, and a thread to you on the most recently accepted story of the batch, with the refusal and what fixes it. It does not try again until you answer the thread or the next story is accepted |

[Publishing a release](../../design/system/strategic-agents.md#publishing-a-release-s-0222) has the whole of it.

## Upgrade to a newer template

```bash
flai upgrade --dry-run      # what would change
flai upgrade                # interactive: keep, replace, or diff each conflict
flai upgrade --keep-all     # scripts and CI: never overwrite a project edit
flai upgrade --relock       # a project assembled by hand: record the current files at this version
flai upgrade --var team=billing   # a variable the template added, or a new value for one
```

`flai new` writes `system-flow.lock.yaml`, a hash of every file the template rendered. On upgrade each template path is classified: **add** when the project lacks it, **merge** for files with the baseline marker such as `CLAUDE.md` and the conventions (template text above the marker, yours below; `topics` you set on a convention's front matter stay yours, and a convention whose topics you left as the template gave them takes the new template's), **replace** when your copy still matches the lock, **unchanged** when identical, otherwise a **conflict** that you decide. Without a lock every difference is a conflict, which is what `--relock` fixes. A dirty git tree is refused unless `--force`, so an upgrade is one reviewable diff, and the manifest's template version is updated only when no conflict is left undecided. Kept conflicts stay divergent and come back next time; replace them or add your rule below a marker instead.

The lock also records the value of every template variable, and upgrade renders with them, so a variable your template's fork added keeps its value. The project's name, key, description, owner, and repository URL come from `system-flow.yaml`; edit them there. A variable the new template version adds takes its default, which the upgrade names, or the value you give with `--var`. `--var` also changes a recorded value, and on a project already at the template's version it re-applies that version. A required variable with no value is named and nothing changes: run the upgrade again with `--var name=value`. A project made by an older flai has no recorded values, so its first upgrade names the variables that took their default; pass `--var` for any whose value it was made with differs.

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
flai dashboard status          # running, not answering, or gone, from a probe of /_health
flai dashboard logs [-f]
flai dashboard stop
flai dashboard token           # print the token and login link
flai dashboard token --rotate  # new token; a running dashboard restarts
flai dashboard --no-serve      # do not register with flai serve, or start flai host
```

### flai host: the process that runs flai serve and the MCP servers

`flai dashboard` also starts `flai host`, one process of yours per machine. It runs `flai serve` and each served project's MCP server over HTTP as its children. It starts them again if they end, and stops them all when it stops ([ADR-0040](../../design/adrs/0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)). It also restarts the dashboard container when it is gone or stops answering, from the image it was running, with a back-off. `flai config set dashboard.no_restart true` turns that off ([The dashboard's watch](../operators/index.md#the-dashboards-watch)). Its state, token, and log are in a folder named `host` beside flai's config file.

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

`flai serve` does what a dashboard asks only among the methods flai offers, and what touches your credentials is off until you turn it on. These *host actions* are yours to enable, by name, in a shell on the host; `push` lets the board publish when you press Publish:

```bash
flai serve actions          # what there is, what each means, and where each is on
flai serve enable push      # for this project; --all-projects for every project
flai serve disable push
flai serve journal          # every host action asked for, and what became of it
```

`auto-publish`, off by default, makes `flai push --pending`, in any shell that uses this configuration, release first. It is the operator's shell tool, outside the workflow (ADR-0067): `flai serve actions` lists it, and the dashboard neither shows nor changes it. Off, acceptances wait, unreleased, to be published together from Publish or `flai release --pending` (S-0144).

A second one, `agent`, starts an agent for each story that becomes ready: the harness and model the story names (see [Who works a story](#who-works-a-story-its-agent)), with the program and permissions you set for that harness on the host (`flai serve agent harness`), or a command you wrote for stories that name none (`flai serve agent set -- <program> [args...]`). Turn it on with `flai serve enable agent`. A story in ready or in progress whose agent dropped or failed gets a new one with `flai serve agent restart <story>`, or **Retry** on its page. So does a story that your flai serve has started no agent for, such as one begun on another host: the board shows it as begun elsewhere with **Start agent** on its page, and its agent is told who began it, where, and when, takes up its branch with `flai stream open` (see [Story branches](#story-branches)), and goes on from what was committed and pushed. For a story in ready while the in-progress limit or review is full, the new agent is queued and starts as soon as there is room. The new agent is told how the last one ended and goes on from the story's narrative. flai serve does not start a story that is held (see [Touches](#touches)), and says why on the board; a restart of one is queued until it is clear. A story in ready gets its agent at once, whatever flai serve's own rules say, held, past the limit, or past a full review, with `flai serve agent start <story>`, or **Start agent** on its page. A story in review whose agent left changes uncommitted in its worktree, which blocks its acceptance, gets an agent to commit them, and do nothing else, with `flai serve agent commit <story>`, or **Have an agent commit them** in the acceptance confirmation or on its review page. A story's agent that runs, or waits for an answer, is stopped with `flai serve agent stop <story>`, or **Stop** on the dashboard's activity page: its process and everything it started end, the story stays where it is with its worktree as the agent left it, and it gets no agent until it is retried or moved back to ready. What a story's agent is saying and doing, read from the log flai serve gives it, is printed by `flai serve agent stream <story>` (`--follow` until it ends), and shown on each card of the dashboard's activity page. From the same logs flai serve measures what the story's agents spent, when an agent ends and when a task of a story whose agent runs enters done, and writes it as the story's and its tasks' `usage` (see [Tokens and cost](#tokens-and-cost)). `flai serve agent usage <story>` prints what the logs say, `--write` records it, and `--all --write` fills in every story flai serve kept a log for, such as those worked before your flai measured them. Run it with the configuration of the flai serve that started the agents (the installed flai, or `--config`), since the logs are in its serve folder. Only Claude Code's logs are read. A third, `host`, lets the dashboard have `flai host` restart `flai serve` or the MCP servers, or upgrade flai and restart on it. A fourth, `settings`, lets the dashboard's Settings page change all of this for you ([the operator guide](../operators/index.md#the-settings-host-action-changing-the-hosts-settings-from-the-dashboard-s-0105) says what that gives the dashboard token). A fifth, `plan`, lets you, the dashboard's **Plan**, and your own agent start the planner for an epic or a story ([Running the planner](#running-the-planner)). A sixth, `orchestrate`, runs the orchestrator for the project for as long as it is on ([Running the orchestrator](#running-the-orchestrator)). `auto-approve`, shell only like `auto-publish`, lets a story's agent write its `.claude/` files without asking you ([Writes under .claude/](#writes-under-claude)). What enabling each means, for who can publish a release and who can start a process on your machine, is in the operator guide; read it first.

The dashboard needs its login token for everything but health and readiness: one token per user, for every project the dashboard serves. `flai dashboard` creates it on first run, as `dashboard.token` in the `serve` folder beside flai's config file, and prints a login link; open the link (or paste the token on the login page) and the browser keeps a session cookie. Tools send it as `Authorization: Bearer`. Details and the exposure table are in the operator guide.

The container runs detached as `flaiover`, one for every project on the host, published on every interface of the host (`--bind`, or `dashboard.bind`, restricts it) on the configured port. Nothing of the project is mounted into it: it is given its port, the login token, and a credential for the host's flai, and everything it shows and changes it asks of `flai serve` on your machine, which `flai dashboard` has `flai host` start alongside it. Commits and acceptances made from the board are therefore made on your machine, as you, wherever the repository lies. Image, tag, port, and bind address come from flags, then the `dashboard` section of `system-flow.yaml`, then `~/.flai/config.json`. If Docker is not installed the command says so with an install pointer. The container holds no git credential: what is accepted from the board is published from the host, by `flai serve` when you press Publish with the `push` host action enabled, or with `flai release --pending` in a shell. The old `--push-key` flag and `dashboard.push_key` setting are retired and ignored.

The image lives on GHCR and is private while the repository is. When the pull is refused, `flai dashboard` logs Docker into the registry with `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` and retries once. The token needs the `read:packages` scope; `gh auth refresh -h github.com -s read:packages` adds it. Inside the monorepo, `--build` sidesteps the registry by building the image from `flaiover/` as `flaiover:local`.

The image is published for `linux/amd64` and `linux/arm64`. `flai dashboard`, `flai dashboard upgrade`, and `flai dashboard check` pull the image for the platform the Docker daemon runs, such as `linux/arm64` on Apple silicon. If the image already present was built for another platform (for example an amd64 image pulled before arm64 was published), `flai dashboard` pulls again without being asked. A tag published without your platform fails with an error that names the platform; pick a newer tag, or use `--build` in the monorepo.

#### Writes under .claude/

Claude Code treats a file in a `.claude/` folder, the project's settings, hooks, and agent definitions and the template's copies of them, as sensitive: run headless, it writes one only when a person or a permission handler approves. `flai serve` starts every Claude Code session, a story's agent's and the planner's, with flai's MCP tool `permission_prompt` as that handler (`--permission-prompt-tool mcp__flai__permission_prompt`, ahead of your harness arguments; S-0257, [ADR-0086](../../design/adrs/0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md)). It approves one thing: an `Edit`, `Write`, `MultiEdit`, or `NotebookEdit` of a file in a `.claude/` folder inside the worktree of a story in progress. Anything else it is asked is refused, as it was before, so harness arguments that narrow an agent's permissions keep their effect.

A write it refuses at once reaches the agent with flai's reason, such as "/tmp/notes.md is not inside a story's worktree", "S-0042 is ready, not in progress", or "… is not in a .claude/ folder of S-0042's worktree", followed by what it approves; no thread is opened. The agent then writes in its story's worktree instead, or moves its story to in progress first. A flai released before S-0283 answers in a shape Claude Code reads as "Permission prompt tool returned an invalid result", so with it an allow never lets a write through and a refusal never says why (I-0082); install a release with S-0283 for the agents' MCP server to answer as described here.

For such a write it opens a thread on the story, such as "Allow Write .claude/settings.json?", that shows the whole content written or the old and new text of each edit, and the agent waits. Reply as the story's owner or as the project's owner, the `owner` in `system-flow.yaml`; the thread names who may answer (S-0284, [ADR-0097](../../design/adrs/0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md)). A reply whose first word is `allow`, `yes`, `approve`, `approved`, or `ok` lets the write through, and any other reply refuses it with your words as the reason. The asking agent's own replies never answer, and replies by anyone else, other agents included, are not answers. When neither owner is set, any reply but the agent's answers. A flai released before S-0284 takes an answer from the story's owner alone, so a reply under another name is never seen and the write waits until Claude Code gives up on it (I-0081). The thread is resolved with what was decided, and the write is refused if the session ends first.

`auto-approve`, off by default, lets these writes through at once, with no thread, and logs each. Turn it on in a shell on the host; it applies to agents already running:

```bash
flai serve enable auto-approve    # for this project; --all-projects for every project
flai serve disable auto-approve
```

With it on, an agent can change its own settings and hooks unattended. Like `auto-publish` it is your shell tool (ADR-0067): `flai serve actions` lists it, and no dashboard shows or changes it, so the dashboard token cannot give agents that right. Only the `flai mcp` an agent starts on stdio reads it; the shared HTTP server (`flai mcp start`) always asks. Other harnesses get no handler.
