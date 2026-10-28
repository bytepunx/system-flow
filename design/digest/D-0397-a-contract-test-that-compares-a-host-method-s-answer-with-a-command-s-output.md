---
id: D-0397
title: A contract test that compares a host method's answer with a command's output runs the command through cobra in the test's own process with a fixed clock, not a built flai
kind: decision
sources: [S-0159]
topics: [testing, flai, go]
paths: [flai/cmd]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0397 A contract test that compares a host method's answer with a command's output runs the command through cobra in the test's own process with a fixed clock, not a built flai

S-0159 put the contract test in `flai/cmd` and ran "each command through cobra in the test's process with a fixed clock, rather than starting a built flai: a separate process cannot share the clock, and `stats` answers from it." A later host method whose answer must match a command's, especially one that reads the time, is tested the same way.

Source: S-0159.
