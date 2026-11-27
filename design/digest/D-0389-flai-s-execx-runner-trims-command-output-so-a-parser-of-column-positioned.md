---
id: D-0389
title: flai's execx runner trims command output, so a parser of column-positioned output such as git status --porcelain must not rely on the first line's leading space
kind: decision
sources: [S-0046]
topics: [flai, go, code]
paths: [flai/internal/execx/execx.go]
date: 2026-10-28
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0389 flai's execx runner trims command output, so a parser of column-positioned output such as git status --porcelain must not rely on the first line's leading space
