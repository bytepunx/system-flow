// The project as the server sees it: nothing on disk. Every answer is asked of flai on the host
// over the channel (ADR-0029), and the container holds no file of the project (ADR-0031). The
// shapes mirror design/system/work-hierarchy.md and repository-layout.md; flai's Go
// implementation is the reference.
import { resolve } from 'node:path';
import { EventEmitter } from 'node:events';
import { randomUUID } from 'node:crypto';
import {
	agentFor,
	AgentError,
	concerns,
	connectedWithin,
	currentProjectKey,
	defaultProjectKey,
	registry
} from './agent';
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
	/** Host actions and whether the operator enabled each for this project, on the host (S-0078). */
	host_actions?: Record<string, boolean>;
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

/**
 * Names the project for tests and development, where flai is asked about a directory on this
 * machine (testing.ts). In the container it names nothing that exists there.
 */
export function projectDir(): string {
	return resolve(process.env.PROJECT_DIR ?? process.cwd());
}

/** Asks flai on the host for a named method (ADR-0029). Tests supply one that reads a fixture through flai. */
export type AskOptions = { timeoutMs?: number; onProgress?: (value: unknown) => void };
export type Ask = <T>(
	method: string,
	params?: Record<string, unknown>,
	opt?: AskOptions
) => Promise<T>;

/** What a command of flai answered: its JSON, and the warnings it logged. */
export type Written<T> = { data: T; warnings: string[] };

/** flai's code for an item or thread that does not exist (internal/hostapi). */
const NOT_FOUND = -32004;
const INVALID_PARAMS = -32602;
const CONFLICT = -32009;
const REFUSED = -32010;
const RULE = -32011;
const DISABLED = -32012; // a host action the operator has not enabled (S-0078)
const NOT_COMMITTED = -32013; // an import applied but not committed: a test failed (S-0098)
const STATUS: Record<number, number> = {
	[NOT_FOUND]: 404,
	[INVALID_PARAMS]: 400,
	[RULE]: 400,
	[CONFLICT]: 409,
	[REFUSED]: 422,
	[DISABLED]: 403,
	[NOT_COMMITTED]: 422
};

/** Asks the named project's flai over the channel, whatever request (if any) is in progress. */
function viaChannel(key: string): Ask {
	return (method, params, opt) =>
		agentFor(key).ask(method, params, opt?.timeoutMs, opt?.onProgress);
}

/**
 * Repo is the dashboard's view of one system-flow repository. The project, its work items, and its
 * threads are asked of flai on the host over the channel and kept until flai says a file changed
 * (S-0073), and so are documents, narratives, the inbox, and search (S-0074): nothing here reads a
 * file of the project. A change is announced as 'change' with the repo-relative path.
 */
export class Repo extends EventEmitter {
	readonly root: string;
	private answers = new Map<string, Promise<unknown>>();
	private listening = false;
	private source: Ask;
	private returned: (ms: number) => Promise<boolean>;
	/** The project this Repo asks over the channel (S-0095); null for one given its own source. */
	private channelKey: string | null;

	constructor(
		root = projectDir(),
		source?: Ask,
		/** Whether flai is back within ms, for the one retry of a write; tests give their own. */
		returned?: (ms: number) => Promise<boolean>,
		key: string = currentProjectKey()
	) {
		super();
		this.root = resolve(root);
		this.channelKey = source ? null : key;
		this.source = source ?? viaChannel(key);
		this.returned =
			returned ?? (source ? async () => false : (ms) => connectedWithin(agentFor(key), ms));
	}

	/** Ask flai, with its refusals as the HTTP statuses the routes answer with. */
	async ask<T>(method: string, params: Record<string, unknown> = {}, opt?: AskOptions): Promise<T> {
		try {
			return await this.source<T>(method, params, opt);
		} catch (e) {
			if (e instanceof AgentError)
				throw new RepoError(
					(e.code !== undefined && STATUS[e.code]) || (e.code === undefined ? e.status : 500),
					e.message,
					e.data
				);
			throw e;
		}
	}

	/** A command of flai that changes nothing: a preview, a template, a diff, the statistics. */
	async run<T>(method: string, params: Record<string, unknown> = {}, opt?: AskOptions) {
		return this.ask<Written<T>>(method, params, opt);
	}

	/**
	 * A write (S-0075). It carries a request ID, so that when the connection is lost before the answer
	 * and flai returns within a few seconds, the same request is sent once more and flai answers from
	 * its journal if it had already done it: a retry can never move, accept, or save twice.
	 */
	async write<T>(method: string, params: Record<string, unknown> = {}, opt?: AskOptions) {
		const body = { ...params, request_id: randomUUID() };
		const timed = { timeoutMs: 60000, ...opt };
		try {
			return await this.source<Written<T>>(method, body, timed).catch(async (e) => {
				const lost = e instanceof AgentError && e.status === 502 && e.code === undefined;
				if (!lost || !(await this.returned(5000))) throw e;
				return this.source<Written<T>>(method, body, timed);
			});
		} catch (e) {
			if (e instanceof AgentError) {
				const status =
					(e.code !== undefined && STATUS[e.code]) || (e.code === undefined ? e.status : 500);
				throw new RepoError(status, e.message, e.data);
			}
			throw e;
		} finally {
			this.answers.clear();
		}
	}

	/** One answer per question until something changes; a failure is not kept. */
	remember<T>(key: string, method: string, params: Record<string, unknown> = {}): Promise<T> {
		let hit = this.answers.get(key) as Promise<T> | undefined;
		if (!hit) {
			hit = this.ask<T>(method, params);
			this.answers.set(key, hit);
			hit.catch(() => this.answers.delete(key));
		}
		return hit;
	}

	/** Forget every answer; the next question is asked of flai again. */
	forget(): void {
		this.answers.clear();
	}

	/** A file of the project changed, by flai's word: forget what was asked and tell the listeners. */
	changed(path: string): void {
		this.answers.clear();
		log().debug({ component: 'watcher', path }, 'file changed');
		this.emit('change', path);
	}

	async manifest(): Promise<Manifest> {
		const m = await this.remember<Manifest>('project', 'project.info');
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
		return this.remember<Thread[]>('threads', 'threads.list', { all: true });
	}

	/** Threads anchored to a repository path or an item ID. */
	async threadsFor(on: string): Promise<Thread[]> {
		const want = on.replace(/\/$/, '');
		return this.remember<Thread[]>(`threads:${want}`, 'threads.list', { on: want, all: true });
	}

	/** The board as flai lays it out, epics and tasks included (board.ts gives it its shape). */
	async boardView<T>(): Promise<T> {
		return this.remember<T>('board', 'board.get', { all: true });
	}

	/** All work items from kanban and archive, sorted by ID, as flai reads them. */
	async items(): Promise<Item[]> {
		const raw = await this.remember<FlaiItem[]>('items', 'items.list', {
			archived: true,
			bodies: true
		});
		return raw.map(fromFlai);
	}

	/** Find an item by ID in any padding (S-32, S-032, S-0032 name the same item). */
	async itemById(id: string): Promise<{ item: Item; children: Item[] }> {
		const got = await this.remember<{ item: FlaiItem; children: FlaiItem[] | null }>(
			`item:${id.trim()}`,
			'item.get',
			{ id: id.trim() }
		);
		return { item: fromFlai(got.item), children: (got.children ?? []).map(fromFlai) };
	}

	/** Documentation trees: design (all types), docs, and wip, each Markdown file with its front matter. */
	async docsTree(): Promise<DocNode[]> {
		return this.remember<DocNode[]>('docs', 'docs.tree');
	}

	/** One markdown file, raw, with parsed front matter. flai refuses anything that is not a document of the project. */
	async docFile(rel: string): Promise<{
		path: string;
		frontMatter: Record<string, unknown> | null;
		body: string;
		raw: string;
	}> {
		return this.remember(`doc:${rel}`, 'doc.get', { path: rel });
	}

	/**
	 * Listen for flai's word that a file changed, and for its return: a flai that comes back may have
	 * missed changes, so everything asked is forgotten and open pages are told to look again.
	 */
	async watch(): Promise<void> {
		const key = this.channelKey;
		if (this.listening || key === null) return;
		this.listening = true;
		// At the registry, not at a hub resolved now (S-0095): the project's flai may not have
		// connected yet, and the default key follows whichever project is the only one when an event
		// arrives, so a Repo made at startup, before any flai, still hears from the first to connect.
		const reg = registry();
		reg.on('change', (from: string, path: string) => {
			if (concerns(key, from)) this.changed(path);
		});
		reg.on('connected', (from: string) => {
			if (concerns(key, from)) this.changed('system-flow.yaml');
		});
		// What was asked of a flai that has gone is not shown as if it were current: without flai the
		// routes answer 503, and the pages say why.
		reg.on('gone', (from: string) => {
			if (concerns(key, from)) this.forget();
		});
		log().info({ component: 'watcher', project: key }, 'listening for changes from the host flai');
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

const repos = new Map<string, Repo>();

/**
 * The current project's Repo (S-0080): one per project key, so that one project's cached answers,
 * watcher, and cursor never leak into another's. Which project is current is set for the length of
 * a request by hooks.server.ts (agent.ts's projectContext); code with none of its own, a test that
 * never sets it included, gets the same one Repo it always has, keyed under defaultProjectKey.
 */
export function repo(): Repo {
	const key = currentProjectKey();
	let r = repos.get(key);
	if (!r) {
		r = new Repo();
		repos.set(key, r);
	}
	return r;
}

/**
 * Replace a project's Repo, or forget it with null. For tests, which ask a flai of their own. With
 * no key, replaces the default project's, as it always has for a single-project dashboard's tests.
 */
export function useRepo(r: Repo | null, key: string = defaultProjectKey): void {
	if (r) repos.set(key, r);
	else repos.delete(key);
}

/** Every Repo made so far, for tests that need to reach into more than the current project's. */
export function knownRepos(): Map<string, Repo> {
	return repos;
}
