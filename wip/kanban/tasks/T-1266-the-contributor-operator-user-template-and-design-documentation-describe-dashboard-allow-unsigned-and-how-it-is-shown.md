---
id: T-1266
type: task
nature: feature
title: The contributor, operator, user, template, and design documentation describe dashboard.allow_unsigned and how it is shown
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:58Z
updated: 2026-10-07T23:09:58Z
transitions: []
stream: S-0239
tags: [cli, dashboard]
touches: [docs/contributors/index.md, docs/operators/index.md, docs/users/flai.md, docs/users/flaiover.md, template/root/docs/operators/index.md.tmpl, template/CHANGELOG.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/release-signing.md]
after: [T-1264, T-1265]
---
# T-1266 The contributor, operator, user, template, and design documentation describe dashboard.allow_unsigned and how it is shown

## Work

Criterion 4 and the documentation of the rest. It waits for T-1264 and T-1265, so that it describes the banner and the statuses in the words they print.

- `docs/contributors/index.md`: this repository builds both components from source, so set `dashboard.allow_unsigned: true` in `.flai-cache/config.json`, and expect the banner and the `warn`.
- `docs/operators/index.md`: beside S-0237's and S-0238's sections on the stamp and the image check, what the allowance lets through, that it is for development builds only, where it is shown, and that the credential is still checked.
- `template/root/docs/operators/index.md.tmpl`: the same for a project made from the template, and an entry in `template/CHANGELOG.md`.
- `docs/users/flai.md`: `flai dashboard --allow-unsigned`, `FLAIOVER_ALLOW_UNSIGNED` with `--build`, and "unsigned allowed" in the three statuses. `docs/users/flaiover.md`: the banner and the connection list's "unsigned allowed".
- `design/system/project-manifest.md`: the manifest's `dashboard.allow_unsigned`. `design/system/flai-cli.md`: the configuration key, the flag, and the statuses. `design/system/flaiover-dashboard.md`: the banner and the variable. `design/system/release-signing.md` § The development case: as built.

## Done when

- Each document above says what the setting does and how it is shown, in the words the code prints.
- The markdown lint passes with `flai test` on the changed paths.

## Notes

Drafted by the planner.
