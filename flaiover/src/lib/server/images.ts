// The project's agent images and image jobs, asked of flai's images methods on the host (S-0757,
// ADR-0295 §5), for the routes under /api/images. Each is asked afresh: a job's record changes in
// flai serve's state, which flai does not announce.

import type { ImageLog, ImageStarted, ImagesState } from '$lib/images';
import { repo, RepoError } from './repo';

/** Which job's log to read, and from where: the query of GET /api/images with job. */
export type ImageLogTarget = { id: string; offset: number };

/** Run fn, a host flai without the images methods told to upgrade. */
async function upgraded<T>(what: string, fn: () => Promise<T>): Promise<T> {
	try {
		return await fn();
	} catch (e) {
		if (e instanceof RepoError && e.status === 501)
			throw new RepoError(
				501,
				`the host flai is too old to ${what} (${e.message}): upgrade flai on the host with flai self-upgrade`
			);
		throw e;
	}
}

/**
 * The images the project's agents need with their state, and its image jobs, flai's images.get: a
 * read, whatever the agent host action says.
 */
export function imagesGet(): Promise<ImagesState> {
	return upgraded('show the agent images', () => repo().ask<ImagesState>('images.get'));
}

/** At most 64 KB of an image job's log from an offset, flai's images.log, so a page follows it live. */
export function imagesLog(t: ImageLogTarget): Promise<ImageLog> {
	return upgraded("read an image job's log", () => repo().ask<ImageLog>('images.log', t));
}

/**
 * A new image job that makes the project's agent images, flai's images.ensure under the agent host
 * action: queued in flai serve, answered at once with how many jobs wait or run before it.
 */
export async function imagesEnsure(): Promise<ImageStarted> {
	const { data } = await upgraded('start an image job', () =>
		repo().write<ImageStarted>('images.ensure')
	);
	return data;
}

/**
 * The job log GET /api/images reads, from its query: undefined with neither job nor offset, the
 * images and jobs then; else job an image job's id as flai serve issues it, such as
 * 20261113T093000Z-001, and offset a whole number of zero or more, 0 when left out. Anything else
 * is refused with 400 before flai is asked.
 */
export function imageLogParam(query: URLSearchParams): ImageLogTarget | undefined {
	const id = query.get('job');
	const raw = query.get('offset');
	if (id === null && raw === null) return undefined;
	if (id === null || !/^\d{8}T\d{6}Z-\d{3,}$/.test(id))
		throw new RepoError(
			400,
			`job must be an image job's id GET /api/images lists, such as 20261113T093000Z-001, not ${JSON.stringify(id)}; leave job and offset out to list the images and jobs`
		);
	const offset = raw ?? '0';
	if (!/^\d+$/.test(offset) || !Number.isSafeInteger(Number(offset)))
		throw new RepoError(
			400,
			`offset must be a whole number of bytes, 0 or more, not ${JSON.stringify(offset)}`
		);
	return { id, offset: Number(offset) };
}
