// Metrics come from `flai stats --json` (ADR-0016): one reference
// implementation of design/system/metrics.md. Results are cached per
// argument set and dropped when the watcher reports a change.
import { flai } from './flai';
import type { Repo } from './repo';

export type StatsQuery = { since?: string; type?: string; by?: string };

const WINDOW = /^\d+[dwh]$/;
const TYPES = new Set(['epic', 'story', 'task']);
const BY = new Set(['nature', 'type', 'parent']);

export function validate(q: StatsQuery): string[] {
	const args = ['stats'];
	if (q.since) {
		if (!WINDOW.test(q.since)) throw new Error('since must look like 30d, 12w, or 720h');
		args.push('--since', q.since);
	}
	if (q.type) {
		if (!TYPES.has(q.type)) throw new Error('type must be epic, story, or task');
		args.push('--type', q.type);
	}
	if (q.by) {
		if (!BY.has(q.by)) throw new Error('by must be nature, type, or parent');
		args.push('--by', q.by);
	}
	return args;
}

const cache = new Map<string, Promise<unknown>>();
const hooked = new WeakSet<Repo>();

export async function stats<T = unknown>(repo: Repo, q: StatsQuery): Promise<T> {
	if (!hooked.has(repo)) {
		hooked.add(repo);
		repo.on('change', () => cache.clear());
	}
	const args = validate(q);
	const key = args.join(' ');
	let p = cache.get(key) as Promise<T> | undefined;
	if (!p) {
		p = flai<T>(repo.root, args).then((r) => r.data);
		cache.set(key, p);
		p.catch(() => cache.delete(key));
	}
	return p;
}

/** For tests. */
export function clearStatsCache() {
	cache.clear();
}
