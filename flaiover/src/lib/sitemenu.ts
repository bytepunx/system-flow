// The nav rail's groups and pages (S-0458, replacing S-0172's two-tier header menu): six groups,
// each with its icon and its pages in order, and which group and page a path belongs to. The Charts
// page carries its charts by group from the catalogue, and a path names the chart it shows (S-0484).
// The Agents group lists a page for each role agent.status reports, at /workflow/<name> (S-0420).
// The Ops group holds the pages that run the project's kluster, Klusters first (S-0580).
import { resolve } from '$app/paths';
import type { LucideIcon } from '@lucide/svelte';
// Each icon by its own module, so dev and tests load the icons used and not the whole library.
import ActivityIcon from '@lucide/svelte/icons/activity';
import BookOpenIcon from '@lucide/svelte/icons/book-open';
import BotIcon from '@lucide/svelte/icons/bot';
import ChartLineIcon from '@lucide/svelte/icons/chart-line';
import CircleAlertIcon from '@lucide/svelte/icons/circle-alert';
import ClipboardListIcon from '@lucide/svelte/icons/clipboard-list';
import ComponentIcon from '@lucide/svelte/icons/component';
import ContainerIcon from '@lucide/svelte/icons/container';
import FileTextIcon from '@lucide/svelte/icons/file-text';
import FolderKanbanIcon from '@lucide/svelte/icons/folder-kanban';
import FolderTreeIcon from '@lucide/svelte/icons/folder-tree';
import InboxIcon from '@lucide/svelte/icons/inbox';
import MailIcon from '@lucide/svelte/icons/mail';
import MessagesSquareIcon from '@lucide/svelte/icons/messages-square';
import MicroscopeIcon from '@lucide/svelte/icons/microscope';
import NetworkIcon from '@lucide/svelte/icons/network';
import PaletteIcon from '@lucide/svelte/icons/palette';
import PlugIcon from '@lucide/svelte/icons/plug';
import RocketIcon from '@lucide/svelte/icons/rocket';
import ScaleIcon from '@lucide/svelte/icons/scale';
import ScrollTextIcon from '@lucide/svelte/icons/scroll-text';
import SearchIcon from '@lucide/svelte/icons/search';
import ServerIcon from '@lucide/svelte/icons/server';
import ServerCogIcon from '@lucide/svelte/icons/server-cog';
import SettingsIcon from '@lucide/svelte/icons/settings';
import ShipWheelIcon from '@lucide/svelte/icons/ship-wheel';
import SquareKanbanIcon from '@lucide/svelte/icons/square-kanban';
import SwatchBookIcon from '@lucide/svelte/icons/swatch-book';
import ToggleRightIcon from '@lucide/svelte/icons/toggle-right';
import UserCogIcon from '@lucide/svelte/icons/user-cog';
import WavesIcon from '@lucide/svelte/icons/waves';
import WorkflowIcon from '@lucide/svelte/icons/workflow';
// The catalogue, not viz/charts, so the layout's bundle does not take in the charts' builders.
import { CHART_GROUPS, TITLES, type ChartGroupKey, type Kind } from '$lib/viz/catalog';
import { ROLE_NAME } from '$lib/strategic';

/** A page the rail links to: the paths it answers, and the href that opens it. */
export type LinkedPage = {
	key: PageKey;
	label: string;
	icon: LucideIcon;
	/** The paths the page answers: it and every path below it. */
	path: string;
	/** Where the rail links, resolve()d. */
	href: string;
	/** An indicator the page carries, shown on its group too. */
	badge?: 'inbox';
	/** What the page lists below it in the rail, by group. */
	children?: PageChildGroup[];
	soon?: false;
};
/** A labelled group of the things a page lists below it. */
export type PageChildGroup = { key: ChartGroupKey; label: string; items: PageChild[] };
/** One thing a page lists below it, a chart, and the href that opens it, resolve()d. */
export type PageChild = { key: Kind; label: string; href: string };
/** A page still to come: shown in the rail, but no link and no path. */
export type SoonPage = {
	key: PageKey;
	label: string;
	icon: LucideIcon;
	soon: true;
	path?: never;
	href?: never;
	badge?: never;
	children?: never;
};
export type MenuPage = LinkedPage | SoonPage;
export type MenuGroup = {
	key: GroupKey;
	label: string;
	icon: LucideIcon;
	pages: MenuPage[];
	/** Paths outside its pages that still belong to the group, such as an item's page. */
	also?: string[];
};
export type GroupKey = 'flow' | 'agents' | 'design' | 'ops' | 'documentation' | 'configuration';
export type PageKey =
	| 'board'
	| 'inbox'
	| 'charts'
	| 'workflow'
	| 'orchestration'
	| 'planning'
	| 'analysis'
	| 'activity'
	| 'threads'
	| 'messages'
	| 'design-system'
	| 'components'
	| 'klusters'
	| 'deployments'
	| 'images'
	| 'docs-system'
	| 'conventions'
	| 'adrs'
	| 'issues'
	| 'search'
	| 'features'
	| 'agents'
	| 'file-system'
	| 'projects'
	| 'harnesses'
	| 'host'
	/** The page of a role other than the orchestrator, the planner, and the analyzer (S-0420). */
	| `role:${string}`;

/**
 * The orchestrator's, the planner's, and the analyzer's pages, by the role each shows, in the rail's
 * order: their state, activity, and runs (S-0228, S-0259), and what the operator can have flai do to
 * each.
 */
const BUILT_IN_ROLES: [string, LinkedPage][] = [
	[
		'orchestrator',
		{
			key: 'orchestration',
			label: 'Orchestration',
			icon: NetworkIcon,
			path: '/workflow/orchestrator',
			href: resolve('/workflow/orchestrator')
		}
	],
	[
		'planner',
		{
			key: 'planning',
			label: 'Planning',
			icon: ClipboardListIcon,
			path: '/workflow/planner',
			href: resolve('/workflow/planner')
		}
	],
	[
		'analyzer',
		{
			key: 'analysis',
			label: 'Analysis',
			icon: MicroscopeIcon,
			path: '/workflow/analyzer',
			href: resolve('/workflow/analyzer')
		}
	]
];

/** The Agents group's pages after its roles'. */
const AGENT_PAGES: LinkedPage[] = [
	{
		key: 'activity',
		label: 'Activity',
		icon: ActivityIcon,
		path: '/activity',
		href: resolve('/activity')
	},
	// Every open thread, where a thread on no item is answered (S-0173, TH-0041).
	{
		key: 'threads',
		label: 'Threads',
		icon: MessagesSquareIcon,
		path: '/threads',
		href: resolve('/threads')
	},
	// The conversations between stories' agents, which never count in the inbox (S-0336).
	{
		key: 'messages',
		label: 'Messages',
		icon: MailIcon,
		path: '/messages',
		href: resolve('/messages')
	}
];

export const SITE_MENU: MenuGroup[] = [
	{
		key: 'flow',
		label: 'Flow',
		icon: WavesIcon,
		pages: [
			{
				key: 'board',
				label: 'Board',
				icon: SquareKanbanIcon,
				path: '/board',
				href: resolve('/board')
			},
			{
				key: 'inbox',
				label: 'Inbox',
				icon: InboxIcon,
				path: '/inbox',
				href: resolve('/inbox'),
				badge: 'inbox'
			},
			{
				key: 'charts',
				label: 'Charts',
				icon: ChartLineIcon,
				path: '/charts',
				href: resolve('/charts/[kind]', { kind: 'cycle-time' }),
				children: CHART_GROUPS.map((g) => ({
					key: g.key,
					label: g.label,
					items: g.kinds.map((kind) => ({
						key: kind,
						label: TITLES[kind],
						href: resolve('/charts/[kind]', { kind })
					}))
				}))
			},
			// The workflow's own page is still to come (S-0458).
			{ key: 'workflow', label: 'Workflow', icon: WorkflowIcon, soon: true }
		],
		also: ['/items', '/review', '/new']
	},
	{
		key: 'agents',
		label: 'Agents',
		icon: BotIcon,
		// Until agent.status names the roles, the three that every project runs (siteMenu).
		pages: [...BUILT_IN_ROLES.map(([, p]) => p), ...AGENT_PAGES],
		// A role's page the rail does not list, such as one flai has not reported yet.
		also: ['/workflow']
	},
	{
		key: 'design',
		label: 'Design',
		icon: PaletteIcon,
		pages: [
			// The project's design system: its tokens, its guidance, and a gallery they style (S-0382).
			{
				key: 'design-system',
				label: 'System',
				icon: SwatchBookIcon,
				path: '/design',
				href: resolve('/design')
			},
			// Every component of the project, rendered live from its own dev server (S-0439).
			{
				key: 'components',
				label: 'Components',
				icon: ComponentIcon,
				path: '/design/components',
				href: resolve('/design/components')
			}
		]
	},
	{
		key: 'ops',
		label: 'Ops',
		icon: ServerCogIcon,
		pages: [
			// The project's kluster: whether it runs, and Up and Down through flai (S-0580).
			{
				key: 'klusters',
				label: 'Klusters',
				icon: ShipWheelIcon,
				path: '/ops/klusters',
				href: resolve('/ops/klusters')
			},
			// Main's namespace and each story's on the kluster, with Deploy,
			// Teardown, and Promote through flai (S-0581).
			{
				key: 'deployments',
				label: 'Deployments',
				icon: RocketIcon,
				path: '/ops/deployments',
				href: resolve('/ops/deployments')
			},
			// The images the project's agents run in and flai serve's jobs that make them, with
			// Retry for a failed one through flai (S-0757).
			{
				key: 'images',
				label: 'Agent images',
				icon: ContainerIcon,
				path: '/ops/images',
				href: resolve('/ops/images')
			}
		],
		// The group's own path, which has no page of its own, and the pages below it the rail does
		// not list, such as the log page, /ops/log, which a log path in a message links (S-0785).
		also: ['/ops']
	},
	{
		key: 'documentation',
		label: 'Documentation',
		icon: BookOpenIcon,
		// A folder's page answers the whole folder and opens at its README or summary: the document
		// viewer opens files, not folders (S-0458).
		pages: [
			{
				key: 'docs-system',
				label: 'System',
				icon: FileTextIcon,
				path: '/docs/design/system',
				href: resolve('/docs/[...path]', { path: 'design/system/README.md' })
			},
			{
				key: 'conventions',
				label: 'Conventions',
				icon: ScrollTextIcon,
				path: '/docs/design/conventions',
				href: resolve('/docs/[...path]', { path: 'design/conventions/README.md' })
			},
			{ key: 'adrs', label: 'ADRs', icon: ScaleIcon, path: '/adrs', href: resolve('/adrs') },
			{
				key: 'issues',
				label: 'Issues',
				icon: CircleAlertIcon,
				path: '/docs/design/issues',
				href: resolve('/docs/[...path]', { path: 'design/issues/summary.md' })
			},
			{
				key: 'search',
				label: 'Search',
				icon: SearchIcon,
				path: '/search',
				href: resolve('/search')
			}
		],
		// The document viewer's other folders, and the editor.
		also: ['/docs', '/edit']
	},
	{
		key: 'configuration',
		label: 'Configuration',
		icon: SettingsIcon,
		// The Configuration pages (S-0459, TH-0553), each a page of the settings below /settings; the
		// settingspages slugs are their keys. /settings itself, which goes on to one of them, and an
		// unknown page below it are in the group with no page. The Host page (S-0489) is the one page
		// of the group that is not a settings page: it stays at /host, which the host-flai pill links.
		pages: [
			{
				key: 'features',
				label: 'Features',
				icon: ToggleRightIcon,
				path: '/settings/features',
				href: resolve('/settings/[page]', { page: 'features' })
			},
			{
				key: 'agents',
				label: 'Agents',
				icon: UserCogIcon,
				path: '/settings/agents',
				href: resolve('/settings/[page]', { page: 'agents' })
			},
			{
				key: 'file-system',
				label: 'File System',
				icon: FolderTreeIcon,
				path: '/settings/file-system',
				href: resolve('/settings/[page]', { page: 'file-system' })
			},
			{
				key: 'projects',
				label: 'Projects',
				icon: FolderKanbanIcon,
				path: '/settings/projects',
				href: resolve('/settings/[page]', { page: 'projects' })
			},
			{
				key: 'harnesses',
				label: 'Harnesses',
				icon: PlugIcon,
				path: '/settings/harnesses',
				href: resolve('/settings/[page]', { page: 'harnesses' })
			},
			{ key: 'host', label: 'Host', icon: ServerIcon, path: '/host', href: resolve('/host') }
		],
		also: ['/settings']
	}
];

/** The page of a role besides the three, labelled by its name, capitalised, its hyphens as spaces. */
function rolePage(name: string): LinkedPage {
	return {
		key: `role:${name}`,
		label: `${name[0].toUpperCase()}${name.slice(1).replaceAll('-', ' ')}`,
		icon: BotIcon,
		path: `/workflow/${name}`,
		href: resolve('/workflow/[role]', { role: name })
	};
}

/**
 * The rail's groups with a page in Agents for each role agent.status reports (S-0420): the
 * orchestrator's, the planner's, and the analyzer's first, as the rail has always listed them, then
 * the rest, the curator and the roles the project defines, by name, each at /workflow/<name>. A name
 * that is not a role's is left out. With no roles known, from a flai older than the role page or
 * before one answers, it is SITE_MENU.
 */
export function siteMenu(roles?: readonly string[] | null): MenuGroup[] {
	if (!roles?.length) return SITE_MENU;
	const reported = new Set(roles);
	const builtIn = new Set(BUILT_IN_ROLES.map(([name]) => name));
	const others = [...reported]
		.filter((name) => !builtIn.has(name) && ROLE_NAME.test(name))
		.sort((a, b) => a.localeCompare(b));
	const pages = [
		...BUILT_IN_ROLES.filter(([name]) => reported.has(name)).map(([, p]) => p),
		...others.map(rolePage),
		...AGENT_PAGES
	];
	return SITE_MENU.map((g) => (g.key === 'agents' ? { ...g, pages } : g));
}

function under(pathname: string, path: string): boolean {
	return pathname === path || pathname.startsWith(`${path}/`);
}

function childAt(page: LinkedPage, pathname: string): Kind | null {
	const items = page.children?.flatMap((g) => g.items) ?? [];
	return items.find((c) => under(pathname, c.href))?.key ?? null;
}

/** Where a path is in the rail: its group, its page, and the page's child it shows. */
export type Located = { group: GroupKey | null; page: PageKey | null; child: Kind | null };

/**
 * The group, page, and child a path belongs to; each is null where the rail has none for it, as for
 * Overview and License, which the top bar and the rail's foot link to. Within a group the page with
 * the longest path that answers wins, so Components at /design/components wins below the design
 * system at /design. The child is the page's child whose href the path is or is under, so
 * /charts/throughput names its chart and /charts/nope none (S-0484). The path is compared with the
 * hrefs as given, as the License link is: both carry the base path, which the page paths leave out.
 * The menu is SITE_MENU, or siteMenu's with the roles agent.status reports.
 */
export function locate(pathname: string, menu: MenuGroup[] = SITE_MENU): Located {
	for (const g of menu) {
		const page = g.pages
			.filter((p): p is LinkedPage => !p.soon)
			.filter((p) => under(pathname, p.path))
			.reduce<LinkedPage | undefined>(
				(a, p) => (a && a.path.length >= p.path.length ? a : p),
				undefined
			);
		if (page) return { group: g.key, page: page.key, child: childAt(page, pathname) };
		if (g.also?.some((a) => under(pathname, a))) return { group: g.key, page: null, child: null };
	}
	return { group: null, page: null, child: null };
}

/** A page's href with the project it was followed from, for a copied link. */
export function withProject(href: string, project: string | null | undefined): string {
	return project ? `${href}?project=${encodeURIComponent(project)}` : href;
}

/** The indicators a group shows: its pages', each once, in page order. */
export function groupBadges(group: MenuGroup): NonNullable<MenuPage['badge']>[] {
	return [...new Set(group.pages.flatMap((p) => (p.badge ? [p.badge] : [])))];
}
