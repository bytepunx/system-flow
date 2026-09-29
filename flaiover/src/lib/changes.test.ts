import { describe, expect, it } from 'vitest';
import { affects, kindOf, type Change } from './changes';

// S-0161: a changed path's kind, read against the manifest's layout.
describe('kindOf', () => {
	const layout = { design: 'design', docs: 'docs', wip: 'wip' };
	it.each([
		['system-flow.yaml', 'project'],
		['wip/kanban/stories/S-0001-a.md', 'item'],
		['wip/kanban/board.md', 'item'],
		['wip/archive/kanban/tasks/T-0001-a.md', 'item'],
		['wip/archive/agents/S-0001.md', 'item'],
		['wip/agents/S-0161.md', 'narrative'],
		['wip/agents/index.md', 'narrative'],
		['wip/threads/TH-0001-a.md', 'thread'],
		['design/adrs/0001-a.md', 'adr'],
		['design/system/overview.md', 'document'],
		['docs/users/flaiover.md', 'document'],
		['wip/README.md', 'document'],
		['scripts/build.sh', 'other'],
		['wipe/kanban/x.md', 'other']
	])('%s is %s', (path, kind) => {
		expect(kindOf(path, layout)).toBe(kind);
	});

	it("follows the project's own folders", () => {
		const own = { design: 'plan/', docs: './guides', wip: 'work' };
		expect(kindOf('work/kanban/stories/S-0001-a.md', own)).toBe('item');
		expect(kindOf('plan/adrs/0001-a.md', own)).toBe('adr');
		expect(kindOf('guides/a.md', own)).toBe('document');
		expect(kindOf('wip/kanban/stories/S-0001-a.md', own)).toBe('other');
	});

	it('reads every path but the manifest as other while the layout is not known', () => {
		expect(kindOf('wip/agents/S-0161.md', null)).toBe('other');
		expect(kindOf('system-flow.yaml', null)).toBe('project');
	});
});

describe('affects', () => {
	const c = (kind: Change['kind']): Change => ({ path: 'x', kind });
	it('is true for a change of a kind asked about, and for the manifest or an unknown file', () => {
		expect(affects([c('narrative'), c('item')], 'item')).toBe(true);
		expect(affects([c('narrative'), c('thread')], 'item')).toBe(false);
		expect(affects([c('project')], 'item')).toBe(true);
		expect(affects([c('other')])).toBe(true);
		expect(affects([], 'item')).toBe(false);
	});
});
