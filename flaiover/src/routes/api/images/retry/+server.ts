import type { ImagesRetried } from '$lib/images';
import { imagesEnsure, imagesGet } from '$lib/server/images';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/images/retry: a new image job that makes the project's agent images, flai's
 * images.ensure under the agent host action (S-0757), then the images and jobs, as GET /api/images
 * answers them, with the job it started as started, and ahead, how many jobs wait or run before it.
 * 403 { error, action, enable } while the action is off; 400 { error, job } with flai's sentence,
 * which names the job, while a job that makes the same images waits or runs; 400 with flai's
 * sentence when the project's agents need no image flai makes; 501 when the host flai is older than
 * images.ensure.
 */
export const POST: RequestHandler = () =>
	respond(async (): Promise<ImagesRetried> => {
		const started = await imagesEnsure();
		return { ...(await imagesGet()), started };
	});
