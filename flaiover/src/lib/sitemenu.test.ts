// S-0458: the nav rail's groups and pages, and which group and page a path belongs to.
import { describe, expect, it, vi } from 'vitest';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

import ActivityIcon from '@lucide/svelte/icons/activity';
import InboxIcon from '@lucide/svelte/icons/inbox';
import { SITE_MENU, groupBadges, locate, withProject } from './sitemenu';

describe('the nav rail (S-0458)', () => {
	it('has Flow, Agents, Design, Documentation, and Configuration, each with its pages', () => {
		expect(SITE_MENU.map((g) => [g.label, g.pages.map((p) => p.label)])).toEqual([
			['Flow', ['Board', 'Inbox', 'Charts', 'Workflow']],
			['Agents', ['Orchestration', 'Planning', 'Analysis', 'Activity', 'Threads', 'Messages']],
			['Design', ['System', 'Components']],
			['Documentation', ['System', 'Conventions', 'ADRs', 'Issues', 'Search']],
			['Configuration', ['Features', 'Agents', 'File System', 'Projects', 'Harnesses']]
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
				'/settings/harnesses'
			]
		]);
		const workflow = SITE_MENU[0].pages[3];
		expect(workflow.soon).toBe(true);
		expect(workflow.path).toBeUndefined();
		expect(SITE_MENU.flatMap((g) => g.pages).filter((p) => p.soon)).toEqual([workflow]);
	});

	it('places each page and the paths below it in its group', () => {
		expect(locate('/board')).toEqual({ group: 'flow', page: 'board' });
		expect(locate('/inbox')).toEqual({ group: 'flow', page: 'inbox' });
		expect(locate('/charts')).toEqual({ group: 'flow', page: 'charts' });
		expect(locate('/charts/cycle-time')).toEqual({ group: 'flow', page: 'charts' });
		expect(locate('/workflow/orchestrator')).toEqual({ group: 'agents', page: 'orchestration' });
		expect(locate('/workflow/planner')).toEqual({ group: 'agents', page: 'planning' });
		expect(locate('/workflow/analyzer')).toEqual({ group: 'agents', page: 'analysis' });
		expect(locate('/activity')).toEqual({ group: 'agents', page: 'activity' });
		expect(locate('/threads')).toEqual({ group: 'agents', page: 'threads' });
		expect(locate('/messages')).toEqual({ group: 'agents', page: 'messages' });
		expect(locate('/design')).toEqual({ group: 'design', page: 'design-system' });
		expect(locate('/design/components')).toEqual({ group: 'design', page: 'components' });
		expect(locate('/docs/design/system')).toEqual({ group: 'documentation', page: 'docs-system' });
		expect(locate('/docs/design/system/overview.md')).toEqual({
			group: 'documentation',
			page: 'docs-system'
		});
		expect(locate('/docs/design/conventions/git.md')).toEqual({
			group: 'documentation',
			page: 'conventions'
		});
		expect(locate('/adrs')).toEqual({ group: 'documentation', page: 'adrs' });
		expect(locate('/docs/design/issues/summary.md')).toEqual({
			group: 'documentation',
			page: 'issues'
		});
		expect(locate('/docs/design/issues/I-0001-golangci-lint-version-mismatch.md')).toEqual({
			group: 'documentation',
			page: 'issues'
		});
		expect(locate('/search')).toEqual({ group: 'documentation', page: 'search' });
	});

	it('selects the longest page path that answers, Components below the design system', () => {
		expect(locate('/design/components/Button')).toEqual({ group: 'design', page: 'components' });
		expect(locate('/design/tokens')).toEqual({ group: 'design', page: 'design-system' });
		expect(locate('/design/componentsx')).toEqual({ group: 'design', page: 'design-system' });
	});

	it('places each Configuration page at its own path (S-0459)', () => {
		for (const page of ['features', 'agents', 'file-system', 'projects', 'harnesses'])
			expect(locate(`/settings/${page}`), page).toEqual({ group: 'configuration', page });
	});

	it('places /settings, which goes on to a page, and an unknown page in Configuration only', () => {
		expect(locate('/settings')).toEqual({ group: 'configuration', page: null });
		expect(locate('/settings/tokens')).toEqual({ group: 'configuration', page: null });
	});

	it('places a path only its group claims in the group, with no page', () => {
		expect(locate('/items/S-0458')).toEqual({ group: 'flow', page: null });
		expect(locate('/review')).toEqual({ group: 'flow', page: null });
		expect(locate('/new')).toEqual({ group: 'flow', page: null });
		expect(locate('/edit/design/system/overview.md')).toEqual({
			group: 'documentation',
			page: null
		});
		expect(locate('/docs')).toEqual({ group: 'documentation', page: null });
		expect(locate('/docs/docs/users/flai.md')).toEqual({ group: 'documentation', page: null });
	});

	it('places nothing for a page the rail does not list, nor one that only shares a prefix', () => {
		expect(locate('/')).toEqual({ group: null, page: null });
		expect(locate('/host')).toEqual({ group: null, page: null });
		expect(locate('/license')).toEqual({ group: null, page: null });
		expect(locate('/login')).toEqual({ group: null, page: null });
		expect(locate('/boards')).toEqual({ group: null, page: null });
		expect(locate('/designs')).toEqual({ group: null, page: null });
		expect(locate('/settingsx')).toEqual({ group: null, page: null });
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
