---
id: I-0222
title: "flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >"
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-10-28T21:41:05Z
last_reported: 2026-10-28T21:41:05Z
updated: 2026-10-29T03:09:30Z
---

# I-0222 flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >

## Description
flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >

## Instances

### 2026-10-28T21:41:05Z
Story: S-0416.
S-0416 (T-1732): an ADR recorded with flai adr new whose title held `<name>-<item>` came out of flai adr topics with the angle brackets written as < and > in its front matter. The task sub-agent deleted the uncommitted file and its index row and recorded the ADR again with a title without angle brackets (ADR-0210).

## Remediation

Story S-0538 remediates this issue, created from it at 2026-10-29T02:01:36Z.
Closed 2026-10-29T03:09:30Z: Fixed by S-0538: flai adr new quoted a title with json.Marshal, which writes <, >, and & as <, >, and &; flai adr topics only showed them, since titleOf printed the title as written. yamlString now quotes with an encoder whose HTML escaping is off, and titleOf reads a double-quoted title as YAML does. TestYamlStringWritesNoHTMLEscapes in flai/internal/adr/adr_test.go and TestAdrNewWritesAngleBracketsAndAmpersandsAsThemselves in flai/cmd/adr_test.go reproduce it. No ADR on main held such an escape, so no file needed repair.
