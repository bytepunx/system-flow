import { repo } from '$lib/server/repo';
import { sseStream } from '$lib/server/sse';
import type { RequestHandler } from './$types';

/** Server-sent events: one "change" event per changed file. */
export const GET: RequestHandler = async ({ request }) => {
	const r = repo();
	await r.watch();
	return new Response(sseStream(r, request.signal), {
		headers: {
			'content-type': 'text/event-stream',
			'cache-control': 'no-cache',
			connection: 'keep-alive'
		}
	});
};
