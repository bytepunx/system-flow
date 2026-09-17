import { repo } from '$lib/server/repo';
import { board } from '$lib/server/board';
import { respond } from '$lib/server/respond';

export const GET = () => respond(() => board(repo()));
