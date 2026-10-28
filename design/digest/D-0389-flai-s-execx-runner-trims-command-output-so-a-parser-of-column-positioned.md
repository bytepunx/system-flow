---
id: D-0389
title: flai's execx runner trims command output, so a parser of column-positioned output such as git status --porcelain must not rely on the first line's leading space
kind: decision
sources: [S-0046]
topics: [flai, go, code]
paths: [flai/internal/execx/execx.go]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0389 flai's execx runner trims command output, so a parser of column-positioned output such as git status --porcelain must not rely on the first line's leading space

`flai/internal/execx/execx.go` runs `strings.TrimSpace` on a command's output. The dashboard's dirty-tree check parsed `git status --porcelain` by column, so a first line that began with a space lost a character, and a `wip/` path read as a foreign change. S-0046's narrative says that would have refused every acceptance from the dashboard.

S-0046 replaced it with a parser that reads each line's shape, with a test. A later parser of positional command output reads the trimmed text, or is tested through the real runner.

Source: S-0046.
