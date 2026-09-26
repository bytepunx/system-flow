---
id: S-0120
type: story
nature: remediation
title: A project imported with flai import on the command line is served by the host flai and shows in the dashboard's switcher
status: done
owner: alex
created: 2026-09-26T05:18:06Z
updated: 2026-09-26T05:39:27Z
transitions:
  - to: ready
    at: 2026-09-26T05:24:11Z
    by: alex
  - to: in-progress
    at: 2026-09-26T05:24:26Z
    by: agent-S-0117
  - to: review
    at: 2026-09-26T05:39:16Z
    by: agent-S-0117
  - to: done
    at: 2026-09-26T05:39:27Z
    by: alex
tags: [cli]
touches: [flai/cmd/import.go, flai/cmd/import_serve.go, flai/cmd/new.go, flai/cmd/repo_url.go, flai/cmd/serve.go, flai/cmd/serve_import.go, flai/cmd/serve_status_below.go, flai/cmd/dashboard_agent.go, flai/internal/serve, design/adrs, docs/users/flai.md, docs/operators, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0120 A project imported with flai import on the command line is served by the host flai and shows in the dashboard's switcher

## Goal

A repository imported with `flai import` on the command line shows up in the dashboard's switcher without any other step, the same as one imported from the board. Remediates [I-0047](../../../design/issues/I-0047-a-project-imported-with-flai-import-on-the-command-line-is-neither-served-nor-offered-so-the-dashboard-never-shows-it.md).

## Acceptance criteria
- [x] When `flai import` (with or without `--commit`) writes the manifest and a `flai host` runs for the config, the project is registered with `flai serve` using the same entry `flai dashboard` writes, and the output says it is now served at the dashboard's address
- [x] When no host runs, `flai import` says the project is not served yet and names the one command that serves it
- [x] A git repository under an `import_roots` folder that has a `system-flow.yaml` but is not registered is no longer dropped from both lists: `flai serve` serves it, or offers it in the switcher to be served. The story records which in `## Decisions` and why
- [x] `flai serve status` lists such repositories and says why each one is not served
- [x] The `repo_url` prompt defaults to the origin remote as an https URL, and a value that is not a URL (such as `https://github.com:owner/repo`) is refused with the reason
- [x] `docs/users/flai.md` and `design/system/flai-cli.md` say how a project imported on the command line reaches the dashboard

## Tasks
- T-0437 flai import registers the project with flai serve when a host runs, and names the command that serves it when none does
- T-0438 flai serve serves the unregistered projects under import_roots, and flai serve status says why any is not served
- T-0439 The repo_url prompt defaults to the origin remote as an https URL and refuses a value that is not a URL
- T-0440 Users guide and flai-cli design say how a project imported on the command line reaches the dashboard

## Notes

S-0121 and S-0122 build on this. Found on 2026-09-25 importing `~/git/blog`. The only code that writes `serve/projects.json` is `flai dashboard` (`cmd/dashboard_agent.go` `connectServe`) and the board's import (`internal/serve/candidates.go`). `cmd/import.go` calls neither. `FindCandidates` skips any repository that has a manifest. So once the CLI import had written one, the blog was neither served nor offered. Restarting the dashboard or the host changes neither list. The workaround was `flai dashboard` in the blog. The blog's manifest also got `repo: https://github.com:arobson/astro-blog`, typed by hand at the prompt, from the origin `git@github.com:arobson/astro-blog.git`.

Delivered on `story/S-0120` (T-0437 to T-0440). Criterion 3: such a repository is served, not offered; why is in the narrative's `## Decisions` and [ADR-0045](../../../design/adrs/0045-a-system-flow-repository-under-a-folder-named-for-import-is-served-whether-or.md). S-0121, in progress at the same time, also edits `flai/internal/serve/serve.go` (`Status` and the status writing in `Run`); whichever is accepted second meets a small conflict there, resolved by keeping both sets of fields and lines. While both are in progress `flai check --strict` warns `wip.overlap`, which fails `make smoke`'s repository check; the template part of the smoke passed.
