import type { LogRead } from '$lib/logs';
import { logParam, logRead } from '$lib/server/logs';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/log?path=&offset=: a log on the host from offset, flai's log.read (S-0785), as
 * { path, exists, offset, next, size, text, reset? }: at most 64 KB of text, and next, the offset to
 * ask from next, so the log page follows a growing log live. A read, whatever the host actions say.
 * A missing path, or an offset that is not a whole number of zero or more, is refused with 400
 * before flai is asked; flai's refusal of a path it does not read keeps its words; 501 when the
 * host flai is older than the method.
 */
export const GET: RequestHandler = ({ url }) =>
	respond((): Promise<LogRead> => logRead(logParam(url.searchParams)));
