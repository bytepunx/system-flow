// Browser-side fetch for the dashboard API: sends the CSRF header on writes
// (the session cookie is HttpOnly and travels on its own) and returns to the
// login page when the session is missing or was rotated away.
export async function api(path: string, init: RequestInit = {}): Promise<Response> {
	const method = (init.method ?? 'GET').toUpperCase();
	const headers = new Headers(init.headers);
	if (method !== 'GET' && method !== 'HEAD') headers.set('x-requested-with', 'flaiover');
	const r = await fetch(path, { ...init, headers });
	if (r.status === 401 && typeof location !== 'undefined' && location.pathname !== '/login') {
		// A full navigation, so page state built with a dead session is dropped.
		location.assign(`/login?next=${encodeURIComponent(location.pathname + location.search)}`);
	}
	return r;
}
