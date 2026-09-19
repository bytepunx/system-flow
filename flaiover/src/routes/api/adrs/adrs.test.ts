import { describe, expect, it } from 'vitest';
import { _acceptArgs, _adrArgs } from './+server';

const TRAILER = '--trailer=Co-Authored-By: flaiover <flaiover@localhost>';
const body = '## Context\nc\n';

describe('arguments for recording an ADR', () => {
	it('builds one flai call with the body on stdin and the title last', () => {
		expect(
			_adrArgs({
				title: 'Tokens: one per project',
				status: 'accepted',
				supersedes: ['ADR-0007'],
				refines: [16, '0018'],
				body
			})
		).toEqual([
			'adr',
			'new',
			'--status=accepted',
			'--supersedes=7',
			'--refines=16',
			'--refines=18',
			'--body-stdin',
			'--autocommit',
			TRAILER,
			'--',
			'Tokens: one per project'
		]);
		expect(_adrArgs({ title: 'A proposal', body }).slice(0, 3)).toEqual([
			'adr',
			'new',
			'--status=proposed'
		]);
	});
	it('requires a title, a body, and a status an author may give', () => {
		expect(() => _adrArgs({ title: ' ', body })).toThrow(/title is required/);
		expect(() => _adrArgs({ title: 'X', body: '\n' })).toThrow(/body is empty/);
		expect(() => _adrArgs({ title: 'X', status: 'superseded', body })).toThrow(
			/proposed or accepted/
		);
	});
	it('lets nothing the designer types be read as a flag', () => {
		const args = _adrArgs({ title: '--autocommit  --json', body });
		expect(args.at(-2)).toBe('--');
		expect(args.at(-1)).toBe('--autocommit --json');
		expect(() => _adrArgs({ title: 'X', supersedes: ['--yes'], body })).toThrow(/not an ADR/);
		expect(() => _adrArgs({ title: 'X', refines: ['ADR-0000'], body })).toThrow(/not an ADR/);
		expect(() => _acceptArgs('../../x')).toThrow(/not an ADR/);
	});
	it('accepts by number', () => {
		expect(_acceptArgs('ADR-0027')).toEqual(['adr', 'accept', '27', '--autocommit', TRAILER]);
	});
});
