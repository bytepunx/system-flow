---
id: I-0086
title: The close-out's install smoke test failed once and passed when run alone, with no cause in its output
class: efficiency
status: open
count: 1
cost: 6m
first_reported: 2026-10-06T10:22:10Z
last_reported: 2026-10-06T10:22:10Z
updated: 2026-10-06T10:31:55Z
---

# I-0086 The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Description
The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Instances

### 2026-10-06T10:22:10Z
Story: S-0283.
S-0283's close-out stopped at smoke's last sub-step, scripts/install-test.sh, after gofmt, vet, lint, the full go test -race, the template render and check, and the markdown lint passed; its output showed self-upgrade --check reaching GitHub for release 1.31.4 but no failure reason. The verifier ran scripts/install-test.sh alone once: exit 0 after over five minutes. Likely the network or a timeout; the cause is not confirmed.

## Remediation

Story S-0291 remediates this issue, created from it at 2026-10-06T10:31:55Z.
