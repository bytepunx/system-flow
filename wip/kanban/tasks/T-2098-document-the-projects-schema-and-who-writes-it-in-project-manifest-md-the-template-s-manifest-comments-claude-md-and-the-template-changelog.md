---
id: T-2098
type: task
nature: improvement
title: Document the projects schema and who writes it in project-manifest.md, the template's manifest comments, CLAUDE.md, and the template changelog
status: backlog
parent: S-0491
owner: alex
created: 2026-10-23T20:41:27Z
updated: 2026-10-23T20:41:27Z
transitions: []
stream: S-0491
tags: [template]
touches: [design/system/project-manifest.md, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
after: [T-2093, T-2096, T-2097]
---
# T-2098 Document the projects schema and who writes it in project-manifest.md, the template's manifest comments, CLAUDE.md, and the template changelog

## Work

The story's goal says no document in the project or the template gives the schema of `projects`. Close that once the commands, the tools, and the permission exist, so the text can name them:

- `design/system/project-manifest.md`: beside the `projects[].components` and `projects[].install` paragraphs, a `projects[]` paragraph giving each key (`name`, `path`, `kind`, `tags`, `components`, `install`), what `name` is used for (the release tag `<name>/vX.Y.Z`, the examples folder under `design/components/`, the components methods' `sub_project`), the kinds flai knows and that others are accepted, the rules `ProjectErrors` enforces and the `manifest.projects` finding, and that `flai subproject`, `settings.subproject`, and `subproject_edit` write the list, creating a missing folder with a README. Correct the line that says flai rewrites `projects` (true only for `flai import`).
- `template/root/system-flow.yaml.tmpl`: the comment block that documents the other keys gains `projects`, saying what an entry holds and that `flai subproject add` writes one.
- `CLAUDE.md` and `template/root/CLAUDE.md.tmpl`: the layout table's row for `system-flow.yaml` says `flai subproject` adds, changes, and removes sub-projects, beside the keys an agent may edit by hand.
- `template/CHANGELOG.md`: an entry for the template's change (the convention's and the agent definition's lines from T-2096, the manifest comment, the CLAUDE row), as the entries above it are written: the story, a patch bump, the ADR, what changed, and what an older flai does without it.

It waits for T-2093, T-2096, and T-2097, whose commands, permission, and section it names.

## Done when

- `design/system/project-manifest.md` describes every key of a `projects[]` entry, the validation, and the three ways the list is written, and no longer says flai rewrites `projects` outside import.
- `template/root/system-flow.yaml.tmpl`'s comments document `projects`, and `CLAUDE.md` and `template/root/CLAUDE.md.tmpl` name `flai subproject`.
- `template/CHANGELOG.md` has the story's entry, and `flai check --strict` and the markdown lint pass for the paths changed.

## Notes

`docs/operators/settings.md` keeps its `projects[].*` rows; T-2093 and T-2096 update it.
