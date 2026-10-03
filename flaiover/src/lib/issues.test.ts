import { describe, expect, it } from 'vitest';
import { acceptLabel, issueDocPath, issueRows, type Issue } from './issues';

const issue = (id: string, over: Partial<Issue> = {}): Issue => ({
	id,
	title: `Title of ${id}`,
	class: 'efficiency',
	status: 'open',
	count: 1,
	path: `/repo/design/issues/${id}-slug.md`,
	stories: [],
	story: '',
	...over
});

describe('issueRows', () => {
	it('lists the story’s own issues first, checked, then the other open ones, unchecked', () => {
		const main = [issue('I-0002'), issue('I-0010'), issue('I-0003')];
		const own = [
			issue('I-0012', { stories: ['S-0198'] }),
			issue('I-0004', { stories: ['S-0198'] })
		];
		expect(issueRows('S-0198', main, own).map((r) => [r.id, r.recorded, r.checked])).toEqual([
			['I-0004', true, true],
			['I-0012', true, true],
			['I-0002', false, false],
			['I-0003', false, false],
			['I-0010', false, false]
		]);
	});

	it('gives an issue in both lists once, as the story’s own copy says it is', () => {
		const main = [issue('I-0005', { count: 2, title: 'old' })];
		const own = [issue('I-0005', { count: 3, title: 'bumped', stories: ['S-0198'] })];
		expect(issueRows('S-0198', main, own)).toEqual([
			{
				id: 'I-0005',
				title: 'bumped',
				class: 'efficiency',
				count: 3,
				doc: 'design/issues/I-0005-slug.md',
				recorded: true,
				checked: true
			}
		]);
	});

	it('counts an issue main says the story bumped as recorded', () => {
		const rows = issueRows('S-0198', [issue('I-0006', { stories: ['S-0100', 'S-0198'] })], []);
		expect(rows[0]).toMatchObject({ id: 'I-0006', recorded: true, checked: true });
	});

	it('leaves out closed issues and those an open story links, wherever they come from', () => {
		const main = [
			issue('I-0001', { status: 'closed' }),
			issue('I-0002', { story: 'S-0150' }),
			issue('I-0007')
		];
		const own = [
			issue('I-0008', { status: 'closed', stories: ['S-0198'] }),
			issue('I-0009', { story: 'S-0151', stories: ['S-0198'] }),
			// main has it open and unlinked; the story's copy, linked, wins
			issue('I-0007', { story: 'S-0152', stories: ['S-0198'] })
		];
		expect(issueRows('S-0198', main, own)).toEqual([]);
	});

	it('is empty when there are no issues', () => {
		expect(issueRows('S-0198', [], [])).toEqual([]);
	});
});

describe('acceptLabel', () => {
	it('says Accept with no issue checked, and that stories are made with one', () => {
		expect(acceptLabel(0, false)).toBe('Accept');
		expect(acceptLabel(1, false)).toBe('Accept and Create Stories');
		expect(acceptLabel(3, false)).toBe('Accept and Create Stories');
	});
	it('says it is accepting while it runs', () => {
		expect(acceptLabel(2, true)).toBe('Accepting…');
		expect(acceptLabel(0, true)).toBe('Accepting…');
	});
});

describe('issueDocPath', () => {
	it('is the issue’s file under design/issues, wherever flai read it', () => {
		expect(issueDocPath('/repo/.flai-cache/worktrees/S-0198/design/issues/I-0012-x.md')).toBe(
			'design/issues/I-0012-x.md'
		);
	});
});
