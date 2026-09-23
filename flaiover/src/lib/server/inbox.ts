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
};
export type Inbox = {
	total: number;
	counts: Record<InboxKind, number>;
	entries: InboxEntry[];
	notes: string[];
};

type FlaiEntry = Omit<InboxEntry, 'href'> & { path?: string };

/** Where an entry leads: the review page for a story in review, a question's own narrative document (even once it also carries an item id, to answer it in place), else its item, else its document. */
export function hrefFor(e: { kind: InboxKind; item?: string; path?: string }): string {
	if (e.kind === 'review' && e.item) return `/review/${e.item}`;
	if (e.kind === 'question' && e.path) return `/docs/${e.path}`;
	if (e.item) return `/items/${e.item}`;
	if (e.path) return `/docs/${e.path}`;
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
		entries: (got.entries ?? []).map(({ path, ...e }) => ({
			...e,
			detail: e.detail || undefined,
			at: e.at || undefined,
			item: e.item || undefined,
			href: hrefFor({ kind: e.kind, item: e.item, path })
		}))
	};
}
