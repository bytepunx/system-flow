import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';
import type { Agent } from '$lib/agent';

/**
 * An item's own words, changed after it was made (S-0085): title, nature, tags, touches, a story's
 * after (S-0130), parent, and the body below its heading. flai does it on the host (item.show, item.edit, which are
 * `flai edit`): a retitle is kept in step in the file's name, the parent's list, the narrative, and
 * links; the repository is checked with the change in place; one commit holds every file touched.
 * What is the item's state (ID, type, status, transitions, blocks, owner, dates) is not reachable
 * from here.
 */

export type ItemView = {
	id: string;
	type: 'epic' | 'story' | 'task';
	status: string;
	title: string;
	nature: string;
	tags: string[];
	touches: string[];
	/** The stories a story waits for until they are done (S-0130). */
	after?: string[];
	parent?: string;
	/** A story's agent, and the project's default a new story would get (S-0103). */
	agent?: Agent;
	default_agent?: Agent;
	body: string;
	path: string;
	hash: string;
	editable: boolean;
	reason?: string;
	natures: string[];
	parents: { id: string; title: string }[];
};

/** GET: the item as an editor loads it, with the hash a save must carry. */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => (await repo().run<ItemView>('item.show', { id: params.id })).data);

const FIELDS = ['title', 'nature', 'tags', 'touches', 'after', 'parent', 'body'] as const;

/**
 * PUT { hash, title?, nature?, tags?, touches?, after?, parent?, agent?, body? }: change what is given and leave the
 * rest. 409 { error, current, hash } when the item changed after it was read; 422 { error, findings }
 * when flai check refuses the change, and then nothing was changed; 400 for a value that is not
 * what it should be.
 */
export const PUT: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		if (typeof body.hash !== 'string' || !body.hash)
			throw new RepoError(400, 'hash is required: the one the item was loaded with');
		const change: Record<string, unknown> = { id: params.id, hash: body.hash };
		for (const f of FIELDS) if (body[f] !== undefined && body[f] !== null) change[f] = body[f];
		// a story's agent replaces what it has; null removes it (S-0103)
		if (body.agent !== undefined) change.agent = body.agent;
		if (Object.keys(change).length === 2) throw new RepoError(400, 'nothing to change');
		const { data, warnings } = await repo().write<Record<string, unknown>>('item.edit', change);
		return { ...data, log: warnings };
	});
