---
id: I-0192
title: component-library.md still says the main checkout is never installed into, though flai accept now installs there by ADR-0185
class: defect
status: closed
count: 1
first_reported: 2026-10-22T20:48:39Z
last_reported: 2026-10-22T20:48:39Z
updated: 2026-10-28T13:06:57Z
---

# I-0192 component-library.md still says the main checkout is never installed into, though flai accept now installs there by ADR-0185

## Description
component-library.md still says the main checkout is never installed into, though flai accept now installs there by ADR-0185

## Instances

### 2026-10-22T20:48:39Z
S-0483 (accepted 2026-10-22T20:48:17Z) made flai accept run a sub-project's install command in the main checkout when the accepted change touched its install paths (ADR-0185). On TH-0578 the operator asked that the rule it breaks be updated with the exception and its reason, but design/system/component-library.md, The install step, still ends: The main checkout is never installed into (TH-0589). The living design now contradicts the shipped behaviour. The fix is that sentence rewritten to name flai accept's install as the one exception, with I-0179 as the reason and a link to ADR-0185. ADR-0185 also states no timeout for the install, which runs under context.Background().

## Remediation

Story S-0526 remediates this issue, created from it at 2026-10-28T12:54:04Z.
Closed 2026-10-28T13:06:57Z: Fixed by S-0526 (T-2204, commit 87cc7f482): design/system/component-library.md § Sources › The install step no longer says the main checkout is never installed into. It says that step never installs there and names flai accept's install as the one install in the main checkout: the manifest's install command, run after the acceptance commit when the accepted change touched its install.paths, with no time limit, a failed install warning without failing the acceptance, with I-0179 as the reason and a link to ADR-0185. No test fits: the defect was one sentence of the living design.
