// Pure helpers for the document editor (S-0040). The editor keeps front
// matter as raw text so that a document it did not change round-trips
// byte for byte; flai compares content by hash (ADR-0023).

/** A document split into its raw front matter (without the --- lines, null when there is none) and body. */
export function splitRaw(doc: string): { frontMatter: string | null; body: string } {
	if (!doc.startsWith('---\n')) return { frontMatter: null, body: doc };
	const rest = doc.slice(4);
	const end = rest.indexOf('\n---\n');
	if (end < 0) return { frontMatter: null, body: doc };
	return { frontMatter: rest.slice(0, end + 1), body: rest.slice(end + 5) };
}

/** The inverse of splitRaw. */
export function compose(frontMatter: string | null, body: string): string {
	if (frontMatter === null) return body;
	const fm = frontMatter.endsWith('\n') ? frontMatter : frontMatter + '\n';
	return `---\n${fm}---\n${body}`;
}

/**
 * The heading whose section holds the caret: the last markdown heading at or
 * before that position, ignoring lines inside fenced code. Null above the
 * first heading.
 */
export function headingAt(body: string, caret: number): string | null {
	let found: string | null = null;
	let offset = 0;
	let fenced = false;
	for (const line of body.split('\n')) {
		if (offset > caret) break;
		if (/^\s*(```|~~~)/.test(line)) fenced = !fenced;
		else if (!fenced) {
			const m = /^#{1,6}\s+(.+?)\s*#*\s*$/.exec(line);
			if (m) found = m[1];
		}
		offset += line.length + 1;
	}
	return found;
}

/** Headings of a document body, in order, as the explorer lists them for threads. */
export function headingsOf(body: string): string[] {
	const out: string[] = [];
	let fenced = false;
	for (const line of body.split('\n')) {
		if (/^\s*(```|~~~)/.test(line)) fenced = !fenced;
		else if (!fenced) {
			const m = /^#{1,6}\s+(.+?)\s*#*\s*$/.exec(line);
			if (m) out.push(m[1]);
		}
	}
	return out;
}
