import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/kluster/workloads: the project's cluster's nodes and each namespace's pods, flai's
 * kluster.workloads on the host (S-0586), asked afresh each time, as { block, no_block?, cluster?,
 * context?, nodes, namespaces, workloads_error? }. 501 when the host flai is older than
 * kluster.workloads.
 */
export const GET: RequestHandler = () => respond(() => repo().klusterWorkloads());
