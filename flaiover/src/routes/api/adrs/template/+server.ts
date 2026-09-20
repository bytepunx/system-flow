import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET: the body a new ADR starts from (flai's adr.template on the host). */
export const GET: RequestHandler = () =>
	respond(async () => (await repo().run<{ body: string }>('adr.template')).data);
