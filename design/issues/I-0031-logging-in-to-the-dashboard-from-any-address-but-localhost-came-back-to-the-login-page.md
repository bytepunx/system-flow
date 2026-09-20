---
id: I-0031
title: Logging in to the dashboard from any address but localhost came back to the login page
class: defect
status: open
count: 1
cost: 20m
first_reported: 2026-09-20T14:06:40Z
last_reported: 2026-09-20T14:06:40Z
updated: 2026-09-20T14:06:40Z
---

# I-0031 Logging in to the dashboard from any address but localhost came back to the login page

## Description
Logging in to the dashboard from any address but localhost came back to the login page

## Instances

### 2026-09-20T14:06:40Z
2026-09-20, reported by the operator. The session cookie was marked Secure over plain HTTP, because adapter-node takes every request for https when no protocol header is configured and the image configured none; a browser keeps a Secure cookie on `http://localhost` and drops it on every other plain-HTTP address. Reproduced against the running dashboard (0.22.3) by the host's LAN address; bearer tokens were not affected. Probably present since authentication landed (S-0036) and unnoticed because the dashboard was used by localhost. Fixed in S-0083: the server entry tells the handler the scheme of each request, and the login page says when a session was not kept.

## Remediation
