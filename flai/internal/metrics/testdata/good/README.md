# metrics fixture

`now` for tests is 2026-09-01T12:00:00Z with a 30 day window (starts 2026-08-02T12:00:00Z).

| Story | Created | Ready | In progress | Review | Done/Cancelled | Notes |
|-------|---------|-------|-------------|--------|----------------|-------|
| S-001 | 08-01 09:00 | 08-01 10:00 | 08-02 10:00 | 08-03 10:00 | done 08-03 12:00 | blocked 08-02 12:00 to 18:00 (6h), estimate 20h |
| S-002 | 08-10 09:00 | 08-10 09:30 | 08-11 09:30 | 08-11 19:30 | done 08-12 09:30 | |
| S-003 | 08-20 09:00 | | | | cancelled 08-21 09:00 | |
| S-004 | 08-25 09:00 | 08-25 10:00 | 08-31 10:00 | | | active, age 26h |
| S-005 | 07-01 09:00 | 07-01 10:00 | 07-02 10:00 | 07-03 10:00 | done 07-05 10:00 | outside window |

Expected (stories): completed 2, cancelled 1, cancellation rate 1/3, throughput 2/(30/7)=0.4667/week, WIP 1.
Cycle seconds [93600, 86400]: p50 86400, p85 93600, max 93600, mean 90000.
Lead seconds [183600, 174600]: p50 174600, p85 183600, mean 179100. Queue: both 86400.
Flow efficiency mean(72000/93600, 1) = 0.884615. S-001 estimate error (93600-72000)/72000 = 0.3.
Time in state totals over lead 358200: backlog 5400, ready 172800, in-progress 122400, review 57600.
Burn-up last day: scope 4 (S-003 excluded), done 3. CFD last day: in-progress 1, done 3, cancelled 1.
