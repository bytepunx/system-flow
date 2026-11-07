import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** The kluster config file's repository path as kluster.status names it; 404 when there is none. */
async function configPath(): Promise<string> {
	const status = await repo().klusterStatus();
	if (!status.block || !status.config)
		throw new RepoError(404, 'this project has no kluster: block naming a kluster config file');
	return status.config;
}

/**
 * GET /api/kluster/config: the kluster config file the manifest's kluster: block names, raw, through
 * flai's file.get on the host (ADR-0147, S-0586): { path, content, hash }, hash its SHA-256. Asked
 * afresh each time. 404 with no block; 501 when the host flai is older than file.get.
 */
export const GET: RequestHandler = () => respond(async () => repo().file(await configPath()));

type SaveBody = { path?: unknown; content?: unknown; hash?: unknown; message?: unknown };

/**
 * PUT /api/kluster/config { path, content, hash, message? }: flai's file.save on the host, on the
 * main checkout (ADR-0147, S-0586); only the config file kluster.status names, else 400. 409
 * { error, current, hash, diff } when the file changed since it was loaded; 422 { error, findings }
 * when flai refuses the content. flai holds the rules; this handler passes through.
 */
export const PUT: RequestHandler = ({ request }) =>
	respond(async () => {
		const { path, content, hash, message } = (await request.json().catch(() => ({}))) as SaveBody;
		if (
			typeof path !== 'string' ||
			!path ||
			typeof content !== 'string' ||
			typeof hash !== 'string'
		)
			throw new RepoError(400, 'path, content, and hash are required');
		const want = await configPath();
		if (path !== want)
			throw new RepoError(
				400,
				`this endpoint saves only the kluster config file ${want}, not ${path}`
			);
		const { data, warnings } = await repo().saveFile({
			path,
			content,
			hash,
			message: typeof message === 'string' && message ? message : undefined
		});
		return { ...data, log: warnings };
	});
