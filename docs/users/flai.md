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

Point `template.repo` at any fork and `template.ref` at any branch, tag, or commit to use your own template. A local directory path also works, which is how the system-flow repo develops against its own `./template`.

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
