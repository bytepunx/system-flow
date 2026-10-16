---
id: T-1894
type: task
nature: feature
title: The documentation convention allows the examples files under design/components, the folder has a README here and in the template, and the layout lists it
status: backlog
parent: S-0443
owner: alex
created: 2026-10-16T21:48:51Z
updated: 2026-10-16T21:48:51Z
transitions: []
stream: S-0443
tags: [template, dashboard]
touches: [design/conventions/documentation.md, template/root/design/conventions/documentation.md, design/system/documentation-standard.md, design/components/README.md, template/root/design/components/README.md, design/README.md, template/root/design/README.md.tmpl, CLAUDE.md, template/root/CLAUDE.md.tmpl, design/system/repository-layout.md, template/CHANGELOG.md, template/template.yaml]
after: [T-1884]
---
# T-1894 The documentation convention allows the examples files under design/components, the folder has a README here and in the template, and the layout lists it

## Work

- The baseline rule "Everything is markdown with YAML front matter, except the design token file" gains the examples files under `design/components/` as a second exception, in `design/conventions/documentation.md`, its template copy, and `design/system/documentation-standard.md`. The operator decided it on TH-0483 answer 4. The convention's text follows that decision as S-0381 followed TH-0445's for the token file. Record it in the story's `## Decisions`.
- `design/components/README.md`, here and in `template/root/design/components/`: what the folder holds, who writes it (flai's save and the story's agents), and its layout.
- `design/README.md` and its template, `CLAUDE.md`'s layout table and `template/root/CLAUDE.md.tmpl`, and `design/system/repository-layout.md` list the folder.
- `template/CHANGELOG.md` and the version in `template/template.yaml`.

Waits for T-1884, whose ADR the convention's words cite.

## Done when

The three convention texts carry the exception, the folder has its README in both places, the indexes list it, and `flai check --strict` and `flai test` on the changed files pass.

## Notes

The documentation convention is a baseline rule. It changes here only because the operator decided the folder; the change lands in `design/` and `template/` in this one story, as `CLAUDE.md` requires.
