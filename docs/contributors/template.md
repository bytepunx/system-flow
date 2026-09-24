---
title: Template guide
updated: 2026-09-24
status: active
---

# Forking and editing the template

The template is what `flai new` and `flai import` render into a repository, what `flai upgrade` brings a repository up to, and where `flai` finds the bodies of new epics, stories, tasks, and narratives. This guide is for anyone who forks it for their organisation or edits it here. The design behind it is [template.md](../../design/system/template.md); the decision to lock rendered files is [ADR-0015](../../design/adrs/0015-template-lock-file.md).

The published template is [bytepunx/system-flow-template](https://github.com/bytepunx/system-flow-template). In this monorepo it is developed in [`template/`](../../template) and published from there.

## What is in a template

```text
<template>/
├── template.yaml      # the manifest: version, variables, layout, render rules
├── CHANGELOG.md       # one entry per version
├── README.md          # about the template; never rendered
├── root/              # rendered into the repository root
│   ├── CLAUDE.md.tmpl
│   ├── system-flow.yaml.tmpl
│   ├── design/ docs/ wip/ scripts/ .github/ ...
│   └── ...
├── items/             # rendered when flai creates a work item or narrative
│   ├── epic.md.tmpl
│   ├── story.md.tmpl
│   ├── task.md.tmpl
│   └── narrative.md.tmpl
└── projects/          # reserved for sub-project skeletons; empty
```

Only `root/` reaches a new repository. Everything beside it belongs to the template.

## Fork it and point flai at the fork

1. Fork or copy the template repository. A local directory works as well as a git remote; flai treats any existing directory as a local template and anything with `://` or starting with `git@` as a remote.
2. Tell flai where it is. Pick the scope you need:

   | Scope | How | Used by |
   |-------|-----|---------|
   | One command | `--template <url-or-dir>` and `--ref <branch-tag-or-commit>` | `flai new`, `flai import`, `flai upgrade`, `flai template show` |
   | Your machine | `flai template use <url-or-dir> --ref <ref>`, which writes `template.repo` and `template.ref` in `~/.flai/config.json` | `flai new`, `flai import`, `flai template show`, `flai template update`, `flai template push` with no directory |
   | One project | `template.repo` and `template.ref` in its `system-flow.yaml`, written when the project is created | `flai upgrade`, and the item templates for new work items |

3. Check the manifest parses and the variables are what you expect:

   ```bash
   flai template show --template ./my-template
   ```

A git template is cloned once into `<cache_dir>/templates/`, one clone per repository and ref, and reused after that. `flai template update` re-fetches the source in your config. A branch ref therefore goes stale in the cache; pin projects to a tag (`--ref v1.2.0`) so a new version is a new clone. Commands that create work items never fetch: they use the project's template only once it is in the cache, and flai's built-in item templates otherwise.

## template.yaml

The manifest is parsed strictly: an unknown key is an error, which catches typos.

| Key | Required | What it does |
|-----|----------|--------------|
| `version` | yes | Semantic version of the template. Rendered as `.template.version` and recorded in the project's `system-flow.yaml` and `system-flow.lock.yaml`. `flai upgrade` compares it with the project's. |
| `min_flai` | no | The oldest flai the template is meant for. `flai template show` prints it; nothing enforces it yet. |
| `description` | no | One line, printed by `flai template show`. |
| `variables` | no | The values asked for at render time. See [Variables](#variables). |
| `layout` | no | Renameable top-level folders, key to default name, such as `design: design`. `flai new --layout design=architecture` renames one. See [Render rules](#render-rules). |
| `render.suffix` | no | Files ending in it are rendered as templates and lose it. Default `.tmpl`. |
| `render.ignore` | no | Paths under `root/` never rendered. See [Render rules](#render-rules). |
| `items` | no | Item type to its template file, relative to the template root: `epic`, `story`, `task`, `narrative`. A missing entry falls back to flai's built-in template. |
| `publish.repo`, `publish.ref` | no | Where `flai template push` publishes this template. |
| `dashboard` | no | Accepted and ignored. The dashboard settings a project gets come from `root/system-flow.yaml.tmpl`. |
| `projects` | no | Accepted and ignored until sub-project skeletons exist. |

## Variables

Each variable has a `name`, and optionally a `prompt`, a `default`, and `required: true`:

```yaml
variables:
  - name: project_name
    prompt: Project name
    required: true
  - name: project_key
    prompt: Short project key
    default: "{{ .project_name | initials }}"
```

`flai new` and `flai import` choose each value, in the order the variables are listed:

1. A `--var name=value` wins. A name the manifest does not define is an error.
2. Otherwise the default. A default containing `{{` is itself a template and can use any variable listed before it and the [functions](#functions). An empty `project_name` defaults to the target directory's name.
3. In a terminal, without `--defaults`, flai prompts with `prompt` and the default filled in. Anywhere else the default is taken without asking.
4. A `required` variable left empty stops the render and names it.

Every variable defined in the manifest is available in every rendered file, as `{{ .name }}`, even when its value is empty.

Keep the five variables the template ships with: `project_name`, `project_key`, `description`, `owner`, `repo_url`. `flai upgrade` does not prompt. It renders with those five only, read back from the project's `system-flow.yaml` (`name`, `key`, `description`, `owner`, `repo`), so a file that uses a variable you added renders on `flai new` and fails on `flai upgrade`. Until upgrade learns more variables, derive anything else from those five in the template itself.

### Built-in data

Available in every file under `root/` and in variable defaults, besides the variables:

| Data | Value |
|------|-------|
| `.layout.<key>` | The folder name chosen for each `layout` key |
| `.today` | Render date, `YYYY-MM-DD`, UTC |
| `.now` | Render time, `YYYY-MM-DDTHH:MM:SSZ`, UTC |
| `.template.repo` | The template source; the absolute directory for a local template. Empty in variable defaults. |
| `.template.ref` | The ref rendered; empty for a local template or the default branch. Empty in variable defaults. |
| `.template.version` | `version` from the manifest |

### Functions

| Function | Does | Example |
|----------|------|---------|
| `quote` | YAML-safe double-quoted string. Use it for every free-text value in a YAML file. | `name: {{ .project_name \| quote }}` |
| `initials` | First letter of each word, lower case | `system-flow` gives `sf` |
| `slug` | Lower case, runs of other characters become `-` | `My Project!` gives `my-project` |
| `upper`, `lower` | Change case | `{{ .project_key \| upper }}` |

Everything else in Go [text/template](https://pkg.go.dev/text/template) works: `if`, `with`, `range`, comparisons. A reference to data that does not exist, such as a misspelt variable, stops the render and names the file.

## Render rules

flai walks `root/` and, for every path under it:

- A path matching a `render.ignore` pattern is skipped; a matching directory is skipped with everything in it. Patterns are relative to `root/`, matched exactly or with Go [`filepath.Match`](https://pkg.go.dev/path/filepath#Match), where `*` does not cross `/`.
- A file ending in `render.suffix` is rendered with the variables and built-in data, and written without the suffix: `root/README.md.tmpl` becomes `README.md`.
- Every other file is copied byte for byte, dotfiles included.
- The file mode is kept, so a script committed as executable in the template is executable in the project. Git records the executable bit; set it with `git update-index --chmod=+x` on hosts that lose it.
- A path whose first folder is a `layout` default is written under the chosen name: with `--layout design=architecture`, `root/design/README.md` becomes `architecture/README.md`. Only the first folder is renamed and only the path; a file that mentions the folder in its text says `{{ .layout.design }}` and ends in `.tmpl` so the text follows.
- An existing file is kept, not overwritten, unless `--force`. `flai new` reports how many it kept.
- Every directory is created, but git does not keep empty ones; put a `README.md` in a folder that must exist.

To add a file, put it under `root/` where it should land. Give it the suffix only if it uses the data.

## Files a project changes: the marker and the lock

Rendering is not the end of a file's life: `flai upgrade` renders the template again, in memory, and decides per path what to do.

| The project's file | Upgrade |
|--------------------|---------|
| Missing | Adds it |
| Identical to the new render | Leaves it |
| Carries `<!-- system-flow:end-of-baseline -->`, and so does the template's | Takes the template's text above the marker and keeps the project's from the marker down |
| Unchanged since it was rendered, by the hash in `system-flow.lock.yaml` | Replaces it |
| Anything else | Reports a conflict: keep, replace, or diff in a terminal; `--keep-all` or `--replace-all` otherwise |

`system-flow.yaml` is the project's from the moment it is rendered; upgrade only updates its `template` fields.

So, when you edit:

- Put text every project must share above the marker in `CLAUDE.md.tmpl` and in each `design/conventions/` file, and leave the space below it (`## Project additions`) for projects. Keep exactly one marker per file; `flai check` fails a convention file without one.
- Expect a conflict in every project that edited a file with no marker. Add a marker to a file projects are meant to extend.
- A file you remove from the template is not removed from projects.

## Item templates

`items/` holds the templates `flai epic new`, `flai story new`, `flai task new`, and `flai stream open` render. They see item data only, not the project's variables:

| Template | Data |
|----------|------|
| `epic`, `story`, `task` | `.id`, `.title`, `.nature`, `.parent`, `.owner`, `.now`, `.today` |
| `narrative` | `.id`, `.title`, `.agent`, `.session`, `.now`, `.today` |

The output must be a valid work item or narrative: the front matter in [work-hierarchy.md](../../design/system/work-hierarchy.md) and [agent-narrative.md](../../design/system/agent-narrative.md), then a `# <ID> <Title>` heading. Quote the title in front matter with `{{ .title | quote }}`. The sections below the heading are yours to choose; the dashboard's new-item form starts from them.

## Testing a change

Render into a scratch directory and check the result:

```bash
flai new /tmp/sample --template ./my-template --defaults --no-git
flai check /tmp/sample --strict
```

`--defaults` exercises every default; add `--var` and `--layout` to exercise the rest. Then test the path existing projects take, from the last published version to your working copy:

```bash
flai new /tmp/old --template <url> --ref v1.2.0 --defaults
cd /tmp/old && git add -A && git commit -qm init
flai upgrade --template /path/to/my-template --dry-run
```

The dry run lists what each file would get. Run it without `--dry-run` and `flai check --strict` again to see the result.

In this monorepo, `make smoke` renders `template/` and checks it, and CI runs it on every change.

## Versioning

The version is the `version` line of `template.yaml`, with an entry in `CHANGELOG.md` for each one.

| Change | Bump |
|--------|------|
| A layout folder, a front matter field, a variable projects rely on, or anything else a project must migrate by hand | Major, with a migration note in the changelog entry |
| A new file, section, or convention rule projects take without editing | Minor |
| Wording, a fix, a rule clarified | Patch |

In a fork, bump the version and add the changelog entry by hand, then publish. Projects move to it with `flai upgrade`, after pointing `template.ref` at the new tag or passing `--ref`.

In this monorepo the release tooling does both. `template/` is a component of kind `template` in `system-flow.yaml`; a story tagged `template` that changes it bumps the version by its nature when the next release is computed, per [git.md](../../design/conventions/git.md). Changes to conventions land in `design/system` first, then in `template/root/design/conventions/`, then above the marker in `design/conventions/`, in the same story.

## Publishing

A template developed inside another repository, as `template/` is here, is published to the repository named by `publish.repo` and `publish.ref`:

```bash
flai template push ./template --dry-run   # list the changes and the commit message
flai template push ./template --tag       # push, and tag v<version>
```

The command clones the branch (creating it if missing), replaces its contents with the local template, commits with the version in the message, and pushes. `--tag` refuses a tag that exists. `--remote` and `--ref` publish somewhere else; `--force` force-pushes and replaces a tag, and needs the same care as any force push. A template that is its own repository needs none of this: commit, tag `v<version>`, push.

In this monorepo publishing is part of releasing: `flai release --pending`, the board's Publish action, and `flai push --pending --publish` bump `template/` when an accepted item released it, push this repository, and then run `flai template push --tag`. `flai push --pending` without `--publish` bumps and pushes this repository but leaves the template repository behind; run `flai template push ./template --tag` after it.
