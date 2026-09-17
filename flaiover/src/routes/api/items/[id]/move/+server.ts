import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** POST { to, reason?, by? } */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as {
			to?: string;
			reason?: string;
			by?: string;
		};
		if (!body.to) throw new RepoError(400, 'to is required');
		const args = ['move', params.id, body.to];
		if (body.reason) args.push('--reason', body.reason);
		if (body.by) args.push('--by', body.by);
		const { data, warnings } = await flai<{ id: string; status: string; warnings: string[] }>(
			repo().root,
			args
		);
		return { ...data, warnings: [...(data?.warnings ?? []), ...warnings] };
	});
