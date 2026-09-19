import { repo, RepoError } from '$lib/server/repo';
import { flaiStream } from '$lib/server/flai';
import { designer } from '$lib/server/threads';
import type { RequestHandler } from './$types';

/**
 * POST { include_uncommitted? }: accept the item as the designer (flai accept <id> --by <owner>) and
 * stream what happens as newline-delimited JSON (S-0041): { event: "progress", step, msg } for each
 * step flai logs, { event: "warning", msg } for its warnings, then { event: "done", result } or
 * { event: "error", error, status } with flai's message verbatim. The HTTP status is 200 once the
 * stream starts; the last line says how it ended. include_uncommitted is the designer's choice to
 * put uncommitted files outside wip into the acceptance commit (S-0051).
 */
export const POST: RequestHandler = async ({ params, request }) => {
	const body = (await request.json().catch(() => ({}))) as { include_uncommitted?: boolean };
	const args = ['accept', params.id, '--by', await designer()];
	if (body.include_uncommitted === true) args.push('--yes');
	const encoder = new TextEncoder();
	const stream = new ReadableStream<Uint8Array>({
		async start(controller) {
			const send = (line: Record<string, unknown>) =>
				controller.enqueue(encoder.encode(JSON.stringify(line) + '\n'));
			try {
				const { data } = await flaiStream<Record<string, unknown>>(repo().root, args, (e) => {
					if (e.msg === 'acceptance step')
						send({ event: 'progress', step: e.step, msg: String(e.detail ?? e.step) });
					else if (e.level === 'WARN') send({ event: 'warning', msg: String(e.detail ?? e.msg) });
				});
				send({ event: 'done', result: data });
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
