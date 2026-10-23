// S-0458: the nav rail's groups and pages, and which group and page a path belongs to.
import { describe, expect, it, vi } from 'vitest';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

import ActivityIcon from '@lucide/svelte/icons/activity';
import InboxIcon from '@lucide/svelte/icons/inbox';
import { KINDS, TITLES } from './viz/catalog';
import { SITE_MENU, groupBadges, locate, withProject } from './sitemenu';

describe('the nav rail (S-0458)', () => {
	it('has Flow, Agents, Design, Documentation, and Configuration, each with its pages', () => {
		expect(SITE_MENU.map((g) => [g.label, g.pages.map((p) => p.label)])).toEqual([
			['Flow', ['Board', 'Inbox', 'Charts', 'Workflow']],
			['Agents', ['Orchestration', 'Planning', 'Analysis', 'Activity', 'Threads', 'Messages']],
			['Design', ['System', 'Components']],
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
	});

	it('gives Flow the Inbox’s indicator, and no other group one', () => {
		expect(SITE_MENU.map((g) => groupBadges(g))).toEqual([['inbox'], [], [], [], []]);
	});

	it('keeps the project in a page’s href', () => {
		expect(withProject('/board', null)).toBe('/board');
		expect(withProject('/board', 'harbour')).toBe('/board?project=harbour');
		expect(withProject('/settings/agents', 'a b')).toBe('/settings/agents?project=a%20b');
	});
});
