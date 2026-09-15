---
title: flai CLI
updated: 2026-09-15
status: draft
---

# flai

The system-flow command line tool. Full command reference will be generated from the CLI in story S-017; the design is in [design/system/flai-cli.md](../../design/system/flai-cli.md). Implemented so far: `version` and `config`.

## Install

```bash
go install github.com/bytepunx/system-flow/flai@latest
```

Release binaries for Linux, macOS, and Windows arrive with story S-010.

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
