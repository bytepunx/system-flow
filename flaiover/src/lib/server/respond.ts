import { json } from '@sveltejs/kit';
import { repo, RepoError } from './repo';
import { projectIdentity } from './project';

/** A JSON object gains `project: { name, key }` (ADR-0024); arrays and scalars keep their shape. */
async function identified(body: unknown): Promise<unknown> {
	if (body === null || typeof body !== 'object' || Array.isArray(body)) return body;
	return { ...(body as Record<string, unknown>), project: await projectIdentity(repo()) };
}

/** Run a handler and map RepoError to an HTTP status with a JSON body. */
export async function respond<T>(fn: () => Promise<T>): Promise<Response> {
	try {
		return json(await identified(await fn()));
	} catch (e) {
		if (e instanceof RepoError)
			return json(await identified({ ...(e.data ?? {}), error: e.message }), { status: e.status });
		throw e;
	}
}
