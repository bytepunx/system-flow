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
};
export type Inbox = {
	total: number;
	counts: Record<InboxKind, number>;
	entries: InboxEntry[];
	notes: string[];
};

type FlaiEntry = Omit<InboxEntry, 'href'> & { item?: string; path?: string };

/** Where an entry leads: the review page for a story in review, else its item, else its document. */
export function hrefFor(e: { kind: InboxKind; item?: string; path?: string }): string {
	if (e.kind === 'review' && e.item) return `/review/${e.item}`;
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
		entries: (got.entries ?? []).map(({ item, path, ...e }) => ({
			...e,
			detail: e.detail || undefined,
			at: e.at || undefined,
			href: hrefFor({ kind: e.kind, item, path })
		}))
	};
}
