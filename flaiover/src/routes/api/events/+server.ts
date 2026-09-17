import { repo } from '$lib/server/repo';
import type { RequestHandler } from './$types';

/** Server-sent events: one "change" event per changed file. */
export const GET: RequestHandler = async ({ request }) => {
	const r = repo();
	await r.watch();
	const stream = new ReadableStream({
		start(controller) {
			const enc = new TextEncoder();
			const send = (path: string) =>
				controller.enqueue(enc.encode(`event: change\ndata: ${JSON.stringify({ path })}\n\n`));
			controller.enqueue(enc.encode(`event: ready\ndata: {}\n\n`));
			r.on('change', send);
			const ping = setInterval(() => controller.enqueue(enc.encode(': ping\n\n')), 25000);
			request.signal.addEventListener('abort', () => {
				clearInterval(ping);
				r.off('change', send);
				controller.close();
			});
		}
	});
	return new Response(stream, {
		headers: {
			'content-type': 'text/event-stream',
			'cache-control': 'no-cache',
			connection: 'keep-alive'
		}
	});
};
