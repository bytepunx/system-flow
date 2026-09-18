import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import { designer } from '$lib/server/threads';
import type { RequestHandler } from './$types';

/** GET ?on=<path|item>&all=1 : threads, unresolved unless all=1. */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const on = url.searchParams.get('on');
		const all = url.searchParams.get('all') === '1';
		const r = repo();
		const list = on ? await r.threadsFor(on) : await r.threads();
		return all ? list : list.filter((t) => t.status !== 'resolved');
	});

/** POST { on, heading?, title, text } : open a thread through flai. */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as {
			on?: string;
			heading?: string;
			title?: string;
			text?: string;
		};
		if (!body.on || !body.title?.trim() || !body.text?.trim())
			throw new RepoError(400, 'on, title, and text are required');
		const args = ['thread', 'new', '--on', body.on, '--by', await designer()];
		if (body.heading) args.push('--heading', body.heading);
		args.push(body.title.trim(), body.text.trim());
		const { data, warnings } = await flai<Record<string, unknown>>(repo().root, args);
		return { ...data, warnings };
	});
