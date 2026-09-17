import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
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
