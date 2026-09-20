---
title: flai CLI
updated: 2026-09-20
status: draft
---

# flai

The system-flow command line tool. Full command reference will be generated from the CLI in story S-0017; the design is in [design/system/flai-cli.md](../../design/system/flai-cli.md). Implemented so far: `version` and `config`.

## Install

One line on macOS or Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/bytepunx/system-flow/main/install.sh | sh
```

The script detects the OS and architecture, resolves the newest `flai/v*` release, downloads the archive and `checksums.txt`, verifies the SHA-256, and installs `flai` into `/usr/local/bin`, using `sudo` only when that directory is not writable. It ends by checking that the directory is on your `PATH`.

| Variable | Effect |
|----------|--------|
| `FLAI_INSTALL_DIR` | Install somewhere else, for example `$HOME/.local/bin` |
| `FLAI_VERSION` | Pin a release, for example `1.0.3` |
| `GITHUB_TOKEN` or `GH_TOKEN` | Authenticate against the GitHub API. While the repository is private one of these, or a `gh auth login` session, is required; the script borrows `gh auth token` when it can |

### Upgrade

```bash
flai self-upgrade --check      # installed and latest versions
flai self-upgrade              # replace this binary with the latest release
flai self-upgrade --version 1.0.3
flai self-upgrade --dir /usr/local/bin
```

`self-upgrade` performs the same steps as the script from inside the binary: resolve, download, verify, and replace the running executable atomically. It uses the same token sources. Without `--version` it does nothing when the installed version is already the latest.

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

`flai` keeps its settings in `~/.flai/config.json`. The first command that needs it creates the file with these defaults:

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
    "bind": "0.0.0.0"
  },
  "cache_dir": "~/.flai/cache",
  "author": "<your username>"
}
```

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

Point `template.repo` at any fork and `template.ref` at any branch, tag, or commit to use your own template. A local directory path also works, which is how the system-flow repo develops against its own `./template`. Git templates are cloned under `cache_dir`; set `FLAI_CACHE_DIR` before the first run to choose where that default lands (for example inside a repository or a CI workspace).

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
```

`import` scans the tree and proposes: the three documentation folders (reusing `docs/`, `design/`, or `wip/` if they exist, or names you choose with `--layout`), whole-folder moves for `adr/`, `adrs/`, `architecture/`, `doc/`, and `documentation/`, a list of loose markdown files to place, and the code sub-projects it found by their build files (`go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`). Applying it creates the structure, renders every template file that does not already exist, performs the moves with `git mv` when the file is tracked, writes `system-flow.yaml` with the sub-projects, and runs `flai check`. Existing files are never overwritten; a conflicting move is reported and the source left in place. A repository that already has `system-flow.yaml` is refused unless `--force`.

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

Items are created from the template's item bodies with the next free ID and linked into their parent's Stories or Tasks list. Natures: `feature`, `improvement`, `remediation`, `research`, `experiment`.

To create an item with its body already written, give the body on standard input. This is what the dashboard's "new" form does, and it is one step that happens or does not:

```bash
flai story new --print-body > story.md          # the sections your project's template gives a story
$EDITOR story.md                                # write the goal, criteria as - [ ] lines, notes
flai story new "Invoice PDF export" --epic E-0001 --body-stdin --autocommit < story.md
```

The heading (`# S-0007 Invoice PDF export`) and the front matter are flai's; the body is everything below the heading. `flai check` runs with the new item in place: if it reports anything the item introduces, nothing is created, the parent is left as it was, the findings are printed, and the exit code is 4. `--autocommit` commits the new item and its parent on their own (`chore: [S-0007] create story: ...`) unless the project sets `dashboard.autocommit: false`; `--trailer` adds trailer lines. Nothing is pushed.

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

`worktrees.relative_paths` is off by default, and flai never turns it on for you, whatever git you have. Set it when the repository cannot be mounted in the dashboard at its own path, or when you move the clone around. With it on and git 2.48 or newer, `flai stream open` creates the worktree with `git worktree add --relative-paths`, so the links work wherever the repository is mounted. With it on and an older git, flai warns, naming your git version and the setting, and creates an ordinary worktree.

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

`touches` is an advisory list of paths or components a story or task is changing. `flai check` warns (`wip.overlap`) when two in-progress items cover the same path, the board prints it under each card, and the dashboard shows a "being worked on" notice on those documents.

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

### Serving agents over MCP

```bash
flai mcp          # an MCP server on stdio; agents start it, you do not
```

`flai mcp` gives an agent session a typed, low-latency view of the repository. Register it once per project in `.mcp.json` (new projects get this from the template):

```json
{ "mcpServers": { "flai": { "command": "flai", "args": ["mcp"] } } }
```

| Tool | What it does |
|------|--------------|
| `inbox` | Threads awaiting the agent (`awaiting: you` when the last entry is not the agent's; `story` filters, `all` includes the rest), `ready`: the stories ready to pull, in pull order, with `can_pull` from the in-progress limit, and `changes`: what others did to work items since this agent last looked (moved, blocked, unblocked, pull order changed), each reported once `unpushed`, on every call while it is true: an acceptance made in this clone and not pushed (items, commits ahead, tags), which the agent pushes from the host with `git fetch` and `flai push --pending` |
| `board` | The board as `flai board --json` prints it; `all` adds epics and tasks |
| `thread_get`, `thread_open`, `thread_reply`, `thread_resolve` | Read, start, answer, and close threads as the agent (`FLAI_AGENT`) |
| `item_get`, `item_move` | Read an item with its children; transition it with the workflow rules. Moving a story or epic to done is refused: acceptance is yours |
| `doc_get` | A markdown document under the design, docs, or wip folders; nothing else in the repository is served |
| `who_touches` | In-progress and in-review items whose `touches` cover a path |
| `wait_for_events` | Returns at once when something changed since this agent last looked, otherwise blocks until a thread, item, or narrative changes, or the timeout passes. Returns `events` in the same shape as `changes`, and the changed paths |

An agent on another machine reaches the same server through the dashboard, over HTTP, with the project token. Its `.mcp.json` names the dashboard's `/mcp` and sends the token as a bearer; `X-Flai-Agent` is the name it works under, as `FLAI_AGENT` is locally:

```json
{
  "mcpServers": {
    "flai": {
      "type": "http",
      "url": "https://dashboard.example/mcp",
      "headers": { "Authorization": "Bearer ${FLAIOVER_TOKEN}", "X-Flai-Agent": "claude@laptop" }
    }
  }
}
```

The tools and their behaviour are the same, because the dashboard starts a `flai mcp` of its own for each connected agent and passes messages through. Keep the token out of the file: Claude Code expands `${VAR}` in `.mcp.json`, as above; for another client, check how it takes a secret. Use `https` through a tunnel or proxy whenever the dashboard is not on a network you trust; see the operator guide.

"Since this agent last looked" is a marker per agent name (`FLAI_AGENT`) under `.flai-cache/mcp/`, outside git. It only decides which changes are news; ready work and open threads are listed on every call, so nothing depends on it. One look reports at most 50 changes, the newest, and says in `changes_omitted` (`events_omitted` for `wait_for_events`) how many older ones it left out; those are not reported later. The first look under a new name covers the last 24 hours and tells of stories and epics only, not task transitions: to an agent that has just arrived, a day of task moves is history, and `board` and `item_get` show how things stand. An agent that ends its turn between your messages calls `inbox` when it starts again and hears what you did in between; one that stays running holds `wait_for_events` and hears within a second. `flai move` records `FLAI_AGENT` as who moved an item when it is set, so an agent is not told about its own moves; a move on the dashboard's board is recorded as the project's `owner`.

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
flai release S-0031 --dry-run          # what acceptance would release
flai accept S-0031 --by alex           # move to done, archive, bump, commit, tag, push
flai accept E-0002 --by alex           # an epic: major release of what it delivered
flai accept S-0016 --by alex --no-release
```

Acceptance is one step, and for a story it is the only way to reach done: `flai move S-0031 done` from review, a card dropped on done in the dashboard, and `flai accept S-0031` all run the same flow with the same flags. It rebases the story branch and fast-forwards it into the main branch, moves the item to done (the same rules as `flai move`), archives it with its children and narrative, computes the release, bumps the template's version file and changelog if the template is involved, commits, creates the tags on that commit, and pushes the branch and tags. `--no-push` keeps everything local; `--dry-run` prints the plan, anything that would block it, and any uncommitted files outside `wip/`, and stops without refusing; `--trailer` appends lines such as co-author attribution to the commit message. The working tree must be clean outside `wip/` so the acceptance commit holds only acceptance, unless you pass `--yes`, which includes those files in it. From the dashboard the same choice is a checkbox in the confirmation.

Acceptance checks what could fail midway before it changes anything: without a git committer identity it refuses and the story stays in review. A push that cannot happen, for example from the dashboard container, which has no credentials unless the operator names a push key (`dashboard.push_key`, [operators](../operators/index.md)), does not undo the acceptance: the result says it was accepted locally and not pushed, and `git push origin HEAD --follow-tags` from a shell finishes the job. `flai board` says so for as long as it is true (`accepted, not pushed: S-0031 (3 commit(s) ahead of origin/main, tags flai/v1.3.0)`), `flai board --json` carries it as `unpushed`, and one command on the host finishes the job:

```bash
flai push --pending             # push the branch and the release tags of acceptances that were not pushed
flai push --pending --dry-run   # say what would be pushed
```

It only acts when the commits ahead of the remote include an acceptance; ordinary commits are yours to push with git. It never forces: when the remote has commits this clone lacks it refuses and tells you to fetch and merge first. It answers from what this clone knows, so a push made from another clone is not seen until you fetch.

A story that is `done` but was never accepted (an older flai, a hand edit) is flagged by `flai check` as `story.unaccepted`, and `flai accept` completes it.

The release follows the git convention. Components are the `projects` in `system-flow.yaml`. The component the item delivers to gets the delivery-type bump: epic major, feature story minor, remediation or improvement patch. It is found from the item's tags (a project name or one of its `tags` aliases), then its epic's tags, and only among the components the item's commits touched: when the tags name several, the one with the most touched files delivers, the earlier tag breaking a tie, and a tag naming a component no commit touched never delivers, so that component gets no release at all. `--deliver` overrides all of that. With no tag deciding, the only touched component delivers, and several ask you for a tag or `--deliver`. Every other component the item's commits touched gets a patch. Code components get an annotated tag `<name>/vX.Y.Z`; a `template` component gets its `template.yaml` version and `CHANGELOG.md` bumped instead. An item whose commits touch no component, such as design or docs work, releases nothing. A research story releases nothing either, whatever it touched: its findings are merged, committed, and pushed like any acceptance, and if its commits changed a component's files the plan says that component lands on main without a release. An experiment story is refused and stays on its branch; `flai accept --no-release` lands one deliberately. The acceptance commit is pushed whether or not a release was cut.

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
flai dashboard --no-serve      # do not register with, or start, flai serve
```

### flai serve: flai on the host, for the dashboards

`flai dashboard` also starts `flai serve`, a small process of yours on the host, and registers the project with it. `flai serve` opens a connection to the project's dashboard and keeps it open, and the dashboard asks flai for what it needs over that connection; the dashboard never connects to the host. It is how the dashboard reads everything it shows (the board, work items, threads, documents, decisions, the inbox, activity, and search) and how it hears that a file changed; writes still use the mounted repository and follow in a later release. Without `flai serve` the dashboard says on every page that it has no flai to ask. One `flai serve` serves every project you start a dashboard for.

```bash
flai serve status     # does it run, which projects, which dashboards have it connected
flai serve            # run it in the foreground to watch it; Ctrl-C stops it
flai serve start      # run it in the background (flai dashboard does this for you)
flai serve stop
```

`flai hostapi` shows what the dashboard can ask, and answers one question on the terminal: `flai hostapi` lists the methods, `flai hostapi board.get '{"all":true}'` prints what the board page is given.

It needs no root and no configuration. `flai dashboard stop` takes the project out of it and leaves it running for your other projects; `flai serve stop` ends it. `flai dashboard status` has a `host flai` line: connected and since when, or why not. The dashboard shows the same at the right of its header, and "host flai: not connected" there means `flai serve` is not running or cannot reach the dashboard: `flai serve status` says which. Its list of projects, its state, and its log (`serve.log`) are in a folder named `serve` beside flai's config file, `~/.flai/serve` unless `FLAI_CONFIG` points elsewhere.

The dashboard needs the project's token for everything but health and readiness. `flai dashboard` creates it at `.flai-cache/dashboard.token` on first run and prints a login link; open the link (or paste the token on the login page) and the browser keeps a session cookie. Tools send it as `Authorization: Bearer`. Details and the exposure table are in the operator guide.

The container runs detached as `flaiover-<project>`, published on every interface of the host (`--bind`, or `dashboard.bind`, restricts it) on the configured port, with the repository mounted read-write at the same absolute path it has on your machine, and running as your user so files it writes keep your ownership. The path matters: git links a story's worktree to the repository with absolute paths, and the dashboard can only accept a story with a branch when those paths exist in the container too. When your path cannot be one in a Linux container (a Windows drive path), the repository is mounted at `/project`, `flai dashboard` warns that stories with a branch must be accepted from a shell with `flai accept`, and the way to make the board work is the relative worktree setting below. Image, tag, port, and bind address come from flags, then the `dashboard` section of `system-flow.yaml`, then `~/.flai/config.json`. If Docker is not installed the command says so with an install pointer. By default the container holds no git credential and an acceptance made from the board is pushed from a shell; `flai config set dashboard.push_key <path>` (or `--push-key`) gives it an SSH key to push with, checked before anything starts. Read [Pushing what the board accepts](../operators/index.md#pushing-what-the-board-accepts) first: it changes what the dashboard token is worth.

The image lives on GHCR and is private while the repository is. When the pull is refused, `flai dashboard` logs Docker into the registry with `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` and retries once. The token needs the `read:packages` scope; `gh auth refresh -h github.com -s read:packages` adds it. Inside the monorepo, `--build` sidesteps the registry by building the image from `flaiover/` as `flaiover:local`.
