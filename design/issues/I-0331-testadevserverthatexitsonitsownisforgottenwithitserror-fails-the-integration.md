---
id: I-0331
title: "TestADevServerThatExitsOnItsOwnIsForgottenWithItsError fails the integration tier under host load: the crash-later fake vite crashes before flai's probe sees it answer"
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-10T08:15:13Z
last_reported: 2026-11-10T08:15:13Z
updated: 2026-11-10T08:20:49Z
---

# I-0331 TestADevServerThatExitsOnItsOwnIsForgottenWithItsError fails the integration tier under host load: the crash-later fake vite crashes before flai's probe sees it answer

## Description
TestADevServerThatExitsOnItsOwnIsForgottenWithItsError fails the integration tier under host load: the crash-later fake vite crashes before flai's probe sees it answer

## Instances

### 2026-11-10T08:15:13Z
Story: S-0652.
S-0652's close-out (job S-0652-20261110T081025Z, head 697a99ef2) stopped at integration on flai/internal/serve/components_test.go:293: 'the Vite dev server ... exited before it answered ... exit status 2 ... crashed', at a host load of about 71. The crash-later fake exits 300ms after it starts, whether or not StartDevServer's probe has seen it answer 200, so under load the crash wins and Start fails. flai/internal/components/devserver_test.go's crash-later fake, used by TestErrSaysWhyAServerExitedOnItsOwn, has the same race.

## Remediation
Closed 2026-11-10T08:20:49Z: Fixed by S-0652's T-2859: the crash-later fake vites in flai/internal/serve/components_test.go and flai/internal/components/devserver_test.go now crash 300ms after they first answer, not 300ms after they start, so flai's probe has always seen them answer. TestADevServerThatExitsOnItsOwnIsForgottenWithItsError holds its fake back 500ms from listening, which failed it on every run with the old fake.
