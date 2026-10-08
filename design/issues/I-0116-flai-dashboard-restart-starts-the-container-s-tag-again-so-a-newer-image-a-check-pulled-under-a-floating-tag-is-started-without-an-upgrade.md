---
id: I-0116
title: flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-10-08T00:09:35Z
last_reported: 2026-10-08T00:09:35Z
updated: 2026-10-08T09:39:22Z
---

# I-0116 flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade

## Description
flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade

## Instances

### 2026-10-08T00:09:35Z
Story: S-0316.
Found reviewing T-1295. runDashboardRestart in flai/cmd/dashboard_upgrade.go restarts from containerInfo's {{.Config.Image}}, the tag the container was started with, such as ghcr.io/bytepunx/flaiover:latest. flai dashboard check pulls that tag, so once a newer latest is pulled a restart starts the newer image, while restart's Long text says a floating tag cannot silently upgrade it. Restarting from {{.Image}}, the image ID, as S-0316's upgrade fallback does, would keep it. Not reproduced on the host; read from the code.

## Remediation

Story S-0344 remediates this issue, created from it at 2026-10-08T08:08:17Z.
Closed 2026-10-08T09:39:22Z: S-0344: flai dashboard restart (T-1360) and flai host's watch (T-1361) start the image ID the container runs, never its tag, so a newer image flai dashboard check pulled under a floating tag is not started without an upgrade; each container carries the io.bytepunx.flai.ref label so the tag is still shown. Tests in flai/cmd/dashboard_test.go and dashboard_watch_test.go reproduce it.
