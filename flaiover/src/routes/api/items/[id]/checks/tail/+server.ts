import { repo, RepoError } from '$lib/server/repo';
import type { RequestHandler } from './$types';

/**
 * GET ?from=N: new output from a story's checks run since byte offset N,
 * streamed as newline-delimited JSON (S-0082, the same shape accept's own
 * progress stream uses): { event: "line", text } for each new line as
 * flai's own checks.tail (a read, progress: true) finds it, then { event:
 * "done", offset, running, outcome? } once flai's own wait (about 20s)
 * elapses or the run ends. Meant to be called again immediately with the
 * returned offset for as long as the page wants to keep watching; closing
 * the tab loses nothing; the run lives on the host, not in this request.
 */
type Line = { msg?: string; text?: string };

export const GET: RequestHandler = async ({ params, url }) => {
	const from = Math.max(0, Number(url.searchParams.get('from') ?? '0') || 0);
	const encoder = new TextEncoder();
	const stream = new ReadableStream<Uint8Array>({
		async start(controller) {
			const send = (line: Record<string, unknown>) =>
				controller.enqueue(encoder.encode(JSON.stringify(line) + '\n'));
			try {
				const { data } = await repo().run<{ offset: number; running: boolean; outcome?: string }>(
					'checks.tail',
					{ id: params.id, from },
					{
						timeoutMs: 30000,
						onProgress: (value) => {
							const e = (value ?? {}) as Line;
							if (e.msg === 'line') send({ event: 'line', text: e.text ?? '' });
						}
					}
				);
				send({ event: 'done', offset: data.offset, running: data.running, outcome: data.outcome });
			} catch (e) {
				const status = e instanceof RepoError ? e.status : 500;
				send({ event: 'error', status, error: e instanceof Error ? e.message : String(e) });
			}
			controller.close();
		}
	});
	return new Response(stream, {
		headers: { 'content-type': 'application/x-ndjson', 'cache-control': 'no-store' }
	});
};
