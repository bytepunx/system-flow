// Repository reader: the server side of flaiover. Reads the mounted
// system-flow repository (PROJECT_DIR) and nothing else. Mirrors the rules in
// design/system/work-hierarchy.md and repository-layout.md; flai's Go
// implementation is the reference.
import { readFile, readdir, stat } from 'node:fs/promises';
import { join, relative, resolve, sep } from 'node:path';
import { EventEmitter } from 'node:events';
import { parse as parseYaml } from 'yaml';
import { watch as chokidarWatch, type FSWatcher } from 'chokidar';
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
const FOLDERS: Record<(typeof ITEM_TYPES)[number], string> = {
	epic: 'epics',
	story: 'stories',
	task: 'tasks'
};

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

/**
 * Repo reads one system-flow repository. Every read goes through a
 * per-path mtime cache; the watcher invalidates and emits 'change'.
 */
export class Repo extends EventEmitter {
	readonly root: string;
	private fileCache = new Map<string, CacheEntry<unknown>>();
	private watcher: FSWatcher | null = null;

	constructor(root = projectDir()) {
		super();
		this.root = resolve(root);
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
		return this.cached('system-flow.yaml', async (abs) => {
			const m = parseYaml(await readFile(abs, 'utf8')) as Manifest;
			if (!m || !m.layout?.design || !m.layout?.docs || !m.layout?.wip) {
				throw new RepoError(
					500,
					'system-flow.yaml is missing layout.design, layout.docs, or layout.wip'
				);
			}
			return m;
		});
	}

	async layout(): Promise<Layout> {
		return (await this.manifest()).layout;
	}

	/** All work items from kanban and archive, sorted by ID. */
	async items(): Promise<Item[]> {
		const layout = await this.layout();
		const out: Item[] = [];
		for (const archived of [false, true]) {
			for (const type of ITEM_TYPES) {
				const dir = archived
					? join(layout.wip, 'archive', 'kanban', FOLDERS[type])
					: join(layout.wip, 'kanban', FOLDERS[type]);
				const abs = this.resolveInside(dir);
				const entries = await readdir(abs).catch(() => [] as string[]);
				for (const name of entries) {
					if (!name.endsWith('.md') || name.startsWith('_')) continue;
					const rel = join(dir, name).split(sep).join('/');
					out.push(await this.item(rel, archived));
				}
			}
		}
		return out.sort((a, b) => rank(a.id) - rank(b.id) || num(a.id) - num(b.id));
	}

	private async item(rel: string, archived: boolean): Promise<Item> {
		return this.cached(rel, async (abs) => {
			const { frontMatter, body } = splitFrontMatter(await readFile(abs, 'utf8'));
			if (!frontMatter) throw new RepoError(500, `${rel}: no front matter`);
			const fm = frontMatter as Partial<Item>;
			return {
				...fm,
				id: String(fm.id),
				type: fm.type as Item['type'],
				nature: String(fm.nature),
				title: String(fm.title),
				status: String(fm.status),
				created: stamp(fm.created),
				updated: stamp(fm.updated),
				transitions: (fm.transitions ?? []).map((t) => ({ ...t, at: stamp(t.at) })),
				blocked: fm.blocked?.map((b) => ({
					...b,
					from: stamp(b.from),
					until: b.until ? stamp(b.until) : undefined
				})),
				path: rel,
				archived,
				body
			};
		});
	}

	/** Find an item by ID in any padding (S-32, S-032, S-0032 name the same item). */
	async itemById(id: string): Promise<{ item: Item; children: Item[] }> {
		const all = await this.items();
		const m = /^([EST])-?0*(\d+)$/i.exec(id.trim());
		const item = m
			? all.find(
					(it) => rank(it.id) === rank(`${m[1].toUpperCase()}-`) && num(it.id) === Number(m[2])
				)
			: undefined;
		if (!item) throw new RepoError(404, `${id} not found`);
		return { item, children: all.filter((it) => it.parent === item.id) };
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

	/** Start watching the documentation and wip folders; emits 'change' with the repo-relative path. */
	async watch(): Promise<void> {
		if (this.watcher) return;
		const layout = await this.layout();
		const targets = [layout.design, layout.docs, layout.wip, 'system-flow.yaml'].map((p) =>
			this.resolveInside(p)
		);
		this.watcher = chokidarWatch(targets, {
			ignoreInitial: true,
			awaitWriteFinish: { stabilityThreshold: 150 }
		});
		const onChange = (kind: string) => (abs: string) => {
			this.fileCache.delete(abs);
			const path = relative(this.root, abs).split(sep).join('/');
			log().debug({ component: 'watcher', event: kind, path }, 'file changed');
			this.emit('change', path);
		};
		this.watcher
			.on('add', onChange('add'))
			.on('change', onChange('change'))
			.on('unlink', onChange('unlink'))
			.on('error', (err) =>
				log().error({ component: 'watcher', err: String(err) }, 'watcher error')
			);
		log().info({ component: 'watcher', targets: targets.length }, 'watcher started');
	}

	async close(): Promise<void> {
		await this.watcher?.close();
		this.watcher = null;
	}
}

export class RepoError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
	}
}

function stamp(v: unknown): string {
	if (v instanceof Date) return v.toISOString().replace(/\.\d{3}Z$/, 'Z');
	return String(v ?? '');
}

function rank(id: string): number {
	return id.startsWith('E-') ? 0 : id.startsWith('S-') ? 1 : 2;
}
function num(id: string): number {
	return Number(id.slice(2)) || 0;
}

let shared: Repo | null = null;
/** The process-wide Repo for the mounted project. */
export function repo(): Repo {
	if (!shared) shared = new Repo();
	return shared;
}
