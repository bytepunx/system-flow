import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { clearCookie } from '$lib/server/auth';

export const POST: RequestHandler = async () =>
	json({ ok: true }, { headers: { 'set-cookie': clearCookie() } });
