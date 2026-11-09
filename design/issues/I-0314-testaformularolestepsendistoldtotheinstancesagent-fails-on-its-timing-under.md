---
id: I-0314
title: TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports
class: defect
status: open
count: 2
cost: 6m
first_reported: 2026-11-08T02:50:07Z
last_reported: 2026-11-09T23:20:40Z
updated: 2026-11-09T23:20:40Z
---

# I-0314 TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports

## Description
TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports

## Instances

### 2026-11-08T02:50:07Z
Story: S-0607.
S-0607: `flai/internal/serve/role_test.go:494` failed twice with "never happened: the worked run's activity logged", once for T-2656 (`hostapi/reads.go`) and once for T-2646 (`statsread/statsread.go`). Neither change touches the formula role steps; `flai test` runs serve's tests because serve imports those packages. Each time the host was at load 50 to 90 on 24 cores, and the next run passed. Each rerun of go-test cost about 3 to 7 minutes. The test waits on a fixed time for an event that a loaded host delivers late.

### 2026-11-09T23:20:40Z
Story: S-0633.
On 2026-11-09, T-2786's flai test runs for flai/internal/itemnew and flai/internal/workitem each failed once in TestAFormulaRoleStepsEndIsToldToTheInstancesAgent, with a different message each time. Two other stories' agents and their go-test runs were loading the host. The test passed in the serve run before and after.

## Remediation
