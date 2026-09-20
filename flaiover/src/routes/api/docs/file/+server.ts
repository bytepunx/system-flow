import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET /api/docs/file?path=design/system/overview.md */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const path = url.searchParams.get('path');
		if (!path) throw new RepoError(400, 'path is required');
		return repo().docFile(path);
	});

type SaveBody = { path?: string; content?: string; hash?: string; message?: string };

/**
 * PUT /api/docs/file { path, content, hash, message? }: flai's doc.save on the host, the content as
 * data and never a file the dashboard writes (ADR-0023, S-0075). 409 { error, current, hash, diff }
 * when the file changed since it was loaded; 422 { error, findings } when flai refuses the content.
 * flai holds the rules; this handler passes through.
 */
export const PUT: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as SaveBody;
		if (!body.path || typeof body.content !== 'string' || !body.hash)
			throw new RepoError(400, 'path, content, and hash are required');
		const { data, warnings } = await repo().write<Record<string, unknown>>('doc.save', {
			path: body.path,
			content: body.content,
			hash: body.hash,
			message: body.message
		});
		return { ...data, log: warnings };
	});
