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

export type PatchKind = 'hunk' | 'add' | 'del' | 'context' | 'note';
/** One line of a patch: `sign` is its + or -, empty for the rest, and `text` is the line without it. */
export type PatchLine = { kind: PatchKind; sign: '+' | '-' | ''; text: string };

/**
 * A unified patch as lines with their kind, the sign of an added or removed line apart from its
 * text, so a view can put it in a margin. A hunk header and a "no newline" note keep their whole
 * text; the empty line the patch's last newline leaves is not a line.
 */
export function patchLines(patch: string): PatchLine[] {
	if (!patch) return [];
	const lines = patch.split('\n');
	if (lines[lines.length - 1] === '') lines.pop();
	return lines.map((raw): PatchLine => {
		if (raw.startsWith('@@')) return { kind: 'hunk', sign: '', text: raw };
		if (raw.startsWith('\\')) return { kind: 'note', sign: '', text: raw };
		if (raw.startsWith('+')) return { kind: 'add', sign: '+', text: raw.slice(1) };
		if (raw.startsWith('-')) return { kind: 'del', sign: '-', text: raw.slice(1) };
		// a context line begins with the space that stands where a sign would
		return { kind: 'context', sign: '', text: raw.startsWith(' ') ? raw.slice(1) : raw };
	});
}

/** Consecutive lines of one kind: a view outlines a run of added or of removed lines as one block. */
export type PatchRun = { kind: PatchKind; lines: PatchLine[] };

/** A unified patch as runs: every line beside its neighbours of the same kind, in the patch's order. */
export function patchRuns(patch: string): PatchRun[] {
	const runs: PatchRun[] = [];
	for (const line of patchLines(patch)) {
		const last = runs[runs.length - 1];
		if (last && last.kind === line.kind) last.lines.push(line);
		else runs.push({ kind: line.kind, lines: [line] });
	}
	return runs;
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

/** Where a step of flai verify got to (S-0270). */
export type VerifyState = 'passed' | 'failed' | 'not-reached';
/** What a failed step found: a test, a check, or a linter's complaint, where it is, and what it says. */
export type VerifyFinding = { name?: string; path?: string; line?: number; message?: string };
/** One step of flai verify: one of its own, or a test or lint tier when `tier` is set. */
export type VerifyStep = {
	name: string;
	tier?: boolean;
	state: VerifyState;
	duration_ms: number;
	duration?: string;
	findings?: VerifyFinding[];
	omitted?: number;
	command?: string[];
	dir?: string;
	exit_code?: number;
};
/** A finding of flai check outside the story: it does not fail the story (S-0249). */
export type VerifyNote = {
	rule: string;
	level: string;
	path?: string;
	line?: number;
	message: string;
};
/** flai verify's report for a story, as it stores the last one and verify.status answers it. */
export type VerifyReport = {
	story: string;
	commit: string;
	base: string;
	ran_at: string;
	duration_ms: number;
	duration?: string;
	passed: boolean;
	stopped_at?: string;
	paths?: string[];
	steps: VerifyStep[];
	notes?: VerifyNote[];
};

/** A finding or a note as a page shows it: where, as path:line and a name, and what it says. */
export type VerifyLine = { where: string; message: string };
/** A step as a row: its state in words, its duration, and the findings when it failed. */
export type VerifyRow = {
	name: string;
	tier: boolean;
	state: VerifyState;
	label: string;
	duration: string;
	findings: VerifyLine[];
	/** How many findings flai left out past its cap. */
	omitted: number;
};
/** A report as a page shows it: the outcome, the commit and when, a row per step, the notes apart. */
export type VerifyView = {
	passed: boolean;
	outcome: string;
	commit: string;
	base: string;
	ranAt: string;
	duration: string;
	rows: VerifyRow[];
	notes: VerifyLine[];
};

const STATE_LABEL: Record<VerifyState, string> = {
	passed: 'passed',
	failed: 'failed',
	'not-reached': 'not reached'
};

/** A duration in milliseconds as flai writes one when the report carries none: 340ms, 5.1s. */
function took(ms: number): string {
	return ms < 1000 ? `${Math.round(ms)}ms` : `${(ms / 1000).toFixed(1)}s`;
}

/** Where a finding or note is, as flai verify prints it: path:line, then its name. */
function where(path?: string, line?: number, name?: string): string {
	const loc = path ? (line ? `${path}:${line}` : path) : '';
	return [loc, name ?? ''].filter((s) => s).join(' ');
}

/**
 * A report as rows for the review page (S-0270): a row per step in the order flai ran them, with
 * the findings under the step that failed only; a step not reached has no duration. The notes from
 * outside the story are apart, since they do not fail it.
 */
export function verifyView(report: VerifyReport): VerifyView {
	const rows = report.steps.map((s): VerifyRow => ({
		name: s.name,
		tier: s.tier === true,
		state: s.state,
		label: STATE_LABEL[s.state] ?? s.state,
		duration: s.state === 'not-reached' ? '' : (s.duration ?? took(s.duration_ms)),
		findings:
			s.state === 'failed'
				? (s.findings ?? []).map((f) => ({
						where: where(f.path, f.line, f.name),
						message: f.message ?? ''
					}))
				: [],
		omitted: s.state === 'failed' ? (s.omitted ?? 0) : 0
	}));
	return {
		passed: report.passed,
		outcome: report.passed
			? 'Passed every step'
			: report.stopped_at
				? `Stopped at ${report.stopped_at}`
				: 'Failed',
		commit: report.commit.slice(0, 12),
		base: report.base,
		ranAt: report.ran_at,
		duration: report.duration ?? took(report.duration_ms),
		rows,
		notes: (report.notes ?? []).map((n) => ({
			where: where(n.path, n.line, n.rule),
			message: `${n.level}: ${n.message}`
		}))
	};
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
