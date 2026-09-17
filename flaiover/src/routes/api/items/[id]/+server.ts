import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = ({ params }) => respond(() => repo().itemById(params.id));
