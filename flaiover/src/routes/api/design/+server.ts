import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/design?story=S-nnnn: the project's design system, flai's design.get on the host (S-0382):
 * its tokens per theme, their problems, the CSS and Tailwind renderings, and each colour role pair's
 * contrast; with no token file, exists false and no_token_file naming flai design init. Kept until a
 * document changes. With story, the design system in that open story's worktree (S-0384), asked
 * afresh each time: 404 when no item has the ID, 400 when it is malformed or names no open story
 * with a worktree. 501 when the host flai is older than design.get; 500 when flai cannot read the files.
 */
export const GET: RequestHandler = ({ url }) =>
	respond(() => repo().designSystem(url.searchParams.get('story') || undefined));
