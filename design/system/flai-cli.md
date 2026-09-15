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
| `flai check` | Validate manifest, front matter, state and transition consistency, parent and child rules, WIP limits, narrative presence. Exit non-zero on errors; warnings with `--strict` too. |
| `flai epic new`, `flai story new --epic E-001`, `flai task new --story S-004` | Create an item from the body template, allocate the next ID, link to parent. |
| `flai move <id> <state>` | Transition an item with rule validation. `--reason` required for `cancelled` and for review to in-progress. |
| `flai block <id> --reason`, `flai unblock <id>` | Open and close blocked intervals. |
| `flai board` | Print the board as a table, with ages and blocked flags. |
| `flai stats [--json] [--since 30d] [--by nature]` | The aggregates in [metrics.md](metrics.md). |
| `flai stream open <story-id>`, `flai stream log <story-id> "<entry>"` | Create a narrative from the template, append a timestamped log entry, update `index.md`. |
| `flai archive [id...]` | Move done and cancelled items and their narratives to `wip/archive`. Default: all eligible. |
| `flai dashboard [--port] [--pull] [--detach]` | Pull the flaiover image if missing, run it with the repo mounted read-write at `/project`, open the browser. `flai dashboard stop`. |
| `flai upgrade [--dry-run] [--force] [--keep-all\|--replace-all]` | Re-integrate the latest template into an existing repo: add new files, replace files unchanged since they were applied, report project-modified files as conflicts, merge `CLAUDE.md` above its marker. Story S-020. |
| `flai template show`, `flai template update`, `flai template use <repo> [--ref]` | Inspect, refresh, and switch the template source. |
| `flai template push [dir] [--remote] [--ref] [--tag] [--dry-run]` | Publish a locally developed template to its git remote: clone, replace contents, commit with the version, push. Git errors surface verbatim; no force push unless `--force`. Story S-021. |
| `flai config get [key]`, `flai config set <key> <value>`, `flai config path` | Read and edit `~/.flai/config.json` by dotted key; `path` prints the resolved file without creating it. |
| `flai version` | Version, commit, build date. |

Global flags: `--config <path>`, `--json` where output is structured, and `--yes` to skip confirmations, so agents can drive it. The first command that needs the config creates it with defaults and says so on stderr.

## Import flow

`flai import` is the deliberate, interactive path for existing repos.

1. Scan the tree. Detect existing `design`, `docs`, `wip`, `adr`, `adrs`, `architecture`, `documentation` folders and any markdown outside code folders. Detect sub-projects by build files (`go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`).
2. Present a proposal: which folders will be created, which existing folders look like they map to `design`, `docs`, or `wip`, and which sub-projects were found.
3. Ask for folder names, offering the defaults and any detected candidates.
4. Create the structure and render the template files that do not conflict.
5. For each existing markdown file outside the new structure, ask: move to `design/system`, `design/adrs`, `docs/<audience>`, leave in place, or skip all remaining. Moves are `git mv` when the repo is git.
6. Write `system-flow.yaml`, run `flai check`, print a summary and next steps.

`--dry-run` prints the proposal and stops. `--yes` accepts every default.

## Internal structure

```
flai/
├── main.go
├── cmd/                 # one file per cobra command
├── internal/
│   ├── config/          # ~/.flai/config.json
│   ├── manifest/        # system-flow.yaml
│   ├── template/        # clone, cache, render
│   ├── workitem/        # parse, validate, transition, ID allocation
│   ├── narrative/       # wip/agents files
│   ├── metrics/         # reference implementation of metrics.md
│   ├── importer/        # analysis and proposal
│   ├── dashboard/       # docker run/stop
│   └── ui/              # prompts and tables
└── testdata/
```

External processes: `git` and `docker` are invoked as subprocesses and must be on `PATH`. See ADR 0010.
