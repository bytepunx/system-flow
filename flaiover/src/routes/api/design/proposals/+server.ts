import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/design/proposals: the open stories whose branches change the design system, flai's
 * design.proposals on the host (S-0384): [{ id, title, status, paths }], [] when none. Asked afresh
 * each time. 501 when the host flai is older than design.proposals; 500 when flai cannot read them.
 */
export const GET: RequestHandler = () => respond(() => repo().designProposals());
