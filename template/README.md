# system-flow template

The template that `flai` renders into a new or imported monorepo. This file describes the template repository itself and is not copied into projects.

- `template.yaml` is the manifest: version, variables, layout defaults, render rules.
- `root/` is rendered into the target repo root. Files ending in `.tmpl` are Go `text/template`; everything else is copied verbatim.
- `items/` holds the templates `flai` renders when creating epics, stories, tasks, and narratives inside a project.
- `projects/` will hold optional sub-project skeletons per kind. Empty in 0.1.0.

Development copy: this folder lives in the `system-flow` monorepo under `./template` until it is published to `github.com/bytepunx/system-flow-template`. The rendering rules and the manifest schema are specified in that monorepo's `design/system/template.md`.

Test it: `flai new /tmp/sample --template /path/to/template --defaults && flai check /tmp/sample`.
