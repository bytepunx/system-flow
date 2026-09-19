// Which project an answer came from (ADR-0024): a front end that spans
// projects routes and labels by it without guessing.
import type { Repo } from './repo';

export type ProjectIdentity = { name: string; key: string };

export const PROJECT_KEY_HEADER = 'x-flai-project-key';
export const PROJECT_NAME_HEADER = 'x-flai-project-name';

/** The manifest's name and key; empty strings when the manifest cannot be read. */
export async function projectIdentity(repo: Repo): Promise<ProjectIdentity> {
	try {
		const m = (await repo.manifest()) as { name?: unknown; key?: unknown };
		return { name: String(m.name ?? ''), key: String(m.key ?? '') };
	} catch {
		return { name: '', key: '' };
	}
}

/** Name the project on a response, whatever its body: arrays, streams, and errors included. */
export function setIdentityHeaders(response: Response, id: ProjectIdentity): Response {
	try {
		response.headers.set(PROJECT_KEY_HEADER, id.key);
		// header values are bytes; a name may hold anything
		response.headers.set(PROJECT_NAME_HEADER, encodeURIComponent(id.name));
		return response;
	} catch {
		// a response with immutable headers: copy it
		const headers = new Headers(response.headers);
		headers.set(PROJECT_KEY_HEADER, id.key);
		headers.set(PROJECT_NAME_HEADER, encodeURIComponent(id.name));
		return new Response(response.body, {
			status: response.status,
			statusText: response.statusText,
			headers
		});
	}
}
