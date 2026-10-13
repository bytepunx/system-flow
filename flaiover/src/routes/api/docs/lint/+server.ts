import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

type LintBody = { path?: unknown; content?: unknown };

/**
 * POST /api/docs/lint { path, content }: flai's doc.lint on the host, the markdown lint findings of
 * the content as the editor holds it (S-0378): { path, findings: [{ line, rule, alias, message }] },
 * findings [] when clean. 400 when path or content is missing or flai refuses the path; 501 when the
 * host flai is older than doc.lint, and the editor shows no live findings.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as LintBody;
		if (typeof body.path !== 'string' || !body.path || typeof body.content !== 'string')
			throw new RepoError(400, 'path and content are required');
		return repo().docLint(body.path, body.content);
	});
