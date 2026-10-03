// Pure helpers for the issues the review page offers as stories (S-0198).

/** An issue as flai's issue.list gives it. */
export type Issue = {
	id: string;
	title: string;
	class: string;
	status: string;
	count: number;
	cost?: string;
	path: string;
	/** The stories its occurrences name. */
	stories: string[];
	/** The open story that links it, or "" when none does. */
	story: string;
};

/** One issue offered on a story's review page. */
export type IssueRow = {
	id: string;
	title: string;
	class: string;
	count: number;
	/** Its document, for a link: design/issues and its file name. */
	doc: string;
	/** The story under review recorded or bumped it. */
	recorded: boolean;
	/** Ticked when the page opens: a story is made for it on acceptance. */
	checked: boolean;
};

const num = (id: string) => Number(id.replace(/^I-/, '')) || 0;

/**
 * The issues a story's review page offers: first every open issue the story recorded or bumped,
 * checked, then every other open issue, unchecked. An issue an open story already links, or a closed
 * one, is left out. `own` is the story's own list (issue.list with its story), read from its worktree,
 * and its copy of an issue wins over `main`'s.
 */
export function issueRows(story: string, main: Issue[], own: Issue[]): IssueRow[] {
	const byId = new Map<string, { issue: Issue; recorded: boolean }>();
	for (const issue of main)
		byId.set(issue.id, { issue, recorded: (issue.stories ?? []).includes(story) });
	for (const issue of own) byId.set(issue.id, { issue, recorded: true });
	const rows = [...byId.values()]
		.filter(({ issue }) => issue.status === 'open' && !issue.story)
		.map(({ issue, recorded }) => ({
			id: issue.id,
			title: issue.title,
			class: issue.class,
			count: issue.count,
			doc: issueDocPath(issue.path),
			recorded,
			checked: recorded
		}));
	return rows.sort((a, b) => Number(b.recorded) - Number(a.recorded) || num(a.id) - num(b.id));
}

/** The accept button's words: it also makes stories when any issue is checked. */
export function acceptLabel(checked: number, running: boolean): string {
	if (running) return 'Accepting…';
	return checked > 0 ? 'Accept and Create Stories' : 'Accept';
}

/** The path of an issue's document, for a link to it: design/issues and its file name. */
export function issueDocPath(path: string): string {
	return `design/issues/${path.split('/').pop()}`;
}
