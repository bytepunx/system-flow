// Who works a story (S-0103): the harness that runs it, the model, and options for the harness.
// The project's default is in system-flow.yaml; a story carries its own copy, which it overrides.

export type Agent = { harness?: string; model?: string; config?: Record<string, string> };

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

/** The agent the fields give, or undefined when they give nothing. */
export function agentFrom(
	harness: string,
	model: string,
	config: Record<string, string>
): Agent | undefined {
	const a: Agent = {};
	if (harness.trim()) a.harness = harness.trim();
	if (model.trim()) a.model = model.trim();
	if (Object.keys(config).length) a.config = config;
	return Object.keys(a).length ? a : undefined;
}

/** Whether two agents say the same thing. */
export function sameAgent(a?: Agent, b?: Agent): boolean {
	return (
		(a?.harness ?? '') === (b?.harness ?? '') &&
		(a?.model ?? '') === (b?.model ?? '') &&
		configText(a?.config) === configText(b?.config)
	);
}

/** The agent in one line: claude-code, claude-opus-5-5, effort=high. */
export function agentLine(a?: Agent): string {
	if (!a) return 'none';
	const parts = [a.harness, a.model].filter(Boolean) as string[];
	for (const k of Object.keys(a.config ?? {}).sort()) parts.push(`${k}=${a.config![k]}`);
	return parts.length ? parts.join(', ') : 'none';
}
