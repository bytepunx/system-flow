import { json } from '@sveltejs/kit';

/** Liveness: the process is up. Never touches the repository. */
export const GET = () => json({ status: 'ok' });
