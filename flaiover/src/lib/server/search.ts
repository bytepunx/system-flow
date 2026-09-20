// Search and the ADR list are flai's on the host (S-0074): it indexes the project's Markdown and
// keeps the index current, and answers what the search page shows. The dashboard adds only its own
// route to each hit.
import type { Repo } from './repo';

export type Hit = {
	path: string;
	kind: 'item' | 'doc';
	itemId?: string;
	title: string;
	scope: string;
	status?: string;
	type?: string;
	nature?: string;
	score: number;
	snippet: string;
	route: string;
};

type FlaiHit = Omit<Hit, 'route'>;

/** Search; docs are excluded unless includeDocs. */
export async function search(
	repo: Repo,
	query: string,
	includeDocs = false,
	limit = 30
): Promise<{ query: string; indexed: number; hits: Hit[] }> {
	const q = query.trim().slice(0, 500);
	const res = await repo.ask<{ query: string; indexed: number; hits: FlaiHit[] | null }>(
		'search.query',
		{ q, docs: includeDocs, limit }
	);
	return {
		query: res.query,
		indexed: res.indexed,
		hits: (res.hits ?? []).map((h) => ({
			...h,
			route: h.kind === 'item' && h.itemId ? `/items/${h.itemId}` : `/docs/${h.path}`
		}))
	};
}

/** ADR front matter for the ADR list. */
export type Adr = {
	path: string;
	file: string;
	id: string;
	title: string;
	status: string;
	date: string;
	supersedes: string[];
	supersededBy: string[];
	/** ADRs this one narrows or extends without replacing them (S-0060). */
	refines: string[];
};

export async function adrs(repo: Repo): Promise<Adr[]> {
	return repo.remember<Adr[]>('adrs', 'adrs.list');
}
