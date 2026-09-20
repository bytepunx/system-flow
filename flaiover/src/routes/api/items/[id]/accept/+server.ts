import { repo, RepoError } from '$lib/server/repo';
import type { RequestHandler } from './$types';

/**
 * POST { include_uncommitted? }: accept the item as the designer (flai's accept.run on the host,
 * S-0075) and stream what happens as newline-delimited JSON (S-0041): { event: "progress", step, msg }
 * for each step flai logs, { event: "warning", msg } for its warnings, then { event: "done", result }
 * or { event: "error", error, status } with flai's message verbatim. The steps arrive as progress
 * notifications over the channel while the acceptance runs. The HTTP status is 200 once the stream
 * starts; the last line says how it ended. When the connection to flai is lost on the way the
 * acceptance may still have completed on the host, and the error says how to find out.
 */
type Step = { level?: string; msg?: string; step?: string; detail?: string };

export const POST: RequestHandler = async ({ params, request }) => {
	const body = (await request.json().catch(() => ({}))) as { include_uncommitted?: boolean };
	const encoder = new TextEncoder();
	const stream = new ReadableStream<Uint8Array>({
		async start(controller) {
			const send = (line: Record<string, unknown>) =>
				controller.enqueue(encoder.encode(JSON.stringify(line) + '\n'));
			try {
				const { data } = await repo().write<Record<string, unknown>>(
					'accept.run',
					{ id: params.id, include_uncommitted: body.include_uncommitted === true },
					{
						timeoutMs: 600000,
						onProgress: (value) => {
							const e = (value ?? {}) as Step;
							if (e.msg === 'acceptance step')
								send({ event: 'progress', step: e.step, msg: String(e.detail ?? e.step) });
							else if (e.level === 'WARN')
								send({ event: 'warning', msg: String(e.detail ?? e.msg) });
						}
					}
				);
				send({ event: 'done', result: data });
			} catch (e) {
				const status = e instanceof RepoError ? e.status : 500;
				const lost = status === 502 || status === 504;
				send({
					event: 'error',
					status,
					error:
						(e instanceof Error ? e.message : String(e)) +
						(lost
							? '. The acceptance may have completed on the host: reload the board, or run flai board there, before trying again.'
							: '')
				});
			}
			controller.close();
		}
	});
	return new Response(stream, {
		headers: { 'content-type': 'application/x-ndjson', 'cache-control': 'no-store' }
	});
};
