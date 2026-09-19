import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/docs/edit?path= : a document as an editor loads it (flai doc show): content, its hash, and
 * the edit mode flai allows, full, body, or none with the reason (ADR-0023).
 */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const path = url.searchParams.get('path');
		if (!path) throw new RepoError(400, 'path is required');
		const { data } = await flai<Record<string, unknown>>(repo().root, ['doc', 'show', path]);
		return data;
	});
