// Which scheme a request arrived over, for adapter-node, which cannot tell
// by itself: with no protocol header configured it reports https for every
// request, and then a session cookie set over plain HTTP is marked Secure,
// which a browser keeps on http://localhost and drops on every other address
// (S-0083). The header is the de facto one a TLS-terminating proxy or tunnel
// sends. When nothing sent it, the connection itself says.

/** The header adapter-node is told to read; PROTOCOL_HEADER overrides it. */
export const schemeHeader = () =>
	(process.env.PROTOCOL_HEADER || 'x-forwarded-proto').toLowerCase();

/**
 * Leave the request with exactly `http` or `https` in the scheme header: what
 * a proxy said, when it said one of the two (the first, when it names a
 * chain), else what the connection is. Anything else a client put there is
 * replaced, so it never reaches adapter-node, which refuses odd values with
 * an error per request. A client that claims https over plain HTTP only
 * affects its own session: its cookie is marked Secure and it loses it.
 * @param {{ headers: Record<string, string | string[] | undefined>, socket?: { encrypted?: boolean } }} req
 * @returns {string}
 */
export function noteScheme(req) {
	const name = schemeHeader();
	const given = req.headers[name];
	const first = String(Array.isArray(given) ? given[0] : (given ?? ''))
		.split(',')[0]
		.trim()
		.toLowerCase();
	req.headers[name] =
		first === 'https' || first === 'http' ? first : req.socket?.encrypted ? 'https' : 'http';
	return req.headers[name];
}
