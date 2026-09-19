import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET /api/docs/file?path=design/system/overview.md */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const path = url.searchParams.get('path');
		if (!path) throw new RepoError(400, 'path is required');
		return repo().docFile(path);
	});

/** The trailer that names the dashboard on commits it makes for the designer (ADR-0023). */
const TRAILER = 'Co-Authored-By: flaiover <flaiover@localhost>';

type SaveBody = { path?: string; content?: string; hash?: string; message?: string };

/**
 * PUT /api/docs/file { path, content, hash, message? } : flai doc save with the content on stdin.
 * 409 { error, current, hash, diff } when the file changed since it was loaded; 422 { error, findings }
 * when flai refuses the content. flai holds the rules; this handler passes through (ADR-0023).
 */
export const PUT: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as SaveBody;
		if (!body.path || typeof body.content !== 'string' || !body.hash)
			throw new RepoError(400, 'path, content, and hash are required');
		const args = ['doc', 'save', body.path, '--hash', body.hash, '--trailer', TRAILER];
		if (body.message?.trim()) args.push('--message', body.message.trim());
		try {
			const { data, warnings } = await flai<Record<string, unknown>>(repo().root, args, {
				input: body.content,
				exitStatus: { 3: 409, 4: 422 }
			});
			return { ...data, log: warnings };
		} catch (e) {
			// flai wraps the payload as { conflict: {...} } or { refused: {...} }; the response is flat.
			if (e instanceof RepoError && e.data) {
				const inner = (e.data.conflict ?? e.data.refused) as Record<string, unknown> | undefined;
				throw new RepoError(e.status, e.message, inner ?? {});
			}
			throw e;
		}
	});
