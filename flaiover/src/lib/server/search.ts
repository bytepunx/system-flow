// Full-text search over the repository's markdown (design and wip, docs on
// request) with MiniSearch. The index lives in the server process and is
// rebuilt, debounced, whenever the watcher reports a change.
import MiniSearch from 'minisearch';
import { readFile, readdir } from 'node:fs/promises';
import { join } from 'node:path';
import { splitFrontMatter, type Repo } from './repo';
import { log } from './log';

export type SearchDoc = {
	id: string; // path
	path: string;
	kind: 'item' | 'doc';
	itemId?: string;
	title: string;
	tags: string;
	headings: string;
	body: string;
	type?: string;
	nature?: string;
	status?: string;
	scope: 'design' | 'wip' | 'docs';
};

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

const FIELDS = ['itemId', 'title', 'tags', 'headings', 'body'];

export class SearchIndex {
	private index: MiniSearch<SearchDoc> | null = null;
	private docs = new Map<string, SearchDoc>();
	private timer: NodeJS.Timeout | null = null;
	private building: Promise<void> | null = null;

	constructor(private repo: Repo) {
		repo.on('change', () => this.schedule());
	}

	private schedule() {
		if (this.timer) clearTimeout(this.timer);
		this.timer = setTimeout(() => void this.build(), 300);
	}

	/** Build (or rebuild) the index from disk. */
	async build(): Promise<void> {
		if (this.building) return this.building;
		this.building = (async () => {
			const startedAt = process.hrtime.bigint();
			const layout = await this.repo.layout();
			const docs: SearchDoc[] = [];
			for (const [scope, dir] of [
				['design', layout.design],
				['wip', layout.wip],
				['docs', layout.docs]
			] as const) {
				await this.collect(scope, dir, docs);
			}
			const index = new MiniSearch<SearchDoc>({
				fields: FIELDS,
				storeFields: ['path', 'kind', 'itemId', 'title', 'scope', 'status', 'type'],
				searchOptions: { boost: { itemId: 4, title: 3, headings: 2 }, prefix: true, fuzzy: 0.2 }
			});
			index.addAll(docs);
			this.index = index;
			this.docs = new Map(docs.map((d) => [d.id, d]));
			log().info(
				{
					component: 'search',
					documents: docs.length,
					duration_ms: Math.round(Number(process.hrtime.bigint() - startedAt) / 1e6)
				},
				'index built'
			);
		})().finally(() => {
			this.building = null;
		});
		return this.building;
	}

	private async collect(scope: SearchDoc['scope'], rel: string, out: SearchDoc[]) {
		const abs = this.repo.resolveInside(rel);
		const entries = await readdir(abs, { withFileTypes: true }).catch(() => []);
		for (const e of entries) {
			if (e.name.startsWith('.')) continue;
			const childRel = `${rel}/${e.name}`;
			if (e.isDirectory()) {
				await this.collect(scope, childRel, out);
				continue;
			}
			if (!e.name.endsWith('.md')) continue;
			const raw = await readFile(join(abs, e.name), 'utf8');
			const { frontMatter, body } = splitFrontMatter(raw);
			const fm = (frontMatter ?? {}) as Record<string, unknown>;
			const isItem =
				typeof fm.id === 'string' && /^[EST]-\d+$/.test(fm.id) && typeof fm.type === 'string';
			const headings = [...body.matchAll(/^#{1,6}\s+(.+)$/gm)].map((m) => m[1]).join(' ');
			out.push({
				id: childRel,
				path: childRel,
				kind: isItem ? 'item' : 'doc',
				itemId: isItem ? String(fm.id) : undefined,
				title: typeof fm.title === 'string' ? fm.title : (firstHeading(body) ?? e.name),
				tags: Array.isArray(fm.tags) ? fm.tags.join(' ') : '',
				headings,
				body,
				type: typeof fm.type === 'string' ? fm.type : undefined,
				nature: typeof fm.nature === 'string' ? fm.nature : undefined,
				status: typeof fm.status === 'string' ? fm.status : undefined,
				scope
			});
		}
	}

	/** Search; docs are excluded unless includeDocs. */
	async search(query: string, includeDocs = false, limit = 30): Promise<Hit[]> {
		if (!this.index) await this.build();
		const q = query.trim();
		if (!q) return [];
		return this.index!.search(q)
			.filter((r) => includeDocs || r.scope !== 'docs')
			.slice(0, limit)
			.map((r) => {
				const doc = this.docs.get(r.id)!;
				return {
					path: doc.path,
					kind: doc.kind,
					itemId: doc.itemId,
					title: doc.title,
					scope: doc.scope,
					status: doc.status,
					type: doc.type,
					nature: doc.nature,
					score: Math.round(r.score * 100) / 100,
					snippet: snippet(doc.body, r.terms),
					route: doc.kind === 'item' && doc.itemId ? `/items/${doc.itemId}` : `/docs/${doc.path}`
				};
			});
	}

	size(): number {
		return this.docs.size;
	}
}

function firstHeading(body: string): string | undefined {
	return body.match(/^#\s+(.+)$/m)?.[1];
}

/** A short window of the body around the first matched term. */
export function snippet(body: string, terms: string[], width = 160): string {
	const text = body.replace(/\s+/g, ' ').trim();
	const lower = text.toLowerCase();
	let at = -1;
	for (const t of terms) {
		const i = lower.indexOf(t.toLowerCase());
		if (i >= 0 && (at < 0 || i < at)) at = i;
	}
	if (at < 0) return text.slice(0, width) + (text.length > width ? '…' : '');
	const start = Math.max(0, at - Math.floor(width / 3));
	const end = Math.min(text.length, start + width);
	return (start > 0 ? '…' : '') + text.slice(start, end) + (end < text.length ? '…' : '');
}

let shared: SearchIndex | null = null;
export function searchIndex(repo: Repo): SearchIndex {
	if (!shared) shared = new SearchIndex(repo);
	return shared;
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
};

export async function adrs(repo: Repo): Promise<Adr[]> {
	const layout = await repo.layout();
	const dir = `${layout.design}/adrs`;
	const abs = repo.resolveInside(dir);
	const out: Adr[] = [];
	for (const name of (await readdir(abs).catch(() => [])).sort()) {
		if (!/^\d{4}-.*\.md$/.test(name) || name.startsWith('0000-')) continue;
		const { frontMatter } = splitFrontMatter(await readFile(join(abs, name), 'utf8'));
		const fm = (frontMatter ?? {}) as Record<string, unknown>;
		const list = (v: unknown) => (Array.isArray(v) ? v.map(String) : []);
		out.push({
			path: `${dir}/${name}`,
			file: name,
			id: String(fm.id ?? name.slice(0, 4)),
			title: String(fm.title ?? name),
			status: String(fm.status ?? ''),
			date: fm.date instanceof Date ? fm.date.toISOString().slice(0, 10) : String(fm.date ?? ''),
			supersedes: list(fm.supersedes),
			supersededBy: list(fm.superseded_by)
		});
	}
	return out;
}
