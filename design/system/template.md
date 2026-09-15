---
title: Template repository
updated: 2026-09-15
status: active
---

# Template repository

The template is the executable form of the standard. It lives in its own git repository so it can be versioned, forked, and pointed at from `flai` config. Until it is published, the prototype lives in [./template](../../template) in this monorepo and `flai` can be pointed at a local path.

Default source: `https://github.com/bytepunx/system-flow-template`, ref `main`. Overridable per user in `~/.flai/config.json` and per project in `system-flow.yaml`.

## Structure

```
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
│   ├── design/ ...
│   ├── docs/ ...
│   └── wip/ ...
├── items/                 # rendered by flai at item creation time, not at project creation
│   ├── epic.md.tmpl
│   ├── story.md.tmpl
│   ├── task.md.tmpl
│   └── narrative.md.tmpl
└── projects/              # optional sub-project skeletons, one per kind
    ├── go/
    └── sveltekit/
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
  ignore: [README.md]      # template-repo-only files at root/, not copied
items:
  epic: items/epic.md.tmpl
  story: items/story.md.tmpl
  task: items/task.md.tmpl
  narrative: items/narrative.md.tmpl
projects:
  - kind: go
    path: projects/go
  - kind: sveltekit
    path: projects/sveltekit
```

## Built-in variables

Available in every rendered file in addition to the manifest variables.

| Variable | Value |
|----------|-------|
| `.layout.design`, `.layout.docs`, `.layout.wip` | Chosen folder names |
| `.today` | `YYYY-MM-DD` UTC at render time |
| `.now` | `YYYY-MM-DDTHH:MM:SSZ` UTC at render time |
| `.template.repo`, `.template.ref`, `.template.version` | Template source being applied |
| `.id`, `.title`, `.nature`, `.parent`, `.owner`, `.agent`, `.session` | Item templates only, supplied by `flai epic|story|task new` and `flai stream open` |

Template functions: `initials`, `slug`, `upper`, `lower`, `quote` (YAML-safe double-quoted scalar, use it for every free-text value in YAML files). For a local template source, `.template.repo` is the absolute directory so the manifest stays meaningful from any working directory.

## Rendering rules

- Files ending in `.tmpl` are rendered with Go `text/template`. Variables are available as `{{ .project_name }}`. Layout folder names are available as `{{ .layout.design }}` and are used inside rendered files so renamed folders stay consistent.
- Every other file is copied byte for byte.
- Paths containing a layout folder name at their first segment are rewritten to the chosen name at render time, so `root/design/README.md` lands in whatever `layout.design` is.
- Rendering never overwrites an existing file unless `--force`. On import, `flai` reports each conflict and offers keep, replace, or show diff.
- Executable bits are preserved.

## The baseline CLAUDE.md

The most important file in the template. It gives an agent, in one read, the repo layout, the work hierarchy, the narrative obligations, and the definition of done. It is deliberately short and links to `design/system` for detail. Projects append their own sections below a marker line; `flai upgrade` (future) replaces only the section above the marker.

## Versioning

Semantic versions in `template.yaml`. `flai` records the applied version in `system-flow.yaml`. Breaking changes to layout or front matter schema bump the major version and ship with a migration note in the template's `CHANGELOG.md`.

## Upgrading a project (S-020)

`flai upgrade` brings a conforming repo to the template version at the configured source. It needs to know which rendered files the project has since changed; the proposed mechanism is a `system-flow.lock.yaml` written by `flai new` and `flai upgrade` with a sha256 per rendered path. Unchanged files are replaced, new files added, changed files reported as conflicts, and `CLAUDE.md` is merged above its marker line. The lock file decision gets an ADR when S-020 starts.

## Publishing a template (S-021)

A template developed inside another repository, as `./template` is here, is published with `flai template push`. The manifest declares its home:

```yaml
publish:
  repo: git@github.com:bytepunx/system-flow-template.git
  ref: main
```

The command clones the remote branch into the cache, replaces its contents with the local template, commits with the template version in the message, pushes, and with `--tag` also pushes `v<version>`. It assumes push permission exists and reports git failures verbatim. It never changes the calling project's `system-flow.yaml`.

## Testing the template

The monorepo's CI renders the prototype template into a temporary directory with `flai new --template ./template` and runs `flai check` on the result. This is the sample repo used by dashboard fixture tests.
