// What needs a human (S-0042): threads awaiting the designer, open questions
// in narratives, stories in review, blocked items, and overlapping touches.
// All of it is read from the files; nothing is recorded to build it.
import { readdir, readFile } from 'node:fs/promises';
import { join } from 'node:path';
import { splitFrontMatter, type Repo } from './repo';
import { flai } from './flai';

export type InboxKind = 'thread' | 'question' | 'review' | 'blocked' | 'overlap';
export type InboxEntry = {
	key: string; // stable: the same thing keeps the same key, so "new" means something
	kind: InboxKind;
	title: string;
	detail?: string;
	href: string; // a dashboard path
	at?: string;
};
export type Inbox = {
	total: number;
	counts: Record<InboxKind, number>;
	entries: InboxEntry[];
	notes: string[];
};

/** The bullets under `## Open questions`, outside the block flai generates for threads. */
export function openQuestions(body: string): string[] {
	const lines = body.split('\n');
	const start = lines.findIndex((l) => l.trim() === '## Open questions');
	if (start < 0) return [];
	const out: string[] = [];
	let generated = false;
	for (const line of lines.slice(start + 1)) {
		if (/^##\s/.test(line)) break;
		if (line.includes('<!-- threads:start -->')) generated = true;
		else if (line.includes('<!-- threads:end -->')) generated = false;
		else if (!generated) {
			const m = /^\s*[-*]\s+(.*\S)\s*$/.exec(line);
			if (m) out.push(m[1]);
			else if (out.length && /^\s{2,}\S/.test(line)) out[out.length - 1] += ' ' + line.trim();
		}
	}
	return out;
}

const hash = (s: string): string => {
	let h = 5381;
	for (let i = 0; i < s.length; i++) h = ((h << 5) + h + s.charCodeAt(i)) >>> 0;
	return h.toString(36);
};
const stamp = (v: unknown): string =>
	v instanceof Date ? v.toISOString().replace(/\.\d{3}Z$/, 'Z') : String(v ?? '');

type Finding = { rule: string; path: string; message: string };
let overlapCache: { repo: Repo; findings: Finding[] | null } | null = null;

/** `wip.overlap` findings from flai check, so the rule has one implementation; cached until the repository changes. */
async function overlaps(repo: Repo): Promise<Finding[] | null> {
	if (overlapCache?.repo === repo) return overlapCache.findings;
	let findings: Finding[] | null;
	try {
		const { data } = await flai<{ findings?: Finding[] }>(repo.root, ['check'], {
			// flai check exits 1 when it has errors; its findings are still the answer
			exitStatus: { 1: 200 }
		}).catch((e) => {
			if (e && typeof e === 'object' && 'data' in e && (e as { data?: unknown }).data)
				return { data: (e as { data: { findings?: Finding[] } }).data };
			throw e;
		});
		findings = (data?.findings ?? []).filter((f) => f.rule === 'wip.overlap');
	} catch {
		findings = null; // flai is not available, or failed: say so rather than guess
	}
	overlapCache = { repo, findings };
	repo.once('change', () => {
		if (overlapCache?.repo === repo) overlapCache = null;
	});
	return findings;
}

/** Forget cached overlaps; for tests. */
export function resetInboxCache(): void {
	overlapCache = null;
}

export async function inbox(repo: Repo): Promise<Inbox> {
	const entries: InboxEntry[] = [];
	const notes: string[] = [];
	// the designer is the manifest's owner, as for threads and board moves
	const who = String(((await repo.manifest()) as { owner?: unknown }).owner || 'designer');

	for (const th of await repo.threads()) {
		if (th.status === 'resolved') continue;
		const last = th.entries.at(-1);
		if (!last || last.author === who) continue;
		const on = th.anchor.item ?? th.anchor.path;
		entries.push({
			key: `thread:${th.id}`,
			kind: 'thread',
			title: th.title,
			detail: `${last.author} wrote last, on ${on}${th.anchor.heading ? ` (${th.anchor.heading})` : ''}`,
			href: th.anchor.item ? `/items/${th.anchor.item}` : `/docs/${th.anchor.path}`,
			at: last.at
		});
	}

	const layout = await repo.layout();
	const dir = `${layout.wip}/agents`;
	const abs = repo.resolveInside(dir);
	for (const name of (await readdir(abs).catch(() => [] as string[])).sort()) {
		if (!name.endsWith('.md') || name === 'index.md' || name === 'README.md') continue;
		const { frontMatter, body } = splitFrontMatter(await readFile(join(abs, name), 'utf8'));
		const fm = (frontMatter ?? {}) as Record<string, unknown>;
		const stream = String(fm.stream ?? name.replace(/\.md$/, ''));
		for (const q of openQuestions(body))
			entries.push({
				key: `question:${stream}:${hash(q)}`,
				kind: 'question',
				title: q,
				detail: `asked in the narrative of ${stream}`,
				href: `/docs/${dir}/${name}`,
				at: stamp(fm.updated) || undefined
			});
	}

	const items = (await repo.items()).filter((it) => !it.archived);
	for (const it of items) {
		if (it.type === 'story' && it.status === 'review')
			entries.push({
				key: `review:${it.id}`,
				kind: 'review',
				title: `${it.id} ${it.title}`,
				detail: 'in review: accept it or send it back',
				href: `/review/${it.id}`,
				at: it.transitions.at(-1)?.at
			});
		const block = (it.blocked ?? []).find((b) => !b.until);
		if (block && it.status !== 'done' && it.status !== 'cancelled')
			entries.push({
				key: `blocked:${it.id}:${block.from}`,
				kind: 'blocked',
				title: `${it.id} ${it.title}`,
				detail: `blocked: ${block.reason}`,
				href: `/items/${it.id}`,
				at: block.from
			});
	}

	const found = await overlaps(repo);
	if (found === null)
		notes.push('Overlapping touches are not listed: flai is not available to this dashboard.');
	else
		for (const f of found) {
			const id = /^([EST]-\d+)/.exec(f.message)?.[1];
			entries.push({
				key: `overlap:${hash(f.message)}`,
				kind: 'overlap',
				title: f.message,
				href: id ? `/items/${id}` : '/board'
			});
		}

	const counts: Record<InboxKind, number> = {
		thread: 0,
		question: 0,
		review: 0,
		blocked: 0,
		overlap: 0
	};
	for (const e of entries) counts[e.kind]++;
	return { total: entries.length, counts, entries, notes };
}
