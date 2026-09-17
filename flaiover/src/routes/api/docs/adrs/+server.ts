import { repo } from '$lib/server/repo';
import { adrs } from '$lib/server/search';
import { respond } from '$lib/server/respond';

export const GET = () => respond(() => adrs(repo()));
