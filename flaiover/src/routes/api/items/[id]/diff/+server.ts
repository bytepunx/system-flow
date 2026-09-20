import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET: what the story's branch changes against main (flai's stream.diff on the host), for review. */
export const GET: RequestHandler = ({ params }) =>
	respond(
		async () =>
			(
				await repo().run<Record<string, unknown>>(
					'stream.diff',
					{ id: params.id },
					{ timeoutMs: 60000 }
				)
			).data
	);
