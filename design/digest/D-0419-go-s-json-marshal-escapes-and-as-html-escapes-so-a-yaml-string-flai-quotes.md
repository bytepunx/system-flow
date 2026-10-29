---
id: D-0419
title: "Go's json.Marshal escapes <, >, and & as HTML escapes, so a YAML string flai quotes through JSON needs an encoder with HTML escaping off"
kind: cause
sources: [S-0538, I-0222]
topics: [flai, go, code]
paths: [flai/internal/adr/adr.go, flai/internal/adr/adr_test.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0419 Go's json.Marshal escapes <, >, and & as HTML escapes, so a YAML string flai quotes through JSON needs an encoder with HTML escaping off

`flai adr topics` rewrote an ADR's front matter with `<` and `>` in its title as backslash-u003c and backslash-u003e (I-0222). S-0538 confirmed the cause in `flai/internal/adr/adr.go`: `yamlString` quoted the title with `json.Marshal`, which writes `<`, `>`, and `&` as JSON's HTML escapes, and `titleOf` printed them as written.

The fix: quote through a `json.Encoder` with `SetEscapeHTML(false)`, and have `titleOf` read a double-quoted title as JSON, falling back to trimming the quotes. Any other flai writer that quotes a YAML or front-matter string with `json.Marshal` has the same fault; use the same encoder.

A related trap from the same story: literal backslash-u escapes typed in an agent's own MCP tool arguments are decoded before flai sees them, so text that must keep them spells them out.

Source: S-0538 (wip/archive/agents/S-0538.md, Decisions), I-0222.
