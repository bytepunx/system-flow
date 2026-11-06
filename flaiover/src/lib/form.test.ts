import { describe, expect, it } from 'vitest';
import { rebase } from './form';

// S-0574: a form read again keeps what the user changed and takes the rest from the new read.
describe('rebase', () => {
	it('keeps a field the user changed and takes an untouched one from the new read', () => {
		const seed = { model: 'opus', effort: 'high' };
		const current = { model: 'sonnet', effort: 'high' };
		const fresh = { model: 'opus', effort: 'low' };
		expect(rebase(seed, current, fresh)).toEqual({ model: 'sonnet', effort: 'low' });
	});

	it('takes the new read whole when nothing was changed', () => {
		const seed = { a: '1', b: { c: '2' } };
		const fresh = { a: '3', b: { c: '4' } };
		expect(rebase(seed, structuredClone(seed), fresh)).toEqual(fresh);
	});

	it('walks nested objects field by field', () => {
		const seed = { agents: { planning: { harness: '', model: 'opus' } } };
		const current = { agents: { planning: { harness: 'claude-code', model: 'opus' } } };
		const fresh = { agents: { planning: { harness: '', model: 'sonnet' } } };
		expect(rebase(seed, current, fresh)).toEqual({
			agents: { planning: { harness: 'claude-code', model: 'sonnet' } }
		});
	});

	it('drops a key the new read lacks and adds one it brings', () => {
		const seed: Record<string, string> = { gone: 'x', kept: 'y' };
		const current: Record<string, string> = { gone: 'typed', kept: 'y' };
		const fresh: Record<string, string> = { kept: 'y', added: 'z' };
		expect(rebase(seed, current, fresh)).toEqual({ kept: 'y', added: 'z' });
	});

	it('keeps a field typed back to what the new read says', () => {
		expect(rebase({ v: 'a' }, { v: 'b' }, { v: 'b' })).toEqual({ v: 'b' });
	});

	it('treats an array and a boolean as one field each', () => {
		const seed = { args: ['a'], on: false };
		const current = { args: ['a', 'b'], on: false };
		const fresh = { args: ['c'], on: true };
		expect(rebase(seed, current, fresh)).toEqual({ args: ['a', 'b'], on: true });
	});

	it('shares no object with the read it takes from', () => {
		const seed: Record<string, { args: string }> = {};
		const fresh: Record<string, { args: string }> = { claude: { args: 'a' } };
		const out = rebase(seed, {}, fresh);
		out.claude.args = 'typed';
		expect(fresh.claude.args).toBe('a');
	});

	it('rebases a single value', () => {
		expect(rebase('a', 'a', 'b')).toBe('b');
		expect(rebase('a', 'typed', 'b')).toBe('typed');
	});
});
