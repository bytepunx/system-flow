import { repo } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: the sections the project's design/adrs/0000-template.md gives a new ADR
 * (`flai adr new --print-body`), which the form starts from (S-0060). It answers only when flai
 * is available, so the ADRs page also uses it to know whether to offer "new ADR".
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		const { data } = await flai<{ body: string }>(repo().root, ['adr', 'new', '--print-body']);
		return data;
	});
