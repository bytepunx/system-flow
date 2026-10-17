// S-0442: which of the catalogue's components a story's branch changes, from stream.diff's paths,
// and the components ordered with those first.
import { describe, expect, it } from 'vitest';
import type { Component, ComponentsCatalogue } from '$lib/server/repo';
import {
	changedComponents,
	changedFirst,
	changesAComponent,
	componentFile,
	diffPaths
} from './changed';

const component = (project: string, path: string): Component => ({
	project,
	path,
	name: path.slice(path.lastIndexOf('/') + 1, -'.svelte'.length),
	props: []
});
const BUTTON = component('web', 'src/lib/Button.svelte');
const CARD = component('web', 'src/lib/Card.svelte');
const NAV = component('apps/admin', 'src/Nav.svelte');
const ROOT = component('.', 'src/App.svelte');
const CATALOGUE: ComponentsCatalogue = {
	projects: [
		{ name: 'web', path: 'web', kind: 'sveltekit', supported: true, components: [BUTTON, CARD] },
		{ name: 'admin', path: 'apps/admin', kind: 'svelte', supported: true, components: [NAV] },
		{ name: 'site', path: '.', kind: 'svelte', supported: true, components: [ROOT] },
		{ name: 'api', path: 'api', kind: 'go', supported: false, components: [] }
	]
};

describe('componentFile', () => {
	it("joins the sub-project's path and the component's", () => {
		expect(componentFile(BUTTON)).toBe('web/src/lib/Button.svelte');
		expect(componentFile(NAV)).toBe('apps/admin/src/Nav.svelte');
	});

	it('adds nothing for a sub-project at the root, and drops a leading ./ and a trailing /', () => {
		expect(componentFile(ROOT)).toBe('src/App.svelte');
		expect(componentFile({ project: '', path: 'src/App.svelte' })).toBe('src/App.svelte');
		expect(componentFile({ project: './web/', path: 'src/lib/Button.svelte' })).toBe(
			'web/src/lib/Button.svelte'
		);
	});
});

describe('diffPaths', () => {
	it("takes each file's path and a rename's old path", () => {
		const diff = {
			story: 'S-0442',
			files: [
				{ path: 'web/src/lib/Button.svelte', status: 'modified' },
				{ path: 'web/src/lib/Tile.svelte', old_path: 'web/src/lib/Card.svelte', status: 'renamed' }
			]
		};
		expect(diffPaths(diff)).toEqual([
			'web/src/lib/Button.svelte',
			'web/src/lib/Tile.svelte',
			'web/src/lib/Card.svelte'
		]);
	});

	it('answers [] for an answer without files', () => {
		expect(diffPaths(null)).toEqual([]);
		expect(diffPaths({})).toEqual([]);
		expect(diffPaths({ files: 'no' })).toEqual([]);
	});
});

describe('changedComponents', () => {
	it("answers the files of the catalogue's components the paths change, in any sub-project", () => {
		const changed = changedComponents(CATALOGUE, [
			'web/src/lib/Card.svelte',
			'apps/admin/src/Nav.svelte',
			'src/App.svelte',
			'web/src/lib/util.ts',
			'design/system/component-library.md'
		]);
		expect([...changed].sort()).toEqual([
			'apps/admin/src/Nav.svelte',
			'src/App.svelte',
			'web/src/lib/Card.svelte'
		]);
	});

	it('does not take a path relative to the sub-project, or a file the catalogue does not list', () => {
		expect(
			changedComponents(CATALOGUE, ['src/lib/Button.svelte', 'web/src/lib/New.svelte']).size
		).toBe(0);
	});
});

describe('changesAComponent', () => {
	it('says whether any path is a component of the catalogue', () => {
		expect(changesAComponent(CATALOGUE, ['README.md', 'web/src/lib/Button.svelte'])).toBe(true);
		expect(changesAComponent(CATALOGUE, ['README.md', 'web/package.json'])).toBe(false);
		expect(changesAComponent(CATALOGUE, [])).toBe(false);
		expect(changesAComponent({ projects: [] }, ['web/src/lib/Button.svelte'])).toBe(false);
	});
});

describe('changedFirst', () => {
	it('puts the changed components first, each part in its own order', () => {
		const all = [BUTTON, CARD, NAV, ROOT];
		const changed = new Set(['src/App.svelte', 'web/src/lib/Card.svelte']);
		expect(changedFirst(all, changed)).toEqual([CARD, ROOT, BUTTON, NAV]);
	});

	it('leaves the order as it is when none is changed', () => {
		expect(changedFirst([BUTTON, CARD], new Set())).toEqual([BUTTON, CARD]);
	});
});
