import { json } from '@sveltejs/kit';
import { repo, RepoError } from '$lib/server/repo';
import type { RequestHandler } from './$types';

/**
 * POST { include_uncommitted?, keep_removals? }: accept the item as the designer (flai's accept.run
 * on the host, S-0075) and stream what happens as newline-delimited JSON (S-0041): { event: "progress", step, msg }
 * for each step flai logs, { event: "warning", msg } for its warnings, then { event: "done", result }
 * or { event: "error", error, status, data? } with flai's message verbatim, and what flai sent with
 * it as data: a refusal for lint the fix left, { id, reason, findings }, each finding a file to
 * open and correct (ADR-0140, S-0378), or a refusal for a branch that undoes main's change,
 * { id, reason, undoes_main } (ADR-0155, S-0388). keep_removals, each { path, reason }, are the
 * removals of main's lines the designer keeps, with why; a list that is not one, or an entry with a
 * blank path or reason, is 400 before flai is asked. The steps arrive as progress
 * notifications over the channel while the acceptance runs. The HTTP status is 200 once the stream
 * starts; the last line says how it ended. When the connection to flai is lost on the way the
 * acceptance may still have completed on the host, and the error says how to find out.
 */
type Step = { level?: string; msg?: string; step?: string; detail?: string };
type Keep = { path: string; reason: string };

/** The removals a body keeps, or why they are not a list of { path, reason } with neither blank. */
function keepRemovals(value: unknown): Keep[] | string {
	if (value === undefined) return [];
	if (!Array.isArray(value))
		return 'keep_removals must be a list of { path, reason }: each a path the branch undoes of main’s change and why its removal is kept';
	const keeps: Keep[] = [];
	for (const entry of value) {
		const { path, reason } = (typeof entry === 'object' && entry !== null ? entry : {}) as {
			path?: unknown;
			reason?: unknown;
		};
		if (typeof path !== 'string' || !path.trim())
			return 'each of keep_removals needs a path: the file the branch undoes of main’s change';
		if (typeof reason !== 'string' || !reason.trim())
			return `the reason for keeping the removal at ${path} is blank: say why the branch’s removal of the main branch’s lines is kept`;
		keeps.push({ path, reason });
	}
	return keeps;
}

export const POST: RequestHandler = async ({ params, request }) => {
	const body = (await request.json().catch(() => ({}))) as {
		include_uncommitted?: boolean;
		keep_removals?: unknown;
	};
	const keeps = keepRemovals(body.keep_removals);
	if (typeof keeps === 'string') return json({ error: keeps }, { status: 400 });
	const encoder = new TextEncoder();
	const stream = new ReadableStream<Uint8Array>({
		async start(controller) {
			const send = (line: Record<string, unknown>) =>
				controller.enqueue(encoder.encode(JSON.stringify(line) + '\n'));
			try {
				const { data } = await repo().write<Record<string, unknown>>(
					'accept.run',
					{
						id: params.id,
						include_uncommitted: body.include_uncommitted === true,
						...(keeps.length ? { keep_removals: keeps } : {})
					},
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
				const data = e instanceof RepoError ? e.data : undefined;
				send({
					event: 'error',
					status,
					error:
						(e instanceof Error ? e.message : String(e)) +
						(lost
							? '. The acceptance may have completed on the host: reload the board, or run flai board there, before trying again.'
							: ''),
					...(data ? { data } : {})
				});
			}
			controller.close();
		}
	});
	return new Response(stream, {
		headers: { 'content-type': 'application/x-ndjson', 'cache-control': 'no-store' }
	});
};
