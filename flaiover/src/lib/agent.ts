// Who works a story (S-0103): the harness that runs it, the model, and options for the harness.
// The project's default is in system-flow.yaml; a story carries its own copy, which it overrides.
// Roles (S-0189) are the agents for the work a story's agent hands to sub-agents, explore and
// verify; the dashboard shows them and keeps them on a save, and flai sets them.

export type Role = { harness?: string; model?: string; config?: Record<string, string> };
export type Agent = Role & { roles?: Record<string, Role> };

/** The config as the form shows it: one key=value a line, keys in order. */
export function configText(config?: Record<string, string>): string {
	return Object.keys(config ?? {})
		.sort()
		.map((k) => `${k}=${config![k]}`)
		.join('\n');
}

/** The config the form's text gives, or the line that is not key=value. */
export function parseConfig(text: string): { config: Record<string, string> } | { error: string } {
	const config: Record<string, string> = {};
	for (const raw of text.split('\n')) {
		const line = raw.trim();
		if (!line) continue;
		const at = line.indexOf('=');
		if (at <= 0) return { error: `"${line}" is not key=value` };
		config[line.slice(0, at).trim()] = line.slice(at + 1).trim();
	}
	return { config };
}

/** The agent the fields give, with the roles it had, or undefined when they give nothing. */
export function agentFrom(
	harness: string,
	model: string,
	config: Record<string, string>,
	roles?: Record<string, Role>
): Agent | undefined {
	const a: Agent = {};
	if (harness.trim()) a.harness = harness.trim();
	if (model.trim()) a.model = model.trim();
	if (Object.keys(config).length) a.config = config;
	if (roles && Object.keys(roles).length) a.roles = roles;
	return Object.keys(a).length ? a : undefined;
}

function sameRole(a?: Role, b?: Role): boolean {
	return (
		(a?.harness ?? '') === (b?.harness ?? '') &&
		(a?.model ?? '') === (b?.model ?? '') &&
		configText(a?.config) === configText(b?.config)
	);
}

/** Whether two agents say the same thing, roles included. */
export function sameAgent(a?: Agent, b?: Agent): boolean {
	const names = new Set([...Object.keys(a?.roles ?? {}), ...Object.keys(b?.roles ?? {})]);
	return sameRole(a, b) && [...names].every((n) => sameRole(a?.roles?.[n], b?.roles?.[n]));
}

function roleLine(r: Role): string {
	const parts = [r.harness, r.model].filter(Boolean) as string[];
	for (const k of Object.keys(r.config ?? {}).sort()) parts.push(`${k}=${r.config![k]}`);
	return parts.join(', ');
}

/** The agent in one line: claude-code, claude-opus-5-5, effort=high; explore: haiku; verify: sonnet. */
export function agentLine(a?: Agent): string {
	if (!a) return 'none';
	const roles = Object.keys(a.roles ?? {})
		.sort()
		.map((n) => `${n}: ${roleLine(a.roles![n])}`);
	const own = roleLine(a);
	if (!own && !roles.length) return 'none';
	return [own || 'default', ...roles].join('; ');
}
