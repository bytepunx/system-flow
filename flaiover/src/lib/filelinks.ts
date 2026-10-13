// Repository file references in plain text (S-0378): a refusal, a report, or a thread names a
// markdown file as `path`, `path:line`, or `path:line:column`. These find them under the layout's
// design, docs, and wip folders, tell whether the document editor can edit one (ADR-0023), and
// turn text into segments a component renders with each reference as a control that opens the
// editor at its line.
import { openEditor } from './editor.svelte';

/** The manifest's layout folders, as system-flow.yaml names them. */
export type Layout = { design: string; docs: string; wip: string };

/** The layout every project starts with; a caller that has the manifest's passes that instead. */
export const DEFAULT_LAYOUT: Layout = { design: 'design', docs: 'docs', wip: 'wip' };

/** A repository markdown file named in text, with the 1-based line and column it names, if any. */
export type FileRef = { path: string; line?: number; column?: number };

/** A run of plain text, or the text of a reference and the file it names. */
export type Segment = { text: string; ref?: FileRef };

/** What a control that opens a reference in the editor needs: its accessible label and its action. */
export type EditLink = FileRef & { label: string; open: () => void };

/** A run of plain text, or the text of a reference the editor can open and the control for it. */
export type LinkSegment = { text: string; edit?: EditLink };

/** Where an opened file saves, the story whose branch it is, and the folders references are under. */
export type LinkOptions = { story?: string; layout?: Layout };

const folder = (dir: string) => dir.replace(/^\.?\/+|\/+$/g, '');
const escape = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

// A path is not part of a longer path or word on either side; a segment never starts with a dot,
// so `..` cannot climb out of the folder; a line and a column are 1-based.
function pattern(layout: Layout): RegExp {
	const roots = [layout.design, layout.docs, layout.wip].map((d) => escape(folder(d))).join('|');
	return new RegExp(
		`(?<![\\w./-])((?:${roots})(?:/[\\w-][\\w.-]*)*?/[\\w-][\\w.-]*\\.md)` +
			`(?::([1-9]\\d*)(?::([1-9]\\d*))?)?(?![\\w/-]|\\.\\w)`,
		'g'
	);
}

function refOf(m: RegExpMatchArray): FileRef {
	const ref: FileRef = { path: m[1] };
	if (m[2]) ref.line = Number(m[2]);
	if (m[3]) ref.column = Number(m[3]);
	return ref;
}

/** Split text into plain runs and the repository markdown files it names under the layout's folders. */
export function parseFileRefs(text: string, layout: Layout = DEFAULT_LAYOUT): Segment[] {
	const out: Segment[] = [];
	let at = 0;
	for (const m of text.matchAll(pattern(layout))) {
		if (m.index > at) out.push({ text: text.slice(at, m.index) });
		out.push({ text: m[0], ref: refOf(m) });
		at = m.index + m[0].length;
	}
	if (at < text.length) out.push({ text: text.slice(at) });
	return out;
}

/** The file a whole string names, such as a code span's `design/x.md:12`, or null when it names none. */
export function fileRef(text: string, layout: Layout = DEFAULT_LAYOUT): FileRef | null {
	const [only, ...rest] = parseFileRefs(text.trim(), layout);
	return only?.ref && rest.length === 0 ? only.ref : null;
}

/**
 * Whether the editor can edit a repository path: not the archive, a thread, an issue file, or the
 * generated narrative index, which flai's docedit.ModeOf gives mode none. An accepted ADR is mode
 * none too, but only its status says so; the editor shows it read-only, with the reason, when opened.
 */
export function editable(path: string, layout: Layout = DEFAULT_LAYOUT): boolean {
	const wip = folder(layout.wip);
	const under = (dir: string) => path.startsWith(`${dir}/`);
	return !(
		under(`${wip}/archive`) ||
		under(`${wip}/threads`) ||
		under(`${folder(layout.design)}/issues`) ||
		path === `${wip}/agents/index.md`
	);
}

/** The accessible name of a control that opens a reference: `Edit <path>`, and `at line N` when it names one. */
export function editLabel(ref: FileRef): string {
	return `Edit ${ref.path}${ref.line ? ` at line ${ref.line}` : ''}`;
}

/** The control for a reference, opening the editor at its line on the story's branch, or null when the editor cannot edit it. */
export function editLink(ref: FileRef, opts: LinkOptions = {}): EditLink | null {
	if (!editable(ref.path, opts.layout)) return null;
	const { path, line } = ref;
	// wip lives on the main checkout, never on a story's branch (ADR-0019): flai refuses it with one.
	const wip = folder((opts.layout ?? DEFAULT_LAYOUT).wip) + '/';
	const story = path.startsWith(wip) ? undefined : opts.story;
	return {
		...ref,
		label: editLabel(ref),
		open: () => openEditor({ path, line, story })
	};
}

// How flai's refusal of an acceptance names lint left in the main checkout's uncommitted files
// (ADR-0140, flai/cmd/accept_markdown.go): those files are fixed there, not on the story's branch.
const MAIN_CHECKOUT_REFUSAL = "the main checkout's uncommitted files";

/**
 * Split plain text, such as a refusal, into runs and the references in it the editor can open; a
 * refusal for the main checkout's uncommitted files opens them there, whatever story is given.
 */
export function linkSegments(text: string, opts: LinkOptions = {}): LinkSegment[] {
	if (opts.story && text.includes(MAIN_CHECKOUT_REFUSAL)) opts = { ...opts, story: undefined };
	return parseFileRefs(text, opts.layout).map((s) => {
		const edit = s.ref && editLink(s.ref, opts);
		return edit ? { text: s.text, edit } : { text: s.text };
	});
}
