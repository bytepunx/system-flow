// @vitest-environment jsdom
// S-0441: the Design page's unsaved edits kept in this browser per project, read back after a
// reload, a corrupt draft dropped, and another tab's change heard through the storage event.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { clearDraft, draftKey, loadDraft, onDraft, parseDraft, saveDraft } from './draft';
import { editKey } from './tokens';

const PRIMARY = { colorSpace: 'srgb', components: [0.2, 0.4, 0.6], hex: '#336699' };
const RADIUS = { value: 0.5, unit: 'rem' };
const edits = {
	[editKey('light', 'color.primary')]: PRIMARY,
	[editKey('dark', 'radius.md')]: RADIUS
};

// What another tab's write looks like to this one.
function fromAnotherTab(key: string | null, newValue: string | null) {
	window.dispatchEvent(new StorageEvent('storage', { key, newValue, storageArea: localStorage }));
}

describe('the design draft', () => {
	afterEach(() => localStorage.clear());

	it('reads back the edits it kept, and clears them', () => {
		expect(loadDraft('alpha')).toEqual({});
		saveDraft('alpha', edits);
		expect(loadDraft('alpha')).toEqual(edits);
		clearDraft('alpha');
		expect(loadDraft('alpha')).toEqual({});
		expect(localStorage.getItem(draftKey('alpha'))).toBeNull();
	});

	it('keeps no draft for no edits', () => {
		saveDraft('alpha', edits);
		saveDraft('alpha', {});
		expect(localStorage.getItem(draftKey('alpha'))).toBeNull();
	});

	it('keeps each project apart, the single project under its own key', () => {
		const other = { [editKey('light', 'color.primary')]: { ...PRIMARY, hex: '#000000' } };
		saveDraft('alpha', edits);
		saveDraft('beta', other);
		saveDraft(null, { [editKey('light', 'radius.md')]: RADIUS });
		expect(loadDraft('alpha')).toEqual(edits);
		expect(loadDraft('beta')).toEqual(other);
		expect(loadDraft(null)).toEqual({ [editKey('light', 'radius.md')]: RADIUS });
		clearDraft('alpha');
		expect(loadDraft('beta')).toEqual(other);
	});

	it('drops a draft that does not parse or is not an object of edits', () => {
		for (const raw of ['{not json', '[1,2]', '"text"', 'null', '{"color.primary":1}']) {
			localStorage.setItem(draftKey('alpha'), raw);
			expect(loadDraft('alpha'), raw).toEqual({});
			expect(localStorage.getItem(draftKey('alpha')), raw).toBeNull();
		}
		expect(parseDraft(JSON.stringify(edits))).toEqual(edits);
	});

	it('keeps the edits for this page when storage refuses them', () => {
		const refuse = () => {
			throw new Error('QuotaExceededError');
		};
		vi.spyOn(Storage.prototype, 'setItem').mockImplementation(refuse);
		expect(() => saveDraft('alpha', edits)).not.toThrow();
		vi.restoreAllMocks();
		expect(loadDraft('alpha')).toEqual({});
	});

	it("hears another tab's change to this project's draft until unsubscribed", () => {
		const heard = vi.fn();
		const stop = onDraft('alpha', heard);
		fromAnotherTab(draftKey('alpha'), JSON.stringify(edits));
		fromAnotherTab(draftKey('beta'), JSON.stringify(edits));
		fromAnotherTab('flaiover-theme', 'dark');
		fromAnotherTab(draftKey('alpha'), null);
		fromAnotherTab(draftKey('alpha'), '{not json');
		fromAnotherTab(null, null);
		expect(heard.mock.calls).toEqual([[edits], [{}], [{}], [{}]]);
		stop();
		fromAnotherTab(draftKey('alpha'), JSON.stringify(edits));
		expect(heard).toHaveBeenCalledTimes(4);
	});
});
