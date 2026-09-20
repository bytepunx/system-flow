---
title: flai CLI
updated: 2026-09-20
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
| `flai check [dir] [--strict]` | Validate manifest and layout, every item in kanban and archive (front matter, ID and file name, parents and children, state history and sequence, acceptance criteria, sections), narratives and index, board limits and order, design/docs front matter including ADRs, and the conventions folder (front matter, unique order, marker and project additions, length, README index), and design/issues (schema, unique IDs, file names, summary freshness), the manifest (`manifest.key` warns when `system-flow.yaml` has no `key`, which API responses use to name the project, ADR-0024), and the clone itself: `git.relative-worktrees` warns when `extensions.relativeWorktrees` is set and the git on `PATH` is older than 2.48, read from `.git/config` because that git refuses the repository (ADR-0022). One `path:line: level: rule: message` per finding. Errors exit 1; `--strict` makes warnings exit 1. |
| `flai epic new`, `flai story new --epic E-0001`, `flai task new --story S-0004` | Create an item from the body template, allocate the next ID, link to parent. `--body-stdin` reads the body below the heading from standard input (S-0059), and creation is then one step (`flai/internal/itemnew`): `flai check` before and after, anything the item introduces refuses it with exit 4 and the findings as `{ refused }`, the file removed, the parent restored, and the ID not spent. `--autocommit` commits the item and its parent, path limited, as `chore: [ID] create <type>: <title>` with `--trailer` lines, unless `dashboard.autocommit` is false. `--print-body` prints the body the project's template gives the type and creates nothing. |
| `flai move <id> <state>` | Transition an item with rule validation. `--reason` required for `cancelled` and for review to in-progress; `--by` defaults to `FLAI_AGENT` when it is set, else the config author, as for `flai block`, `unblock`, and threads, so an agent's moves are recorded as the agent's (S-0058). Warns on WIP limit breaches. A move to `cancelled` cancels everything open under the item too (ADR-0028): the items are listed first, a terminal is asked unless `--yes`, `--dry-run` lists and changes nothing, `--json` answers `{ id, status, cancelled[], left_behind[], warnings[] }`, and the narratives, story branches, and worktrees of cancelled stories are named and left as they are (S-0070). Moving a story from review to done is acceptance and runs the `flai accept` flow with its flags (S-0046). |
| `flai show <id>` | One item with history, blocks, and children. |
| `flai block <id> --reason`, `flai unblock <id>` | Open and close blocked intervals. |
| `flai board [--all]` | Stories per column with nature, age in column, blocked flag, WIP counts, and the pull order, with `backlog` and `ready` laid out in it (S-0057); `--all` adds epics and tasks. |
| `flai order <story> --before <id> \| --after <id> \| --top \| --bottom` | Places a ready or backlog story in its column's pull order (S-0057) and rewrites `order` in `board.md`: ready stories, then backlog stories, stale names dropped, the unplaced tail of the backlog left unnamed. Refuses epics, tasks, other states, and a reference in another column, with the reason. `--json` prints `{ id, status, sequence, order }`. Rules in [workflow.md](workflow.md). |
| `flai push --pending [--dry-run]` | On the host, push an acceptance that was made and not pushed (S-0063, ADR-0026): when the main checkout's branch is ahead of its remote-tracking branch and the commits ahead include an acceptance (`chore: [ID] accept and archive`), push the branch and the tags on unpushed commits with the host's credentials. "nothing pending" otherwise; ordinary commits are left to `git push`; refuses when the two have diverged; never forces. The question is asked offline (`flai/internal/pending`): commits ahead, the items accepted, the tags reachable from the branch and not from the remote-tracking branch, commits behind. `flai board` prints it and `board --json` carries it as `unpushed`. |
| `flai stats [--since 30d] [--type story] [--by nature\|type\|parent] [--json]` | The aggregates in [metrics.md](metrics.md) as a table; `--json` adds per-item values, weekly throughput, burn-up (whole set and per parent), cumulative flow, and aging. |
| `flai stream open <story-id> [--no-branch]`, `flai stream log <story-id> "<entry>"` | Create a narrative from the template and, in a git repository, branch `story/<id>` from the main branch checked out in a worktree under `.flai-cache/worktrees/<id>` (ADR-0019; `wip/` is written in the main checkout from anywhere); append a timestamped log entry. With the config key `worktrees.relative_paths` on and git 2.48 or newer the worktree is created with `--relative-paths`; on an older git flai warns and creates an ordinary one (ADR-0022). `FLAI_AGENT` and `FLAI_SESSION` identify the writer. `index.md` is regenerated from the active narratives after every open, log, move, and archive. |
| `flai stream diff <story-id>` | The story branch against its merge base with the main branch (S-0041): per file the path, status (added, modified, deleted, renamed), additions, deletions, and unified hunks; binary files without hunks; a file's patch and the whole result are capped and marked `truncated`. Read with git in the main checkout, so it does not need the worktree. `--json` is what the dashboard's review page shows. |
| `flai migrate ids [--dry-run]` | Widen every work item and issue ID to four digits: rename items and narratives in kanban and archive, and issues under `design/issues`, with `git mv` (plain rename outside git) and rewrite references in the layout folders, the root markdown and yaml files, and each project's root markdown files. Idempotent. |
| `flai archive [id...] [--dry-run]` | Move done and cancelled items and their narratives to `wip/archive`. Default: every closed epic whose stories are archived, every closed story with its tasks, and closed tasks whose story is gone from the board. |
| `flai dashboard [--image] [--tag] [--port] [--bind] [--push-key] [--push-known-hosts] [--pull] [--attach] [--open] [--build]`, `flai dashboard stop\|status\|logs` | Pull the flaiover image if missing (on an unauthorized pull from a non-Hub registry, `docker login` with `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` and one retry, else an error naming the `read:packages` scope; `--build` builds `flaiover:local` from `flaiover/` in the monorepo instead) and run it detached as `flaiover-<project>`, bound to localhost, repository mounted read-write at its own host path with `PROJECT_DIR` set to it (ADR-0022; `/project` only when the host path cannot be a container path, with a warning that stories with a branch must then be accepted from a shell), as the host user, with the host's git identity and its global git excludes file (read-only, named through `GIT_CONFIG_*`) passed in so git in the container commits as the operator and ignores what the host ignores (S-0046, S-0051); status and the already-running message inspect the real container. Precedence: flags, manifest `dashboard`, config., published on `dashboard.bind` (default every interface) With `--push-key <path>` or the host config key `dashboard.push_key`, never the manifest and off by default, the container is given an SSH key to push acceptances with ([ADR-0026](../adrs/0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md), S-0062): before anything starts the key is refused when it is missing, looser than 0600, not a private key, or protected by a passphrase, when `origin` is not an SSH URL, and when this host has no `known_hosts` entry for the remote (`--push-known-hosts` or `dashboard.push_known_hosts` names the file to take it from); then the key, the pinned host keys (`.flai-cache/dashboard.known_hosts`), and a passwd entry for the host's user (`.flai-cache/dashboard.passwd`) are mounted read-only and `GIT_SSH_COMMAND` names only them, with strict host key checking, batch mode, no agent, and no user configuration. The start message and `status` say which key the container holds, by fingerprint; `--json` carries `push_key`. Refused on Windows hosts for now. Every start also mounts `.git/hooks`, `.git/info`, `.git/config`, the token's file, and an in-clone flai config read-only over the clone (S-0064, [ADR-0027](../adrs/0027-git-hooks-config-and-info-are-read-only-in-the-dashboard-container.md)), creates the two directories when missing, and lists enabled hooks and executing git settings already in the clone; `status` reports whether the running container has the read-only paths (`git_read_only`) and the same list (`already_in_clone`). |
| `flai upgrade [--dry-run] [--force] [--keep-all\|--replace-all] [--relock]` | Bring the project to the template version at the configured source: add new files, merge marker files above the marker, replace files unchanged since applied (per `system-flow.lock.yaml`, ADR-0015), report the rest as conflicts with keep, replace, or diff; refuses a dirty tree; `--relock` records a hand-assembled project. |
| `flai prime [--cat] [--json]` | Print `design/conventions` in read order, README first, as paths or contents, so an agent or hook loads the norms in one call; the open-issues summary follows when non-empty. |
| `flai release <id> [--dry-run] [--deliver] [--apply]` | Compute the semver release for an item per the git convention (delivered component gets the delivery-type bump, touched components a patch; the delivered component is `--deliver`, else the one the item's tags and then its epic's tags name among the components its commits touched, the most touched files winning and the earlier tag breaking a tie, else the only touched one; a tag naming an untouched component never delivers, S-0047) and with `--apply` create the tags and version bumps. A `research` story plans no release ([ADR-0025](../adrs/0025-research-is-accepted-without-a-release.md)): `skipped` gives the reason and `unreleased` lists `{ component, files }` for each component its commits touched; an `experiment` story is an error. |
| `flai accept <id> [--by] [--deliver] [--no-release] [--no-push] [--trailer]` | The operator's acceptance in one step: preflight, before anything changes (committer identity; the workflow's rules for done, open tasks and unticked criteria, checked on a copy since S-0041; the story worktree readable by git), rebase and fast-forward the story branch into main and remove its worktree, move to done, archive, template bump, commit, tags on that commit, push. A failed push is reported (`push_error`), not fatal. Each step is logged as an info event as it completes (`step`: merged, done, archived, committed, tagged, pushed or not-pushed), which the dashboard streams as progress (S-0041). Completes a story that is done but unarchived. `--dry-run --json` is the dashboard's preview (plan, branch, blockers, uncommitted). Uncommitted changes under `wip/` are expected; anything else is refused unless `--yes` includes it in the acceptance commit. A dry run never refuses for them: it lists them in `uncommitted`, with or without `--yes`, so the choice is made before confirming (S-0051). |
| `flai issue new "<title>" --class <c> [--cost] [--note]`, `flai issue bump <id> [--cost] [--note]`, `flai issue close <id> --reason`, `flai issue list [--all]`, `flai issue summary` | Record, increment, and close recurring friction in `design/issues`; every command regenerates `summary.md` (average and total cost, most expensive first). |
| `flai template show`, `flai template update`, `flai template use <repo> [--ref]` | Inspect, refresh, and switch the template source. |
| `flai template push [dir] [--remote] [--ref] [--tag] [--dry-run] [--force]` | Publish a locally developed template to its git remote: clone the branch, replace contents, commit with the version, push, and with `--tag` push `v<version>` (refused if it exists). Defaults from `publish` in `template.yaml`. Git errors surface verbatim; no retries; no force push unless `--force`. `flai accept` calls it for template components. |
| `flai config get [key]`, `flai config set <key> <value>`, `flai config path` | Read and edit `~/.flai/config.json` by dotted key; `path` prints the resolved file without creating it. Keys: `template.*`, `dashboard.*`, `cache_dir`, `author`, and `worktrees.relative_paths` (boolean, default false, ADR-0022). |
| `flai version` | Version, commit, build date. |
| `flai dashboard token [--rotate]` (E-0006) | Print the per-project dashboard token and login link; `--rotate` replaces it and restarts the container. |
| `flai stream sync <story-id>` | Rebase `story/<id>` onto the main checkout's branch inside its worktree (`--autostash`); conflicts stop there and are listed. Run at every task transition (S-0037). |
| `flai touches <id> [paths...] [--clear]` | Set or show the advisory `touches` list on a story or task; `flai check` warns `wip.overlap` between in-progress items (S-0037). |
| `flai thread new --on <path\|id> [--heading] "<title>" "<text>"`, `reply <id> "<text>"`, `resolve <id> [--reason]`, `list [--on] [--all]`, `show <id>` | Threads under `wip/threads/` (`TH-nnnn`, ADR-0020) anchored to a document, a heading in it, or a work item; anchors are validated; replies from anyone but the opener set `answered`, the opener's follow-up sets `open`, `resolve` closes; unresolved threads on a story or its tasks are mirrored into the narrative's `## Open questions` (S-0038). |
| `flai doc show <path>`, `flai doc save <path> --hash <sha256> [--message] [--trailer] [--no-commit]` | Read and save one markdown document under design, docs, or wip for an editor (ADR-0023, S-0040). `show` prints the content, its SHA-256, and the edit mode: `full`, `body` (flai owns the front matter: work items, narratives, `board.md`), or `none` with the reason (generated files, threads, issues, archive). `save` reads the new content on stdin; fails with `conflict:` and the current content, hash, and a unified diff when the hash is stale; refuses a front matter change in `body` mode; runs the check before and after writing, and when the edit introduces any finding, or the file has an error, restores the file and fails with `refused:` (exit 4) and the findings; a conflict is exit 3; sets `updated` to today on design and docs files unless the designer changed it; commits that path only (`docs:` or `chore:`, item ID in brackets for items and narratives) unless `--no-commit` or `dashboard.autocommit: false`. Nothing is pushed. |
| `flai adr new "<title>" [--status] [--supersedes N]... [--refines N]... [--body-stdin] [--autocommit] [--trailer] [--print-body]`, `flai adr accept <N>` | Record an architecture decision (S-0060, `flai/internal/adr`). The number is one more than the highest `NNNN-*.md` in `design/adrs`; the file is `NNNN-slug.md` with `id`, `title`, `status` (`proposed` by default), `date`, `supersedes`, `superseded_by`, and `refines` when given; the body is the project's `0000-template.md` sections or standard input; the row is added to `design/adrs/README.md` with the index's notes ("accepted, refines 0018"); a superseded ADR gets `superseded_by` and its row noted. `flai check` runs with everything in place, and a finding the ADR introduces, or any failure after files were written, puts every file back (exit 4 with `{ refused }` for a check refusal). `--autocommit` makes one `docs: ADR-NNNN <title>` commit of what was written unless `dashboard.autocommit` is false. `accept` sets a proposed ADR to accepted with today's date and updates its row. `flai check` rule `adr.index`: a file with no row, a row with no file. `flai doc` treats an accepted ADR as not editable. |
| `flai mcp` | MCP server on stdio for agent sessions (ADR-0020, S-0039), built on the official Go SDK. Tools: `inbox` (unresolved threads, `awaiting: you` when the last entry is not the agent's; the stories ready to pull in pull order with `can_pull` from the in-progress limit; and the changes others made since this agent last looked, from a per-agent cursor under `.flai-cache/mcp/`, S-0058; at most 50 per look, newest kept, the rest counted in `changes_omitted` and never reported later, and a first look with no cursor covers 24 hours of stories and epics only, S-0061), `board` (what `flai board --json` prints), `thread_get`, `thread_open`, `thread_reply`, `thread_resolve`, `item_get`, `item_move` (workflow rules enforced; refuses to move a story or epic to done), `doc_get` (markdown under design, docs, wip only), `who_touches`, and `wait_for_events` (returns at once when the cursor is behind, otherwise polls thread, kanban, and narrative files every 250 ms and returns on the first change or the timeout, capped at five minutes; returns events as well as paths and advances the cursor). Resources: `flai://design/{+path}` and `flai://docs/{+path}`. The agent is `FLAI_AGENT`; every write goes through the same packages as the CLI; stdout carries the protocol only. |
| `flai self-upgrade [--check] [--version X.Y.Z] [--dir DIR]` | Resolve the newest `flai/v*` release (or the pinned one), download the platform archive and `checksums.txt` through the release asset API, verify the SHA-256, and replace the running executable atomically (`--dir` installs elsewhere). Token from `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token`. `install.sh` at the repository root is the same operation for a shell one-liner; `scripts/install-test.sh` exercises both against the real latest release in the smoke tier. |

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
