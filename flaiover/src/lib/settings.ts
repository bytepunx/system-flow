// The host's settings as the dashboard's settings page sees them (S-0105), flai's settings.get.
import type { Agent } from '$lib/agent';

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
	'dashboard_token'
] as const;
export type SettingsKind = (typeof SETTINGS_KINDS)[number];

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
		error?: string;
	};
};

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
