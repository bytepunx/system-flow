// Repository reader: the server side of flaiover. Reads the mounted
// system-flow repository (PROJECT_DIR) and nothing else. Mirrors the rules in
// design/system/work-hierarchy.md and repository-layout.md; flai's Go
// implementation is the reference.
import { readFile, readdir, stat } from 'node:fs/promises';
import { join, resolve, sep } from 'node:path';
import { EventEmitter } from 'node:events';
import { parse as parseYaml } from 'yaml';
import { agent, AgentError } from './agent';
import { log } from './log';

export type Layout = { design: string; docs: string; wip: string };

export type Manifest = {
	version: number;
	name: string;
	key?: string;
	description?: string;
	owner?: string;
	repo?: string;
	template?: { repo?: string; ref?: string; version?: string; applied?: string };
	layout: Layout;
	projects?: { name: string; path: string; kind: string; tags?: string[] }[];
	dashboard?: { image?: string; tag?: string; port?: number };
};

export type Transition = { to: string; at: string; by: string };
export type Block = { from: string; until?: string; reason: string };

export type ThreadEntry = { at: string; author: string; text: string };
export type Thread = {
	id: string;
	title: string;
	anchor: { path: string; heading?: string; item?: string };
	status: 'open' | 'answered' | 'resolved';
	participants: string[];
	created: string;
	updated: string;
	path: string; // repo-relative
	entries: ThreadEntry[];
};

export type Item = {
	id: string;
	type: 'epic' | 'story' | 'task';
	nature: string;
	title: string;
	status: string;
	parent?: string;
	owner?: string;
	created: string;
	updated: string;
	transitions: Transition[];
	blocked?: Block[];
	estimate?: string;
	stream?: string;
	tags?: string[];
	touches?: string[]; // paths or components the work changes (ADR-0019)
	path: string; // repo-relative
	archived: boolean;
	body: string;
};

export type DocNode = {
	name: string;
	path: string; // repo-relative, forward slashes
	kind: 'dir' | 'file';
	title?: string;
	frontMatter?: Record<string, unknown>;
	children?: DocNode[];
};

export const ITEM_TYPES = ['epic', 'story', 'task'] as const;

export function projectDir(): string {
	return resolve(process.env.PROJECT_DIR ?? process.cwd());
}

/** Split a markdown document into front matter (parsed) and body. */
export function splitFrontMatter(doc: string): {
	frontMatter: Record<string, unknown> | null;
	body: string;
} {
	if (!doc.startsWith('---\n')) return { frontMatter: null, body: doc };
	const rest = doc.slice(4);
	const end = rest.indexOf('\n---\n');
	if (end < 0) return { frontMatter: null, body: doc };
	const fm = parseYaml(rest.slice(0, end + 1)) as Record<string, unknown> | null;
	return { frontMatter: fm ?? {}, body: rest.slice(end + 5) };
}

type CacheEntry<T> = { mtimeMs: number; value: T };

/** Asks flai on the host for a named method (ADR-0029). Tests supply one that reads a fixture through flai. */
export type Ask = <T>(method: string, params?: Record<string, unknown>) => Promise<T>;

/** flai's code for an item or thread that does not exist (internal/hostapi). */
const NOT_FOUND = -32004;
const INVALID_PARAMS = -32602;

const viaChannel: Ask = (method, params) => agent().ask(method, params);

/**
 * Repo is the dashboard's view of one system-flow repository. The project, its work items, and its
 * threads are asked of flai on the host over the channel and kept until flai says a file changed
 * (S-0073); documents, narratives, and search still read the mount through a per-path mtime cache
 * until S-0074 moves them. Either way a change is announced as 'change' with the repo-relative path.
 */
export class Repo extends EventEmitter {
	readonly root: string;
	private fileCache = new Map<string, CacheEntry<unknown>>();
	private answers = new Map<string, Promise<unknown>>();
	private listening = false;

	constructor(
		root = projectDir(),
		private source: Ask = viaChannel
	) {
		super();
		this.root = resolve(root);
	}

	/** Ask flai, with its refusals as the HTTP statuses the routes answer with. */
	async ask<T>(method: string, params: Record<string, unknown> = {}): Promise<T> {
		try {
			return await this.source<T>(method, params);
		} catch (e) {
			if (e instanceof AgentError) {
				const status = e.code === NOT_FOUND ? 404 : e.code === INVALID_PARAMS ? 400 : e.status;
				throw new RepoError(status, e.message);
			}
			throw e;
		}
	}

	/** One answer per question until something changes; a failure is not kept. */
	private remembered<T>(key: string, method: string, params: Record<string, unknown> = {}) {
		let hit = this.answers.get(key) as Promise<T> | undefined;
		if (!hit) {
			hit = this.ask<T>(method, params);
			this.answers.set(key, hit);
			hit.catch(() => this.answers.delete(key));
		}
		return hit;
	}

	/** A file of the project changed, by flai's word: forget what was asked and tell the listeners. */
	changed(path: string): void {
		this.answers.clear();
		try {
			this.fileCache.delete(this.resolveInside(path));
		} catch {
			// not a path inside the project: nothing cached under it
		}
		log().debug({ component: 'watcher', path }, 'file changed');
		this.emit('change', path);
	}

	/** Resolve a repo-relative path and refuse anything outside the root. */
	resolveInside(rel: string): string {
		const abs = resolve(this.root, rel);
		if (abs !== this.root && !abs.startsWith(this.root + sep)) {
			throw new RepoError(400, `path escapes the repository: ${rel}`);
		}
		return abs;
	}

	private async cached<T>(rel: string, load: (abs: string) => Promise<T>): Promise<T> {
		const abs = this.resolveInside(rel);
		const st = await stat(abs).catch(() => null);
		if (!st) throw new RepoError(404, `not found: ${rel}`);
		const hit = this.fileCache.get(abs) as CacheEntry<T> | undefined;
		if (hit && hit.mtimeMs === st.mtimeMs) return hit.value;
		const value = await load(abs);
		this.fileCache.set(abs, { mtimeMs: st.mtimeMs, value });
		return value;
	}

	async manifest(): Promise<Manifest> {
		const m = await this.remembered<Manifest>('project', 'project.info');
		if (!m || !m.layout?.design || !m.layout?.docs || !m.layout?.wip) {
			throw new RepoError(
				500,
				'system-flow.yaml is missing layout.design, layout.docs, or layout.wip'
			);
		}
		return m;
	}

	async layout(): Promise<Layout> {
		return (await this.manifest()).layout;
	}

	/** Every thread, resolved ones included, sorted by ID (ADR-0020), as flai reads them. */
	async threads(): Promise<Thread[]> {
		return this.remembered<Thread[]>('threads', 'threads.list', { all: true });
	}

	/** Threads anchored to a repository path or an item ID. */
	async threadsFor(on: string): Promise<Thread[]> {
		const want = on.replace(/\/$/, '');
		return this.remembered<Thread[]>(`threads:${want}`, 'threads.list', { on: want, all: true });
	}

	/** All work items from kanban and archive, sorted by ID, as flai reads them. */
	async items(): Promise<Item[]> {
		const raw = await this.remembered<FlaiItem[]>('items', 'items.list', {
			archived: true,
			bodies: true
		});
		return raw.map(fromFlai);
	}

	/** Find an item by ID in any padding (S-32, S-032, S-0032 name the same item). */
	async itemById(id: string): Promise<{ item: Item; children: Item[] }> {
		const got = await this.remembered<{ item: FlaiItem; children: FlaiItem[] | null }>(
			`item:${id.trim()}`,
			'item.get',
			{ id: id.trim() }
		);
		return { item: fromFlai(got.item), children: (got.children ?? []).map(fromFlai) };
	}

	/** Documentation trees: design (all types), docs, and wip. */
	async docsTree(): Promise<DocNode[]> {
		const layout = await this.layout();
		const roots: DocNode[] = [];
		for (const dir of [layout.design, layout.docs, layout.wip]) {
			roots.push(await this.tree(dir));
		}
		return roots;
	}

	private async tree(rel: string): Promise<DocNode> {
		const abs = this.resolveInside(rel);
		const node: DocNode = {
			name: rel.split('/').pop() ?? rel,
			path: rel,
			kind: 'dir',
			children: []
		};
		const entries = (await readdir(abs, { withFileTypes: true }).catch(() => [])).sort((a, b) =>
			a.name.localeCompare(b.name)
		);
		for (const e of entries) {
			if (e.name.startsWith('.')) continue;
			const childRel = `${rel}/${e.name}`;
			if (e.isDirectory()) {
				node.children!.push(await this.tree(childRel));
			} else if (e.name.endsWith('.md')) {
				const { frontMatter } = splitFrontMatter(await readFile(join(abs, e.name), 'utf8'));
				node.children!.push({
					name: e.name,
					path: childRel,
					kind: 'file',
					title: typeof frontMatter?.title === 'string' ? frontMatter.title : undefined,
					frontMatter: frontMatter ?? undefined
				});
			}
		}
		return node;
	}

	/** One markdown file, raw, with parsed front matter. */
	async docFile(rel: string): Promise<{
		path: string;
		frontMatter: Record<string, unknown> | null;
		body: string;
		raw: string;
	}> {
		if (!rel.endsWith('.md')) throw new RepoError(400, 'only markdown files are served');
		return this.cached(rel, async (abs) => {
			const raw = await readFile(abs, 'utf8');
			const { frontMatter, body } = splitFrontMatter(raw);
			return { path: rel, frontMatter, body, raw };
		});
	}

	/**
	 * Listen for flai's word that a file changed, and for its return: a flai that comes back may have
	 * missed changes, so everything asked is forgotten and open pages are told to look again.
	 */
	async watch(): Promise<void> {
		if (this.listening || this.source !== viaChannel) return;
		this.listening = true;
		const hub = agent();
		hub.on('change', (path: string) => this.changed(path));
		hub.on('connected', () => this.changed('system-flow.yaml'));
		log().info({ component: 'watcher' }, 'listening for changes from the host flai');
	}

	async close(): Promise<void> {
		this.answers.clear();
	}
}

export class RepoError extends Error {
	constructor(
		public status: number,
		message: string,
		/** Extra fields for the response body, beside `error` (a conflict's current version, a refusal's findings). */
		public data?: Record<string, unknown>
	) {
		super(message);
	}
}

/** An item as flai marshals it: empty strings and nulls where the front matter had nothing. */
type FlaiItem = Omit<Item, 'blocked' | 'tags' | 'touches' | 'transitions'> & {
	transitions: Transition[] | null;
	blocked: (Omit<Block, 'until'> & { until?: string })[] | null;
	tags: string[] | null;
	touches?: string[] | null;
};

function fromFlai(it: FlaiItem): Item {
	return {
		...it,
		parent: it.parent || undefined,
		owner: it.owner || undefined,
		estimate: it.estimate || undefined,
		stream: it.stream || undefined,
		transitions: it.transitions ?? [],
		blocked: it.blocked?.length
			? it.blocked.map((b) => ({ ...b, until: b.until || undefined }))
			: undefined,
		tags: it.tags ?? undefined,
		touches: it.touches ?? undefined
	};
}

let shared: Repo | null = null;
/** The process-wide Repo for the mounted project. */
export function repo(): Repo {
	if (!shared) shared = new Repo();
	return shared;
}

/** Replace the process-wide Repo, or forget it with null. For tests, which ask a flai of their own. */
export function useRepo(r: Repo | null): void {
	shared = r;
}
