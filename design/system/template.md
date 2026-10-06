---
title: Template repository
updated: 2026-10-06
status: active
topics: [template]
---

# Template repository

The template is the executable form of the standard. It lives in its own git repository, [bytepunx/system-flow-template](https://github.com/bytepunx/system-flow-template) (private, published 2026-09-17 at 1.0.0), so it can be versioned, forked, and pointed at from `flai` config. The development copy is [./template](../../template) in this monorepo; this repository points `flai` at that local path and publishes it with `flai template push` when a release bumps it (see [Publishing a template](#publishing-a-template)). How to fork, edit, test, and version it is the [template guide](../../docs/contributors/template.md).

Default source: `https://github.com/bytepunx/system-flow-template`, ref `main`. Overridable per user in `~/.flai/config.json` and per project in `system-flow.yaml`. `main` follows releases: with no `--ref`, flai applies the template's newest version tag ([Which version is applied](#which-version-is-applied), [ADR-0103](../adrs/0103-flai-new-import-and-upgrade-follow-the-template-s-releases-the-newest-version.md)).

## Structure

```text
<template repo>/
├── README.md              # about the template itself, not copied
├── template.yaml          # manifest: version, variables, prompts, layout, file rules
├── root/                  # everything below is rendered into the target repo root
│   ├── CLAUDE.md.tmpl
│   ├── README.md.tmpl
│   ├── system-flow.yaml.tmpl
│   ├── .editorconfig
│   ├── .gitignore
│   ├── .gitattributes
│   ├── .markdownlint.yaml
│   ├── Makefile
│   ├── .github/workflows/system-flow-check.yml
│   ├── design/ ...         # includes conventions/, the baseline agent norms, see conventions.md
│   ├── docs/ ...
│   └── wip/ ...
├── items/                 # rendered by flai at item creation time, not at project creation
│   ├── epic.md.tmpl
│   ├── story.md.tmpl
│   ├── task.md.tmpl
│   └── narrative.md.tmpl
└── projects/              # reserved for sub-project skeletons, one per kind; empty
```

## Manifest

```yaml
# template.yaml
version: 0.1.0
min_flai: 0.1.0
variables:
  - name: project_name
    prompt: Project name
    required: true
  - name: project_key
    prompt: Short project key
    default: "{{ .project_name | initials }}"
  - name: description
    prompt: One line description
    default: ""
  - name: owner
    prompt: Owner or team
    default: ""
layout:                    # renameable folders and their defaults
  design: design
  docs: docs
  wip: wip
render:
  suffix: .tmpl            # files with this suffix are rendered with Go text/template and the suffix removed
  ignore: []               # paths under root/ never rendered: exact, or filepath.Match globs; a matched directory is skipped whole
items:
  epic: items/epic.md.tmpl
  story: items/story.md.tmpl
  task: items/task.md.tmpl
  narrative: items/narrative.md.tmpl
projects: []               # parsed, not yet used: sub-project skeletons do not exist
publish:                   # the template's git home, used by flai template push
  repo: git@github.com:bytepunx/system-flow-template.git
  ref: main
```

The manifest is parsed strictly; an unknown key is an error. `version` is required. `min_flai` is printed by `flai template show` and not enforced. A `dashboard` section is accepted and ignored: a project's dashboard settings come from `root/system-flow.yaml.tmpl`.

## Variables

`flai new` and `flai import` take each variable, in manifest order, from `--var`, else its default (an empty `project_name` defaults to the target directory's name), else a prompt in a terminal. A required variable left empty, whichever way it was set, stops the command before anything is written or, for `flai import`, moved (S-0185, I-0041). `flai new` records every value in `system-flow.lock.yaml` under `vars`; `flai upgrade` renders with them, as [Upgrading a project](#upgrading-a-project) says, so a variable a fork adds survives upgrades. `flai import` writes no lock, so an imported project records its values at its first `flai upgrade` or `flai upgrade --relock`.

## Built-in variables

Available in every rendered file in addition to the manifest variables.

| Variable | Value |
|----------|-------|
| `.layout.design`, `.layout.docs`, `.layout.wip` | Chosen folder names |
| `.today` | `YYYY-MM-DD` UTC at render time |
| `.now` | `YYYY-MM-DDTHH:MM:SSZ` UTC at render time |
| `.template.repo`, `.template.ref`, `.template.version` | Template source being applied |
| `.id`, `.title`, `.nature`, `.parent`, `.owner` | Epic, story, and task templates only, supplied by `flai epic new`, `flai story new`, `flai task new`. Pipe `.title` through `quote` in YAML. |
| `.id`, `.title`, `.agent`, `.session` | The narrative template only, supplied by `flai stream open`. |

Item templates see `.now` and `.today` besides these, and not the manifest variables or `.layout`. A variable default that is itself a template sees the variables listed before it and the built-ins, with `.template.repo` and `.template.ref` empty. A reference to data that does not exist stops the render with the file name (`missingkey=error`).

Template functions: `initials`, `slug`, `upper`, `lower`, `quote` (YAML-safe double-quoted scalar, use it for every free-text value in YAML files). For a local template source, `.template.repo` is the absolute directory so the manifest stays meaningful from any working directory.

## Rendering rules

- Files ending in `.tmpl` are rendered with Go `text/template`. Variables are available as `{{ .project_name }}`. Layout folder names are available as `{{ .layout.design }}` and are used inside rendered files so renamed folders stay consistent.
- Every other file is copied byte for byte.
- Paths containing a layout folder name at their first segment are rewritten to the chosen name at render time, so `root/design/README.md` lands in whatever `layout.design` is.
- Rendering never overwrites an existing file unless `--force`. `flai new` and `flai import` keep an existing file and report how many they kept.
- Executable bits are preserved.

## The baseline CLAUDE.md

The entry point for an agent. It opens with a priming section that tells the agent to read `conventions/` in order, then the streams index, then the board, before any change (S-0024). The rest is a map: the layout table and pointers to `design/system`. Norms themselves live in `design/conventions/`, not here. Projects append their own sections below a marker line; `flai upgrade` replaces only the section above the marker.

## Conventions

`root/design/conventions/` holds the baseline agent norms, one file per topic with a marker line for project additions, exactly as specified in [conventions.md](conventions.md). They are copied verbatim (no `.tmpl`) so they read the same in the template repository and in projects.

## Versioning

Semantic versions in `template.yaml`. Each release is tagged `v<version>` by `flai template push --tag`. `flai` records the applied ref and version in `system-flow.yaml` and in `system-flow.lock.yaml`. Breaking changes to layout or front matter schema bump the major version and ship with a migration note in the template's `CHANGELOG.md`.

## Which version is applied

`flai new`, `flai import`, and `flai upgrade` follow the template's releases ([ADR-0103](../adrs/0103-flai-new-import-and-upgrade-follow-the-template-s-releases-the-newest-version.md)).

- **The ref.** `flai new` and `flai import` take `template.ref` from the config; `flai upgrade` takes it from `system-flow.yaml`. `--ref` wins over both, and `--template` over `template.repo`.
- **Refs that follow releases.** With no `--ref`, a ref that is empty, the template's default branch (`main`), or a version tag resolves to the newest version tag of a git template: `vX.Y.Z`, the newest in semantic version order, pre-releases excluded. A template with no version tags falls back to the ref. Another branch or a commit is used as given. A local template directory has no tags and is used as it is. The config's default `template.ref` is `main`, so a new project starts at the newest release.
- **`--ref`.** It names what to apply. `--ref 1.0.60` matches the tag `v1.0.60`. It is written to `template.ref` in `system-flow.yaml` and, by `flai new` and `flai upgrade`, to the lock.
- **The version a project is at** is the lock's, or the manifest's when there is no lock. The lock records what was applied; it is never a target.
- **The cache.** A git template is cloned into the cache once per repo and ref. A clone at a branch is fetched again before each use; a clone at a tag is reused. When the fetch fails, as offline, flai warns and uses the clone it has.

## Upgrading a project

`flai upgrade` brings a conforming repo from the version its lock records to the version [Which version is applied](#which-version-is-applied) resolves, at `template.repo` in `system-flow.yaml` or `--template` ([ADR-0015](../adrs/0015-template-lock-file.md), [ADR-0103](../adrs/0103-flai-new-import-and-upgrade-follow-the-template-s-releases-the-newest-version.md)).

With no `--ref`, it asks the operator which version to apply when `system-flow.yaml` names a different `template.ref` or `template.version` than the lock recorded, because the operator edited it, and that differs from the newest tag. A changed version names the tag `v<version>`. The choices are the version the manifest names, the newest tag, and changing nothing. An edit that equals the newest tag is applied without asking. A recorded tag equal to the lock's is not a pin: the next upgrade takes the newest tag without asking. Without a terminal, or with `--yes`, a run that would ask changes nothing and exits non-zero, naming the choices and the `flai upgrade --ref <tag>` to run. `--dry-run` says which version it would apply and whether it would ask. The version applied, chosen or not, is written to `template.ref` and `template.version` and to the lock.

`system-flow.lock.yaml`, written by `flai new` and `flai upgrade`, records a sha256 per rendered path, the `topics` the template gave each marker file, and under `vars` the value of every template variable the project was last rendered with. On upgrade each template path is classified: added when absent, merged when both sides carry the baseline marker (template above, project below, except that a project's own `topics` on the file stay: [conventions.md](conventions.md#topics)), replaced when the project file still matches the lock, skipped when identical, otherwise a conflict. In a terminal each conflict offers keep, replace, or a diff; non-interactive runs need `--keep-all` or `--replace-all` and otherwise change nothing. The manifest's `template.version` and `template.applied` and the lock are updated only when no conflict is unresolved. A dirty git tree is refused unless `--force`, so the upgrade is reviewable as one diff. Upgrade never prompts for a variable. It renders with each variable the template defines, taken from `--var`, then, for the five standard variables, from `system-flow.yaml` (`name`, `key`, `description`, `owner`, `repo`; `--var` refuses these), then from the values the lock recorded (a required variable recorded empty is passed over), then from the variable's default, naming each variable that took it; a required variable left empty is named and nothing is changed (S-0185, I-0040). `--var` on a project already at the template's version re-applies it. See [Variables](#variables). `system-flow.yaml` itself is never re-rendered; only its `template` fields are updated. A file removed from the template stays in projects. `--relock` writes the lock at the current version for projects assembled by hand.

## Publishing a template

A template developed inside another repository, as `./template` is here, is published with `flai template push`. The manifest's `publish` section names its home. The command clones the remote branch into the cache (creating the branch from the default branch if it does not exist), replaces its contents with the local template, commits with the template version in the message, pushes, and with `--tag` also pushes `v<version>`, refusing if that tag exists. It assumes push permission exists, reports git failures verbatim, never retries, and never force-pushes without `--force`. It never changes the calling project's `system-flow.yaml`. `flai release --pending` runs it with `--tag` for a template component after bumping the version, as part of publishing everything accumulated (S-0087); `flai accept` no longer bumps or publishes anything. The first publish (S-0003) used a subtree split to seed the history; every push since is a commit on top.

## Testing the template

The monorepo's CI (`make smoke`, `scripts/template-test.sh`) renders the template into a temporary directory with `flai new --template ./template` and runs `flai check` on the result. This is the sample repo used by dashboard fixture tests.
