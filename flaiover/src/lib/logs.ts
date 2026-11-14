// A log on the host as flai's log.read answers it (S-0785), for GET /api/log and the log page,
// /ops/log, that shows one, and the log paths in free text from flai that link to that page.
import { resolve } from '$app/paths';

/**
 * A read of a log on the host from an offset (Go's hostapi.LogRead), as KlusterLog is of a step's:
 * at most 64 KB of text from offset, and next, the offset to ask from next, so a page follows a
 * growing log live: GET /api/log?path=&offset=.
 */
export type LogRead = {
	/** The path as it was asked for. */
	path: string;
	/** False while there is no file at the path. */
	exists: boolean;
	/** Where text starts in the log. */
	offset: number;
	/** The offset to ask from next: offset and the bytes of text. */
	next: number;
	/** The log's size when it was read. */
	size: number;
	text: string;
	/** The log was shorter than the offset asked, as when it was written anew; text starts at 0. */
	reset?: boolean;
};

/** A run of plain text, or a log path, in text from flai. */
export type LogSegment = { text: string } | { path: string };

// A candidate starts with /, ~/, or .flai-cache/, not inside a longer path, word, or URL, and runs
// to the next whitespace, quote, or backtick; it is a log path when, without trailing punctuation,
// it ends in .log.
const CANDIDATE = /(?<![\w./~:-])(?:\/(?!\/)|~\/|\.flai-cache\/)[^\s'"`]*/g;
const TRAILING = /[.,;:)]+$/;

/**
 * Split text into plain runs and the log paths it names: an absolute, ~/, or .flai-cache/ path
 * ending in .log, cut at whitespace, quotes, and backticks, without trailing punctuation.
 */
export function splitLogPaths(text: string): LogSegment[] {
	const out: LogSegment[] = [];
	let at = 0;
	for (const m of text.matchAll(CANDIDATE)) {
		const path = m[0].replace(TRAILING, '');
		if (!path.endsWith('.log')) continue;
		if (m.index > at) out.push({ text: text.slice(at, m.index) });
		out.push({ path });
		at = m.index + path.length;
	}
	if (at < text.length) out.push({ text: text.slice(at) });
	return out;
}

/** The log page's href for a log path, with the project the tab views when one is chosen. */
export function logHref(path: string, project?: string | null): string {
	const query = new URLSearchParams({ path });
	if (project) query.set('project', project);
	return `${resolve('/ops/log')}?${query}`;
}
