// Pure helpers for the review page (S-0041).

export type Criterion = { text: string; checked: boolean };

/** The checkboxes under "## Acceptance criteria" in a story body. */
export function criteriaOf(body: string): Criterion[] {
	const section = sectionOf(body, 'Acceptance criteria');
	const out: Criterion[] = [];
	for (const line of section.split('\n')) {
		const m = /^\s*[-*]\s+\[( |x|X)\]\s+(.*\S)\s*$/.exec(line);
		if (m) out.push({ checked: m[1] !== ' ', text: m[2] });
	}
	return out;
}

/** The text under a level-two heading, up to the next one; empty when the heading is missing. */
export function sectionOf(markdown: string, heading: string): string {
	const lines = markdown.split('\n');
	const start = lines.findIndex((l) => l.trim().toLowerCase() === `## ${heading}`.toLowerCase());
	if (start < 0) return '';
	const rest = lines.slice(start + 1);
	const end = rest.findIndex((l) => /^##\s/.test(l));
	return (end < 0 ? rest : rest.slice(0, end)).join('\n').trim();
}

export type PatchLine = { kind: 'hunk' | 'add' | 'del' | 'context' | 'note'; text: string };

/** A unified patch as lines with their kind, so a view can colour them and still show + and -. */
export function patchLines(patch: string): PatchLine[] {
	if (!patch) return [];
	return patch.split('\n').map((text) => ({
		text,
		kind: text.startsWith('@@')
			? 'hunk'
			: text.startsWith('+')
				? 'add'
				: text.startsWith('-')
					? 'del'
					: text.startsWith('\\')
						? 'note'
						: 'context'
	}));
}

/**
 * Read a newline-delimited JSON response as it arrives, calling `each` for every complete line.
 * Lines that are not JSON are skipped; a last line without a newline is still delivered.
 */
export async function readNdjson(
	response: Response,
	each: (line: Record<string, unknown>) => void
): Promise<void> {
	const emit = (raw: string) => {
		const text = raw.trim();
		if (!text) return;
		try {
			each(JSON.parse(text));
		} catch {
			// not a JSON line
		}
	};
	if (!response.body) {
		for (const l of (await response.text()).split('\n')) emit(l);
		return;
	}
	const reader = response.body.getReader();
	const decoder = new TextDecoder();
	let pending = '';
	for (;;) {
		const { value, done } = await reader.read();
		if (done) break;
		pending += decoder.decode(value, { stream: true });
		const lines = pending.split('\n');
		pending = lines.pop() ?? '';
		lines.forEach(emit);
	}
	emit(pending + decoder.decode());
}

/**
 * The body with its nth acceptance criterion ticked or unticked (S-0085). Only lines of the
 * "Acceptance criteria" section count, in order, as criteriaOf lists them; null when there is no
 * such criterion. Everything else in the body is left exactly as it was.
 */
export function toggleCriterion(body: string, index: number): string | null {
	const lines = body.split('\n');
	const start = lines.findIndex((l) => l.trim().toLowerCase() === '## acceptance criteria');
	if (start < 0 || index < 0) return null;
	let n = -1;
	for (let i = start + 1; i < lines.length && !/^##\s/.test(lines[i]); i++) {
		const m = /^(\s*[-*]\s+\[)( |x|X)(\]\s+.*\S\s*)$/.exec(lines[i]);
		if (!m) continue;
		if (++n === index) {
			lines[i] = m[1] + (m[2] === ' ' ? 'x' : ' ') + m[3];
			return lines.join('\n');
		}
	}
	return null;
}
