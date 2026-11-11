---
id: I-0258
title: A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-11-02T19:06:43Z
last_reported: 2026-11-02T19:06:43Z
updated: 2026-11-11T16:04:11Z
---

# I-0258 A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot

## Description
A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot

## Instances

### 2026-11-02T19:06:43Z
Story: S-0509.
S-0509 starts the four built-in roles in a container. pilot reads their role file from the mounted checkout with `roles.Builtin`. A role a project defines arrives with `harness.Request.RoleFile` set, and that field is `json:"-"`, as are its `Inputs`. The container adapter therefore refuses it (`containerSubject`) and names another harness. Remedy: have pilot load the role file by name from the mounted checkout (`design/workflows/roles/<name>.md`), as it does a built-in one, and carry the inputs in the start request.

## Remediation

Story S-0705 remediates this issue, created from it at 2026-11-10T09:11:33Z.
Closed 2026-11-11T16:04:11Z: S-0705: pilot's start request now carries a defined role's name (role_name) and its inputs, never its file; pilot reads the file from the mounted checkout's design/workflows/roles and starts the inner harness as on the host, the container adapter no longer refuses it, its container and cache folder are named for the role's name, and flai serve passes the start request on standard input for a defined role's start, which it had dropped.
