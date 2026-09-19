import { describe, expect, it } from 'vitest';
import { _newArgs } from './+server';

const story = {
	type: 'story',
	title: 'Created from the board',
	parent: 'E-0006',
	body: '## Goal\nx\n'
};
const TRAILER = '--trailer=Co-Authored-By: flaiover <flaiover@localhost>';

describe('arguments for creating an item', () => {
	it('builds one flai call, with the designer as owner, the body on stdin, and the title last', () => {
		expect(
			_newArgs(
				{ ...story, nature: 'improvement', tags: ['dashboard'], touches: ['flaiover/src'] },
				'alex'
			)
		).toEqual([
			'story',
			'new',
			'--nature=improvement',
			'--owner=alex',
			'--epic=E-0006',
			'--tag=dashboard',
			'--touches=flaiover/src',
			'--body-stdin',
			'--autocommit',
			TRAILER,
			'--',
			'Created from the board'
		]);
		expect(_newArgs({ type: 'epic', title: 'An epic', body: '## Outcome\nx\n' }, 'alex')).toEqual([
			'epic',
			'new',
			'--nature=feature',
			'--owner=alex',
			'--body-stdin',
			'--autocommit',
			TRAILER,
			'--',
			'An epic'
		]);
	});
	it('creates epics and stories only', () => {
		expect(() => _newArgs({ ...story, type: 'task' }, 'alex')).toThrow(/epic or story/);
		expect(() => _newArgs({ ...story, type: undefined }, 'alex')).toThrow(/epic or story/);
	});
	it('requires a title, a body, a known nature, and a story’s epic', () => {
		expect(() => _newArgs({ ...story, title: '   ' }, 'alex')).toThrow(/title is required/);
		expect(() => _newArgs({ ...story, body: ' \n' }, 'alex')).toThrow(/body is empty/);
		expect(() => _newArgs({ ...story, nature: 'chore' }, 'alex')).toThrow(/nature must be one of/);
		expect(() => _newArgs({ ...story, parent: undefined }, 'alex')).toThrow(/parent epic/);
		expect(() => _newArgs({ ...story, parent: 'S-0001' }, 'alex')).toThrow(/parent epic/);
		expect(() =>
			_newArgs({ type: 'epic', title: 'E', body: 'x', parent: 'E-0001' }, 'alex')
		).toThrow(/no parent/);
	});
	it('lets nothing the designer types be read as a flag', () => {
		const args = _newArgs({ ...story, title: '--json  and\nmore', tags: ['--yes'] }, 'alex');
		expect(args.at(-2)).toBe('--');
		expect(args.at(-1)).toBe('--json and more');
		expect(args).toContain('--tag=--yes');
		expect(() => _newArgs({ ...story, parent: '--epic' }, 'alex')).toThrow(/parent epic/);
		expect(() => _newArgs({ ...story, tags: ['a,b'] }, 'alex')).toThrow(/without commas/);
		expect(() => _newArgs({ ...story, touches: [''] }, 'alex')).toThrow(/without commas/);
	});
});
