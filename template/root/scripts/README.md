# scripts

Purpose-named shell scripts for common tasks. The `Makefile` calls these; CI calls the same ones. Each script is safe to run from any directory, prints what it does, and exits non-zero on failure.

| Script | Does |
|--------|------|
| `check.sh` | Runs `flai check --strict` on this repository |
| `test.sh` | Behavior tests, fast, every iteration (`make test`) |
| `integration.sh` | Integration tests against real adapters (`make integration`) |
| `smoke.sh` | End-to-end smoke tests (`make smoke`) |
