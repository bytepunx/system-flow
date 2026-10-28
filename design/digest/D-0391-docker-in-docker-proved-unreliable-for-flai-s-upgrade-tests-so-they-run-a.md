---
id: D-0391
title: Docker-in-Docker proved unreliable for flai's upgrade tests, so they run a throwaway flai binary against the host's real Docker daemon
kind: decision
sources: [S-0081]
topics: [testing, docker, host]
paths: [design/tech/docker.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0391 Docker-in-Docker proved unreliable for flai's upgrade tests, so they run a throwaway flai binary against the host's real Docker daemon

S-0081 (T-0325) first isolated its upgrade tests in a nested Docker daemon, and switched because that path "proved unreliable". It ran a throwaway flai binary against the host's real daemon instead, and diffed the worktree source against the patched scratch copy before each use, so the code tested stayed the code that ships.

A later test that needs a real daemon follows the same pattern, and never touches the operator's own containers.

Source: S-0081.
