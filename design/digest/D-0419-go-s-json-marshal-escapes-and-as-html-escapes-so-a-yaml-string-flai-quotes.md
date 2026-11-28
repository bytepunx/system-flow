---
id: D-0419
title: "Go's json.Marshal escapes <, >, and & as HTML escapes, so a YAML string flai quotes through JSON needs an encoder with HTML escaping off"
kind: cause
sources: [S-0538, I-0222]
topics: [flai, go, code]
paths: [flai/internal/adr/adr.go, flai/internal/adr/adr_test.go]
date: 2026-10-29
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0419 Go's json.Marshal escapes <, >, and & as HTML escapes, so a YAML string flai quotes through JSON needs an encoder with HTML escaping off
