import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** The token file's repository path as design.get names it; 404 when the project has none. */
async function tokenPath(): Promise<string> {
	const design = await repo().designSystem();
	if (!design.exists)
		throw new RepoError(
			404,
			design.no_token_file || `there is no design token file at ${design.path}`
		);
	return design.path;
}

/** Run fn, a host flai without design.get, file.get, or file.save told to upgrade (ADR-0147). */
async function upgraded<T>(what: string, fn: () => Promise<T>): Promise<T> {
	try {
		return await fn();
	} catch (e) {
		if (e instanceof RepoError && e.status === 501)
			throw new RepoError(
				501,
				`the host flai is too old to ${what} the design token file (${e.message}): upgrade flai on the host with flai self-upgrade`
			);
		throw e;
	}
}

/**
 * GET /api/design/tokens: the design token file, raw, through flai's file.get on the host (ADR-0147,
 * S-0383), at the path design.get names: { path, content, hash }, hash its SHA-256. Asked afresh each
 * time. 404 when the project has no token file; 501 when the host flai is older than file.get.
 */
export const GET: RequestHandler = () =>
	respond(() => upgraded('read', async () => repo().file(await tokenPath())));

type SaveBody = { path?: unknown; content?: unknown; hash?: unknown; message?: unknown };

/**
 * PUT /api/design/tokens { path, content, hash, message? }: flai's file.save on the host, on the main
 * checkout, the content as data (ADR-0147, S-0383); only the token file design.get names, else 400.
 * 409 { error, current, hash, diff } when the file changed since it was loaded; 422 { error, findings }
 * when flai refuses the content, such as a design.tokens finding; 501 when the host flai is older
 * than file.save. flai holds the rules; this handler passes through.
 */
export const PUT: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as SaveBody;
		const { path, content, hash, message } = body;
		if (
			typeof path !== 'string' ||
			!path ||
			typeof content !== 'string' ||
			typeof hash !== 'string' ||
			!hash
		)
			throw new RepoError(400, 'path, content, and hash are required');
		return upgraded('save', async () => {
			const want = (await repo().designSystem()).path;
			if (path !== want)
				throw new RepoError(
					400,
					`this endpoint saves only the design token file ${want}, not ${path}`
				);
			const { data, warnings } = await repo().saveFile({
				path,
				content,
				hash,
				message: typeof message === 'string' && message ? message : undefined
			});
			return { ...data, log: warnings };
		});
	});
