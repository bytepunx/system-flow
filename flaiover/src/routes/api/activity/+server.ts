import { repo } from '$lib/server/repo';
import { activity } from '$lib/server/activity';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET: active streams with their agent, the task in progress, and the last log entry (S-0042). */
export const GET: RequestHandler = () => respond(() => activity(repo()));
