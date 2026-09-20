import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { title, status?, supersedes?, refines?, body }: record an architecture decision with the
 * designer's markdown as its body (S-0060; flai's adr.new on the host, S-0075). flai does it as one
 * step (ADR-0016): number, file, front matter, index row, superseded_by, the check, and one commit.
 * 400 when an argument is not what it should be; 422 { error, findings } when the check refuses it,
 * and then nothing was created.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		const { data, warnings } = await repo().write<Record<string, unknown>>('adr.new', {
			title: body.title,
			status: body.status,
			supersedes: body.supersedes,
			refines: body.refines,
			body: body.body
		});
		return { ...data, log: warnings };
	});
