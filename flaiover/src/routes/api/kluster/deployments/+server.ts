import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/kluster/deployments: the namespaces deployed to the project's kluster, flai's
 * kluster.deployments on the host (S-0581), asked afresh each time, as { enabled, enable, block,
 * no_block?, status_error?, cluster?, default_namespace?, deployments, review, running?, last? }:
 * main's namespace first, then each story's, with its commit, outcome, steps and their logs, and
 * how to reach it for a demo. 501 when the host flai is older than kluster.deployments.
 */
export const GET: RequestHandler = () => respond(() => repo().klusterDeployments());
