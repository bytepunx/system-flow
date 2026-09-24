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

/** Whether the page may change a setting of this kind, and if not the command that allows it. */
export function allowed(view: SettingsView, kind: SettingsKind): { ok: boolean; enable: string } {
	return HOSTWIDE.includes(kind)
		? { ok: view.everywhere, enable: view.enable_everywhere }
		: { ok: view.here, enable: view.enable };
}
