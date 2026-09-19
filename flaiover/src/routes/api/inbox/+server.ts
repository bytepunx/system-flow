import { repo } from '$lib/server/repo';
import { inbox } from '$lib/server/inbox';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET: what needs a human: threads, questions, reviews, blocked items, overlaps (S-0042). */
export const GET: RequestHandler = () => respond(() => inbox(repo()));
