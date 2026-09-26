---
id: S-0117
type: story
nature: remediation
title: A project imported with flai import on the command line is served by the host flai and shows in the dashboard's switcher
status: backlog
owner: alex
created: 2026-09-26T05:18:06Z
updated: 2026-09-26T05:18:06Z
transitions: []
tags: [cli]
touches: [flai/cmd/import.go, flai/internal/serve]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0117 A project imported with flai import on the command line is served by the host flai and shows in the dashboard's switcher

## Goal

A repository imported with `flai import` on the command line shows up in the dashboard's switcher without any other step, the same as one imported from the board. Remediates [I-0045](../../../design/issues/I-0045-a-project-imported-with-flai-import-on-the-command-line-is-neither-served-nor-offered-so-the-dashboard-never-shows-it.md).

## Acceptance criteria
- [ ] When `flai import` (with or without `--commit`) writes the manifest and a `flai host` runs for the config, the project is registered with `flai serve` using the same entry `flai dashboard` writes, and the output says it is now served at the dashboard's address
- [ ] When no host runs, `flai import` says the project is not served yet and names the one command that serves it
- [ ] A git repository under an `import_roots` folder that has a `system-flow.yaml` but is not registered is no longer dropped from both lists: `flai serve` serves it, or offers it in the switcher to be served. The story records which in `## Decisions` and why
- [ ] `flai serve status` lists such repositories and says why each one is not served
- [ ] The `repo_url` prompt defaults to the origin remote as an https URL, and a value that is not a URL (such as `https://github.com:owner/repo`) is refused with the reason
- [ ] `docs/users/flai.md` and `design/system/flai-cli.md` say how a project imported on the command line reaches the dashboard

## Tasks

## Notes

S-0118 and S-0119 build on this. Found on 2026-09-25 importing `~/git/blog`. The only code that writes `serve/projects.json` is `flai dashboard` (`cmd/dashboard_agent.go` `connectServe`) and the board's import (`internal/serve/candidates.go`). `cmd/import.go` calls neither. `FindCandidates` skips any repository that has a manifest. So once the CLI import had written one, the blog was neither served nor offered. Restarting the dashboard or the host changes neither list. The workaround was `flai dashboard` in the blog. The blog's manifest also got `repo: https://github.com:arobson/astro-blog`, typed by hand at the prompt, from the origin `git@github.com:arobson/astro-blog.git`.
