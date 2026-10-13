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

/** The offset of the start of a 1-based line of text, clamped to the first and the last line. */
export function offsetOfLine(text: string, line: number): number {
	let offset = 0;
	for (let n = 1; n < line; n++) {
		const next = text.indexOf('\n', offset);
		if (next < 0) break;
		offset = next + 1;
	}
	return offset;
}

/** The 1-based line of text that holds an offset. */
export function lineAt(text: string, offset: number): number {
	let line = 1;
	const end = Math.min(Math.max(offset, 0), text.length);
	for (let i = 0; i < end; i++) if (text[i] === '\n') line++;
	return line;
}

/** Where a line of the whole file falls in the editor: in the front matter, or in the body, as a 1-based line of that part. */
export type Located = { part: 'frontMatter' | 'body'; line: number };

/**
 * A 1-based line of the file, as a link or a lint finding names it, located in
 * the parts splitRaw makes: the opening and closing --- lines count as the
 * front matter's first and last line.
 */
export function locate(frontMatter: string | null, fileLine: number): Located {
	if (frontMatter === null) return { part: 'body', line: Math.max(fileLine, 1) };
	const fmLines = lineAt(frontMatter, frontMatter.length) - (frontMatter.endsWith('\n') ? 1 : 0);
	if (fileLine <= fmLines + 2)
		return { part: 'frontMatter', line: Math.min(Math.max(fileLine - 1, 1), Math.max(fmLines, 1)) };
	return { part: 'body', line: fileLine - fmLines - 2 };
}

/** What a save that flai accepted answers, for whoever opened the editor. */
export type Saved = { path: string; hash?: string; commit?: string; committed?: boolean };
