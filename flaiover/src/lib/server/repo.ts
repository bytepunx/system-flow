// The project as the server sees it: nothing on disk. Every answer is asked of flai on the host
// over the channel (ADR-0029), and the container holds no file of the project (ADR-0031). The
// shapes mirror design/system/work-hierarchy.md and repository-layout.md; flai's Go
// implementation is the reference.
import type { Agent } from '$lib/agent';
import type { CostOfDelay, Forecast } from '$lib/planning';
import { planFrom, type FlaiPlan, type TaskPlan } from '$lib/taskplan';
import type { ExpectedCost } from '$lib/usage';
import { posix, resolve } from 'node:path';
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
import { kindOf, type Change, type ChangeKind } from '$lib/changes';

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
	/** The project's default agent, which a story created now gets (S-0103). */
	agent?: Agent;
};

export type Transition = { to: string; at: string; by: string };
export type Block = { from: string; until?: string; reason: string };

// operator is set by /api/threads: whether the operator wrote the entry rather than an agent (S-0127).
export type ThreadEntry = { at: string; author: string; text: string; operator?: boolean };
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

/** One dated contribution to a conversation; story is empty on the entry that closes it. */
export type ConversationEntry = { at: string; author: string; story: string; text: string };
/**
 * A conversation between two stories (S-0336), as flai's messages.View shapes it: closed and its
 * reason are read from the file and the two stories' states, and awaiting names the story whose
 * reply it waits for while open, empty once closed.
 */
export type Conversation = {
	id: string;
	title: string;
	from: string;
	to: string;
	about: string[];
	status: 'open' | 'closed';
	closed: boolean;
	closed_reason: string;
	awaiting: string;
	participants: string[];
	created: string;
	updated: string;
	path: string; // repo-relative
	entries: ConversationEntry[];
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
	topics?: string[]; // what a story or epic is about beyond its components (S-0135, ADR-0047)
	agent?: Agent; // who works a story (S-0103)
	draft?: boolean; // a story an agent wrote that the operator has not finalized (S-0199)
	cost_of_delay?: CostOfDelay; // epics and stories (S-0199)
	forecast?: Forecast; // stories (S-0199)
	path: string; // repo-relative
	archived: boolean;
	body: string;
};

/** Which items to ask flai for: by type and state, and the archive only where it is shown (S-0162). */
export type ItemQuery = { type?: string; status?: string; archive?: boolean };

/** How many items are active and how many archived (flai's items.count, S-0162). */
export type ItemCount = { active: number; archived: number };

/** A folder or Markdown file of the documentation trees: a file's title, not its front matter (S-0162). */
export type DocNode = {
	name: string;
	path: string; // repo-relative, forward slashes
	kind: 'dir' | 'file';
	title?: string;
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
/**
 * A write's options. retry: false is for a write whose success ends the connection it came on
 * (restarting flai serve, S-0107): losing the connection is then the expected outcome, and the
 * flai that answers the retry has no record of the request, so it would do it again.
 */
export type WriteOptions = AskOptions & { retry?: boolean };
export type Ask = <T>(
	method: string,
	params?: Record<string, unknown>,
	opt?: AskOptions
) => Promise<T>;

/** What a command of flai answered: its JSON, and the warnings it logged. */
export type Written<T> = { data: T; warnings: string[] };

/** One markdown lint finding of a document's content, as flai's doc.lint gives it (S-0378). */
export type LintFinding = { line: number; rule: string; alias: string; message: string };

/** What flai's doc.lint answers: the document's path, cleaned, and its findings, [] when clean. */
export type DocLint = { path: string; findings: LintFinding[] };

/** A CSS custom property a design token renders as, such as --color-primary (ADR-0146). */
export type DesignVariable = { name: string; value: string };

/** One design token in one theme, as flai's design.get gives it (S-0381). */
export type DesignToken = {
	/** Its dotted path, such as color.primary. */
	path: string;
	/** Its DTCG type: color, dimension, fontFamily, typography, shadow, and so on. */
	type: string;
	/** Its $value as JSON; an alias's is its target's. */
	value: unknown;
	/** Its first CSS property's value; '' when it renders nothing. */
	css: string;
	/** The CSS custom properties it renders as; [] when none. */
	variables: DesignVariable[];
	/** The token path an alias names; absent for a token that is not an alias. */
	ref?: string;
	/** Its line in the token file, from 1. */
	line: number;
};

/** One theme's tokens, in file order. */
export type DesignTheme = { name: 'light' | 'dark'; tokens: DesignToken[] };

/** A problem of the token file: path is the token's, absent for the whole file. */
export type DesignProblem = { path?: string; line: number; message: string };

/** One colour role pair's contrast in one theme: a report, never a problem. */
export type DesignContrastPair = {
	theme: 'light' | 'dark';
	/** Token paths, such as color.text and color.background. */
	foreground: string;
	background: string;
	/** The two colours in sRGB, #rrggbb. */
	foreground_hex: string;
	background_hex: string;
	/** Truncated to two decimals, so that it never reads as passing when it fails. */
	ratio: number;
	/** The least ratio that passes. */
	threshold: number;
	passes: boolean;
	/** Either colour is an oklch colour outside sRGB, clipped to it. */
	clipped: boolean;
};

/** One theme's role pairs. */
export type DesignContrast = { theme: 'light' | 'dark'; pairs: DesignContrastPair[] };

/**
 * The project's design system as flai's design.get answers it, what flai design show --json prints
 * (S-0381, Go's designsystem.Shown). With no token file, exists is false, no_token_file says so and
 * names flai design init, and every list is empty. Lists are never null.
 */
export type DesignSystem = {
	/** The token file's repository path, the manifest's design.tokens or its default. */
	path: string;
	exists: boolean;
	/** Set only when there is no token file: so, and the command that writes one. */
	no_token_file?: string;
	/** The guidance document's repository path; '' when there is no token file. */
	guidance: string;
	/** The generated stylesheet the manifest's design.stylesheet names; absent when it names none (S-0384). */
	stylesheet?: string;
	/** The guidance's headings; [] when it is missing. */
	headings: string[];
	/** Whether the token file is a DTCG resolver document. */
	resolver: boolean;
	/** Light, then dark; [] when the file does not parse. */
	themes: DesignTheme[];
	/** By line, then path. */
	problems: DesignProblem[];
	/** The tokens rendered as CSS custom properties; '' when the file does not parse. */
	css: string;
	/** The tokens rendered as a Tailwind v4 @theme block; '' when the file does not parse. */
	tailwind: string;
	contrast: DesignContrast[];
};

/**
 * An open story whose branch, or whose worktree uncommitted, changes the design system, as flai's
 * design.proposals gives it (S-0384, Go's designsystem.Proposal).
 */
export type DesignProposal = {
	id: string;
	title: string;
	/** Its state, never done or cancelled. */
	status: string;
	/** The token file, the guidance, or both, in that order. */
	paths: string[];
};

/** A text file that is not markdown as flai's file.get answers it (ADR-0147): raw, with its SHA-256. */
export type ProjectFile = { path: string; content: string; hash: string };

/** What flai's file.save is given (ADR-0147); with no story it saves and commits on the main checkout. */
export type FileSave = {
	path: string;
	content: string;
	/** The SHA-256 the file was loaded with: a stale one is a conflict. */
	hash: string;
	message?: string;
	story?: string;
};

/** The kind of value a component's prop takes, as flai components list reads it (S-0437). */
export type ComponentPropKind =
	| 'string'
	| 'number'
	| 'boolean'
	| 'union'
	| 'array'
	| 'object'
	| 'function'
	| 'snippet'
	| 'unknown';

/** One prop of a component, in the order its $props() destructures them (Go's components.Prop). */
export type ComponentProp = {
	name: string;
	kind: ComponentPropKind;
	/** Its type as the source writes it; absent when it has none. */
	type?: string;
	/** Its default as the source writes it; absent when it has none. */
	default?: string;
	/** The type does not mark it optional and it has no default. */
	required: boolean;
	/** A union's members, each as JSON: "a", 1, true. */
	literals?: string[];
};

/** One .svelte file of a sub-project (Go's components.Component). */
export type Component = {
	/** The sub-project's path relative to the project root, in slash form. */
	project: string;
	/** The file's path relative to the sub-project, in slash form. */
	path: string;
	/** The file's base name without .svelte. */
	name: string;
	/** [] when it has no $props() or they could not be read. */
	props: ComponentProp[];
	/** Why its props could not be read, such as a parse error; absent when they were. */
	error?: string;
};

/** One sub-project of the manifest in the catalogue; only a supported one has components. */
export type ComponentsProject = {
	name: string;
	path: string;
	kind: string;
	/** Whether its kind is sveltekit or svelte, so that flai lists its components. */
	supported: boolean;
	/** The module its harness imports first, the manifest's components.setup; absent when none. */
	setup?: string;
	/** Sorted by path; [] for a sub-project that is not supported. */
	components: Component[];
};

/** The component catalogue as flai's components.list answers it, what flai components list --json prints. */
export type ComponentsCatalogue = { projects: ComponentsProject[] };

/** How one sub-project's component dev server stands (Go's hostapi.ComponentServer, S-0438). */
export type ComponentServer = {
	/** The sub-project's name in the manifest. */
	project: string;
	/** Its dev server runs and answers. */
	running: boolean;
	/** Its dev server was started and does not answer yet. */
	starting?: boolean;
	/** Set by a start that answered the dev server that ran rather than start one. */
	already_running?: boolean;
	/** The harness page, http://127.0.0.1:<port>/__flai/components/. */
	url?: string;
	port?: number;
	/** The process ID of the sub-project's vite. */
	pid?: number;
	/** When vite was started, RFC 3339, UTC. */
	since?: string;
	/** Why its last dev server failed to start or exited on its own, until one starts again. */
	last_error?: string;
};

/** What flai's components.stop answers: stopped is false when no dev server ran. */
export type ComponentStopped = { project: string; stopped: boolean; pid?: number };

/** The component dev servers as flai's components.status answers them (Go's hostapi.ComponentsStatus). */
export type ComponentsStatus = {
	/** Whether the components host action is on for the project. */
	enabled: boolean;
	/** What the operator runs on the host to turn it on: flai serve enable components. */
	enable: string;
	/** How long a dev server may go unasked before flai stops it; absent when flai runs no dev servers. */
	idle_seconds?: number;
	/** Every supported sub-project's, in the manifest's order. */
	servers: ComponentServer[];
};

/** JSON-RPC's code for a method the host flai does not offer: it is older than this dashboard. */
const METHOD_NOT_FOUND = -32601;
/** flai's code for an item or thread that does not exist (internal/hostapi). */
const NOT_FOUND = -32004;
const INVALID_PARAMS = -32602;
const CONFLICT = -32009;
const REFUSED = -32010;
const RULE = -32011;
const DISABLED = -32012; // a host action the operator has not enabled (S-0078)
const NOT_COMMITTED = -32013; // an import applied but not committed: a test failed (S-0098)
const STATUS: Record<number, number> = {
	[METHOD_NOT_FOUND]: 501,
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
 * file of the project. A change forgets only the answers read from its kind of file (S-0161). It is
 * announced as 'change' with the repo-relative path and its kind, and a story's agent starting or
 * ending as 'agent' with the story's ID (S-0154).
 */
export class Repo extends EventEmitter {
	readonly root: string;
	private answers = new Map<string, Promise<unknown>>();
	/** The manifest's layout as last read, to tell what kind of file a change is (S-0161). */
	private seen: Layout | null = null;
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
	async write<T>(method: string, params: Record<string, unknown> = {}, opt?: WriteOptions) {
		const body = { ...params, request_id: randomUUID() };
		const { retry = true, ...rest } = opt ?? {};
		const timed = { timeoutMs: 60000, ...rest };
		try {
			return await this.source<Written<T>>(method, body, timed).catch(async (e) => {
				const lost = e instanceof AgentError && e.status === 502 && e.code === undefined;
				if (!lost || !retry || !(await this.returned(5000))) throw e;
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

	/**
	 * A file of the project changed, by flai's word: forget the answers read from its kind of file
	 * (S-0161), and tell the listeners the path and its kind. Until the layout is known every answer
	 * is forgotten, and the manifest is asked for so that the next change is told apart.
	 */
	changed(path: string): void {
		const change: Change = { path, kind: kindOf(path, this.seen) };
		if (change.kind === 'project') this.seen = null;
		for (const key of [...this.answers.keys()]) if (forgets(key, change)) this.answers.delete(key);
		if (!this.seen) void this.manifest().catch(() => {});
		log().debug({ component: 'watcher', path, kind: change.kind }, 'file changed');
		this.emit('change', path, change.kind);
	}

	async manifest(): Promise<Manifest> {
		const m = await this.remember<Manifest>('project', 'project.info');
		if (!m || !m.layout?.design || !m.layout?.docs || !m.layout?.wip) {
			throw new RepoError(
				500,
				'system-flow.yaml is missing layout.design, layout.docs, or layout.wip'
			);
		}
		this.seen = m.layout;
		return m;
	}

	async layout(): Promise<Layout> {
		return (await this.manifest()).layout;
	}

	/** Who the operator is: the manifest's owner, else "designer", as flai decides for its writes and inbox. */
	async operator(): Promise<string> {
		return (await this.manifest()).owner || 'designer';
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

	/** Every conversation between stories, closed ones included, sorted by ID, as flai reads them (S-0336). */
	async messages(): Promise<Conversation[]> {
		const raw = await this.remember<FlaiConversation[]>('messages', 'messages.list', { all: true });
		return raw.map(fromFlaiConversation);
	}

	/** The conversations a story is one of the two sides of, closed ones included, its ID in any padding. */
	async messagesFor(story: string): Promise<Conversation[]> {
		const want = story.trim();
		const raw = await this.remember<FlaiConversation[]>(`messages:${want}`, 'messages.list', {
			story: want,
			all: true
		});
		return raw.map(fromFlaiConversation);
	}

	/** The board as flai lays it out, epics and tasks included (board.ts gives it its shape). */
	async boardView<T>(): Promise<T> {
		return this.remember<T>('board', 'board.get', { all: true });
	}

	/**
	 * Work items from kanban, and the archive with them when asked, sorted by ID, as flai reads
	 * them, without their bodies: a body is asked for when its item is opened (S-0162). Each query
	 * is kept on its own.
	 */
	async items(q: ItemQuery = {}): Promise<Item[]> {
		const params: Record<string, unknown> = { archived: q.archive === true };
		if (q.type) params.type = q.type;
		if (q.status) params.status = q.status;
		const key = `items:${q.archive ? 'archive' : 'active'}:${q.type ?? ''}:${q.status ?? ''}`;
		const raw = await this.remember<FlaiItem[]>(key, 'items.list', params);
		return raw.map(fromFlai);
	}

	/** How many items are active and how many archived, without asking for the archive (S-0162). */
	async itemCount(): Promise<ItemCount> {
		return this.remember<ItemCount>('items:count', 'items.count');
	}

	/**
	 * Find an item by ID in any padding (S-32, S-032, S-0032 name the same item), with its children,
	 * for a story with tasks its task plan (S-0176), and for an open story the paths its branch
	 * undoes of main's change (S-0388).
	 */
	async itemById(id: string): Promise<{
		item: Item;
		children: Item[];
		plan?: TaskPlan;
		expected_cost?: ExpectedCost;
		undoes_main?: string[];
	}> {
		const got = await this.remember<{
			item: FlaiItem;
			children: FlaiItem[] | null;
			plan?: FlaiPlan;
			expected_cost?: ExpectedCost;
			undoes_main?: string[] | null;
		}>(`item:${id.trim()}`, 'item.get', { id: id.trim() });
		const plan = planFrom(got.plan);
		return {
			item: fromFlai(got.item),
			children: (got.children ?? []).map(fromFlai),
			...(plan ? { plan } : {}),
			// what the item is expected to cost, priced by flai (ADR-0083)
			...(got.expected_cost ? { expected_cost: got.expected_cost } : {}),
			// the paths the open story's branch undoes of main's change, from flai's record (ADR-0155)
			...(got.undoes_main?.length ? { undoes_main: got.undoes_main } : {})
		};
	}

	/** Documentation trees: design (all types), docs, and wip, each Markdown file with its title. */
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
	 * The markdown lint findings of a document's content as the editor holds it (S-0378), asked
	 * afresh each time: the content is not a file of the project, so there is nothing to remember.
	 */
	async docLint(rel: string, content: string): Promise<DocLint> {
		return this.ask<DocLint>('doc.lint', { path: rel, content });
	}

	/**
	 * The project's design system (S-0382): its tokens per theme, their problems, renderings, and
	 * contrast, or that there is no token file, as flai's design.get answers it. With story, from that
	 * open story's worktree (S-0384), asked afresh each time: a branch's files change without flai's
	 * word, which is of the main checkout only. flai refuses a story that is not open or has no worktree.
	 */
	async designSystem(story?: string): Promise<DesignSystem> {
		if (story) return this.ask<DesignSystem>('design.get', { story });
		return this.remember<DesignSystem>('design', 'design.get');
	}

	/**
	 * The open stories whose branches propose a change to the design system (S-0384), as flai's
	 * design.proposals answers them, [] when none. Asked afresh each time, as stream.diff is: a
	 * branch's commits are not announced as changes.
	 */
	async designProposals(): Promise<DesignProposal[]> {
		return this.ask<DesignProposal[]>('design.proposals');
	}

	/**
	 * A text file that is not markdown, such as the design token file, through flai's file.get
	 * (ADR-0147), from the story's worktree when story names one. Asked afresh each time, as doc.show
	 * is, so that the hash a save is given is the file's as it is now.
	 */
	async file(rel: string, story?: string): Promise<ProjectFile> {
		return (await this.run<ProjectFile>('file.get', { path: rel, story })).data;
	}

	/**
	 * Save a text file that is not markdown through flai's file.save (ADR-0147), by doc.save's rules:
	 * 409 with current, hash, and diff when the hash is stale; 422 with findings when flai refuses the
	 * content. As every write, it forgets every answer, design.get's included.
	 */
	async saveFile(save: FileSave): Promise<Written<Record<string, unknown>>> {
		return this.write<Record<string, unknown>>('file.save', { ...save });
	}

	/**
	 * The component catalogue, flai's components.list (S-0437): every sub-project, a Svelte one with
	 * its components and their props. Asked afresh each time: a component file may be added at any time.
	 */
	async components(): Promise<ComponentsCatalogue> {
		return upgraded('list the components', () => this.ask<ComponentsCatalogue>('components.list'));
	}

	/**
	 * How every supported sub-project's component dev server stands, flai's components.status (S-0438),
	 * and whether the components host action is on. Asked afresh each time; asking keeps a running
	 * dev server from being stopped as idle.
	 */
	async componentsStatus(): Promise<ComponentsStatus> {
		return upgraded('show the component dev servers', () =>
			this.ask<ComponentsStatus>('components.status')
		);
	}

	/**
	 * Start a sub-project's component dev server, flai's components.start under the components host
	 * action (S-0438), and answer it once its harness answers; one that runs is answered as it is.
	 * project may be left out when one sub-project is supported. flai waits a minute for Vite, so the
	 * ask waits longer, to answer flai's own message, which carries Vite's last lines of output.
	 */
	async componentsStart(project?: string): Promise<ComponentServer> {
		const { data } = await upgraded('start a component dev server', () =>
			this.write<ComponentServer>('components.start', project ? { project } : {}, {
				timeoutMs: 90000
			})
		);
		return data;
	}

	/** Stop a sub-project's component dev server when one runs, flai's components.stop (S-0438). */
	async componentsStop(project?: string): Promise<ComponentStopped> {
		const { data } = await upgraded('stop a component dev server', () =>
			this.write<ComponentStopped>('components.stop', project ? { project } : {})
		);
		return data;
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
		// an agent's run is asked of flai each time, not remembered: nothing to forget (S-0154)
		reg.on('agent', (from: string, story: string) => {
			if (concerns(key, from)) this.emit('agent', story);
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

/**
 * What each remembered answer is read from in flai (S-0161), by its key up to the first colon: a
 * change of another kind leaves it kept. project.info is read from the manifest alone, which
 * forgets everything; a thread's view names its story, read from the items; docs.tree carries
 * every file's front matter under the three folders. doc:<path> is read from its own file only.
 * A conversation's file under the wip folder's messages is a document, and whether it reads as
 * closed follows its two stories' states, read from the items. design.get reads the token file and
 * the guidance, both documents under the design folder by default; flai serve watches only the
 * manifest and the three folders, so a token file design.tokens places outside them is not heard of.
 */
const READS: Record<string, ChangeKind[]> = {
	project: [],
	board: ['item'],
	items: ['item'],
	item: ['item'],
	stats: ['item'],
	threads: ['thread', 'item'],
	messages: ['document', 'item'],
	inbox: ['item', 'thread', 'narrative'],
	activity: ['item', 'narrative'],
	adrs: ['adr'],
	design: ['document'],
	docs: ['item', 'narrative', 'thread', 'adr', 'document']
};

/** Whether change makes the answer kept under key stale; an answer the table does not name always is. */
function forgets(key: string, change: Change): boolean {
	if (change.kind === 'project' || change.kind === 'other') return true;
	const at = key.indexOf(':');
	const name = at < 0 ? key : key.slice(0, at);
	if (name === 'doc') return clean(key.slice(at + 1)) === clean(change.path);
	const reads = READS[name];
	return !reads || reads.includes(change.kind);
}

function clean(rel: string): string {
	return posix.normalize(rel.replaceAll('\\', '/')).replace(/^(\.\/|\/)+/, '');
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

/**
 * The story a document endpoint reads or saves on (S-0378): undefined when none is named, the
 * main checkout's then, else an ID of S- and digits; anything else is refused with 400.
 */
export function storyParam(value: unknown): string | undefined {
	if (value === undefined || value === null || value === '') return undefined;
	if (typeof value !== 'string' || !/^S-\d+$/.test(value))
		throw new RepoError(
			400,
			`story must be a story's ID such as S-0042, not ${JSON.stringify(value)}; leave it out to edit the main checkout`
		);
	return value;
}

/**
 * The sub-project a component dev server endpoint starts or stops (S-0439): undefined when none is
 * named, flai's one supported sub-project then, else its name in the manifest; anything but a string
 * is refused with 400.
 */
export function subProjectParam(value: unknown): string | undefined {
	if (value === undefined || value === null) return undefined;
	if (typeof value !== 'string')
		throw new RepoError(
			400,
			`project must be a sub-project's name in the manifest, not ${JSON.stringify(value)}; leave it out when one sub-project is Svelte`
		);
	return value.trim() || undefined;
}

/** Run fn, a host flai without the method it asks for told to upgrade (501), as /api/design/tokens does. */
async function upgraded<T>(what: string, fn: () => Promise<T>): Promise<T> {
	try {
		return await fn();
	} catch (e) {
		if (e instanceof RepoError && e.status === 501)
			throw new RepoError(
				501,
				`the host flai is too old to ${what} (${e.message}): upgrade flai on the host with flai self-upgrade`
			);
		throw e;
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
		body: it.body ?? '',
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

/** A conversation as flai marshals it: participants is null when its front matter names none. */
type FlaiConversation = Omit<Conversation, 'participants'> & { participants: string[] | null };

function fromFlaiConversation(c: FlaiConversation): Conversation {
	return { ...c, participants: c.participants ?? [] };
}

const repos = new Map<string, Repo>();

/**
 * The current project's Repo (S-0080): one per project key, so that one project's cached answers,
 * watcher, and cursor never leak into another's. Each listens for its own flai's changes from the
 * moment it is made (S-0117). Which project is current is set for the length of
 * a request by hooks.server.ts (agent.ts's projectContext); code with none of its own, a test that
 * never sets it included, gets the same one Repo it always has, keyed under defaultProjectKey.
 */
export function repo(): Repo {
	const key = currentProjectKey();
	let r = repos.get(key);
	if (!r) {
		r = new Repo();
		repos.set(key, r);
		// Every project's answers are forgotten when its flai says a file changed, not only the
		// default project's, which hooks.server.ts watches at startup (S-0117).
		void r.watch();
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
