---
title: flai CLI
updated: 2026-09-15
status: active
---

# flai CLI

`flai` is a single static Go binary. It creates and imports projects, manages work items, validates a repo against the standard, prints metrics, and runs the dashboard. It is designed to be used by humans at a terminal and by agents through non-interactive flags.

Source: `./flai` in this monorepo. Module `github.com/bytepunx/system-flow/flai`.

## Configuration

`~/.flai/config.json`, created on first run with defaults. Overridable with `FLAI_CONFIG` and `--config`.

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
  "author": "alex"
}
```

Template repo cloning goes into `cache_dir/templates/<hash of repo+ref>` and is refreshed with `flai template update`. A local path in `template.repo` is used directly without cloning, which is how this monorepo develops against `./template`.

## Commands

| Command | Purpose |
|---------|---------|
| `flai new <dir>` | Create a new conforming monorepo from the template. Prompts for variables when stdin is a terminal; `--var k=v` sets them, `--defaults` never prompts. `--template` and `--ref` override the config source, `--layout key=name` renames a folder, `--force` overwrites, `--no-git` skips `git init`. |
| `flai import [dir]` | Analyse an existing repo, propose the layout, prompt for folder names, create missing structure, offer to move existing markdown into it, write `system-flow.yaml`. |
| `flai check [dir] [--strict]` | Validate manifest and layout, every item in kanban and archive (front matter, ID and file name, parents and children, state history and sequence, acceptance criteria, sections), narratives and index, board limits and order, design/docs front matter including ADRs, and the conventions folder (front matter, unique order, marker and project additions, length, README index), and design/issues (schema, unique IDs, file names, summary freshness). One `path:line: level: rule: message` per finding. Errors exit 1; `--strict` makes warnings exit 1. |
| `flai epic new`, `flai story new --epic E-0001`, `flai task new --story S-0004` | Create an item from the body template, allocate the next ID, link to parent. |
| `flai move <id> <state>` | Transition an item with rule validation. `--reason` required for `cancelled` and for review to in-progress; `--by` defaults to the config author. Warns on WIP limit breaches. |
| `flai show <id>` | One item with history, blocks, and children. |
| `flai block <id> --reason`, `flai unblock <id>` | Open and close blocked intervals. |
| `flai board [--all]` | Stories per column with nature, age in column, blocked flag, WIP counts, and the pull order; `--all` adds epics and tasks. |
| `flai stats [--since 30d] [--type story] [--by nature\|type\|parent] [--json]` | The aggregates in [metrics.md](metrics.md) as a table; `--json` adds per-item values, weekly throughput, burn-up (whole set and per parent), cumulative flow, and aging. |
| `flai stream open <story-id>`, `flai stream log <story-id> "<entry>"` | Create a narrative from the template, append a timestamped log entry. `FLAI_AGENT` and `FLAI_SESSION` identify the writer. `index.md` is regenerated from the active narratives after every open, log, move, and archive. |
| `flai migrate ids [--dry-run]` | Widen every work item and issue ID to four digits: rename items and narratives in kanban and archive, and issues under `design/issues`, with `git mv` (plain rename outside git) and rewrite references in the layout folders, the root markdown and yaml files, and each project's root markdown files. Idempotent. |
| `flai archive [id...] [--dry-run]` | Move done and cancelled items and their narratives to `wip/archive`. Default: every closed epic whose stories are archived, every closed story with its tasks, and closed tasks whose story is gone from the board. |
| `flai dashboard [--image] [--tag] [--port] [--pull] [--attach] [--open]`, `flai dashboard stop\|status\|logs` | Pull the flaiover image if missing and run it detached as `flaiover-<project>`, bound to localhost, repo mounted read-write at `/project`, as the host user; status and the already-running message inspect the real container. Precedence: flags, manifest `dashboard`, config. |
| `flai upgrade [--dry-run] [--force] [--keep-all\|--replace-all] [--relock]` | Bring the project to the template version at the configured source: add new files, merge marker files above the marker, replace files unchanged since applied (per `system-flow.lock.yaml`, ADR-0015), report the rest as conflicts with keep, replace, or diff; refuses a dirty tree; `--relock` records a hand-assembled project. |
| `flai prime [--cat] [--json]` | Print `design/conventions` in read order, README first, as paths or contents, so an agent or hook loads the norms in one call; the open-issues summary follows when non-empty. |
| `flai release <id> [--dry-run] [--deliver] [--apply]` | Compute the semver release for an item per the git convention (delivered component gets the delivery-type bump, touched components a patch) and with `--apply` create the tags and version bumps. |
| `flai accept <id> [--by] [--deliver] [--no-release] [--no-push] [--trailer]` | The operator's acceptance in one step: move to done, archive, template bump, commit, tags on that commit, push. |
| `flai issue new "<title>" --class <c> [--cost] [--note]`, `flai issue bump <id> [--cost] [--note]`, `flai issue close <id> --reason`, `flai issue list [--all]`, `flai issue summary` | Record, increment, and close recurring friction in `design/issues`; every command regenerates `summary.md` (average and total cost, most expensive first). |
| `flai template show`, `flai template update`, `flai template use <repo> [--ref]` | Inspect, refresh, and switch the template source. |
| `flai template push [dir] [--remote] [--ref] [--tag] [--dry-run] [--force]` | Publish a locally developed template to its git remote: clone the branch, replace contents, commit with the version, push, and with `--tag` push `v<version>` (refused if it exists). Defaults from `publish` in `template.yaml`. Git errors surface verbatim; no retries; no force push unless `--force`. `flai accept` calls it for template components. |
| `flai config get [key]`, `flai config set <key> <value>`, `flai config path` | Read and edit `~/.flai/config.json` by dotted key; `path` prints the resolved file without creating it. |
| `flai version` | Version, commit, build date. |

Global flags: `--config <path>`, `--json` where output is structured, `--yes` to skip confirmations, and `--verbose` for debug log events, so agents can drive it. stdout is command output; stderr carries structured log events per `design/conventions/logging.md` (`log/slog`, text on a terminal, JSON otherwise, `LOG_LEVEL` and `LOG_FORMAT`), and a failing command logs exactly one fatal event at the boundary.

## Import flow

`flai import` is the deliberate, interactive path for existing repos.

1. Scan the tree. Detect existing `design`, `docs`, `wip`, `adr`, `adrs`, `architecture`, `documentation` folders and any markdown outside code folders. Detect sub-projects by build files (`go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`).
2. Present a proposal: which folders will be created, which existing folders look like they map to `design`, `docs`, or `wip`, which candidate folders (`adr`, `adrs`, `architecture`, `doc`, `documentation`) move whole into the structure, which loose markdown needs a decision, and which sub-projects were found.
3. Ask for folder names, offering existing folders and the template defaults; `--layout key=name` sets them without asking.
4. Move candidate folders (confirmed one by one in a terminal), then render the template files that do not already exist. Nothing is overwritten; a conflicting file in a moved folder stays where it was and is reported.
5. For each loose markdown file, ask: `design/system`, `design/adrs`, `docs/<audience>`, leave in place, or skip all remaining. Non-interactive runs leave files in place. Moves are `git mv` when the file is tracked.
6. Write `system-flow.yaml` with the detected sub-projects, run `flai check`, print a summary and next steps.

`--dry-run` prints the proposal and stops (`--json` gives the analysis and plan). `--yes` accepts every default. A tree that already has a manifest is refused unless `--force`.

## Internal structure

```text
flai/
├── main.go
├── cmd/                 # one file per cobra command
├── internal/
│   ├── config/          # ~/.flai/config.json
│   ├── manifest/        # system-flow.yaml
│   ├── template/        # clone, cache, render
│   ├── workitem/        # parse, validate, transition, ID allocation, board, narratives, archive
│   ├── execx/           # git and docker behind a Runner interface
│   ├── logx/            # slog setup per the logging convention: levels incl. fatal, env and flag
│   ├── narrative/       # wip/agents files
│   ├── check/           # reference validator, rule names are stable identifiers
│   ├── conventions/     # loads and validates design/conventions
│   ├── issues/          # design/issues: record, bump, close, summary
│   ├── release/         # semver plan from item, commits, and manifest components; tags and version files
│   ├── lock/            # system-flow.lock.yaml
│   ├── upgrade/         # classify add, merge, replace, conflict; apply with a policy
│   ├── publish/         # push a local template to its remote: clone, replace, commit, tag, push
│   ├── metrics/         # reference implementation of metrics.md
│   ├── importer/        # scan, plan, moves with git mv
│   ├── dashboard/       # docker run/stop
│   └── ui/              # prompts and tables
└── testdata/
```

External processes: `git` and `docker` are invoked as subprocesses and must be on `PATH`. See ADR 0010.

Front matter is parsed with goccy/go-yaml but written by a small purpose-built emitter, so timestamps stay unquoted, sequences stay indented, and a `flai` edit never rewrites lines it did not change. Item bodies come from the template's `items/` when the project's template is in the cache, else from copies embedded in the binary.
