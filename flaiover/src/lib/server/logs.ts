// A log on the host by its path, asked of flai's log.read (S-0785), for GET /api/log: each read is
// asked afresh, since a log grows without flai announcing it.

import type { LogRead } from '$lib/logs';
import { repo, RepoError } from './repo';

/** Which log to read, and from where: the query of GET /api/log. */
export type LogTarget = { path: string; offset: number };

/**
 * At most 64 KB of a log on the host from an offset, flai's log.read, so a page follows it live; a
 * host flai without the method told to upgrade (501).
 */
export async function logRead(t: LogTarget): Promise<LogRead> {
	try {
		return await repo().ask<LogRead>('log.read', t);
	} catch (e) {
		if (e instanceof RepoError && e.status === 501)
			throw new RepoError(
				501,
				`the host flai is too old to read a log (${e.message}): upgrade flai on the host with flai self-upgrade`
			);
		throw e;
	}
}

/**
 * The log GET /api/log reads, from its query: path a log's path as flai names it, absolute,
 * ~/-prefixed, or from the project's root, and offset a whole number of zero or more, 0 when left
 * out. A missing path or another offset is refused with 400 before flai is asked; what path flai
 * reads is flai's to say.
 */
export function logParam(query: URLSearchParams): LogTarget {
	const path = query.get('path');
	if (!path)
		throw new RepoError(
			400,
			'path must name the log to read, such as .flai-cache/kluster/logs/s-0595/build-api.log; open a log from a link in a message'
		);
	const offset = query.get('offset') ?? '0';
	if (!/^\d+$/.test(offset) || !Number.isSafeInteger(Number(offset)))
		throw new RepoError(
			400,
			`offset must be a whole number of bytes, 0 or more, not ${JSON.stringify(offset)}`
		);
	return { path, offset: Number(offset) };
}
