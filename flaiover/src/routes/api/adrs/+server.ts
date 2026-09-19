import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** The trailer that names the dashboard on commits it makes for the designer (ADR-0023). */
const TRAILER = 'Co-Authored-By: flaiover <flaiover@localhost>';

type NewAdr = {
	title?: string;
	status?: string;
	supersedes?: unknown[];
	refines?: unknown[];
	body?: string;
};

/** An ADR reference as the list gives it (ADR-0007) or as a number; flai gets the number. */
function adrNumber(v: unknown): string {
	const m = /^(?:ADR-)?(\d{1,4})$/.exec(String(v).trim());
	if (!m || Number(m[1]) === 0) throw new RepoError(400, `${String(v)} is not an ADR`);
	return String(Number(m[1]));
}

/**
 * The flai arguments for recording an ADR. Exported for tests. Values go to flai as
 * `--flag=value` and the title after `--`, so nothing the designer types is read as a flag.
 */
export function _adrArgs(body: NewAdr): string[] {
	const title = (body.title ?? '').replace(/\s+/g, ' ').trim();
	if (!title) throw new RepoError(400, 'a title is required: the decision, as a sentence');
	if (title.length > 200) throw new RepoError(400, 'the title is longer than 200 characters');
	const status = body.status ?? 'proposed';
	if (status !== 'proposed' && status !== 'accepted')
		throw new RepoError(400, 'status must be proposed or accepted');
	if (typeof body.body !== 'string' || body.body.trim() === '')
		throw new RepoError(400, 'write the decision: the body is empty');
	const args = ['adr', 'new', `--status=${status}`];
	for (const n of body.supersedes ?? []) args.push(`--supersedes=${adrNumber(n)}`);
	for (const n of body.refines ?? []) args.push(`--refines=${adrNumber(n)}`);
	args.push('--body-stdin', '--autocommit', `--trailer=${TRAILER}`, '--', title);
	return args;
}

/** The flai arguments for accepting a proposed ADR. Exported for tests. */
export function _acceptArgs(id: string): string[] {
	return ['adr', 'accept', adrNumber(id), '--autocommit', `--trailer=${TRAILER}`];
}

/**
 * POST { title, status?, supersedes?, refines?, body }: record an architecture decision with the
 * designer's markdown as its body (S-0060). flai does it as one step (ADR-0016): number, file,
 * front matter, index row, superseded_by, the check, and one commit. 422 { error, findings } when
 * the check refuses it, and then nothing was created.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as NewAdr;
		try {
			const { data, warnings } = await flai<Record<string, unknown>>(repo().root, _adrArgs(body), {
				input: body.body,
				exitStatus: { 4: 422 }
			});
			return { ...data, log: warnings };
		} catch (e) {
			if (e instanceof RepoError && e.data?.refused)
				throw new RepoError(e.status, e.message, e.data.refused as Record<string, unknown>);
			throw e;
		}
	});
