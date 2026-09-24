---
id: I-0040
title: flai upgrade renders with the five standard variables only, so a template variable a fork adds fails every upgrade
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-09-24T08:16:44Z
last_reported: 2026-09-24T08:16:44Z
updated: 2026-09-24T08:16:44Z
---

# I-0040 flai upgrade renders with the five standard variables only, so a template variable a fork adds fails every upgrade

## Description
flai upgrade renders with the five standard variables only, so a template variable a fork adds fails every upgrade

## Instances

### 2026-09-24T08:16:44Z
S-0018: a fork of template/ adding team_channel, used in root/docs/contact.md.tmpl, rendered with flai new; flai upgrade --template to it failed with map has no entry for key team_channel. cmd/upgrade.go manifestVars supplies project_name, project_key, description, owner, repo_url from system-flow.yaml and nothing else; documented as a limitation in docs/contributors/template.md.

## Remediation
