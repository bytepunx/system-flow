// The hosts a browser's Origin may name beside the request's own (ADR-0193): FLAIOVER_ALLOWED_HOSTS,
// which flai dashboard sets from dashboard.allowed_hosts, so that a page served through a tunnel or a
// domain that rewrites Host can still open its live events.

/** The port a browser leaves out of an Origin, by its scheme. */
const DEFAULT_PORT: Record<string, string> = { 'http:': '80', 'https:': '443' };

/** raw, a comma-separated list of hosts, each trimmed and lower-cased, the empty ones dropped. */
export function parseAllowedHosts(raw: string | undefined): string[] {
	return (raw ?? '')
		.split(',')
		.map((h) => h.trim().toLowerCase())
		.filter((h) => h !== '');
}

/** The hosts FLAIOVER_ALLOWED_HOSTS names in env; none when it is unset. */
export function allowedHosts(env: Record<string, string | undefined> = process.env): string[] {
	return parseAllowedHosts(env.FLAIOVER_ALLOWED_HOSTS);
}

/** An entry's host name and its port, '' when it gives none; an IPv6 address is in brackets. */
function hostAndPort(entry: string): [string, string] {
	const at = entry.lastIndexOf(':');
	if (at < 0 || entry.endsWith(']')) return [entry, ''];
	return [entry.slice(0, at), entry.slice(at + 1)];
}

/**
 * Whether origin names one of hosts: the same host name and the same port, where an entry with no port
 * means the default port of the origin's scheme, as an Origin with no port does.
 */
export function allowsOrigin(origin: URL, hosts: readonly string[]): boolean {
	const fallback = DEFAULT_PORT[origin.protocol] ?? '';
	const want = `${origin.hostname}:${origin.port || fallback}`;
	return hosts.some((h) => {
		const [name, port] = hostAndPort(h);
		return `${name}:${port || fallback}` === want;
	});
}
