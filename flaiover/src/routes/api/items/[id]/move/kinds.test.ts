import { afterEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask, type AskOptions } from '$lib/server/repo';
import { baselineFlow, renamedFlow } from '$lib/lanes.fixture';
import type { Flow } from '$lib/lanes';
import { POST } from './+server';

type Handler = (event: never) => Promise<Response>;
type Asked = { method: string; params?: Record<string, unknown>; opt?: AskOptions };

// I-0399: the route told a cancellation and an acceptance by the state's name, so on a workflow
// whose cancelled or accepted lane has another name either was a plain move. It reads the lane's
// kind from the board's flow instead (S-0838).
function fake(flow: Flow): Asked[] {
	const asked: Asked[] = [];
	const answers: Record<string, unknown> = {
		'project.info': {
			name: 'Harbour',
			key: 'harbour',
			layout: { design: 'design', docs: 'docs', wip: 'wip' }
		},
		'board.get': { ...flow, columns: {}, wip_limits: {}, order: [] },
		'item.move': { data: { id: 'S-0001', status: 'moved' }, warnings: [] },
		'item.move.preview': {
			data: { id: 'S-0001', status: 'cancelled', cancelled: [] },
			warnings: []
		}
	};
	const ask = (async (method: string, params?: Record<string, unknown>, opt?: AskOptions) => {
		asked.push({ method, params, opt });
		return answers[method];
	}) as Ask;
	useRepo(new Repo('/nowhere', ask));
	return asked;
}

const move = (body: unknown) =>
	(POST as unknown as Handler)({
		params: { id: 'S-0001' },
		request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
	} as never);

/** The asks other than reading the board and the project. */
const writes = (asked: Asked[]) =>
	asked.filter((a) => a.method !== 'board.get' && a.method !== 'project.info');

describe("a move is what its lane's kind says, not its name (I-0399)", () => {
	afterEach(() => useRepo(null));

	for (const [name, flow, accepted, cancelled, other] of [
		['renamed', renamedFlow, 'shipped', 'dropped', 'queued'],
		['baseline', baselineFlow, 'done', 'cancelled', 'ready']
	] as const) {
		describe(`on the ${name} workflow`, () => {
			it(`previews a dry run to ${cancelled}, the cancelled lane, without moving`, async () => {
				const asked = fake(flow);
				const r = await move({ to: cancelled, dry_run: true });
				expect(r.status).toBe(200);
				expect(writes(asked).map((a) => a.method)).toEqual(['item.move.preview']);
			});
			it(`refuses a dry run to ${other}, which is no cancellation, and moves nothing`, async () => {
				const asked = fake(flow);
				expect((await move({ to: other, dry_run: true })).status).toBe(400);
				expect(writes(asked)).toEqual([]);
			});
			it(`gives a move to ${accepted}, the accepted lane, the acceptance's timeout`, async () => {
				const asked = fake(flow);
				expect((await move({ to: accepted })).status).toBe(200);
				const [w] = writes(asked);
				expect(w.method).toBe('item.move');
				expect(w.params?.to).toBe(accepted);
				expect(w.opt?.timeoutMs).toBe(600000);
			});
			it(`gives a move to ${other} a plain move's timeout`, async () => {
				const asked = fake(flow);
				expect((await move({ to: other })).status).toBe(200);
				const [w] = writes(asked);
				expect(w.method).toBe('item.move');
				expect(w.opt?.timeoutMs).toBe(60000);
			});
		});
	}
});
