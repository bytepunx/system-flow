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
	/** The empty value of its type, which flai gives an array, object, or unknown prop: such as
	 * [{"id": ""}], or [], {}, or null where flai could not resolve the type (S-0449); absent from
	 * a host flai older than S-0449. */
	skeleton?: unknown;
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
	/** Its examples, in its examples file's order, [] when it has none (S-0443, ADR-0164); absent
	 * from a host flai older than S-0443. */
	examples?: ComponentExample[];
	/** Its examples file's repository path, given whether or not the file exists; absent from a
	 * host flai older than S-0443. */
	examples_path?: string;
};

/** One named set of a component's props, as its examples file keeps it (ADR-0164). */
export type ComponentExample = {
	name: string;
	/** Prop name to value, as the props panel sends them; a function prop is left out. */
	props: Record<string, unknown>;
	/** The one theme the page shows; absent, the page keeps the theme or themes it shows. */
	theme?: 'light' | 'dark';
	/** The frame's width in CSS pixels; absent, the page keeps the width it shows. */
	width?: number;
	/** What the harness answers the component's requests to the app's own server with, in the order
	 * it tries them (S-0476, ADR-0187); absent, every such request gets the harness's 501. */
	answers?: ComponentAnswer[];
};

/** One answer an example gives a framed component's fetch, or with events its EventSource, to the
 * app's own server (ADR-0187). flai check's components.examples holds it to the format. */
export type ComponentAnswer = {
	/** The request's path, starting with /, with its query when it answers that query alone. */
	path: string;
	/** An HTTP method, compared without case; absent, it matches every method. */
	method?: string;
	/** The fetch's status, 200 to 599; 200 when absent. */
	status?: number;
	/** The fetch's response headers, name to value. */
	headers?: Record<string, string>;
	/** The fetch's body: a string sent as it is, anything else as JSON; absent, empty. */
	body?: unknown;
	/** Makes it an EventSource's answer: the messages it sends, in order. */
	events?: ComponentAnswerEvent[];
};

/** One message an answered EventSource sends: data as it is when a string, as JSON otherwise. */
export type ComponentAnswerEvent = {
	/** The event's name; absent, message. */
	event?: string;
	data: unknown;
	id?: string;
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
	/** The story whose worktree the dev server runs in (S-0442); absent for the main checkout's. */
	story?: string;
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

/** The project's cluster as flai kluster status reports it (Go's kluster.Status, S-0577). */
export type KlusterCluster = {
	name: string;
	profile?: string;
	provider: string;
	/** The kube context kluster names for it. */
	context: string;
	/** Whether kluster status lists the cluster; false when it was never brought up or is down. */
	exists: boolean;
	/** Whether kluster says it runs. */
	running: boolean;
	/** Its age as kluster prints it, such as 2h50m39s. */
	age?: string;
	age_seconds?: number;
};

/**
 * One job flai serve ran, or runs, for the project's kluster (Go's hostapi.KlusterJob): an up or a
 * down of the cluster (S-0580), a story's deploy or teardown or main's promotion (S-0581), or the
 * checks run against a story's namespace the operator started (S-0663).
 */
export type KlusterJob = {
	action: 'up' | 'down' | 'deploy' | 'teardown' | 'promote' | 'checks';
	cluster: string;
	/** The story a deploy, a teardown, or a checks run is of; absent otherwise. */
	story?: string;
	/** The namespace a deploy, a teardown, or a promotion acts on; absent for an up or a down. */
	namespace?: string;
	/** The folder each step's own log goes to, for a deploy, a teardown, or a promotion. */
	logs?: string;
	/** RFC 3339, UTC. */
	started_at: string;
	/** Absent while it runs. */
	ended_at?: string;
	/** done or failed, once it ended. */
	outcome?: 'done' | 'failed';
	/** Why it failed, kluster's or flai's own sentence. */
	error?: string;
	/** The file kluster's output goes to, from the project's root. */
	log: string;
	/** The end of that file, at most 4 KB; only kluster.status sends it. */
	log_tail?: string;
};

/** What flai's kluster.status answers (Go's hostapi.KlusterStatus, S-0580). */
export type KlusterStatus = {
	/** Whether the kluster host action is on for the project. */
	enabled: boolean;
	/** What the operator runs on the host to turn it on: flai serve enable kluster. */
	enable: string;
	/** Whether the manifest has a kluster: block. */
	block: boolean;
	/** flai's sentence naming the block to add, when there is none. */
	no_block?: string;
	/** Absent when there is no block or its config cannot be read. */
	cluster?: KlusterCluster;
	/** Why the cluster could not be read, such as kluster missing from the host's PATH. */
	status_error?: string;
	/** The up or down that runs now. */
	running?: KlusterJob;
	/** The last up or down to end since flai serve started. */
	last?: KlusterJob;
	// What setup still lacks (S-0586), so the Klusters page can lead the operator through it; a
	// host flai older than S-0586 sends none of these, so each is optional.
	/** The kluster: block's keys flai manifest set writes, listed with no block too. */
	settings?: KlusterSetting[];
	/** Whether the settings host action is on: the block's and the config file's writes need it. */
	settings_on?: boolean;
	/** What the operator runs on the host to turn settings on: flai serve enable settings. */
	settings_enable?: string;
	/** The config file the block names, from the repository root; absent with no block. */
	config?: string;
	/** What that file names; absent with no block or no config named. */
	config_file?: KlusterConfigFile;
	/** What flai check finds wrong with the block, one sentence each beginning with its key. */
	findings?: string[];
	/** Whether the block names a deploy command; without one flai deploys nothing. */
	deployable?: boolean;
	/** kluster, docker, and kubectl, each found on the host's PATH or not, with how to install it. */
	tools?: KlusterTool[];
};

/** One key of the kluster: block (Go's manifest.SettingValue, S-0586), with its value. */
export type KlusterSetting = {
	key: string;
	/**
	 * text for config and namespace, command for build, deploy, and teardown, duration for timeout,
	 * components for components, a list of {name, build} (ADR-0263).
	 */
	kind: 'text' | 'command' | 'duration' | 'components';
	/** The value in effect when the key is unset, as it would be written. */
	default?: string;
	/** A sentence on what the setting does. */
	meaning: string;
	/**
	 * The value as the manifest has it, a command as its argument list, components as their list;
	 * absent when unset.
	 */
	value?: string | string[] | KlusterComponent[];
	set: boolean;
};

/** One component of the kluster: block (Go's manifest.KlusterComponent, ADR-0263): its name and build. */
export type KlusterComponent = { name: string; build: string[] };

/** What the kluster config file names (Go's hostapi.KlusterConfigFile, S-0586). */
export type KlusterConfigFile = {
	/** Whether the file is there. */
	exists: boolean;
	name?: string;
	provider?: string;
	/** The profile it names; empty when none, and kluster then installs default_profile. */
	profile?: string;
	default_profile: string;
	/** The profiles kluster offers: spire, signet, authstar. */
	profiles: string[];
	/** The opt-in addons it names. */
	addons: string[];
	/** Why it could not be read. */
	error?: string;
};

/** One program the project's kluster needs on the host (Go's hostapi.KlusterTool, S-0586). */
export type KlusterTool = {
	name: 'kluster' | 'docker' | 'kubectl';
	found: boolean;
	path?: string;
	/** What flai needs it for. */
	need: string;
	/** The command the operator runs on the host to install it; absent for docker. */
	install?: string;
	/** Where its installation is described. */
	docs: string;
};

/** One node of the cluster (Go's hostapi.KlusterNode, S-0586). */
export type KlusterNode = { name: string; ready: boolean; roles: string[]; version?: string };

/** One pod of a namespace (Go's hostapi.KlusterPod, S-0586). */
export type KlusterPod = {
	name: string;
	/** Pending, Running, Succeeded, Failed, or Unknown. */
	phase: string;
	/** How many of its containers are ready, of how many. */
	ready: number;
	containers: number;
	/** Its containers' restarts, summed. */
	restarts: number;
};

/** What flai's kluster.workloads answers (Go's hostapi.KlusterWorkloads, S-0586). */
export type KlusterWorkloads = {
	block: boolean;
	no_block?: string;
	cluster?: string;
	context?: string;
	nodes: KlusterNode[];
	/** Each namespace's pods, namespaces and pods by name. */
	namespaces: { namespace: string; pods: KlusterPod[] }[];
	/** Why they could not be read: no kubeconfig yet, kubectl missing, or kubectl's failure. */
	workloads_error?: string;
};

/**
 * One step of a deployment (Go's hostapi.KlusterStep): build, build-<name> per component, load,
 * deploy, or teardown. A deploy records every step it will run as pending before the first starts
 * (ADR-0263), so a pending step has no started_at yet.
 */
export type KlusterStep = {
	name: string;
	started_at?: string;
	ended_at?: string;
	/** pending, running, passed, failed, or skipped. */
	outcome: string;
	exit_code?: number;
	/** The step's log, from the project's root. */
	log?: string;
	/** The end of that log, at most 4 KB, cut to whole lines. */
	log_tail?: string;
};

/**
 * A read of one step's log from an offset (Go's hostapi.KlusterLog, ADR-0263): at most 64 KB of
 * text from offset, and next, the offset to ask from next, so a page follows a running step live.
 */
export type KlusterLog = {
	namespace: string;
	step: string;
	/** The log's path, from the project's root. */
	log: string;
	/** False when the step has no log yet. */
	exists: boolean;
	/** Where text starts in the log. */
	offset: number;
	/** The offset to ask from next: offset and the bytes of text. */
	next: number;
	/** The log's size when it was read. */
	size: number;
	text: string;
	/** The log was shorter than the offset asked, as when a new deploy cleared it; text starts at 0. */
	reset?: boolean;
	/** The step's outcome in the namespace's record; absent when the record has none. */
	outcome?: string;
};

/** Which step's log to read, and from where: the query of GET /api/kluster/log. */
export type KlusterLogTarget = { namespace: string; step: string; offset: number };

/** One port of a service, as kubectl reports it. */
export type KlusterPort = { name?: string; port: number; protocol?: string };

/** One service of a namespace, with the host it answers on and the command that forwards it. */
export type KlusterService = {
	name: string;
	/** The name it answers on inside the cluster, <name>.<namespace>.svc.cluster.local. */
	host?: string;
	ports: KlusterPort[];
	/** The kubectl port-forward command, run in the project's main checkout. */
	port_forward?: string;
};

/** How to reach a namespace for a demo: through a port-forward, as ADR-0254 decides. */
export type KlusterAccess = {
	/** flai's kubeconfig of the cluster, from the project's root. */
	kubeconfig: string;
	context: string;
	namespace: string;
	/** The command that lists the namespace's services, run in the project's main checkout. */
	list: string;
	services: KlusterService[];
	/** Why the services could not be read. */
	services_error?: string;
	/** The port-forwards flai serve runs, or ran, into the namespace (S-0586); absent from an older flai. */
	forwards?: KlusterForward[];
};

/**
 * One port-forward flai serve runs, or ran, to a service port of a namespace (Go's
 * hostapi.KlusterForward, S-0586): kubectl port-forward on local_port of the host, bound to address,
 * 0.0.0.0, so anything that reaches the host's network reaches it, with no sign-in of its own.
 */
export type KlusterForward = {
	namespace: string;
	service: string;
	/** The service's port in the cluster. */
	port: number;
	/** The host's port it listens on. */
	local_port: number;
	address: string;
	started_at: string;
	/** Whether kubectl still runs; an ended forward is kept with the end of its output. */
	running: boolean;
	ended_at?: string;
	output?: string;
};

/** The job that holds the project's kluster lock (Go's kluster.Job). */
export type KlusterLockJob = {
	kind: 'deploy' | 'teardown' | 'promote';
	story: string;
	namespace: string;
	started_at: string;
	pid: number;
	/** Who asked for it, as the record's asked_by; absent for a job run in a shell. */
	by?: string;
};

/** One namespace's deployment as kluster.deployments answers it (Go's hostapi.KlusterDeployment). */
export type KlusterDeployment = {
	/** The story deployed; empty for the default namespace. */
	story: string;
	namespace: string;
	/** Whether it is the default namespace, main's. */
	default: boolean;
	/** The commit deployed, the whole hash. */
	commit: string;
	tag: string;
	/** The checkout the commands ran in, from the project's root. */
	checkout: string;
	images: string[] | null;
	started_at: string;
	ended_at?: string;
	/** running, passed, failed, or interrupted. */
	outcome: string;
	failed_step?: string;
	exit_code?: number;
	error?: string;
	steps: KlusterStep[];
	/** The job running against the namespace now, if any. */
	job?: KlusterLockJob;
	/** The story's title and status, for a story's namespace. */
	title?: string;
	status?: string;
	/** The story branch's head now, or main's for the default namespace. */
	head?: string;
	head_error?: string;
	// What the page shows of each story namespace (S-0663); a host flai older than S-0663 sends
	// none of these, so each is optional.
	/** The story's column on the board, in-progress or review; absent for the default namespace. */
	column?: string;
	/** Whether commit is the story branch's head; false or absent for the default namespace. */
	at_head?: boolean;
	/** Where the last job's output is, from the project's root: its failed step's log, else its logs' folder. */
	log?: string;
	/** Who asked for the last job: the story's agent, or flai serve for the review pass; absent from a shell. */
	asked_by?: string;
	/** The commit the last checks run since the deploy ran against; absent until one ends. */
	checks_commit?: string;
	/** passed, failed, timed-out, or cancelled. */
	checks_outcome?: string;
	/** The last checks run's log, from the project's root. */
	checks_log?: string;
	/** Who asked for that checks run; absent for a run from a shell or the dashboard. */
	checks_asked_by?: string;
	/** The end of that log, at most 4 KB, cut to whole lines. */
	checks_log_tail?: string;
	access: KlusterAccess;
};

/** A story in review with no deployment recorded. */
export type KlusterReview = { story: string; title: string; head?: string; head_error?: string };

/** A story in progress or in review with no deployment recorded (Go's hostapi.KlusterUndeployed, S-0663). */
export type KlusterUndeployed = KlusterReview & {
	/** in-progress or review. */
	column: string;
};

/** A deploy flai refused for kluster.agent_namespaces, the limit (Go's kluster.Refusal, S-0663). */
export type KlusterRefusal = {
	story: string;
	/** The agent that asked. */
	by: string;
	/** RFC 3339. */
	at: string;
	/** kluster.agent_namespaces when it was refused. */
	limit: number;
	/** flai's words, naming the stories that held a namespace. */
	reason: string;
};

/** What flai's kluster.deployments answers (Go's hostapi.KlusterDeployments, S-0581). */
export type KlusterDeployments = {
	enabled: boolean;
	enable: string;
	block: boolean;
	no_block?: string;
	/** Why the block's config could not be read. */
	status_error?: string;
	cluster?: string;
	/** The namespace main is promoted into. */
	default_namespace?: string;
	/** The default namespace's first, then the stories' by number. */
	deployments: KlusterDeployment[];
	/** The stories in review with no deployment recorded. */
	review: KlusterReview[];
	// The stories in progress beside those in review, and the count against the limit (S-0663); a
	// host flai older than S-0663 sends none of these, so each is optional.
	/** The stories in progress or in review with no deployment recorded, by ID. */
	undeployed?: KlusterUndeployed[];
	/** How many stories in progress hold a namespace. */
	live?: number;
	/** Which stories in progress hold one. */
	holding?: string[];
	/** kluster.agent_namespaces, the most stories in progress that may hold one; 0 with no block. */
	limit?: number;
	/** The deploys refused for the limit, newest first. */
	refusals?: KlusterRefusal[];
	/** The deploy, teardown, or promotion that runs now, whoever started it. */
	running?: KlusterJob;
	/** The last deploy, teardown, or promotion this flai serve ran. */
	last?: KlusterJob;
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

	/**
	 * The unresolved threads, sorted by ID (ADR-0020), as flai reads them; all adds the resolved ones,
	 * which only a page that shows them asks for, since they are most of the list (S-0611). The two
	 * lists are kept apart, and a thread or item change forgets both.
	 */
	async threads(all = false): Promise<Thread[]> {
		return all
			? this.remember<Thread[]>('threads:all=1', 'threads.list', { all: true })
			: this.remember<Thread[]>('threads', 'threads.list', {});
	}

	/** The unresolved threads anchored to a repository path or an item ID; all adds the resolved ones. */
	async threadsFor(on: string, all = false): Promise<Thread[]> {
		const want = on.replace(/\/$/, '');
		return all
			? this.remember<Thread[]>(`threads:all=1:${want}`, 'threads.list', { on: want, all: true })
			: this.remember<Thread[]>(`threads:${want}`, 'threads.list', { on: want });
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
	 * With story, the catalogue in that open story's worktree (S-0442); flai refuses a story that is
	 * not open or has no worktree.
	 */
	async components(story?: string): Promise<ComponentsCatalogue> {
		return upgraded('list the components', () =>
			this.ask<ComponentsCatalogue>('components.list', story ? { story } : {})
		);
	}

	/**
	 * How every supported sub-project's component dev server stands, flai's components.status (S-0438),
	 * and whether the components host action is on. Asked afresh each time; asking keeps a running
	 * dev server from being stopped as idle. With story, that story's dev servers (S-0442).
	 */
	async componentsStatus(story?: string): Promise<ComponentsStatus> {
		return upgraded('show the component dev servers', () =>
			this.ask<ComponentsStatus>('components.status', story ? { story } : {})
		);
	}

	/**
	 * Start a sub-project's component dev server, flai's components.start under the components host
	 * action (S-0438), and answer it once its harness answers; one that runs is answered as it is.
	 * project may be left out when one sub-project is supported; with story, the dev server runs in
	 * that story's worktree (S-0442). flai waits a minute for Vite, so the ask waits longer, to answer
	 * flai's own message, which carries Vite's last lines of output. The sub-project reaches flai as
	 * sub_project, since project is the channel's key (S-0448).
	 */
	async componentsStart(project?: string, story?: string): Promise<ComponentServer> {
		const { data } = await upgraded('start a component dev server', () =>
			this.write<ComponentServer>(
				'components.start',
				{ ...(project ? { sub_project: project } : {}), ...(story ? { story } : {}) },
				{ timeoutMs: 90000 }
			)
		);
		return data;
	}

	/**
	 * Stop a sub-project's component dev server when one runs, flai's components.stop (S-0438); with
	 * story, that story's (S-0442). The sub-project reaches flai as sub_project, since project is the
	 * channel's key (S-0448).
	 */
	async componentsStop(project?: string, story?: string): Promise<ComponentStopped> {
		const { data } = await upgraded('stop a component dev server', () =>
			this.write<ComponentStopped>('components.stop', {
				...(project ? { sub_project: project } : {}),
				...(story ? { story } : {})
			})
		);
		return data;
	}

	/**
	 * The project's kluster, flai's kluster.status (S-0580): the cluster its manifest's kluster: block
	 * names, whether it runs and its age, whether the kluster host action is on, and the up or down
	 * flai serve runs or last ran, with the end of its log. Asked afresh each time: kluster answers
	 * from the cluster, not from a file flai watches.
	 */
	async klusterStatus(): Promise<KlusterStatus> {
		return upgraded('show the kluster', () => this.ask<KlusterStatus>('kluster.status'));
	}

	/**
	 * Start kluster up of the project's cluster, flai's kluster.up under the kluster host action
	 * (S-0580). flai answers the job at once and runs it in the background; kluster.status says how
	 * it ends.
	 */
	async klusterUp(): Promise<KlusterJob> {
		const { data } = await upgraded('bring the kluster up', () =>
			this.write<KlusterJob>('kluster.up')
		);
		return data;
	}

	/** Start kluster down of the project's cluster, flai's kluster.down, as klusterUp starts up. */
	async klusterDown(): Promise<KlusterJob> {
		const { data } = await upgraded('take the kluster down', () =>
			this.write<KlusterJob>('kluster.down')
		);
		return data;
	}

	/**
	 * The project's cluster's nodes and each namespace's pods, flai's kluster.workloads (S-0586),
	 * read with kubectl once flai has the cluster's kubeconfig. Asked afresh each time.
	 */
	async klusterWorkloads(): Promise<KlusterWorkloads> {
		return upgraded('show the kluster workloads', () =>
			this.ask<KlusterWorkloads>('kluster.workloads')
		);
	}

	/**
	 * The namespaces deployed to the project's kluster, flai's kluster.deployments (S-0581): main's
	 * first, then each story's, with its commit, outcome, each step's log, and how to reach it for a
	 * demo; the stories in review with no deployment; and the deploy that runs. Asked afresh each
	 * time: a deploy's record changes under .flai-cache, which flai does not announce.
	 */
	async klusterDeployments(): Promise<KlusterDeployments> {
		return upgraded('show the deployments', () =>
			this.ask<KlusterDeployments>('kluster.deployments')
		);
	}

	/**
	 * Read one step's log of a namespace from an offset, flai's kluster.log (ADR-0263): a read,
	 * whatever the kluster host action says, asked afresh each time so a page can follow it live.
	 */
	async klusterLog(t: KlusterLogTarget): Promise<KlusterLog> {
		return upgraded("read a deployment step's log", () => this.ask<KlusterLog>('kluster.log', t));
	}

	/**
	 * Start a story's deploy into its namespace, flai's kluster.deploy under the kluster host action
	 * (S-0581). flai answers the job at once and runs it in the background.
	 */
	async klusterDeploy(story: string): Promise<KlusterJob> {
		const { data } = await upgraded(`deploy ${story}`, () =>
			this.write<KlusterJob>('kluster.deploy', { story })
		);
		return data;
	}

	/**
	 * Start the project's checks run against a story's namespace, flai's kluster.checks under the
	 * kluster host action (S-0663), as klusterDeploy starts a deploy: flai answers the run at once.
	 */
	async klusterChecks(story: string): Promise<KlusterJob> {
		const { data } = await upgraded(`run the checks of ${story}`, () =>
			this.write<KlusterJob>('kluster.checks', { story })
		);
		return data;
	}

	/** Start the teardown of a story's namespace, flai's kluster.teardown, as klusterDeploy starts. */
	async klusterTeardown(story: string): Promise<KlusterJob> {
		const { data } = await upgraded(`tear ${story} down`, () =>
			this.write<KlusterJob>('kluster.teardown', { story })
		);
		return data;
	}

	/**
	 * Start a port-forward of a service port of a namespace, flai's kluster.forward.start under the
	 * kluster host action (S-0586): flai serve runs kubectl port-forward bound to 0.0.0.0 on a free
	 * port of the host, and answers the forward once it forwards, or the one that runs.
	 */
	async klusterForwardStart(f: KlusterForwardTarget): Promise<KlusterForward> {
		const { data } = await upgraded(`forward ${f.service} port ${f.port}`, () =>
			this.write<KlusterForward>('kluster.forward.start', f)
		);
		return data;
	}

	/** Stop a port-forward, flai's kluster.forward.stop, as klusterForwardStart starts one. */
	async klusterForwardStop(f: KlusterForwardTarget): Promise<KlusterForward> {
		const { data } = await upgraded(`stop forwarding ${f.service} port ${f.port}`, () =>
			this.write<KlusterForward>('kluster.forward.stop', f)
		);
		return data;
	}

	/** Start main's promotion into the default namespace, flai's kluster.promote. */
	async klusterPromote(): Promise<KlusterJob> {
		const { data } = await upgraded('promote main', () =>
			this.write<KlusterJob>('kluster.promote')
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
 * The story a kluster deploy or teardown is of (S-0581): a story's ID, which it must name; anything
 * else is refused with 400 before flai is asked.
 */
/** The service port of a namespace a port-forward is of (S-0586). */
export type KlusterForwardTarget = { namespace: string; service: string; port: number };

/**
 * A port-forward's target from a request's body (S-0586): a namespace and a service that are DNS
 * labels and a port from 1 to 65535, else 400, before flai is asked.
 */
export function klusterForwardParam(body: unknown): KlusterForwardTarget {
	const b = (body ?? {}) as Record<string, unknown>;
	const label = /^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$/;
	if (typeof b.namespace !== 'string' || !label.test(b.namespace))
		throw new RepoError(
			400,
			`namespace must be a namespace kluster.deployments lists, not ${JSON.stringify(b.namespace ?? null)}`
		);
	if (typeof b.service !== 'string' || !label.test(b.service))
		throw new RepoError(
			400,
			`service must be a service kluster.deployments lists, not ${JSON.stringify(b.service ?? null)}`
		);
	if (typeof b.port !== 'number' || !Number.isInteger(b.port) || b.port < 1 || b.port > 65535)
		throw new RepoError(
			400,
			`port must be one of the service's ports, not ${JSON.stringify(b.port ?? null)}`
		);
	return { namespace: b.namespace, service: b.service, port: b.port };
}

/**
 * The step log GET /api/kluster/log reads (ADR-0263), from its query: namespace a DNS label, step
 * build, load, deploy, teardown, or build-<name> with name a DNS label, and offset a whole number of
 * zero or more, 0 when left out; anything else is refused with 400 before flai is asked.
 */
export function klusterLogParam(query: URLSearchParams): KlusterLogTarget {
	const label = /^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$/;
	const namespace = query.get('namespace');
	if (namespace === null || !label.test(namespace))
		throw new RepoError(
			400,
			`namespace must be a namespace kluster.deployments lists, not ${JSON.stringify(namespace)}`
		);
	const step = query.get('step');
	const component = step?.startsWith('build-') ? step.slice('build-'.length) : null;
	if (
		step === null ||
		!(
			['build', 'load', 'deploy', 'teardown'].includes(step) ||
			(component !== null && label.test(component))
		)
	)
		throw new RepoError(
			400,
			`step must be build, build-<component>, load, deploy, or teardown, not ${JSON.stringify(step)}`
		);
	const raw = query.get('offset') ?? '0';
	if (!/^\d+$/.test(raw) || !Number.isSafeInteger(Number(raw)))
		throw new RepoError(
			400,
			`offset must be a whole number of bytes, 0 or more, not ${JSON.stringify(raw)}`
		);
	return { namespace, step, offset: Number(raw) };
}

export function klusterStoryParam(value: unknown): string {
	if (typeof value !== 'string' || !/^S-\d+$/.test(value))
		throw new RepoError(
			400,
			`story must be the ID of the story to deploy or tear down, such as S-0042, not ${JSON.stringify(value ?? null)}`
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
