import type { ImageLog, ImagesState } from '$lib/images';
import { imageLogParam, imagesGet, imagesLog } from '$lib/server/images';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/images: the images the project's agents need with their state, and its image jobs,
 * flai's images.get on the host (S-0757), asked afresh each time, as { flai, images, unjudged?,
 * running, queued, ended }. GET /api/images?job=&offset=: one job's log from offset, flai's
 * images.log, as { id, state, outcome?, exists, offset, next, size, text, reset? }: at most 64 KB of
 * text, and next, the offset to ask from next, so the page follows a running job live. Both are
 * reads, whatever the agent host action says. A job that is not a job's id, or an offset that is not
 * a whole number of zero or more, is refused with 400 before flai is asked; 404 when the project has
 * no such job; 501 when the host flai is older than the images methods.
 */
export const GET: RequestHandler = ({ url }) =>
	respond((): Promise<ImagesState | ImageLog> => {
		const log = imageLogParam(url.searchParams);
		return log ? imagesLog(log) : imagesGet();
	});
