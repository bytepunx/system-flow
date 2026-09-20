import { connectedProjects } from '$lib/server/agent';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: every project flai on the host has ever named here, for a switcher (S-0080). No `?project=`
 * scopes this one, on purpose: it is how a page learns what projects exist at all, before it can
 * name one.
 */
export const GET: RequestHandler = () => respond(async () => ({ projects: connectedProjects() }));
