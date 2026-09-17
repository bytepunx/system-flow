import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { isSecure, sessionCookie, tokenMatches } from '$lib/server/auth';

/** Exchange the token (from the login page's URL fragment or a pasted value) for the session cookie. */
export const POST: RequestHandler = async ({ request, url }) => {
	let body: { token?: unknown };
	try {
		body = await request.json();
	} catch {
		return json({ error: 'unauthorized' }, { status: 401 });
	}
	const token = typeof body.token === 'string' ? body.token.trim() : '';
	if (!tokenMatches(token)) return json({ error: 'unauthorized' }, { status: 401 });
	return json(
		{ ok: true },
		{ headers: { 'set-cookie': sessionCookie(token, isSecure(url, request.headers)) } }
	);
};
