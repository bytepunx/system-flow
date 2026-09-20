---
id: I-0028
title: An acceptance from the board pushed and published with the operator's credentials, with nothing enabling it
class: defect
status: open
count: 1
cost: 30m
first_reported: 2026-09-20T13:05:49Z
last_reported: 2026-09-20T13:05:49Z
updated: 2026-09-20T13:05:49Z
---

# I-0028 An acceptance from the board pushed and published with the operator's credentials, with nothing enabling it

## Description
An acceptance from the board pushed and published with the operator's credentials, with nothing enabling it

## Instances

### 2026-09-20T13:05:49Z
2026-09-20, found while starting S-0078. S-0075 moved acceptance to flai on the host: hostapi accept.run runs flai accept without --no-push. In the container that pushed nothing, because the container held no credential; on the host it pushes main and the release tags and publishes the template with the operator's own credentials. Released in flai 1.5.3. The acceptances of S-0076 and S-0077 from the board were pushed two seconds after each click. The operators' page written in S-0075 and S-0077 said nothing pushes unasked and that the dashboard token is not the power to publish; both were false as released. ADR-0029 says host actions are off until enabled. Remedy in S-0078: accept.run always passes --no-push, and pushing is a host action of its own, off by default, enabled by name on the host, and journalled.

## Remediation
