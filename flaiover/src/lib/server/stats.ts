// Metrics come from flai's stats on the host (ADR-0016, S-0075): one reference implementation of
// design/system/metrics.md. flai validates the window, type, and grouping; the answer is kept per
// question until flai says a file changed, like everything else the Repo is told. The bucket is
// what spend over time is laid out in (S-0163); a flai older than that ignores it and sends no
// spend, which the charts say.
import type { Repo } from './repo';

export type StatsQuery = { since?: string; type?: string; by?: string; bucket?: string };

export async function stats<T = unknown>(repo: Repo, q: StatsQuery): Promise<T> {
	const params = {
		since: q.since || undefined,
		type: q.type || undefined,
		by: q.by || undefined,
		bucket: q.bucket || undefined
	};
	const got = await repo.remember<{ data: T }>(
		`stats:${JSON.stringify(params)}`,
		'stats.get',
		params
	);
	return got.data;
}
