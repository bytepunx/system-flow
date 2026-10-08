import { repo, type Conversation } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET ?story=<id>&all=1 : the conversations between stories (S-0336), the project's or one story's,
 * open ones only unless all=1. Open before closed, then the newest activity (updated) first, then by
 * ID.
 */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const story = url.searchParams.get('story');
		const all = url.searchParams.get('all') === '1';
		const r = repo();
		const list = story ? await r.messagesFor(story) : await r.messages();
		return list.filter((c) => all || !c.closed).sort(byActivity);
	});

/** Open before closed, then updated descending, then ID ascending. */
function byActivity(a: Conversation, b: Conversation): number {
	if (a.closed !== b.closed) return a.closed ? 1 : -1;
	if (a.updated !== b.updated) return a.updated < b.updated ? 1 : -1;
	return a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
}
