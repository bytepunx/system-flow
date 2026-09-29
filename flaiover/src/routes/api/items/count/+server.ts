import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';

/** GET: { active, archived }, how many items there are of each, without the archive (S-0162). */
export const GET = () => respond(() => repo().itemCount());
