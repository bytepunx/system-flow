---
id: I-0314
title: TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports
class: defect
status: open
count: 1
cost: 7m
first_reported: 2026-11-08T02:50:07Z
last_reported: 2026-11-08T02:50:07Z
updated: 2026-11-08T02:50:07Z
---

# I-0314 TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports

## Description
TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports

## Instances

### 2026-11-08T02:50:07Z
Story: S-0607.
S-0607: `flai/internal/serve/role_test.go:494` failed twice with "never happened: the worked run's activity logged", once for T-2656 (`hostapi/reads.go`) and once for T-2646 (`statsread/statsread.go`). Neither change touches the formula role steps; `flai test` runs serve's tests because serve imports those packages. Each time the host was at load 50 to 90 on 24 cores, and the next run passed. Each rerun of go-test cost about 3 to 7 minutes. The test waits on a fixed time for an event that a loaded host delivers late.

## Remediation
