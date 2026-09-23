import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * A repository offered for import (S-0098): `?project=` names the candidate the page picked.
 *
 * GET: what an import would do, changing nothing (flai import --dry-run on the host): the layout,
 * the folders it would move, the sub-projects it found, and the tests it would run.
 *
 * POST: import it (flai import --yes --commit on the host), streamed as newline-delimited JSON as
 * an acceptance is: { event: "progress", msg } for each test as it starts and ends, then
 * { event: "done", result, committed } with flai's answer, or { event: "error", error, status }.
 * committed is false when a test failed: the import was applied, nothing was committed, and the
 * result says which test failed and what it printed. Either way the repository is a project from
 * then on, under the key in result.key.
 */
export const GET: RequestHandler = () =>
	respond(async () => (await repo().run<Record<string, unknown>>('import.preview')).data);

type Step = { level?: string; msg?: string; name?: string; component?: string; seconds?: number };

export const POST: RequestHandler = async () => {
	const encoder = new TextEncoder();
	const stream = new ReadableStream<Uint8Array>({
		async start(controller) {
			const send = (line: Record<string, unknown>) =>
				controller.enqueue(encoder.encode(JSON.stringify(line) + '\n'));
			send({ event: 'progress', msg: 'importing' });
			try {
				const { data } = await repo().write<Record<string, unknown>>(
					'import.run',
					{},
					{
						timeoutMs: 2 * 60 * 60 * 1000,
						onProgress: (value) => {
							const e = (value ?? {}) as Step;
							if (e.component !== 'import' || !e.name) return;
							if (e.msg === 'running tests') send({ event: 'progress', msg: `running ${e.name}` });
							else if (e.msg === 'tests passed')
								send({ event: 'progress', msg: `${e.name} passed` });
							else if (e.msg === 'tests failed')
								send({ event: 'progress', msg: `${e.name} failed` });
						}
					}
				);
				send({ event: 'done', result: data, committed: true });
			} catch (e) {
				if (e instanceof RepoError && e.status === 422 && e.data?.commit) {
					send({ event: 'done', result: e.data, committed: false });
				} else {
					const status = e instanceof RepoError ? e.status : 500;
					send({ event: 'error', status, error: e instanceof Error ? e.message : String(e) });
				}
			}
			controller.close();
		}
	});
	return new Response(stream, {
		headers: { 'content-type': 'application/x-ndjson', 'cache-control': 'no-store' }
	});
};
