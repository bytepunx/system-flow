# Runbooks

Step by step, for an operator at a shell on the host. Each runbook has a section for flai (the command-line tool, and `flai host` and `flai serve`, which run on the host) and one for flaiover (the dashboard container). Settings are in the [settings index](../settings.md); what each part does is in the [operators guide](../index.md).

| Runbook | When |
|---------|------|
| [install.md](install.md) | A new machine, or the first project on one |
| [update.md](update.md) | A newer flai or dashboard is published |
| [backup.md](backup.md) | Before a risky change, or on a schedule |
| [restore.md](restore.md) | After losing the host's state, or a project's clone |
| [migrate.md](migrate.md) | Moving to another machine, bringing a project to a newer template, or crossing a release that changed how things are kept |
| [delete.md](delete.md) | Taking flai or the dashboard off a machine |
