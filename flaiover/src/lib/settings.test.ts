// S-0595, ADR-0263: the kluster: block's components as the page edits them, a name and its build
// one argument a line, and as settings.manifest takes them, a list of {name, build}.
import { describe, expect, it } from 'vitest';
import { componentList, componentRows, sameComponents } from './settings';

describe('the components of the kluster: block (S-0595)', () => {
	it('reads the rows from the list the manifest has', () => {
		expect(
			componentRows([
				{ name: 'api', build: ['scripts/build.sh', 'api'] },
				{ name: 'web', build: [] }
			])
		).toEqual([
			{ name: 'api', build: 'scripts/build.sh\napi' },
			{ name: 'web', build: '' }
		]);
	});

	it('reads no rows from what is not a list of components', () => {
		expect(componentRows(undefined)).toEqual([]);
		expect(componentRows('api')).toEqual([]);
		expect(componentRows([1, null, { build: ['x'] }])).toEqual([]);
	});

	it('gives the list the rows say, names trimmed, builds as argument lists', () => {
		expect(
			componentList([
				{ name: ' api ', build: 'scripts/build.sh\n\napi\r\n' },
				{ name: '', build: '' },
				{ name: 'web', build: '' }
			])
		).toEqual([
			{ name: 'api', build: ['scripts/build.sh', 'api'] },
			{ name: 'web', build: [] }
		]);
	});

	it('tells two lists apart by name, build, and order', () => {
		const a = [
			{ name: 'api', build: ['b', 'api'] },
			{ name: 'web', build: ['b', 'web'] }
		];
		expect(
			sameComponents(
				a,
				a.map((c) => ({ ...c, build: [...c.build] }))
			)
		).toBe(true);
		expect(sameComponents(a, [a[1], a[0]])).toBe(false);
		expect(sameComponents(a, [a[0]])).toBe(false);
		expect(sameComponents(a, [a[0], { name: 'web', build: ['b'] }])).toBe(false);
		expect(sameComponents([], [])).toBe(true);
	});
});
