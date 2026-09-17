import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';

export const GET = () => respond(() => repo().manifest());
