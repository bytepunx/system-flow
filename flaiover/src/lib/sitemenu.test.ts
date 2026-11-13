// S-0458: the nav rail's groups and pages, and which group and page a path belongs to.
import { describe, expect, it, vi } from 'vitest';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

import ActivityIcon from '@lucide/svelte/icons/activity';
import ContainerIcon from '@lucide/svelte/icons/container';
import InboxIcon from '@lucide/svelte/icons/inbox';
import ServerCogIcon from '@lucide/svelte/icons/server-cog';
import RocketIcon from '@lucide/svelte/icons/rocket';
import ShipWheelIcon from '@lucide/svelte/icons/ship-wheel';
import { KINDS, TITLES } from './viz/catalog';
import { SITE_MENU, groupBadges, locate, siteMenu, withProject } from './sitemenu';

describe('the nav rail (S-0458)', () => {
	it('has Flow, Agents, Design, Ops, Documentation, and Configuration, each with its pages', () => {
		expect(SITE_MENU.map((g) => [g.label, g.pages.map((p) => p.label)])).toEqual([
			['Flow', ['Board', 'Inbox', 'Charts', 'Workflow']],
			['Agents', ['Orchestration', 'Planning', 'Analysis', 'Activity', 'Threads', 'Messages']],
			['Design', ['System', 'Components']],
			['Ops', ['Klusters', 'Deployments', 'Agent images']],
			['Documentation', ['System', 'Conventions', 'ADRs', 'Issues', 'Search']],
			['Configuration', ['Features', 'Agents', 'File System', 'Projects', 'Harnesses', 'Host']]
		]);
	});

	it('gives every group and page an icon, and every page a unique key', () => {
		for (const g of SITE_MENU) {
			expect(g.icon).toBeTypeOf('function');
			for (const p of g.pages) expect(p.icon, p.label).toBeTypeOf('function');
		}
		const keys = SITE_MENU.flatMap((g) => g.pages.map((p) => p.key));
		expect(new Set(keys).size).toBe(keys.length);
		const flow = SITE_MENU[0].pages;
		expect(flow[1].icon).toBe(InboxIcon);
		expect(SITE_MENU[1].pages[3].icon).toBe(ActivityIcon);
	});

	it('links every page to its route, and Workflow to none, as still to come', () => {
		expect(SITE_MENU.map((g) => g.pages.map((p) => p.href ?? null))).toEqual([
			['/board', '/inbox', '/charts/cycle-time', null],
			[
				'/workflow/orchestrator',
				'/workflow/planner',
				'/workflow/analyzer',
				'/activity',
				'/threads',
				'/messages'
			],
			['/design', '/design/components'],
			['/ops/klusters', '/ops/deployments', '/ops/images'],
			[
				'/docs/design/system/README.md',
				'/docs/design/conventions/README.md',
				'/adrs',
				'/docs/design/issues/summary.md',
				'/search'
			],
			[
				'/settings/features',
				'/settings/agents',
				'/settings/file-system',
				'/settings/projects',
				'/settings/harnesses',
				'/host'
			]
		]);
		const workflow = SITE_MENU[0].pages[3];
		expect(workflow.soon).toBe(true);
		expect(workflow.path).toBeUndefined();
		expect(SITE_MENU.flatMap((g) => g.pages).filter((p) => p.soon)).toEqual([workflow]);
	});

	it('places each page and the paths below it in its group', () => {
		expect(locate('/board')).toEqual({ group: 'flow', page: 'board', child: null });
		expect(locate('/inbox')).toEqual({ group: 'flow', page: 'inbox', child: null });
		expect(locate('/charts')).toEqual({ group: 'flow', page: 'charts', child: null });
		expect(locate('/workflow/orchestrator')).toEqual({
			group: 'agents',
			page: 'orchestration',
			child: null
		});
		expect(locate('/workflow/planner')).toEqual({ group: 'agents', page: 'planning', child: null });
		expect(locate('/workflow/analyzer')).toEqual({
			group: 'agents',
			page: 'analysis',
			child: null
		});
		expect(locate('/activity')).toEqual({ group: 'agents', page: 'activity', child: null });
		expect(locate('/threads')).toEqual({ group: 'agents', page: 'threads', child: null });
		expect(locate('/messages')).toEqual({ group: 'agents', page: 'messages', child: null });
		expect(locate('/design')).toEqual({ group: 'design', page: 'design-system', child: null });
		expect(locate('/design/components')).toEqual({
			group: 'design',
			page: 'components',
			child: null
		});
		expect(locate('/docs/design/system')).toEqual({
			group: 'documentation',
			page: 'docs-system',
			child: null
		});
		expect(locate('/docs/design/system/overview.md')).toEqual({
			group: 'documentation',
			page: 'docs-system',
			child: null
		});
		expect(locate('/docs/design/conventions/git.md')).toEqual({
			group: 'documentation',
			page: 'conventions',
			child: null
		});
		expect(locate('/adrs')).toEqual({ group: 'documentation', page: 'adrs', child: null });
		expect(locate('/docs/design/issues/summary.md')).toEqual({
			group: 'documentation',
			page: 'issues',
			child: null
		});
		expect(locate('/docs/design/issues/I-0001-golangci-lint-version-mismatch.md')).toEqual({
			group: 'documentation',
			page: 'issues',
			child: null
		});
		expect(locate('/search')).toEqual({ group: 'documentation', page: 'search', child: null });
	});

	it('lists the charts below Charts by group, every kind once (S-0484)', () => {
		const charts = SITE_MENU[0].pages[2];
		expect(charts.key).toBe('charts');
		expect(charts.path).toBe('/charts');
		const groups = charts.children ?? [];
		expect(groups.map((g) => [g.key, g.label])).toEqual([
			['flow', 'flow'],
			['usage', 'usage'],
			['planning', 'planning'],
			['strategic', 'strategic']
		]);
		const items = groups.flatMap((g) => g.items);
		expect(items.map((c) => c.key).sort()).toEqual([...KINDS].sort());
		for (const c of items) {
			expect(c.label, c.key).toBe(TITLES[c.key]);
			expect(c.href, c.key).toBe(`/charts/${c.key}`);
		}
		const others = SITE_MENU.flatMap((g) => g.pages).filter((p) => p !== charts);
		expect(others.filter((p) => p.children)).toEqual([]);
	});

	it('names the chart a /charts path shows, and none for an unknown kind (S-0484)', () => {
		expect(locate('/charts/cycle-time')).toEqual({
			group: 'flow',
			page: 'charts',
			child: 'cycle-time'
		});
		expect(locate('/charts/throughput')).toEqual({
			group: 'flow',
			page: 'charts',
			child: 'throughput'
		});
		expect(locate('/charts/nope')).toEqual({ group: 'flow', page: 'charts', child: null });
		expect(locate('/charts/cfdx')).toEqual({ group: 'flow', page: 'charts', child: null });
		expect(locate('/charts')).toEqual({ group: 'flow', page: 'charts', child: null });
		for (const kind of KINDS) expect(locate(`/charts/${kind}`).child).toBe(kind);
		expect(locate('/board').child).toBeNull();
		expect(locate('/design/components').child).toBeNull();
		expect(locate('/license').child).toBeNull();
	});

	it('selects the longest page path that answers, Components below the design system', () => {
		expect(locate('/design/components/Button')).toEqual({
			group: 'design',
			page: 'components',
			child: null
		});
		expect(locate('/design/tokens')).toEqual({
			group: 'design',
			page: 'design-system',
			child: null
		});
		expect(locate('/design/componentsx')).toEqual({
			group: 'design',
			page: 'design-system',
			child: null
		});
	});

	it('puts Ops between Design and Documentation, with Klusters its first page (S-0580)', () => {
		const keys = SITE_MENU.map((g) => g.key);
		expect(keys.indexOf('ops')).toBe(keys.indexOf('design') + 1);
		expect(keys.indexOf('documentation')).toBe(keys.indexOf('ops') + 1);
		const ops = SITE_MENU.find((g) => g.key === 'ops')!;
		expect(ops.icon).toBe(ServerCogIcon);
		expect(ops.pages[0]).toMatchObject({
			key: 'klusters',
			label: 'Klusters',
			icon: ShipWheelIcon,
			path: '/ops/klusters',
			href: '/ops/klusters'
		});
	});

	it('places /ops/klusters on Klusters, and /ops or an unknown /ops page in Ops only (S-0580)', () => {
		expect(locate('/ops/klusters')).toEqual({ group: 'ops', page: 'klusters', child: null });
		expect(locate('/ops/klusters/k3d')).toEqual({ group: 'ops', page: 'klusters', child: null });
		expect(locate('/ops')).toEqual({ group: 'ops', page: null, child: null });
		expect(locate('/ops/unknown')).toEqual({ group: 'ops', page: null, child: null });
	});

	it('gives Ops Deployments after Klusters, at /ops/deployments (S-0581)', () => {
		const ops = SITE_MENU.find((g) => g.key === 'ops')!;
		expect(ops.pages.map((p) => p.key)).toEqual(['klusters', 'deployments', 'images']);
		expect(ops.pages[1]).toMatchObject({
			key: 'deployments',
			label: 'Deployments',
			icon: RocketIcon,
			path: '/ops/deployments',
			href: '/ops/deployments'
		});
		expect(locate('/ops/deployments')).toEqual({ group: 'ops', page: 'deployments', child: null });
		expect(locate('/ops/deployments/s-0581')).toEqual({
			group: 'ops',
			page: 'deployments',
			child: null
		});
	});

	it('gives Ops Agent images after Deployments, at /ops/images (S-0757)', () => {
		const ops = SITE_MENU.find((g) => g.key === 'ops')!;
		expect(ops.pages[2]).toMatchObject({
			key: 'images',
			label: 'Agent images',
			icon: ContainerIcon,
			path: '/ops/images',
			href: '/ops/images'
		});
		expect(locate('/ops/images')).toEqual({ group: 'ops', page: 'images', child: null });
	});

	it('places each Configuration page at its own path (S-0459)', () => {
		for (const page of ['features', 'agents', 'file-system', 'projects', 'harnesses'])
			expect(locate(`/settings/${page}`), page).toEqual({
				group: 'configuration',
				page,
				child: null
			});
	});

	it('places the Host page, which is not a settings page, in Configuration (S-0489)', () => {
		expect(locate('/host')).toEqual({ group: 'configuration', page: 'host', child: null });
		expect(locate('/hosts')).toEqual({ group: null, page: null, child: null });
	});

	it('places /settings, which goes on to a page, and an unknown page in Configuration only', () => {
		expect(locate('/settings')).toEqual({ group: 'configuration', page: null, child: null });
		expect(locate('/settings/tokens')).toEqual({ group: 'configuration', page: null, child: null });
	});

	it('places a path only its group claims in the group, with no page', () => {
		expect(locate('/items/S-0458')).toEqual({ group: 'flow', page: null, child: null });
		expect(locate('/review')).toEqual({ group: 'flow', page: null, child: null });
		expect(locate('/new')).toEqual({ group: 'flow', page: null, child: null });
		expect(locate('/edit/design/system/overview.md')).toEqual({
			group: 'documentation',
			page: null,
			child: null
		});
		expect(locate('/docs')).toEqual({ group: 'documentation', page: null, child: null });
		expect(locate('/docs/docs/users/flai.md')).toEqual({
			group: 'documentation',
			page: null,
			child: null
		});
	});

	it('places nothing for a page the rail does not list, nor one that only shares a prefix', () => {
		expect(locate('/')).toEqual({ group: null, page: null, child: null });
		expect(locate('/license')).toEqual({ group: null, page: null, child: null });
		expect(locate('/login')).toEqual({ group: null, page: null, child: null });
		expect(locate('/boards')).toEqual({ group: null, page: null, child: null });
		expect(locate('/designs')).toEqual({ group: null, page: null, child: null });
		expect(locate('/settingsx')).toEqual({ group: null, page: null, child: null });
		expect(locate('/opsx')).toEqual({ group: null, page: null, child: null });
	});

	it('gives Flow the Inbox’s indicator, and no other group one', () => {
		expect(SITE_MENU.map((g) => groupBadges(g))).toEqual([['inbox'], [], [], [], [], []]);
	});

	it('keeps the project in a page’s href', () => {
		expect(withProject('/board', null)).toBe('/board');
		expect(withProject('/board', 'harbour')).toBe('/board?project=harbour');
		expect(withProject('/settings/agents', 'a b')).toBe('/settings/agents?project=a%20b');
	});
});

// S-0420: the Agents group lists a page for each role agent.status reports, the orchestrator's, the
// planner's, and the analyzer's first, then the rest by name, each at /workflow/<name>.
describe('the Agents group with the roles agent.status reports (S-0420)', () => {
	/** The roles of a fixture agent.status, in its map's order: a defined role, reviewer, among them. */
	const ROLES = ['reviewer', 'curator', 'analyzer', 'orchestrator', 'planner', 'release-notes'];
	const agents = (roles: string[] | null) => siteMenu(roles).find((g) => g.key === 'agents')!;

	it('is the menu as it was while no roles are known', () => {
		expect(siteMenu(null)).toBe(SITE_MENU);
		expect(siteMenu(undefined)).toBe(SITE_MENU);
		expect(siteMenu([])).toBe(SITE_MENU);
	});

	it('lists the three first, then the curator and the defined roles by name', () => {
		expect(agents(ROLES).pages.map((p) => [p.key, p.label, p.href])).toEqual([
			['orchestration', 'Orchestration', '/workflow/orchestrator'],
			['planning', 'Planning', '/workflow/planner'],
			['analysis', 'Analysis', '/workflow/analyzer'],
			['role:curator', 'Curator', '/workflow/curator'],
			['role:release-notes', 'Release notes', '/workflow/release-notes'],
			['role:reviewer', 'Reviewer', '/workflow/reviewer'],
			['activity', 'Activity', '/activity'],
			['threads', 'Threads', '/threads'],
			['messages', 'Messages', '/messages']
		]);
		for (const p of agents(ROLES).pages) expect(p.icon, p.label).toBeTypeOf('function');
		const keys = siteMenu(ROLES).flatMap((g) => g.pages.map((p) => p.key));
		expect(new Set(keys).size).toBe(keys.length);
	});

	it('changes no other group', () => {
		const others = (menu: typeof SITE_MENU) => menu.filter((g) => g.key !== 'agents');
		expect(others(siteMenu(ROLES))).toEqual(others(SITE_MENU));
	});

	it('lists only the roles reported, and leaves out a name that is not a role’s', () => {
		expect(agents(['planner', 'reviewer', 'Bad Name', '../x']).pages.map((p) => p.label)).toEqual([
			'Planning',
			'Reviewer',
			'Activity',
			'Threads',
			'Messages'
		]);
	});

	it('places a role’s page in Agents, and one it does not list in the group with no page', () => {
		const menu = siteMenu(ROLES);
		expect(locate('/workflow/reviewer', menu)).toEqual({
			group: 'agents',
			page: 'role:reviewer',
			child: null
		});
		expect(locate('/workflow/planner', menu)).toEqual({
			group: 'agents',
			page: 'planning',
			child: null
		});
		expect(locate('/workflow/reviewer')).toEqual({ group: 'agents', page: null, child: null });
		expect(locate('/workflow/nobody', menu)).toEqual({ group: 'agents', page: null, child: null });
	});
});
