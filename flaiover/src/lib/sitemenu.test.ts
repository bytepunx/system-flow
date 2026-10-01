// S-0172: the site menu's groups, and which group and page a path belongs to.
import { describe, expect, it } from 'vitest';
import { SITE_MENU, groupBadges, locate } from './sitemenu';

describe('the site menu (S-0172)', () => {
	it('has Workflow, Status, and Host, with their pages in order', () => {
		expect(SITE_MENU.map((g) => [g.label, g.pages.map((p) => p.label)])).toEqual([
			['Workflow', ['Overview', 'Board', 'Inbox', 'Threads', 'Activity']],
			['Status', ['Charts', 'ADRs', 'Docs', 'Search']],
			['Host', ['Updates', 'Settings']]
		]);
	});

	it('places a page and the paths below it in its group', () => {
		expect(locate('/')).toEqual({ group: 'workflow', page: 'overview' });
		expect(locate('/inbox')).toEqual({ group: 'workflow', page: 'inbox' });
		expect(locate('/threads')).toEqual({ group: 'workflow', page: 'threads' });
		expect(locate('/charts/cycle-time')).toEqual({ group: 'status', page: 'charts' });
		expect(locate('/docs/design/system/overview.md')).toEqual({ group: 'status', page: 'docs' });
		expect(locate('/host')).toEqual({ group: 'host', page: 'updates' });
		expect(locate('/settings')).toEqual({ group: 'host', page: 'settings' });
	});

	it('places a path only its group claims in the group, with no page', () => {
		expect(locate('/items/S-0172')).toEqual({ group: 'workflow', page: null });
		expect(locate('/edit/design/system/overview.md')).toEqual({ group: 'status', page: null });
	});

	it('places nothing for a path no group has, nor for one that only shares a prefix', () => {
		expect(locate('/login')).toEqual({ group: null, page: null });
		expect(locate('/boards')).toEqual({ group: null, page: null });
	});

	it('gives a group the indicators of its pages', () => {
		expect(SITE_MENU.map((g) => groupBadges(g))).toEqual([['inbox'], [], []]);
	});
});
