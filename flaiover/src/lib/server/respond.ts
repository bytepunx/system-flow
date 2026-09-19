import { json } from '@sveltejs/kit';
import { RepoError } from './repo';

/** Run a handler and map RepoError to an HTTP status with a JSON body. */
export async function respond<T>(fn: () => Promise<T>): Promise<Response> {
	try {
		return json(await fn());
	} catch (e) {
		if (e instanceof RepoError)
			return json({ ...(e.data ?? {}), error: e.message }, { status: e.status });
		throw e;
	}
}
