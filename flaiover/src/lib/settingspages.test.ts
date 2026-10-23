// S-0459: the five Configuration pages, the sections each shows, and where an old anchor lands.
import { describe, expect, it } from 'vitest';
import {
	SETTINGS_PAGES,
	SETTINGS_SECTIONS,
	pageOfAnchor,
	settingsPage,
	type SettingsPage
} from './settingspages';

/** How many pages hold each section, and the sections pages name that are none of the eleven. */
function homes(pages: readonly SettingsPage[]) {
	const count = Object.fromEntries(SETTINGS_SECTIONS.map((s) => [s, 0]));
	const unknown: string[] = [];
	for (const p of pages)
		for (const s of p.sections) {
			if (s in count) count[s]++;
			else unknown.push(s);
		}
	return { count, unknown };
}

describe('the Configuration pages (S-0459)', () => {
	it('are Features, Agents, File System, Projects, and Harnesses, each with its sections in order', () => {
		expect(SETTINGS_PAGES.map((p) => [p.slug, p.title, p.sections])).toEqual([
			['features', 'Features', ['actions', 'checks']],
			['agents', 'Agents', ['default-agent', 'providers', 'planning']],
			['file-system', 'File System', ['shared', 'import']],
			['projects', 'Projects', ['projects', 'sub-projects', 'tokens']],
			['harnesses', 'Harnesses', ['agent']]
		]);
	});

	it('puts each of the eleven sections on exactly one page', () => {
		expect(SETTINGS_SECTIONS).toHaveLength(11);
		expect(new Set(SETTINGS_SECTIONS).size).toBe(11);
		const { count, unknown } = homes(SETTINGS_PAGES);
		expect(unknown).toEqual([]);
		for (const s of SETTINGS_SECTIONS) expect(count[s], s).toBe(1);
	});

	it('counts a section on no page, or on two, so the check above fails for either', () => {
		const none = SETTINGS_PAGES.map((p) => ({
			...p,
			sections: p.sections.filter((s) => s !== 'tokens')
		}));
		expect(homes(none).count.tokens).toBe(0);
		const two = SETTINGS_PAGES.map((p) =>
			p.slug === 'features' ? { ...p, sections: [...p.sections, 'agent' as const] } : p
		);
		expect(homes(two).count.agent).toBe(2);
	});

	it('gives each page its own slug', () => {
		const slugs = SETTINGS_PAGES.map((p) => p.slug);
		expect(new Set(slugs).size).toBe(slugs.length);
	});

	it('finds a page by its slug, and none for an unknown one', () => {
		for (const p of SETTINGS_PAGES) expect(settingsPage(p.slug)).toBe(p);
		expect(settingsPage('tokens')).toBeUndefined();
		expect(settingsPage('')).toBeUndefined();
		expect(settingsPage('Features')).toBeUndefined();
	});

	it('lands an old /settings anchor on its page, and none or an unknown one on Features', () => {
		expect(pageOfAnchor('#features')).toBe('features');
		expect(pageOfAnchor('#agents')).toBe('agents');
		expect(pageOfAnchor('#file-system')).toBe('file-system');
		expect(pageOfAnchor('#projects')).toBe('projects');
		expect(pageOfAnchor('#harnesses')).toBe('harnesses');
		expect(pageOfAnchor('harnesses')).toBe('harnesses');
		expect(pageOfAnchor('')).toBe('features');
		expect(pageOfAnchor('#')).toBe('features');
		expect(pageOfAnchor('#tokens')).toBe('features');
		expect(pageOfAnchor('#sub-projects')).toBe('features');
	});
});
