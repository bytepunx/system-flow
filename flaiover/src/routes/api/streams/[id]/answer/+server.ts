import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { question, answer }: answer a hand-written open question in the
 * story's narrative (flai's stream.answer, S-0090). It moves from Open
 * questions to Decisions, with the answer recorded beside it, and so stops
 * showing in the inbox. `question` must read exactly as the inbox entry's
 * own title, so flai can find the bullet to remove.
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as {
			question?: string;
			answer?: string;
		};
		if (!body.question?.trim()) throw new RepoError(400, 'question is required');
		if (!body.answer?.trim()) throw new RepoError(400, 'answer is required');
		return (
			await repo().write('stream.answer', {
				id: params.id,
				question: body.question,
				answer: body.answer
			})
		).data;
	});
