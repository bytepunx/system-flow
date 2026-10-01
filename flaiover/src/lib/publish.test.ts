import { describe, expect, it } from 'vitest';
import { doneLane, lagging } from './publish';

const behind = {
	remote: 'origin',
	behind: [{ component: 'cli', local: 'cli/v1.0.0', remote: 'cli/v1.4.0' }],
	fix: 'git fetch --tags origin',
	message: 'missing'
};
const unchecked = {
	remote: 'origin',
	unchecked: 'unable to access',
	fix: 'git fetch --tags origin',
	message: 'could not ask'
};

describe('publish (S-0174)', () => {
	it('lags only when the remote names a newer tag', () => {
		expect(lagging(behind)).toBe(true);
		expect(lagging(unchecked)).toBe(false);
		expect(lagging(null)).toBe(false);
	});

	it('leaves archived done cards out only while lagging', () => {
		const cards = [
			{ id: 'S-0001', archived: true },
			{ id: 'S-0002', archived: false },
			{ id: 'S-0003' }
		];
		expect(doneLane(cards, behind).map((c) => c.id)).toEqual(['S-0002', 'S-0003']);
		expect(doneLane(cards, unchecked)).toEqual(cards);
		expect(doneLane(cards, null)).toEqual(cards);
	});
});
