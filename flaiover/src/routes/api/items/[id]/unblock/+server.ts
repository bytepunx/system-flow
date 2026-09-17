import { repo } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = ({ params }) =>
	respond(async () => (await flai(repo().root, ['unblock', params.id])).data);
