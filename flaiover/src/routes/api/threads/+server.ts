import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET ?on=<path|item>&all=1 : threads, unresolved unless all=1, each entry marked operator when the operator wrote it. */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const on = url.searchParams.get('on');
		const all = url.searchParams.get('all') === '1';
		const r = repo();
		const [list, operator] = await Promise.all([on ? r.threadsFor(on) : r.threads(), r.operator()]);
		return (all ? list : list.filter((t) => t.status !== 'resolved')).map((t) => ({
			...t,
			entries: (t.entries ?? []).map((e) => ({ ...e, operator: e.author === operator }))
		}));
	});

/** POST { on, heading?, title, text }: open a thread as the designer (flai's thread.new on the host). */
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
		const { data, warnings } = await repo().write<Record<string, unknown>>('thread.new', {
			on: body.on,
			heading: body.heading,
			title: body.title,
			text: body.text
		});
		return { ...data, warnings };
	});
