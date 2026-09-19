import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import { designer } from '$lib/server/threads';
import { NATURE_NAMES } from '$lib/natures';
import type { RequestHandler } from './$types';

/** GET /api/items?type=story&status=review&archived=false */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const type = url.searchParams.get('type');
		const status = url.searchParams.get('status');
		const archived = url.searchParams.get('archived');
		let items = await repo().items();
		if (type) items = items.filter((it) => it.type === type);
		if (status) items = items.filter((it) => it.status === status);
		if (archived === 'true' || archived === 'false')
			items = items.filter((it) => it.archived === (archived === 'true'));
		return items.map((it) => {
			const { body, ...rest } = it;
			void body;
			return rest;
		});
	});

/** The trailer that names the dashboard on commits it makes for the designer (ADR-0023). */
const TRAILER = 'Co-Authored-By: flaiover <flaiover@localhost>';

const ID = /^[EST]-\d+$/;

type NewBody = {
	type?: string;
	title?: string;
	nature?: string;
	parent?: string;
	tags?: string[];
	touches?: string[];
	body?: string;
};

/** One tag or path: no commas (flai splits on them), no control characters, not empty. */
const listValue = (v: unknown): v is string =>
	typeof v === 'string' && v.trim() !== '' && !/[,\n\r\0]/.test(v) && v.length <= 200;

/**
 * The flai arguments for creating an item. Exported for tests. Only epics and stories are created
 * here: tasks are written by the agent that pulls a story (ADR-0021). Values go to flai as
 * `--flag=value` and the title after `--`, so nothing the designer types can be read as a flag.
 */
export function _newArgs(body: NewBody, owner: string): string[] {
	if (body.type !== 'epic' && body.type !== 'story')
		throw new RepoError(
			400,
			'type must be epic or story; tasks are written by the agent that pulls a story'
		);
	const title = (body.title ?? '').replace(/\s+/g, ' ').trim();
	if (!title) throw new RepoError(400, 'a title is required');
	if (title.length > 200) throw new RepoError(400, 'the title is longer than 200 characters');
	const nature = body.nature ?? 'feature';
	if (!NATURE_NAMES.includes(nature))
		throw new RepoError(400, `nature must be one of ${NATURE_NAMES.join(', ')}`);
	if (typeof body.body !== 'string' || body.body.trim() === '')
		throw new RepoError(400, 'write what the item is for: the body is empty');
	const args = [body.type, 'new', `--nature=${nature}`, `--owner=${owner}`];
	if (body.type === 'story') {
		if (!body.parent || !ID.test(body.parent) || !body.parent.startsWith('E-'))
			throw new RepoError(400, 'a story needs its parent epic');
		args.push(`--epic=${body.parent}`);
	} else if (body.parent) {
		throw new RepoError(400, 'an epic has no parent');
	}
	for (const [flag, values] of [
		['tag', body.tags],
		['touches', body.touches]
	] as const) {
		for (const v of values ?? []) {
			if (!listValue(v))
				throw new RepoError(400, `${flag} values are single words or paths without commas`);
			args.push(`--${flag}=${v.trim()}`);
		}
	}
	args.push('--body-stdin', '--autocommit', `--trailer=${TRAILER}`, '--', title);
	return args;
}

/**
 * POST { type, title, nature?, parent?, tags?, touches?, body }: create an epic or a story with the
 * designer's markdown as its body (S-0059). flai does it as one step (ADR-0016): the item from the
 * project's template with the next ID, linked into its parent, checked with it in place, committed
 * on its own. 422 { error, findings } when the check refuses it, and then nothing was created.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as NewBody;
		const args = _newArgs(body, await designer());
		try {
			const { data, warnings } = await flai<Record<string, unknown>>(repo().root, args, {
				input: body.body,
				exitStatus: { 4: 422 }
			});
			return { ...data, log: warnings };
		} catch (e) {
			if (e instanceof RepoError && e.data?.refused)
				throw new RepoError(e.status, e.message, e.data.refused as Record<string, unknown>);
			throw e;
		}
	});
