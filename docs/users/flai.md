---
title: flai CLI
updated: 2026-09-15
status: draft
---

# flai

The system-flow command line tool. Full command reference will be generated from the CLI in story S-017; the design is in [design/system/flai-cli.md](../../design/system/flai-cli.md). Implemented so far: `version` and `config`.

## Install

With Go 1.26 or newer:

```bash
go install github.com/bytepunx/system-flow/flai@latest
```

Or download a binary from the [releases page](https://github.com/bytepunx/system-flow/releases): releases named `flai vX.Y.Z` carry `flai_X.Y.Z_<os>_<arch>.tar.gz` (a zip on Windows) for Linux, macOS, and Windows on amd64 and arm64, plus `checksums.txt`. Unpack and put `flai` on your `PATH`.

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
    "port": 4242
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
flai story new "Invoice PDF export" --epic E-001
flai task new "Render invoice template" --story S-001 --tag pdf
flai show S-001
```

Items are created from the template's item bodies with the next free ID and linked into their parent's Stories or Tasks list. Natures: `feature`, `improvement`, `remediation`, `research`, `experiment`.

### Moving work

```bash
flai move S-001 ready          # needs at least one task and acceptance criteria
flai move S-001 in-progress    # warns if the WIP limit is exceeded
flai move T-001 in-progress
flai move T-001 done           # tasks may skip review
flai move S-001 review
flai move S-001 done --by alex # needs every task closed and every criterion checked
flai move S-001 in-progress --reason "tests missing"     # from review
flai move S-002 cancelled --reason "superseded by S-005"
```

Every move appends to the item's `transitions` with a timestamp and who made it (`--by`, default the config author). Reasons land under the item's Notes.

### Blocking

```bash
flai block T-001 --reason "waiting on API keys"
flai unblock T-001
```

Blocked items keep their column; the interval is recorded so blocked time shows in the charts.

### The board

```bash
flai board          # stories by column with age, nature, blocked flag, WIP counts
flai board --all    # epics and tasks too
flai board --json
```

### Agent narratives

```bash
export FLAI_AGENT=claude FLAI_SESSION=abc123
flai stream open S-001
flai stream log S-001 "T-001 done, starting T-002"
```

`wip/agents/index.md` is regenerated after every open, log, move, and archive.

### Archiving

```bash
flai archive --dry-run
flai archive               # everything done or cancelled that is safe to move
flai archive S-001         # one story with its tasks and narrative
```

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
flai issue bump I-001 --cost 8m --note "reinstalled again in S-008"
flai issue close I-001 --reason "scripts/install-tools.sh pins v2"
flai issue list [--all]
flai issue summary
```

Issues live in `design/issues/`, one file per recurring problem with a class (`defect`, `blocker`, `efficiency`, `impression`), a count, an average cost per occurrence, and one dated instance per occurrence. `bump` increments the count, updates the average, and appends the instance. Every command regenerates `summary.md`, the table of open issues most expensive first, and `flai prime` lists it after the conventions when anything is open. `flai check` validates the files and warns when the summary is stale.

## Accept and release

```bash
flai release S-031 --dry-run          # what acceptance would release
flai accept S-031 --by alex           # move to done, archive, bump, commit, tag, push
flai accept E-002 --by alex           # an epic: major release of what it delivered
flai accept S-016 --by alex --no-release
```

Acceptance is one command. It moves the item to done (the same rules as `flai move`), archives it with its children and narrative, computes the release, bumps the template's version file and changelog if the template is involved, commits, creates the tags on that commit, and pushes the branch and tags. `--no-push` keeps everything local; `--dry-run` prints the plan and stops; `--trailer` appends lines such as co-author attribution to the commit message. The working tree must be clean so the acceptance commit holds only acceptance, unless you pass `--yes`.

The release follows the git convention. Components are the `projects` in `system-flow.yaml`. The component the item delivers to, found from the item's tags (a project name or one of its `tags` aliases), its epic's tags, or `--deliver`, gets the delivery-type bump: epic major, feature story minor, remediation or improvement patch. Every other component the item's commits touched gets a patch. Code components get an annotated tag `<name>/vX.Y.Z`; a `template` component gets its `template.yaml` version and `CHANGELOG.md` bumped instead. An item whose commits touch no component, such as design or docs work, releases nothing. Research and experiment stories are refused.

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
