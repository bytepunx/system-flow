---
id: T-1025
type: task
nature: improvement
title: The manifest carries the shared paths as glob patterns, validated, with a matcher and an edit function the CLI, HTTP, and MCP share
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:15:09Z
updated: 2026-10-06T12:15:09Z
transitions: []
stream: S-0295
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/shared.go, flai/internal/manifest/shared_test.go, flai/go.mod, flai/go.sum, design/tech/go-libraries.md, system-flow.yaml, template/root/system-flow.yaml.tmpl]
after: [T-1023]
---
# T-1025 The manifest carries the shared paths as glob patterns, validated, with a matcher and an edit function the CLI, HTTP, and MCP share

## Work

Add the key T-1023's ADR names (proposed `claims.shared`) to the manifest in `flai/internal/manifest/manifest.go`: a list of glob patterns. In `flai/internal/manifest/shared.go`, add:

- **Validation.** It runs at load, beside the other `Errors()` checks, and rejects an empty pattern, an absolute path, `..`, and a malformed glob.
- **A matcher.** It says whether a path, or a touches entry that is a folder, lies wholly inside a shared pattern, and which pattern matched. Use the ADR's dialect: `*` within one segment, `**` across segments, and a plain path covering itself and everything below it.
- **Add and remove.** These rewrite the list in `system-flow.yaml`, keep the file's other keys and comments, refuse a duplicate or an invalid pattern, and say what changed. The CLI, the HTTP API, and the MCP tools all call them.

Either take a glob library such as `github.com/bmatcuk/doublestar/v4` and record it in `design/tech/go-libraries.md`, or write the matcher on `path.Match` with `**` handled by hand. Say which in the narrative's `## Decisions`.

Set this project's list in `system-flow.yaml` and the template's default in `template/root/system-flow.yaml.tmpl`, as the ADR records them. This task waits for T-1023 for the key and the dialect.

## Done when

- Table tests cover the matcher: `design/adrs` and `design/adrs/**` both match `design/adrs/0099-x.md`, `docs/users/*.md` matches `docs/users/flai.md` but not `docs/users/a/b.md`, and a folder touch counts as inside only when the whole folder is.
- Tests show validation rejecting each bad pattern, and add and remove keeping the rest of a manifest byte for byte.
- `scripts/flai-test.sh` passes and `flai check --strict` is clean with the new key in `system-flow.yaml`.

## Notes
