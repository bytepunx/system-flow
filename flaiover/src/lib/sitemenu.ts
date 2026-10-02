// The site menu's two tiers (S-0172): three groups in the header, each with its pages in a row
// below, and which group and page a path belongs to.

/** A page of the menu, named by key so the header can resolve its route. */
export type MenuPage = {
	key: PageKey;
	label: string;
	/** The paths the page answers: this one exactly for '/', else it and every path below it. */
	path: string;
	/** An indicator the page carries, shown on its group too. */
	badge?: 'inbox';
};
export type MenuGroup = {
	key: GroupKey;
	label: string;
	pages: MenuPage[];
	/** Paths outside its pages that still belong to the group, such as an item's page. */
	also?: string[];
};
export type GroupKey = 'workflow' | 'status' | 'host';
export type PageKey =
	| 'overview'
	| 'board'
	| 'inbox'
	| 'threads'
	| 'activity'
	| 'charts'
	| 'adrs'
	| 'docs'
	| 'search'
	| 'updates'
	| 'settings'
	| 'license';

export const SITE_MENU: MenuGroup[] = [
	{
		key: 'workflow',
		label: 'Workflow',
		pages: [
			{ key: 'overview', label: 'Overview', path: '/' },
			{ key: 'board', label: 'Board', path: '/board' },
			{ key: 'inbox', label: 'Inbox', path: '/inbox', badge: 'inbox' },
			// Every open thread, where a thread on no item is answered (S-0173, TH-0041).
			{ key: 'threads', label: 'Threads', path: '/threads' },
			{ key: 'activity', label: 'Activity', path: '/activity' }
		],
		also: ['/items', '/review', '/new']
	},
	{
		key: 'status',
		label: 'Status',
		pages: [
			{ key: 'charts', label: 'Charts', path: '/charts' },
			{ key: 'adrs', label: 'ADRs', path: '/adrs' },
			{ key: 'docs', label: 'Docs', path: '/docs' },
			{ key: 'search', label: 'Search', path: '/search' }
		],
		also: ['/edit']
	},
	{
		key: 'host',
		label: 'Host',
		pages: [
			// The host page is where the dashboard and flai are upgraded (S-0081, S-0107).
			{ key: 'updates', label: 'Updates', path: '/host' },
			{ key: 'settings', label: 'Settings', path: '/settings' },
			// The license the image carries, the same text flai license prints (S-0231).
			{ key: 'license', label: 'License', path: '/license' }
		]
	}
];

function under(pathname: string, path: string): boolean {
	if (path === '/') return pathname === '/';
	return pathname === path || pathname.startsWith(`${path}/`);
}

/** The group and page a path belongs to; either is null where the menu has none for it. */
export function locate(
	pathname: string,
	menu: MenuGroup[] = SITE_MENU
): { group: GroupKey | null; page: PageKey | null } {
	for (const g of menu) {
		const page = g.pages.find((p) => under(pathname, p.path));
		if (page) return { group: g.key, page: page.key };
		if (g.also?.some((a) => under(pathname, a))) return { group: g.key, page: null };
	}
	return { group: null, page: null };
}

/** The indicators a group shows: its pages', each once, in page order. */
export function groupBadges(group: MenuGroup): NonNullable<MenuPage['badge']>[] {
	return [...new Set(group.pages.flatMap((p) => (p.badge ? [p.badge] : [])))];
}
