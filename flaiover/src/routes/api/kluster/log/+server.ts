import { klusterLogParam, repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/kluster/log?namespace=&step=&offset=: one step's log of a namespace from offset, flai's
 * kluster.log on the host (ADR-0263), asked afresh each time, as { namespace, step, log, exists,
 * offset, next, size, text, reset?, outcome? }: at most 64 KB of text, and next, the offset to ask
 * from next, so the Deployments page follows a running step live. A query that is not a namespace,
 * a step, and an offset of zero or more is refused with 400 before flai is asked; 501 when the host
 * flai is older than kluster.log.
 */
export const GET: RequestHandler = ({ url }) =>
	respond(() => repo().klusterLog(klusterLogParam(url.searchParams)));
