// The host's settings as the dashboard's settings page sees them (S-0105), flai's settings.get.
import { agentFrom, configText, parseConfig, type Agent } from '$lib/agent';

export const SETTINGS_KINDS = [
	'action',
	'default_agent',
	'agent',
	'harness',
	'check',
	'checks_timeout',
	'import',
	'serve',
	'unserve',
	'mcp_token',
	'dashboard_token',
	'shared',
	'manifest'
] as const;
export type SettingsKind = (typeof SETTINGS_KINDS)[number];

/** Settings questions that change nothing: asked of flai with no host action and no request ID. */
export const SETTINGS_READS = ['shared_check'] as const;
export type SettingsRead = (typeof SETTINGS_READS)[number];

/** Settings kept in the host's configuration for every project: they need settings on everywhere. */
export const HOSTWIDE: readonly SettingsKind[] = [
	'agent',
	'harness',
	'check',
	'checks_timeout',
	'import',
	'dashboard_token'
];

export type HostAction = { name: string; means: string; here: boolean; everywhere: boolean };
export type Harness = { program: string; args: string[]; set: boolean };
export type NamedCommand = { name: string; command: string[] };

/**
 * A project flai serve serves, or one below a folder named for import that it does not (S-0122):
 * from is registry, folder, or import; state is connected, connecting, not-connected,
 * unavailable, or not-running.
 */
export type ServedProject = {
	key: string;
	name: string;
	root: string;
	from?: 'registry' | 'folder' | 'import';
	/** The folder flai serve serves from that it is below, if any: removing it lists it as removed. */
	below?: string;
	state?: 'connected' | 'connecting' | 'not-connected' | 'unavailable' | 'not-running';
	since?: string;
	last_error?: string;
	/** Why it is unavailable, or why it is not served. */
	reason?: string;
	/** Not served because the operator removed it from below a folder; Serve brings it back (S-0123). */
	removed?: boolean;
	/** Whether the settings action is on for it, so that the dashboard may serve or remove it. */
	settings: boolean;
};

/**
 * When flai serve plans again on its own (S-0211, ADR-0084), read from system-flow.yaml's
 * planning and the plan host action; set by hand in the manifest, never from the page.
 */
export type PlanningTriggers = {
	/**
	 * Whether the plan host action is on for this project: every trigger needs it, and while it is
	 * on an edit to a planned story queues the planner.
	 */
	plan: boolean;
	/** planning.replan, or deterministic when it is unset: never, deterministic, or agent. */
	replan: string;
	/** Whether the manifest sets planning.replan. */
	replan_set: boolean;
	/** Why planning.replan cannot be read, when it cannot. */
	replan_error?: string;
	/** planning.schedule as written, a five-field cron expression in UTC or daily; empty when unset. */
	schedule: string;
	/** The schedule's next run, RFC 3339 in UTC, when it is set and can be read. */
	next?: string;
	/** Why planning.schedule cannot be read, when it cannot. */
	schedule_error?: string;
};

/** What a replan policy does when a story is accepted or cancelled or the pull order changes. */
export function replanMeans(policy: string): string {
	switch (policy) {
		case 'never':
			return 'nothing';
		case 'deterministic':
			return 'forecast deliveries follow, with no agent';
		case 'agent':
			return 'forecast deliveries follow, and the planner runs for each story whose delivery moved';
		default:
			return '';
	}
}

export type ProjectsView = {
	running: boolean;
	served: ServedProject[];
	unserved: ServedProject[];
	error?: string;
};

export type SettingsView = {
	/** Whether the dashboard may change this project's own settings. */
	here: boolean;
	/** Whether it may change the settings kept for every project. */
	everywhere: boolean;
	enable: string;
	enable_everywhere: string;
	host?: {
		actions: HostAction[];
		default_agent?: Agent | null;
		agent: {
			name: string;
			command: string[];
			harnesses: Record<string, Harness>;
		};
		checks: { commands: NamedCommand[]; timeout_minutes: number };
		manifest_checks?: NamedCommand[] | null;
		import_roots: string[];
		projects?: ProjectsView;
		mcp?: { running: boolean; url?: string; pid?: number };
		planning?: PlanningTriggers;
		/** The project's shared paths, claims.shared in its manifest (S-0295, ADR-0096), as globs. */
		shared?: string[];
		/** The strategic agents' settings in the manifest (S-0229), present when the project opens. */
		strategic?: StrategicSettings;
		error?: string;
	};
};

/** The blocks of the manifest that hold the strategic agents' settings, a key's first segment. */
export type StrategicBlock = 'orchestration' | 'planning' | 'analysis';

/** What a strategic setting's value is, and so how it is edited. */
export type StrategicKind =
	'boolean' | 'choice' | 'number' | 'duration' | 'cron' | 'agent' | 'text';

/** One key flai manifest set writes, with its value in the manifest (settings.get, S-0229). */
export type StrategicSetting = {
	key: string;
	kind: StrategicKind;
	/** The values a choice or a boolean takes, in order. */
	values?: string[];
	/** The value in effect when the key is unset, as it would be written. */
	default?: string;
	/** A sentence on what the setting does. */
	meaning: string;
	/** For a permission of the orchestrator's, what can go wrong while it is on. */
	risk?: string;
	/** A number that is a count, a whole number. */
	whole?: boolean;
	/** The value as the manifest has it; absent when unset. */
	value?: boolean | string | number | Agent;
	set: boolean;
};

/** settings.get's strategic block: each setting, and whether the dashboard may change them. */
export type StrategicSettings = {
	/** Whether the settings host action is on for this project. */
	editable: boolean;
	/** The command, run on the host in the project, that makes them editable. */
	enable: string;
	settings: StrategicSetting[];
};

/** One field flai refused, and why. */
export type Refusal = { field: string; reason: string };

/** An agent setting as its fields are typed: harness, model, and config one key=value a line. */
export type AgentForm = { harness: string; model: string; config: string };

/** A block's settings as typed: a boolean as true or false, an agent by its fields, the rest as text. */
export type StrategicForm = { values: Record<string, string>; agents: Record<string, AgentForm> };

/** What settings.manifest is asked: the keys to write with their values, and the keys to remove. */
export type ManifestChange = {
	set: Record<string, boolean | number | string | Agent>;
	unset: string[];
};

/** The settings of one block, in catalog order. */
export function inBlock(s: StrategicSettings, block: StrategicBlock): StrategicSetting[] {
	return s.settings.filter((x) => x.key.split('.')[0] === block);
}

/** The form as the manifest has the settings: an unset key is empty, so that its default shows. */
export function strategicForm(settings: StrategicSetting[]): StrategicForm {
	const form: StrategicForm = { values: {}, agents: {} };
	for (const s of settings) {
		if (s.kind === 'agent') {
			const a = (s.value ?? {}) as Agent;
			form.agents[s.key] = {
				harness: a.harness ?? '',
				model: a.model ?? '',
				config: configText(a.config)
			};
		} else if (s.kind === 'boolean') {
			form.values[s.key] = String(s.value ?? s.default === 'true');
		} else {
			form.values[s.key] = s.value === undefined || s.value === null ? '' : String(s.value);
		}
	}
	return form;
}

/**
 * The change from before to now: only the keys whose fields changed, an emptied one as an unset; or,
 * for what cannot be sent at all (a number that is not one, a config line that is not key=value),
 * each field and why.
 */
export function strategicChange(
	settings: StrategicSetting[],
	before: StrategicForm,
	now: StrategicForm
): ManifestChange | { refused: Refusal[] } {
	const change: ManifestChange = { set: {}, unset: [] };
	const refused: Refusal[] = [];
	for (const s of settings) {
		if (s.kind === 'agent') {
			const b = before.agents[s.key];
			const n = now.agents[s.key];
			if (!n || (b && b.harness === n.harness && b.model === n.model && b.config === n.config))
				continue;
			const parsed = parseConfig(n.config);
			if ('error' in parsed) {
				refused.push({ field: s.key, reason: `config: ${parsed.error}` });
				continue;
			}
			const roles = (s.value as Agent | undefined)?.roles;
			const agent = agentFrom(n.harness, n.model, parsed.config, roles);
			if (agent) change.set[s.key] = agent;
			else change.unset.push(s.key);
			continue;
		}
		const typed = (now.values[s.key] ?? '').trim();
		if (typed === (before.values[s.key] ?? '').trim()) continue;
		if (s.kind === 'boolean') change.set[s.key] = typed === 'true';
		else if (typed === '') change.unset.push(s.key);
		else if (s.kind === 'number') {
			const n = Number(typed);
			if (Number.isFinite(n)) change.set[s.key] = n;
			else refused.push({ field: s.key, reason: `"${typed}" is not a number` });
		} else change.set[s.key] = typed;
	}
	return refused.length ? { refused } : change;
}

/** Whether a change changes nothing. */
export function unchanged(c: ManifestChange): boolean {
	return Object.keys(c.set).length === 0 && c.unset.length === 0;
}

/** A command as the page edits it: one argument a line, so that nothing is quoted or split. */
export function lines(list: string[] | undefined): string {
	return (list ?? []).join('\n');
}

/** The argument list of a command written one argument a line; blank lines are dropped. */
export function argv(text: string): string[] {
	return text
		.split('\n')
		.map((l) => l.replace(/\r$/, ''))
		.filter((l) => l.trim() !== '');
}

/**
 * One entry flai shared check reports (S-0295): a path given, or an entry of a story's claim, and
 * the first shared path's pattern it lies wholly inside, if any.
 */
export type SharedCheck = {
	entry: string;
	/** The entry as it is matched: a component read as its path, without ./ or a trailing /. */
	path: string;
	shared: boolean;
	pattern?: string;
	/** The story whose claim the entry is, when a story was checked. */
	story?: string;
};

/** What settings.shared_check is asked for the text typed: a story's claim for its ID, else a path. */
export function sharedQuery(text: string): { story: string } | { paths: string[] } {
	const t = text.trim();
	return /^[Ss]-\d+$/.test(t) ? { story: t } : { paths: [t] };
}

/**
 * flai's refusal of a shared path without the manifest's absolute path it begins with, which says
 * nothing on the page: `/…/system-flow.yaml: cannot add "x" to claims.shared: …` reads from cannot.
 */
export function sharedRefusal(text: string): string {
	return text.replace(/^\S*system-flow\.yaml: /, '');
}

/** What a served project's state says, in a few words. */
export function health(p: ServedProject): string {
	switch (p.state) {
		case 'connected':
			return p.since ? `connected since ${p.since}` : 'connected';
		case 'not-connected':
			return `not connected: ${p.last_error ?? 'no answer yet'}`;
		case 'unavailable':
			return `not served: ${p.reason ?? 'unavailable'}`;
		case 'not-running':
			return 'flai serve is not running';
		default:
			return 'connecting';
	}
}

/** Whether the page may change a setting of this kind, and if not the command that allows it. */
export function allowed(view: SettingsView, kind: SettingsKind): { ok: boolean; enable: string } {
	return HOSTWIDE.includes(kind)
		? { ok: view.everywhere, enable: view.enable_everywhere }
		: { ok: view.here, enable: view.enable };
}
