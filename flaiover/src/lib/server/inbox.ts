// What needs a human (S-0042): threads awaiting the designer, open questions in narratives, stories
// in review, blocked items, and overlapping touches. flai on the host composes it from the files
// (S-0074, inbox.designer; not the agent's MCP inbox); the dashboard adds its own link to each entry.
import type { Repo } from './repo';

export type InboxKind = 'thread' | 'question' | 'review' | 'blocked' | 'overlap';
export type InboxEntry = {
	key: string; // stable: the same thing keeps the same key, so "new" means something
	kind: InboxKind;
	title: string;
	detail?: string;
	href: string; // a dashboard path
	at?: string;
	// The story or item this entry is about, when it has one: a question
	// entry needs it to answer in place (S-0090), not just to link out.
	item?: string;
	// On a thread, the recommendation awaiting the designer's confirmation, so it can be confirmed
	// in place (ADR-0090): its text without the source line, and the source apart.
	recommendation?: Recommendation;
};
export type Recommendation = {
	author: string;
	at: string;
	text: string;
	source?: { path: string; heading?: string };
};
export type Inbox = {
	total: number;
	counts: Record<InboxKind, number>;
	entries: InboxEntry[];
	notes: string[];
};

type FlaiEntry = Omit<InboxEntry, 'href' | 'recommendation'> & {
	path?: string;
	pending_recommendation?: Omit<Recommendation, 'source'> & {
		source?: Recommendation['source'] | null;
	};
};

/** An entry's text without the `Source:` paragraph flai ends it with when it cites one. */
export function withoutSource(text: string): string {
	return text.replace(/\n\nSource: [^\n]*$/, '').trim();
}

// A work item's file or a narrative, live or archived, names its item: wip/kanban/stories/S-0001-a.md, wip/agents/S-0001.md.
const itemPath =
	/^wip\/(?:archive\/)?(?:kanban\/(?:epics|stories|tasks)\/([A-Z]+-\d+)-[^/]*|agents\/([A-Z]+-\d+))\.md$/;

/** The item a repository path is the file or narrative of, if it is one. */
export function itemOf(path: string | undefined): string | undefined {
	const m = path ? itemPath.exec(path) : null;
	return m ? (m[1] ?? m[2]) : undefined;
}

/**
 * Where an entry leads, never to a document, where it could not be answered (S-0173): the review
 * page for a story in review; a question to its story's page, opened on it; a thread on an item, or
 * on an item's file or narrative, to that item's page opened on the thread (S-0155); any other
 * thread to the threads page opened on it (TH-0041); anything else to its item, else the board.
 */
export function hrefFor(e: {
	kind: InboxKind;
	key?: string;
	item?: string;
	path?: string;
}): string {
	const item = e.item || itemOf(e.path);
	if (e.kind === 'review' && item) return `/review/${item}`;
	if (e.kind === 'question' && item && e.key)
		return `/items/${item}?question=${encodeURIComponent(e.key)}`;
	const thread = e.kind === 'thread' ? e.key?.replace(/^thread:/, '') : undefined;
	if (thread)
		return item
			? `/items/${item}?thread=${encodeURIComponent(thread)}`
			: `/threads?thread=${encodeURIComponent(thread)}`;
	if (item) return `/items/${item}`;
	return '/board';
}

export async function inbox(repo: Repo): Promise<Inbox> {
	const got = await repo.remember<Omit<Inbox, 'entries'> & { entries: FlaiEntry[] | null }>(
		'inbox',
		'inbox.designer'
	);
	return {
		total: got.total,
		counts: got.counts,
		notes: got.notes ?? [],
		entries: (got.entries ?? []).map(({ path, pending_recommendation: p, ...e }) => ({
			...e,
			detail: e.detail || undefined,
			at: e.at || undefined,
			item: e.item || undefined,
			href: hrefFor({ kind: e.kind, key: e.key, item: e.item, path }),
			...(p && {
				recommendation: {
					author: p.author,
					at: p.at,
					text: withoutSource(p.text),
					...(p.source?.path && { source: p.source })
				}
			})
		}))
	};
}
