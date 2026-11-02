---
id: I-0257
title: A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-02T19:06:40Z
last_reported: 2026-11-02T19:06:40Z
updated: 2026-11-02T19:06:40Z
---

# I-0257 A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host

## Description
A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host

## Instances

### 2026-11-02T19:06:40Z
Story: S-0509.
ADR-0247 mounts the main checkout read-write for the planner, the orchestrator, the analyzer, and the curator, with `.git/hooks` read-only over it. `.git/config` stays writable, because the orchestrator's `flai accept` deletes a merged branch's section from it. Keys such as `core.fsmonitor`, `core.sshCommand`, or `core.hooksPath` name commands that git runs on the host at its next command there. The guard refuses each role `git config`, but it does not see a write through the shell. That is no worse than a host run, but it is a way back out of the container. Remedy to weigh: mount `.git/config` read-only and move the branch deletion of an acceptance to the host (`flai serve`), or have `flai serve` check `.git/config` for command-naming keys after each strategic container run.

## Remediation
