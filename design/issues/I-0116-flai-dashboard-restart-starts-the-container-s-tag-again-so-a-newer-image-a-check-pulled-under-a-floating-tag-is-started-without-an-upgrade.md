---
id: I-0116
title: flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-08T00:09:35Z
last_reported: 2026-10-08T00:09:35Z
updated: 2026-10-08T00:09:35Z
---

# I-0116 flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade

## Description
flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade

## Instances

### 2026-10-08T00:09:35Z
Story: S-0316.
Found reviewing T-1295. runDashboardRestart in flai/cmd/dashboard_upgrade.go restarts from containerInfo's {{.Config.Image}}, the tag the container was started with, such as ghcr.io/bytepunx/flaiover:latest. flai dashboard check pulls that tag, so once a newer latest is pulled a restart starts the newer image, while restart's Long text says a floating tag cannot silently upgrade it. Restarting from {{.Image}}, the image ID, as S-0316's upgrade fallback does, would keep it. Not reproduced on the host; read from the code.

## Remediation
