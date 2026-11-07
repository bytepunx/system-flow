import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/kluster: the project's kluster, flai's kluster.status on the host (S-0580), asked afresh
 * each time, as { enabled, enable, block, no_block?, cluster?, status_error?, running?, last? }.
 * It answers whether the kluster host action is on or off. 501 when the host flai is older than
 * kluster.status.
 */
export const GET: RequestHandler = () => respond(() => repo().klusterStatus());
