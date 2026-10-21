// S-0462: the board's lanes at a phone's width and a desktop's. Below Tailwind's sm width the lanes
// stack, each collapsed to its heading, a button that opens or closes its cards: Ready, In progress,
// and Review are open until the user changes it, and the choice outlasts a reload. From sm the six
// lanes are one row, every one showing its cards, and when they do not fit the row scrolls sideways
// inside the board, never the page. The page's own request for the board is answered in the browser
// with one story in every lane: a fake flai would make a project known to the server for every test
// file that runs after it, and the pages layout.e2e.ts checks render with none.
import { expect, test, type Page } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';

const STATES = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'];
const OPEN = ['ready', 'in-progress', 'review'];

// what playwright.config.ts made for the run and gave the server
const dir = process.env.FLAIOVER_E2E_DIR ?? '';
const token = () => readFileSync(join(dir, 'token'), 'utf8').trim();

/** /api/board's answer: one story in each lane, so an open lane's cards have a size. */
const board = {
	wip_limits: { 'in-progress': 3, review: 2 },
	order: ['S-9002'],
	writable: true,
	columns: Object.fromEntries(
		STATES.map((state, n) => [
			state,
			[
				{
					id: `S-900${n + 1}`,
					type: 'story',
					title: `A story in ${state}`,
					nature: 'feature',
					status: state,
					blocked: false,
					age_seconds: 3600
				}
			]
		])
	)
};

test.beforeEach(async ({ context, baseURL }) => {
	await context.addCookies([{ name: 'flaiover_session', value: token(), url: baseURL! }]);
	// the board itself, not its archive or limit writes
	await context.route(/\/api\/board(\?|$)/, (route) => route.fulfill({ json: board }));
});

/** How far the page reaches past the window sideways: 0 when it never scrolls sideways. */
const overflow = (page: Page) =>
	page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);

const lane = (page: Page, state: string) => page.locator(`section[data-lane="${state}"]`);
const toggle = (page: Page, state: string) =>
	lane(page, state).locator('[data-testid="lane-toggle"]');
const list = (page: Page, state: string) => page.locator(`#lane-cards-${state}`);

async function arrive(page: Page) {
	await page.goto('/board');
	await expect(page.locator('[data-testid="lanes"]')).toBeVisible();
	for (const state of STATES) await expect(lane(page, state)).toBeVisible();
}

/** Each lane's heading reads as open or closed, and its cards are shown or hidden to match. */
async function expectOpen(page: Page, open: string[]) {
	for (const state of STATES) {
		const isOpen = open.includes(state);
		await expect(toggle(page, state), state).toHaveAttribute('aria-expanded', String(isOpen));
		await expect(toggle(page, state)).toHaveAttribute('aria-controls', `lane-cards-${state}`);
		await expect(list(page, state), state).toHaveCSS('display', isOpen ? 'block' : 'none');
		if (isOpen) await expect(list(page, state).locator('[data-card]')).toBeVisible();
		else await expect(list(page, state).locator('[data-card]')).toBeHidden();
	}
}

/** Each lane's box, in the order of the lanes. */
const boxes = (page: Page) =>
	Promise.all(STATES.map(async (state) => (await lane(page, state).boundingBox())!));

/** The six lanes share one row: one top, and each to the right of the one before. */
async function expectOneRow(page: Page) {
	const all = await boxes(page);
	for (let n = 1; n < all.length; n++) {
		expect(Math.abs(all[n].y - all[0].y), `${STATES[n]}'s top`).toBeLessThanOrEqual(1);
		expect(all[n].x, `${STATES[n]} after ${STATES[n - 1]}`).toBeGreaterThan(
			all[n - 1].x + all[n - 1].width - 1
		);
	}
}

test.describe('at 390px', () => {
	test.use({ viewport: { width: 390, height: 844 } });

	test('the lanes stack, Ready, In progress, and Review open, and a heading toggles its lane', async ({
		page
	}) => {
		await arrive(page);
		await expectOpen(page, OPEN);
		// stacked: each lane below the one before
		const all = await boxes(page);
		for (let n = 1; n < all.length; n++)
			expect(all[n].y).toBeGreaterThanOrEqual(all[n - 1].y + all[n - 1].height);
		expect(await overflow(page)).toBeLessThanOrEqual(0);

		// tapping a closed lane's heading opens it, an open one's closes it
		await toggle(page, 'backlog').click();
		await expectOpen(page, ['backlog', ...OPEN]);
		await toggle(page, 'ready').click();
		await expectOpen(page, ['backlog', 'in-progress', 'review']);
		expect(await overflow(page)).toBeLessThanOrEqual(0);

		// the choice outlasts a reload
		await page.reload();
		await expect(page.locator('[data-testid="lanes"]')).toBeVisible();
		await expectOpen(page, ['backlog', 'in-progress', 'review']);
		expect(await overflow(page)).toBeLessThanOrEqual(0);

		// and toggles back
		await toggle(page, 'backlog').click();
		await toggle(page, 'ready').click();
		await expectOpen(page, OPEN);
	});
});

test.describe('at 1440px', () => {
	test.use({ viewport: { width: 1440, height: 900 } });

	test('the six lanes are one row, every one open, and the row scrolls sideways when narrowed', async ({
		page
	}) => {
		await arrive(page);
		await expectOneRow(page);
		await expectOpen(page, STATES);
		expect(await overflow(page)).toBeLessThanOrEqual(0);

		// from sm a heading toggles nothing
		await toggle(page, 'ready').click();
		await toggle(page, 'backlog').click();
		await expectOpen(page, STATES);

		// just above sm the six lanes do not fit: the row scrolls sideways inside the board, the
		// lanes still one row and the page no wider than the window
		await page.setViewportSize({ width: 700, height: 900 });
		const lanes = page.locator('[data-testid="lanes"]');
		await expect
			.poll(() => lanes.evaluate((el) => el.scrollWidth - el.clientWidth))
			.toBeGreaterThan(0);
		await expectOneRow(page);
		await expectOpen(page, STATES);
		expect(await overflow(page)).toBeLessThanOrEqual(0);
	});
});
